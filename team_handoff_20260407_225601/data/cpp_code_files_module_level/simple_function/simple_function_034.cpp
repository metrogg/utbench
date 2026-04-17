#ifndef XLA_HLO_IR_HLO_REACHABILITY_H_
#define XLA_HLO_IR_HLO_REACHABILITY_H_
#include <memory>
#include <utility>
#include <vector>
#include "absl/algorithm/container.h"
#include "absl/container/flat_hash_map.h"
#include "absl/types/span.h"
#include "xla/hlo/ir/hlo_computation.h"
#include "xla/hlo/ir/hlo_instruction.h"
#include "xla/hlo/ir/hlo_module.h"
#include "xla/types.h"
namespace xla {
class HloReachabilityMap {
 public:
  using Index = size_t;
  explicit HloReachabilityMap(
      absl::Span<const HloInstruction* const> instructions);
  static std::unique_ptr<HloReachabilityMap> Build(
      const HloComputation* computation);
  static std::unique_ptr<HloReachabilityMap> BuildWithRestrictions(
      const HloComputation* computation,
      absl::FunctionRef<void(const HloInstruction*,
                             std::vector<HloInstruction*>*)>
          add_dependencies);
  bool SetReachabilityToUnion(absl::Span<const HloInstruction* const> inputs,
                              const HloInstruction* instruction);
  void FastSetReachabilityToUnion(
      absl::Span<const HloInstruction* const> inputs,
      const HloInstruction* instruction);
  void FastSetReachabilityToUnion(absl::Span<const Index> input_indices,
                                  Index index);
  Index GetIndex(const HloInstruction* instruction) const {
    return indices_.at(GetKey(instruction));
  }
  void SetReachable(const HloInstruction* a, const HloInstruction* b) {
    SetReachable(GetIndex(a), GetIndex(b));
  }
  void SetReachable(Index a, Index b) { bit_sets_[b].Set(a); }
  void UpdateReachabilityThroughInstruction(const HloInstruction* instruction);
  bool IsReachable(const HloInstruction* a, const HloInstruction* b) const {
    return IsReachable(GetIndex(a), GetIndex(b));
  }
  bool IsReachable(Index a, Index b) const { return bit_sets_[b].Get(a); }
  bool IsConnected(const HloInstruction* a, const HloInstruction* b) const {
    return IsConnected(GetIndex(a), GetIndex(b));
  }
  bool IsConnected(Index a, Index b) const {
    return IsReachable(a, b) || IsReachable(b, a);
  }
  bool IsPresent(const HloInstruction* instruction) const {
    return indices_.contains(GetKey(instruction));
  }
  void Replace(const HloInstruction* original,
               const HloInstruction* replacement);
 private:
  class BitSet {
   public:
    BitSet() = default;
    explicit BitSet(size_t size)
        : size_(size), vector_((size + kBits - 1) / kBits, 0) {}
    bool Get(Index index) const {
      DCHECK(index >= 0 && index < size_);
      return vector_[index / kBits] & (1ull << (index % kBits));
    }
    void Set(Index index) {
      DCHECK(index >= 0 && index < size_);
      vector_[index / kBits] |= 1ull << (index % kBits);
    }
    void operator|=(const BitSet& other) {
      if (this == &other) return;
      DCHECK(size_ == other.size_);
      const Word* a = vector_.data();
      const Word* b = other.vector_.data();
      Word* __restrict out = vector_.data();
      size_t num_words = vector_.size();
      for (size_t i = 0; i < num_words; ++i) {
        out[i] = a[i] | b[i];
      }
    }
    void SetToZero() { absl::c_fill(vector_, 0); }
    bool operator==(const BitSet& other) const {
      return vector_ == other.vector_;
    }
    bool operator!=(const BitSet& other) const { return !(*this == other); }
   private:
    using Word = uint64_t;
    static constexpr size_t kBits = 64;
    size_t size_;  
    std::vector<Word> vector_;
  };
  friend class HloReachabilityMapBitSetBenchmark;
  using Key = std::pair<int, int>;  
  static Key GetKey(const HloInstruction* instruction) {
    return {instruction->GetModule()->unique_id(), instruction->unique_id()};
  }
  void SetReachabilityToUnionHelper(
      absl::Span<const HloInstruction* const> inputs, Index index);
  void SetReachabilityToUnionHelper(absl::Span<const Index> input_indices,
                                    Index index);
  absl::flat_hash_map<Key, Index> indices_;
  std::vector<BitSet> bit_sets_;
  BitSet tmp_bit_set_;
};
}  
#endif  
#include "xla/hlo/ir/hlo_reachability.h"
#include <memory>
#include <queue>
#include <vector>
#include "absl/algorithm/container.h"
#include "xla/hlo/ir/hlo_instruction.h"
namespace xla {
HloReachabilityMap::HloReachabilityMap(
    absl::Span<const HloInstruction* const> instructions)
    : bit_sets_(instructions.size(), BitSet(instructions.size())) {
  indices_.reserve(instructions.size());
  for (size_t i = 0; i < instructions.size(); ++i) {
    bit_sets_[i].Set(i);  
    indices_[GetKey(instructions[i])] = i;
  }
}
bool HloReachabilityMap::SetReachabilityToUnion(
    absl::Span<const HloInstruction* const> inputs,
    const HloInstruction* instruction) {
  Index index = GetIndex(instruction);
  BitSet& bit_set = bit_sets_[index];
  tmp_bit_set_ = bit_set;
  SetReachabilityToUnionHelper(inputs, index);
  return bit_set != tmp_bit_set_;
}
void HloReachabilityMap::FastSetReachabilityToUnion(
    absl::Span<const HloInstruction* const> inputs,
    const HloInstruction* instruction) {
  SetReachabilityToUnionHelper(inputs, GetIndex(instruction));
}
void HloReachabilityMap::FastSetReachabilityToUnion(
    absl::Span<const Index> input_indices, Index index) {
  SetReachabilityToUnionHelper(input_indices, index);
}
void HloReachabilityMap::SetReachabilityToUnionHelper(
    absl::Span<const HloInstruction* const> inputs, Index index) {
  absl::InlinedVector<Index, 16> input_indices;
  input_indices.reserve(inputs.size());
  for (const HloInstruction* input : inputs) {
    input_indices.push_back(GetIndex(input));
  }
  SetReachabilityToUnionHelper(input_indices, index);
}
void HloReachabilityMap::SetReachabilityToUnionHelper(
    absl::Span<const Index> input_indices, Index index) {
  BitSet& bit_set = bit_sets_[index];
  if (!absl::c_linear_search(input_indices, index)) {
    bit_set.SetToZero();
  }
  bit_set.Set(index);
  for (Index input_index : input_indices) {
    if (input_index != index) {
      bit_set |= bit_sets_[input_index];
    }
  }
}
void HloReachabilityMap::Replace(const HloInstruction* original,
                                 const HloInstruction* replacement) {
  if (GetKey(original) != GetKey(replacement)) {
    indices_[GetKey(replacement)] = GetIndex(original);
    indices_.erase(GetKey(original));
  }
}
std::unique_ptr<HloReachabilityMap> HloReachabilityMap::BuildWithRestrictions(
    const HloComputation* computation,
    absl::FunctionRef<void(const HloInstruction*,
                           std::vector<HloInstruction*>*)>
        add_dependencies) {
  const auto& all = computation->MakeInstructionPostOrder();
  auto result = std::make_unique<HloReachabilityMap>(all);
  std::vector<HloInstruction*> inputs;
  for (const HloInstruction* hlo : all) {
    inputs.clear();
    add_dependencies(hlo, &inputs);
    result->FastSetReachabilityToUnion(inputs, hlo);
  }
  return result;
}
std::unique_ptr<HloReachabilityMap> HloReachabilityMap::Build(
    const HloComputation* computation) {
  HloComputation::ChannelDependencies channel_dependencies =
      computation->ComputeChannelDependencies();
  std::vector<HloInstruction*> instructions =
      computation->MakeInstructionPostOrder(channel_dependencies);
  auto result = std::make_unique<HloReachabilityMap>(instructions);
  auto get_bit_set = [&](const HloInstruction* instruction) -> BitSet& {
    return result->bit_sets_[result->GetIndex(instruction)];
  };
  for (const HloInstruction* instruction : instructions) {
    BitSet& bit_set = get_bit_set(instruction);
    auto add_dependencies = [&](const HloInstruction* instruction) {
      for (const HloInstruction* operand : instruction->operands()) {
        bit_set |= get_bit_set(operand);
      }
      for (const HloInstruction* predecessor :
           instruction->control_predecessors()) {
        bit_set |= get_bit_set(predecessor);
      }
    };
    add_dependencies(instruction);
    auto it = channel_dependencies.find(instruction);
    if (it != channel_dependencies.end()) {
      absl::c_for_each(it->second, add_dependencies);
    }
  }
  return result;
}
void HloReachabilityMap::UpdateReachabilityThroughInstruction(
    const HloInstruction* instruction) {
  std::queue<const HloInstruction*> worklist;
  worklist.push(instruction);
  std::vector<HloInstruction*> inputs;
  while (!worklist.empty()) {
    const HloInstruction* item = worklist.front();
    worklist.pop();
    inputs.assign(item->operands().begin(), item->operands().end());
    inputs.insert(inputs.end(), item->control_predecessors().begin(),
                  item->control_predecessors().end());
    if (SetReachabilityToUnion(inputs, item)) {
      for (const HloInstruction* user : item->users()) {
        worklist.push(user);
      }
      for (const HloInstruction* succ : item->control_successors()) {
        worklist.push(succ);
      }
    }
  }
}
}  