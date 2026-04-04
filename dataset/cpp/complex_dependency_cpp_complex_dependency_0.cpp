#ifndef THIRD_PARTY_CEL_CPP_CODELAB_EXERCISE1_H_
#define THIRD_PARTY_CEL_CPP_CODELAB_EXERCISE1_H_
#include <string>
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
namespace google::api::expr::codelab {
absl::StatusOr<std::string> ParseAndEvaluate(absl::string_view cel_expr);
}  
#endif  
#include "codelab/exercise1.h"
#include <memory>
#include <string>
#include "google/api/expr/v1alpha1/syntax.pb.h"
#include "google/protobuf/arena.h"
#include "absl/status/status.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "eval/public/activation.h"
#include "eval/public/builtin_func_registrar.h"
#include "eval/public/cel_expr_builder_factory.h"
#include "eval/public/cel_expression.h"
#include "eval/public/cel_options.h"
#include "eval/public/cel_value.h"
#include "internal/status_macros.h"
#include "parser/parser.h"
namespace google::api::expr::codelab {
namespace {
using ::google::api::expr::v1alpha1::ParsedExpr;
using ::google::api::expr::parser::Parse;
using ::google::api::expr::runtime::Activation;
using ::google::api::expr::runtime::CelExpression;
using ::google::api::expr::runtime::CelExpressionBuilder;
using ::google::api::expr::runtime::CelValue;
using ::google::api::expr::runtime::CreateCelExpressionBuilder;
using ::google::api::expr::runtime::InterpreterOptions;
using ::google::api::expr::runtime::RegisterBuiltinFunctions;
absl::StatusOr<std::string> ConvertResult(const CelValue& value) {
  if (CelValue::StringHolder inner_value; value.GetValue(&inner_value)) {
    return std::string(inner_value.value());
  } else {
    return absl::InvalidArgumentError(absl::StrCat(
        "expected string result got '", CelValue::TypeName(value.type()), "'"));
  }
}
}  
absl::StatusOr<std::string> ParseAndEvaluate(absl::string_view cel_expr) {
  InterpreterOptions options;
  std::unique_ptr<CelExpressionBuilder> builder =
      CreateCelExpressionBuilder(options);
  CEL_RETURN_IF_ERROR(
      RegisterBuiltinFunctions(builder->GetRegistry(), options));
  ParsedExpr parsed_expr;
  CEL_ASSIGN_OR_RETURN(parsed_expr, Parse(cel_expr));
  google::protobuf::Arena arena;
  Activation activation;
  CEL_ASSIGN_OR_RETURN(std::unique_ptr<CelExpression> expression_plan,
                       builder->CreateExpression(&parsed_expr.expr(),
                                                 &parsed_expr.source_info()));
  CEL_ASSIGN_OR_RETURN(CelValue result,
                       expression_plan->Evaluate(activation, &arena));
  return ConvertResult(result);
}
}  