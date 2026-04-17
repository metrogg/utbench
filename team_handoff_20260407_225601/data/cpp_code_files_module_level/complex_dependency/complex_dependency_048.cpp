#ifndef THIRD_PARTY_CEL_CPP_EVAL_PUBLIC_CEL_FUNCTION_REGISTRY_H_
#define THIRD_PARTY_CEL_CPP_EVAL_PUBLIC_CEL_FUNCTION_REGISTRY_H_
#include <initializer_list>
#include <memory>
#include <string>
#include <utility>
#include <vector>
#include "absl/base/thread_annotations.h"
#include "absl/container/flat_hash_map.h"
#include "absl/container/node_hash_map.h"
#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/mutex.h"
#include "base/function.h"
#include "base/function_descriptor.h"
#include "base/kind.h"
#include "eval/public/cel_function.h"
#include "eval/public/cel_options.h"
#include "eval/public/cel_value.h"
#include "runtime/function_overload_reference.h"
#include "runtime/function_registry.h"
namespace google::api::expr::runtime {
class CelFunctionRegistry {
 public:
  using LazyOverload = cel::FunctionRegistry::LazyOverload;
  CelFunctionRegistry() = default;
  ~CelFunctionRegistry() = default;
  using Registrar = absl::Status (*)(CelFunctionRegistry*,
                                     const InterpreterOptions&);
  absl::Status Register(std::unique_ptr<CelFunction> function) {
    auto descriptor = function->descriptor();
    return Register(descriptor, std::move(function));
  }
  absl::Status Register(const cel::FunctionDescriptor& descriptor,
                        std::unique_ptr<cel::Function> implementation) {
    return modern_registry_.Register(descriptor, std::move(implementation));
  }
  absl::Status RegisterAll(std::initializer_list<Registrar> registrars,
                           const InterpreterOptions& opts);
  absl::Status RegisterLazyFunction(const CelFunctionDescriptor& descriptor) {
    return modern_registry_.RegisterLazyFunction(descriptor);
  }
  std::vector<const CelFunction*> FindOverloads(
      absl::string_view name, bool receiver_style,
      const std::vector<CelValue::Type>& types) const;
  std::vector<cel::FunctionOverloadReference> FindStaticOverloads(
      absl::string_view name, bool receiver_style,
      const std::vector<cel::Kind>& types) const {
    return modern_registry_.FindStaticOverloads(name, receiver_style, types);
  }
  std::vector<const CelFunctionDescriptor*> FindLazyOverloads(
      absl::string_view name, bool receiver_style,
      const std::vector<CelValue::Type>& types) const;
  std::vector<LazyOverload> ModernFindLazyOverloads(
      absl::string_view name, bool receiver_style,
      const std::vector<CelValue::Type>& types) const {
    return modern_registry_.FindLazyOverloads(name, receiver_style, types);
  }
  absl::node_hash_map<std::string, std::vector<const cel::FunctionDescriptor*>>
  ListFunctions() const {
    return modern_registry_.ListFunctions();
  }
  const cel::FunctionRegistry& InternalGetRegistry() const {
    return modern_registry_;
  }
  cel::FunctionRegistry& InternalGetRegistry() { return modern_registry_; }
 private:
  cel::FunctionRegistry modern_registry_;
  mutable absl::Mutex mu_;
  mutable absl::flat_hash_map<const cel::Function*,
                              std::unique_ptr<CelFunction>>
      functions_ ABSL_GUARDED_BY(mu_);
};
}  
#endif  
#include "eval/public/cel_function_registry.h"
#include <algorithm>
#include <initializer_list>
#include <iterator>
#include <memory>
#include <utility>
#include <vector>
#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/mutex.h"
#include "absl/types/span.h"
#include "base/function.h"
#include "base/function_descriptor.h"
#include "base/type_provider.h"
#include "common/type_manager.h"
#include "common/value.h"
#include "common/value_manager.h"
#include "common/values/legacy_value_manager.h"
#include "eval/internal/interop.h"
#include "eval/public/cel_function.h"
#include "eval/public/cel_options.h"
#include "eval/public/cel_value.h"
#include "extensions/protobuf/memory_manager.h"
#include "internal/status_macros.h"
#include "runtime/function_overload_reference.h"
#include "google/protobuf/arena.h"
namespace google::api::expr::runtime {
namespace {
using ::cel::extensions::ProtoMemoryManagerRef;
class ProxyToModernCelFunction : public CelFunction {
 public:
  ProxyToModernCelFunction(const cel::FunctionDescriptor& descriptor,
                           const cel::Function& implementation)
      : CelFunction(descriptor), implementation_(&implementation) {}
  absl::Status Evaluate(absl::Span<const CelValue> args, CelValue* result,
                        google::protobuf::Arena* arena) const override {
    auto memory_manager = ProtoMemoryManagerRef(arena);
    cel::common_internal::LegacyValueManager manager(
        memory_manager, cel::TypeProvider::Builtin());
    cel::FunctionEvaluationContext context(manager);
    std::vector<cel::Value> modern_args =
        cel::interop_internal::LegacyValueToModernValueOrDie(arena, args);
    CEL_ASSIGN_OR_RETURN(auto modern_result,
                         implementation_->Invoke(context, modern_args));
    *result = cel::interop_internal::ModernValueToLegacyValueOrDie(
        arena, modern_result);
    return absl::OkStatus();
  }
 private:
  const cel::Function* implementation_;
};
}  
absl::Status CelFunctionRegistry::RegisterAll(
    std::initializer_list<Registrar> registrars,
    const InterpreterOptions& opts) {
  for (Registrar registrar : registrars) {
    CEL_RETURN_IF_ERROR(registrar(this, opts));
  }
  return absl::OkStatus();
}
std::vector<const CelFunction*> CelFunctionRegistry::FindOverloads(
    absl::string_view name, bool receiver_style,
    const std::vector<CelValue::Type>& types) const {
  std::vector<cel::FunctionOverloadReference> matched_funcs =
      modern_registry_.FindStaticOverloads(name, receiver_style, types);
  std::vector<const CelFunction*> results;
  results.reserve(matched_funcs.size());
  {
    absl::MutexLock lock(&mu_);
    for (cel::FunctionOverloadReference entry : matched_funcs) {
      std::unique_ptr<CelFunction>& legacy_impl =
          functions_[&entry.implementation];
      if (legacy_impl == nullptr) {
        legacy_impl = std::make_unique<ProxyToModernCelFunction>(
            entry.descriptor, entry.implementation);
      }
      results.push_back(legacy_impl.get());
    }
  }
  return results;
}
std::vector<const CelFunctionDescriptor*>
CelFunctionRegistry::FindLazyOverloads(
    absl::string_view name, bool receiver_style,
    const std::vector<CelValue::Type>& types) const {
  std::vector<LazyOverload> lazy_overloads =
      modern_registry_.FindLazyOverloads(name, receiver_style, types);
  std::vector<const CelFunctionDescriptor*> result;
  result.reserve(lazy_overloads.size());
  for (const LazyOverload& overload : lazy_overloads) {
    result.push_back(&overload.descriptor);
  }
  return result;
}
}  