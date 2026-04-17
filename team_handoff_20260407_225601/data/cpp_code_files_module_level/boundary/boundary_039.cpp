#ifndef TENSORFLOW_COMPILER_TF2XLA_CONST_ANALYSIS_H_
#define TENSORFLOW_COMPILER_TF2XLA_CONST_ANALYSIS_H_
#include <vector>
#include "tensorflow/core/graph/graph.h"
#include "tensorflow/core/lib/core/status.h"
namespace tensorflow {
Status BackwardsConstAnalysis(
    const Graph& g, std::vector<bool>* compile_time_const_arg_indices,
    std::vector<bool>* compile_time_const_nodes,
    FunctionLibraryRuntime* flib_runtime,
    std::function<bool(const Edge&)> edge_filter_input = nullptr);
Status GetCompileTimeConstInputs(const OpKernel* op_kernel,
                                 std::vector<int>* const_input_idxs,
                                 FunctionLibraryRuntime* flib_runtime);
}  
#endif  
#include "tensorflow/compiler/tf2xla/const_analysis.h"
#include <unordered_map>
#include <unordered_set>
#include "absl/algorithm/container.h"
#include "tensorflow/compiler/tf2xla/tf2xla_util.h"
#include "tensorflow/compiler/tf2xla/xla_op_registry.h"
#include "xla/status_macros.h"
#include "tensorflow/core/common_runtime/function.h"
#include "tensorflow/core/framework/attr_value.pb.h"
#include "tensorflow/core/framework/function.h"
#include "tensorflow/core/framework/node_def_util.h"
#include "tensorflow/core/graph/algorithm.h"
#include "tensorflow/core/lib/core/errors.h"
namespace tensorflow {
namespace {
Status GetFunctionBody(FunctionLibraryRuntime* flib_runtime,
                       const NodeDef& node, StringPiece func_attr_name,
                       const FunctionBody** fbody) {
  NameAttrList name_attr_list;
  TF_RETURN_IF_ERROR(GetNodeAttr(node, func_attr_name, &name_attr_list));
  FunctionLibraryRuntime::Handle func_handle;
  TF_RETURN_IF_ERROR(flib_runtime->Instantiate(
      name_attr_list.name(), AttrSlice(&name_attr_list.attr()), &func_handle));
  *fbody = flib_runtime->GetFunctionBody(func_handle);
  return absl::OkStatus();
}
Status GetFunctionBodies(FunctionLibraryRuntime* flib_runtime,
                         const NodeDef& node, StringPiece func_list_attr_name,
                         std::vector<const FunctionBody*>* fbodies) {
  std::vector<NameAttrList> name_attr_lists;
  TF_RETURN_IF_ERROR(GetNodeAttr(node, func_list_attr_name, &name_attr_lists));
  for (const NameAttrList& name_attr_list : name_attr_lists) {
    FunctionLibraryRuntime::Handle func_handle;
    TF_RETURN_IF_ERROR(flib_runtime->Instantiate(
        name_attr_list.name(), AttrSlice(&name_attr_list.attr()),
        &func_handle));
    fbodies->push_back(flib_runtime->GetFunctionBody(func_handle));
  }
  return absl::OkStatus();
}
Status CondConstInputIndices(
    absl::Span<const FunctionBody* const> branch_bodies,
    std::vector<int>* const_input_idxs, FunctionLibraryRuntime* flib_runtime) {
  TF_RET_CHECK(!branch_bodies.empty());
  TF_RET_CHECK(branch_bodies[0] != nullptr);
  int num_inputs =
      branch_bodies[0]->record->fdef().signature().input_arg_size();
  std::vector<bool> compile_time_const_arg_indices(num_inputs);
  for (auto fbody : branch_bodies) {
    TF_RET_CHECK(fbody != nullptr);
    TF_RETURN_IF_ERROR(BackwardsConstAnalysis(
        *(fbody->graph), &compile_time_const_arg_indices,
        nullptr, flib_runtime));
  }
  for (int i = 0, end = compile_time_const_arg_indices.size(); i < end; i++) {
    if (compile_time_const_arg_indices[i]) {
      const_input_idxs->push_back(i + 1);
    }
  }
  return absl::OkStatus();
}
Status GetCompileTimeConstInputs(const NodeDef& node, const OpKernel* op_kernel,
                                 const OpDef* op_def,
                                 std::vector<int>* const_input_idxs,
                                 FunctionLibraryRuntime* flib_runtime) {
  DCHECK(op_def != nullptr || op_kernel != nullptr);
  if (node.op() == "While" || node.op() == "StatelessWhile") {
    const FunctionBody* fcond = nullptr;
    const FunctionBody* fbody = nullptr;
    TF_RETURN_IF_ERROR(GetFunctionBody(flib_runtime, node, "cond", &fcond));
    TF_RETURN_IF_ERROR(GetFunctionBody(flib_runtime, node, "body", &fbody));
    TF_RET_CHECK(fcond);
    TF_RET_CHECK(fbody);
    int num_inputs = fbody->record->fdef().signature().input_arg_size();
    std::vector<bool> compile_time_const_arg_indices(num_inputs);
    TF_RETURN_IF_ERROR(BackwardsConstAnalysis(
        *(fcond->graph), &compile_time_const_arg_indices,
        nullptr, flib_runtime));
    TF_RETURN_IF_ERROR(BackwardsConstAnalysis(
        *(fbody->graph), &compile_time_const_arg_indices,
        nullptr, flib_runtime));
    for (int i = 0; i < num_inputs; i++) {
      if (compile_time_const_arg_indices[i]) {
        TF_ASSIGN_OR_RETURN(
            bool is_loop_invariant,
            IsLoopInvariant(fbody, i,
                            flib_runtime->GetFunctionLibraryDefinition()));
        if (is_loop_invariant) {
          const_input_idxs->push_back(i);
        } else {
          Node* arg_i = fbody->arg_nodes[i];
          Node* ret_i = fbody->ret_nodes[i];
          VLOG(1) << "Argument " << i << " to while-loop " << node.name()
                  << " has to be constant, but it's not a loop invariant, "
                     "cluster compilation likely to fail at compile time: "
                  << arg_i->DebugString() << " vs. " << ret_i->DebugString();
          VLOG(1) << node.ShortDebugString();
        }
      }
    }
    return absl::OkStatus();
  } else if (node.op() == "If" || node.op() == "StatelessIf") {
    const FunctionBody* fthen = nullptr;
    const FunctionBody* felse = nullptr;
    TF_RETURN_IF_ERROR(
        GetFunctionBody(flib_runtime, node, "then_branch", &fthen));
    TF_RETURN_IF_ERROR(
        GetFunctionBody(flib_runtime, node, "else_branch", &felse));
    return CondConstInputIndices({fthen, felse}, const_input_idxs,
                                 flib_runtime);
  } else if (node.op() == "Case" || node.op() == "StatelessCase") {
    std::vector<const FunctionBody*> branch_bodies;
    TF_RETURN_IF_ERROR(
        GetFunctionBodies(flib_runtime, node, "branches", &branch_bodies));
    return CondConstInputIndices(branch_bodies, const_input_idxs, flib_runtime);
  } else if (node.op() == "PartitionedCall" ||
             node.op() == "StatefulPartitionedCall") {
    const FunctionBody* fbody;
    TF_RETURN_IF_ERROR(GetFunctionBody(flib_runtime, node, "f", &fbody));
    int num_inputs = fbody->record->fdef().signature().input_arg_size();
    std::vector<bool> compile_time_const_arg_indices(num_inputs);
    TF_RETURN_IF_ERROR(BackwardsConstAnalysis(
        *(fbody->graph), &compile_time_const_arg_indices,
        nullptr, flib_runtime));
    for (int i = 0; i < num_inputs; i++) {
      if (compile_time_const_arg_indices[i]) {
        const_input_idxs->push_back(i);
      }
    }
    return absl::OkStatus();
  } else if (op_def != nullptr) {
    return XlaOpRegistry::CompileTimeConstantInputs(node, *op_def,
                                                    const_input_idxs);
  } else {
    return XlaOpRegistry::CompileTimeConstantInputs(*op_kernel,
                                                    const_input_idxs);
  }
}
Status GetCompileTimeConstInputs(const Node* node,
                                 std::vector<int>* const_input_idxs,
                                 FunctionLibraryRuntime* flib_runtime) {
  return GetCompileTimeConstInputs(node->def(), nullptr,
                                   &node->op_def(), const_input_idxs,
                                   flib_runtime);
}
}  
Status BackwardsConstAnalysis(
    const Graph& g, std::vector<bool>* compile_time_const_arg_indices,
    std::vector<bool>* compile_time_const_nodes,
    FunctionLibraryRuntime* flib_runtime,
    std::function<bool(const Edge&)> edge_filter_input) {
  if (!compile_time_const_nodes && g.GetConstArgIndicesCache().has_value() &&
      !edge_filter_input) {
    VLOG(5) << "Using cached argument indices on graph " << &g;
    *compile_time_const_arg_indices = g.GetConstArgIndicesCache().value();
    return absl::OkStatus();
  }
  auto edge_filter = [&](const Edge& e) {
    return edge_filter_input ? edge_filter_input(e) : true;
  };
  std::vector<bool> compile_time_const_nodes_impl;
  if (compile_time_const_nodes) {
    CHECK_EQ(compile_time_const_nodes->size(), g.num_node_ids());
  } else {
    compile_time_const_nodes_impl.resize(g.num_node_ids());
    compile_time_const_nodes = &compile_time_const_nodes_impl;
  }
  Status status;
  auto visit = [&](Node* node) {
    if (!status.ok()) return;
    if (XlaOpRegistry::IsMetadataOp(node->type_string())) {
      VLOG(3) << "must-be-const node is metadata op: " << node->name();
      return;
    }
    if ((*compile_time_const_nodes)[node->id()]) {
      VLOG(3) << "marking consts for must-be-const node " << node->name();
      if (node->type_string() == "_Arg") {
        int index;
        status = GetNodeAttr(node->attrs(), "index", &index);
        if (!status.ok()) return;
        if (compile_time_const_arg_indices) {
          (*compile_time_const_arg_indices)[index] = true;
        }
        VLOG(3) << "  const _Arg " << index << ": " << node->name();
        return;
      }
      for (const Edge* pred : node->in_edges()) {
        if (!pred->IsControlEdge() && edge_filter(*pred)) {
          while (edge_filter(*pred) && IsConstTraversableOpType(pred->src())) {
            status = pred->src()->input_edge(pred->src_output(), &pred);
            if (!status.ok()) return;
          }
          if (edge_filter(*pred)) {
            VLOG(4) << "  " << pred->src()->name() << " must be const (is "
                    << pred->src()->type_string() << ")";
            (*compile_time_const_nodes)[pred->src()->id()] = true;
          }
        }
      }
      return;
    }
    std::vector<int> const_input_idxs;
    status = GetCompileTimeConstInputs(node, &const_input_idxs, flib_runtime);
    if (!status.ok() || const_input_idxs.empty()) {
      return;
    }
    VLOG(3) << "marking consts for must-be-const inputs of " << node->name();
    for (Edge const* edge : node->in_edges()) {
      if (!edge->IsControlEdge() &&
          absl::c_binary_search(const_input_idxs, edge->dst_input()) &&
          edge_filter(*edge)) {
        while (edge_filter(*edge) && IsConstTraversableOpType(edge->src())) {
          status = edge->src()->input_edge(edge->src_output(), &edge);
          if (!status.ok()) return;
        }
        if (edge_filter(*edge)) {
          VLOG(4) << "  input " << edge->dst_input() << ": "
                  << edge->src()->name() << " must be const (is "
                  << edge->src()->type_string() << ")";
          (*compile_time_const_nodes)[edge->src()->id()] = true;
        }
      }
    }
  };
  DFS(g, {}, visit, NodeComparatorName{},
      [](const Edge& edge) { return !edge.src()->IsNextIteration(); });
  if (compile_time_const_arg_indices && !edge_filter_input) {
    VLOG(5) << "Setting the cache on the graph: " << &g;
    g.GetConstArgIndicesCache() = *compile_time_const_arg_indices;
  }
  return status;
}
Status GetCompileTimeConstInputs(const OpKernel* op_kernel,
                                 std::vector<int>* const_input_idxs,
                                 FunctionLibraryRuntime* flib_runtime) {
  return GetCompileTimeConstInputs(op_kernel->def(), op_kernel,
                                   nullptr, const_input_idxs,
                                   flib_runtime);
}
}  