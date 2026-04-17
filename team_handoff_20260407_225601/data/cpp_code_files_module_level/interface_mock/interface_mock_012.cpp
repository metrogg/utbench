#ifndef TENSORFLOW_LITE_GRAPH_INFO_H_
#define TENSORFLOW_LITE_GRAPH_INFO_H_
#include <stddef.h>
#include <cstdint>
#include <utility>
#include <vector>
#include "tensorflow/lite/core/c/common.h"
namespace tflite {
class GraphInfo {
 public:
  virtual ~GraphInfo() {}
  virtual size_t num_tensors() const = 0;
  virtual TfLiteTensor* tensor(size_t index) = 0;
  virtual TfLiteTensor* tensors() = 0;
  virtual size_t num_execution_nodes() const = 0;
  virtual size_t num_total_nodes() const = 0;
  virtual const TfLiteNode& node(size_t index) const = 0;
  virtual const TfLiteRegistration& registration(size_t index) const = 0;
  virtual size_t node_index(size_t index) const = 0;
  virtual const std::vector<int>& inputs() const = 0;
  virtual const std::vector<int>& outputs() const = 0;
  virtual const std::vector<int>& variables() const = 0;
};
struct NodeSubset {
  enum Type {
    kTfUnexplored = 0,  
    kTfPartition,
    kTfNonPartition
  };
  Type type = kTfUnexplored;
  std::vector<int> nodes;
  std::vector<int> input_tensors;
  std::vector<int> output_tensors;
};
using ControlEdge = std::pair<int32_t, int32_t>;
using ControlEdges = std::vector<ControlEdge>;
TfLiteStatus PartitionGraphIntoIndependentNodeSubsets(
    const GraphInfo* info, const TfLiteIntArray* nodes_to_partition,
    std::vector<NodeSubset>* node_subsets, bool greedily,
    const ControlEdges* control_edges = nullptr);
}  
#endif  
#include "tensorflow/lite/graph_info.h"
#include <algorithm>
#include <vector>
#include "tensorflow/lite/context_util.h"
#include "tensorflow/lite/core/c/common.h"
namespace tflite {
namespace {
template <class T>
void Uniquefy(std::vector<T>* items) {
  std::sort(items->begin(), items->end());
  items->erase(std::unique(items->begin(), items->end()), items->end());
}
class PartitionGraphIntoIndependentNodeSubsetsImpl {
 public:
  PartitionGraphIntoIndependentNodeSubsetsImpl(
      const GraphInfo* info, const TfLiteIntArray* nodes_to_partition,
      std::vector<NodeSubset>* node_subsets, bool greedily,
      const ControlEdges& control_edges)
      : info_(info),
        node_subsets_(node_subsets),
        node_type_(info_->num_total_nodes(), NodeSubset::kTfNonPartition),
        greedily_(greedily),
        control_edges_(control_edges),
        num_incoming_control_edges_(info_->num_execution_nodes(), 0) {
    for (auto node_index : TfLiteIntArrayView(nodes_to_partition)) {
      node_type_[node_index] = NodeSubset::kTfPartition;
    }
    Uniquefy(&control_edges_);
  }
  void Partition() {
    node_subsets_->clear();
    tensor_epochs_.clear();
    tensor_epochs_.resize(info_->num_tensors(), kEpochAlwaysReady);
    node_epochs_.clear();
    node_epochs_.resize(info_->num_execution_nodes(), kEpochNotReady);
    num_incoming_control_edges_.clear();
    num_incoming_control_edges_.resize(info_->num_execution_nodes(), 0);
    for (const auto& edge : control_edges_) {
      ++num_incoming_control_edges_[edge.second];
    }
    for (int node_index = 0; node_index < info_->num_execution_nodes();
         node_index++) {
      const TfLiteNode& node = info_->node(node_index);
      for (int output_tensor_index : TfLiteIntArrayView(node.outputs)) {
        if (output_tensor_index == kTfLiteOptionalTensor) continue;
        tensor_epochs_[output_tensor_index] = kEpochNotReady;
      }
    }
    while (true) {
      BuildNodeSubset();
      if (node_subsets_->back().nodes.empty()) {
        node_subsets_->pop_back();
        break;
      }
    }
    for (int output_index : info_->outputs()) {
      int output_epoch = tensor_epochs_[output_index];
      if (output_epoch == kEpochAlwaysReady) {
        continue;
      }
      NodeSubset& output_subset = (*node_subsets_)[output_epoch];
      output_subset.output_tensors.push_back(output_index);
    }
    for (NodeSubset& node_subset : *node_subsets_) {
      Uniquefy(&node_subset.input_tensors);
      Uniquefy(&node_subset.output_tensors);
    }
  }
 private:
  enum {
    kEpochNotReady = -1,
    kEpochAlwaysReady = -2
  };
  bool UpdateNode(int node_index) {
    const TfLiteNode& node = info_->node(node_index);
    NodeSubset& current_subset = node_subsets_->back();
    int current_epoch = node_subsets_->size() - 1;
    if (node_epochs_[node_index] != kEpochNotReady) {
      return false;
    }
    for (int input_tensor_index : TfLiteIntArrayView(node.inputs)) {
      if (input_tensor_index != kTfLiteOptionalTensor &&
          tensor_epochs_[input_tensor_index] == kEpochNotReady) {
        return false;
      }
    }
    if (num_incoming_control_edges_[node_index] != 0) {
      return false;
    }
    int original_node_idx = info_->node_index(node_index);
    if (current_subset.type == NodeSubset::kTfUnexplored) {
      current_subset.type = node_type_[original_node_idx];
    }
    if (current_subset.type == node_type_[original_node_idx]) {
      node_epochs_[node_index] = current_epoch;
      current_subset.nodes.push_back(original_node_idx);
      for (int output_tensor_index : TfLiteIntArrayView(node.outputs)) {
        if (output_tensor_index == kTfLiteOptionalTensor) continue;
        tensor_epochs_[output_tensor_index] = current_epoch;
      }
      for (int input_tensor_index : TfLiteIntArrayView(node.inputs)) {
        if (input_tensor_index == kTfLiteOptionalTensor) {
          continue;
        }
        int input_epoch = tensor_epochs_[input_tensor_index];
        int node_epoch = current_epoch;
        if (input_epoch != node_epoch) {
          current_subset.input_tensors.push_back(input_tensor_index);
          if (input_epoch >= 0) {
            NodeSubset& input_subset = (*node_subsets_)[input_epoch];
            input_subset.output_tensors.push_back(input_tensor_index);
          }
        }
      }
      for (auto edge_iter =
               std::lower_bound(control_edges_.begin(), control_edges_.end(),
                                ControlEdge(node_index, 0));
           edge_iter != control_edges_.end() && edge_iter->first == node_index;
           ++edge_iter) {
        --num_incoming_control_edges_[edge_iter->second];
      }
      return true;
    } else {
      return false;
    }
  }
  void BuildNodeSubset() {
    node_subsets_->emplace_back(NodeSubset());
    while (true) {
      bool did_something = false;
      for (int node_index = 0; node_index < info_->num_execution_nodes();
           node_index++) {
        if (UpdateNode(node_index)) {
          did_something = true;
        } else {
          if (did_something && !greedily_) {
            return;
          }
        }
      }
      if (!did_something) return;
    }
  }
  const GraphInfo* info_;
  std::vector<NodeSubset>* node_subsets_;
  std::vector<NodeSubset::Type> node_type_;
  std::vector<int> tensor_epochs_;
  std::vector<int> node_epochs_;
  const bool greedily_;
  ControlEdges control_edges_;
  std::vector<int> num_incoming_control_edges_;
};
}  
TfLiteStatus PartitionGraphIntoIndependentNodeSubsets(
    const GraphInfo* info, const TfLiteIntArray* nodes_to_partition,
    std::vector<NodeSubset>* node_subsets, bool greedily,
    const ControlEdges* control_edges) {
  ControlEdges my_control_edges;
  if (control_edges == nullptr) {
    control_edges = &my_control_edges;
    if (greedily) {
      for (int last_op_with_side_effect = -1, node_index = 0;
           node_index < info->num_execution_nodes(); ++node_index) {
        const auto& node = info->node(node_index);
        if (node.might_have_side_effect) {
          if (last_op_with_side_effect != -1) {
            my_control_edges.emplace_back(last_op_with_side_effect, node_index);
          }
          last_op_with_side_effect = node_index;
        }
      }
    }
  }
  PartitionGraphIntoIndependentNodeSubsetsImpl(
      info, nodes_to_partition, node_subsets, greedily, *control_edges)
      .Partition();
  return kTfLiteOk;
}
}  