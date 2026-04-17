#ifndef ABSL_STRINGS_STR_REPLACE_H_
#define ABSL_STRINGS_STR_REPLACE_H_
#include <string>
#include <utility>
#include <vector>
#include "absl/base/attributes.h"
#include "absl/base/nullability.h"
#include "absl/strings/string_view.h"
namespace absl {
ABSL_NAMESPACE_BEGIN
ABSL_MUST_USE_RESULT std::string StrReplaceAll(
    absl::string_view s,
    std::initializer_list<std::pair<absl::string_view, absl::string_view>>
        replacements);
template <typename StrToStrMapping>
std::string StrReplaceAll(absl::string_view s,
                          const StrToStrMapping& replacements);
int StrReplaceAll(
    std::initializer_list<std::pair<absl::string_view, absl::string_view>>
        replacements,
    absl::Nonnull<std::string*> target);
template <typename StrToStrMapping>
int StrReplaceAll(const StrToStrMapping& replacements,
                  absl::Nonnull<std::string*> target);
namespace strings_internal {
struct ViableSubstitution {
  absl::string_view old;
  absl::string_view replacement;
  size_t offset;
  ViableSubstitution(absl::string_view old_str,
                     absl::string_view replacement_str, size_t offset_val)
      : old(old_str), replacement(replacement_str), offset(offset_val) {}
  bool OccursBefore(const ViableSubstitution& y) const {
    if (offset != y.offset) return offset < y.offset;
    return old.size() > y.old.size();
  }
};
template <typename StrToStrMapping>
std::vector<ViableSubstitution> FindSubstitutions(
    absl::string_view s, const StrToStrMapping& replacements) {
  std::vector<ViableSubstitution> subs;
  subs.reserve(replacements.size());
  for (const auto& rep : replacements) {
    using std::get;
    absl::string_view old(get<0>(rep));
    size_t pos = s.find(old);
    if (pos == s.npos) continue;
    if (old.empty()) continue;
    subs.emplace_back(old, get<1>(rep), pos);
    size_t index = subs.size();
    while (--index && subs[index - 1].OccursBefore(subs[index])) {
      std::swap(subs[index], subs[index - 1]);
    }
  }
  return subs;
}
int ApplySubstitutions(absl::string_view s,
                       absl::Nonnull<std::vector<ViableSubstitution>*> subs_ptr,
                       absl::Nonnull<std::string*> result_ptr);
}  
template <typename StrToStrMapping>
std::string StrReplaceAll(absl::string_view s,
                          const StrToStrMapping& replacements) {
  auto subs = strings_internal::FindSubstitutions(s, replacements);
  std::string result;
  result.reserve(s.size());
  strings_internal::ApplySubstitutions(s, &subs, &result);
  return result;
}
template <typename StrToStrMapping>
int StrReplaceAll(const StrToStrMapping& replacements,
                  absl::Nonnull<std::string*> target) {
  auto subs = strings_internal::FindSubstitutions(*target, replacements);
  if (subs.empty()) return 0;
  std::string result;
  result.reserve(target->size());
  int substitutions =
      strings_internal::ApplySubstitutions(*target, &subs, &result);
  target->swap(result);
  return substitutions;
}
ABSL_NAMESPACE_END
}  
#endif  
#include "absl/strings/str_replace.h"
#include <cstddef>
#include <initializer_list>
#include <string>
#include <utility>
#include <vector>
#include "absl/base/config.h"
#include "absl/base/nullability.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
namespace absl {
ABSL_NAMESPACE_BEGIN
namespace strings_internal {
using FixedMapping =
    std::initializer_list<std::pair<absl::string_view, absl::string_view>>;
int ApplySubstitutions(
    absl::string_view s,
    absl::Nonnull<std::vector<strings_internal::ViableSubstitution>*> subs_ptr,
    absl::Nonnull<std::string*> result_ptr) {
  auto& subs = *subs_ptr;
  int substitutions = 0;
  size_t pos = 0;
  while (!subs.empty()) {
    auto& sub = subs.back();
    if (sub.offset >= pos) {
      if (pos <= s.size()) {
        StrAppend(result_ptr, s.substr(pos, sub.offset - pos), sub.replacement);
      }
      pos = sub.offset + sub.old.size();
      substitutions += 1;
    }
    sub.offset = s.find(sub.old, pos);
    if (sub.offset == s.npos) {
      subs.pop_back();
    } else {
      size_t index = subs.size();
      while (--index && subs[index - 1].OccursBefore(subs[index])) {
        std::swap(subs[index], subs[index - 1]);
      }
    }
  }
  result_ptr->append(s.data() + pos, s.size() - pos);
  return substitutions;
}
}  
std::string StrReplaceAll(absl::string_view s,
                          strings_internal::FixedMapping replacements) {
  return StrReplaceAll<strings_internal::FixedMapping>(s, replacements);
}
int StrReplaceAll(strings_internal::FixedMapping replacements,
                  absl::Nonnull<std::string*> target) {
  return StrReplaceAll<strings_internal::FixedMapping>(replacements, target);
}
ABSL_NAMESPACE_END
}  