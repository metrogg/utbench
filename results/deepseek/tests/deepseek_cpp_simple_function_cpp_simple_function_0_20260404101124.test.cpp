#include "common/kind.h"
#include "absl/strings/string_view.h"
#include "gtest/gtest.h"

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
  EXPECT_EQ(KindToString(static_cast<Kind>(999)), "*error*");
}

TEST(KindToStringTest, ReturnsErrorForNegativeKind) {
  EXPECT_EQ(KindToString(static_cast<Kind>(-1)), "*error*");
}

TEST(KindToStringTest, ReturnsErrorForKindJustBeforeValidRange) {
  EXPECT_EQ(KindToString(static_cast<Kind>(0)), "*error*");
}

TEST(KindToStringTest, ReturnsErrorForKindJustAfterValidRange) {
  EXPECT_EQ(KindToString(static_cast<Kind>(1000)), "*error*");
}

TEST(KindToStringTest, ReturnsNonEmptyStringForAllValidKinds) {
  EXPECT_FALSE(KindToString(Kind::kNullType).empty());
  EXPECT_FALSE(KindToString(Kind::kDyn).empty());
  EXPECT_FALSE(KindToString(Kind::kAny).empty());
  EXPECT_FALSE(KindToString(Kind::kType).empty());
  EXPECT_FALSE(KindToString(Kind::kTypeParam).empty());
  EXPECT_FALSE(KindToString(Kind::kFunction).empty());
  EXPECT_FALSE(KindToString(Kind::kBool).empty());
  EXPECT_FALSE(KindToString(Kind::kInt).empty());
  EXPECT_FALSE(KindToString(Kind::kUint).empty());
  EXPECT_FALSE(KindToString(Kind::kDouble).empty());
  EXPECT_FALSE(KindToString(Kind::kString).empty());
  EXPECT_FALSE(KindToString(Kind::kBytes).empty());
  EXPECT_FALSE(KindToString(Kind::kDuration).empty());
  EXPECT_FALSE(KindToString(Kind::kTimestamp).empty());
  EXPECT_FALSE(KindToString(Kind::kList).empty());
  EXPECT_FALSE(KindToString(Kind::kMap).empty());
  EXPECT_FALSE(KindToString(Kind::kStruct).empty());
  EXPECT_FALSE(KindToString(Kind::kUnknown).empty());
  EXPECT_FALSE(KindToString(Kind::kOpaque).empty());
  EXPECT_FALSE(KindToString(Kind::kBoolWrapper).empty());
  EXPECT_FALSE(KindToString(Kind::kIntWrapper).empty());
  EXPECT_FALSE(KindToString(Kind::kUintWrapper).empty());
  EXPECT_FALSE(KindToString(Kind::kDoubleWrapper).empty());
  EXPECT_FALSE(KindToString(Kind::kStringWrapper).empty());
  EXPECT_FALSE(KindToString(Kind::kBytesWrapper).empty());
}

TEST(KindToStringTest, ReturnsCorrectStringViewType) {
  auto result = KindToString(Kind::kString);
  EXPECT_EQ(typeid(result), typeid(absl::string_view));
}

}  // namespace
}  // namespace cel