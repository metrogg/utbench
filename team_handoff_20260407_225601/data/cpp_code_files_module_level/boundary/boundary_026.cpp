#ifndef TENSORFLOW_CORE_COMMON_RUNTIME_LOWER_CASE_OP_H_
#define TENSORFLOW_CORE_COMMON_RUNTIME_LOWER_CASE_OP_H_
#include "tensorflow/core/lib/core/status.h"
namespace tensorflow {
class Graph;
class Node;
Status RewriteCaseNode(Node* n, Graph* g, bool keep_node_fetchable);
}  
#endif  
#include "tensorflow/core/common_runtime/lower_case_op.h"
#include "tensorflow/core/common_runtime/inline_function_utils.h"
#include "tensorflow/core/framework/node_def_builder.h"
#include "tensorflow/core/graph/graph.h"
#include "tensorflow/core/graph/node_builder.h"
#include "tensorflow/core/lib/core/errors.h"
namespace tensorflow {
namespace {
using NodeOut = NodeBuilder::NodeOut;
constexpr const char* const kLowerAsMultiDeviceFunctionAttr =
    LowerFunctionalOpsConstants::kLowerAsMultiDeviceFunctionAttr;
class CaseBuilder {
 public:
  CaseBuilder(Node* case_op, const std::vector<string>& branch_fn_names,
              bool keep_node_fetchable, Graph* graph);
  Status CreatePivotNodes();
  Status AddInputs();
  Status AddOutputs();
  Status BuildLoweredCaseOutput();
 private:
  string NewName(const string& infix);
  Status AddInput(Node* src, int src_output);
  std::vector<NodeOut> outputs_;
  Node* control_predecessor_;
  Node* case_op_;
  Node* lowered_case_output_;
  OutputTensor branch_index_;
  int num_branches_;
  std::vector<Node*> pivots_;
  std::vector<Node*> call_nodes_;
  Node* branch_executed_node_;
  Graph* graph_;
  string name_;
  bool keep_node_fetchable_;
  NodeDebugInfo debug_info_;
  std::vector<NodeBuilder> branch_call_builders_;
};
CaseBuilder::CaseBuilder(Node* case_op,
                         const std::vector<string>& branch_fn_names,
                         bool keep_node_fetchable, Graph* graph)
    : case_op_(case_op),
      num_branches_(branch_fn_names.size()),
      graph_(graph),
      name_(case_op->name()),
      keep_node_fetchable_(keep_node_fetchable),
      debug_info_(*case_op_) {
  branch_call_builders_.reserve(num_branches_);
  for (int b = 0; b < num_branches_; b++) {
    branch_call_builders_.emplace_back(NewName(strings::StrCat("branch", b)),
                                       branch_fn_names[b], graph->op_registry(),
                                       &debug_info_);
    branch_call_builders_[b].Device(case_op_->requested_device());
    branch_call_builders_[b].Attr(kLowerAsMultiDeviceFunctionAttr, true);
  }
  TF_CHECK_OK(case_op_->input_tensor(0, &branch_index_));
}
Status CaseBuilder::CreatePivotNodes() {
  Node* branch_index;
  TF_RETURN_IF_ERROR(NodeBuilder(NewName("branch_index"), "_SwitchN",
                                 graph_->op_registry(), &debug_info_)
                         .Input(NodeOut(branch_index_))
                         .Input(NodeOut(branch_index_))
                         .Attr("num_outs", num_branches_)
                         .Device(case_op_->requested_device())
                         .Finalize(graph_, &branch_index));
  control_predecessor_ = branch_index;
  pivots_.resize(num_branches_, nullptr);
  for (int b = 0; b < num_branches_; b++) {
    TF_RETURN_IF_ERROR(NodeBuilder(NewName(strings::StrCat("pivot_", b)),
                                   "Identity", graph_->op_registry(),
                                   &debug_info_)
                           .Input(branch_index, b)
                           .Device(case_op_->requested_device())
                           .Finalize(graph_, &pivots_[b]));
  }
  return absl::OkStatus();
}
string CaseBuilder::NewName(const string& infix) {
  return graph_->NewName(strings::StrCat(name_, "/", infix));
}
Status CaseBuilder::AddInput(Node* src, int src_output) {
  Node* input;
  NodeDebugInfo debug_info(*src);
  TF_RETURN_IF_ERROR(NodeBuilder(NewName(src->name()), "_SwitchN",
                                 graph_->op_registry(), &debug_info)
                         .Input(src, src_output)
                         .Input(branch_index_)
                         .Device(src->requested_device())
                         .Attr("_class", {src->name()})
                         .Attr("num_outs", num_branches_)
                         .Finalize(graph_, &input));
  for (int b = 0; b < num_branches_; b++) {
    branch_call_builders_[b].Input(input, b);
  }
  return absl::OkStatus();
}
Status CaseBuilder::AddInputs() {
  std::vector<const Edge*> edges;
  TF_RETURN_IF_ERROR(case_op_->input_edges(&edges));
  for (int i = 1; i < edges.size(); ++i) {
    const Edge* e = edges[i];
    TF_RETURN_IF_ERROR(AddInput(e->src(), e->src_output()));
  }
  for (const Edge* e : case_op_->in_edges()) {
    if (e->IsControlEdge()) {
      graph_->AddControlEdge(e->src(), control_predecessor_);
    }
  }
  return absl::OkStatus();
}
Status CaseBuilder::AddOutputs() {
  call_nodes_.resize(num_branches_, nullptr);
  for (int b = 0; b < num_branches_; b++) {
    TF_RETURN_IF_ERROR(
        branch_call_builders_[b].Finalize(graph_, &call_nodes_[b]));
    graph_->AddControlEdge(pivots_[b], call_nodes_[b]);
  }
  const int num_outputs = call_nodes_[0]->num_outputs();
  std::vector<Node*> merges(num_outputs);
  outputs_.resize(merges.size());
  for (int i = 0; i < num_outputs; ++i) {
    std::vector<NodeOut> merge_input;
    merge_input.reserve(num_branches_);
    for (int j = 0; j < num_branches_; j++) {
      merge_input.emplace_back(call_nodes_[j], i);
    }
    TF_RETURN_IF_ERROR(NodeBuilder(NewName("merge"), "Merge",
                                   graph_->op_registry(), &debug_info_)
                           .Input(merge_input)
                           .Device(case_op_->requested_device())
                           .Finalize(graph_, &merges[i]));
    outputs_[i] = NodeOut(merges[i], 0);
  }
  std::vector<NodeOut> pivots(num_branches_);
  for (int j = 0; j < num_branches_; j++) {
    pivots[j] = NodeOut(pivots_[j]);
  }
  TF_RETURN_IF_ERROR(NodeBuilder(NewName("branch_executed"), "Merge",
                                 graph_->op_registry(), &debug_info_)
                         .Input(pivots)
                         .ControlInputs(call_nodes_)
                         .Device(case_op_->requested_device())
                         .Finalize(graph_, &branch_executed_node_));
  TF_RETURN_IF_ERROR(BuildLoweredCaseOutput());
  for (const Edge* e : case_op_->out_edges()) {
    if (e->IsControlEdge()) {
      graph_->AddControlEdge(branch_executed_node_, e->dst());
    } else {
      graph_->AddEdge(merges[e->src_output()], 0, e->dst(), e->dst_input());
    }
  }
  return absl::OkStatus();
}
Status CaseBuilder::BuildLoweredCaseOutput() {
  NodeBuilder builder = keep_node_fetchable_ && !outputs_.empty()
                            ? NodeBuilder(name_, "IdentityN").Input(outputs_)
                            : NodeBuilder(name_, "NoOp");
  return builder.Device(case_op_->requested_device())
      .ControlInput(branch_executed_node_)
      .Finalize(graph_, &lowered_case_output_);
}
}  
Status RewriteCaseNode(Node* n, Graph* g, bool keep_node_fetchable) {
  VLOG(2) << "Lower Case node (keep_node_fetchable=" << keep_node_fetchable
          << "): " << SummarizeNode(*n);
  const AttrValue* branches_attr = n->attrs().Find("branches");
  if (branches_attr == nullptr) {
    return errors::InvalidArgument("branch functions missing");
  }
  int num_branches = branches_attr->list().func_size();
  std::vector<string> branch_fn_names;
  branch_fn_names.reserve(num_branches);
  for (int b = 0; b < num_branches; b++) {
    branch_fn_names.emplace_back(branches_attr->list().func(b).name());
  }
  CaseBuilder cb(n, branch_fn_names, keep_node_fetchable, g);
  TF_RETURN_IF_ERROR(cb.CreatePivotNodes());
  TF_RETURN_IF_ERROR(cb.AddInputs());
  TF_RETURN_IF_ERROR(cb.AddOutputs());
  g->RemoveNode(n);
  return absl::OkStatus();
}
}  