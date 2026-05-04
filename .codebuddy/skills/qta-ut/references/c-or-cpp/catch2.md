# Catch2

Catch2 是现代化的 C++ 测试框架，语法更简洁：

- **基础宏** ：
  - `TEST_CASE("description")` - 定义测试用例
  - `SECTION("description")` - 测试用例内的子测试
  - `REQUIRE(expression)` - 断言失败时终止
  - `CHECK(expression)` - 断言失败时继续

- **常用断言** ：
  - `REQUIRE(actual == expected)`
  - `REQUIRE_THAT(value, Predicate)`
  - `REQUIRE_THROWS_AS(expression, exception_type)`

## 测试代码结构

```cpp
#include <catch2/catch.hpp>
#include "module_under_test.h"

TEST_CASE("Module: function description", "[module]") {
    SECTION("Normal case") {
        int result = function_under_test(42);
        REQUIRE(result == expected_value);
    }

    SECTION("Edge case") {
        int result = function_under_test(0);
        REQUIRE(result == 0);
    }
}
```

## 测试代码命名规范

- **测试文件命名** ： `{module}.test.cpp` 或 `test_{module}.cpp`
- **TEST_CASE 命名** ： 使用描述性字符串
  - 格式： `"Module: function - scenario"`
  - 示例： `TEST_CASE("Calculator: add - positive numbers", "[calculator]")`
- **标签** ： 使用 `[tag]` 组织测试
