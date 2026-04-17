#ifndef TENSORSTORE_STATUS_H_
#define TENSORSTORE_STATUS_H_
#include <optional>
#include <string>
#include <string_view>
#include <type_traits>
#include <utility>
#include "absl/base/optimization.h"
#include "absl/status/status.h"
#include "absl/strings/str_format.h"
#include "tensorstore/internal/preprocessor/expand.h"
#include "tensorstore/internal/source_location.h"
#include "tensorstore/internal/type_traits.h"
namespace tensorstore {
namespace internal {
void MaybeAddSourceLocationImpl(absl::Status& status, SourceLocation loc);
absl::Status MaybeAnnotateStatusImpl(absl::Status source,
                                     std::string_view prefix_message,
                                     std::optional<absl::StatusCode> new_code,
                                     std::optional<SourceLocation> loc);
[[noreturn]] void FatalStatus(const char* message, const absl::Status& status,
                              SourceLocation loc);
inline absl::Status MaybeConvertStatusTo(
    absl::Status status, absl::StatusCode code,
    SourceLocation loc = tensorstore::SourceLocation::current()) {
  if (status.code() == code) {
    if (!status.message().empty()) MaybeAddSourceLocationImpl(status, loc);
    return status;
  }
  return MaybeAnnotateStatusImpl(std::move(status), {}, code, loc);
}
inline absl::Status ConvertInvalidArgumentToFailedPrecondition(
    absl::Status status,
    SourceLocation loc = tensorstore::SourceLocation::current()) {
  if (status.code() == absl::StatusCode::kInvalidArgument ||
      status.code() == absl::StatusCode::kOutOfRange) {
    return MaybeAnnotateStatusImpl(std::move(status), {},
                                   absl::StatusCode::kFailedPrecondition, loc);
  }
  return status;
}
template <typename F, typename... Args>
inline absl::Status InvokeForStatus(F&& f, Args&&... args) {
  using R = std::invoke_result_t<F&&, Args&&...>;
  static_assert(std::is_void_v<R> ||
                std::is_same_v<internal::remove_cvref_t<R>, absl::Status>);
  if constexpr (std::is_void_v<R>) {
    std::invoke(static_cast<F&&>(f), static_cast<Args&&>(args)...);
    return absl::OkStatus();
  } else {
    return std::invoke(static_cast<F&&>(f), static_cast<Args&&>(args)...);
  }
}
}  
inline void MaybeAddSourceLocation(
    absl::Status& status,
    SourceLocation loc = tensorstore::SourceLocation::current()) {
  if (status.message().empty()) return;
  internal::MaybeAddSourceLocationImpl(status, loc);
}
std::optional<std::string> AddStatusPayload(absl::Status& status,
                                            std::string_view prefix,
                                            absl::Cord value);
inline absl::Status MaybeAnnotateStatus(
    absl::Status source, std::string_view message,
    SourceLocation loc = tensorstore::SourceLocation::current()) {
  return internal::MaybeAnnotateStatusImpl(std::move(source), message,
                                           std::nullopt, loc);
}
inline absl::Status MaybeAnnotateStatus(
    absl::Status source, std::string_view message, absl::StatusCode new_code,
    SourceLocation loc = tensorstore::SourceLocation::current()) {
  return internal::MaybeAnnotateStatusImpl(std::move(source), message, new_code,
                                           loc);
}
inline const absl::Status& GetStatus(const absl::Status& status) {
  return status;
}
inline absl::Status GetStatus(absl::Status&& status) {
  return std::move(status);
}
}  
#define TENSORSTORE_RETURN_IF_ERROR(...) \
  TENSORSTORE_PP_EXPAND(                 \
      TENSORSTORE_INTERNAL_RETURN_IF_ERROR_IMPL(__VA_ARGS__, _))
#define TENSORSTORE_INTERNAL_RETURN_IF_ERROR_IMPL(expr, error_expr, ...) \
  for (absl::Status _ = ::tensorstore::GetStatus(expr);                  \
       ABSL_PREDICT_FALSE(!_.ok());)                                     \
  return ::tensorstore::MaybeAddSourceLocation(_), error_expr 
#define TENSORSTORE_CHECK_OK(...)                                           \
  do {                                                                      \
    [](const ::absl::Status& tensorstore_check_ok_condition) {              \
      if (ABSL_PREDICT_FALSE(!tensorstore_check_ok_condition.ok())) {       \
        ::tensorstore::internal::FatalStatus(                               \
            "Status not ok: " #__VA_ARGS__, tensorstore_check_ok_condition, \
            tensorstore::SourceLocation::current());                        \
      }                                                                     \
    }(::tensorstore::GetStatus((__VA_ARGS__)));                             \
  } while (false)
#endif  
#if !defined(TENSORSTORE_INTERNAL_STATUS_TEST_HACK)
#include "tensorstore/util/status.h"
#endif
#include <array>
#include <cstdio>
#include <exception>
#include <optional>
#include <string>
#include <string_view>
#include <utility>
#include "absl/status/status.h"
#include "absl/strings/cord.h"
#include "absl/strings/str_format.h"
#include "absl/strings/str_join.h"
#include "tensorstore/internal/source_location.h"
namespace tensorstore {
namespace internal {
void MaybeAddSourceLocationImpl(absl::Status& status, SourceLocation loc) {
  constexpr const char kSourceLocationKey[] = "source locations";
#if TENSORSTORE_HAVE_SOURCE_LOCATION_CURRENT
  if (loc.line() <= 1) return;
  std::string_view filename(loc.file_name());
  if (auto idx = filename.find("tensorstore"); idx != std::string::npos) {
    filename.remove_prefix(idx);
  }
  auto payload = status.GetPayload(kSourceLocationKey);
  if (!payload.has_value()) {
    status.SetPayload(kSourceLocationKey, absl::Cord(absl::StrFormat(
                                              "%s:%d", filename, loc.line())));
  } else {
    payload->Append(absl::StrFormat("\n%s:%d", filename, loc.line()));
    status.SetPayload(kSourceLocationKey, std::move(*payload));
  }
#endif
}
absl::Status MaybeAnnotateStatusImpl(absl::Status source,
                                     std::string_view prefix_message,
                                     std::optional<absl::StatusCode> new_code,
                                     std::optional<SourceLocation> loc) {
  if (source.ok()) return source;
  if (!new_code) new_code = source.code();
  size_t index = 0;
  std::array<std::string_view, 3> to_join = {};
  if (!prefix_message.empty()) {
    to_join[index++] = prefix_message;
  }
  if (!source.message().empty()) {
    to_join[index++] = source.message();
  }
  absl::Status dest(*new_code, (index > 1) ? std::string_view(absl::StrJoin(
                                                 to_join.begin(),
                                                 to_join.begin() + index, ": "))
                                           : to_join[0]);
  source.ForEachPayload([&](auto name, const absl::Cord& value) {
    dest.SetPayload(name, value);
  });
  if (loc) {
    MaybeAddSourceLocation(dest, *loc);
  }
  return dest;
}
[[noreturn]] void FatalStatus(const char* message, const absl::Status& status,
                              SourceLocation loc) {
  std::fprintf(stderr, "%s:%d: %s: %s\n", loc.file_name(), loc.line(), message,
               status.ToString().c_str());
  std::terminate();
}
}  
std::optional<std::string> AddStatusPayload(absl::Status& status,
                                            std::string_view prefix,
                                            absl::Cord value) {
  std::string payload_id(prefix);
  int i = 1;
  while (true) {
    auto p = status.GetPayload(payload_id);
    if (!p.has_value()) {
      break;
    }
    if (p.value() == value) return std::nullopt;
    payload_id = absl::StrFormat("%s[%d]", prefix, i++);
  }
  status.SetPayload(payload_id, std::move(value));
  return payload_id;
}
}  