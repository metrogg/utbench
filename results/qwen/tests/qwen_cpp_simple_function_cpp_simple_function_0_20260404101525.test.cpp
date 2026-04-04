#include <gtest/gtest.h>
#include "absl/strings/string_view.h"
#include "common/kind.h"

namespace cel {
// Forward declaration to satisfy compilation in isolated test environments.
absl::string_view KindToString(Kind kind);
}  // namespace cel

struct KindToStringTestCase {
  cel::Kind kind;
  const char* expected_str;
};

class KindToStringTest : public ::testing::TestWithParam<KindToStringTestCase> {};

TEST_P(KindToStringTest, MapsValidKindToCorrectString) {
  EXPECT_EQ(GetParam().expected_str, cel::KindToString(GetParam().kind));
}

INSTANTIATE_TEST_SUITE_P(
    AllValidKinds,
    KindToStringTest,
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
        KindToStringTestCase{cel::Kind::kBytesWrapper, "google.protobuf.BytesValue"}
    )
);

TEST(KindToStringTest, ReturnsErrorForInvalidOrOutOfRangeKind) {
  // Test lower boundary out-of-range value
  cel::Kind invalid_low = static_cast<cel::Kind>(-1);
  EXPECT_EQ("*error*", cel::KindToString(invalid_low));

  // Test upper boundary out-of-range value
  cel::Kind invalid_high = static_cast<cel::Kind>(9999);
  EXPECT_EQ("*error*", cel::KindToString(invalid_high));
}