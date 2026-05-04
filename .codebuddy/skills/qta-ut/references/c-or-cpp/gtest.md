# Google Test

Google Test (gtest/gmock) 是最流行的 C++ 测试框架，也可用于测试 C 代码：

- **基础宏** ：
  - `TEST(TestSuiteName, TestName)` - 定义测试用例
  - `TEST_F(TestFixtureName, TestName)` - 定义使用 fixture 的测试用例
  - `ASSERT_*` - 断言失败时终止当前测试
  - `EXPECT_*` - 断言失败时继续执行

- **常用断言** ：
  - `EXPECT_EQ(expected, actual)` - 相等断言
  - `EXPECT_NE(val1, val2)` - 不等断言
  - `EXPECT_TRUE(condition)` / `EXPECT_FALSE(condition)` - 布尔断言
  - `EXPECT_LT/LE/GT/GE(val1, val2)` - 比较断言
  - `EXPECT_STREQ(str1, str2)` - C 字符串相等断言
  - `EXPECT_NEAR(val1, val2, tolerance)` - 浮点数近似相等
  - `EXPECT_THROW(statement, exception_type)` - 异常断言（C++ 专用）

## 测试代码结构

```cpp
#include <gtest/gtest.h>

// 如果是 C 代码
extern "C" {
#include "module_under_test.h"
}

// 简单测试
TEST(ModuleName, FunctionName_Scenario) {
    // Arrange: 准备测试数据
    int input = 42;

    // Act: 执行被测函数
    int result = function_under_test(input);

    // Assert: 验证结果
    EXPECT_EQ(expected_value, result);
}

// 使用 Fixture 的测试
class ModuleTest : public ::testing::Test {
protected:
    void SetUp() override {
        // 初始化
        test_data = create_test_data();
    }

    void TearDown() override {
        // 清理
        cleanup_test_data(test_data);
    }

    void* test_data;
};

TEST_F(ModuleTest, FunctionName_Scenario) {
    // 使用 test_data
    EXPECT_TRUE(validate(test_data));
}
```

## 测试代码命名规范

- **测试文件命名** ： `{module}_test.cc` 或 `{module}_unittest.cc`
  - 示例： `calculator.c` → `calculator_test.cc`
- **测试文件路径** ：
  - 通常与被测文件同目录或放在 `tests/` 目录下
  - CMake 项目常见： `tests/unit/{module}_test.cc`
- **TEST 宏命名** ： `TEST(TestSuiteName, TestCaseName)`
  - `TestSuiteName` ： 通常是模块名或类名，使用 PascalCase
  - `TestCaseName` ： 描述测试场景，格式为 `FunctionName_Scenario`
  - 示例： `TEST(Calculator, Add_PositiveNumbers)`
- **TEST_F 宏命名** ： `TEST_F(FixtureName, TestCaseName)`
  - Fixture 类名通常是 `{Module}Test` 或 `{Module}Fixture`

## 测试 Fixture ：用于共享测试环境

```cpp
class MyTestFixture : public ::testing::Test {
protected:
    void SetUp() override {
        // 每个测试前的初始化
    }

    void TearDown() override {
        // 每个测试后的清理
    }

    // 共享的测试数据
    int shared_data;
};

TEST_F(MyTestFixture, TestCase1) {
    // 使用 shared_data
}
```

## 参数化测试

```cpp
class MyParamTest : public ::testing::TestWithParam<int> {};

TEST_P(MyParamTest, TestWithParam) {
    int param = GetParam();
    // 使用参数进行测试
}

INSTANTIATE_TEST_SUITE_P(MyTests, MyParamTest,
                         ::testing::Values(1, 2, 3, 4));
```

## Mock 框架 (Google Mock)

- 用于 C++ 接口的 Mock
- `MOCK_METHOD` - 定义 Mock 方法
- `EXPECT_CALL` - 设置期望调用
- `ON_CALL` - 设置默认行为

```cpp
class MockDependency : public DependencyInterface {
public:
    MOCK_METHOD(int, getValue, (), (override));
    MOCK_METHOD(void, setValue, (int value), (override));
};

TEST(MyTest, UseMock) {
    MockDependency mock;
    EXPECT_CALL(mock, getValue())
        .WillOnce(::testing::Return(42));

    int result = function_using_dependency(&mock);
    EXPECT_EQ(42, result);
}
```
