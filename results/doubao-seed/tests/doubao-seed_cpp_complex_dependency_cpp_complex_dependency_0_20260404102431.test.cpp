#include <string>
#include <gtest/gtest.h>
#include "absl/status/status.h"
#include "absl/status/status_matchers.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "codelab/exercise1.h"

namespace google::api::expr::codelab::test {
using ::absl_testing::IsOkAndHolds;
using ::absl_testing::StatusIs;
using ::testing::HasSubstr;

TEST(ParseAndEvaluateTest, SimpleStringLiteralReturnsCorrectValue) {
  EXPECT_THAT(ParseAndEvaluate(R"("hello world")"), IsOkAndHolds("hello world"));
  EXPECT_THAT(ParseAndEvaluate(R"("test123!@#")"), IsOkAndHolds("test123!@#"));
}

TEST(ParseAndEvaluateTest, EmptyStringLiteralReturnsEmptyString) {
  EXPECT_THAT(ParseAndEvaluate(R"("")"), IsOkAndHolds(""));
}

TEST(ParseAndEvaluateTest, StringConcatenationReturnsCombinedResult) {
  EXPECT_THAT(ParseAndEvaluate(R"("foo" + "bar" + "_suffix")"), IsOkAndHolds("foobar_suffix"));
}

TEST(ParseAndEvaluateTest, BuiltinStringFunctionsWorkCorrectly) {
  EXPECT_THAT(ParseAndEvaluate(R"(lower("HELLO_WORLD"))"), IsOkAndHolds("hello_world"));
  EXPECT_THAT(ParseAndEvaluate(R"(upper("test_input"))"), IsOkAndHolds("TEST_INPUT"));
  EXPECT_THAT(ParseAndEvaluate(R"(strip("  spaced text  "))"), IsOkAndHolds("spaced text"));
}

TEST(ParseAndEvaluateTest, InvalidCelSyntaxReturnsParseError) {
  EXPECT_THAT(ParseAndEvaluate("1 + "), StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(ParseAndEvaluate("'unclosed string literal"), StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(ParseAndEvaluate("invalid!@#syntax"), StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(ParseAndEvaluateTest, NonStringIntResultReturnsTypeError) {
  EXPECT_THAT(ParseAndEvaluate("123 + 456"), 
              StatusIs(absl::StatusCode::kInvalidArgument, 
                       HasSubstr("expected string got 'int64_t'")));
}

TEST(ParseAndEvaluateTest, NonStringBoolResultReturnsTypeError) {
  EXPECT_THAT(ParseAndEvaluate("10 > 5"), 
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("expected string got 'bool'")));
}

TEST(ParseAndEvaluateTest, NonStringNullResultReturnsTypeError) {
  EXPECT_THAT(ParseAndEvaluate("null"),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("expected string got 'null_type'")));
}

TEST(ParseAndEvaluateTest, UndefinedVariableReturnsEvaluationError) {
  EXPECT_THAT(ParseAndEvaluate("undefined_variable"), StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(ParseAndEvaluate("x + y"), StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(ParseAndEvaluateTest, LongStringLiteralReturnsCorrectValue) {
  const std::string long_content(1024, 'x');
  const std::string cel_expr = absl::StrCat("\"", long_content, "\"");
  EXPECT_THAT(ParseAndEvaluate(cel_expr), IsOkAndHolds(long_content));
}

}  // namespace google::api::expr::codelab::test