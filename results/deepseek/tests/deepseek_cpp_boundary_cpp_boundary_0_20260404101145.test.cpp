#include "extensions/protobuf/internal/constant.h"
#include "google/api/expr/v1alpha1/syntax.pb.h"
#include "absl/time/time.h"
#include "gtest/gtest.h"

namespace cel::extensions::protobuf_internal {
namespace {

using ConstantProto = google::api::expr::v1alpha1::Constant;

TEST(ConstantProtoTest, ConstantToProto_HandlesNullValue) {
    Constant constant;
    constant.set_null_value();
    ConstantProto proto;
    
    auto status = ConstantToProto(constant, &proto);
    
    EXPECT_TRUE(status.ok());
    EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kNullValue);
    EXPECT_EQ(proto.null_value(), google::protobuf::NULL_VALUE);
}

TEST(ConstantProtoTest, ConstantToProto_HandlesBoolValue) {
    Constant constant;
    constant.set_bool_value(true);
    ConstantProto proto;
    
    auto status = ConstantToProto(constant, &proto);
    
    EXPECT_TRUE(status.ok());
    EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kBoolValue);
    EXPECT_TRUE(proto.bool_value());
}

TEST(ConstantProtoTest, ConstantToProto_HandlesInt64Value) {
    Constant constant;
    constant.set_int_value(INT64_MAX);
    ConstantProto proto;
    
    auto status = ConstantToProto(constant, &proto);
    
    EXPECT_TRUE(status.ok());
    EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kInt64Value);
    EXPECT_EQ(proto.int64_value(), INT64_MAX);
}

TEST(ConstantProtoTest, ConstantToProto_HandlesUint64Value) {
    Constant constant;
    constant.set_uint_value(UINT64_MAX);
    ConstantProto proto;
    
    auto status = ConstantToProto(constant, &proto);
    
    EXPECT_TRUE(status.ok());
    EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kUint64Value);
    EXPECT_EQ(proto.uint64_value(), UINT64_MAX);
}

TEST(ConstantProtoTest, ConstantToProto_HandlesDoubleValue) {
    Constant constant;
    constant.set_double_value(3.14159);
    ConstantProto proto;
    
    auto status = ConstantToProto(constant, &proto);
    
    EXPECT_TRUE(status.ok());
    EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kDoubleValue);
    EXPECT_DOUBLE_EQ(proto.double_value(), 3.14159);
}

TEST(ConstantProtoTest, ConstantToProto_HandlesStringValue) {
    Constant constant;
    constant.set_string_value("test_string");
    ConstantProto proto;
    
    auto status = ConstantToProto(constant, &proto);
    
    EXPECT_TRUE(status.ok());
    EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kStringValue);
    EXPECT_EQ(proto.string_value(), "test_string");
}

TEST(ConstantProtoTest, ConstantToProto_HandlesBytesValue) {
    Constant constant;
    constant.set_bytes_value("test_bytes");
    ConstantProto proto;
    
    auto status = ConstantToProto(constant, &proto);
    
    EXPECT_TRUE(status.ok());
    EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kBytesValue);
    EXPECT_EQ(proto.bytes_value(), "test_bytes");
}

TEST(ConstantProtoTest, ConstantToProto_HandlesDurationValue) {
    Constant constant;
    constant.set_duration_value(absl::Seconds(10));
    ConstantProto proto;
    
    auto status = ConstantToProto(constant, &proto);
    
    EXPECT_TRUE(status.ok());
    EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kDurationValue);
}

TEST(ConstantProtoTest, ConstantToProto_HandlesTimestampValue) {
    Constant constant;
    constant.set_timestamp_value(absl::UnixEpoch() + absl::Seconds(100));
    ConstantProto proto;
    
    auto status = ConstantToProto(constant, &proto);
    
    EXPECT_TRUE(status.ok());
    EXPECT_EQ(proto.constant_kind_case(), ConstantProto::kTimestampValue);
}

TEST(ConstantProtoTest, ConstantToProto_HandlesEmptyConstant) {
    Constant constant;
    ConstantProto proto;
    
    auto status = ConstantToProto(constant, &proto);
    
    EXPECT_TRUE(status.ok());
    EXPECT_EQ(proto.constant_kind_case(), ConstantProto::CONSTANT_KIND_NOT_SET);
}

TEST(ConstantProtoTest, ConstantFromProto_HandlesNullValue) {
    ConstantProto proto;
    proto.set_null_value(google::protobuf::NULL_VALUE);
    Constant constant;
    
    auto status = ConstantFromProto(proto, constant);
    
    EXPECT_TRUE(status.ok());
    EXPECT_TRUE(constant.is_null_value());
}

TEST(ConstantProtoTest, ConstantFromProto_HandlesBoolValue) {
    ConstantProto proto;
    proto.set_bool_value(false);
    Constant constant;
    
    auto status = ConstantFromProto(proto, constant);
    
    EXPECT_TRUE(status.ok());
    EXPECT_TRUE(constant.is_bool_value());
    EXPECT_FALSE(constant.bool_value());
}

TEST(ConstantProtoTest, ConstantFromProto_HandlesInt64Value) {
    ConstantProto proto;
    proto.set_int64_value(INT64_MIN);
    Constant constant;
    
    auto status = ConstantFromProto(proto, constant);
    
    EXPECT_TRUE(status.ok());
    EXPECT_TRUE(constant.is_int_value());
    EXPECT_EQ(constant.int_value(), INT64_MIN);
}

TEST(ConstantProtoTest, ConstantFromProto_HandlesUint64Value) {
    ConstantProto proto;
    proto.set_uint64_value(0);
    Constant constant;
    
    auto status = ConstantFromProto(proto, constant);
    
    EXPECT_TRUE(status.ok());
    EXPECT_TRUE(constant.is_uint_value());
    EXPECT_EQ(constant.uint_value(), 0);
}

TEST(ConstantProtoTest, ConstantFromProto_HandlesDoubleValue) {
    ConstantProto proto;
    proto.set_double_value(-1.5);
    Constant constant;
    
    auto status = ConstantFromProto(proto, constant);
    
    EXPECT_TRUE(status.ok());
    EXPECT_TRUE(constant.is_double_value());
    EXPECT_DOUBLE_EQ(constant.double_value(), -1.5);
}

TEST(ConstantProtoTest, ConstantFromProto_HandlesStringValue) {
    ConstantProto proto;
    proto.set_string_value("");
    Constant constant;
    
    auto status = ConstantFromProto(proto, constant);
    
    EXPECT_TRUE(status.ok());
    EXPECT_TRUE(constant.is_string_value());
    EXPECT_EQ(constant.string_value(), "");
}

TEST(ConstantProtoTest, ConstantFromProto_HandlesBytesValue) {
    ConstantProto proto;
    proto.set_bytes_value("");
    Constant constant;
    
    auto status = ConstantFromProto(proto, constant);
    
    EXPECT_TRUE(status.ok());
    EXPECT_TRUE(constant.is_bytes_value());
    EXPECT_EQ(constant.bytes_value(), "");
}

TEST(ConstantProtoTest, ConstantFromProto_HandlesDurationValue) {
    ConstantProto proto;
    proto.mutable_duration_value()->set_seconds(30);
    proto.mutable_duration_value()->set_nanos(500);
    Constant constant;
    
    auto status = ConstantFromProto(proto, constant);
    
    EXPECT_TRUE(status.ok());
    EXPECT_TRUE(constant.is_duration_value());
}

TEST(ConstantProtoTest, ConstantFromProto_HandlesTimestampValue) {
    ConstantProto proto;
    proto.mutable_timestamp_value()->set_seconds(1000);
    proto.mutable_timestamp_value()->set_nanos(0);
    Constant constant;
    
    auto status = ConstantFromProto(proto, constant);
    
    EXPECT_TRUE(status.ok());
    EXPECT_TRUE(constant.is_timestamp_value());
}

TEST(ConstantProtoTest, ConstantFromProto_HandlesEmptyProto) {
    ConstantProto proto;
    Constant constant;
    
    auto status = ConstantFromProto(proto, constant);
    
    EXPECT_TRUE(status.ok());
    EXPECT_TRUE(constant.is_null_value());
}

TEST(ConstantProtoTest, ConstantFromProto_ReturnsErrorForInvalidCase) {
    ConstantProto proto;
    Constant constant;
    
    proto.set_constant_kind_case(static_cast<ConstantProto::ConstantKindCase>(999));
    
    auto status = ConstantFromProto(proto, constant);
    
    EXPECT_FALSE(status.ok());
    EXPECT_EQ(status.code(), absl::StatusCode::kInvalidArgument);
}

TEST(ConstantProtoTest, RoundTripConversion_NullValue) {
    Constant original;
    original.set_null_value();
    ConstantProto proto;
    Constant result;
    
    auto status1 = ConstantToProto(original, &proto);
    auto status2 = ConstantFromProto(proto, result);
    
    EXPECT_TRUE(status1.ok());
    EXPECT_TRUE(status2.ok());
    EXPECT_TRUE(result.is_null_value());
}

TEST(ConstantProtoTest, RoundTripConversion_Int64Value) {
    Constant original;
    original.set_int_value(42);
    ConstantProto proto;
    Constant result;
    
    auto status1 = ConstantToProto(original, &proto);
    auto status2 = ConstantFromProto(proto, result);
    
    EXPECT_TRUE(status1.ok());
    EXPECT_TRUE(status2.ok());
    EXPECT_TRUE(result.is_int_value());
    EXPECT_EQ(result.int_value(), 42);
}

TEST(ConstantProtoTest, RoundTripConversion_StringValue) {
    Constant original;
    original.set_string_value("round_trip_test");
    ConstantProto proto;
    Constant result;
    
    auto status1 = ConstantToProto(original, &proto);
    auto status2 = ConstantFromProto(proto, result);
    
    EXPECT_TRUE(status1.ok());
    EXPECT_TRUE(status2.ok());
    EXPECT_TRUE(result.is_string_value());
    EXPECT_EQ(result.string_value(), "round_trip_test");
}

TEST(ConstantProtoTest, ConstantToProto_WithNullPointer) {
    Constant constant;
    constant.set_bool_value(true);
    
    auto status = ConstantToProto(constant, nullptr);
    
    EXPECT_FALSE(status.ok());
}

TEST(ConstantProtoTest, BoundaryValues_Int64) {
    Constant constant_min;
    constant_min.set_int_value(INT64_MIN);
    Constant constant_max;
    constant_max.set_int_value(INT64_MAX);
    ConstantProto proto_min, proto_max;
    
    auto status_min = ConstantToProto(constant_min, &proto_min);
    auto status_max = ConstantToProto(constant_max, &proto_max);
    
    EXPECT_TRUE(status_min.ok());
    EXPECT_TRUE(status_max.ok());
    EXPECT_EQ(proto_min.int64_value(), INT64_MIN);
    EXPECT_EQ(proto_max.int64_value(), INT64_MAX);
}

TEST(ConstantProtoTest, BoundaryValues_Uint64) {
    Constant constant_zero;
    constant_zero.set_uint_value(0);
    Constant constant_max;
    constant_max.set_uint_value(UINT64_MAX);
    ConstantProto proto_zero, proto_max;
    
    auto status_zero = ConstantToProto(constant_zero, &proto_zero);
    auto status_max = ConstantToProto(constant_max, &proto_max);
    
    EXPECT_TRUE(status_zero.ok());
    EXPECT_TRUE(status_max.ok());
    EXPECT_EQ(proto_zero.uint64_value(), 0);
    EXPECT_EQ(proto_max.uint64_value(), UINT64_MAX);
}

TEST(ConstantProtoTest, BoundaryValues_Double) {
    Constant constant_inf;
    constant_inf.set_double_value(std::numeric_limits<double>::infinity());
    Constant constant_nan;
    constant_nan.set_double_value(std::numeric_limits<double>::quiet_NaN());
    ConstantProto proto_inf, proto_nan;
    
    auto status_inf = ConstantToProto(constant_inf, &proto_inf);
    auto status_nan = ConstantToProto(constant_nan, &proto_nan);
    
    EXPECT_TRUE(status_inf.ok());
    EXPECT_TRUE(status_nan.ok());
    EXPECT_TRUE(std::isinf(proto_inf.double_value()));
    EXPECT_TRUE(std::isnan(proto_nan.double_value()));
}

}  // namespace
}  // namespace cel::extensions::protobuf_internal