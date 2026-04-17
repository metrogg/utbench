#ifndef TENSORFLOW_CORE_GRAPPLER_UTILS_TRANSITIVE_FANIN_H_
#define TENSORFLOW_CORE_GRAPPLER_UTILS_TRANSITIVE_FANIN_H_
#include <unordered_map>
#include <vector>
#include "tensorflow/core/framework/graph.pb.h"
#include "tensorflow/core/lib/core/status.h"
namespace tensorflow {
namespace grappler {
Status ComputeTransitiveFanin(
    const GraphDef& graph, const std::vector<string>& terminal_nodes,
    std::unordered_map<string, const NodeDef*>* name_to_fanin_node,
    std::vector<const NodeDef*>* fanin_nodes);
Status ComputeTransitiveFanin(const GraphDef& graph,
                              const std::vector<string>& terminal_nodes,
                              std::vector<const NodeDef*>* fanin_nodes);
Status SetTransitiveFaninGraph(const GraphDef& input_graph,
                               GraphDef* output_graph,
                               const std::vector<string>& terminal_nodes);
}  
}  
#endif  
#include "tensorflow/core/grappler/utils/transitive_fanin.h"
#include <queue>
#include <vector>
#include "tensorflow/core/framework/node_def_util.h"
#include "tensorflow/core/grappler/utils.h"
#include "tensorflow/core/platform/errors.h"
namespace tensorflow {
namespace grappler {
Status ComputeTransitiveFanin(
    const GraphDef& graph, const std::vector<string>& terminal_nodes,
    std::unordered_map<string, const NodeDef*>* name_to_fanin_node,
    std::vector<const NodeDef*>* fanin_nodes) {
  std::unordered_map<string, const NodeDef*> name_to_node;
  std::unordered_map<string, const NodeDef*> name_to_send;
  for (const auto& node : graph.node()) {
    name_to_node[node.name()] = &node;
    if (node.op() == "_Send") {
      const auto& attr = node.attr();
      name_to_send[attr.at("tensor_name").s()] = &node;
    }
  }
  std::vector<const NodeDef*> queue;
  for (const string& root : terminal_nodes) {
    const NodeDef* node = name_to_node[NodeName(root)];
    if (!node) {
      return errors::InvalidArgument("Graph does not contain terminal node ",
                                     root, ".");
    }
    queue.push_back(node);
  }
  std::unordered_set<const NodeDef*> visited;
  while (!queue.empty()) {
    const NodeDef* node = queue.back();
    queue.pop_back();
    if (!visited.insert(node).second) {
      continue;
    }
    fanin_nodes->push_back(node);
    if (name_to_fanin_node) {
      name_to_fanin_node->insert(
          std::pair<string, const NodeDef*>(node->name(), node));
    }
    for (const string& input : node->input()) {
      const NodeDef* in = name_to_node[NodeName(input)];
      if (!in) {
        return errors::InvalidArgument("Graph does not contain input ",
                                       NodeName(input), " of node ",
                                       node->name(), ".");
      }
      queue.push_back(in);
    }
    if (node->op() == "_Recv") {
      const auto& attr = node->attr();
      const NodeDef* send = name_to_send[attr.at("tensor_name").s()];
      if (send) {
        queue.push_back(send);
      }
    }
  }
  return absl::OkStatus();
}
Status ComputeTransitiveFanin(const GraphDef& graph,
                              const std::vector<string>& terminal_nodes,
                              std::vector<const NodeDef*>* fanin_nodes) {
  return ComputeTransitiveFanin(graph, terminal_nodes, nullptr, fanin_nodes);
}
Status SetTransitiveFaninGraph(const GraphDef& input_graph,
                               GraphDef* output_graph,
                               const std::vector<string>& terminal_nodes) {
  std::vector<const NodeDef*> keep;
  TF_RETURN_IF_ERROR(
      ComputeTransitiveFanin(input_graph, terminal_nodes, &keep));
  output_graph->mutable_node()->Reserve(keep.size());
  for (int i = keep.size() - 1; i >= 0; --i) {
    *output_graph->add_node() = *keep[i];
  }
  return absl::OkStatus();
}
}  
}  