#include <gtest/gtest.h>
#include "absl/status/status.h"
#include "absl/status/status_matchers.h"
#include "absl/strings/string_view.h"
#include "codelab/exercise1.h"

namespace google::api::expr::codelab::test {
using ::absl_testing::IsOkAndHolds;
using ::absl_testing::StatusIs;
using ::testing::HasSubstr;

TEST(ParseAndEvaluateTest, ValidStringLiteralReturnsCorrectValue) {
  EXPECT_THAT(ParseAndEvaluate(R"("hello world")"), IsOkAndHolds("hello world"));
}

TEST(ParseAndEvaluateTest, ValidStringConcatenationReturnsCorrectValue) {
  EXPECT_THAT(ParseAndEvaluate(R"("foo" + "bar" + "123")"), IsOkAndHolds("foobar123"));
}

TEST(ParseAndEvaluateTest, BuiltinStringFunctionReturnsCorrectValue) {
  EXPECT_THAT(ParseAndEvaluate(R"("test string".upper())"), IsOkAndHolds("TEST STRING"));
}

TEST(ParseAndEvaluateTest, EmptyStringLiteralReturnsEmptyString) {
  EXPECT_THAT(ParseAndEvaluate(R"("")"), IsOkAndHolds(""));
}

TEST(ParseAndEvaluateTest, InvalidSyntaxReturnsParseError) {
  EXPECT_THAT(ParseAndEvaluate("1 + "), StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(ParseAndEvaluateTest, IntegerResultReturnsTypeError) {
  auto result = ParseAndEvaluate("1 + 2");
  EXPECT_THAT(result, StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(result.status().message(), HasSubstr("expected string got 'int'"));
}

TEST(ParseAndEvaluateTest, BooleanResultReturnsTypeError) {
  auto result = ParseAndEvaluate("10 > 5");
  EXPECT_THAT(result, StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(result.status().message(), HasSubstr("expected string got 'bool'"));
}

TEST(ParseAndEvaluateTest, DoubleResultReturnsTypeError) {
  auto result = ParseAndEvaluate("3.14159");
  EXPECT_THAT(result, StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(result.status().message(), HasSubstr("expected string got 'double'"));
}

TEST(ParseAndEvaluateTest, NullResultReturnsTypeError) {
  auto result = ParseAndEvaluate("null");
  EXPECT_THAT(result, StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(result.status().message(), HasSubstr("expected string got 'null_type'"));
}

TEST(ParseAndEvaluateTest, UndefinedVariableReturnsEvaluationError) {
  auto result = ParseAndEvaluate(R"(undefined_var + "test")");
  EXPECT_FALSE(result.ok());
}

}  // namespace google::api::expr::codelab::test