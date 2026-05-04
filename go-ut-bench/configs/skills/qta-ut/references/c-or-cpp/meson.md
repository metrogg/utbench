# Meson 构建配置规范

本指南用于指导 AI 在生成 C/C++ 单元测试后，如何更新项目的 Meson 构建配置，确保测试文件被正确包含并可执行。

生成测试文件后，**必须**更新对应目录的 `meson.build`，将 `_test.c` 文件添加到测试源列表。

## 📁 文件结构

```text
src/
├── meson.build          # 顶层配置
├── parser.c
├── parser_test.c          # 测试文件与源文件同目录
└── core/
    ├── meson.build      # 子模块配置
    ├── engine.c
    └── engine_test.c
```

**命名规则**: 源文件 `foo.c` → 测试文件 `foo_test.c`

## 🔧 meson.build 配置

### 子模块配置

```meson
# 源文件
core_sources = files('engine.c', 'parser.c')

# 测试文件 - AI 生成后追加到此列表
core_test_sources = files('engine_test.c', 'parser_test.c')
```

### 顶层配置

```meson
subdir('core')
subdir('utils')

all_sources = src_sources + core_sources + utils_sources
all_test_sources = src_test_sources + core_test_sources + utils_test_sources

if get_option('tests')
  test_exe = executable('test_all',
    all_sources + all_test_sources,
    dependencies: [test_deps],
    include_directories: include_directories('.', 'core', 'utils')
  )
  test('all_tests', test_exe)
endif
```

## ✅ 验证命令

```shell
meson setup builddir -Dtests=true  # 初始化
meson compile -C builddir          # 编译
meson test -C builddir             # 测试
```
