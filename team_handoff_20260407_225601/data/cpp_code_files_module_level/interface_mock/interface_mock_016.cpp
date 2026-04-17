#ifndef AROLLA_EXPR_EVAL_CASTING_H_
#define AROLLA_EXPR_EVAL_CASTING_H_
#include "absl/status/statusor.h"
#include "arolla/expr/eval/eval.h"
#include "arolla/expr/expr_node.h"
namespace arolla::expr::eval_internal {
absl::StatusOr<ExprNodePtr> CastingTransformation(
    const DynamicEvaluationEngineOptions& options, ExprNodePtr expr);
}  
#endif  
#include "arolla/expr/eval/casting.h"
#include <cstddef>
#include <memory>
#include <optional>
#include <utility>
#include <vector>
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_format.h"
#include "absl/types/span.h"
#include "arolla/expr/derived_qtype_cast_operator.h"
#include "arolla/expr/eval/eval.h"
#include "arolla/expr/expr.h"
#include "arolla/expr/expr_debug_string.h"
#include "arolla/expr/expr_node.h"
#include "arolla/expr/expr_operator.h"
#include "arolla/expr/operators/casting_registry.h"
#include "arolla/expr/registered_expr_operator.h"
#include "arolla/qexpr/operators.h"
#include "arolla/qtype/array_like/array_like_qtype.h"
#include "arolla/qtype/derived_qtype.h"
#include "arolla/qtype/qtype.h"
#include "arolla/util/status_macros_backport.h"
namespace arolla::expr::eval_internal {
namespace {
absl::StatusOr<std::vector<QTypePtr>> GetQTypesFromNodeDeps(
    const ExprNodePtr& expr) {
  std::vector<QTypePtr> qtypes;
  qtypes.reserve(expr->node_deps().size());
  for (int i = 0; i < expr->node_deps().size(); i++) {
    const auto* qtype = expr->node_deps()[i]->qtype();
    if (qtype == nullptr) {
      return absl::InternalError(
          absl::StrFormat("QType not set for %i-th argument of node %s", i,
                          ToDebugString(expr)));
    }
    qtypes.push_back(qtype);
  }
  return qtypes;
}
absl::StatusOr<std::optional<ExprNodePtr>> GetShapeForBroadcasting(
    absl::Span<const ExprNodePtr> deps) {
  ExprNodePtr array_dep_or = nullptr;
  for (const auto& dep : deps) {
    if (IsArrayLikeQType(dep->qtype())) {
      array_dep_or = dep;
      break;
    }
  }
  if (array_dep_or != nullptr) {
    return CallOp("core.shape_of", {array_dep_or});
  }
  return std::nullopt;
}
absl::StatusOr<std::vector<ExprNodePtr>> BuildNodeDepsWithCasts(
    absl::Span<const ExprNodePtr> deps, absl::Span<const QTypePtr> dep_types,
    absl::Span<const QTypePtr> required_types) {
  const auto* casting_registry = expr_operators::CastingRegistry::GetInstance();
  ASSIGN_OR_RETURN(std::optional<ExprNodePtr> shape_for_broadcasting,
                   GetShapeForBroadcasting(deps));
  std::vector<ExprNodePtr> new_deps;
  new_deps.reserve(deps.size());
  for (size_t i = 0; i < deps.size(); ++i) {
    auto dep = deps[i];
    if (dep->qtype() != required_types[i]) {
      ASSIGN_OR_RETURN(dep,
                       casting_registry->GetCast(dep, required_types[i],
                                                 true,
                                                 shape_for_broadcasting),
                       _ << "while casting arguments "
                         << FormatTypeVector(dep_types) << " into "
                         << FormatTypeVector(required_types));
    }
    new_deps.push_back(std::move(dep));
  }
  return new_deps;
}
}  
absl::StatusOr<ExprNodePtr> CastingTransformation(
    const DynamicEvaluationEngineOptions& options, ExprNodePtr expr) {
  const OperatorDirectory& backend_operators =
      options.operator_directory != nullptr ? *options.operator_directory
                                            : *OperatorRegistry::GetInstance();
  if (!expr->is_op()) {
    return expr;
  }
  ASSIGN_OR_RETURN(auto op, DecayRegisteredOperator(expr->op()));
  if (!HasBackendExprOperatorTag(op)) {
    return expr;
  }
  auto backend_op_name = op->display_name();
  ASSIGN_OR_RETURN(auto dep_types, GetQTypesFromNodeDeps(expr));
  QTypePtr result_qtype = expr->qtype();
  if (result_qtype == nullptr) {
    return absl::InternalError(
        "all QTypes must be known before the casting compilation step");
  }
  ASSIGN_OR_RETURN(
      auto backend_op,
      backend_operators.LookupOperator(backend_op_name, dep_types,
                                       result_qtype),
      expr);
  auto* backend_op_signature = backend_op->signature();
  if (backend_op_signature->input_types() != dep_types) {
    ASSIGN_OR_RETURN(auto cast_deps, BuildNodeDepsWithCasts(
                                         expr->node_deps(), dep_types,
                                         backend_op_signature->input_types()));
    ASSIGN_OR_RETURN(expr, WithNewDependencies(expr, std::move(cast_deps)));
    if (expr->qtype() != DecayDerivedQType(result_qtype)) {
      return absl::FailedPreconditionError(absl::StrFormat(
          "expr output QType changed after input casting: was %s, became %s",
          result_qtype->name(), expr->qtype()->name()));
    }
  }
  if (backend_op_signature->output_type() == result_qtype) {
    return expr;
  }
  if (backend_op_signature->output_type() == DecayDerivedQType(result_qtype)) {
    auto downcast_op =
        std::make_shared<expr::DerivedQTypeDowncastOperator>(result_qtype);
    return CallOp(downcast_op, {expr});
  }
  return absl::FailedPreconditionError(absl::StrFormat(
      "inconsistent output types for QExpr and expr %s operator: %s",
      backend_op_name,
      JoinTypeNames({result_qtype, backend_op_signature->output_type()})));
}
}  