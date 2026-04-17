#ifndef THIRD_PARTY_CEL_CPP_TOOLS_BRANCH_COVERAGE_H_
#define THIRD_PARTY_CEL_CPP_TOOLS_BRANCH_COVERAGE_H_
#include <cstdint>
#include <memory>
#include "google/api/expr/v1alpha1/checked.pb.h"
#include "absl/base/attributes.h"
#include "common/value.h"
#include "eval/public/cel_value.h"
#include "tools/navigable_ast.h"
namespace cel {
class BranchCoverage {
 public:
  struct NodeCoverageStats {
    bool is_boolean;
    int evaluation_count;
    int boolean_true_count;
    int boolean_false_count;
    int error_count;
  };
  virtual ~BranchCoverage() = default;
  virtual void Record(int64_t expr_id, const Value& value) = 0;
  virtual void RecordLegacyValue(
      int64_t expr_id, const google::api::expr::runtime::CelValue& value) = 0;
  virtual NodeCoverageStats StatsForNode(int64_t expr_id) const = 0;
  virtual const NavigableAst& ast() const ABSL_ATTRIBUTE_LIFETIME_BOUND = 0;
  virtual const google::api::expr::v1alpha1::CheckedExpr& expr() const
      ABSL_ATTRIBUTE_LIFETIME_BOUND = 0;
};
std::unique_ptr<BranchCoverage> CreateBranchCoverage(
    const google::api::expr::v1alpha1::CheckedExpr& expr);
}  
#endif  
#include "tools/branch_coverage.h"
#include <cstdint>
#include <memory>
#include "google/api/expr/v1alpha1/checked.pb.h"
#include "absl/base/no_destructor.h"
#include "absl/base/nullability.h"
#include "absl/base/thread_annotations.h"
#include "absl/container/flat_hash_map.h"
#include "absl/container/flat_hash_set.h"
#include "absl/functional/overload.h"
#include "absl/status/status.h"
#include "absl/synchronization/mutex.h"
#include "absl/types/variant.h"
#include "common/value.h"
#include "eval/internal/interop.h"
#include "eval/public/cel_value.h"
#include "tools/navigable_ast.h"
#include "google/protobuf/arena.h"
namespace cel {
namespace {
using ::google::api::expr::v1alpha1::CheckedExpr;
using ::google::api::expr::v1alpha1::Type;
using ::google::api::expr::runtime::CelValue;
const absl::Status& UnsupportedConversionError() {
  static absl::NoDestructor<absl::Status> kErr(
      absl::StatusCode::kInternal, "Conversion to legacy type unsupported.");
  return *kErr;
}
struct ConstantNode {};
struct BoolNode {
  int result_true;
  int result_false;
  int result_error;
};
struct OtherNode {
  int result_error;
};
struct CoverageNode {
  int evaluate_count;
  absl::variant<ConstantNode, OtherNode, BoolNode> kind;
};
absl::Nullable<const Type*> FindCheckerType(const CheckedExpr& expr,
                                            int64_t expr_id) {
  if (auto it = expr.type_map().find(expr_id); it != expr.type_map().end()) {
    return &it->second;
  }
  return nullptr;
}
class BranchCoverageImpl : public BranchCoverage {
 public:
  explicit BranchCoverageImpl(const CheckedExpr& expr) : expr_(expr) {}
  void Record(int64_t expr_id, const Value& value) override {
    auto value_or = interop_internal::ToLegacyValue(&arena_, value);
    if (!value_or.ok()) {
      RecordImpl(expr_id, CelValue::CreateError(&UnsupportedConversionError()));
    } else {
      return RecordImpl(expr_id, *value_or);
    }
  }
  void RecordLegacyValue(int64_t expr_id, const CelValue& value) override {
    return RecordImpl(expr_id, value);
  }
  BranchCoverage::NodeCoverageStats StatsForNode(
      int64_t expr_id) const override;
  const NavigableAst& ast() const override;
  const CheckedExpr& expr() const override;
  void Init();
 private:
  friend class BranchCoverage;
  void RecordImpl(int64_t expr_id, const CelValue& value);
  bool InferredBoolType(const AstNode& node) const;
  CheckedExpr expr_;
  NavigableAst ast_;
  mutable absl::Mutex coverage_nodes_mu_;
  absl::flat_hash_map<int64_t, CoverageNode> coverage_nodes_
      ABSL_GUARDED_BY(coverage_nodes_mu_);
  absl::flat_hash_set<int64_t> unexpected_expr_ids_
      ABSL_GUARDED_BY(coverage_nodes_mu_);
  google::protobuf::Arena arena_;
};
BranchCoverage::NodeCoverageStats BranchCoverageImpl::StatsForNode(
    int64_t expr_id) const {
  BranchCoverage::NodeCoverageStats stats{
      false,
      0,
      0,
      0,
      0,
  };
  absl::MutexLock lock(&coverage_nodes_mu_);
  auto it = coverage_nodes_.find(expr_id);
  if (it != coverage_nodes_.end()) {
    const CoverageNode& coverage_node = it->second;
    stats.evaluation_count = coverage_node.evaluate_count;
    absl::visit(absl::Overload([&](const ConstantNode& cov) {},
                               [&](const OtherNode& cov) {
                                 stats.error_count = cov.result_error;
                               },
                               [&](const BoolNode& cov) {
                                 stats.is_boolean = true;
                                 stats.boolean_true_count = cov.result_true;
                                 stats.boolean_false_count = cov.result_false;
                                 stats.error_count = cov.result_error;
                               }),
                coverage_node.kind);
    return stats;
  }
  return stats;
}
const NavigableAst& BranchCoverageImpl::ast() const { return ast_; }
const CheckedExpr& BranchCoverageImpl::expr() const { return expr_; }
bool BranchCoverageImpl::InferredBoolType(const AstNode& node) const {
  int64_t expr_id = node.expr()->id();
  const auto* checker_type = FindCheckerType(expr_, expr_id);
  if (checker_type != nullptr) {
    return checker_type->has_primitive() &&
           checker_type->primitive() == Type::BOOL;
  }
  return false;
}
void BranchCoverageImpl::Init() ABSL_NO_THREAD_SAFETY_ANALYSIS {
  ast_ = NavigableAst::Build(expr_.expr());
  for (const AstNode& node : ast_.Root().DescendantsPreorder()) {
    int64_t expr_id = node.expr()->id();
    CoverageNode& coverage_node = coverage_nodes_[expr_id];
    coverage_node.evaluate_count = 0;
    if (node.node_kind() == NodeKind::kConstant) {
      coverage_node.kind = ConstantNode{};
    } else if (InferredBoolType(node)) {
      coverage_node.kind = BoolNode{0, 0, 0};
    } else {
      coverage_node.kind = OtherNode{0};
    }
  }
}
void BranchCoverageImpl::RecordImpl(int64_t expr_id, const CelValue& value) {
  absl::MutexLock lock(&coverage_nodes_mu_);
  auto it = coverage_nodes_.find(expr_id);
  if (it == coverage_nodes_.end()) {
    unexpected_expr_ids_.insert(expr_id);
    it = coverage_nodes_.insert({expr_id, CoverageNode{0, {}}}).first;
    if (value.IsBool()) {
      it->second.kind = BoolNode{0, 0, 0};
    }
  }
  CoverageNode& coverage_node = it->second;
  coverage_node.evaluate_count++;
  bool is_error = value.IsError() &&
                  value.ErrorOrDie() != &UnsupportedConversionError();
  absl::visit(absl::Overload([&](ConstantNode& node) {},
                             [&](OtherNode& cov) {
                               if (is_error) {
                                 cov.result_error++;
                               }
                             },
                             [&](BoolNode& cov) {
                               if (value.IsBool()) {
                                 bool held_value = value.BoolOrDie();
                                 if (held_value) {
                                   cov.result_true++;
                                 } else {
                                   cov.result_false++;
                                 }
                               } else if (is_error) {
                                 cov.result_error++;
                               }
                             }),
              coverage_node.kind);
}
}  
std::unique_ptr<BranchCoverage> CreateBranchCoverage(const CheckedExpr& expr) {
  auto result = std::make_unique<BranchCoverageImpl>(expr);
  result->Init();
  return result;
}
}  