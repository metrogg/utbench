#include <gtest/gtest.h>
#include "absl/strings/string_view.h"
#include "common/kind.h"

namespace cel {
namespace test {

TEST(KindToStringTest, ValidKindValuesReturnExpectedStrings) {
  EXPECT_EQ(KindToString(Kind::kNullType), "null_type");
  EXPECT_EQ(KindToString(Kind::kDyn), "dyn");
  EXPECT_EQ(KindToString(Kind::kAny), "any");
  EXPECT_EQ(KindToString(Kind::kType), "type");
  EXPECT_EQ(KindToString(Kind::kTypeParam), "type_param");
  EXPECT_EQ(KindToString(Kind::kFunction), "function");
  EXPECT_EQ(KindToString(Kind::kBool), "bool");
  EXPECT_EQ(KindToString(Kind::kInt), "int");
  EXPECT_EQ(KindToString(Kind::kUint), "uint");
  EXPECT_EQ(KindToString(Kind::kDouble), "double");
  EXPECT_EQ(KindToString(Kind::kString), "string");
  EXPECT_EQ(KindToString(Kind::kBytes), "bytes");
  EXPECT_EQ(KindToString(Kind::kDuration), "duration");
  EXPECT_EQ(KindToString(Kind::kTimestamp), "timestamp");
  EXPECT_EQ(KindToString(Kind::kList), "list");
  EXPECT_EQ(KindToString(Kind::kMap), "map");
  EXPECT_EQ(KindToString(Kind::kStruct), "struct");
  EXPECT_EQ(KindToString(Kind::kUnknown), "*unknown*");
  EXPECT_EQ(KindToString(Kind::kOpaque), "*opaque*");
  EXPECT_EQ(KindToString(Kind::kBoolWrapper), "google.protobuf.BoolValue");
  EXPECT_EQ(KindToString(Kind::kIntWrapper), "google.protobuf.Int64Value");
  EXPECT_EQ(KindToString(Kind::kUintWrapper), "google.protobuf.UInt64Value");
  EXPECT_EQ(KindToString(Kind::kDoubleWrapper), "google.protobuf.DoubleValue");
  EXPECT_EQ(KindToString(Kind::kStringWrapper), "google.protobuf.StringValue");
  EXPECT_EQ(KindToString(Kind::kBytesWrapper), "google.protobuf.BytesValue");
}

TEST(KindToStringTest, InvalidKindValueReturnsErrorString) {
  EXPECT_EQ(KindToString(static_cast<Kind>(-1)), "*error*");
  EXPECT_EQ(KindToString(static_cast<Kind>(255)), "*error*");
  EXPECT_EQ(KindToString(static_cast<Kind>(1000)), "*error*");
}

}  // namespace test
}  // namespace cel