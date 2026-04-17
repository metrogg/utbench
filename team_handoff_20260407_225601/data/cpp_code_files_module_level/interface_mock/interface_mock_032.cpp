#ifndef THIRD_PARTY_CEL_CPP_EVAL_COMPILER_INSTRUMENTATION_H_
#define THIRD_PARTY_CEL_CPP_EVAL_COMPILER_INSTRUMENTATION_H_
#include <cstdint>
#include <functional>
#include "absl/functional/any_invocable.h"
#include "absl/status/status.h"
#include "base/ast_internal/ast_impl.h"
#include "common/value.h"
#include "eval/compiler/flat_expr_builder_extensions.h"
namespace google::api::expr::runtime {
using Instrumentation =
    std::function<absl::Status(int64_t expr_id, const cel::Value&)>;
using InstrumentationFactory = absl::AnyInvocable<Instrumentation(
    const cel::ast_internal::AstImpl&) const>;
ProgramOptimizerFactory CreateInstrumentationExtension(
    InstrumentationFactory factory);
}  
#endif  
#include "eval/compiler/instrumentation.h"
#include <cstdint>
#include <memory>
#include <utility>
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "base/ast_internal/ast_impl.h"
#include "base/ast_internal/expr.h"
#include "eval/compiler/flat_expr_builder_extensions.h"
#include "eval/eval/evaluator_core.h"
#include "eval/eval/expression_step_base.h"
namespace google::api::expr::runtime {
namespace {
class InstrumentStep : public ExpressionStepBase {
 public:
  explicit InstrumentStep(int64_t expr_id, Instrumentation instrumentation)
      : ExpressionStepBase(expr_id, false),
        expr_id_(expr_id),
        instrumentation_(std::move(instrumentation)) {}
  absl::Status Evaluate(ExecutionFrame* frame) const override {
    if (!frame->value_stack().HasEnough(1)) {
      return absl::InternalError("stack underflow in instrument step.");
    }
    return instrumentation_(expr_id_, frame->value_stack().Peek());
    return absl::OkStatus();
  }
 private:
  int64_t expr_id_;
  Instrumentation instrumentation_;
};
class InstrumentOptimizer : public ProgramOptimizer {
 public:
  explicit InstrumentOptimizer(Instrumentation instrumentation)
      : instrumentation_(std::move(instrumentation)) {}
  absl::Status OnPreVisit(PlannerContext& context,
                          const cel::ast_internal::Expr& node) override {
    return absl::OkStatus();
  }
  absl::Status OnPostVisit(PlannerContext& context,
                           const cel::ast_internal::Expr& node) override {
    if (context.GetSubplan(node).empty()) {
      return absl::OkStatus();
    }
    return context.AddSubplanStep(
        node, std::make_unique<InstrumentStep>(node.id(), instrumentation_));
  }
 private:
  Instrumentation instrumentation_;
};
}  
ProgramOptimizerFactory CreateInstrumentationExtension(
    InstrumentationFactory factory) {
  return [fac = std::move(factory)](PlannerContext&,
                                    const cel::ast_internal::AstImpl& ast)
             -> absl::StatusOr<std::unique_ptr<ProgramOptimizer>> {
    Instrumentation ins = fac(ast);
    if (ins) {
      return std::make_unique<InstrumentOptimizer>(std::move(ins));
    }
    return nullptr;
  };
}
}  