#ifndef ABSL_STATUS_INTERNAL_STATUS_MATCHERS_H_
#define ABSL_STATUS_INTERNAL_STATUS_MATCHERS_H_
#include <ostream>  
#include <string>
#include <type_traits>
#include <utility>
#include "gmock/gmock.h"  
#include "absl/base/config.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
namespace absl_testing {
ABSL_NAMESPACE_BEGIN
namespace status_internal {
inline const absl::Status& GetStatus(const absl::Status& status) {
  return status;
}
template <typename T>
inline const absl::Status& GetStatus(const absl::StatusOr<T>& status) {
  return status.status();
}
template <typename StatusOrType>
class IsOkAndHoldsMatcherImpl
    : public ::testing::MatcherInterface<StatusOrType> {
 public:
  typedef
      typename std::remove_reference<StatusOrType>::type::value_type value_type;
  template <typename InnerMatcher>
  explicit IsOkAndHoldsMatcherImpl(InnerMatcher&& inner_matcher)
      : inner_matcher_(::testing::SafeMatcherCast<const value_type&>(
            std::forward<InnerMatcher>(inner_matcher))) {}
  void DescribeTo(std::ostream* os) const override {
    *os << "is OK and has a value that ";
    inner_matcher_.DescribeTo(os);
  }
  void DescribeNegationTo(std::ostream* os) const override {
    *os << "isn't OK or has a value that ";
    inner_matcher_.DescribeNegationTo(os);
  }
  bool MatchAndExplain(
      StatusOrType actual_value,
      ::testing::MatchResultListener* result_listener) const override {
    if (!actual_value.ok()) {
      *result_listener << "which has status " << actual_value.status();
      return false;
    }
    return inner_matcher_.MatchAndExplain(*actual_value, result_listener);
  }
 private:
  const ::testing::Matcher<const value_type&> inner_matcher_;
};
template <typename InnerMatcher>
class IsOkAndHoldsMatcher {
 public:
  explicit IsOkAndHoldsMatcher(InnerMatcher inner_matcher)
      : inner_matcher_(std::forward<InnerMatcher>(inner_matcher)) {}
  template <typename StatusOrType>
  operator ::testing::Matcher<StatusOrType>() const {  
    return ::testing::Matcher<StatusOrType>(
        new IsOkAndHoldsMatcherImpl<const StatusOrType&>(inner_matcher_));
  }
 private:
  const InnerMatcher inner_matcher_;
};
class StatusCode {
 public:
   StatusCode(int code)  
      : code_(static_cast<::absl::StatusCode>(code)) {}
   StatusCode(::absl::StatusCode code) : code_(code) {}  
  explicit operator int() const { return static_cast<int>(code_); }
  friend inline void PrintTo(const StatusCode& code, std::ostream* os) {
    *os << static_cast<int>(code);
  }
 private:
  ::absl::StatusCode code_;
};
inline bool operator==(const StatusCode& lhs, const StatusCode& rhs) {
  return static_cast<int>(lhs) == static_cast<int>(rhs);
}
inline bool operator!=(const StatusCode& lhs, const StatusCode& rhs) {
  return static_cast<int>(lhs) != static_cast<int>(rhs);
}
class StatusIsMatcherCommonImpl {
 public:
  StatusIsMatcherCommonImpl(
      ::testing::Matcher<StatusCode> code_matcher,
      ::testing::Matcher<absl::string_view> message_matcher)
      : code_matcher_(std::move(code_matcher)),
        message_matcher_(std::move(message_matcher)) {}
  void DescribeTo(std::ostream* os) const;
  void DescribeNegationTo(std::ostream* os) const;
  bool MatchAndExplain(const absl::Status& status,
                       ::testing::MatchResultListener* result_listener) const;
 private:
  const ::testing::Matcher<StatusCode> code_matcher_;
  const ::testing::Matcher<absl::string_view> message_matcher_;
};
template <typename T>
class MonoStatusIsMatcherImpl : public ::testing::MatcherInterface<T> {
 public:
  explicit MonoStatusIsMatcherImpl(StatusIsMatcherCommonImpl common_impl)
      : common_impl_(std::move(common_impl)) {}
  void DescribeTo(std::ostream* os) const override {
    common_impl_.DescribeTo(os);
  }
  void DescribeNegationTo(std::ostream* os) const override {
    common_impl_.DescribeNegationTo(os);
  }
  bool MatchAndExplain(
      T actual_value,
      ::testing::MatchResultListener* result_listener) const override {
    return common_impl_.MatchAndExplain(GetStatus(actual_value),
                                        result_listener);
  }
 private:
  StatusIsMatcherCommonImpl common_impl_;
};
class StatusIsMatcher {
 public:
  template <typename StatusCodeMatcher, typename StatusMessageMatcher>
  StatusIsMatcher(StatusCodeMatcher&& code_matcher,
                  StatusMessageMatcher&& message_matcher)
      : common_impl_(::testing::MatcherCast<StatusCode>(
                         std::forward<StatusCodeMatcher>(code_matcher)),
                     ::testing::MatcherCast<absl::string_view>(
                         std::forward<StatusMessageMatcher>(message_matcher))) {
  }
  template <typename T>
   operator ::testing::Matcher<T>() const {  
    return ::testing::Matcher<T>(
        new MonoStatusIsMatcherImpl<const T&>(common_impl_));
  }
 private:
  const StatusIsMatcherCommonImpl common_impl_;
};
template <typename T>
class MonoIsOkMatcherImpl : public ::testing::MatcherInterface<T> {
 public:
  void DescribeTo(std::ostream* os) const override { *os << "is OK"; }
  void DescribeNegationTo(std::ostream* os) const override {
    *os << "is not OK";
  }
  bool MatchAndExplain(T actual_value,
                       ::testing::MatchResultListener*) const override {
    return GetStatus(actual_value).ok();
  }
};
class IsOkMatcher {
 public:
  template <typename T>
   operator ::testing::Matcher<T>() const {  
    return ::testing::Matcher<T>(new MonoIsOkMatcherImpl<const T&>());
  }
};
}  
ABSL_NAMESPACE_END
}  
#endif  
#include "absl/status/internal/status_matchers.h"
#include <ostream>
#include <string>
#include "gmock/gmock.h"  
#include "absl/base/config.h"
#include "absl/status/status.h"
namespace absl_testing {
ABSL_NAMESPACE_BEGIN
namespace status_internal {
void StatusIsMatcherCommonImpl::DescribeTo(std::ostream* os) const {
  *os << ", has a status code that ";
  code_matcher_.DescribeTo(os);
  *os << ", and has an error message that ";
  message_matcher_.DescribeTo(os);
}
void StatusIsMatcherCommonImpl::DescribeNegationTo(std::ostream* os) const {
  *os << ", or has a status code that ";
  code_matcher_.DescribeNegationTo(os);
  *os << ", or has an error message that ";
  message_matcher_.DescribeNegationTo(os);
}
bool StatusIsMatcherCommonImpl::MatchAndExplain(
    const ::absl::Status& status,
    ::testing::MatchResultListener* result_listener) const {
  ::testing::StringMatchResultListener inner_listener;
  if (!code_matcher_.MatchAndExplain(status.code(), &inner_listener)) {
    *result_listener << (inner_listener.str().empty()
                             ? "whose status code is wrong"
                             : "which has a status code " +
                                   inner_listener.str());
    return false;
  }
  if (!message_matcher_.Matches(std::string(status.message()))) {
    *result_listener << "whose error message is wrong";
    return false;
  }
  return true;
}
}  
ABSL_NAMESPACE_END
}  