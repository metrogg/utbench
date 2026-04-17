#ifndef TENSORSTORE_UTIL_STATUS_TESTUTIL_H_
#define TENSORSTORE_UTIL_STATUS_TESTUTIL_H_
#include <ostream>
#include <string>
#include <system_error>  
#include <type_traits>
#include <utility>
#include <gmock/gmock.h>
#include <gtest/gtest.h>
#include "absl/status/status.h"
#include "tensorstore/util/future.h"
#include "tensorstore/util/result.h"
namespace tensorstore {
template <typename T>
void PrintTo(const Result<T>& result, std::ostream* os) {
  if (result) {
    if constexpr (std::is_void_v<T>) {
      *os << "Result{}";
    } else {
      *os << "Result{" << ::testing::PrintToString(*result) << "}";
    }
  } else {
    *os << result.status();
  }
}
template <typename T>
void PrintTo(const Future<T>& future, std::ostream* os) {
  if (!future.ready()) {
    *os << "Future{<not ready>}";
    return;
  }
  if (future.status().ok()) {
    if constexpr (std::is_void_v<T>) {
      *os << "Future{}";
    } else {
      *os << "Future{" << ::testing::PrintToString(future.value()) << "}";
    }
  } else {
    *os << future.status();
  }
}
namespace internal_status {
template <typename StatusType>
class IsOkAndHoldsMatcherImpl : public ::testing::MatcherInterface<StatusType> {
 public:
  typedef
      typename std::remove_reference<StatusType>::type::value_type value_type;
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
      StatusType actual_value,
      ::testing::MatchResultListener* result_listener) const override {
    auto status = ::tensorstore::GetStatus(actual_value);  
    if (!status.ok()) {
      *result_listener << "whose status code is "
                       << absl::StatusCodeToString(status.code());
      return false;
    }
    ::testing::StringMatchResultListener inner_listener;
    if (!inner_matcher_.MatchAndExplain(actual_value.value(),
                                        &inner_listener)) {
      *result_listener << "whose value "
                       << ::testing::PrintToString(actual_value.value())
                       << " doesn't match";
      if (!inner_listener.str().empty()) {
        *result_listener << ", " << inner_listener.str();
      }
      return false;
    }
    return true;
  }
 private:
  const ::testing::Matcher<const value_type&> inner_matcher_;
};
template <typename InnerMatcher>
class IsOkAndHoldsMatcher {
 public:
  explicit IsOkAndHoldsMatcher(InnerMatcher inner_matcher)
      : inner_matcher_(std::move(inner_matcher)) {}
  template <typename StatusType>
  operator ::testing::Matcher<StatusType>() const {  
    return ::testing::Matcher<StatusType>(
        new IsOkAndHoldsMatcherImpl<const StatusType&>(inner_matcher_));
  }
 private:
  const InnerMatcher inner_matcher_;
};
template <typename StatusType>
class MonoIsOkMatcherImpl : public ::testing::MatcherInterface<StatusType> {
 public:
  void DescribeTo(std::ostream* os) const override { *os << "is OK"; }
  void DescribeNegationTo(std::ostream* os) const override {
    *os << "is not OK";
  }
  bool MatchAndExplain(
      StatusType actual_value,
      ::testing::MatchResultListener* result_listener) const override {
    return ::tensorstore::GetStatus(actual_value).ok();  
  }
};
class IsOkMatcher {
 public:
  template <typename StatusType>
  operator ::testing::Matcher<StatusType>() const {  
    return ::testing::Matcher<StatusType>(
        new MonoIsOkMatcherImpl<const StatusType&>());
  }
};
template <typename StatusType>
class StatusIsMatcherImpl : public ::testing::MatcherInterface<StatusType> {
 public:
  explicit StatusIsMatcherImpl(
      testing::Matcher<absl::StatusCode> code_matcher,
      testing::Matcher<const std::string&> message_matcher)
      : code_matcher_(std::move(code_matcher)),
        message_matcher_(std::move(message_matcher)) {}
  void DescribeTo(std::ostream* os) const override {
    *os << "has a status code that ";
    code_matcher_.DescribeTo(os);
    *os << ", and has an error message that ";
    message_matcher_.DescribeTo(os);
  }
  void DescribeNegationTo(std::ostream* os) const override {
    *os << "has a status code that ";
    code_matcher_.DescribeNegationTo(os);
    *os << ", or has an error message that ";
    message_matcher_.DescribeNegationTo(os);
  }
  bool MatchAndExplain(
      StatusType actual_value,
      ::testing::MatchResultListener* result_listener) const override {
    auto status = ::tensorstore::GetStatus(actual_value);  
    testing::StringMatchResultListener inner_listener;
    if (!code_matcher_.MatchAndExplain(status.code(), &inner_listener)) {
      *result_listener << "whose status code "
                       << absl::StatusCodeToString(status.code())
                       << " doesn't match";
      const std::string inner_explanation = inner_listener.str();
      if (!inner_explanation.empty()) {
        *result_listener << ", " << inner_explanation;
      }
      return false;
    }
    if (!message_matcher_.Matches(std::string(status.message()))) {
      *result_listener << "whose error message is wrong";
      return false;
    }
    return true;
  }
 private:
  const testing::Matcher<absl::StatusCode> code_matcher_;
  const testing::Matcher<const std::string&> message_matcher_;
};
class StatusIsMatcher {
 public:
  StatusIsMatcher(testing::Matcher<absl::StatusCode> code_matcher,
                  testing::Matcher<const std::string&> message_matcher)
      : code_matcher_(std::move(code_matcher)),
        message_matcher_(std::move(message_matcher)) {}
  template <typename StatusType>
  operator ::testing::Matcher<StatusType>() const {  
    return ::testing::Matcher<StatusType>(
        new StatusIsMatcherImpl<const StatusType&>(code_matcher_,
                                                   message_matcher_));
  }
 private:
  const testing::Matcher<absl::StatusCode> code_matcher_;
  const testing::Matcher<const std::string&> message_matcher_;
};
}  
inline internal_status::IsOkMatcher IsOk() {
  return internal_status::IsOkMatcher();
}
template <typename InnerMatcher>
internal_status::IsOkAndHoldsMatcher<typename std::decay<InnerMatcher>::type>
IsOkAndHolds(InnerMatcher&& inner_matcher) {
  return internal_status::IsOkAndHoldsMatcher<
      typename std::decay<InnerMatcher>::type>(
      std::forward<InnerMatcher>(inner_matcher));
}
template <typename CodeMatcher, typename MessageMatcher>
internal_status::StatusIsMatcher StatusIs(CodeMatcher code_matcher,
                                          MessageMatcher message_matcher) {
  return internal_status::StatusIsMatcher(std::move(code_matcher),
                                          std::move(message_matcher));
}
template <typename CodeMatcher>
internal_status::StatusIsMatcher StatusIs(CodeMatcher code_matcher) {
  return internal_status::StatusIsMatcher(std::move(code_matcher),
                                          ::testing::_);
}
inline internal_status::StatusIsMatcher MatchesStatus(
    absl::StatusCode status_code) {
  return internal_status::StatusIsMatcher(status_code, ::testing::_);
}
internal_status::StatusIsMatcher MatchesStatus(
    absl::StatusCode status_code, const std::string& message_pattern);
}  
#define TENSORSTORE_EXPECT_OK(expr) EXPECT_THAT(expr, ::tensorstore::IsOk())
#define TENSORSTORE_ASSERT_OK(expr) ASSERT_THAT(expr, ::tensorstore::IsOk())
#define TENSORSTORE_ASSERT_OK_AND_ASSIGN(decl, expr)                      \
  TENSORSTORE_ASSIGN_OR_RETURN(decl, expr,                                \
                               ([&] { FAIL() << #expr << ": " << _; })()) \
#endif  
#include "tensorstore/util/status_testutil.h"
#include <ostream>
#include <regex>  
#include <string>
#include <system_error>  
#include <gmock/gmock.h>
#include "absl/status/status.h"
namespace tensorstore {
namespace internal_status {
namespace {
template <typename StringType>
class RegexMatchImpl : public ::testing::MatcherInterface<StringType> {
 public:
  RegexMatchImpl(const std::string& message_pattern)
      : message_pattern_(message_pattern) {}
  void DescribeTo(std::ostream* os) const override {
    *os << "message matches pattern ";
    ::testing::internal::UniversalPrint(message_pattern_, os);
  }
  void DescribeNegationTo(std::ostream* os) const override {
    *os << "message doesn't match pattern ";
    ::testing::internal::UniversalPrint(message_pattern_, os);
  }
  bool MatchAndExplain(
      StringType message,
      ::testing::MatchResultListener* result_listener) const override {
    return std::regex_match(message, std::regex(message_pattern_));
  }
 private:
  const std::string message_pattern_;
};
}  
}  
internal_status::StatusIsMatcher MatchesStatus(
    absl::StatusCode status_code, const std::string& message_pattern) {
  return internal_status::StatusIsMatcher(
      status_code, ::testing::Matcher<const std::string&>(
                       new internal_status::RegexMatchImpl<const std::string&>(
                           message_pattern)));
}
}  