#ifndef THIRD_PARTY_CEL_CPP_EVAL_COMPILER_CEL_EXPRESSION_BUILDER_FLAT_IMPL_H_
#define THIRD_PARTY_CEL_CPP_EVAL_COMPILER_CEL_EXPRESSION_BUILDER_FLAT_IMPL_H_
#include <memory>
#include <string>
#include <utility>
#include <vector>
#include "google/api/expr/v1alpha1/checked.pb.h"
#include "google/api/expr/v1alpha1/syntax.pb.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "base/ast.h"
#include "eval/compiler/flat_expr_builder.h"
#include "eval/public/cel_expression.h"
#include "runtime/runtime_options.h"
namespace google::api::expr::runtime {
class CelExpressionBuilderFlatImpl : public CelExpressionBuilder {
 public:
  explicit CelExpressionBuilderFlatImpl(const cel::RuntimeOptions& options)
      : flat_expr_builder_(GetRegistry()->InternalGetRegistry(),
                           *GetTypeRegistry(), options) {}
  CelExpressionBuilderFlatImpl()
      : flat_expr_builder_(GetRegistry()->InternalGetRegistry(),
                           *GetTypeRegistry()) {}
  absl::StatusOr<std::unique_ptr<CelExpression>> CreateExpression(
      const google::api::expr::v1alpha1::Expr* expr,
      const google::api::expr::v1alpha1::SourceInfo* source_info) const override;
  absl::StatusOr<std::unique_ptr<CelExpression>> CreateExpression(
      const google::api::expr::v1alpha1::Expr* expr,
      const google::api::expr::v1alpha1::SourceInfo* source_info,
      std::vector<absl::Status>* warnings) const override;
  absl::StatusOr<std::unique_ptr<CelExpression>> CreateExpression(
      const google::api::expr::v1alpha1::CheckedExpr* checked_expr) const override;
  absl::StatusOr<std::unique_ptr<CelExpression>> CreateExpression(
      const google::api::expr::v1alpha1::CheckedExpr* checked_expr,
      std::vector<absl::Status>* warnings) const override;
  FlatExprBuilder& flat_expr_builder() { return flat_expr_builder_; }
  void set_container(std::string container) override {
    CelExpressionBuilder::set_container(container);
    flat_expr_builder_.set_container(std::move(container));
  }
 private:
  absl::StatusOr<std::unique_ptr<CelExpression>> CreateExpressionImpl(
      std::unique_ptr<cel::Ast> converted_ast,
      std::vector<absl::Status>* warnings) const;
  FlatExprBuilder flat_expr_builder_;
};
}  
#endif  
#include "eval/compiler/cel_expression_builder_flat_impl.h"
#include <memory>
#include <utility>
#include <vector>
#include "google/api/expr/v1alpha1/checked.pb.h"
#include "google/api/expr/v1alpha1/syntax.pb.h"
#include "absl/base/macros.h"
#include "absl/log/check.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "base/ast.h"
#include "common/native_type.h"
#include "eval/eval/cel_expression_flat_impl.h"
#include "eval/eval/direct_expression_step.h"
#include "eval/eval/evaluator_core.h"
#include "eval/public/cel_expression.h"
#include "extensions/protobuf/ast_converters.h"
#include "internal/status_macros.h"
#include "runtime/runtime_issue.h"
namespace google::api::expr::runtime {
using ::cel::Ast;
using ::cel::RuntimeIssue;
using ::google::api::expr::v1alpha1::CheckedExpr;
using ::google::api::expr::v1alpha1::Expr;  
using ::google::api::expr::v1alpha1::SourceInfo;
absl::StatusOr<std::unique_ptr<CelExpression>>
CelExpressionBuilderFlatImpl::CreateExpression(
    const Expr* expr, const SourceInfo* source_info,
    std::vector<absl::Status>* warnings) const {
  ABSL_ASSERT(expr != nullptr);
  CEL_ASSIGN_OR_RETURN(
      std::unique_ptr<Ast> converted_ast,
      cel::extensions::CreateAstFromParsedExpr(*expr, source_info));
  return CreateExpressionImpl(std::move(converted_ast), warnings);
}
absl::StatusOr<std::unique_ptr<CelExpression>>
CelExpressionBuilderFlatImpl::CreateExpression(
    const Expr* expr, const SourceInfo* source_info) const {
  return CreateExpression(expr, source_info,
                          nullptr);
}
absl::StatusOr<std::unique_ptr<CelExpression>>
CelExpressionBuilderFlatImpl::CreateExpression(
    const CheckedExpr* checked_expr,
    std::vector<absl::Status>* warnings) const {
  ABSL_ASSERT(checked_expr != nullptr);
  CEL_ASSIGN_OR_RETURN(
      std::unique_ptr<Ast> converted_ast,
      cel::extensions::CreateAstFromCheckedExpr(*checked_expr));
  return CreateExpressionImpl(std::move(converted_ast), warnings);
}
absl::StatusOr<std::unique_ptr<CelExpression>>
CelExpressionBuilderFlatImpl::CreateExpression(
    const CheckedExpr* checked_expr) const {
  return CreateExpression(checked_expr, nullptr);
}
absl::StatusOr<std::unique_ptr<CelExpression>>
CelExpressionBuilderFlatImpl::CreateExpressionImpl(
    std::unique_ptr<Ast> converted_ast,
    std::vector<absl::Status>* warnings) const {
  std::vector<RuntimeIssue> issues;
  auto* issues_ptr = (warnings != nullptr) ? &issues : nullptr;
  CEL_ASSIGN_OR_RETURN(FlatExpression impl,
                       flat_expr_builder_.CreateExpressionImpl(
                           std::move(converted_ast), issues_ptr));
  if (issues_ptr != nullptr) {
    for (const auto& issue : issues) {
      warnings->push_back(issue.ToStatus());
    }
  }
  if (flat_expr_builder_.options().max_recursion_depth != 0 &&
      !impl.subexpressions().empty() &&
      impl.subexpressions().front().size() == 1 &&
      impl.subexpressions().front().front()->GetNativeTypeId() ==
          cel::NativeTypeId::For<WrappedDirectStep>()) {
    return CelExpressionRecursiveImpl::Create(std::move(impl));
  }
  return std::make_unique<CelExpressionFlatImpl>(std::move(impl));
}
}  