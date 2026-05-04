# C/C++ 单元测试规范

## 单元测试编写规范

### 基本原则

- 对于每个被测函数，所有补充的测试逻辑都编写在同一个测试文件中
- 仅编写单元测试代码，不修改被测试代码
- 尽量做最小变更
- 尽量接近已有代码编写风格、遵循项目构建工具、测试框架规范
- 限制测试代码的副作用：
  - 测试逻辑应该尽量不依赖外部环境，也尽量少对外部环境产生影响
  - 如果测试中涉及读写文件，应使用临时目录（如 `/tmp`）并在测试结束后清理
  - 被测逻辑依赖外部 HTTP 服务或其他 I/O 的，应通过 Mock 或 Fake 对象替代
  - 对依赖的外部函数、全局变量等使用适当的隔离技术

### 测试框架选择

优先使用项目中已经使用的测试框架。如果项目中未明确使用测试框架，推荐使用：

- **推荐 GoogleTest (又称 GTest )** : 最流行的 C++ 测试框架，也可用于测试 C 代码。具体规范参考 `./gtest.md`
- Unity: 是轻量级的 C 测试框架，专用于 C 。具体规范参考 `./unity.md`
- Catch2: 现代化的 C++ 测试框架，语法更简洁。具体规范参考 `./catch2.md`

### C 代码测试的特殊考虑

对于纯 C 代码的测试：

1. **使用 `extern "C"` 包装**（如果使用 C++ 测试框架）：

   ```cpp
   extern "C" {
   #include "module_under_test.h"
   }
   ```

2. **函数指针替换** ：可以通过函数指针实现依赖注入

   ```c
   // 被测代码
   int (*external_func)(int) = real_external_func;

   // 测试代码中替换
   external_func = mock_external_func;
   ```

3. **链接时替换** ：使用 `--wrap` 链接选项 Mock 函数
   ```bash
   -Wl,--wrap=function_to_mock
   ```

### 内存管理测试

C/C++ 测试特别需要关注内存管理：

1. **内存泄漏检测** ：使用 Valgrind 或 AddressSanitizer
2. **测试中的内存分配** ：确保测试结束时释放所有分配的内存
3. **测试清理** ：在 TearDown 或测试结束时清理资源

```cpp
TEST_F(MemoryTest, AllocateAndFree) {
    // Arrange
    void* ptr = malloc(100);
    ASSERT_NE(nullptr, ptr);

    // Act
    process_data(ptr);

    // Cleanup
    free(ptr);  // 确保释放
}
```

### 头文件依赖处理

- 测试文件需要包含被测模块的头文件
- **static 函数测试原则**：
  - **禁止直接测试**：`static` 函数是内部实现细节。禁止直接针对 `static` 函数编写单元测试。其正确性必须通过测试调用它的公开接口（Public API）来间接验证。
  - **寻找调用路径**：如果目标函数是 `static` 的，应首先查找该函数在当前文件内的调用者（Callers）。测试目标应转为这些公开调用者，通过构造输入触发 `static` 函数的逻辑。
  - **处理不可触达代码**：若发现 `static` 函数无法通过任何公开接口触达，应将其识别为“疑似死代码”并报告给用户，不应为其生成测试。
  - **不修改可见性**：禁止为了测试而将 `static` 改为非 `static`，也禁止使用宏（如 `#define static`）在测试中强行暴露 `static` 函数。

## 测试代码和函数命名规范

- 优先遵循测试框架规范
- 测试函数名应清晰表达测试意图
- 遵循项目已有的命名风格
- C 代码的测试文件可以是 `.c` 或 `.cpp` (使用 C++ 测试框架时)
- 测试代码一般与被测代码在同一个目录内添加 `_test` 后缀

## 运行编译检查方法

### CMake 项目

在项目构建目录运行编译检查：

```shell
# 配置（如果未配置）
cmake -S . -B build

# 编译测试
cmake --build build --target {test_target}
```

其中 `{test_target}` 是测试目标名，通常类似 `{module}_test`

### Makefile 项目

```shell
# 编译测试
make tests

# 或编译特定测试目标
make {test_target}
```

### Bazel 项目

```shell
# 编译测试
bazel build //path/to:test_target
```

### Meson 项目

```shell
# 配置（如果未配置）
meson setup builddir -Dtests=true

# 编译测试
meson compile -C builddir
```

具体规范参考 `./meson.md`

### CMake + CTest

#### 基本运行

```shell
# 运行所有测试
ctest --test-dir build

# 运行特定测试
ctest --test-dir build -R {test_name_regex}

# 详细输出
ctest --test-dir build --verbose
```

#### 生成覆盖率数据

**方法一：使用 gcov + lcov**

```shell
# 1. 配置时启用覆盖率
cmake -DCMAKE_CXX_FLAGS="--coverage" -DCMAKE_C_FLAGS="--coverage" -B build

# 2. 编译测试
cmake --build build

# 3. 运行测试
ctest --test-dir build

# 4. 生成 lcov 格式覆盖率数据
lcov --capture --directory build --output-file /tmp/coverage.info
lcov --remove /tmp/coverage.info '/usr/*' '*/test/*' --output-file /tmp/coverage.info
```

**方法二：如果项目有 coverage target**

```shell
# 1. 运行 coverage target（会自动运行测试并生成覆盖率）
cmake --build build --target coverage
```

`lcov` 格式覆盖率数据输出路径：`/tmp/coverage.info`

### 直接运行测试可执行文件

#### 基本运行

```shell
# Google Test
./build/tests/{test_executable}

# 运行特定测试
./build/tests/{test_executable} --gtest_filter=TestSuite.TestCase

# Catch2
./build/tests/{test_executable}

# 运行特定测试
./build/tests/{test_executable} "test case name"
```

#### 生成覆盖率数据

```shell
# 1. 编译时启用覆盖率（如果尚未启用）
cmake -DCMAKE_CXX_FLAGS="--coverage" -DCMAKE_C_FLAGS="--coverage" -B build
cmake --build build

# 2. 运行测试可执行文件
./build/tests/{test_executable}

# 3. 生成 lcov 格式覆盖率数据
lcov --capture --directory build --output-file /tmp/coverage.info
lcov --remove /tmp/coverage.info '/usr/*' '*/test/*' --output-file /tmp/coverage.info
```

`lcov` 格式覆盖率数据输出路径：`/tmp/coverage.info`

### Makefile 项目

#### 基本运行

```shell
# 运行测试
make test

# 或运行特定测试
./tests/{test_executable}
```

#### 生成覆盖率数据

```shell
# 1. 清理之前的覆盖率数据
make clean

# 2. 使用覆盖率标志编译
CFLAGS="--coverage" CXXFLAGS="--coverage" LDFLAGS="--coverage" make

# 3. 运行测试
make test
# 或直接运行测试可执行文件
./tests/{test_executable}

# 4. 生成 lcov 格式覆盖率数据
lcov --capture --directory . --output-file /tmp/coverage.info
lcov --remove /tmp/coverage.info '/usr/*' '*/test/*' --output-file /tmp/coverage.info
```

`lcov` 格式覆盖率数据输出路径：`/tmp/coverage.info`

### Bazel 项目

#### 基本运行

```shell
# 运行测试
bazel test //path/to:test_target

# 运行并显示输出
bazel test //path/to:test_target --test_output=all
```

#### 生成覆盖率数据

```shell
# 1. 运行测试并生成覆盖率
bazel coverage //path/to:test_target --combined_report=lcov

# 2. 覆盖率数据位于
# bazel-out/_coverage/_coverage_report.dat (这是 lcov 格式)

# 3. 复制到标准路径
cp bazel-out/_coverage/_coverage_report.dat /tmp/coverage.info
```

`lcov` 格式覆盖率数据输出路径：`/tmp/coverage.info`

### Meson 项目

```shell
# 运行测试
meson test -C builddir
```
