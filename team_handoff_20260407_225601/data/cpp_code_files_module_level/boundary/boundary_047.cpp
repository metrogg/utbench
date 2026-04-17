#ifndef XLA_SERVICE_COLLECTIVE_PERMUTE_DECOMPOSER_H_
#define XLA_SERVICE_COLLECTIVE_PERMUTE_DECOMPOSER_H_
#include "xla/hlo/ir/hlo_module.h"
#include "xla/service/hlo_pass_interface.h"
namespace xla {
class CollectivePermuteDecomposer : public HloModulePass {
 public:
  explicit CollectivePermuteDecomposer(int64_t threshold_in_bytes)
      : threshold_in_bytes_(threshold_in_bytes) {}
  absl::string_view name() const override {
    return "collective-permute-decomposer";
  }
  using HloPassInterface::Run;
  absl::StatusOr<bool> Run(
      HloModule* module,
      const absl::flat_hash_set<absl::string_view>& execution_threads) override;
 private:
  int64_t threshold_in_bytes_;
};
}  
#endif  
#include "xla/service/collective_permute_decomposer.h"
#include <cstdint>
#include <optional>
#include <string>
#include <utility>
#include <vector>
#include "absl/status/status.h"
#include "absl/strings/str_join.h"
#include "xla/hlo/ir/hlo_casting_utils.h"
#include "xla/hlo/ir/hlo_computation.h"
#include "xla/hlo/ir/hlo_instruction.h"
#include "xla/hlo/ir/hlo_instructions.h"
#include "xla/hlo/ir/hlo_opcode.h"
#include "xla/service/collective_ops_utils.h"
#include "xla/service/gpu/backend_configs.pb.h"
#include "xla/service/graphcycles/graphcycles.h"
#include "xla/shape_util.h"
#include "xla/xla_data.pb.h"
#include "tsl/platform/errors.h"
namespace xla {
namespace {
using SourceTargetPair = std::pair<int64_t, int64_t>;
using SourceTargetPairs = std::vector<SourceTargetPair>;
bool HasCycles(const SourceTargetPairs& pairs) {
  tensorflow::GraphCycles graph;
  absl::flat_hash_map<int64_t, int32_t> replica_to_node_id;
  auto get_node_id = [&](int64_t replica) {
    auto it_and_inserted = replica_to_node_id.emplace(replica, -1);
    auto it = it_and_inserted.first;
    auto inserted = it_and_inserted.second;
    if (inserted) {
      it->second = graph.NewNode();
    }
    return it->second;
  };
  for (auto pair : pairs) {
    auto source = get_node_id(pair.first);
    auto target = get_node_id(pair.second);
    VLOG(3) << "See source " << source << " -> target " << target;
    if (!graph.InsertEdge(source, target)) {
      VLOG(3) << "Detected cycles";
      return true;
    }
  }
  return false;
}
bool ShouldDecompose(const HloCollectivePermuteInstruction& collective_permute,
                     int64_t threshold_in_bytes) {
  if (!collective_permute.channel_id().has_value()) {
    return false;
  }
  const Shape& result_shape = collective_permute.shape();
  if (!result_shape.IsArray()) {
    return false;
  }
  if (ShapeUtil::ByteSizeOf(result_shape) < threshold_in_bytes) {
    return false;
  }
  return !HasCycles(collective_permute.source_target_pairs());
}
bool MayPipeline(const HloCollectivePermuteInstruction& collective_permute) {
  const HloInstruction* data = collective_permute.operand(0);
  return (data->opcode() == HloOpcode::kGetTupleElement &&
          data->operand(0)->opcode() == HloOpcode::kParameter);
}
absl::Status DecomposeCollectivePermute(
    HloCollectivePermuteInstruction* collective_permute,
    HloComputation* computation, const std::string& pipeline_decision) {
  int64_t channel_id = collective_permute->channel_id().value();
  HloInstruction* data = collective_permute->mutable_operand(0);
  const Shape& data_shape = data->shape();
  const OpMetadata& metadata = collective_permute->metadata();
  const xla::FrontendAttributes& old_attributes =
      collective_permute->frontend_attributes();
  xla::FrontendAttributes attributes;
  std::string source_target_pairs_string =
      "{" +
      absl::StrJoin(collective_permute->source_target_pairs(), ",",
                    absl::PairFormatter(
                        [](std::string* out, int64_t value) {
                          absl::StrAppend(out, "{", value);
                        },
                        ",",
                        [](std::string* out, int64_t value) {
                          absl::StrAppend(out, value, "}");
                        })) +
      "}";
  attributes.mutable_map()->insert(old_attributes.map().begin(),
                                   old_attributes.map().end());
  (*attributes.mutable_map())[kSendRecvSourceTargetPairsAttr] =
      source_target_pairs_string;
  HloInstruction* after_all =
      computation->AddInstruction(HloInstruction::CreateToken());
  HloInstruction* recv = computation->AddInstruction(
      HloInstruction::CreateRecv(data_shape, after_all, channel_id));
  recv->add_frontend_attributes(attributes);
  recv->set_metadata(metadata);
  HloInstruction* send = computation->AddInstruction(
      HloInstruction::CreateSend(data, after_all, channel_id));
  send->add_frontend_attributes(attributes);
  send->set_metadata(metadata);
  HloInstruction* recv_done =
      computation->AddInstruction(HloInstruction::CreateRecvDone(recv));
  HloInstruction* send_done =
      computation->AddInstruction(HloInstruction::CreateSendDone(send));
  TF_RETURN_IF_ERROR(send->AddControlDependencyTo(recv_done));
  HloInstruction* recv_data = computation->AddInstruction(
      HloInstruction::CreateGetTupleElement(recv_done, 0));
  TF_RETURN_IF_ERROR(collective_permute->ReplaceAllUsesWith(recv_data));
  TF_RETURN_IF_ERROR(
      computation->RemoveInstructionAndUnusedOperands(collective_permute));
  if (!pipeline_decision.empty()) {
    xla::FrontendAttributes attributes;
    (*attributes.mutable_map())[kSendRecvPipelineAttr] = pipeline_decision;
    send->add_frontend_attributes(attributes);
    send_done->add_frontend_attributes(attributes);
    recv->add_frontend_attributes(attributes);
    recv_done->add_frontend_attributes(attributes);
  }
  return absl::OkStatus();
}
bool IsForwardCycle(const SourceTargetPair& backedge,
                    const SourceTargetPairs& others) {
  int64_t num_pairs = others.size() + 1;
  if (backedge.first != num_pairs - 1 || backedge.second != 0) {
    return false;
  }
  for (int64_t i = 0; i < num_pairs - 1; ++i) {
    const SourceTargetPair& pair = others[i];
    if (pair.first != i || pair.second != i + 1) {
      return false;
    }
  }
  return true;
}
bool IsBackwardCycle(const SourceTargetPair& backedge,
                     const SourceTargetPairs& others) {
  int64_t num_pairs = others.size() + 1;
  if (backedge.first != 0 || backedge.second != num_pairs - 1) {
    return false;
  }
  for (int64_t i = 0; i < num_pairs - 1; ++i) {
    const SourceTargetPair& pair = others[i];
    if (pair.first != i + 1 || pair.second != i) {
      return false;
    }
  }
  return true;
}
std::optional<std::pair<HloCollectivePermuteInstruction*,
                        HloCollectivePermuteInstruction*>>
CheckCyclePatterns(HloCollectivePermuteInstruction* cp0,
                   HloCollectivePermuteInstruction* cp1) {
  const SourceTargetPairs& cp0_pairs = cp0->source_target_pairs();
  const SourceTargetPairs& cp1_pairs = cp1->source_target_pairs();
  if (cp0_pairs.size() == 1) {
    if (IsForwardCycle(cp0_pairs.front(), cp1_pairs) ||
        IsBackwardCycle(cp0_pairs.front(), cp1_pairs)) {
      return std::make_pair(cp0, cp1);
    }
  }
  if (cp1_pairs.size() == 1) {
    if (IsForwardCycle(cp1_pairs.front(), cp0_pairs) ||
        IsBackwardCycle(cp1_pairs.front(), cp0_pairs)) {
      return std::make_pair(cp1, cp0);
    }
  }
  return std::nullopt;
}
}  
absl::StatusOr<bool> CollectivePermuteDecomposer::Run(
    HloModule* module,
    const absl::flat_hash_set<absl::string_view>& execution_threads) {
  bool changed = false;
  std::vector<HloComputation*> all_computations =
      module->MakeComputationPostOrder(execution_threads);
  absl::flat_hash_set<HloComputation*> while_bodies;
  for (auto iter = all_computations.rbegin(); iter != all_computations.rend();
       ++iter) {
    HloComputation* computation = *iter;
    bool may_pipeline = while_bodies.contains(computation);
    std::vector<HloCollectivePermuteInstruction*> cps_to_decompose;
    HloCollectivePermuteInstruction* cp0_to_pipeline = nullptr;
    HloCollectivePermuteInstruction* cp1_to_pipeline = nullptr;
    for (HloInstruction* hlo : computation->MakeInstructionPostOrder()) {
      if (hlo->opcode() == HloOpcode::kWhile) {
        while_bodies.insert(hlo->while_body());
        continue;
      }
      if (hlo->opcode() != HloOpcode::kCollectivePermute) {
        continue;
      }
      HloCollectivePermuteInstruction* cp =
          Cast<HloCollectivePermuteInstruction>(hlo);
      if (!ShouldDecompose(*cp, threshold_in_bytes_)) {
        continue;
      }
      cps_to_decompose.push_back(cp);
      if (!while_bodies.contains(computation) || !may_pipeline) {
        continue;
      }
      if (cp0_to_pipeline != nullptr && cp1_to_pipeline != nullptr) {
        continue;
      }
      if (!MayPipeline(*cp)) {
        continue;
      }
      if (cp0_to_pipeline == nullptr) {
        cp0_to_pipeline = cp;
        continue;
      }
      auto optional_pair = CheckCyclePatterns(cp0_to_pipeline, cp);
      if (optional_pair.has_value()) {
        cp0_to_pipeline = optional_pair.value().first;
        cp1_to_pipeline = optional_pair.value().second;
      }
    }
    for (HloCollectivePermuteInstruction* cp : cps_to_decompose) {
      std::string pipeline_decision;
      if (cp0_to_pipeline == cp) {
        pipeline_decision = "0";
      } else if (cp1_to_pipeline == cp) {
        pipeline_decision = "1";
      }
      TF_RETURN_IF_ERROR(
          DecomposeCollectivePermute(cp, computation, pipeline_decision));
    }
    if (!cps_to_decompose.empty()) {
      changed = true;
    }
  }
  return changed;
}
}  