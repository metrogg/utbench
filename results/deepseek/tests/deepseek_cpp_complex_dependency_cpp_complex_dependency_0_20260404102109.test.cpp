#include "codelab/exercise1.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "eval/public/cel_value.h"
#include "eval/public/activation.h"
#include "eval/public/cel_expression.h"
#include "eval/public/cel_expr_builder.h"
#include "eval/public/builtin_func_registrar.h"
#include "parser/parser.h"
#include "google/protobuf/arena.h"
#include <memory>
#include <string>

namespace google::api::expr::codelab {
namespace {

using ::testing::HasSubstr;
using ::testing::Return;
using ::testing::_;

// Mock classes for testing error paths
class MockCelExpressionBuilder : public ::google::api::expr::runtime::CelExpressionBuilder {
 public:
  MOCK_METHOD(absl::StatusOr<std::unique_ptr<::google::api::expr::runtime::CelExpression>>,
              CreateExpression,
              (const ::google::api::expr::v1alpha1::Expr*,
               const ::google::api::expr::v1alpha1::SourceInfo*),
              (override));
  MOCK_METHOD(::google::api::expr::runtime::CelExpressionBuilder&, set_container,
              (absl::string_view), (override));
  MOCK_METHOD(absl::string_view, container, (), (const override));
  MOCK_METHOD(::google::api::expr::runtime::FunctionRegistry*, GetRegistry, (), (override));
  MOCK_METHOD(void, set_short_circuiting, (bool), (override));
  MOCK_METHOD(bool, short_circuiting, (), (const override));
  MOCK_METHOD(void, set_enable_qualified_type_identifiers, (bool), (override));
  MOCK_METHOD(bool, enable_qualified_type_identifiers, (), (const override));
  MOCK_METHOD(void, set_enable_comprehension, (bool), (override));
  MOCK_METHOD(bool, enable_comprehension, (), (const override));
  MOCK_METHOD(void, set_enable_comprehension_list_append, (bool), (override));
  MOCK_METHOD(bool, enable_comprehension_list_append, (), (const override));
  MOCK_METHOD(void, set_enable_regex, (bool), (override));
  MOCK_METHOD(bool, enable_regex, (), (const override));
  MOCK_METHOD(void, set_enable_string_conversion, (bool), (override));
  MOCK_METHOD(bool, enable_string_conversion, (), (const override));
  MOCK_METHOD(void, set_enable_string_concat, (bool), (override));
  MOCK_METHOD(bool, enable_string_concat, (), (const override));
  MOCK_METHOD(void, set_enable_list_concat, (bool), (override));
  MOCK_METHOD(bool, enable_list_concat, (), (const override));
  MOCK_METHOD(void, set_enable_timestamp_duration_arithmetic, (bool), (override));
  MOCK_METHOD(bool, enable_timestamp_duration_arithmetic, (), (const override));
  MOCK_METHOD(void, set_enable_heterogeneous_equality, (bool), (override));
  MOCK_METHOD(bool, enable_heterogeneous_equality, (), (const override));
  MOCK_METHOD(void, set_enable_null_to_proto_cast, (bool), (override));
  MOCK_METHOD(bool, enable_null_to_proto_cast, (), (const override));
  MOCK_METHOD(void, set_enable_wrapper_type_null_unboxing, (bool), (override));
  MOCK_METHOD(bool, enable_wrapper_type_null_unboxing, (), (const override));
  MOCK_METHOD(void, set_enable_qualified_identifier_unparsing, (bool), (override));
  MOCK_METHOD(bool, enable_qualified_identifier_unparsing, (), (const override));
};

class MockCelExpression : public ::google::api::expr::runtime::CelExpression {
 public:
  MOCK_METHOD(absl::StatusOr<::google::api::expr::runtime::CelValue>,
              Evaluate,
              (const ::google::api::expr::runtime::Activation&,
               ::google::protobuf::Arena*),
              (const override));
  MOCK_METHOD(absl::StatusOr<::google::api::expr::runtime::CelValue>,
              Trace,
              (const ::google::api::expr::runtime::Activation&,
               ::google::protobuf::Arena*,
               ::google::api::expr::runtime::CelEvaluationListener),
              (const override));
};

// Test fixture for ParseAndEvaluate
class ParseAndEvaluateTest : public ::testing::Test {
 protected:
  void SetUp() override {}
  void TearDown() override {}
};

// Normal path tests
TEST_F(ParseAndEvaluateTest, ReturnsStringForStringLiteral) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("\"hello world\"");
  ASSERT_TRUE(result.ok());
  EXPECT_EQ(*result, "hello world");
}

TEST_F(ParseAndEvaluateTest, ReturnsStringForStringConcat) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("\"hello\" + \" \" + \"world\"");
  ASSERT_TRUE(result.ok());
  EXPECT_EQ(*result, "hello world");
}

TEST_F(ParseAndEvaluateTest, ReturnsStringForEmptyString) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("\"\"");
  ASSERT_TRUE(result.ok());
  EXPECT_EQ(*result, "");
}

TEST_F(ParseAndEvaluateTest, ReturnsStringForStringWithSpecialChars) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("\"line1\\nline2\\ttab\"");
  ASSERT_TRUE(result.ok());
  EXPECT_EQ(*result, "line1\nline2\ttab");
}

// Boundary condition tests
TEST_F(ParseAndEvaluateTest, HandlesVeryLongString) {
  std::string long_string(10000, 'a');
  absl::StatusOr<std::string> result = ParseAndEvaluate("\"" + long_string + "\"");
  ASSERT_TRUE(result.ok());
  EXPECT_EQ(*result, long_string);
}

TEST_F(ParseAndEvaluateTest, HandlesMinimalValidExpression) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("\"a\"");
  ASSERT_TRUE(result.ok());
  EXPECT_EQ(*result, "a");
}

// Error path tests
TEST_F(ParseAndEvaluateTest, ReturnsErrorForNonStringResult) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("42");
  ASSERT_FALSE(result.ok());
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
  EXPECT_THAT(result.status().message(), HasSubstr("int"));
}

TEST_F(ParseAndEvaluateTest, ReturnsErrorForBooleanResult) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("true");
  ASSERT_FALSE(result.ok());
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
  EXPECT_THAT(result.status().message(), HasSubstr("bool"));
}

TEST_F(ParseAndEvaluateTest, ReturnsErrorForNullResult) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("null");
  ASSERT_FALSE(result.ok());
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
  EXPECT_THAT(result.status().message(), HasSubstr("null_type"));
}

TEST_F(ParseAndEvaluateTest, ReturnsErrorForListResult) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("[1, 2, 3]");
  ASSERT_FALSE(result.ok());
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
  EXPECT_THAT(result.status().message(), HasSubstr("list"));
}

TEST_F(ParseAndEvaluateTest, ReturnsErrorForMapResult) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("{'key': 'value'}");
  ASSERT_FALSE(result.ok());
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
  EXPECT_THAT(result.status().message(), HasSubstr("map"));
}

TEST_F(ParseAndEvaluateTest, ReturnsErrorForInvalidCELSyntax) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("invalid syntax!!");
  ASSERT_FALSE(result.ok());
  EXPECT_THAT(result.status().message(), HasSubstr("Syntax error"));
}

TEST_F(ParseAndEvaluateTest, ReturnsErrorForEmptyExpression) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("");
  ASSERT_FALSE(result.ok());
}

TEST_F(ParseAndEvaluateTest, ReturnsErrorForWhitespaceOnlyExpression) {
  absl::StatusOr<std::string> result = ParseAndEvaluate("   \t\n  ");
  ASSERT_FALSE(result.ok());
}

// Test ConvertResult helper function directly
TEST(ConvertResultTest, ConvertsStringHolderToString) {
  CelValue string_value = CelValue::CreateStringView("test string");
  absl::StatusOr<std::string> result = ConvertResult(string_value);
  ASSERT_TRUE(result.ok());
  EXPECT_EQ(*result, "test string");
}

TEST(ConvertResultTest, ReturnsErrorForIntType) {
  CelValue int_value = CelValue::CreateInt64(42);
  absl::StatusOr<std::string> result = ConvertResult(int_value);
  ASSERT_FALSE(result.ok());
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
  EXPECT_THAT(result.status().message(), HasSubstr("int"));
}

TEST(ConvertResultTest, ReturnsErrorForUintType) {
  CelValue uint_value = CelValue::CreateUint64(42);
  absl::StatusOr<std::string> result = ConvertResult(uint_value);
  ASSERT_FALSE(result.ok());
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
  EXPECT_THAT(result.status().message(), HasSubstr("uint"));
}

TEST(ConvertResultTest, ReturnsErrorForDoubleType) {
  CelValue double_value = CelValue::CreateDouble(3.14);
  absl::StatusOr<std::string> result = ConvertResult(double_value);
  ASSERT_FALSE(result.ok());
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
  EXPECT_THAT(result.status().message(), HasSubstr("double"));
}

TEST(ConvertResultTest, ReturnsErrorForBoolType) {
  CelValue bool_value = CelValue::CreateBool(true);
  absl::StatusOr<std::string> result = ConvertResult(bool_value);
  ASSERT_FALSE(result.ok());
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
  EXPECT_THAT(result.status().message(), HasSubstr("bool"));
}

TEST(ConvertResultTest, ReturnsErrorForBytesType) {
  std::string bytes_data = "bytes";
  CelValue bytes_value = CelValue::CreateBytesView(bytes_data);
  absl::StatusOr<std::string> result = ConvertResult(bytes_value);
  ASSERT_FALSE(result.ok());
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
  EXPECT_THAT(result.status().message(), HasSubstr("bytes"));
}

TEST(ConvertResultTest, ReturnsErrorForDurationType) {
  CelValue duration_value = CelValue::CreateDuration(absl::Seconds(10));
  absl::StatusOr<std::string> result = ConvertResult(duration_value);
  ASSERT_FALSE(result.ok());
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
  EXPECT_THAT(result.status().message(), HasSubstr("duration"));
}

TEST(ConvertResultTest, ReturnsErrorForTimestampType) {
  CelValue timestamp_value = CelValue::CreateTimestamp(absl::Now());
  absl::StatusOr<std::string> result = ConvertResult(timestamp_value);
  ASSERT_FALSE(result.ok());
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
  EXPECT_THAT(result.status().message(), HasSubstr("timestamp"));
}

TEST(ConvertResultTest, ReturnsErrorForNullType) {
  CelValue null_value = CelValue::CreateNull();
  absl::StatusOr<std::string> result = ConvertResult(null_value);
  ASSERT_FALSE(result.ok());
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result"));
  EXPECT_THAT(result.status().message(), HasSubstr("null_type"));
}

}  // namespace
}  // namespace google::api::expr::codelab