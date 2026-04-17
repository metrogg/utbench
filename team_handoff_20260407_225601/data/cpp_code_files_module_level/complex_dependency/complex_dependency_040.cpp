#ifndef THIRD_PARTY_CEL_CPP_EVAL_COMPILER_RESOLVER_H_
#define THIRD_PARTY_CEL_CPP_EVAL_COMPILER_RESOLVER_H_
#include <cstdint>
#include <string>
#include <utility>
#include <vector>
#include "absl/container/flat_hash_map.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/types/optional.h"
#include "base/kind.h"
#include "common/value.h"
#include "common/value_manager.h"
#include "runtime/function_overload_reference.h"
#include "runtime/function_registry.h"
#include "runtime/type_registry.h"
namespace google::api::expr::runtime {
class Resolver {
 public:
  Resolver(
      absl::string_view container,
      const cel::FunctionRegistry& function_registry,
      const cel::TypeRegistry& type_registry, cel::ValueManager& value_factory,
      const absl::flat_hash_map<std::string, cel::TypeRegistry::Enumeration>&
          resolveable_enums,
      bool resolve_qualified_type_identifiers = true);
  ~Resolver() = default;
  absl::optional<cel::Value> FindConstant(absl::string_view name,
                                          int64_t expr_id) const;
  absl::StatusOr<absl::optional<std::pair<std::string, cel::Type>>> FindType(
      absl::string_view name, int64_t expr_id) const;
  std::vector<cel::FunctionRegistry::LazyOverload> FindLazyOverloads(
      absl::string_view name, bool receiver_style,
      const std::vector<cel::Kind>& types, int64_t expr_id = -1) const;
  std::vector<cel::FunctionOverloadReference> FindOverloads(
      absl::string_view name, bool receiver_style,
      const std::vector<cel::Kind>& types, int64_t expr_id = -1) const;
  std::vector<std::string> FullyQualifiedNames(absl::string_view base_name,
                                               int64_t expr_id = -1) const;
 private:
  std::vector<std::string> namespace_prefixes_;
  absl::flat_hash_map<std::string, cel::Value> enum_value_map_;
  const cel::FunctionRegistry& function_registry_;
  cel::ValueManager& value_factory_;
  const absl::flat_hash_map<std::string, cel::TypeRegistry::Enumeration>&
      resolveable_enums_;
  bool resolve_qualified_type_identifiers_;
};
inline std::vector<cel::Kind> ArgumentsMatcher(int argument_count) {
  std::vector<cel::Kind> argument_matcher(argument_count);
  for (int i = 0; i < argument_count; i++) {
    argument_matcher[i] = cel::Kind::kAny;
  }
  return argument_matcher;
}
}  
#endif  
#include "eval/compiler/resolver.h"
#include <cstdint>
#include <memory>
#include <string>
#include <utility>
#include <vector>
#include "absl/base/nullability.h"
#include "absl/container/flat_hash_map.h"
#include "absl/status/statusor.h"
#include "absl/strings/match.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/str_split.h"
#include "absl/strings/string_view.h"
#include "absl/strings/strip.h"
#include "absl/types/optional.h"
#include "base/kind.h"
#include "common/memory.h"
#include "common/type.h"
#include "common/value.h"
#include "common/value_manager.h"
#include "internal/status_macros.h"
#include "runtime/function_overload_reference.h"
#include "runtime/function_registry.h"
#include "runtime/type_registry.h"
namespace google::api::expr::runtime {
using ::cel::Value;
Resolver::Resolver(
    absl::string_view container, const cel::FunctionRegistry& function_registry,
    const cel::TypeRegistry&, cel::ValueManager& value_factory,
    const absl::flat_hash_map<std::string, cel::TypeRegistry::Enumeration>&
        resolveable_enums,
    bool resolve_qualified_type_identifiers)
    : namespace_prefixes_(),
      enum_value_map_(),
      function_registry_(function_registry),
      value_factory_(value_factory),
      resolveable_enums_(resolveable_enums),
      resolve_qualified_type_identifiers_(resolve_qualified_type_identifiers) {
  auto container_elements = absl::StrSplit(container, '.');
  std::string prefix = "";
  namespace_prefixes_.push_back(prefix);
  for (const auto& elem : container_elements) {
    if (elem.empty()) {
      continue;
    }
    absl::StrAppend(&prefix, elem, ".");
    namespace_prefixes_.insert(namespace_prefixes_.begin(), prefix);
  }
  for (const auto& prefix : namespace_prefixes_) {
    for (auto iter = resolveable_enums_.begin();
         iter != resolveable_enums_.end(); ++iter) {
      absl::string_view enum_name = iter->first;
      if (!absl::StartsWith(enum_name, prefix)) {
        continue;
      }
      auto remainder = absl::StripPrefix(enum_name, prefix);
      const auto& enum_type = iter->second;
      for (const auto& enumerator : enum_type.enumerators) {
        auto key = absl::StrCat(remainder, !remainder.empty() ? "." : "",
                                enumerator.name);
        enum_value_map_[key] = value_factory.CreateIntValue(enumerator.number);
      }
    }
  }
}
std::vector<std::string> Resolver::FullyQualifiedNames(absl::string_view name,
                                                       int64_t expr_id) const {
  std::vector<std::string> names;
  if (absl::StartsWith(name, ".")) {
    std::string fully_qualified_name = std::string(name.substr(1));
    names.push_back(fully_qualified_name);
    return names;
  }
  for (const auto& prefix : namespace_prefixes_) {
    std::string fully_qualified_name = absl::StrCat(prefix, name);
    names.push_back(fully_qualified_name);
  }
  return names;
}
absl::optional<cel::Value> Resolver::FindConstant(absl::string_view name,
                                                  int64_t expr_id) const {
  auto names = FullyQualifiedNames(name, expr_id);
  for (const auto& name : names) {
    auto enum_entry = enum_value_map_.find(name);
    if (enum_entry != enum_value_map_.end()) {
      return enum_entry->second;
    }
    if (resolve_qualified_type_identifiers_ || !absl::StrContains(name, ".")) {
      auto type_value = value_factory_.FindType(name);
      if (type_value.ok() && type_value->has_value()) {
        return value_factory_.CreateTypeValue(**type_value);
      }
    }
  }
  return absl::nullopt;
}
std::vector<cel::FunctionOverloadReference> Resolver::FindOverloads(
    absl::string_view name, bool receiver_style,
    const std::vector<cel::Kind>& types, int64_t expr_id) const {
  std::vector<cel::FunctionOverloadReference> funcs;
  auto names = FullyQualifiedNames(name, expr_id);
  for (auto it = names.begin(); it != names.end(); it++) {
    funcs = function_registry_.FindStaticOverloads(*it, receiver_style, types);
    if (!funcs.empty()) {
      return funcs;
    }
  }
  return funcs;
}
std::vector<cel::FunctionRegistry::LazyOverload> Resolver::FindLazyOverloads(
    absl::string_view name, bool receiver_style,
    const std::vector<cel::Kind>& types, int64_t expr_id) const {
  std::vector<cel::FunctionRegistry::LazyOverload> funcs;
  auto names = FullyQualifiedNames(name, expr_id);
  for (const auto& name : names) {
    funcs = function_registry_.FindLazyOverloads(name, receiver_style, types);
    if (!funcs.empty()) {
      return funcs;
    }
  }
  return funcs;
}
absl::StatusOr<absl::optional<std::pair<std::string, cel::Type>>>
Resolver::FindType(absl::string_view name, int64_t expr_id) const {
  auto qualified_names = FullyQualifiedNames(name, expr_id);
  for (auto& qualified_name : qualified_names) {
    CEL_ASSIGN_OR_RETURN(auto maybe_type,
                         value_factory_.FindType(qualified_name));
    if (maybe_type.has_value()) {
      return std::make_pair(std::move(qualified_name), std::move(*maybe_type));
    }
  }
  return absl::nullopt;
}
}  