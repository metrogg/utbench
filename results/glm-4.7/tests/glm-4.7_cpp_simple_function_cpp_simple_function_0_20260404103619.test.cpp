#include "gtest/gtest.h"
#include "base/kind.h"
#include "absl/strings/string_view.h"

namespace {

using KindToStringTestCase = std::pair<cel::Kind, absl::string_view>;

class KindToStringTest : public ::testing::TestWithParam<KindToStringTestCase> {};

TEST_P(KindToStringTest, ReturnsCorrectStringForKnownKinds) {
  const auto& [kind, expected] = GetParam();
  EXPECT_EQ(cel::KindToString(kind), expected);
}

INSTANTIATE_TEST_SUITE_P(
    KindToStringValidKinds, KindToStringTest,
    ::testing::Values(
        KindToStringTestCase{cel::Kind::kNullType, "null_type"},
        KindToStringTestCase{cel::Kind::kDyn, "dyn"},
        KindToStringTestCase{cel::Kind::kAny, "any"},
        KindToStringTestCase{cel::Kind::kType, "type"},
        KindToStringTestCase{cel::Kind::kTypeParam, "type_param"},
        KindToStringTestCase{cel::Kind::kFunction, "function"},
        KindToStringTestCase{cel::Kind::kBool, "bool"},
        KindToStringTestCase{cel::Kind::kInt, "int"},
        KindToStringTestCase{cel::Kind::kUint, "uint"},
        KindToStringTestCase{cel::Kind::kDouble, "double"},
        KindToStringTestCase{cel::Kind::kString, "string"},
        KindToStringTestCase{cel::Kind::kBytes, "bytes"},
        KindToStringTestCase{cel::Kind::kDuration, "duration"},
        KindToStringTestCase{cel::Kind::kTimestamp, "timestamp"},
        KindToStringTestCase{cel::Kind::kList, "list"},
        KindToStringTestCase{cel::Kind::kMap, "map"},
        KindToStringTestCase{cel::Kind::kStruct, "struct"},
        KindToStringTestCase{cel::Kind::kUnknown, "*unknown*"},
        KindToStringTestCase{cel::Kind::kOpaque, "*opaque*"},
        KindToStringTestCase{cel::Kind::kBoolWrapper, "google.protobuf.BoolValue"},
        KindToStringTestCase{cel::Kind::kIntWrapper, "google.protobuf.Int64Value"},
        KindToStringTestCase{cel::Kind::kUintWrapper, "google.protobuf.UInt64Value"},
        KindToStringTestCase{cel::Kind::kDoubleWrapper, "google.protobuf.DoubleValue"},
        KindToStringTestCase{cel::Kind::kStringWrapper, "google.protobuf.StringValue"},
        KindToStringTestCase{cel::Kind::kBytesWrapper, "google.protobuf.BytesValue"}),
    [](const testing::TestParamInfo<KindToStringTestCase>& info) {
      // Use the integer value of the Kind for the test case suffix
      return std::to_string(static_cast<int>(info.param.first));
    });

TEST(KindToStringTest, ReturnsErrorForInvalidKind) {
  // Note: Passing a value not present in the enum definition is technically
  // undefined behavior in C++, but it is a common practice to test the
  // 'default' case in a switch statement handling enums to ensure robustness.
  EXPECT_EQ(cel::KindToString(static_cast<cel::Kind>(-1)), "*error*");
  EXPECT_EQ(cel::KindToString(static_cast<cel::Kind>(9999)), "*error*");
}

}  // namespace