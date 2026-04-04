#include <gtest/gtest.h>
#include "common/kind.h"
#include "base/kind.h"

namespace cel {

class KindToStringTest : public ::testing::Test {
 protected:
  void SetUp() override {}
};

TEST_F(KindToStringTest, NullType) {
  EXPECT_EQ(KindToString(Kind::kNullType), "null_type");
}

TEST_F(KindToStringTest, Dyn) {
  EXPECT_EQ(KindToString(Kind::kDyn), "dyn");
}

TEST_F(KindToStringTest, Any) {
  EXPECT_EQ(KindToString(Kind::kAny), "any");
}

TEST_F(KindToStringTest, Type) {
  EXPECT_EQ(KindToString(Kind::kType), "type");
}

TEST_F(KindToStringTest, TypeParam) {
  EXPECT_EQ(KindToString(Kind::kTypeParam), "type_param");
}

TEST_F(KindToStringTest, Function) {
  EXPECT_EQ(KindToString(Kind::kFunction), "function");
}