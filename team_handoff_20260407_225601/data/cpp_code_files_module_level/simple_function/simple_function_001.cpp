#ifndef XLA_HLO_TRANSFORMS_HLO_CONSTANT_SPLITTER_H_
#define XLA_HLO_TRANSFORMS_HLO_CONSTANT_SPLITTER_H_
#include "absl/container/flat_hash_set.h"
#include "absl/functional/function_ref.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "xla/hlo/ir/hlo_instruction.h"
#include "xla/hlo/ir/hlo_module.h"
#include "xla/service/hlo_pass_interface.h"
namespace xla {
class HloConstantSplitter : public HloModulePass {
 public:
  explicit HloConstantSplitter(
      bool split_expressions = false,
      absl::FunctionRef<bool(const HloInstruction*)> extra_constraints =
          [](const HloInstruction* instruction) { return true; })
      : split_expressions_(split_expressions),
        extra_constraints_(extra_constraints) {}
  absl::string_view name() const override { return "hlo-constant-splitter"; }
  using HloPassInterface::Run;
  absl::StatusOr<bool> Run(
      HloModule* module,
      const absl::flat_hash_set<absl::string_view>& execution_threads) override;
 private:
  bool split_expressions_;
  absl::FunctionRef<bool(const HloInstruction*)> extra_constraints_;
};
}  
#endif  
#include "xla/hlo/transforms/hlo_constant_splitter.h"
#include <iterator>
#include <utility>
#include <vector>
#include "absl/container/flat_hash_map.h"
#include "absl/container/flat_hash_set.h"
#include "absl/container/inlined_vector.h"
#include "absl/log/check.h"
#include "absl/log/log.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "xla/hlo/ir/hlo_instruction.h"
#include "xla/hlo/ir/hlo_opcode.h"
#include "tsl/platform/errors.h"
#include "tsl/platform/statusor.h"
namespace xla {
namespace {
bool IsSupportedConstant(const HloInstruction* instruction,
                         bool split_expressions) {
  return instruction->opcode() == HloOpcode::kConstant ||
         (split_expressions && instruction->opcode() == HloOpcode::kIota);
}
bool IsSupportedConstantExpression(const HloInstruction* instruction) {
  if (instruction->HasSideEffect()) {
    return false;
  }
  if (instruction->IsElementwise()) {
    return true;
  }
  switch (instruction->opcode()) {
    case HloOpcode::kBroadcast:
    case HloOpcode::kSlice:
      return true;
    default:
      return false;
  }
}
absl::StatusOr<bool> DuplicateConstantExpressionPerUser(
    HloComputation* computation, HloInstruction* to_clone,
    HloInstruction* user) {
  absl::InlinedVector<std::pair<const HloInstruction*, int>, 8> worklist(
      1, std::make_pair(to_clone, 0));
  absl::InlinedVector<const HloInstruction*, 8> to_clone_vec;
  absl::flat_hash_set<const HloInstruction*> visited;
  bool changed = false;
  VLOG(10) << "Duplicating: " << to_clone->ToString() << " for user "
           << user->ToString();
  while (!worklist.empty()) {
    auto& [to_clone_i, index] = worklist.back();
    if (index >= to_clone_i->operand_count()) {
      to_clone_vec.push_back(to_clone_i);
      worklist.pop_back();
      continue;
    }
    int64_t prev_idx = index++;
    if (visited.insert(to_clone_i->operands()[prev_idx]).second) {
      VLOG(10) << "Adding operand to worklist: "
               << to_clone_i->operands()[prev_idx]->ToString();
      worklist.push_back(std::make_pair(to_clone_i->operands()[prev_idx], 0));
    }
  }
  absl::flat_hash_map<const HloInstruction*, HloInstruction*>
      cloned_instructions_map;
  for (auto* i : to_clone_vec) {
    absl::InlinedVector<HloInstruction*, 4> new_operand_vector;
    for (auto* op : i->operands()) {
      auto it = cloned_instructions_map.find(op);
      CHECK(it != cloned_instructions_map.end())
          << "Expected already cloned instruction for operand: "
          << op->ToString() << " Instruction to clone: " << i->ToString();
      new_operand_vector.push_back(it->second);
    }
    HloInstruction* cloned_instr = computation->AddInstruction(
        i->CloneWithNewOperands(i->shape(), new_operand_vector));
    cloned_instructions_map[i] = cloned_instr;
    if (i == to_clone) {
      TF_RETURN_IF_ERROR(to_clone->ReplaceUseWith(user, cloned_instr));
      changed = true;
    }
  }
  return changed;
}
}  
absl::StatusOr<bool> HloConstantSplitter::Run(
    HloModule* module,
    const absl::flat_hash_set<absl::string_view>& execution_threads) {
  bool changed = false;
  for (HloComputation* computation : module->computations(execution_threads)) {
    absl::flat_hash_set<HloInstruction*> constants_set;
    std::vector<HloInstruction*> constants_list;
    std::vector<HloInstruction*> worklist;
    for (HloInstruction* instruction :
         computation->MakeInstructionPostOrder()) {
      VLOG(10) << "Considering: " << instruction->ToString();
      if (IsSupportedConstant(instruction, split_expressions_) &&
          extra_constraints_(instruction)) {
        VLOG(10) << "Adding to constant list: " << instruction->ToString();
        constants_set.insert(instruction);
        constants_list.push_back(instruction);
      }
    }
    int64_t previous_total_constants = 0;
    while (constants_list.size() != previous_total_constants) {
      VLOG(10) << "Previous total: " << previous_total_constants
               << " current constants: " << constants_list.size();
      previous_total_constants = constants_list.size();
      worklist.clear();
      worklist.insert(worklist.end(), constants_list.begin(),
                      constants_list.end());
      while (!worklist.empty()) {
        auto* i = worklist.back();
        worklist.pop_back();
        bool is_constant = true;
        for (auto* ops : i->operands()) {
          if (!constants_set.contains(ops)) {
            is_constant = false;
            break;
          }
        }
        if (is_constant) {
          if (constants_set.insert(i).second) {
            constants_list.push_back(i);
          }
          if (split_expressions_) {
            for (auto* u : i->users()) {
              if (IsSupportedConstantExpression(u) &&
                  !constants_set.contains(u)) {
                worklist.push_back(u);
              }
            }
          }
        }
      }
    }
    if (VLOG_IS_ON(5)) {
      VLOG(5) << "For computation: " << computation->ToString();
      for (HloInstruction* instruction : constants_list) {
        VLOG(5) << "Is a constant: " << instruction->ToString();
      }
    }
    for (HloInstruction* instruction : constants_list) {
      if (IsSupportedConstant(instruction, split_expressions_) &&
          instruction->user_count() <= 1) {
        continue;
      }
      absl::InlinedVector<HloInstruction*, 8> users;
      users.reserve(instruction->user_count());
      for (HloInstruction* user : instruction->users()) {
        if (instruction->opcode() == HloOpcode::kConstant ||
            !constants_set.contains(user)) {
          users.push_back(user);
        }
      }
      for (auto* u : users) {
        TF_ASSIGN_OR_RETURN(bool duplicated, DuplicateConstantExpressionPerUser(
                                                 computation, instruction, u));
        changed |= duplicated;
      }
    }
  }
  return changed;
}
}  