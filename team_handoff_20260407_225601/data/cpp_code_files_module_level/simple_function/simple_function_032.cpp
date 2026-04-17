#ifndef XLA_SERVICE_HLO_PHI_GRAPH_H_
#define XLA_SERVICE_HLO_PHI_GRAPH_H_
#include <iterator>
#include <memory>
#include <string>
#include <vector>
#include "absl/container/flat_hash_map.h"
#include "absl/container/flat_hash_set.h"
#include "xla/hlo/ir/hlo_instruction.h"
#include "xla/hlo/ir/hlo_module.h"
#include "xla/service/hlo_value.h"
namespace xla {
class PhiGraph {
 public:
  void RegisterPhi(const HloValue& value,
                   absl::Span<const HloValue* const> inputs);
  HloValue::Id GetOptimizedId(const HloValue& value);
  bool InputsEqualTo(const HloValue& value,
                     absl::Span<const HloValue* const> inputs);
  HloValue::Id FindOptimizedValue(const HloValue::Id id);
  void Optimize();
  std::string ToString();
 private:
  struct Node {
    bool is_phi;
    std::vector<Node*> users;
    std::vector<Node*> operands;
    HloValue::Id value_id;
    bool mark_as_dead = false;
  };
  Node* CreateOrReuseNode(const HloValue& value);
  void ReplaceNodeWith(Node* node, Node* replace);
  absl::flat_hash_map<Node*, std::vector<HloValue::Id>> node_to_value_id_;
  absl::flat_hash_map<HloValue::Id, Node*> value_id_to_node_;
  std::vector<std::unique_ptr<Node>> node_storage_;
};
}  
#endif  
#include "xla/service/hlo_phi_graph.h"
#include <queue>
namespace xla {
HloValue::Id PhiGraph::GetOptimizedId(const HloValue& value) {
  Node* node = value_id_to_node_[value.id()];
  CHECK(!node->mark_as_dead);
  return node->value_id;
}
bool PhiGraph::InputsEqualTo(const HloValue& value,
                             absl::Span<const HloValue* const> inputs) {
  auto iter = value_id_to_node_.find(value.id());
  CHECK(iter != value_id_to_node_.end());
  absl::flat_hash_set<HloValue::Id> existing_set;
  for (Node* operand : iter->second->operands) {
    existing_set.insert(operand->value_id);
  }
  absl::flat_hash_set<HloValue::Id> new_set;
  for (const HloValue* input : inputs) {
    new_set.insert(input->id());
  }
  return existing_set == new_set;
}
HloValue::Id PhiGraph::FindOptimizedValue(const HloValue::Id id) {
  auto iter = value_id_to_node_.find(id);
  CHECK(iter != value_id_to_node_.end());
  CHECK(!iter->second->mark_as_dead);
  return iter->second->value_id;
}
PhiGraph::Node* PhiGraph::CreateOrReuseNode(const HloValue& value) {
  auto iter = value_id_to_node_.find(value.id());
  if (iter == value_id_to_node_.end()) {
    node_storage_.emplace_back(std::make_unique<Node>());
    Node* node = node_storage_.back().get();
    node->value_id = value.id();
    value_id_to_node_[value.id()] = node;
    node_to_value_id_[node].push_back(value.id());
    return node;
  } else {
    CHECK_NE(iter->second, nullptr);
    CHECK_EQ(iter->second->value_id, value.id());
    return iter->second;
  }
}
void PhiGraph::ReplaceNodeWith(PhiGraph::Node* node, PhiGraph::Node* replace) {
  CHECK(node->is_phi);
  if (node->mark_as_dead) {
    return;
  }
  if (replace->mark_as_dead) {
    auto iter = value_id_to_node_.find(replace->value_id);
    CHECK(iter != value_id_to_node_.end());
    return ReplaceNodeWith(node, iter->second);
  }
  CHECK(!replace->mark_as_dead);
  for (Node* user : node->users) {
    absl::c_replace(user->operands, node, replace);
  }
  for (Node* operand : node->operands) {
    absl::c_replace(operand->users, node, replace);
  }
  for (HloValue::Id value_id : node_to_value_id_[node]) {
    CHECK(value_id_to_node_.contains(value_id));
    value_id_to_node_[value_id] = replace;
  }
  absl::c_copy(node_to_value_id_[node],
               std::back_inserter(node_to_value_id_[replace]));
  node_to_value_id_[node].clear();
  node->mark_as_dead = true;
}
void PhiGraph::RegisterPhi(const HloValue& value,
                           absl::Span<const HloValue* const> inputs) {
  Node* node = CreateOrReuseNode(value);
  CHECK(value.is_phi());
  node->is_phi = true;
  node->operands.clear();
  for (auto input : inputs) {
    CHECK(input != nullptr);
    Node* input_node = CreateOrReuseNode(*input);
    node->operands.push_back(input_node);
  }
}
std::string PhiGraph::ToString() {
  std::string out = "PhiGraph: \n";
  for (auto& node : node_storage_) {
    absl::StrAppend(&out, node->value_id);
    if (node->is_phi) {
      absl::StrAppend(&out, ", phi");
    }
    if (node->mark_as_dead) {
      absl::StrAppend(&out, ", dead", ":\n");
    }
    for (Node* input : node->operands) {
      absl::StrAppend(&out, "  ", input->value_id, "\n");
    }
  }
  return out;
}
void PhiGraph::Optimize() {
  VLOG(2) << "Optimizing phi graph:";
  XLA_VLOG_LINES(2, ToString());
  for (auto& node : node_storage_) {
    for (Node* input : node->operands) {
      input->users.push_back(node.get());
    }
  }
  bool changed = true;
  while (changed) {
    changed = false;
    absl::flat_hash_set<Node*> checked_for_closure;
    for (auto& node : node_storage_) {
      if (!node->is_phi) {
        continue;
      }
      if (node->mark_as_dead) {
        continue;
      }
      Node* node_ptr = node.get();
      VLOG(2) << "Optimizing: " << node_ptr->value_id;
      CHECK_GE(node_ptr->operands.size(), 1);
      auto it = absl::c_find(node_ptr->operands, node_ptr);
      while (it != node_ptr->operands.end()) {
        node_ptr->operands.erase(it);
        it = absl::c_find(node_ptr->operands, node_ptr);
      }
      it = absl::c_find(node_ptr->users, node_ptr);
      while (it != node_ptr->users.end()) {
        node_ptr->users.erase(it);
        it = absl::c_find(node_ptr->users, node_ptr);
      }
      CHECK_GE(node_ptr->operands.size(), 1);
      bool all_inputs_are_same = absl::c_all_of(
          node_ptr->operands,
          [&](Node* elem) { return elem == node_ptr->operands[0]; });
      if (all_inputs_are_same) {
        VLOG(1) << "All inputs to node " << node_ptr->value_id
                << " are the same, replacing it with "
                << node_ptr->operands[0]->value_id;
        ReplaceNodeWith(node_ptr, node_ptr->operands[0]);
        changed = true;
        continue;
      }
      if (checked_for_closure.contains(node_ptr)) {
        continue;
      }
      absl::flat_hash_set<Node*> workset;
      std::queue<Node*> worklist;
      Node* non_phi = nullptr;
      worklist.push(node_ptr);
      while (!worklist.empty()) {
        Node* todo = worklist.front();
        worklist.pop();
        if (workset.contains(todo)) {
          continue;
        }
        checked_for_closure.insert(todo);
        workset.insert(todo);
        for (Node* operand : todo->operands) {
          worklist.push(operand);
        }
        if (!todo->is_phi) {
          if (non_phi != nullptr && non_phi != todo) {
            non_phi = nullptr;
            break;
          } else {
            non_phi = todo;
          }
        }
      }
      if (non_phi != nullptr) {
        for (Node* node : workset) {
          if (!node->is_phi) {
            CHECK_EQ(node, non_phi);
            continue;
          }
          VLOG(1) << "Replace node " << node->value_id
                  << " in the closure with node " << non_phi->value_id;
          ReplaceNodeWith(node, non_phi);
          changed = true;
        }
      }
    }
  }
}
}  