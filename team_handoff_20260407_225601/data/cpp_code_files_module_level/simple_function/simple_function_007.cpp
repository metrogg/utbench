#ifndef TENSORFLOW_CORE_FRAMEWORK_CONTROL_FLOW_H_
#define TENSORFLOW_CORE_FRAMEWORK_CONTROL_FLOW_H_
#include "tensorflow/core/lib/hash/hash.h"
#include "tensorflow/core/platform/logging.h"
#include "tensorflow/core/platform/types.h"
namespace tensorflow {
const uint64 kIllegalFrameId = ~0uLL;
const int64_t kIllegalIterId = -1;
struct FrameAndIter {
  uint64 frame_id = kIllegalFrameId;
  int64_t iter_id = kIllegalIterId;
  FrameAndIter() {}
  FrameAndIter(uint64 frame, int64_t iter) {
    frame_id = frame;
    iter_id = iter;
  }
  bool operator==(const FrameAndIter& other) const {
    return (frame_id == other.frame_id && iter_id == other.iter_id);
  }
};
struct FrameAndIterHash {
  size_t operator()(const FrameAndIter& key) const {
    CHECK_EQ(sizeof(uint64) + sizeof(int64_t), sizeof(FrameAndIter));
    return Hash64(reinterpret_cast<const char*>(&key), sizeof(FrameAndIter));
  }
};
}  
#endif  
#include "tensorflow/core/graph/control_flow.h"
#include <deque>
#include <vector>
#include "tensorflow/core/framework/node_def_util.h"
#include "tensorflow/core/framework/types.h"
#include "tensorflow/core/graph/node_builder.h"
#include "tensorflow/core/lib/core/errors.h"
namespace tensorflow {
namespace {
struct Frame {
  string name;
  Frame* parent = nullptr;
  const Node* loop_cond = nullptr;
};
Status ValidateControlFlowInfo(const Graph* graph,
                               const std::vector<ControlFlowInfo>& cf_info) {
  std::unordered_map<string, Frame> frames;
  for (const Node* node : graph->op_nodes()) {
    const ControlFlowInfo& cf = cf_info[node->id()];
    if (!cf.frame || !cf.parent_frame) {
      continue;
    }
    Frame& frame = frames[cf.frame_name];
    Frame* parent = &frames[cf_info[cf.parent_frame->id()].frame_name];
    if (frame.parent == nullptr) {
      frame.parent = parent;
      frame.name = cf.frame_name;
    } else if (frame.parent != parent) {
      return errors::Internal(
          "Invalid loop structure: Mismatched parent frames for \"",
          cf.frame_name, "\": \"", parent->name, "\" vs \"", frame.parent->name,
          "\". The node giving this error: ", FormatNodeForError(*node),
          ". This is an internal bug, please file a bug report with "
          "instructions on how to reproduce the error.");
    }
    if (IsLoopCond(node)) {
      if (frame.loop_cond &&
          !absl::StrContains(frame.loop_cond->name(), "LoopCounter") &&
          !absl::StrContains(node->name(), "LoopCounter")) {
        return errors::InvalidArgument(
            "Invalid loop structure: Loop \"", cf.frame_name,
            "\" has more than one LoopCond node: ", FormatNodeForError(*node),
            " and ", FormatNodeForError(*frame.loop_cond),
            ". This is an internal bug, please file a bug report with "
            "instructions on how to reproduce the error.");
      }
      frame.loop_cond = node;
    }
  }
  return absl::OkStatus();
}
}  
Status BuildControlFlowInfo(const Graph* g, std::vector<ControlFlowInfo>* info,
                            std::vector<string>* unreachable_nodes) {
  info->clear();
  info->resize(g->num_node_ids());
  std::vector<const Node*> parent_nodes;
  parent_nodes.resize(g->num_node_ids());
  const Node* src_node = g->source_node();
  ControlFlowInfo& src_info = (*info)[src_node->id()];
  src_info.frame = src_node;
  src_info.parent_frame = src_node;
  string frame_name;
  std::deque<const Node*> ready;
  ready.push_back(src_node);
  while (!ready.empty()) {
    const Node* curr_node = ready.front();
    ready.pop_front();
    const ControlFlowInfo& curr_info = (*info)[curr_node->id()];
    const Node* frame = curr_info.frame;
    const Node* parent = curr_info.parent_frame;
    frame_name = curr_info.frame_name;
    if (IsExit(curr_node)) {
      const ControlFlowInfo& parent_info = (*info)[parent->id()];
      frame = parent_info.frame;
      parent = parent_info.parent_frame;
      frame_name = parent_info.frame_name;
    }
    for (const Edge* out_edge : curr_node->out_edges()) {
      const Node* out = out_edge->dst();
      int out_id = out->id();
      ControlFlowInfo* out_info = &(*info)[out_id];
      const Node* out_parent = out_info->parent_frame;
      bool is_visited = (parent_nodes[out_id] != nullptr);
      if (!out->IsOp()) continue;
      if (!is_visited) {
        parent_nodes[out->id()] = curr_node;
        ready.push_back(out);
      }
      if (IsEnter(out)) {
        if (is_visited) {
          const string& parent_frame = (*info)[out_parent->id()].frame_name;
          if (parent_frame != frame_name) {
            return errors::InvalidArgument(
                FormatNodeForError(*out),
                " has inputs from different frames. The input ",
                FormatNodeForError(*curr_node), " is in frame '", frame_name,
                "'. The input ", FormatNodeForError(*parent_nodes[out->id()]),
                " is in frame '", parent_frame, "'.");
          }
        } else {
          out_info->frame = out;
          out_info->parent_frame = frame;
          TF_RETURN_IF_ERROR(
              GetNodeAttr(out->attrs(), "frame_name", &out_info->frame_name));
          if (out_info->frame_name.empty()) {
            return errors::InvalidArgument("The Enter ",
                                           FormatNodeForError(*out),
                                           " must have a frame name.");
          }
        }
      } else {
        if (is_visited) {
          if (out_info->frame_name != frame_name) {
            return errors::InvalidArgument(
                FormatNodeForError(*out),
                " has inputs from different frames. The input ",
                FormatNodeForError(*curr_node), " is in frame '", frame_name,
                "'. The input ", FormatNodeForError(*parent_nodes[out->id()]),
                " is in frame '", out_info->frame_name, "'.");
          }
        } else {
          out_info->frame = frame;
          out_info->parent_frame = parent;
          out_info->frame_name = frame_name;
        }
      }
    }
  }
  if (unreachable_nodes) {
    for (const Node* node : g->op_nodes()) {
      if (!parent_nodes[node->id()]) {
        unreachable_nodes->push_back(node->name());
      }
    }
  }
  TF_RETURN_IF_ERROR(ValidateControlFlowInfo(g, *info));
  return absl::OkStatus();
}
}  