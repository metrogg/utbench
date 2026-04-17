#ifndef XLA_SERVICE_FLOAT_SUPPORT_H_
#define XLA_SERVICE_FLOAT_SUPPORT_H_
#include "xla/hlo/ir/hlo_instruction.h"
#include "xla/hlo/ir/hlo_opcode.h"
#include "xla/xla_data.pb.h"
namespace xla {
class FloatSupport {
 public:
  explicit FloatSupport(PrimitiveType low_precision_type,
                        PrimitiveType high_precision_type = F32)
      : low_precision_type_(low_precision_type),
        high_precision_type_(high_precision_type) {}
  virtual ~FloatSupport() = default;
  PrimitiveType LowPrecisionType() const { return low_precision_type_; }
  PrimitiveType HighPrecisionType() const { return high_precision_type_; }
  virtual bool SupportsLowPrecisionOperand(const HloInstruction& hlo,
                                           int64_t operand_index) const;
  virtual bool SupportsLowPrecisionOutput(const HloInstruction& hlo) const;
  virtual bool SupportsMixedPrecisions(const HloInstruction& hlo) const;
  static bool EffectiveOperandPrecisionIsOutputPrecision(
      const HloInstruction& hlo, int64_t operand_index);
  virtual bool EffectiveOperandPrecisionIsLowPrecision(
      const HloInstruction& hlo, int64_t operand_index) const;
 private:
  PrimitiveType low_precision_type_;
  PrimitiveType high_precision_type_;
};
}  
#endif  
#include "xla/service/float_support.h"
#include "xla/hlo/ir/hlo_computation.h"
#include "xla/hlo/ir/hlo_instruction.h"
#include "xla/hlo/ir/hlo_opcode.h"
namespace xla {
bool FloatSupport::SupportsLowPrecisionOperand(const HloInstruction& hlo,
                                               int64_t operand_index) const {
  switch (hlo.opcode()) {
    case HloOpcode::kCall:
    case HloOpcode::kConditional:
    case HloOpcode::kCustomCall:
    case HloOpcode::kDomain:
    case HloOpcode::kGetTupleElement:
    case HloOpcode::kTuple:
    case HloOpcode::kWhile:
    case HloOpcode::kOptimizationBarrier:
      return true;
    case HloOpcode::kConvert:
      CHECK_EQ(operand_index, 0);
      return hlo.operand(0)->shape().element_type() == low_precision_type_;
    default:
      break;
  }
  return false;
}
bool FloatSupport::SupportsLowPrecisionOutput(const HloInstruction& hlo) const {
  switch (hlo.opcode()) {
    case HloOpcode::kCall:
    case HloOpcode::kConditional:
    case HloOpcode::kCustomCall:
    case HloOpcode::kDomain:
    case HloOpcode::kGetTupleElement:
    case HloOpcode::kTuple:
    case HloOpcode::kWhile:
    case HloOpcode::kOptimizationBarrier:
      return true;
    case HloOpcode::kConvert:
      return hlo.shape().element_type() == low_precision_type_;
    default:
      break;
  }
  return false;
}
bool FloatSupport::SupportsMixedPrecisions(const HloInstruction& hlo) const {
  switch (hlo.opcode()) {
    case HloOpcode::kCall:
    case HloOpcode::kConditional:
    case HloOpcode::kConvert:
    case HloOpcode::kCustomCall:
    case HloOpcode::kGetTupleElement:
    case HloOpcode::kTuple:
    case HloOpcode::kWhile:
    case HloOpcode::kOptimizationBarrier:
      return true;
    default:
      break;
  }
  return false;
}
bool FloatSupport::EffectiveOperandPrecisionIsOutputPrecision(
    const HloInstruction& hlo, int64_t operand_index) {
  switch (hlo.opcode()) {
    case HloOpcode::kAbs:
    case HloOpcode::kAllGather:
    case HloOpcode::kAllToAll:
    case HloOpcode::kBroadcast:
    case HloOpcode::kClamp:
    case HloOpcode::kCollectiveBroadcast:
    case HloOpcode::kCollectivePermute:
    case HloOpcode::kConcatenate:
    case HloOpcode::kConvert:
    case HloOpcode::kCopy:
    case HloOpcode::kDomain:
    case HloOpcode::kGetTupleElement:
    case HloOpcode::kMaximum:
    case HloOpcode::kMinimum:
    case HloOpcode::kPad:
    case HloOpcode::kReshape:
    case HloOpcode::kReverse:
    case HloOpcode::kSlice:
    case HloOpcode::kSort:
    case HloOpcode::kTranspose:
    case HloOpcode::kTuple:
    case HloOpcode::kOptimizationBarrier:
      return true;
    case HloOpcode::kBitcast:
      return hlo.shape().element_type() ==
             hlo.operand(0)->shape().element_type();
    case HloOpcode::kDynamicSlice:
      return operand_index == 0;
    case HloOpcode::kDynamicUpdateSlice:
      return operand_index == 0 || operand_index == 1;
    case HloOpcode::kGather:
      return operand_index == 0;
    case HloOpcode::kSelect:
      return operand_index == 1 || operand_index == 2;
    case HloOpcode::kReduce:
    case HloOpcode::kReduceWindow: {
      HloComputation* reduce_comp = hlo.called_computations()[0];
      for (HloInstruction* inst : reduce_comp->instructions()) {
        if (inst->opcode() == HloOpcode::kParameter) {
          continue;
        }
        for (int64_t i = 0; i < inst->operand_count(); ++i) {
          if (!EffectiveOperandPrecisionIsOutputPrecision(*inst, i)) {
            return false;
          }
        }
      }
      return true;
    }
    default:
      break;
  }
  return false;
}
bool FloatSupport::EffectiveOperandPrecisionIsLowPrecision(
    const HloInstruction& hlo, int64_t operand_index) const {
  return false;
}
}  