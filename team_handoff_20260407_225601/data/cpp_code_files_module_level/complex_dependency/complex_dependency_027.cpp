#ifndef THIRD_PARTY_CEL_CPP_EVAL_EVAL_TERNARY_STEP_H_
#define THIRD_PARTY_CEL_CPP_EVAL_EVAL_TERNARY_STEP_H_
#include <cstdint>
#include <memory>
#include "absl/status/statusor.h"
#include "eval/eval/direct_expression_step.h"
#include "eval/eval/evaluator_core.h"
namespace google::api::expr::runtime {
std::unique_ptr<DirectExpressionStep> CreateDirectTernaryStep(
    std::unique_ptr<DirectExpressionStep> condition,
    std::unique_ptr<DirectExpressionStep> left,
    std::unique_ptr<DirectExpressionStep> right, int64_t expr_id,
    bool shortcircuiting = true);
absl::StatusOr<std::unique_ptr<ExpressionStep>> CreateTernaryStep(
    int64_t expr_id);
}  
#endif  
#include "eval/eval/ternary_step.h"
#include <cstddef>
#include <cstdint>
#include <memory>
#include <utility>
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "base/builtins.h"
#include "common/casting.h"
#include "common/value.h"
#include "eval/eval/attribute_trail.h"
#include "eval/eval/direct_expression_step.h"
#include "eval/eval/evaluator_core.h"
#include "eval/eval/expression_step_base.h"
#include "eval/internal/errors.h"
#include "internal/status_macros.h"
namespace google::api::expr::runtime {
namespace {
using ::cel::BoolValue;
using ::cel::Cast;
using ::cel::ErrorValue;
using ::cel::InstanceOf;
using ::cel::UnknownValue;
using ::cel::builtin::kTernary;
using ::cel::runtime_internal::CreateNoMatchingOverloadError;
inline constexpr size_t kTernaryStepCondition = 0;
inline constexpr size_t kTernaryStepTrue = 1;
inline constexpr size_t kTernaryStepFalse = 2;
class ExhaustiveDirectTernaryStep : public DirectExpressionStep {
 public:
  ExhaustiveDirectTernaryStep(std::unique_ptr<DirectExpressionStep> condition,
                              std::unique_ptr<DirectExpressionStep> left,
                              std::unique_ptr<DirectExpressionStep> right,
                              int64_t expr_id)
      : DirectExpressionStep(expr_id),
        condition_(std::move(condition)),
        left_(std::move(left)),
        right_(std::move(right)) {}
  absl::Status Evaluate(ExecutionFrameBase& frame, cel::Value& result,
                        AttributeTrail& attribute) const override {
    cel::Value condition;
    cel::Value lhs;
    cel::Value rhs;
    AttributeTrail condition_attr;
    AttributeTrail lhs_attr;
    AttributeTrail rhs_attr;
    CEL_RETURN_IF_ERROR(condition_->Evaluate(frame, condition, condition_attr));
    CEL_RETURN_IF_ERROR(left_->Evaluate(frame, lhs, lhs_attr));
    CEL_RETURN_IF_ERROR(right_->Evaluate(frame, rhs, rhs_attr));
    if (InstanceOf<ErrorValue>(condition) ||
        InstanceOf<UnknownValue>(condition)) {
      result = std::move(condition);
      attribute = std::move(condition_attr);
      return absl::OkStatus();
    }
    if (!InstanceOf<BoolValue>(condition)) {
      result = frame.value_manager().CreateErrorValue(
          CreateNoMatchingOverloadError(kTernary));
      return absl::OkStatus();
    }
    if (Cast<BoolValue>(condition).NativeValue()) {
      result = std::move(lhs);
      attribute = std::move(lhs_attr);
    } else {
      result = std::move(rhs);
      attribute = std::move(rhs_attr);
    }
    return absl::OkStatus();
  }
 private:
  std::unique_ptr<DirectExpressionStep> condition_;
  std::unique_ptr<DirectExpressionStep> left_;
  std::unique_ptr<DirectExpressionStep> right_;
};
class ShortcircuitingDirectTernaryStep : public DirectExpressionStep {
 public:
  ShortcircuitingDirectTernaryStep(
      std::unique_ptr<DirectExpressionStep> condition,
      std::unique_ptr<DirectExpressionStep> left,
      std::unique_ptr<DirectExpressionStep> right, int64_t expr_id)
      : DirectExpressionStep(expr_id),
        condition_(std::move(condition)),
        left_(std::move(left)),
        right_(std::move(right)) {}
  absl::Status Evaluate(ExecutionFrameBase& frame, cel::Value& result,
                        AttributeTrail& attribute) const override {
    cel::Value condition;
    AttributeTrail condition_attr;
    CEL_RETURN_IF_ERROR(condition_->Evaluate(frame, condition, condition_attr));
    if (InstanceOf<ErrorValue>(condition) ||
        InstanceOf<UnknownValue>(condition)) {
      result = std::move(condition);
      attribute = std::move(condition_attr);
      return absl::OkStatus();
    }
    if (!InstanceOf<BoolValue>(condition)) {
      result = frame.value_manager().CreateErrorValue(
          CreateNoMatchingOverloadError(kTernary));
      return absl::OkStatus();
    }
    if (Cast<BoolValue>(condition).NativeValue()) {
      return left_->Evaluate(frame, result, attribute);
    }
    return right_->Evaluate(frame, result, attribute);
  }
 private:
  std::unique_ptr<DirectExpressionStep> condition_;
  std::unique_ptr<DirectExpressionStep> left_;
  std::unique_ptr<DirectExpressionStep> right_;
};
class TernaryStep : public ExpressionStepBase {
 public:
  explicit TernaryStep(int64_t expr_id) : ExpressionStepBase(expr_id) {}
  absl::Status Evaluate(ExecutionFrame* frame) const override;
};
absl::Status TernaryStep::Evaluate(ExecutionFrame* frame) const {
  if (!frame->value_stack().HasEnough(3)) {
    return absl::Status(absl::StatusCode::kInternal, "Value stack underflow");
  }
  auto args = frame->value_stack().GetSpan(3);
  const auto& condition = args[kTernaryStepCondition];
  if (frame->enable_unknowns()) {
    if (condition->Is<cel::UnknownValue>()) {
      frame->value_stack().Pop(2);
      return absl::OkStatus();
    }
  }
  if (condition->Is<cel::ErrorValue>()) {
    frame->value_stack().Pop(2);
    return absl::OkStatus();
  }
  cel::Value result;
  if (!condition->Is<cel::BoolValue>()) {
    result = frame->value_factory().CreateErrorValue(
        CreateNoMatchingOverloadError(kTernary));
  } else if (condition.As<cel::BoolValue>().NativeValue()) {
    result = args[kTernaryStepTrue];
  } else {
    result = args[kTernaryStepFalse];
  }
  frame->value_stack().PopAndPush(args.size(), std::move(result));
  return absl::OkStatus();
}
}  
std::unique_ptr<DirectExpressionStep> CreateDirectTernaryStep(
    std::unique_ptr<DirectExpressionStep> condition,
    std::unique_ptr<DirectExpressionStep> left,
    std::unique_ptr<DirectExpressionStep> right, int64_t expr_id,
    bool shortcircuiting) {
  if (shortcircuiting) {
    return std::make_unique<ShortcircuitingDirectTernaryStep>(
        std::move(condition), std::move(left), std::move(right), expr_id);
  }
  return std::make_unique<ExhaustiveDirectTernaryStep>(
      std::move(condition), std::move(left), std::move(right), expr_id);
}
absl::StatusOr<std::unique_ptr<ExpressionStep>> CreateTernaryStep(
    int64_t expr_id) {
  return std::make_unique<TernaryStep>(expr_id);
}
}  