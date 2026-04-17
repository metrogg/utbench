#ifndef TENSORFLOW_CORE_COMMON_RUNTIME_LOWER_IF_OP_H_
#define TENSORFLOW_CORE_COMMON_RUNTIME_LOWER_IF_OP_H_
#include "tensorflow/core/lib/core/status.h"
namespace tensorflow {
class Graph;
class Node;
Status RewriteIfNode(Node* n, Graph* g, bool keep_node_fetchable);
}  
#endif  
#include "tensorflow/core/common_runtime/lower_if_op.h"
#include "tensorflow/core/common_runtime/inline_function_utils.h"
#include "tensorflow/core/framework/node_def_builder.h"
#include "tensorflow/core/graph/graph.h"
#include "tensorflow/core/graph/node_builder.h"
namespace tensorflow {
namespace {
using NodeOut = NodeBuilder::NodeOut;
constexpr const char* const kLowerAsMultiDeviceFunctionAttr =
    LowerFunctionalOpsConstants::kLowerAsMultiDeviceFunctionAttr;
class CondBuilder {
 public:
  enum Branch { kElseBranch = 0, kThenBranch = 1 };
  CondBuilder(Node* if_op, const NameAttrList& then_fn,
              const NameAttrList& else_fn, bool keep_node_fetchable,
              Graph* graph);
  Status CreatePivotNodes();
  Status AddInputs();
  Status AddOutputs();
  Status BuildLoweredIfOutput();
 private:
  string NewName(const string& infix);
  Status AddInput(Node* src, int src_output);
  Status SetColocationAndFinalize(NodeBuilder node_builder, Graph* graph,
                                  Node** created_node);
  std::vector<NodeOut> outputs_;
  Node* control_predecessor_;
  Node* if_op_;
  const AttrValue* coloc_attr_;
  Node* lowered_if_output_;
  OutputTensor pred_;
  Node* pivot_f_;
  Node* pivot_t_;
  Node* then_call_node_;
  Node* else_call_node_;
  Node* branch_executed_node_;
  Graph* graph_;
  string name_;
  bool keep_node_fetchable_;
  NodeDebugInfo debug_info_;
  NodeBuilder then_call_builder_;
  NodeBuilder else_call_builder_;
};
CondBuilder::CondBuilder(Node* if_op, const NameAttrList& then_fn,
                         const NameAttrList& else_fn, bool keep_node_fetchable,
                         Graph* graph)
    : if_op_(if_op),
      coloc_attr_(if_op_->attrs().Find(kColocationAttrName)),
      graph_(graph),
      name_(if_op->name()),
      keep_node_fetchable_(keep_node_fetchable),
      debug_info_(*if_op_),
      then_call_builder_(NewName("then"), then_fn.name(), graph->op_registry(),
                         &debug_info_),
      else_call_builder_(NewName("else"), else_fn.name(), graph->op_registry(),
                         &debug_info_) {
  TF_CHECK_OK(if_op_->input_tensor(0, &pred_));
  then_call_builder_.Device(if_op_->requested_device());
  then_call_builder_.Attr(kLowerAsMultiDeviceFunctionAttr, true);
  for (const auto& i : then_fn.attr()) {
    then_call_builder_.Attr(i.first, i.second);
  }
  else_call_builder_.Device(if_op_->requested_device());
  else_call_builder_.Attr(kLowerAsMultiDeviceFunctionAttr, true);
  for (const auto& i : else_fn.attr()) {
    else_call_builder_.Attr(i.first, i.second);
  }
}
Status CondBuilder::SetColocationAndFinalize(NodeBuilder node_builder,
                                             Graph* graph,
                                             Node** created_node) {
  if (coloc_attr_ != nullptr) {
    node_builder = node_builder.Attr(kColocationAttrName, *coloc_attr_);
  }
  return node_builder.Finalize(graph, created_node);
}
Status CondBuilder::CreatePivotNodes() {
  Node* switch_pred;
  TF_RETURN_IF_ERROR(
      SetColocationAndFinalize(NodeBuilder(NewName("switch_pred"), "Switch",
                                           graph_->op_registry(), &debug_info_)
                                   .Input(NodeOut(pred_))
                                   .Input(NodeOut(pred_))
                                   .Device(if_op_->requested_device()),
                               graph_, &switch_pred));
  control_predecessor_ = switch_pred;
  TF_RETURN_IF_ERROR(
      SetColocationAndFinalize(NodeBuilder(NewName("pivot_f"), "Identity",
                                           graph_->op_registry(), &debug_info_)
                                   .Input(switch_pred, kElseBranch)
                                   .Device(if_op_->requested_device()),
                               graph_, &pivot_f_));
  TF_RETURN_IF_ERROR(
      SetColocationAndFinalize(NodeBuilder(NewName("pivot_t"), "Identity",
                                           graph_->op_registry(), &debug_info_)
                                   .Input(switch_pred, kThenBranch)
                                   .Device(if_op_->requested_device()),
                               graph_, &pivot_t_));
  return absl::OkStatus();
}
string CondBuilder::NewName(const string& infix) {
  return graph_->NewName(strings::StrCat(name_, "/", infix));
}
Status CondBuilder::AddInput(Node* src, int src_output) {
  Node* input;
  NodeDebugInfo debug_info(*src);
  TF_RETURN_IF_ERROR(
      NodeBuilder(NewName(src->name()), "Switch", graph_->op_registry(),
                  &debug_info)
          .Input(src, src_output)
          .Input(pred_)
          .Device(src->requested_device())
          .Attr(kColocationAttrName,
                {absl::StrCat(kColocationGroupPrefix, src->name())})
          .Finalize(graph_, &input));
  then_call_builder_.Input(input, kThenBranch);
  else_call_builder_.Input(input, kElseBranch);
  return absl::OkStatus();
}
Status CondBuilder::AddInputs() {
  std::vector<const Edge*> edges;
  TF_RETURN_IF_ERROR(if_op_->input_edges(&edges));
  for (int i = 1; i < edges.size(); ++i) {
    const Edge* e = edges[i];
    TF_RETURN_IF_ERROR(AddInput(e->src(), e->src_output()));
  }
  for (const Edge* e : if_op_->in_edges()) {
    if (e->IsControlEdge()) {
      graph_->AddControlEdge(e->src(), control_predecessor_);
    }
  }
  return absl::OkStatus();
}
Status CondBuilder::AddOutputs() {
  TF_RETURN_IF_ERROR(then_call_builder_.Finalize(graph_, &then_call_node_));
  graph_->AddControlEdge(pivot_t_, then_call_node_);
  TF_RETURN_IF_ERROR(else_call_builder_.Finalize(graph_, &else_call_node_));
  graph_->AddControlEdge(pivot_f_, else_call_node_);
  std::vector<Node*> merges(then_call_node_->num_outputs());
  outputs_.resize(merges.size());
  for (int i = 0; i < then_call_node_->num_outputs(); ++i) {
    TF_RETURN_IF_ERROR(SetColocationAndFinalize(
        NodeBuilder(NewName("output"), "Merge", graph_->op_registry(),
                    &debug_info_)
            .Input({NodeOut(then_call_node_, i), NodeOut(else_call_node_, i)})
            .Device(if_op_->requested_device()),
        graph_, &merges[i]));
    outputs_[i] = NodeOut(merges[i], 0);
  }
  TF_RETURN_IF_ERROR(SetColocationAndFinalize(
      NodeBuilder(NewName("branch_executed"), "Merge", graph_->op_registry(),
                  &debug_info_)
          .Input({pivot_t_, pivot_f_})
          .ControlInputs({then_call_node_, else_call_node_})
          .Device(if_op_->requested_device()),
      graph_, &branch_executed_node_));
  TF_RETURN_IF_ERROR(BuildLoweredIfOutput());
  for (const Edge* e : if_op_->out_edges()) {
    if (e->IsControlEdge()) {
      graph_->AddControlEdge(branch_executed_node_, e->dst());
    } else {
      graph_->AddEdge(merges[e->src_output()], 0, e->dst(), e->dst_input());
    }
  }
  return absl::OkStatus();
}
Status CondBuilder::BuildLoweredIfOutput() {
  NodeBuilder builder = keep_node_fetchable_ && !outputs_.empty()
                            ? NodeBuilder(name_, "IdentityN").Input(outputs_)
                            : NodeBuilder(name_, "NoOp");
  return builder.Device(if_op_->requested_device())
      .ControlInput(branch_executed_node_)
      .Finalize(graph_, &lowered_if_output_);
}
}  
Status RewriteIfNode(Node* n, Graph* g, bool keep_node_fetchable) {
  VLOG(2) << "Lower If node (keep_node_fetchable=" << keep_node_fetchable
          << "): " << SummarizeNode(*n);
  const AttrValue* then_attr = n->attrs().Find("then_branch");
  if (then_attr == nullptr) {
    return errors::InvalidArgument("Then branch function missing");
  }
  const AttrValue* else_attr = n->attrs().Find("else_branch");
  if (else_attr == nullptr) {
    return errors::InvalidArgument("Else branch function missing");
  }
  CondBuilder cb(n, then_attr->func(), else_attr->func(), keep_node_fetchable,
                 g);
  TF_RETURN_IF_ERROR(cb.CreatePivotNodes());
  TF_RETURN_IF_ERROR(cb.AddInputs());
  TF_RETURN_IF_ERROR(cb.AddOutputs());
  g->RemoveNode(n);
  return absl::OkStatus();
}
}  