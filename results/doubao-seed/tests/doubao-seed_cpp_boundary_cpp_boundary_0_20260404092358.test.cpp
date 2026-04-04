#include <gtest/gtest.h>
#include <gmock/gmock.h>
#include <cstdint>
#include <string>
#include <cmath>
#include <limits>
#include "absl/status/status.h"
#include "absl/time/time.h"
#include "absl/types/variant.h"
#include "google/api/expr/v1alpha1/syntax.pb.h"
#include "google/protobuf/descriptor.h"
#include "common/constant.h"
#include "extensions/protobuf/internal/constant.h"

namespace cel::extensions::protobuf_internal {
namespace {

using ::testing::HasSubstr;
using ConstantProto = ::google::api::expr::v1alpha1::Constant;

TEST(ConstantProtoTest, ToProtoMonostate) {
  Constant constant;
  ConstantProto proto;
  ASSERT_TRUE(ConstantToProto(constant, &proto).ok());
  EXPECT_EQ(proto.constant_kind_case(), ConstantProto::CONSTANT_KIND_NOT_SET);
}

TEST(ConstantProtoTest, ToProtoNull) {
  Constant constant(nullptr);
  ConstantProto proto;
  ASSERT_TRUE(ConstantToProto(constant, &proto).ok());
  EXPECT_EQ(proto.null_value(), ::google::protobuf::NULL_VALUE);
}

TEST(ConstantProtoTest, ToProtoBool) {
  Constant constant_true(true);
  ConstantProto proto_true;
  ASSERT_TRUE(ConstantToProto(constant_true, &proto_true).ok());
  EXPECT_TRUE(proto_true.bool_value());

  Constant constant_false(false);
  ConstantProto proto_false;
  ASSERT_TRUE(ConstantToProto(constant_false, &proto_false).ok());
  EXPECT_FALSE(proto_false.bool_value());
}

TEST(ConstantProtoTest, ToProtoInt64BoundaryValues) {
  Constant normal(INT64_C(123456789));
  ConstantProto proto_normal;
  ASSERT_TRUE(ConstantToProto(normal, &proto_normal).ok());
  EXPECT_EQ(proto_normal.int64_value(), 123456789);

  Constant min_val(INT64_MIN);
  ConstantProto proto_min;
  ASSERT_TRUE(ConstantToProto(min_val, &proto_min).ok());
  EXPECT_EQ(proto_min.int64_value(), INT64_MIN);

  Constant max_val(INT64_MAX);
  ConstantProto proto_max;
  ASSERT_TRUE(ConstantToProto(max_val, &proto_max).ok());
  EXPECT_EQ(proto_max.int64_value(), INT64_MAX);
}

TEST(ConstantProtoTest, ToProtoUint64BoundaryValues) {
  Constant zero(UINT64_C(0));
  ConstantProto proto_zero;
  ASSERT_TRUE(ConstantToProto(zero, &proto_zero).ok());
  EXPECT_EQ(proto_zero.uint64_value(), 0);

  Constant max_val(UINT64_MAX);
  ConstantProto proto_max;
  ASSERT_TRUE(ConstantToProto(max_val, &proto_max).ok());
  EXPECT_EQ(proto_max.uint64_value(), UINT64_MAX);
}

TEST(ConstantProtoTest, ToProtoDoubleSpecialValues) {
  Constant normal(123.456);
  ConstantProto proto_normal;
  ASSERT_TRUE(ConstantToProto(normal, &proto_normal).ok());
  EXPECT_DOUBLE_EQ(proto_normal.double_value(), 123.456);

  Constant nan_val(std::nan(""));
  ConstantProto proto_nan;
  ASSERT_TRUE(ConstantToProto(nan_val, &proto_nan).ok());
  EXPECT_TRUE(std::isnan(proto_nan.double_value()));

  Constant inf_val(std::numeric_limits<double>::infinity());
  ConstantProto proto_inf;
  ASSERT_TRUE(ConstantToProto(inf_val, &proto_inf).ok());
  EXPECT_TRUE(std::isinf(proto_inf.double_value()));
  EXPECT_GT(proto_inf.double_value(), 0);
}

TEST(ConstantProtoTest, ToProtoStringEmptyAndUnicode) {
  Constant empty("");
  ConstantProto proto_empty;
  ASSERT_TRUE(ConstantToProto(empty, &proto_empty).ok());
  EXPECT_EQ(proto_empty.string_value(), "");

  Constant unicode(u8"测试 🌍 hello!");
  ConstantProto proto_unicode;
  ASSERT_TRUE(ConstantToProto(unicode, &proto_unicode).ok());
  EXPECT_EQ(proto_unicode.string_value(), u8"测试 🌍 hello!");
}

TEST(ConstantProtoTest, ToProtoBytesWithNullChars) {
  Constant bytes(BytesConstant("\x00\x01\xFF\xFE\x00"));
  ConstantProto proto;
  ASSERT_TRUE(ConstantToProto(bytes, &proto).ok());
  EXPECT_EQ(proto.bytes_value(), "\x00\x01\xFF\xFE\x00");
}

TEST(ConstantProtoTest, ToProtoDuration) {
  absl::Duration test_duration = absl::Seconds(123) + absl::Nanoseconds(456789);
  Constant constant(test_duration);
  ConstantProto proto;
  ASSERT_TRUE(ConstantToProto(constant, &proto).ok());
  EXPECT_EQ(proto.duration_value().seconds(), 123);
  EXPECT_EQ(proto.duration_value().nanos(), 456789);
}

TEST(ConstantProtoTest, ToProtoTimestamp) {
  absl::Time test_time = absl::UnixEpoch() + absl::Hours(72) + absl::Nanoseconds(123456);
  Constant constant(test_time);
  ConstantProto proto;
  ASSERT_TRUE(ConstantToProto(constant, &proto).ok());
  EXPECT_EQ(proto.timestamp_value().seconds(), 72 * 3600);
  EXPECT_EQ(proto.timestamp_value().nanos(), 123456);
}

TEST(ConstantProtoTest, FromProtoNotSet) {
  ConstantProto proto;
  Constant constant;
  ASSERT_TRUE(ConstantFromProto(proto, constant).ok());
  EXPECT_TRUE(absl::holds_alternative<absl::monostate>(constant.kind()));
}

TEST(ConstantProtoTest, FromProtoNull) {
  ConstantProto proto;
  proto.set_null_value(::google::protobuf::NULL_VALUE);
  Constant constant;
  ASSERT_TRUE(ConstantFromProto(proto, constant).ok());
  EXPECT_TRUE(absl::holds_alternative<std::nullptr_t>(constant.kind()));
}

TEST(ConstantProtoTest, FromProtoBool) {
  ConstantProto proto_true;
  proto_true.set_bool_value(true);
  Constant constant_true;
  ASSERT_TRUE(ConstantFromProto(proto_true, constant_true).ok());
  EXPECT_TRUE(constant_true.bool_value());

  ConstantProto proto_false;
  proto_false.set_bool_value(false);
  Constant constant_false;
  ASSERT_TRUE(ConstantFromProto(proto_false, constant_false).ok());
  EXPECT_FALSE(constant_false.bool_value());
}

TEST(ConstantProtoTest, FromProtoInt64Boundary) {
  ConstantProto proto;
  proto.set_int64_value(INT64_MIN);
  Constant constant;
  ASSERT_TRUE(ConstantFromProto(proto, constant).ok());
  EXPECT_EQ(constant.int_value(), INT64_MIN);

  proto.set_int64_value(INT64_MAX);
  ASSERT_TRUE(ConstantFromProto(proto, constant).ok());
  EXPECT_EQ(constant.int_value(), INT64_MAX);
}

TEST(ConstantProtoTest, FromProtoUnknownKindReturnsError) {
  ConstantProto proto;
  const std::string serialized = "\xba\x1f\x04test";
  ASSERT_TRUE(proto.ParseFromString(serialized));
  Constant constant;
  absl::Status status = ConstantFromProto(proto, constant);
  EXPECT_EQ(status.code(), absl::StatusCode::kInvalidArgument);
  EXPECT_THAT(status.message(), HasSubstr("unexpected ConstantKindCase: 999"));
}

TEST(ConstantProtoTest, RoundTripMonostate) {
  Constant original;
  ConstantProto proto;
  ASSERT_TRUE(ConstantToProto(original, &proto).ok());
  Constant result;
  ASSERT_TRUE(ConstantFromProto(proto, result).ok());
  EXPECT_TRUE(absl::holds_alternative<absl::monostate>(result.kind()));
}

TEST(ConstantProtoTest, RoundTripNull) {
  Constant original(nullptr);
  ConstantProto proto;
  ASSERT_TRUE(ConstantToProto(original, &proto).ok());
  Constant result;
  ASSERT_TRUE(ConstantFromProto(proto, result).ok());
  EXPECT_TRUE(absl::holds_alternative<std::nullptr_t>(result.kind()));
}

TEST(ConstantProtoTest, RoundTripInt64) {
  Constant original(INT64_C(-987654321));
  ConstantProto proto;
  ASSERT_TRUE(ConstantToProto(original, &proto).ok());
  Constant result;
  ASSERT_TRUE(ConstantFromProto(proto, result).ok());
  EXPECT_EQ(result.int_value(), original.int_value());
}

TEST(ConstantProtoTest, RoundTripDoubleNaN) {
  Constant original(std::nan(""));
  ConstantProto proto;
  ASSERT_TRUE(ConstantToProto(original, &proto).ok());
  Constant result;
  ASSERT_TRUE(ConstantFromProto(proto, result).ok());
  EXPECT_TRUE(std::isnan(result.double_value()));
}

TEST(ConstantProtoTest, RoundTripDuration) {
  Constant original(absl::Minutes(500) + absl::Nanoseconds(123456789));
  ConstantProto proto;
  ASSERT_TRUE(ConstantToProto(original, &proto).ok());
  Constant result;
  ASSERT_TRUE(ConstantFromProto(proto, result).ok());
  EXPECT_EQ(result.duration_value(), original.duration_value());
}

TEST(ConstantProtoTest, RoundTripTimestamp) {
  Constant original(absl::UnixEpoch() + absl::Days(365) + absl::Seconds(12345));
  ConstantProto proto;
  ASSERT_TRUE(ConstantToProto(original, &proto).ok());
  Constant result;
  ASSERT_TRUE(ConstantFromProto(proto, result).ok());
  EXPECT_EQ(result.timestamp_value(), original.timestamp_value());
}

}  // namespace
}  // namespace cel::extensions::protobuf_internal