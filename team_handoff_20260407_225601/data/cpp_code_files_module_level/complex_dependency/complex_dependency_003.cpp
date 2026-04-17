#ifndef XLA_BACKENDS_INTERPRETER_EXECUTABLE_H_
#define XLA_BACKENDS_INTERPRETER_EXECUTABLE_H_
#include <memory>
#include "absl/status/statusor.h"
#include "absl/types/span.h"
#include "xla/backends/interpreter/executable_base.h"
#include "xla/hlo/evaluator/hlo_evaluator.h"
#include "xla/hlo/ir/hlo_module.h"
#include "xla/service/executable.h"
#include "xla/service/hlo_cost_analysis.h"
#include "xla/service/hlo_execution_profile.h"
#include "xla/service/hlo_module_config.h"
#include "xla/service/service_executable_run_options.h"
#include "xla/service/shaped_buffer.h"
#include "xla/stream_executor/stream_executor.h"
#include "xla/types.h"
#include "xla/xla_data.pb.h"
namespace xla {
namespace interpreter {
class InterpreterExecutable : public InterpreterExecutableBase {
 public:
  InterpreterExecutable(
      std::unique_ptr<HloModule> hlo_module,
      std::unique_ptr<HloEvaluator> evaluator,
      std::optional<DynamicDimensionInference> dynamic_dymension_inference);
  static int64_t ShapeSizeBytes(const Shape& shape);
 protected:
  absl::StatusOr<Literal> Evaluate(
      const ServiceExecutableRunOptions* run_options,
      const HloComputation& computation,
      absl::Span<const Literal> arg_literals) override
      ABSL_LOCKS_EXCLUDED(evaluator_lock_);
  std::unique_ptr<HloEvaluator> evaluator_ ABSL_PT_GUARDED_BY(evaluator_lock_);
  mutable absl::Mutex evaluator_lock_;
 private:
  std::optional<DynamicDimensionInference> dynamic_dimension_inference_;
  InterpreterExecutable(const InterpreterExecutable&) = delete;
  InterpreterExecutable& operator=(const InterpreterExecutable&) = delete;
};
}  
}  
#endif  
#include "xla/backends/interpreter/executable.h"
#include <algorithm>
#include <cstring>
#include <memory>
#include <string>
#include <utility>
#include <vector>
#include "xla/backends/interpreter/executable_base.h"
#include "xla/backends/interpreter/executor.h"
#include "xla/hlo/ir/hlo_computation.h"
#include "xla/hlo/ir/hlo_instruction.h"
#include "xla/literal.h"
#include "xla/service/maybe_owning_device_memory.h"
#include "xla/service/transfer_manager.h"
#include "xla/shape_util.h"
#include "xla/status_macros.h"
#include "xla/stream_executor/stream_executor.h"
#include "tsl/platform/env.h"
#include "tsl/platform/errors.h"
namespace xla {
namespace interpreter {
InterpreterExecutable::InterpreterExecutable(
    std::unique_ptr<HloModule> hlo_module,
    std::unique_ptr<HloEvaluator> evaluator,
    std::optional<DynamicDimensionInference> dynamic_dymension_inference)
    : InterpreterExecutableBase(std::move(hlo_module)),
      evaluator_(std::move(evaluator)),
      dynamic_dimension_inference_(std::move(dynamic_dymension_inference)) {
  if (dynamic_dimension_inference_.has_value()) {
    evaluator_->set_dynamic_dimension_inference(
        &dynamic_dimension_inference_.value());
  }
}
absl::StatusOr<Literal> InterpreterExecutable::Evaluate(
    const ServiceExecutableRunOptions* run_options,
    const HloComputation& computation, absl::Span<const Literal> arg_literals) {
  absl::MutexLock lock(&evaluator_lock_);
  evaluator_->ResetVisitStates();
  return evaluator_->Evaluate(computation, arg_literals);
}
 int64_t InterpreterExecutable::ShapeSizeBytes(const Shape& shape) {
  if (shape.IsOpaque()) {
    return sizeof(void*);
  }
  return ShapeUtil::ByteSizeOf(shape, sizeof(void*));
}
}  
}  