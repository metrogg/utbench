#ifndef TENSORFLOW_TSL_PLATFORM_STATUS_MATCHERS_H_
#define TENSORFLOW_TSL_PLATFORM_STATUS_MATCHERS_H_
#include <ostream>
#include <string>
#include <utility>
#include "tsl/platform/status.h"
#include "tsl/platform/statusor.h"
#include "tsl/platform/test.h"
#include "tsl/protobuf/error_codes.pb.h"
namespace tsl {
inline void PrintTo(const tsl::error::Code code, std::ostream* os) {
  *os << Code_Name(code);
}
template <typename T>
void PrintTo(const StatusOr<T>& status_or, std::ostream* os) {
  *os << ::testing::PrintToString(status_or.status());
  if (status_or.ok()) {
    *os << ": " << ::testing::PrintToString(status_or.value());
  }
}
namespace testing {
namespace internal_status {
inline const absl::Status& GetStatus(const absl::Status& status) {
  return status;
}
template <typename T>
inline const absl::Status& GetStatus(const StatusOr<T>& status) {
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
    ::testing::StringMatchResultListener inner_listener;
    const bool matches =
        inner_matcher_.MatchAndExplain(*actual_value, &inner_listener);
    const std::string inner_explanation = inner_listener.str();
    if (!inner_explanation.empty()) {
      *result_listener << "which contains value "
                       << ::testing::PrintToString(*actual_value) << ", "
                       << inner_explanation;
    }
    return matches;
  }
 private:
  const ::testing::Matcher<const value_type&> inner_matcher_;
};
template <typename InnerMatcher>
class IsOkAndHoldsMatcher {
 public:
  explicit IsOkAndHoldsMatcher(InnerMatcher inner_matcher)
      : inner_matcher_(std::move(inner_matcher)) {}
  template <typename StatusOrType>
  operator ::testing::Matcher<StatusOrType>() const {  
    return ::testing::Matcher<StatusOrType>(
        new IsOkAndHoldsMatcherImpl<const StatusOrType&>(inner_matcher_));
  }
 private:
  const InnerMatcher inner_matcher_;
};
class StatusIsMatcherCommonImpl {
 public:
  StatusIsMatcherCommonImpl(
      ::testing::Matcher<const absl::StatusCode> code_matcher,
      ::testing::Matcher<const std::string&> message_matcher)
      : code_matcher_(std::move(code_matcher)),
        message_matcher_(std::move(message_matcher)) {}
  void DescribeTo(std::ostream* os) const;
  void DescribeNegationTo(std::ostream* os) const;
  bool MatchAndExplain(const absl::Status& status,
                       ::testing::MatchResultListener* result_listener) const;
 private:
  const ::testing::Matcher<const absl::StatusCode> code_matcher_;
  const ::testing::Matcher<const std::string&> message_matcher_;
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
  StatusIsMatcher(::testing::Matcher<const absl::StatusCode> code_matcher,
                  ::testing::Matcher<const std::string&> message_matcher)
      : common_impl_(
            ::testing::MatcherCast<const absl::StatusCode>(code_matcher),
            ::testing::MatcherCast<const std::string&>(message_matcher)) {}
  template <typename T>
  operator ::testing::Matcher<T>() const {  
    return ::testing::MakeMatcher(new MonoStatusIsMatcherImpl<T>(common_impl_));
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
template <typename MessageMatcher>
internal_status::StatusIsMatcher StatusIs(tensorflow::error::Code code_matcher,
                                          MessageMatcher message_matcher) {
  return internal_status::StatusIsMatcher(
      static_cast<absl::StatusCode>(code_matcher), std::move(message_matcher));
}
template <typename CodeMatcher>
internal_status::StatusIsMatcher StatusIs(CodeMatcher code_matcher) {
  return StatusIs(std::move(code_matcher), ::testing::_);
}
template <>
inline internal_status::StatusIsMatcher StatusIs(
    tensorflow::error::Code code_matcher) {
  return StatusIs(static_cast<absl::StatusCode>(code_matcher), ::testing::_);
}
inline internal_status::IsOkMatcher IsOk() {
  return internal_status::IsOkMatcher();
}
}  
}  
#endif  
#include "tsl/platform/status_matchers.h"
#include <ostream>
#include <string>
#include "tsl/platform/status.h"
#include "tsl/platform/test.h"
#include "tsl/protobuf/error_codes.pb.h"
namespace tsl {
namespace testing {
namespace internal_status {
void StatusIsMatcherCommonImpl::DescribeTo(std::ostream* os) const {
  *os << "has a status code that ";
  code_matcher_.DescribeTo(os);
  *os << ", and has an error message that ";
  message_matcher_.DescribeTo(os);
}
void StatusIsMatcherCommonImpl::DescribeNegationTo(std::ostream* os) const {
  *os << "has a status code that ";
  code_matcher_.DescribeNegationTo(os);
  *os << ", or has an error message that ";
  message_matcher_.DescribeNegationTo(os);
}
bool StatusIsMatcherCommonImpl::MatchAndExplain(
    const absl::Status& status,
    ::testing::MatchResultListener* result_listener) const {
  ::testing::StringMatchResultListener inner_listener;
  inner_listener.Clear();
  if (!code_matcher_.MatchAndExplain(
          static_cast<absl::StatusCode>(status.code()), &inner_listener)) {
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
}  
}  