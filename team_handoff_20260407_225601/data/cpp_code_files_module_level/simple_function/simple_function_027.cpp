#ifndef TENSORFLOW_CORE_COMMON_RUNTIME_DEVICE_PROPAGATION_H_
#define TENSORFLOW_CORE_COMMON_RUNTIME_DEVICE_PROPAGATION_H_
#include <functional>
#include <string>
#include "absl/container/flat_hash_set.h"
#include "tensorflow/core/graph/graph.h"
#include "tensorflow/core/platform/stringpiece.h"
namespace tensorflow {
namespace device_propagation {
typedef std::function<bool(StringPiece)> DeviceFilter;
typedef std::function<bool(const Node&)> NodeFilter;
}  
void PropagateDevices(const device_propagation::NodeFilter& node_filter,
                      const device_propagation::DeviceFilter& device_filter,
                      Graph* graph);
}  
#endif  
#include "tensorflow/core/common_runtime/device_propagation.h"
#include <string>
#include <utility>
#include "absl/container/flat_hash_set.h"
#include "tensorflow/core/framework/types.h"
#include "tensorflow/core/graph/algorithm.h"
#include "tensorflow/core/graph/graph.h"
namespace tensorflow {
namespace {
const std::string& AssignedOrRequestedDevice(const Node& node) {
  if (!node.assigned_device_name().empty()) {
    return node.assigned_device_name();
  }
  return node.requested_device();
}
bool UpdateDeviceFromInputs(
    const device_propagation::NodeFilter& node_filter,
    const device_propagation::DeviceFilter& device_filter, Node* node) {
  if (!AssignedOrRequestedDevice(*node).empty() || !node_filter(*node)) {
    return false;
  }
  string proposed_device = "";
  Node* proposed_src = nullptr;
  for (const Edge* e : node->in_edges()) {
    if (e->IsControlEdge()) {
      continue;
    }
    Node* src = e->src();
    const string& src_device = AssignedOrRequestedDevice(*src);
    if ((node->IsSwitch() && src->IsLoopCond()) ||
        (node->IsMerge() && src->IsEnter())) {
      continue;
    }
    if (!device_filter(src_device)) return false;
    if (proposed_src == nullptr) {
      proposed_device = src_device;
      proposed_src = src;
    } else if (proposed_device != src_device) {
      return false;
    }
  }
  if (proposed_src) {
    node->set_assigned_device_name(proposed_src->assigned_device_name());
    node->set_requested_device(proposed_src->requested_device());
    return true;
  } else {
    return false;
  }
}
}  
void PropagateDevices(const device_propagation::NodeFilter& node_filter,
                      const device_propagation::DeviceFilter& device_filter,
                      Graph* graph) {
  bool nodes_changed = true;
  while (nodes_changed) {
    nodes_changed = false;
    BreadthFirstTraversal(
        *graph, {}, [&nodes_changed, &node_filter, &device_filter](Node* node) {
          nodes_changed |=
              UpdateDeviceFromInputs(node_filter, device_filter, node);
        });
  }
}
}  