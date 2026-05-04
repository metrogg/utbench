# Unity & FFF (Fixture 模式) 单元测试规范

> 本文档为 AI 辅助生成 C 语言单元测试的指令集。优先保证规则明确、命名一致、避免冲突。

## 1. 核心指令 (Mandatory)

- **命名强制性**：`TEST_GROUP` 必须使用路径下划线命名（如 `pal_base_conf`），严禁使用 PascalCase 或短名。
- **Mock 标注**：若使用 Mock，必须在文件头部添加 `// @mock_wrap: <func>`。
- **禁止系统 Mock**：严禁 Mock 系统底层函数（`malloc`, `free`, `printf`, `open` 等）。
- **隔离性**：每个 `TEST` 必须是原子的，且在 `TEST_SETUP` 中重置所有 Mock 状态。

## 2. 命名规范

| 类型           | 格式             | 示例 (源文件: `pal/base/conf.c`)             |
| :------------- | :--------------- | :------------------------------------------- |
| **TEST_GROUP** | `路径_文件名`    | `TEST_GROUP(pal_base_conf);`                 |
| **TEST**       | `函数_场景_预期` | `TEST(pal_base_conf, read_valid_returns_ok)` |

## 3. 标准测试模板 (Fixture)

```c
#include "unity_fixture.h"
#include "module_under_test.h"

TEST_GROUP(路径_文件名);

TEST_SETUP(路径_文件名) {
    // 调用 setUp 或重置状态
}

TEST_TEAR_DOWN(路径_文件名) {
    // 清理资源
}

TEST(路径_文件名, 函数_场景_预期) {
    // Arrange -> Act -> Assert
    TEST_ASSERT_EQUAL_INT(0, 0);
}

TEST_GROUP_RUNNER(路径_文件名) {
    RUN_TEST_CASE(路径_文件名, 函数_场景_预期);
}
```

## 4. Mock (FFF Wrap 方案)

### 规则

1. **优先复用**：检查项目公共目录是否已有 Mock 实现。
2. **避免冲突**：使用 `__wrap_` 前缀定义 Mock。
3. **标注头部**：文件最上方声明注解。

### 示例

```c
// @mock_wrap: socket_send
// @mock_wrap: socket_recv

#include "unity.h"
#include "fff.h"

DEFINE_FFF_GLOBALS;

/* 使用 __wrap_ 避免链接冲突 */
FAKE_VALUE_FUNC(int, __wrap_socket_send, int, const void*, size_t);
FAKE_VALUE_FUNC(int, __wrap_socket_recv, int, void*, size_t);

void setUp(void) {
    RESET_FAKE(__wrap_socket_send);
    RESET_FAKE(__wrap_socket_recv);
    FFF_RESET_HISTORY();
}
```

## 5. 常用断言 (Cheat Sheet)

| 类别       | 宏                                                                           |
| :--------- | :--------------------------------------------------------------------------- |
| **整数**   | `TEST_ASSERT_EQUAL_INT(exp, act)`, `TEST_ASSERT_INT_WITHIN(delta, exp, act)` |
| **布尔**   | `TEST_ASSERT_TRUE(cond)`, `TEST_ASSERT_FALSE(cond)`                          |
| **字符串** | `TEST_ASSERT_EQUAL_STRING(exp, act)`                                         |
| **指针**   | `TEST_ASSERT_NULL(ptr)`, `TEST_ASSERT_NOT_NULL(ptr)`                         |
| **内存**   | `TEST_ASSERT_EQUAL_MEMORY(exp, act, len)`                                    |
| **自定义** | 所有断言后缀 `_MESSAGE` 均可添加描述                                         |

## 6. 主运行器 (RunAllTests)

```c
#include "unity_fixture.h"

extern void TEST_GROUP_RUNNER(路径_文件名_1);
extern void TEST_GROUP_RUNNER(路径_文件名_2);

static void RunAllTests(void) {
    RUN_TEST_GROUP(路径_文件名_1);
    RUN_TEST_GROUP(路径_文件名_2);
}

int main(int argc, const char* argv[]) {
    return UnityMain(argc, argv, RunAllTests);
}
```
