#include "gtest/gtest.h"
#include "gmock/gmock.h"
#include "absl/status/status_matchers.h"
#include "codelab/exercise1.h"

namespace google::api::expr::codelab {
namespace {

using ::absl_testing::IsOkAndHolds;
using ::absl_testing::StatusIs;
using ::testing::HasSubstr;
using ::testing::Eq;

// Test parsing and evaluating a simple string literal.
TEST(ParseAndEvaluateTest, SimpleStringLiteral) {
  auto result = ParseAndEvaluate("\"hello world\"");
  EXPECT_THAT(result, IsOkAndHolds("hello world"));
}

// Test string concatenation.
TEST(ParseAndEvaluateTest, StringConcatenation) {
  auto result = ParseAndEvaluate("\"foo\" + \" \" + \"bar\"");
  EXPECT_THAT(result, IsOkAndHolds("foo bar"));
}

// Test evaluating an empty string.
TEST(ParseAndEvaluateTest, EmptyString) {
  auto result = ParseAndEvaluate("\"\"");
  EXPECT_THAT(result, IsOkAndHolds(""));
}

// Test a ternary operator that results in a string.
TEST(ParseAndEvaluateTest, TernaryOperatorReturningString) {
  auto result = ParseAndEvaluate("true ? \"yes\" : \"no\"");
  EXPECT_THAT(result, IsOkAndHolds("yes"));
}

// Test a function call that results in a string (e.g., contains check wrapped in ternary).
TEST(ParseAndEvaluateTest, FunctionCallWithTernary) {
  auto result = ParseAndEvaluate("\"apple\".contains(\"p\") ? \"contains p\" : \"no p\"");
  EXPECT_THAT(result, IsOkAndHolds("contains p"));
}

// Test handling of syntax errors in the expression.
TEST(ParseAndEvaluateTest, SyntaxError) {
  auto result = ParseAndEvaluate("1 + + 2");
  EXPECT_THAT(result, StatusIs(absl::StatusCode::kInvalidArgument));
}

// Test handling of invalid expression (incomplete).
TEST(ParseAndEvaluateTest, IncompleteExpression) {
  auto result = ParseAndEvaluate("\"unclosed string");
  EXPECT_THAT(result, StatusIs(absl::StatusCode::kInvalidArgument));
}

// Test type mismatch: expression evaluates to an integer, but string is expected.
TEST(ParseAndEvaluateTest, ResultIsInteger) {
  auto result = ParseAndEvaluate("123");
  EXPECT_THAT(result, StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result got 'int'"));
}

// Test type mismatch: expression evaluates to a boolean, but string is expected.
TEST(ParseAndEvaluateTest, ResultIsBoolean) {
  auto result = ParseAndEvaluate("false");
  EXPECT_THAT(result, StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result got 'bool'"));
}

// Test type mismatch: expression evaluates to a list, but string is expected.
TEST(ParseAndEvaluateTest, ResultIsList) {
  auto result = ParseAndEvaluate("[1, 2, 3]");
  EXPECT_THAT(result, StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(result.status().message(), HasSubstr("expected string result got 'list'"));
}

// Test evaluation failure due to undefined variable (Activation is empty).
TEST(ParseAndEvaluateTest, UndefinedVariable) {
  auto result = ParseAndEvaluate("undefined_variable");
  EXPECT_THAT(result, StatusIs(absl::StatusCode::kNotFound));
}

} // namespace
} // namespace google::api::expr::codelab