#include "codelab/exercise1.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"
#include "absl/status/statusor.h"
#include "eval/public/cel_value.h"
#include "eval/public/activation.h"
#include "google/protobuf/arena.h"

namespace google::api::expr::codelab {
namespace {

using ::testing::HasSubstr;
using ::testing::Not;

TEST(ParseAndEvaluateTest, ValidStringExpressionReturnsString) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("\"hello world\"");
  ASSERT_TRUE(result.ok()) << "Status: " << result.status();
  EXPECT_EQ(*result, "hello world");
}

TEST(ParseAndEvaluateTest, ValidStringConcatenationReturnsString) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("\"hello\" + \" world\"");
  ASSERT_TRUE(result.ok()) << "Status: " << result.status();
  EXPECT_EQ(*result, "hello world");
}

TEST(ParseAndEvaluateTest, EmptyStringExpressionReturnsEmptyString) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("\"\"");
  ASSERT_TRUE(result.ok()) << "Status: " << result.status();
  EXPECT_EQ(*result, "");
}

TEST(ParseAndEvaluateTest, StringWithSpecialCharactersReturnsCorrectString) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("\"test\\nstring\\twith\\\"quotes\\\"\"");
  ASSERT_TRUE(result.ok()) << "Status: " << result.status();
  EXPECT_EQ(*result, "test\nstring\twith\"quotes\"");
}

TEST(ParseAndEvaluateTest, NonStringResultReturnsInvalidArgumentError) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("42");
  ASSERT_FALSE(result.ok());
  EXPECT_EQ(result.status().code(), absl::StatusCode::kInvalidArgument);
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
}

TEST(ParseAndEvaluateTest, BooleanResultReturnsInvalidArgumentError) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("true");
  ASSERT_FALSE(result.ok());
  EXPECT_EQ(result.status().code(), absl::StatusCode::kInvalidArgument);
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
}

TEST(ParseAndEvaluateTest, ListResultReturnsInvalidArgumentError) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("[1, 2, 3]");
  ASSERT_FALSE(result.ok());
  EXPECT_EQ(result.status().code(), absl::StatusCode::kInvalidArgument);
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
}

TEST(ParseAndEvaluateTest, MapResultReturnsInvalidArgumentError) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("{'key': 'value'}");
  ASSERT_FALSE(result.ok());
  EXPECT_EQ(result.status().code(), absl::StatusCode::kInvalidArgument);
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
}

TEST(ParseAndEvaluateTest, InvalidSyntaxReturnsError) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("\"unclosed string");
  ASSERT_FALSE(result.ok());
  EXPECT_THAT(result.status().message(), Not(HasSubstr("expected string result")));
}

TEST(ParseAndEvaluateTest, MalformedExpressionReturnsError) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("1 + ");
  ASSERT_FALSE(result.ok());
  EXPECT_THAT(result.status().message(), Not(HasSubstr("expected string result")));
}

TEST(ParseAndEvaluateTest, ComplexStringExpressionReturnsCorrectResult) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("\"prefix_\" + \"suffix\"");
  ASSERT_TRUE(result.ok()) << "Status: " << result.status();
  EXPECT_EQ(*result, "prefix_suffix");
}

TEST(ParseAndEvaluateTest, StringWithUnicodeCharacters) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("\"café 🍵\"");
  ASSERT_TRUE(result.ok()) << "Status: " << result.status();
  EXPECT_EQ(*result, "café 🍵");
}

TEST(ParseAndEvaluateTest, NullResultReturnsInvalidArgumentError) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("null");
  ASSERT_FALSE(result.ok());
  EXPECT_EQ(result.status().code(), absl::StatusCode::kInvalidArgument);
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
}

TEST(ParseAndEvaluateTest, VeryLongStringExpression) {
  std::string long_string(1000, 'a');
  absl::StatusOr<std::string> result = ParseAndEvaluate("\"" + long_string + "\"");
  ASSERT_TRUE(result.ok()) << "Status: " << result.status();
  EXPECT_EQ(*result, long_string);
}

TEST(ConvertResultTest, StringValueReturnsString) {
  CelValue value = CelValue::CreateStringView("test string");
  absl::StatusOr<std::string> result = ConvertResult(value);
  ASSERT_TRUE(result.ok());
  EXPECT_EQ(*result, "test string");
}

TEST(ConvertResultTest, IntValueReturnsInvalidArgument) {
  CelValue value = CelValue::CreateInt64(42);
  absl::StatusOr<std::string> result = ConvertResult(value);
  ASSERT_FALSE(result.ok());
  EXPECT_EQ(result.status().code(), absl::StatusCode::kInvalidArgument);
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
}

TEST(ConvertResultTest, DoubleValueReturnsInvalidArgument) {
  CelValue value = CelValue::CreateDouble(3.14);
  absl::StatusOr<std::string> result = ConvertResult(value);
  ASSERT_FALSE(result.ok());
  EXPECT_EQ(result.status().code(), absl::StatusCode::kInvalidArgument);
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
}

TEST(ConvertResultTest, BoolValueReturnsInvalidArgument) {
  CelValue value = CelValue::CreateBool(true);
  absl::StatusOr<std::string> result = ConvertResult(value);
  ASSERT_FALSE(result.ok());
  EXPECT_EQ(result.status().code(), absl::StatusCode::kInvalidArgument);
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
}

TEST(ConvertResultTest, NullValueReturnsInvalidArgument) {
  CelValue value = CelValue::CreateNull();
  absl::StatusOr<std::string> result = ConvertResult(value);
  ASSERT_FALSE(result.ok());
  EXPECT_EQ(result.status().code(), absl::StatusCode::kInvalidArgument);
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
}

TEST(ConvertResultTest, BytesValueReturnsInvalidArgument) {
  std::string bytes_data = "bytes";
  CelValue value = CelValue::CreateBytes(&bytes_data);
  absl::StatusOr<std::string> result = ConvertResult(value);
  ASSERT_FALSE(result.ok());
  EXPECT_EQ(result.status().code(), absl::StatusCode::kInvalidArgument);
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
}

TEST(ConvertResultTest, EmptyStringValueReturnsEmptyString) {
  CelValue value = CelValue::CreateStringView("");
  absl::StatusOr<std::string> result = ConvertResult(value);
  ASSERT_TRUE(result.ok());
  EXPECT_EQ(*result, "");
}

}  // namespace
}  // namespace google::api::expr::codelab