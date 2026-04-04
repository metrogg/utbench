#include <gtest/gtest.h>
#include <string>
#include "codelab/exercise1.h"
#include "absl/status/status.h"

namespace google::api::expr::codelab {
namespace {

// Normal path: valid string expressions
TEST(ParseAndEvaluateTest, DoubleQuotedStringLiteral) {
  auto result = ParseAndEvaluate("\"hello\"");
  ASSERT_TRUE(result.ok()) << result.status().message();
  EXPECT_EQ(result.value(), "hello");
}

TEST(ParseAndEvaluateTest, SingleQuotedStringLiteral) {
  auto result = ParseAndEvaluate("'world'");
  ASSERT_TRUE(result.ok()) << result.status().message();
  EXPECT_EQ(result.value(), "world");
}

TEST(ParseAndEvaluateTest, StringConcatenation) {
  auto result = ParseAndEvaluate("\"foo\" + \"bar\"");
  ASSERT_TRUE(result.ok()) << result.status().message();
  EXPECT_EQ(result.value(), "foobar");
}

TEST(ParseAndEvaluateTest, EmptyStringLiteral) {
  auto result = ParseAndEvaluate("\"\"");
  ASSERT_TRUE(result.ok()) << result.status().message();
  EXPECT_EQ(result.value(), "");
}

// Boundary conditions: type mismatches (non-string results)
TEST(ParseAndEvaluateTest, IntegerResultReturnsInvalidArgument) {
  auto result = ParseAndEvaluate("1 + 2");
  ASSERT_FALSE(result.ok());
  EXPECT_EQ(result.status().code(), absl::StatusCode::kInvalidArgument);
}

TEST(ParseAndEvaluateTest, BooleanResultReturnsInvalidArgument) {
  auto result = ParseAndEvaluate("true");
  ASSERT_FALSE(result.ok());
  EXPECT_EQ(result.status().code(), absl::StatusCode::kInvalidArgument);
}

TEST(ParseAndEvaluateTest, DoubleResultReturnsInvalidArgument) {
  auto result = ParseAndEvaluate("3.14");
  ASSERT_FALSE(result.ok());
  EXPECT_EQ(result.status().code(), absl::StatusCode::kInvalidArgument);
}

TEST(ParseAndEvaluateTest, NullResultReturnsInvalidArgument) {
  auto result = ParseAndEvaluate("null");
  ASSERT_FALSE(result.ok());
  EXPECT_EQ(result.status().code(), absl::StatusCode::kInvalidArgument);
}

// Error paths: syntax and parsing failures
TEST(ParseAndEvaluateTest, SyntaxErrorReturnsFailure) {
  auto result = ParseAndEvaluate("1 + + 2");
  ASSERT_FALSE(result.ok());
}

TEST(ParseAndEvaluateTest, UnclosedStringLiteralReturnsFailure) {
  auto result = ParseAndEvaluate("\"unclosed string");
  ASSERT_FALSE(result.ok());
}

TEST(ParseAndEvaluateTest, EmptyInputReturnsFailure) {
  auto result = ParseAndEvaluate("");
  ASSERT_FALSE(result.ok());
}

TEST(ParseAndEvaluateTest, InvalidTokenReturnsFailure) {
  auto result = ParseAndEvaluate("@invalid_token");
  ASSERT_FALSE(result.ok());
}

TEST(ParseAndEvaluateTest, UnbalancedParenthesesReturnsFailure) {
  auto result = ParseAndEvaluate("(1 + 2");
  ASSERT_FALSE(result.ok());
}

}  // namespace
}  // namespace google::api::expr::codelab