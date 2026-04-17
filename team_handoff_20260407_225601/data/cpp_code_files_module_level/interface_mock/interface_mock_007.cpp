#ifndef THIRD_PARTY_CEL_CPP_EVAL_PUBLIC_TESTING_MATCHERS_H_
#define THIRD_PARTY_CEL_CPP_EVAL_PUBLIC_TESTING_MATCHERS_H_
#include <ostream>
#include "google/protobuf/message.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"
#include "absl/strings/string_view.h"
#include "absl/time/time.h"
#include "eval/public/cel_value.h"
#include "eval/public/set_util.h"
#include "eval/public/unknown_set.h"
namespace google {
namespace api {
namespace expr {
namespace runtime {
void PrintTo(const CelValue& value, std::ostream* os);
namespace test {
using CelValueMatcher = testing::Matcher<CelValue>;
CelValueMatcher EqualsCelValue(const CelValue& v);
CelValueMatcher IsCelNull();
CelValueMatcher IsCelBool(testing::Matcher<bool> m);
CelValueMatcher IsCelInt64(testing::Matcher<int64_t> m);
CelValueMatcher IsCelUint64(testing::Matcher<uint64_t> m);
CelValueMatcher IsCelDouble(testing::Matcher<double> m);
CelValueMatcher IsCelString(testing::Matcher<absl::string_view> m);
CelValueMatcher IsCelBytes(testing::Matcher<absl::string_view> m);
CelValueMatcher IsCelMessage(testing::Matcher<const google::protobuf::Message*> m);
CelValueMatcher IsCelDuration(testing::Matcher<absl::Duration> m);
CelValueMatcher IsCelTimestamp(testing::Matcher<absl::Time> m);
CelValueMatcher IsCelError(testing::Matcher<absl::Status> m);
template <typename ContainerMatcher>
class CelListMatcher : public testing::MatcherInterface<const CelValue&> {
 public:
  explicit CelListMatcher(ContainerMatcher m) : container_matcher_(m) {}
  bool MatchAndExplain(const CelValue& v,
                       testing::MatchResultListener* listener) const override {
    const CelList* cel_list;
    if (!v.GetValue(&cel_list) || cel_list == nullptr) return false;
    std::vector<CelValue> cel_vector;
    cel_vector.reserve(cel_list->size());
    for (int i = 0; i < cel_list->size(); ++i) {
      cel_vector.push_back((*cel_list)[i]);
    }
    return container_matcher_.Matches(cel_vector);
  }
  void DescribeTo(std::ostream* os) const override {
    CelValue::Type type =
        static_cast<CelValue::Type>(CelValue::IndexOf<const CelList*>::value);
    *os << absl::StrCat("type is ", CelValue::TypeName(type), " and ");
    container_matcher_.DescribeTo(os);
  }
 private:
  const testing::Matcher<std::vector<CelValue>> container_matcher_;
};
template <typename ContainerMatcher>
CelValueMatcher IsCelList(ContainerMatcher m) {
  return CelValueMatcher(new CelListMatcher(m));
}
}  
}  
}  
}  
}  
#endif  
#include "eval/public/testing/matchers.h"
#include <ostream>
#include <utility>
#include "google/protobuf/message.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"
#include "absl/strings/string_view.h"
#include "eval/public/set_util.h"
#include "internal/casts.h"
namespace google::api::expr::runtime {
void PrintTo(const CelValue& value, std::ostream* os) {
  *os << value.DebugString();
}
namespace test {
namespace {
using testing::_;
using testing::MatcherInterface;
using testing::MatchResultListener;
class CelValueEqualImpl : public MatcherInterface<CelValue> {
 public:
  explicit CelValueEqualImpl(const CelValue& v) : value_(v) {}
  bool MatchAndExplain(CelValue arg,
                       MatchResultListener* listener) const override {
    return CelValueEqual(arg, value_);
  }
  void DescribeTo(std::ostream* os) const override {
    *os << value_.DebugString();
  }
 private:
  const CelValue& value_;
};
template <typename UnderlyingType>
class CelValueMatcherImpl : public testing::MatcherInterface<const CelValue&> {
 public:
  explicit CelValueMatcherImpl(testing::Matcher<UnderlyingType> m)
      : underlying_type_matcher_(std::move(m)) {}
  bool MatchAndExplain(const CelValue& v,
                       testing::MatchResultListener* listener) const override {
    UnderlyingType arg;
    return v.GetValue(&arg) && underlying_type_matcher_.Matches(arg);
  }
  void DescribeTo(std::ostream* os) const override {
    CelValue::Type type =
        static_cast<CelValue::Type>(CelValue::IndexOf<UnderlyingType>::value);
    *os << absl::StrCat("type is ", CelValue::TypeName(type), " and ");
    underlying_type_matcher_.DescribeTo(os);
  }
 private:
  const testing::Matcher<UnderlyingType> underlying_type_matcher_;
};
template <>
class CelValueMatcherImpl<const google::protobuf::Message*>
    : public testing::MatcherInterface<const CelValue&> {
 public:
  explicit CelValueMatcherImpl(testing::Matcher<const google::protobuf::Message*> m)
      : underlying_type_matcher_(std::move(m)) {}
  bool MatchAndExplain(const CelValue& v,
                       testing::MatchResultListener* listener) const override {
    CelValue::MessageWrapper arg;
    return v.GetValue(&arg) && arg.HasFullProto() &&
           underlying_type_matcher_.Matches(
               cel::internal::down_cast<const google::protobuf::Message*>(
                   arg.message_ptr()));
  }
  void DescribeTo(std::ostream* os) const override {
    *os << absl::StrCat("type is ",
                        CelValue::TypeName(CelValue::Type::kMessage), " and ");
    underlying_type_matcher_.DescribeTo(os);
  }
 private:
  const testing::Matcher<const google::protobuf::Message*> underlying_type_matcher_;
};
}  
CelValueMatcher EqualsCelValue(const CelValue& v) {
  return CelValueMatcher(new CelValueEqualImpl(v));
}
CelValueMatcher IsCelNull() {
  return CelValueMatcher(new CelValueMatcherImpl<CelValue::NullType>(_));
}
CelValueMatcher IsCelBool(testing::Matcher<bool> m) {
  return CelValueMatcher(new CelValueMatcherImpl<bool>(std::move(m)));
}
CelValueMatcher IsCelInt64(testing::Matcher<int64_t> m) {
  return CelValueMatcher(new CelValueMatcherImpl<int64_t>(std::move(m)));
}
CelValueMatcher IsCelUint64(testing::Matcher<uint64_t> m) {
  return CelValueMatcher(new CelValueMatcherImpl<uint64_t>(std::move(m)));
}
CelValueMatcher IsCelDouble(testing::Matcher<double> m) {
  return CelValueMatcher(new CelValueMatcherImpl<double>(std::move(m)));
}
CelValueMatcher IsCelString(testing::Matcher<absl::string_view> m) {
  return CelValueMatcher(new CelValueMatcherImpl<CelValue::StringHolder>(
      testing::Property(&CelValue::StringHolder::value, m)));
}
CelValueMatcher IsCelBytes(testing::Matcher<absl::string_view> m) {
  return CelValueMatcher(new CelValueMatcherImpl<CelValue::BytesHolder>(
      testing::Property(&CelValue::BytesHolder::value, m)));
}
CelValueMatcher IsCelMessage(testing::Matcher<const google::protobuf::Message*> m) {
  return CelValueMatcher(
      new CelValueMatcherImpl<const google::protobuf::Message*>(std::move(m)));
}
CelValueMatcher IsCelDuration(testing::Matcher<absl::Duration> m) {
  return CelValueMatcher(new CelValueMatcherImpl<absl::Duration>(std::move(m)));
}
CelValueMatcher IsCelTimestamp(testing::Matcher<absl::Time> m) {
  return CelValueMatcher(new CelValueMatcherImpl<absl::Time>(std::move(m)));
}
CelValueMatcher IsCelError(testing::Matcher<absl::Status> m) {
  return CelValueMatcher(
      new CelValueMatcherImpl<const google::api::expr::runtime::CelError*>(
          testing::AllOf(testing::NotNull(), testing::Pointee(m))));
}
}  
}  