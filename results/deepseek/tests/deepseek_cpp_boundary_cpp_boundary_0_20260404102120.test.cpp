#include "extensions/protobuf/internal/constant.h"
#include "google/api/expr/v1alpha1/syntax.pb.h"
#include "absl/time/time.h"
#include "gtest/gtest.h"

namespace cel::extensions::protobuf_internal {
namespace {

using ConstantProto = google::api::expr::v1alpha1::Constant;

TEST(ConstantProtoTest, ConstantToProto_NullValue) {
  Constant constant;
  constant.set_null_value();
  ConstantProto proto;
  
  auto status = ConstantToProto(constant, &proto);
  EXPECT_TRUE(status.ok());
  EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kNullValue);
  EXPECT_EQ(proto.null_value(), google::protobuf::NULL_VALUE);
}

TEST(ConstantProtoTest, ConstantToProto_BoolTrue) {
  Constant constant;
  constant.set_bool_value(true);
  ConstantProto proto;
  
  auto status = ConstantToProto(constant, &proto);
  EXPECT_TRUE(status.ok());
  EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kBoolValue);
  EXPECT_TRUE(proto.bool_value());
}

TEST(ConstantProtoTest, ConstantToProto_BoolFalse) {
  Constant constant;
  constant.set_bool_value(false);
  ConstantProto proto;
  
  auto status = ConstantToProto(constant, &proto);
  EXPECT_TRUE(status.ok());
  EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kBoolValue);
  EXPECT_FALSE(proto.bool_value());
}

TEST(ConstantProtoTest, ConstantToProto_Int64Positive) {
  Constant constant;
  constant.set_int_value(1234567890);
  ConstantProto proto;
  
  auto status = ConstantToProto(constant, &proto);
  EXPECT_TRUE(status.ok());
  EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kInt64Value);
  EXPECT_EQ(proto.int64_value(), 1234567890);
}

TEST(ConstantProtoTest, ConstantToProto_Int64Negative) {
  Constant constant;
  constant.set_int_value(-987654321);
  ConstantProto proto;
  
  auto status = ConstantToProto(constant, &proto);
  EXPECT_TRUE(status.ok());
  EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kInt64Value);
  EXPECT_EQ(proto.int64_value(), -987654321);
}

TEST(ConstantProtoTest, ConstantToProto_Int64Zero) {
  Constant constant;
  constant.set_int_value(0);
  ConstantProto proto;
  
  auto status = ConstantToProto(constant, &proto);
  EXPECT_TRUE(status.ok());
  EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kInt64Value);
  EXPECT_EQ(proto.int64_value(), 0);
}

TEST(ConstantProtoTest, ConstantToProto_Uint64Positive) {
  Constant constant;
  constant.set_uint_value(18446744073709551615ULL);
  ConstantProto proto;
  
  auto status = ConstantToProto(constant, &proto);
  EXPECT_TRUE(status.ok());
  EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kUint64Value);
  EXPECT_EQ(proto.uint64_value(), 18446744073709551615ULL);
}

TEST(ConstantProtoTest, ConstantToProto_Uint64Zero) {
  Constant constant;
  constant.set_uint_value(0);
  ConstantProto proto;
  
  auto status = ConstantToProto(constant, &proto);
  EXPECT_TRUE(status.ok());
  EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kUint64Value);
  EXPECT_EQ(proto.uint64_value(), 0);
}

TEST(ConstantProtoTest, ConstantToProto_DoublePositive) {
  Constant constant;
  constant.set_double_value(3.141592653589793);
  ConstantProto proto;
  
  auto status = ConstantToProto(constant, &proto);
  EXPECT_TRUE(status.ok());
  EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kDoubleValue);
  EXPECT_DOUBLE_EQ(proto.double_value(), 3.141592653589793);
}

TEST(ConstantProtoTest, ConstantToProto_DoubleNegative) {
  Constant constant;
  constant.set_double_value(-2.718281828459045);
  ConstantProto proto;
  
  auto status = ConstantToProto(constant, &proto);
  EXPECT_TRUE(status.ok());
  EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kDoubleValue);
  EXPECT_DOUBLE_EQ(proto.double_value(), -2.718281828459045);
}

TEST(ConstantProtoTest, ConstantToProto_DoubleZero) {
  Constant constant;
  constant.set_double_value(0.0);
  ConstantProto proto;
  
  auto status = ConstantToProto(constant, &proto);
  EXPECT_TRUE(status.ok());
  EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kDoubleValue);
  EXPECT_DOUBLE_EQ(proto.double_value(), 0.0);
}

TEST(ConstantProtoTest, ConstantToProto_StringValue) {
  Constant constant;
  constant.set_string_value("Hello, World!");
  ConstantProto proto;
  
  auto status = ConstantToProto(constant, &proto);
  EXPECT_TRUE(status.ok());
  EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kStringValue);
  EXPECT_EQ(proto.string_value(), "Hello, World!");
}

TEST(ConstantProtoTest, ConstantToProto_StringEmpty) {
  Constant constant;
  constant.set_string_value("");
  ConstantProto proto;
  
  auto status = ConstantToProto(constant, &proto);
  EXPECT_TRUE(status.ok());
  EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kStringValue);
  EXPECT_EQ(proto.string_value(), "");
}

TEST(ConstantProtoTest, ConstantToProto_BytesValue) {
  Constant constant;
  constant.set_bytes_value("\x00\x01\x02\x03");
  ConstantProto proto;
  
  auto status = ConstantToProto(constant, &proto);
  EXPECT_TRUE(status.ok());
  EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kBytesValue);
  EXPECT_EQ(proto.bytes_value(), "\x00\x01\x02\x03");
}

TEST(ConstantProtoTest, ConstantToProto_BytesEmpty) {
  Constant constant;
  constant.set_bytes_value("");
  ConstantProto proto;
  
  auto status = ConstantToProto(constant, &proto);
  EXPECT_TRUE(status.ok());
  EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kBytesValue);
  EXPECT_EQ(proto.bytes_value(), "");
}

TEST(ConstantProtoTest, ConstantToProto_Duration) {
  Constant constant;
  constant.set_duration_value(absl::Seconds(3600) + absl::Milliseconds(500));
  ConstantProto proto;
  
  auto status = ConstantToProto(constant, &proto);
  EXPECT_TRUE(status.ok());
  EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kDurationValue);
}

TEST(ConstantProtoTest, ConstantToProto_Timestamp) {
  Constant constant;
  constant.set_timestamp_value(absl::FromUnixSeconds(1234567890));
  ConstantProto proto;
  
  auto status = ConstantToProto(constant, &proto);
  EXPECT_TRUE(status.ok());
  EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kTimestampValue);
}

TEST(ConstantProtoTest, ConstantToProto_EmptyConstant) {
  Constant constant;
  ConstantProto proto;
  
  auto status = ConstantToProto(constant, &proto);
  EXPECT_TRUE(status.ok());
  EXPECT_EQ(proto.constant_kind_case(), ConstantProto::CONSTANT_KIND_NOT_SET);
}

TEST(ConstantProtoTest, ConstantFromProto_NullValue) {
  ConstantProto proto;
  proto.set_null_value(google::protobuf::NULL_VALUE);
  Constant constant;
  
  auto status = ConstantFromProto(proto, constant);
  EXPECT_TRUE(status.ok());
  EXPECT_TRUE(constant.is_null_value());
}

TEST(ConstantProtoTest, ConstantFromProto_BoolTrue) {
  ConstantProto proto;
  proto.set_bool_value(true);
  Constant constant;
  
  auto status = ConstantFromProto(proto, constant);
  EXPECT_TRUE(status.ok());
  EXPECT_TRUE(constant.is_bool_value());
  EXPECT_TRUE(constant.bool_value());
}

TEST(ConstantProtoTest, ConstantFromProto_BoolFalse) {
  ConstantProto proto;
  proto.set_bool_value(false);
  Constant constant;
  
  auto status = ConstantFromProto(proto, constant);
  EXPECT_TRUE(status.ok());
  EXPECT_TRUE(constant.is_bool_value());
  EXPECT_FALSE(constant.bool_value());
}

TEST(ConstantProtoTest, ConstantFromProto_Int64Value) {
  ConstantProto proto;
  proto.set_int64_value(-123456789);
  Constant constant;
  
  auto status = ConstantFromProto(proto, constant);
  EXPECT_TRUE(status.ok());
  EXPECT_TRUE(constant.is_int_value());
  EXPECT_EQ(constant.int_value(), -123456789);
}

TEST(ConstantProtoTest, ConstantFromProto_Uint64Value) {
  ConstantProto proto;
  proto.set_uint64_value(9999999999ULL);
  Constant constant;
  
  auto status = ConstantFromProto(proto, constant);
  EXPECT_TRUE(status.ok());
  EXPECT_TRUE(constant.is_uint_value());
  EXPECT_EQ(constant.uint_value(), 9999999999ULL);
}

TEST(ConstantProtoTest, ConstantFromProto_DoubleValue) {
  ConstantProto proto;
  proto.set_double_value(1.23456789);
  Constant constant;
  
  auto status = ConstantFromProto(proto, constant);
  EXPECT_TRUE(status.ok());
  EXPECT_TRUE(constant.is_double_value());
  EXPECT_DOUBLE_EQ(constant.double_value(), 1.23456789);
}

TEST(ConstantProtoTest, ConstantFromProto_StringValue) {
  ConstantProto proto;
  proto.set_string_value("Test String");
  Constant constant;
  
  auto status = ConstantFromProto(proto, constant);
  EXPECT_TRUE(status.ok());
  EXPECT_TRUE(constant.is_string_value());
  EXPECT_EQ(constant.string_value(), "Test String");
}

TEST(ConstantProtoTest, ConstantFromProto_BytesValue) {
  ConstantProto proto;
  proto.set_bytes_value("\xFF\xFE\xFD");
  Constant constant;
  
  auto status = ConstantFromProto(proto, constant);
  EXPECT_TRUE(status.ok());
  EXPECT_TRUE(constant.is_bytes_value());
  EXPECT_EQ(constant.bytes_value(), "\xFF\xFE\xFD");
}

TEST(ConstantProtoTest, ConstantFromProto_DurationValue) {
  ConstantProto proto;
  proto.mutable_duration_value()->set_seconds(3600);
  proto.mutable_duration_value()->set_nanos(500000000);
  Constant constant;
  
  auto status = ConstantFromProto(proto, constant);
  EXPECT_TRUE(status.ok());
  EXPECT_TRUE(constant.is_duration_value());
}

TEST(ConstantProtoTest, ConstantFromProto_TimestampValue) {
  ConstantProto proto;
  proto.mutable_timestamp_value()->set_seconds(1234567890);
  proto.mutable_timestamp_value()->set_nanos(123456789);
  Constant constant;
  
  auto status = ConstantFromProto(proto, constant);
  EXPECT_TRUE(status.ok());
  EXPECT_TRUE(constant.is_timestamp_value());
}

TEST(ConstantProtoTest, ConstantFromProto_EmptyProto) {
  ConstantProto proto;
  Constant constant;
  
  auto status = ConstantFromProto(proto, constant);
  EXPECT_TRUE(status.ok());
}

TEST(ConstantProtoTest, ConstantFromProto_InvalidCase) {
  ConstantProto proto;
  Constant constant;
  
  proto.set_constant_kind_case(static_cast<ConstantProto::ConstantKindCase>(999));
  auto status = ConstantFromProto(proto, constant);
  EXPECT_FALSE(status.ok());
  EXPECT_EQ(status.code(), absl::StatusCode::kInvalidArgument);
}

TEST(ConstantProtoTest, RoundTrip_NullValue) {
  Constant original;
  original.set_null_value();
  
  ConstantProto proto;
  auto status1 = ConstantToProto(original, &proto);
  EXPECT_TRUE(status1.ok());
  
  Constant recovered;
  auto status2 = ConstantFromProto(proto, recovered);
  EXPECT_TRUE(status2.ok());
  EXPECT_TRUE(recovered.is_null_value());
}

TEST(ConstantProtoTest, RoundTrip_BoolValue) {
  Constant original;
  original.set_bool_value(true);
  
  ConstantProto proto;
  auto status1 = ConstantToProto(original, &proto);
  EXPECT_TRUE(status1.ok());
  
  Constant recovered;
  auto status2 = ConstantFromProto(proto, recovered);
  EXPECT_TRUE(status2.ok());
  EXPECT_TRUE(recovered.is_bool_value());
  EXPECT_TRUE(recovered.bool_value());
}

TEST(ConstantProtoTest, RoundTrip_Int64Value) {
  Constant original;
  original.set_int_value(-987654321);
  
  ConstantProto proto;
  auto status1 = ConstantToProto(original, &proto);
  EXPECT_TRUE(status1.ok());
  
  Constant recovered;
  auto status2 = ConstantFromProto(proto, recovered);
  EXPECT_TRUE(status2.ok());
  EXPECT_TRUE(recovered.is_int_value());
  EXPECT_EQ(recovered.int_value(), -987654321);
}

TEST(ConstantProtoTest, RoundTrip_Uint64Value) {
  Constant original;
  original.set_uint_value(18446744073709551615ULL);
  
  ConstantProto proto;
  auto status1 = ConstantToProto(original, &proto);
  EXPECT_TRUE(status1.ok());
  
  Constant recovered;
  auto status2 = ConstantFromProto(proto, recovered);
  EXPECT_TRUE(status2.ok());
  EXPECT_TRUE(recovered.is_uint_value());
  EXPECT_EQ(recovered.uint_value(), 18446744073709551615ULL);
}

TEST(ConstantProtoTest, RoundTrip_DoubleValue) {
  Constant original;
  original.set_double_value(2.718281828459045);
  
  ConstantProto proto;
  auto status1 = ConstantToProto(original, &proto);
  EXPECT_TRUE(status1.ok());
  
  Constant recovered;
  auto status2 = ConstantFromProto(proto, recovered);
  EXPECT_TRUE(status2.ok());
  EXPECT_TRUE(recovered.is_double_value());
  EXPECT_DOUBLE_EQ(recovered.double_value(), 2.718281828459045);
}

TEST(ConstantProtoTest, RoundTrip_StringValue) {
  Constant original;
  original.set_string_value("Round Trip Test");
  
  ConstantProto proto;
  auto status1 = ConstantToProto(original, &proto);
  EXPECT_TRUE(status1.ok());
  
  Constant recovered;
  auto status2 = ConstantFromProto(proto, recovered);
  EXPECT_TRUE(status2.ok());
  EXPECT_TRUE(recovered.is_string_value());
  EXPECT_EQ(recovered.string_value(), "Round Trip Test");
}

TEST(ConstantProtoTest, RoundTrip_BytesValue) {
  Constant original;
  original.set_bytes_value("\x01\x02\x03\x04\x05");
  
  ConstantProto proto;
  auto status1 = ConstantToProto(original, &proto);
  EXPECT_TRUE(status1.ok());
  
  Constant recovered;
  auto status2 = ConstantFromProto(proto, recovered);
  EXPECT_TRUE(status2.ok());
  EXPECT_TRUE(recovered.is_bytes_value());
  EXPECT_EQ(recovered.bytes_value(), "\x01\x02\x03\x04\x05");
}

}  // namespace
}  // namespace cel::extensions::protobuf_internal