#ifndef THIRD_PARTY_CEL_CPP_EVAL_PUBLIC_PORTABLE_CEL_EXPR_BUILDER_FACTORY_H_
#define THIRD_PARTY_CEL_CPP_EVAL_PUBLIC_PORTABLE_CEL_EXPR_BUILDER_FACTORY_H_
#include "eval/public/cel_expression.h"
#include "eval/public/cel_options.h"
#include "eval/public/structs/legacy_type_provider.h"
namespace google {
namespace api {
namespace expr {
namespace runtime {
std::unique_ptr<CelExpressionBuilder> CreatePortableExprBuilder(
    std::unique_ptr<LegacyTypeProvider> type_provider,
    const InterpreterOptions& options = InterpreterOptions());
}  
}  
}  
}  
#endif  
#include "eval/public/portable_cel_expr_builder_factory.h"
#include <memory>
#include <utility>
#include "absl/log/absl_log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "base/ast_internal/ast_impl.h"
#include "base/kind.h"
#include "common/memory.h"
#include "common/values/legacy_type_reflector.h"
#include "eval/compiler/cel_expression_builder_flat_impl.h"
#include "eval/compiler/comprehension_vulnerability_check.h"
#include "eval/compiler/constant_folding.h"
#include "eval/compiler/flat_expr_builder.h"
#include "eval/compiler/flat_expr_builder_extensions.h"
#include "eval/compiler/qualified_reference_resolver.h"
#include "eval/compiler/regex_precompilation_optimization.h"
#include "eval/public/cel_expression.h"
#include "eval/public/cel_function.h"
#include "eval/public/cel_options.h"
#include "eval/public/structs/legacy_type_provider.h"
#include "extensions/protobuf/memory_manager.h"
#include "extensions/select_optimization.h"
#include "runtime/runtime_options.h"
namespace google::api::expr::runtime {
namespace {
using ::cel::MemoryManagerRef;
using ::cel::ast_internal::AstImpl;
using ::cel::extensions::CreateSelectOptimizationProgramOptimizer;
using ::cel::extensions::kCelAttribute;
using ::cel::extensions::kCelHasField;
using ::cel::extensions::ProtoMemoryManagerRef;
using ::cel::extensions::SelectOptimizationAstUpdater;
using ::cel::runtime_internal::CreateConstantFoldingOptimizer;
struct ArenaBackedConstfoldingFactory {
  MemoryManagerRef memory_manager;
  absl::StatusOr<std::unique_ptr<ProgramOptimizer>> operator()(
      PlannerContext& ctx, const AstImpl& ast) const {
    return CreateConstantFoldingOptimizer(memory_manager)(ctx, ast);
  }
};
}  
std::unique_ptr<CelExpressionBuilder> CreatePortableExprBuilder(
    std::unique_ptr<LegacyTypeProvider> type_provider,
    const InterpreterOptions& options) {
  if (type_provider == nullptr) {
    ABSL_LOG(ERROR) << "Cannot pass nullptr as type_provider to "
                       "CreatePortableExprBuilder";
    return nullptr;
  }
  cel::RuntimeOptions runtime_options = ConvertToRuntimeOptions(options);
  auto builder =
      std::make_unique<CelExpressionBuilderFlatImpl>(runtime_options);
  builder->GetTypeRegistry()
      ->InternalGetModernRegistry()
      .set_use_legacy_container_builders(options.use_legacy_container_builders);
  builder->GetTypeRegistry()->RegisterTypeProvider(std::move(type_provider));
  FlatExprBuilder& flat_expr_builder = builder->flat_expr_builder();
  flat_expr_builder.AddAstTransform(NewReferenceResolverExtension(
      (options.enable_qualified_identifier_rewrites)
          ? ReferenceResolverOption::kAlways
          : ReferenceResolverOption::kCheckedOnly));
  if (options.enable_comprehension_vulnerability_check) {
    builder->flat_expr_builder().AddProgramOptimizer(
        CreateComprehensionVulnerabilityCheck());
  }
  if (options.constant_folding) {
    builder->flat_expr_builder().AddProgramOptimizer(
        ArenaBackedConstfoldingFactory{
            ProtoMemoryManagerRef(options.constant_arena)});
  }
  if (options.enable_regex_precompilation) {
    flat_expr_builder.AddProgramOptimizer(
        CreateRegexPrecompilationExtension(options.regex_max_program_size));
  }
  if (options.enable_select_optimization) {
    flat_expr_builder.AddAstTransform(
        std::make_unique<SelectOptimizationAstUpdater>());
    absl::Status status =
        builder->GetRegistry()->RegisterLazyFunction(CelFunctionDescriptor(
            kCelAttribute, false, {cel::Kind::kAny, cel::Kind::kList}));
    if (!status.ok()) {
      ABSL_LOG(ERROR) << "Failed to register " << kCelAttribute << ": "
                      << status;
    }
    status = builder->GetRegistry()->RegisterLazyFunction(CelFunctionDescriptor(
        kCelHasField, false, {cel::Kind::kAny, cel::Kind::kList}));
    if (!status.ok()) {
      ABSL_LOG(ERROR) << "Failed to register " << kCelHasField << ": "
                      << status;
    }
    flat_expr_builder.AddProgramOptimizer(
        CreateSelectOptimizationProgramOptimizer());
  }
  return builder;
}
}  