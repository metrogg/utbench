#ifndef TENSORFLOW_CORE_COMMON_RUNTIME_PARTITIONING_UTILS_H_
#define TENSORFLOW_CORE_COMMON_RUNTIME_PARTITIONING_UTILS_H_
#include <unordered_map>
#include <vector>
#include "tensorflow/core/common_runtime/device_set.h"
#include "tensorflow/core/framework/function.h"
#include "tensorflow/core/lib/core/status.h"
namespace tensorflow {
Status PartitionFunctionGraph(
    const DeviceSet& device_set, std::unique_ptr<Graph> graph,
    std::unordered_map<string, std::unique_ptr<Graph>>* subgraphs,
    std::function<string(const Edge*)> get_tensor_name_attr = nullptr);
absl::StatusOr<std::unique_ptr<Graph>> InsertTransferOps(
    const DeviceSet& device_set, std::unique_ptr<Graph> graph);
Status UpdateArgAndRetvalMetadata(
    Graph* graph, std::vector<FunctionArgIndex>* arg_indices,
    std::vector<int>* ret_indices,
    std::vector<AllocatorAttributes>* arg_alloc_attrs,
    std::vector<AllocatorAttributes>* ret_alloc_attrs, bool ints_on_device);
class FunctionNameGenerator {
 public:
  FunctionNameGenerator(const FunctionLibraryDefinition* flib_def,
                        const string& name)
      : flib_def_(flib_def), name_(name), counter_(0) {}
  string GetName();
 private:
  const FunctionLibraryDefinition* flib_def_;
  const string name_;
  uint32 counter_;
};
}  
#endif  
#include "tensorflow/core/common_runtime/partitioning_utils.h"
#include <algorithm>
#include <functional>
#include <memory>
#include <optional>
#include <string>
#include <unordered_map>
#include <utility>
#include "tensorflow/core/common_runtime/arg_ret_placement.h"
#include "tensorflow/core/common_runtime/graph_constructor.h"
#include "tensorflow/core/framework/function.h"
#include "tensorflow/core/framework/op.h"
#include "tensorflow/core/framework/types.h"
#include "tensorflow/core/graph/graph.h"
#include "tensorflow/core/graph/graph_partition.h"
namespace tensorflow {
namespace {
Status PartitionFunctionGraph(
    const DeviceSet& device_set, Graph* graph,
    std::unordered_map<string, GraphDef>* partitions,
    std::function<string(const Node*)> node_to_loc,
    std::function<string(const Edge*)> get_tensor_name_attr) {
  PartitionOptions partition_options;
  if (node_to_loc != nullptr) {
    partition_options.node_to_loc = node_to_loc;
  } else {
    partition_options.node_to_loc = [](const Node* node) {
      return node->assigned_device_name();
    };
  }
  int64_t edge_name_counter = 0;
  partition_options.new_name = [&edge_name_counter](const string& prefix) {
    return strings::StrCat(prefix, "/_", ++edge_name_counter);
  };
  partition_options.get_incarnation =
      [&device_set](const string& name) -> int64 {
    const Device* d = device_set.FindDeviceByName(name);
    if (d == nullptr) {
      return PartitionOptions::kIllegalIncarnation;
    } else {
      return d->attributes().incarnation();
    }
  };
  partition_options.control_flow_added = false;
  partition_options.get_tensor_name_attr = get_tensor_name_attr;
  partition_options.can_make_destructive_changes = true;
  return Partition(partition_options, graph, partitions);
}
struct SendRecvPair {
  Node* send_node = nullptr;
  Node* recv_node = nullptr;
};
constexpr char kTensorNameAttr[] = "tensor_name";
Status MakeSendRecvDependencyExplicit(Graph* graph) {
  absl::flat_hash_map<std::string, SendRecvPair> send_recv_pairs;
  for (Node* node : graph->op_nodes()) {
    if (node->IsSend() || node->IsRecv()) {
      auto tensor_name_it = node->def().attr().find(kTensorNameAttr);
      if (tensor_name_it == node->def().attr().end()) {
        return errors::Internal(
            "'", kTensorNameAttr,
            "' attribute is not found from node: ", node->DebugString());
      }
      if (node->IsSend()) {
        send_recv_pairs[tensor_name_it->second.s()].send_node = node;
      } else {
        send_recv_pairs[tensor_name_it->second.s()].recv_node = node;
      }
    }
  }
  for (const auto& [tensor_name, send_recv_pair] : send_recv_pairs) {
    if (send_recv_pair.send_node == nullptr ||
        send_recv_pair.recv_node == nullptr) {
      return errors::Internal(
          "No matching Send/Recv nodes found for tensor_name = ", tensor_name);
    }
    graph->AddControlEdge(send_recv_pair.send_node, send_recv_pair.recv_node);
  }
  return absl::OkStatus();
}
}  
Status PartitionFunctionGraph(
    const DeviceSet& device_set, std::unique_ptr<Graph> graph,
    std::unordered_map<string, std::unique_ptr<Graph>>* subgraphs,
    std::function<string(const Edge*)> get_tensor_name_attr) {
  std::unordered_map<string, GraphDef> partitions;
  TF_RETURN_IF_ERROR(
      PartitionFunctionGraph(device_set, graph.get(), &partitions,
                             nullptr, get_tensor_name_attr));
  const OpRegistryInterface* default_registry =
      graph->flib_def().default_registry();
  graph.reset();
  for (auto& partition : partitions) {
    const string& device = partition.first;
    GraphDef& graph_def = partition.second;
    auto subgraph = std::make_unique<Graph>(default_registry);
    GraphConstructorOptions opts;
    opts.allow_internal_ops = true;
    opts.expect_device_spec = true;
    TF_RETURN_IF_ERROR(
        ConvertGraphDefToGraph(opts, std::move(graph_def), subgraph.get()));
    subgraphs->emplace(device, std::move(subgraph));
  }
  return absl::OkStatus();
}
absl::StatusOr<std::unique_ptr<Graph>> InsertTransferOps(
    const DeviceSet& device_set, std::unique_ptr<Graph> graph) {
  auto node_to_loc = [](const Node* node) {
    return node->assigned_device_name();
  };
  bool has_multiple_devices = false;
  absl::optional<std::string> location;
  for (const Node* node : graph->op_nodes()) {
    if (location) {
      if (*location != node_to_loc(node)) {
        has_multiple_devices = true;
        break;
      }
    } else {
      location = node_to_loc(node);
    }
  }
  if (!has_multiple_devices) {
    return graph;
  }
  auto new_graph = std::make_unique<Graph>(graph->flib_def());
  std::unordered_map<string, GraphDef> partitions;
  TF_RETURN_IF_ERROR(PartitionFunctionGraph(device_set, graph.get(),
                                            &partitions, node_to_loc,
                                            nullptr));
  GraphDef merged_graph_def;
  if (!partitions.empty()) {
    auto iter = partitions.begin();
    merged_graph_def = std::move(iter->second);
    while (++iter != partitions.end()) {
      merged_graph_def.MergeFrom(iter->second);
    }
  }
  GraphConstructorOptions opts;
  opts.allow_internal_ops = true;
  opts.expect_device_spec = true;
  TF_RETURN_IF_ERROR(ConvertGraphDefToGraph(opts, std::move(merged_graph_def),
                                            new_graph.get()));
  TF_RETURN_IF_ERROR(MakeSendRecvDependencyExplicit(new_graph.get()));
  return std::move(new_graph);
}
Status UpdateArgAndRetvalMetadata(
    Graph* graph, std::vector<FunctionArgIndex>* arg_indices,
    std::vector<int>* ret_indices,
    std::vector<AllocatorAttributes>* arg_alloc_attrs,
    std::vector<AllocatorAttributes>* ret_alloc_attrs, bool ints_on_device) {
  std::vector<std::pair<Node*, FunctionArgIndex>> arg_nodes;
  std::vector<std::pair<Node*, int>> ret_nodes;
  const AttrValue* attr_value;
  for (Node* node : graph->op_nodes()) {
    if (node->IsArg()) {
      TF_RETURN_IF_ERROR(node->attrs().Find("index", &attr_value));
      int index = static_cast<int>(attr_value->i());
      int sub_index = -1;
      if (node->attrs().Find("sub_index", &attr_value).ok()) {
        sub_index = static_cast<int>(attr_value->i());
      }
      arg_nodes.emplace_back(node, FunctionArgIndex(index, sub_index));
    } else if (node->IsRetval()) {
      TF_RETURN_IF_ERROR(node->attrs().Find("index", &attr_value));
      int index = static_cast<int>(attr_value->i());
      ret_nodes.emplace_back(node, index);
    }
  }
  auto arg_comparator = [](std::pair<Node*, FunctionArgIndex> a,
                           std::pair<Node*, FunctionArgIndex> b) {
    return std::tie(a.second.index, a.second.sub_index) <
           std::tie(b.second.index, b.second.sub_index);
  };
  std::sort(arg_nodes.begin(), arg_nodes.end(), arg_comparator);
  auto ret_comparator = [](std::pair<Node*, int> a, std::pair<Node*, int> b) {
    return a.second < b.second;
  };
  std::sort(ret_nodes.begin(), ret_nodes.end(), ret_comparator);
  arg_indices->reserve(arg_nodes.size());
  for (const auto& pair : arg_nodes) arg_indices->push_back(pair.second);
  ret_indices->reserve(ret_nodes.size());
  for (const auto& pair : ret_nodes) ret_indices->push_back(pair.second);
  for (int i = 0; i < arg_nodes.size(); ++i) {
    Node* arg = arg_nodes[i].first;
    arg->AddAttr("index", i);
  }
  if (arg_alloc_attrs != nullptr) {
    TF_RETURN_IF_ERROR(full_type::SingleDeviceSetAllocAttrsForArgs(
        arg_nodes, ints_on_device, *arg_alloc_attrs));
  }
  for (int i = 0; i < ret_nodes.size(); ++i) {
    Node* ret = ret_nodes[i].first;
    ret->AddAttr("index", i);
  }
  if (ret_alloc_attrs) {
    TF_RETURN_IF_ERROR(full_type::SingleDeviceSetAllocAttrsForRets(
        ret_nodes, ints_on_device, *ret_alloc_attrs));
  }
  return absl::OkStatus();
}
string FunctionNameGenerator::GetName() {
  while (true) {
    const string candidate = strings::StrCat(name_, "_", counter_++);
    if (flib_def_->Find(candidate) == nullptr) {
      return candidate;
    }
  }
}
}  