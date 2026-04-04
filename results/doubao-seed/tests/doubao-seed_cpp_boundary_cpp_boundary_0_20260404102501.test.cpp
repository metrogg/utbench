#include <cstdint>
#include <cmath>
#include "gtest/gtest.h"
#include "absl/status/status.h"
#include "absl/time/time.h"
#include "google/api/expr/v1alpha1/syntax.pb.h"
#include "common/constant.h"
#include "extensions/protobuf/internal/constant.h"

namespace cel::extensions::protobuf_internal {
namespace {

using ::google::api::expr::v1alpha1::Constant;

TEST(ConstantToProtoTest, MonostateCase) {
  ConstantProto proto;
  cel::Constant constant;
  EXPECT_OK(ConstantToProto(constant, &proto));
  EXPECT_EQ(proto.constant_kind_case(), ConstantProto::CONSTANT_KIND_NOT_SET);
}

TEST(ConstantToProtoTest, NullValueCase) {
  ConstantProto proto;
  cel::Constant constant;
  constant.set_null_value();
  EXPECT_OK(ConstantToProto(constant, &proto));
  EXPECT_EQ(proto.null_value(), google::protobuf::NULL_VALUE);
}

TEST(ConstantToProtoTest, BoolValueCases) {
  ConstantProto proto;
  cel::Constant constant_true(true);
  EXPECT_OK(ConstantToProto(constant_true, &proto));
  EXPECT_TRUE(proto.bool_value());
  
  cel::Constant constant_false(false);
  EXPECT_OK(ConstantToProto(constant_false, &proto));
  EXPECT_FALSE(proto.bool_value());
}

TEST(ConstantToProtoTest, Int64ValueCases) {
  ConstantProto proto;
  constexpr int64_t kNormalVal = 123456789;
  cel::Constant normal(kNormalVal);
  EXPECT_OK(ConstantToProto(normal, &proto));
  EXPECT_EQ(proto.int64_value(), kNormalVal);
  
  cel::Constant min_val(INT64_MIN);
  EXPECT_OK(ConstantToProto(min_val, &proto));
  EXPECT_EQ(proto.int64_value(), INT64_MIN);
  
  cel::Constant max_val(INT64_MAX);
  EXPECT_OK(ConstantToProto(max_val, &proto));
  EXPECT_EQ(proto.int64_value(), INT64_MAX);
}

TEST(ConstantToProtoTest, Uint64ValueCases) {
  ConstantProto proto;
  constexpr uint64_t kNormalVal = 987654321ULL;
  cel::Constant normal(kNormalVal);
  EXPECT_OK(ConstantToProto(normal, &proto));
  EXPECT_EQ(proto.uint64_value(), kNormalVal);
  
  cel::Constant min_val(0ULL);
  EXPECT_OK(ConstantToProto(min_val, &proto));
  EXPECT_EQ(proto.uint64_value(), 0ULL);
  
  cel::Constant max_val(UINT64_MAX);
  EXPECT_OK(ConstantToProto(max_val, &proto));
  EXPECT_EQ(proto.uint64_value(), UINT64_MAX);
}

TEST(ConstantToProtoTest, DoubleValueCases) {
  ConstantProto proto;
  constexpr double kNormalVal = 3.1415926535;
  cel::Constant normal(kNormalVal);
  EXPECT_OK(ConstantToProto(normal, &proto));
  EXPECT_DOUBLE_EQ(proto.double_value(), kNormalVal);
  
  cel::Constant zero(0.0);
  EXPECT_OK(ConstantToProto(zero, &proto));
  EXPECT_DOUBLE_EQ(proto.double_value(), 0.0);
  
  cel::Constant inf(INFINITY);
  EXPECT_OK(ConstantToProto(inf, &proto));
  EXPECT_TRUE(std::isinf(proto.double_value()));
  
  cel::Constant nan(NAN);
  EXPECT_OK(ConstantToProto(nan, &proto));
  EXPECT_TRUE(std::isnan(proto.double_value()));
}

TEST(ConstantToProtoTest, StringValueCases) {
  ConstantProto proto;
  constexpr absl::string_view kEmptyStr = "";
  cel::Constant empty(kEmptyStr);
  EXPECT_OK(ConstantToProto(empty, &proto));
  EXPECT_EQ(proto.string_value(), kEmptyStr);
  
  constexpr absl::string_view kTestStr = "test_string_with_special_chars_!@#$%^&*()";
  cel::Constant normal(kTestStr);
  EXPECT_OK(ConstantToProto(normal, &proto));
  EXPECT_EQ(proto.string_value(), kTestStr);
}

TEST(ConstantToProtoTest, BytesValueCases) {
  ConstantProto proto;
  const BytesConstant empty = "";
  cel::Constant empty_const(empty);
  EXPECT_OK(ConstantToProto(empty_const, &proto));
  EXPECT_EQ(proto.bytes_value(), empty);
  
  const BytesConstant test_bytes = "\x00\x01\x02\x03\xFF\xFE\xFD";
  cel::Constant bytes_const(test_bytes);
  EXPECT_OK(ConstantToProto(bytes_const, &proto));
  EXPECT_EQ(proto.bytes_value(), test_bytes);
}

TEST(ConstantToProtoTest, DurationValidCase) {
  ConstantProto proto;
  const absl::Duration test_duration = absl::Seconds(12345) + absl::Nanoseconds(678);
  cel::Constant constant(test_duration);
  EXPECT_OK(ConstantToProto(constant, &proto));
  EXPECT_EQ(proto.duration_value().seconds(), 12345);
  EXPECT_EQ(proto.duration_value().nanos(), 678);
}

TEST(ConstantToProtoTest, TimestampValidCase) {
  ConstantProto proto;
  const absl::Time test_time = absl::UnixEpoch() + absl::Hours(24*365) + absl::Nanoseconds(12345);
  cel::Constant constant(test_time);
  EXPECT_OK(ConstantToProto(constant, &proto));
  const auto& ts = proto.timestamp_value();
  EXPECT_EQ(ts.seconds(), 31536000);
  EXPECT_EQ(ts.nanos(), 12345);
}

TEST(ConstantFromProtoTest, NotSetCase) {
  ConstantProto proto;
  cel::Constant constant;
  EXPECT_OK(ConstantFromProto(proto, constant));
  EXPECT_TRUE(absl::holds_alternative<absl::monostate>(constant.kind()));
}

TEST(ConstantFromProtoTest, NullValueCase) {
  ConstantProto proto;
  proto.set_null_value(google::protobuf::NULL_VALUE);
  cel::Constant constant;
  EXPECT_OK(ConstantFromProto(proto, constant));
  EXPECT_TRUE(absl::holds_alternative<std::nullptr_t>(constant.kind()));
}

TEST(ConstantFromProtoTest, BoolValueCases) {
  ConstantProto proto;
  cel::Constant constant;
  proto.set_bool_value(true);
  EXPECT_OK(ConstantFromProto(proto, constant));
  EXPECT_EQ(constant.bool_value(), true);
  
  proto.set_bool_value(false);
  EXPECT_OK(ConstantFromProto(proto, constant));
  EXPECT_EQ(constant.bool_value(), false);
}

TEST(ConstantFromProtoTest, Int64ValueCases) {
  ConstantProto proto;
  cel::Constant constant;
  proto.set_int64_value(INT64_MIN);
  EXPECT_OK(ConstantFromProto(proto, constant));
  EXPECT_EQ(constant.int_value(), INT64_MIN);
  
  proto.set_int64_value(INT64_MAX);
  EXPECT_OK(ConstantFromProto(proto, constant));
  EXPECT_EQ(constant.int_value(), INT64_MAX);
}

TEST(ConstantFromProtoTest, Uint64ValueCases) {
  ConstantProto proto;
  cel::Constant constant;
  proto.set_uint64_value(0ULL);
  EXPECT_OK(ConstantFromProto(proto, constant));
  EXPECT_EQ(constant.uint_value(), 0ULL);
  
  proto.set_uint64_value(UINT64_MAX);
  EXPECT_OK(ConstantFromProto(proto, constant));
  EXPECT_EQ(constant.uint_value(), UINT64_MAX);
}

TEST(ConstantFromProtoTest, DoubleValueCases) {
  ConstantProto proto;
  cel::Constant constant;
  proto.set_double_value(1.2345);
  EXPECT_OK(ConstantFromProto(proto, constant));
  EXPECT_DOUBLE_EQ(constant.double_value(), 1.2345);
  
  proto.set_double_value(NAN);
  EXPECT_OK(ConstantFromProto(proto, constant));
  EXPECT_TRUE(std::isnan(constant.double_value()));
}

TEST(ConstantFromProtoTest, StringValueCases) {
  ConstantProto proto;
  cel::Constant constant;
  proto.set_string_value("");
  EXPECT_OK(ConstantFromProto(proto, constant));
  EXPECT_EQ(constant.string_value(), "");
  
  proto.set_string_value("test_from_proto_string");
  EXPECT_OK(ConstantFromProto(proto, constant));
  EXPECT_EQ(constant.string_value(), "test_from_proto_string");
}

TEST(ConstantFromProtoTest, BytesValueCases) {
  ConstantProto proto;
  cel::Constant constant;
  proto.set_bytes_value("");
  EXPECT_OK(ConstantFromProto(proto, constant));
  EXPECT_EQ(constant.bytes_value(), "");
  
  proto.set_bytes_value("\x00\xFF\x11\x22");
  EXPECT_OK(ConstantFromProto(proto, constant));
  EXPECT_EQ(constant.bytes_value(), "\x00\xFF\x11\x22");
}

TEST(ConstantFromProtoTest, DurationValidCase) {
  ConstantProto proto;
  auto* duration = proto.mutable_duration_value();
  duration->set_seconds(9999);
  duration->set_nanos(123456789);
  cel::Constant constant;
  EXPECT_OK(ConstantFromProto(proto, constant));
  EXPECT_EQ(constant.duration_value(), absl::Seconds(9999) + absl::Nanoseconds(123456789));
}

TEST(ConstantFromProtoTest, TimestampValidCase) {
  ConstantProto proto;
  auto* ts = proto.mutable_timestamp_value();
  ts->set_seconds(123456789);
  ts->set_nanos(987654321);
  cel::Constant constant;
  EXPECT_OK(ConstantFromProto(proto, constant));
  EXPECT_EQ(constant.timestamp_value(), absl::FromUnixSeconds(123456789) + absl::Nanoseconds(987654321));
}

TEST(ConstantFromProtoTest, UnknownKindError) {
  ConstantProto proto;
  const google::protobuf::Reflection* reflection = proto.GetReflection();
  const google::protobuf::OneofDescriptor* oneof = proto.GetDescriptor()->FindOneofByName("constant_kind");
  reflection->SetOneofCase(&proto, oneof, 9999);
  cel::Constant constant;
  absl::Status status = ConstantFromProto(proto, constant);
  EXPECT_EQ(status.code(), absl::StatusCode::kInvalidArgument);
  EXPECT_TRUE(absl::StrContains(status.message(), "unexpected ConstantKindCase: 9999"));
}

TEST(ConstantRoundTripTest, AllTypeRoundTrip) {
  // Null
  {
    cel::Constant original;
    original.set_null_value();
    ConstantProto proto;
    EXPECT_OK(ConstantToProto(original, &proto));
    cel::Constant converted;
    EXPECT_OK(ConstantFromProto(proto, converted));
    EXPECT_TRUE(absl::holds_alternative<std::nullptr_t>(converted.kind()));
  }
  // Bool
  {
    cel::Constant original(true);
    ConstantProto proto;
    EXPECT_OK(ConstantToProto(original, &proto));
    cel::Constant converted;
    EXPECT_OK(ConstantFromProto(proto, converted));
    EXPECT_EQ(converted.bool_value(), true);
  }
  // Int
  {
    cel::Constant original(INT64_MIN);
    ConstantProto proto;
    EXPECT_OK(ConstantToProto(original, &proto));
    cel::Constant converted;
    EXPECT_OK(ConstantFromProto(proto, converted));
    EXPECT_EQ(converted.int_value(), INT64_MIN);
  }
  // Uint
  {
    cel::Constant original(UINT64_MAX);
    ConstantProto proto;
    EXPECT_OK(ConstantToProto(original, &proto));
    cel::Constant converted;
    EXPECT_OK(ConstantFromProto(proto, converted));
    EXPECT_EQ(converted.uint_value(), UINT64_MAX);
  }
  // Double
  {
    cel::Constant original(3.14159);
    ConstantProto proto;
    EXPECT_OK(ConstantToProto(original, &proto));
    cel::Constant converted;
    EXPECT_OK(ConstantFromProto(proto, converted));
    EXPECT_DOUBLE_EQ(converted.double_value(), 3.14159);
  }
  // String
  {
    cel::Constant original("round_trip_test_string");
    ConstantProto proto;
    EXPECT_OK(ConstantToProto(original, &proto));
    cel::Constant converted;
    EXPECT_OK(ConstantFromProto(proto, converted));
    EXPECT_EQ(converted.string_value(), "round_trip_test_string");
  }
  // Bytes
  {
    cel::Constant original(BytesConstant("\x01\x02\x03\x04"));
    ConstantProto proto;
    EXPECT_OK(ConstantToProto(original, &proto));
    cel::Constant converted;
    EXPECT_OK(ConstantFromProto(proto, converted));
    EXPECT_EQ(converted.bytes_value(), "\x01\x02\x03\x04");
  }
  // Duration
  {
    cel::Constant original(absl::Minutes(123) + absl::Microseconds(456));
    ConstantProto proto;
    EXPECT_OK(ConstantToProto(original, &proto));
    cel::Constant converted;
    EXPECT_OK(ConstantFromProto(proto, converted));
    EXPECT_EQ(converted.duration_value(), absl::Minutes(123) + absl::Microseconds(456));
  }
  // Timestamp
  {
    cel::Constant original(absl::FromUnixMillis(1680000000000));
    ConstantProto proto;
    EXPECT_OK(ConstantToProto(original, &proto));
    cel::Constant converted;
    EXPECT_OK(ConstantFromProto(proto, converted));
    EXPECT_EQ(converted.timestamp_value(), absl::FromUnixMillis(1680000000000));
  }
}

}  // namespace
}  // namespace cel::extensions::protobuf_internal