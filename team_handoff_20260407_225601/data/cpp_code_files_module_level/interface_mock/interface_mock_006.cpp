#ifndef THIRD_PARTY_CEL_CPP_EXTENSIONS_PROTOBUF_VALUE_TESTING_H_
#define THIRD_PARTY_CEL_CPP_EXTENSIONS_PROTOBUF_VALUE_TESTING_H_
#include <ostream>
#include "absl/status/status.h"
#include "common/value.h"
#include "extensions/protobuf/internal/message.h"
#include "extensions/protobuf/value.h"
#include "internal/testing.h"
namespace cel::extensions::test {
template <typename MessageType>
class StructValueAsProtoMatcher {
 public:
  using is_gtest_matcher = void;
  explicit StructValueAsProtoMatcher(testing::Matcher<MessageType>&& m)
      : m_(std::move(m)) {}
  bool MatchAndExplain(cel::Value v,
                       testing::MatchResultListener* result_listener) const {
    MessageType msg;
    absl::Status s = ProtoMessageFromValue(v, msg);
    if (!s.ok()) {
      *result_listener << "cannot convert to "
                       << MessageType::descriptor()->full_name() << ": " << s;
      return false;
    }
    return m_.MatchAndExplain(msg, result_listener);
  }
  void DescribeTo(std::ostream* os) const {
    *os << "matches proto message " << m_;
  }
  void DescribeNegationTo(std::ostream* os) const {
    *os << "does not match proto message " << m_;
  }
 private:
  testing::Matcher<MessageType> m_;
};
template <typename MessageType>
inline StructValueAsProtoMatcher<MessageType> StructValueAsProto(
    testing::Matcher<MessageType>&& m) {
  static_assert(
      cel::extensions::protobuf_internal::IsProtoMessage<MessageType>);
  return StructValueAsProtoMatcher<MessageType>(std::move(m));
}
}  
#endif  
#include "common/value_testing.h"
#include <cstdint>
#include <ostream>
#include <string>
#include <utility>
#include "gtest/gtest.h"
#include "absl/status/status.h"
#include "absl/time/time.h"
#include "common/casting.h"
#include "common/value.h"
#include "common/value_kind.h"
#include "internal/testing.h"
namespace cel {
void PrintTo(const Value& value, std::ostream* os) { *os << value << "\n"; }
namespace test {
namespace {
using testing::Matcher;
template <typename Type>
constexpr ValueKind ToValueKind() {
  if constexpr (std::is_same_v<Type, BoolValue>) {
    return ValueKind::kBool;
  } else if constexpr (std::is_same_v<Type, IntValue>) {
    return ValueKind::kInt;
  } else if constexpr (std::is_same_v<Type, UintValue>) {
    return ValueKind::kUint;
  } else if constexpr (std::is_same_v<Type, DoubleValue>) {
    return ValueKind::kDouble;
  } else if constexpr (std::is_same_v<Type, StringValue>) {
    return ValueKind::kString;
  } else if constexpr (std::is_same_v<Type, BytesValue>) {
    return ValueKind::kBytes;
  } else if constexpr (std::is_same_v<Type, DurationValue>) {
    return ValueKind::kDuration;
  } else if constexpr (std::is_same_v<Type, TimestampValue>) {
    return ValueKind::kTimestamp;
  } else if constexpr (std::is_same_v<Type, ErrorValue>) {
    return ValueKind::kError;
  } else if constexpr (std::is_same_v<Type, MapValue>) {
    return ValueKind::kMap;
  } else if constexpr (std::is_same_v<Type, ListValue>) {
    return ValueKind::kList;
  } else if constexpr (std::is_same_v<Type, StructValue>) {
    return ValueKind::kStruct;
  } else if constexpr (std::is_same_v<Type, OpaqueValue>) {
    return ValueKind::kOpaque;
  } else {
    return ValueKind::kError;
  }
}
template <typename Type, typename NativeType>
class SimpleTypeMatcherImpl : public testing::MatcherInterface<const Value&> {
 public:
  using MatcherType = Matcher<NativeType>;
  explicit SimpleTypeMatcherImpl(MatcherType&& matcher)
      : matcher_(std::forward<MatcherType>(matcher)) {}
  bool MatchAndExplain(const Value& v,
                       testing::MatchResultListener* listener) const override {
    return InstanceOf<Type>(v) &&
           matcher_.MatchAndExplain(Cast<Type>(v).NativeValue(), listener);
  }
  void DescribeTo(std::ostream* os) const override {
    *os << absl::StrCat("kind is ", ValueKindToString(ToValueKind<Type>()),
                        " and ");
    matcher_.DescribeTo(os);
  }
 private:
  MatcherType matcher_;
};
template <typename Type>
class StringTypeMatcherImpl : public testing::MatcherInterface<const Value&> {
 public:
  using MatcherType = Matcher<std::string>;
  explicit StringTypeMatcherImpl(MatcherType matcher)
      : matcher_((std::move(matcher))) {}
  bool MatchAndExplain(const Value& v,
                       testing::MatchResultListener* listener) const override {
    return InstanceOf<Type>(v) && matcher_.Matches(Cast<Type>(v).ToString());
  }
  void DescribeTo(std::ostream* os) const override {
    *os << absl::StrCat("kind is ", ValueKindToString(ToValueKind<Type>()),
                        " and ");
    matcher_.DescribeTo(os);
  }
 private:
  MatcherType matcher_;
};
template <typename Type>
class AbstractTypeMatcherImpl : public testing::MatcherInterface<const Value&> {
 public:
  using MatcherType = Matcher<Type>;
  explicit AbstractTypeMatcherImpl(MatcherType&& matcher)
      : matcher_(std::forward<MatcherType>(matcher)) {}
  bool MatchAndExplain(const Value& v,
                       testing::MatchResultListener* listener) const override {
    return InstanceOf<Type>(v) && matcher_.Matches(Cast<Type>(v));
  }
  void DescribeTo(std::ostream* os) const override {
    *os << absl::StrCat("kind is ", ValueKindToString(ToValueKind<Type>()),
                        " and ");
    matcher_.DescribeTo(os);
  }
 private:
  MatcherType matcher_;
};
class OptionalValueMatcherImpl
    : public testing::MatcherInterface<const Value&> {
 public:
  explicit OptionalValueMatcherImpl(ValueMatcher matcher)
      : matcher_(std::move(matcher)) {}
  bool MatchAndExplain(const Value& v,
                       testing::MatchResultListener* listener) const override {
    if (!InstanceOf<OptionalValue>(v)) {
      *listener << "wanted OptionalValue, got " << ValueKindToString(v.kind());
      return false;
    }
    const auto& optional_value = Cast<OptionalValue>(v);
    if (!optional_value->HasValue()) {
      *listener << "OptionalValue is not engaged";
      return false;
    }
    return matcher_.MatchAndExplain(optional_value->Value(), listener);
  }
  void DescribeTo(std::ostream* os) const override {
    *os << "is OptionalValue that is engaged with value whose ";
    matcher_.DescribeTo(os);
  }
 private:
  ValueMatcher matcher_;
};
MATCHER(OptionalValueIsEmptyImpl, "is empty OptionalValue") {
  const Value& v = arg;
  if (!InstanceOf<OptionalValue>(v)) {
    *result_listener << "wanted OptionalValue, got "
                     << ValueKindToString(v.kind());
    return false;
  }
  const auto& optional_value = Cast<OptionalValue>(v);
  *result_listener << (optional_value.HasValue() ? "is not empty" : "is empty");
  return !optional_value->HasValue();
}
}  
ValueMatcher BoolValueIs(Matcher<bool> m) {
  return ValueMatcher(new SimpleTypeMatcherImpl<BoolValue, bool>(std::move(m)));
}
ValueMatcher IntValueIs(Matcher<int64_t> m) {
  return ValueMatcher(
      new SimpleTypeMatcherImpl<IntValue, int64_t>(std::move(m)));
}
ValueMatcher UintValueIs(Matcher<uint64_t> m) {
  return ValueMatcher(
      new SimpleTypeMatcherImpl<UintValue, uint64_t>(std::move(m)));
}
ValueMatcher DoubleValueIs(Matcher<double> m) {
  return ValueMatcher(
      new SimpleTypeMatcherImpl<DoubleValue, double>(std::move(m)));
}
ValueMatcher TimestampValueIs(Matcher<absl::Time> m) {
  return ValueMatcher(
      new SimpleTypeMatcherImpl<TimestampValue, absl::Time>(std::move(m)));
}
ValueMatcher DurationValueIs(Matcher<absl::Duration> m) {
  return ValueMatcher(
      new SimpleTypeMatcherImpl<DurationValue, absl::Duration>(std::move(m)));
}
ValueMatcher ErrorValueIs(Matcher<absl::Status> m) {
  return ValueMatcher(
      new SimpleTypeMatcherImpl<ErrorValue, absl::Status>(std::move(m)));
}
ValueMatcher StringValueIs(Matcher<std::string> m) {
  return ValueMatcher(new StringTypeMatcherImpl<StringValue>(std::move(m)));
}
ValueMatcher BytesValueIs(Matcher<std::string> m) {
  return ValueMatcher(new StringTypeMatcherImpl<BytesValue>(std::move(m)));
}
ValueMatcher MapValueIs(Matcher<MapValue> m) {
  return ValueMatcher(new AbstractTypeMatcherImpl<MapValue>(std::move(m)));
}
ValueMatcher ListValueIs(Matcher<ListValue> m) {
  return ValueMatcher(new AbstractTypeMatcherImpl<ListValue>(std::move(m)));
}
ValueMatcher StructValueIs(Matcher<StructValue> m) {
  return ValueMatcher(new AbstractTypeMatcherImpl<StructValue>(std::move(m)));
}
ValueMatcher OptionalValueIs(ValueMatcher m) {
  return ValueMatcher(new OptionalValueMatcherImpl(std::move(m)));
}
ValueMatcher OptionalValueIsEmpty() { return OptionalValueIsEmptyImpl(); }
}  
}  