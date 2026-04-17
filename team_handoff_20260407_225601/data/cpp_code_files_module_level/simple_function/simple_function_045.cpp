#ifndef XLA_HLO_IR_HLO_DFS_REACHABILITY_H_
#define XLA_HLO_IR_HLO_DFS_REACHABILITY_H_
#include <cstddef>
#include <memory>
#include "llvm/ADT/DenseMap.h"
#include "xla/hlo/ir/hlo_computation.h"
#include "xla/hlo/ir/hlo_instruction.h"
namespace xla {
class HloDfsReachability {
 public:
  bool IsPresent(const HloInstruction* instruction) const;
  bool IsReachable(const HloInstruction* from, const HloInstruction* to) const;
  bool IsConnected(const HloInstruction* a, const HloInstruction* b) const;
  static std::unique_ptr<HloDfsReachability> Build(
      const HloComputation* computation);
 private:
  llvm::DenseMap<const HloInstruction*, size_t> instruction_to_idx_;
};
}  
#endif  
#include "xla/hlo/ir/hlo_dfs_reachability.h"
#include <cstddef>
#include <memory>
#include <vector>
#include "absl/algorithm/container.h"
#include "llvm/ADT/BitVector.h"
#include "llvm/ADT/SmallVector.h"
#include "xla/hlo/ir/hlo_computation.h"
#include "xla/hlo/ir/hlo_instruction.h"
namespace xla {
bool HloDfsReachability::IsPresent(const HloInstruction* instruction) const {
  return instruction_to_idx_.contains(instruction);
}
bool HloDfsReachability::IsReachable(const HloInstruction* from,
                                     const HloInstruction* to) const {
  if (from == to) {
    return true;
  }
  if (to->operand_count() == 0 && from->control_predecessors().empty()) {
    return false;
  }
  size_t target_node_idx = instruction_to_idx_.at(from);
  size_t dfs_root_idx = instruction_to_idx_.at(to);
  if (dfs_root_idx < target_node_idx) {
    return false;
  }
  llvm::SmallVector<const HloInstruction*> stack{to};
  llvm::BitVector visited_idxs(1 + (dfs_root_idx - target_node_idx));
  visited_idxs.set(dfs_root_idx - target_node_idx);
  auto check_and_enqueue = [&](const HloInstruction* instr) {
    if (instr == from) {
      return true;
    }
    size_t instr_idx = instruction_to_idx_.at(instr);
    if (instr_idx < target_node_idx) {
      return false;
    }
    size_t visited_idx = instr_idx - target_node_idx;
    if (visited_idxs.test(visited_idx)) {
      return false;
    }
    visited_idxs.set(visited_idx);
    stack.push_back(instr);
    return false;
  };
  while (!stack.empty()) {
    const HloInstruction* instr = stack.pop_back_val();
    if (absl::c_any_of(instr->operands(), check_and_enqueue) ||
        absl::c_any_of(instr->control_predecessors(), check_and_enqueue)) {
      return true;
    }
  }
  return false;
}
bool HloDfsReachability::IsConnected(const HloInstruction* a,
                                     const HloInstruction* b) const {
  return IsReachable(a, b) || IsReachable(b, a);
}
std::unique_ptr<HloDfsReachability> HloDfsReachability::Build(
    const HloComputation* computation) {
  auto res = std::make_unique<HloDfsReachability>();
  HloComputation::ChannelDependencies empty_channel_dependencies;
  std::vector<HloInstruction*> instructions =
      computation->MakeInstructionPostOrder(empty_channel_dependencies);
  res->instruction_to_idx_.reserve(instructions.size());
  for (size_t i = 0; i < instructions.size(); ++i) {
    res->instruction_to_idx_[instructions[i]] = i;
  }
  return res;
}
}  