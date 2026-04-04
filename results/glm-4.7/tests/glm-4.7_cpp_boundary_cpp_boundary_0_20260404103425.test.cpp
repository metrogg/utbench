#include <gtest/gtest.h>
#include <google/protobuf/duration.pb.h>
#include <google/protobuf/timestamp.pb.h>
#include <google/protobuf/wrappers.pb.h>
#include <google/api/expr/v1alpha1/syntax.pb.h>
#include <absl/base/nullability.h>
#include <absl/status/status.h>
#include <absl/strings/str_cat.h>
#include <absl/time/time.h>
#include <absl/types/variant.h>
#include <cstddef>
#include <cstdint>
#include <string>
#include <variant>

// --- Stubs for Dependencies ---

namespace cel {
using BytesConstant = std::string;
using StringConstant = std::string;

struct Constant {
  using VariantType = std::variant<absl::monostate, std::nullptr_t, bool, int64_t,
                                   uint64_t, double, BytesConstant,
                                   StringConstant, absl::Duration, absl::Time>;
  VariantType kind_;

  Constant() : kind_(absl::monostate{}) {}
  Constant(std::nullptr_t) : kind_(nullptr) {}
  Constant(bool v) : kind_(v) {}
  Constant(int64_t v) : kind_(v) {}
  Constant(uint64_t v) : kind_(v) {}
  Constant(double v) : kind_(v) {}
  Constant(const BytesConstant& v) : kind_(v) {}
  Constant(const StringConstant& v) : kind_(v) {}
  Constant(absl::Duration v) : kind_(v) {}
  Constant(absl::Time v) : kind_(v) {}

  const VariantType& kind() const { return kind_; }

  void set_null_value() { kind_ = nullptr; }
  void set_bool_value(bool v) { kind_ = v; }
  void set_int_value(int64_t v) { kind_ = v; }
  void set_uint_value(uint64_t v) { kind_ = v; }
  void set_double_value(double v) { kind_ = v; }
  void set_string_value(const StringConstant& v) { kind_ = v; }
  void set_bytes_value(const BytesConstant& v) { kind_ = v; }
  void set_duration_value(absl::Duration v) { kind_ = v; }
  void set_timestamp_value(absl::Time v) { kind_ = v; }
};
} // namespace cel

namespace cel::extensions::protobuf_internal::internal {
absl::Status EncodeDuration(absl::Duration d,
                             google::protobuf::Duration* proto) {
  proto->set_seconds(absl::ToInt64Seconds(d));
  proto->set_nanos(
      absl::ToInt64Nanoseconds(d - absl::Seconds(absl::ToInt64Seconds(d))));
  return absl::OkStatus();
}

absl::Status EncodeTime(absl::Time t, google::protobuf::Timestamp* proto) {
  proto->set_seconds(absl::ToInt64Seconds(t - absl::UnixEpoch()));
  proto->set_nanos(absl::ToInt64Nanoseconds(
      (t - absl::UnixEpoch()) - absl::Seconds(absl::ToInt64Seconds(t - absl::UnixEpoch()))));
  return absl::OkStatus();
}

absl::Duration DecodeDuration(const google::protobuf::Duration& proto) {
  return absl::Seconds(proto.seconds()) + absl::Nanos(proto.nanos());
}

absl::Time DecodeTime(const google::protobuf::Timestamp& proto) {
  return absl::UnixEpoch() + absl::Seconds(proto.seconds()) +
         absl::Nanos(proto.nanos());
}
} // namespace cel::extensions::protobuf_internal::internal

// --- Source Code Under Test ---

namespace cel::extensions::protobuf_internal {
using ConstantProto = google::api::expr::v1alpha1::Constant;

absl::Status ConstantToProto(const Constant& constant,
                             absl::Nonnull<ConstantProto*> proto) {
  return absl::visit(
      absl::Overload(
          [proto](absl::monostate) -> absl::Status {
            proto->clear_constant_kind();
            return absl::OkStatus();
          },
          [proto](std::nullptr_t) -> absl::Status {
            proto->set_null_value(google::protobuf::NULL_VALUE);
            return absl::OkStatus();
          },
          [proto](bool value) -> absl::Status {
            proto->set_bool_value(value);
            return absl::OkStatus();
          },
          [proto](int64_t value) -> absl::Status {
            proto->set_int64_value(value);
            return absl::OkStatus();
          },
          [proto](uint64_t value) -> absl::Status {
            proto->set_uint64_value(value);
            return absl::OkStatus();
          },
          [proto](double value) -> absl::Status {
            proto->set_double_value(value);
            return absl::OkStatus();
          },
          [proto](const BytesConstant& value) -> absl::Status {
            proto->set_bytes_value(value);
            return absl::OkStatus();
          },
          [proto](const StringConstant& value) -> absl::Status {
            proto->set_string_value(value);
            return absl::OkStatus();
          },
          [proto](absl::Duration value) -> absl::Status {
            return internal::EncodeDuration(value,
                                           proto->mutable_duration_value());
          },
          [proto](absl::Time value) -> absl::Status {
            return internal::EncodeTime(value, proto->mutable_timestamp_value());
          }),
      constant.kind());
}

absl::Status ConstantFromProto(const ConstantProto& proto, Constant& constant) {
  switch (proto.constant_kind_case()) {
  case ConstantProto::CONSTANT_KIND_NOT_SET:
    constant = Constant{};
    break;
  case ConstantProto::kNullValue:
    constant.set_null_value();
    break;
  case ConstantProto::kBoolValue:
    constant.set_bool_value(proto.bool_value());
    break;
  case ConstantProto::kInt64Value:
    constant.set_int_value(proto.int64_value());
    break;
  case ConstantProto::kUint64Value:
    constant.set_uint_value(proto.uint64_value());
    break;
  case ConstantProto::kDoubleValue:
    constant.set_double_value(proto.double_value());
    break;
  case ConstantProto::kStringValue:
    constant.set_string_value(proto.string_value());
    break;
  case ConstantProto::kBytesValue:
    constant.set_bytes_value(proto.bytes_value());
    break;
  case ConstantProto::kDurationValue:
    constant.set_duration_value(internal::DecodeDuration(proto.duration_value()));
    break;
  case ConstantProto::kTimestampValue:
    constant.set_timestamp_value(internal::DecodeTime(proto.timestamp_value()));
    break;
  default:
    return absl::InvalidArgumentError(
        absl::StrCat("unexpected ConstantKindCase: ",
                     static_cast<int>(proto.constant_kind_case())));
  }
  return absl::OkStatus();
}
} // namespace cel::extensions::protobuf_internal

// --- Tests ---

namespace {
using namespace cel::extensions::protobuf_internal;
using namespace cel;

class ConstantConversionTest : public ::testing::Test {
protected:
  void SetUp() override {}
  void TearDown() override {}
};

TEST_F(ConstantConversionTest, ConstantToProto_Monostate) {
  Constant c;
  ConstantProto p;
  auto status = ConstantToProto(c, &p);
  ASSERT_TRUE(status.ok());
  EXPECT_EQ(ConstantProto::CONSTANT_KIND_NOT_SET, p.constant_kind_case());
}

TEST_F(ConstantConversionTest, ConstantToProto_Null) {
  Constant c(nullptr);
  ConstantProto p;
  auto status = ConstantToProto(c, &p);
  ASSERT_TRUE(status.ok());
  EXPECT_EQ(ConstantProto::kNullValue, p.constant_kind_case());
  EXPECT_EQ(google::protobuf::NULL_VALUE, p.null_value());
}

TEST ConstantConversionTest, ConstantToProto_Bool) {
  Constant c(true);
  ConstantProto p;
  auto status = ConstantToProto(c, &p);
  ASSERT_TRUE(status.ok());
  EXPECT_EQ(ConstantProto::kBoolValue, p.constant_kind_case());
  EXPECT_TRUE(p.bool_value());

  Constant c_false(false);
  ConstantProto p_false;
  status = ConstantToProto(c_false, &p_false);
  ASSERT_TRUE(status.ok());
  EXPECT_FALSE(p_false.bool_value());
}

TEST_F(ConstantConversionTest, ConstantToProto_Int64) {
  int64_t val = -12345678901234LL;
  Constant c(val);
  ConstantProto p;
  auto status = ConstantToProto(c, &p);
  ASSERT_TRUE(status.ok());
  EXPECT_EQ(ConstantProto::kInt64Value, p.constant_kind_case());
  EXPECT_EQ(val, p.int64_value());
}

TEST_F(ConstantConversionTest, ConstantToProto_Uint64) {
  uint64_t val = 12345678901234ULL;
  Constant c(val);
  ConstantProto p;
  auto status = ConstantToProto(c, &p);
  ASSERT_TRUE(status.ok());
  EXPECT_EQ(ConstantProto::kUint64Value, p.constant_kind_case());
  EXPECT_EQ(val, p.uint64_value());
}

TEST_F(ConstantConversionTest, ConstantToProto_Double) {
  double val = 3.1415926535;
  Constant c(val);
  ConstantProto p;
  auto status = ConstantToProto(c, &p);
  ASSERT_TRUE(status.ok());
  EXPECT_EQ(ConstantProto::kDoubleValue, p.constant_kind_case());
  EXPECT_DOUBLE_EQ(val, p.double_value());
}

TEST_F(ConstantConversionTest, ConstantToProto_Bytes) {
  BytesConstant val = "binary_data\x01\x02\x03";
  Constant c(val);
  ConstantProto p;
  auto status = ConstantToProto(c, &p);
  ASSERT_TRUE(status.ok());
  EXPECT_EQ(ConstantProto::kBytesValue, p.constant_kind_case());
  EXPECT_EQ(val, p.bytes_value());
}

TEST_F(ConstantConversionTest, ConstantToProto_String) {
  StringConstant val = "hello world";
  Constant c(val);
  ConstantProto p;
  auto status = ConstantToProto(c, &p);
  ASSERT_TRUE(status.ok());
  EXPECT_EQ(ConstantProto::kStringValue, p.constant_kind_case());
  EXPECT_EQ(val, p.string_value());
}

TEST_F(ConstantConversionTest, ConstantToProto_Duration) {
  absl::Duration val = absl::Seconds(100) + absl::Nanos(500);
  Constant c(val);
  ConstantProto p;
  auto status = ConstantToProto(c, &p);
  ASSERT_TRUE(status.ok());
  EXPECT_EQ(ConstantProto::kDurationValue, p.constant_kind_case());
  EXPECT_EQ(100, p.duration_value().seconds());
  EXPECT_EQ(500, p.duration_value().nanos());
}

TEST_F(ConstantConversionTest, ConstantToProto_Time) {
  absl::Time val = absl::UnixEpoch() + absl::Seconds(200) + absl::Nanos(300);
  Constant c(val);
  ConstantProto p;
  auto status = ConstantToProto(c, &p);
  ASSERT_TRUE(status.ok());
  EXPECT_EQ(ConstantProto::kTimestampValue, p.constant_kind_case());
  EXPECT_EQ(200, p.timestamp_value().seconds());
  EXPECT_EQ(300, p.timestamp_value().nanos());
}

TEST_F(ConstantConversionTest, ConstantFromProto_NotSet) {
  ConstantProto p;
  Constant c;
  // Initialize c to something else to ensure reset
  c = Constant(true); 
  auto status = ConstantFromProto(p, c);
  ASSERT_TRUE(status.ok());
  EXPECT_TRUE(std::holds_alternative<absl::monostate>(c.kind()));
}

TEST_F(ConstantConversionTest, ConstantFromProto_Null) {
  ConstantProto p;
  p.set_null_value(google::protobuf::NULL_VALUE);
  Constant c;
  auto status = ConstantFromProto(p, c);
  ASSERT_TRUE(status.ok());
  EXPECT_TRUE(std::holds_alternative<std::nullptr_t>(c.kind()));
}

TEST_F(ConstantConversionTest, ConstantFromProto_Bool) {
  ConstantProto p;
  p.set_bool_value(true);
  Constant c;
  auto status = ConstantFromProto(p, c);
  ASSERT_TRUE(status.ok());
  ASSERT_TRUE(std::holds_alternative<bool>(c.kind()));
  EXPECT_TRUE(std::get<bool>(c.kind()));

  p.set_bool_value(false);
  status = ConstantFromProto(p, c);
  ASSERT_TRUE(status.ok());
  EXPECT_FALSE(std::get<bool>(c.kind()));
}

TEST_F(ConstantConversionTest, ConstantFromProto_Int64) {
  ConstantProto p;
  int64_t val = -9876543210;
  p.set_int64_value(val);
  Constant c;
  auto status = ConstantFromProto(p, c);
  ASSERT_TRUE(status.ok());
  ASSERT_TRUE(std::holds_alternative<int64_t>(c.kind()));
  EXPECT_EQ(val, std::get<int64_t>(c.kind()));
}

TEST_F(ConstantConversionTest, ConstantFromProto_Uint64) {
  ConstantProto p;
  uint64_t val = 9876543210;
  p.set_uint64_value(val);
  Constant c;
  auto status = ConstantFromProto(p, c);
  ASSERT_TRUE(status.ok());
  ASSERT_TRUE(std::holds_alternative<uint64_t>(c.kind()));
  EXPECT_EQ(val, std::get<uint64_t>(c.kind()));
}

TEST_F(ConstantConversionTest, ConstantFromProto_Double) {
  ConstantProto p;
  double val = 2.71828;
  p.set_double_value(val);
  Constant c;
  auto status = ConstantFromProto(p, c);
  ASSERT_TRUE(status.ok());
  ASSERT_TRUE(std::holds_alternative<double>(c.kind()));
  EXPECT_DOUBLE_EQ(val, std::get<double>(c.kind()));
}

TEST_F(ConstantConversionTest, ConstantFromProto_String) {
  ConstantProto p;
  std::string val = "test_string";
  p.set_string_value(val);
  Constant c;
  auto status = ConstantFromProto(p, c);
  ASSERT_TRUE(status.ok());
  ASSERT_TRUE(std::holds_alternative<StringConstant>(c.kind()));
  EXPECT_EQ(val, std::get<StringConstant>(c.kind()));
}

TEST_F(ConstantConversionTest, ConstantFromProto_Bytes) {
  ConstantProto p;
  std::string val = "test_bytes\xFF";
  p.set_bytes_value(val);
  Constant c;
  auto status = ConstantFromProto(p, c);
  ASSERT_TRUE(status.ok());
  ASSERT_TRUE(std::holds_alternative<BytesConstant>(c.kind()));
  EXPECT_EQ(val, std::get<BytesConstant>(c.kind()));
}

TEST_F(ConstantConversionTest, ConstantFromProto_Duration) {
  ConstantProto p;
  auto* dur = p.mutable_duration_value();
  dur->set_seconds(500);
  dur->set_nanos(100);
  Constant c;
  auto status = ConstantFromProto(p, c);
  ASSERT_TRUE(status.ok());
  ASSERT_TRUE(std::holds_alternative<absl::Duration>(c.kind()));
  auto result = std::get<absl::Duration>(c.kind());
  EXPECT_EQ(500, absl::ToInt64Seconds(result));
  EXPECT_EQ(100, absl::ToInt64Nanoseconds(result - absl::Seconds(500)));
}

TEST_F(ConstantConversionTest, ConstantFromProto_Time) {
  ConstantProto p;
  auto* ts = p.mutable_timestamp_value();
  ts->set_seconds(1000);
  ts->set_nanos(200);
  Constant c;
  auto status = ConstantFromProto(p, c);
  ASSERT_TRUE(status.ok());
  ASSERT_TRUE(std::holds_alternative<absl::Time>(c.kind()));
  auto result = std::get<absl::Time>(c.kind());
  auto delta = result - absl::UnixEpoch();
  EXPECT_EQ(1000, absl::ToInt64Seconds(delta));
  EXPECT_EQ(200, absl::ToInt64Nanoseconds(delta - absl::Seconds(1000)));
}

} // namespace