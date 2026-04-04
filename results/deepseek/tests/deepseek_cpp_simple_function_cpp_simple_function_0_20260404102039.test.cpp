#include "gtest/gtest.h"
#include "common/kind.h"
#include "absl/strings/string_view.h"

namespace cel {
namespace {

TEST(KindToStringTest, ReturnsCorrectStringForNullType) {
  EXPECT_EQ(KindToString(Kind::kNullType), "null_type");
}

TEST(KindToStringTest, ReturnsCorrectStringForDyn) {
  EXPECT_EQ(KindToString(Kind::kDyn), "dyn");
}

TEST(KindToStringTest, ReturnsCorrectStringForAny) {
  EXPECT_EQ(KindToString(Kind::kAny), "any");
}

TEST(KindToStringTest, ReturnsCorrectStringForType) {
  EXPECT_EQ(KindToString(Kind::kType), "type");
}

TEST(KindToStringTest, ReturnsCorrectStringForTypeParam) {
  EXPECT_EQ(KindToString(Kind::kTypeParam), "type_param");
}

TEST(KindToStringTest, ReturnsCorrectStringForFunction) {
  EXPECT_EQ(KindToString(Kind::kFunction), "function");
}

TEST(KindToStringTest, ReturnsCorrectStringForBool) {
  EXPECT_EQ(KindToString(Kind::kBool), "bool");
}

TEST(KindToStringTest, ReturnsCorrectStringForInt) {
  EXPECT_EQ(KindToString(Kind::kInt), "int");
}

TEST(KindToStringTest, ReturnsCorrectStringForUint) {
  EXPECT_EQ(KindToString(Kind::kUint), "uint");
}

TEST(KindToStringTest, ReturnsCorrectStringForDouble) {
  EXPECT_EQ(KindToString(Kind::kDouble), "double");
}

TEST(KindToStringTest, ReturnsCorrectStringForString) {
  EXPECT_EQ(KindToString(Kind::kString), "string");
}

TEST(KindToStringTest, ReturnsCorrectStringForBytes) {
  EXPECT_EQ(KindToString(Kind::kBytes), "bytes");
}

TEST(KindToStringTest, ReturnsCorrectStringForDuration) {
  EXPECT_EQ(KindToString(Kind::kDuration), "duration");
}

TEST(KindToStringTest, ReturnsCorrectStringForTimestamp) {
  EXPECT_EQ(KindToString(Kind::kTimestamp), "timestamp");
}

TEST(KindToStringTest, ReturnsCorrectStringForList) {
  EXPECT_EQ(KindToString(Kind::kList), "list");
}

TEST(KindToStringTest, ReturnsCorrectStringForMap) {
  EXPECT_EQ(KindToString(Kind::kMap), "map");
}

TEST(KindToStringTest, ReturnsCorrectStringForStruct) {
  EXPECT_EQ(KindToString(Kind::kStruct), "struct");
}

TEST(KindToStringTest, ReturnsCorrectStringForUnknown) {
  EXPECT_EQ(KindToString(Kind::kUnknown), "*unknown*");
}

TEST(KindToStringTest, ReturnsCorrectStringForOpaque) {
  EXPECT_EQ(KindToString(Kind::kOpaque), "*opaque*");
}

TEST(KindToStringTest, ReturnsCorrectStringForBoolWrapper) {
  EXPECT_EQ(KindToString(Kind::kBoolWrapper), "google.protobuf.BoolValue");
}

TEST(KindToStringTest, ReturnsCorrectStringForIntWrapper) {
  EXPECT_EQ(KindToString(Kind::kIntWrapper), "google.protobuf.Int64Value");
}

TEST(KindToStringTest, ReturnsCorrectStringForUintWrapper) {
  EXPECT_EQ(KindToString(Kind::kUintWrapper), "google.protobuf.UInt64Value");
}

TEST(KindToStringTest, ReturnsCorrectStringForDoubleWrapper) {
  EXPECT_EQ(KindToString(Kind::kDoubleWrapper), "google.protobuf.DoubleValue");
}

TEST(KindToStringTest, ReturnsCorrectStringForStringWrapper) {
  EXPECT_EQ(KindToString(Kind::kStringWrapper), "google.protobuf.StringValue");
}

TEST(KindToStringTest, ReturnsCorrectStringForBytesWrapper) {
  EXPECT_EQ(KindToString(Kind::kBytesWrapper), "google.protobuf.BytesValue");
}

TEST(KindToStringTest, ReturnsErrorForInvalidKind) {
  // Test with an invalid enum value that's not in the switch statement
  EXPECT_EQ(KindToString(static_cast<Kind>(999)), "*error*");
}

TEST(KindToStringTest, ReturnsErrorForNegativeKind) {
  // Test with a negative enum value
  EXPECT_EQ(KindToString(static_cast<Kind>(-1)), "*error*");
}

TEST(KindToStringTest, ReturnsErrorForMaxEnumValuePlusOne) {
  // Test with value just beyond the last valid enum value
  // Assuming Kind::kBytesWrapper is the last valid value
  int last_valid = static_cast<int>(Kind::kBytesWrapper);
  EXPECT_EQ(KindToString(static_cast<Kind>(last_valid + 1)), "*error*");
}

TEST(KindToStringTest, ReturnsNonEmptyStringForAllValidKinds) {
  // Verify that all valid kinds return non-empty strings
  for (int i = 0; i <= static_cast<int>(Kind::kBytesWrapper); ++i) {
    Kind kind = static_cast<Kind>(i);
    absl::string_view result = KindToString(kind);
    EXPECT_FALSE(result.empty()) << "Kind " << i << " returned empty string";
  }
}

TEST(KindToStringTest, ReturnsStringViewForAllCases) {
  // Verify function returns absl::string_view for all code paths
  absl::string_view result = KindToString(Kind::kNullType);
  EXPECT_EQ(result, "null_type");
  
  result = KindToString(static_cast<Kind>(999));
  EXPECT_EQ(result, "*error*");
}

}  // namespace
}  // namespace cel