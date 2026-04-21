# go-ut-bench

多语言单元测试生成效果横向评测 CLI 工具。

## 支持语言

| 语言 | 测试框架 | 覆盖率工具 | 变异测试 |
|------|---------|-----------|---------|
| Python | pytest | coverage | mutmut |
| Go | go test | go test -cover | avito-tech/go-mutesting |
| Java | Maven | JaCoCo | pitest |
| C++ | GoogleTest | gcov | mull-15 |

## 安装

```bash
go build -o utbench ./cmd/utbench/
```

## 快速开始

```bash
# 查看帮助
./utbench --help

# 完整评测流程
./utbench run \
  --models deepseek \
  --langs python \
  --dataset-root ./datasets \
  --class self_contained \
  --max-samples 10

# 分步执行：生成 -> 评测 -> 报告
./utbench generate --models deepseek --langs python
./utbench evaluate --manifest ./artifacts/runs/.../generated_manifest.json
./utbench report --input ./artifacts/runs/.../evaluation_result.json
```

详细文档：[docs/quickstart.md](docs/quickstart.md)

## 命令

| 命令 | 说明 |
|------|------|
| `utbench run` | 完整流程 (generate → evaluate → report) |
| `utbench generate` | 仅生成单元测试 |
| `utbench evaluate` | 评测已生成的单元测试 |
| `utbench report` | 生成评测报告 |
| `utbench ingest` | 结果导入 SQLite |
| `utbench dataset` | 数据集管理 |

详细命令文档：[docs/cli-spec.md](docs/cli-spec.md)

## 数据集

数据集按语言和复杂度组织：

```
datasets/
  python/
    python_code_files_self_contained/
      boundary/          # 边界条件测试
      simple_function/   # 简单函数测试
      complex_dependency/ # 复杂依赖测试
      interface_mock/     # 接口模拟测试
  java/
  go/
  cpp/
```

## 输出结构

```
artifacts/
  runs/
    <run-id>/
      generated/          # 生成的测试文件
      evaluation/         # 评测结果
      reports/           # 报告文件
```

## 环境要求

### Python
```bash
pip install pytest coverage mutmut
```

### Go
```bash
go install github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest
```

### Java
- Maven 3+
- pitest 1.15+

### C++
```bash
sudo apt-get install cmake clang-15 libgtest-dev g++-15
sudo apt-get install mull-14  # 或 mull-15
```

## 配置

模型配置位于 `../benchmark/config/models.yaml`：

```yaml
models:
  deepseek:
    enabled: true
    provider: deepseek
    config:
      api_key_env: DEEPSEEK_API_KEY
      model: deepseek-coder
```
