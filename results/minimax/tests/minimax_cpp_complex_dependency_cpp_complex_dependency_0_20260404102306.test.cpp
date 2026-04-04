#include "codelab/exercise1.h"
#include <gtest/gtest.h>
#include "absl/status/status.h"
#include "absl/strings/string_view.h"

namespace google::api::expr::codelab {
namespace {

class ParseAndEvaluateTest : public ::testing::Test {
 protected:
  void SetUp() override {
    // Setup if needed
  }
};

// Test valid expressions that return strings
TEST_F(ParseAndEvaluateTest, SimpleStringLiteral) {
  auto result = ParseAndEvaluate("'hello world'");
  ASSERT_TRUE(result.ok()) << result.status();
  EXPECT_EQ(*result, "hello world");
}

TEST_F(ParseAndEvaluateTest, StringConcatenation) {
  auto result = ParseAndEvaluate("'hello' + ' ' + 'world'");
  ASSERT_TRUE(result.ok()) << result.status();
  EXPECT_EQ(*result, "hello world");
}

TEST_F(ParseAndEvaluateTest, ArithmeticExpression) {
  auto result = ParseAndEvaluate("1 + 2");
  ASSERT_TRUE(result.ok()) << result.status();
  EXPECT_EQ(*result, "3");
}

TEST_F(ParseAndEvaluateTest, BooleanExpression) {
  auto result = ParseAndEvaluate("true == true");
  ASSERT_TRUE(result.ok()) << result.status();
  EXPECT_EQ(*result, "true");
}