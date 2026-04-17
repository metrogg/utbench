#ifndef AROLLA_EXPR_EVAL_SIDE_OUTPUT_H_
#define AROLLA_EXPR_EVAL_SIDE_OUTPUT_H_
#include <string>
#include "absl/container/flat_hash_map.h"
#include "absl/status/statusor.h"
#include "arolla/expr/expr_node.h"
#include "arolla/io/slot_listener.h"
namespace arolla::expr {
struct ExprWithSideOutputs {
  ExprNodePtr expr;
  absl::flat_hash_map<std::string, ExprNodePtr> side_outputs;
};
absl::StatusOr<ExprWithSideOutputs> ExtractSideOutputs(ExprNodePtr expr);
absl::StatusOr<absl::flat_hash_map<std::string, ExprNodePtr>>
PrepareSideOutputsForListener(
    const absl::flat_hash_map<std::string, ExprNodePtr>& side_outputs,
    const SlotListenerBase& slot_listener);
}  
#endif  
#include "arolla/expr/eval/side_output.h"
#include <string>
#include <utility>
#include "absl/container/flat_hash_map.h"
#include "absl/log/check.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "arolla/expr/annotation_utils.h"
#include "arolla/expr/expr.h"
#include "arolla/expr/expr_debug_string.h"
#include "arolla/expr/expr_node.h"
#include "arolla/expr/expr_visitor.h"
#include "arolla/expr/operators/bootstrap_operators.h"
#include "arolla/io/slot_listener.h"
#include "arolla/util/status_macros_backport.h"
namespace arolla::expr {
absl::StatusOr<ExprWithSideOutputs> ExtractSideOutputs(ExprNodePtr expr) {
  ExprWithSideOutputs result;
  ASSIGN_OR_RETURN(
      result.expr,
      Transform(expr, [&](ExprNodePtr node) -> absl::StatusOr<ExprNodePtr> {
        if (!IsExportAnnotation(node)) {
          return node;
        }
        DCHECK_GE(node->node_deps().size(), 2);  
        auto unwrapped_node = node->node_deps()[0];
        auto tag = ReadExportAnnotationTag(node);
        auto value_expr = ReadExportAnnotationValue(node);
        DCHECK_NE(unwrapped_node, nullptr);
        DCHECK_NE(value_expr, nullptr);
        if (auto [it, inserted] = result.side_outputs.emplace(tag, value_expr);
            !inserted) {
          return absl::FailedPreconditionError(absl::StrCat(
              "duplicated export name ", tag, ": ", GetDebugSnippet(value_expr),
              " vs ", GetDebugSnippet(it->second)));
        }
        return unwrapped_node;
      }));
  return result;
}
absl::StatusOr<absl::flat_hash_map<std::string, ExprNodePtr>>
PrepareSideOutputsForListener(
    const absl::flat_hash_map<std::string, ExprNodePtr>& side_outputs,
    const SlotListenerBase& slot_listener) {
  absl::flat_hash_map<std::string, ExprNodePtr> result;
  for (auto [name, expr] : side_outputs) {
    if (auto qtype = slot_listener.GetQTypeOf(name); qtype != nullptr) {
      ASSIGN_OR_RETURN(expr, expr_operators::CoreCast(expr, Literal(qtype)));
    }
    result.emplace(name, std::move(expr));
  }
  return result;
}
}  