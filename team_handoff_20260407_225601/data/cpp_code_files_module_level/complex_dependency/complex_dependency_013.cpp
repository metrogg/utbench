#ifndef XLA_SERVICE_GPU_FUSION_PROCESS_DUMP_H_
#define XLA_SERVICE_GPU_FUSION_PROCESS_DUMP_H_
#include <cstdint>
#include <memory>
#include <string>
#include <utility>
#include "absl/container/flat_hash_map.h"
#include "absl/container/inlined_vector.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "xla/hlo/ir/hlo_instruction.h"
#include "xla/hlo/ir/hlo_module.h"
#include "xla/service/gpu/fusion_process_dump.pb.h"
#include "xla/stream_executor/device_description.h"
namespace xla {
namespace gpu {
class FusionProcessDump {
 public:
  static absl::StatusOr<FusionProcessDump> LoadFromFile(
      const std::string& path);
  static absl::StatusOr<FusionProcessDump> LoadFromData(
      const std::string& data, absl::string_view format);
  static absl::StatusOr<FusionProcessDump> LoadFromProto(
      const FusionProcessDumpProto& fusion_process_dump_proto);
  const FusionProcessDumpProto& proto() { return fusion_process_dump_proto_; }
  HloModule* module() { return hlo_module_.get(); }
  const se::DeviceDescription& device_info() { return device_info_; }
  int64_t current_step_idx() { return current_step_idx_; }
  HloComputation* GetCurrentComputation();
  HloInstruction* GetInstructionWithName(absl::string_view name);
  HloInstruction* GetProducer();
  absl::InlinedVector<HloInstruction*, 2> GetConsumers();
  HloInstruction* GetLastFusion() { return last_fusion_; }
  const FusionStep& CurrentStep();
  bool HasNext();
  void Advance();
 private:
  FusionProcessDump(FusionProcessDumpProto fusion_process_dump_proto,
                    std::unique_ptr<HloModule> hlo_module,
                    se::DeviceDescription device_info,
                    absl::flat_hash_map<std::string, HloComputation*>
                        instruction_name_to_computation_map)
      : fusion_process_dump_proto_(std::move(fusion_process_dump_proto)),
        hlo_module_(std::move(hlo_module)),
        device_info_(std::move(device_info)),
        instruction_name_to_computation_map_(
            std::move(instruction_name_to_computation_map)) {}
  FusionProcessDumpProto fusion_process_dump_proto_;
  std::unique_ptr<HloModule> hlo_module_;
  se::DeviceDescription device_info_;
  absl::flat_hash_map<std::string, HloComputation*>
      instruction_name_to_computation_map_;
  int64_t current_step_idx_ = 0;
  HloInstruction* last_fusion_ = nullptr;
};
}  
}  
#endif  
#include "xla/service/gpu/fusion_process_dump.h"
#include <string>
#include <string_view>
#include <utility>
#include "absl/container/flat_hash_map.h"
#include "absl/container/inlined_vector.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "xla/hlo/ir/hlo_computation.h"
#include "xla/hlo/ir/hlo_instruction.h"
#include "xla/hlo/ir/hlo_opcode.h"
#include "xla/service/gpu/fusion_process_dump.pb.h"
#include "xla/stream_executor/device_description.h"
#include "xla/tools/hlo_module_loader.h"
#include "xla/util.h"
#include "tsl/platform/env.h"
#include "tsl/platform/errors.h"
#include "tsl/platform/logging.h"
#include "tsl/platform/path.h"
#include "tsl/platform/protobuf.h"  
#include "tsl/platform/status.h"
#include "tsl/platform/statusor.h"
namespace xla {
namespace gpu {
namespace {
HloInstruction* AddFusionInstruction(HloInstruction* producer,
                                     HloInstruction* consumer,
                                     HloComputation* computation,
                                     std::string_view fusion_name) {
  if (consumer->opcode() == HloOpcode::kFusion) {
    return consumer;
  }
  auto kind = HloInstruction::FusionKind::kLoop;
  auto fusion_instruction = computation->AddInstruction(
      HloInstruction::CreateFusion(consumer->shape(), kind, consumer),
      fusion_name);
  TF_CHECK_OK(computation->ReplaceInstruction(consumer, fusion_instruction));
  return fusion_instruction;
}
HloInstruction* Fuse(HloInstruction* producer, HloInstruction* consumer,
                     HloComputation* computation,
                     std::string_view fusion_name) {
  HloInstruction* fusion_instruction =
      AddFusionInstruction(producer, consumer, computation, fusion_name);
  if (producer->opcode() == HloOpcode::kFusion) {
    fusion_instruction->MergeFusionInstruction(producer);
  } else {
    fusion_instruction->FuseInstruction(producer);
  }
  if (producer->user_count() == 0) {
    TF_CHECK_OK(computation->RemoveInstruction(producer));
  }
  return fusion_instruction;
}
absl::string_view GetProducerName(const FusionStep& step) {
  if (step.has_fusion()) {
    return step.fusion().producer_name();
  }
  if (step.has_update_priority()) {
    return step.update_priority().producer_name();
  }
  if (step.has_producer_ineligible()) {
    return step.producer_ineligible().producer_name();
  }
  LOG(FATAL) << "Producer name not found in the current step.";
}
}  
absl::StatusOr<FusionProcessDump> FusionProcessDump::LoadFromFile(
    const std::string& path) {
  std::string format = std::string(tsl::io::Extension(path));
  std::string data;
  TF_RETURN_IF_ERROR(tsl::ReadFileToString(tsl::Env::Default(), path, &data));
  return FusionProcessDump::LoadFromData(data, format);
}
absl::StatusOr<FusionProcessDump> FusionProcessDump::LoadFromData(
    const std::string& data, absl::string_view format) {
  FusionProcessDumpProto fusion_process_dump_proto;
  if (format == "txt" || format == "pbtxt") {
    if (!tsl::protobuf::TextFormat::ParseFromString(
            data, &fusion_process_dump_proto)) {
      return InvalidArgument("Failed to parse input as HLO protobuf text");
    }
  } else if (format == "pb") {
    if (!fusion_process_dump_proto.ParseFromString(data)) {
      return InvalidArgument("Failed to parse input as HLO protobuf binary");
    }
  } else {
    return InvalidArgument(
        "Invalid format from file extension: '%s'. Expected: txt, pb, or pbtxt",
        format);
  }
  return FusionProcessDump::LoadFromProto(fusion_process_dump_proto);
}
absl::StatusOr<FusionProcessDump> FusionProcessDump::LoadFromProto(
    const FusionProcessDumpProto& fusion_process_dump_proto) {
  TF_ASSIGN_OR_RETURN(
      auto module,
      LoadModuleFromData(fusion_process_dump_proto.hlo_module_before_fusion(),
                         "txt"));
  se::DeviceDescription gpu_device_info(
      fusion_process_dump_proto.gpu_device_info());
  absl::flat_hash_map<std::string, HloComputation*>
      instruction_name_to_computation_map;
  for (HloComputation* computation : module->MakeNonfusionComputations()) {
    for (HloInstruction* instr : computation->instructions()) {
      instruction_name_to_computation_map[instr->name()] = computation;
    }
  }
  return FusionProcessDump(std::move(fusion_process_dump_proto),
                           std::move(module), std::move(gpu_device_info),
                           std::move(instruction_name_to_computation_map));
}
HloComputation* FusionProcessDump::GetCurrentComputation() {
  return instruction_name_to_computation_map_.at(
      GetProducerName(CurrentStep()));
}
HloInstruction* FusionProcessDump::GetInstructionWithName(
    absl::string_view name) {
  return instruction_name_to_computation_map_[name]->GetInstructionWithName(
      name);
}
HloInstruction* FusionProcessDump::GetProducer() {
  return GetInstructionWithName(GetProducerName(CurrentStep()));
}
absl::InlinedVector<HloInstruction*, 2> FusionProcessDump::GetConsumers() {
  auto& step = CurrentStep();
  if (step.has_fusion()) {
    return {GetInstructionWithName(step.fusion().consumer_name())};
  }
  if (step.has_update_priority()) {
    absl::InlinedVector<HloInstruction*, 2> consumers;
    for (const auto& consumer_name : step.update_priority().consumer_names()) {
      consumers.push_back(GetInstructionWithName(consumer_name));
    }
    return consumers;
  }
  return {};
}
const FusionStep& FusionProcessDump::CurrentStep() {
  CHECK(HasNext());
  return fusion_process_dump_proto_.fusion_steps(current_step_idx_);
}
bool FusionProcessDump::HasNext() {
  return current_step_idx_ < fusion_process_dump_proto_.fusion_steps_size();
}
void FusionProcessDump::Advance() {
  auto step = CurrentStep();
  if (step.has_fusion()) {
    const auto& fusion_step = step.fusion();
    auto* computation = GetCurrentComputation();
    HloInstruction* producer =
        computation->GetInstructionWithName(fusion_step.producer_name());
    HloInstruction* consumer =
        computation->GetInstructionWithName(fusion_step.consumer_name());
    HloInstruction* fusion =
        Fuse(producer, consumer, computation, fusion_step.fusion_name());
    instruction_name_to_computation_map_[fusion->name()] = computation;
    last_fusion_ = fusion;
  }
  ++current_step_idx_;
}
}  
}  