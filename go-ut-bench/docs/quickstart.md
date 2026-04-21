# 快速上手指南

## 一、环境准备

### 1.1 安装依赖工具

**Python:**
```bash
pip install pytest coverage mutmut
```

**Go:**
```bash
go install github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest
```

**Java:**
- Maven 3+
- pitest 1.15+ (通过 Maven 插件)

**C++:**
```bash
sudo apt-get install cmake clang-15 libgtest-dev g++-15
# Mull 突变测试
sudo apt-get install mull-14  # Ubuntu 22.04
```

### 1.2 构建工具

```bash
cd go-ut-bench
go build -o utbench ./cmd/utbench/
```

### 1.3 设置 API Key

```bash
export DEEPSEEK_API_KEY="sk-xxx"
export MINIMAX_API_KEY="xxx"
export VOLCENGINE_API_KEY="xxx"  # doubao
```

---

## 二、跑一次完整流程（推荐新手先 Dry Run）

### 2.1 Dry Run 测试（不调用真实 API）

```bash
./utbench run \
  --models deepseek \
  --langs python \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --max-samples 2 \
  --dry-run
```

确认输出正常后，再去掉 `--dry-run`。

### 2.2 完整命令格式

```bash
./utbench run \
  --models deepseek,minimax,doubao \
  --langs python,java,go,cpp \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --class self_contained \
  --scenario boundary,simple_function,complex_dependency,interface_mock \
  --max-samples 5 \
  --mutation-enabled \
  --mutation-timeout 1800
```

**参数说明：**

| 参数 | 必须 | 说明 |
|------|------|------|
| `--models` | 是 | 模型名，多个逗号分隔。可选：`deepseek`, `minimax`, `doubao` |
| `--langs` | 是 | 语言，多个逗号分隔。可选：`python`, `java`, `go`, `cpp` |
| `--dataset-root` | 是 | 数据集路径 |
| `--output-root` | 是 | 结果输出路径 |
| `--class` | 否 | 数据集类型：`self_contained` 或 `module_level`，默认全部 |
| `--scenario` | 否 | 测试场景：`boundary`, `simple_function`, `complex_dependency`, `interface_mock`，默认全部 |
| `--max-samples` | 否 | 每个场景采样数量，默认全部（0） |
| `--mutation-enabled` | 否 | 是否开启变异测试，默认关闭 |
| `--mutation-timeout` | 否 | 变异测试超时（秒），默认 1800 |
| `--dry-run` | 否 | 不调用 API，用于测试 |

---

## 三、分步执行（方便调试）

### Step 1: 生成测试

```bash
./utbench generate \
  --models deepseek,minimax \
  --langs python,java \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --max-samples 5
```

生成完成后会输出：
```
Generated 20 cases
Manifest: ./artifacts/runs/run_1234567890/generated/generated_manifest.json
```

### Step 2: 评测

```bash
./utbench evaluate \
  --manifest ./artifacts/runs/run_1234567890/generated/generated_manifest.json \
  --output-root ./artifacts \
  --mutation-enabled
```

### Step 3: 生成报告

```bash
./utbench report \
  --input ./artifacts/runs/run_1234567890/evaluation/evaluation_result.json \
  --output-root ./artifacts
```

---

## 四、只跑特定语言/模型

### 只跑 Python
```bash
./utbench run \
  --models deepseek \
  --langs python \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --max-samples 5
```

### 只跑 boundary 场景
```bash
./utbench run \
  --models deepseek \
  --langs python,java,go,cpp \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --scenario boundary \
  --max-samples 5
```

### 只跑 self_contained 类
```bash
./utbench run \
  --models deepseek \
  --langs python,java,go,cpp \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --class self_contained \
  --max-samples 5
```

---

## 五、查看结果

```bash
# 评测结果
cat ./artifacts/runs/run_1234567890/evaluation/evaluation_result.json

# 生成的测试文件
ls ./artifacts/runs/run_1234567890/generated/tests/

# HTML 报告
cat ./artifacts/runs/run_1234567890/reports/reporter_*.html
```

---

## 六、常见问题

### Q: 报 "no such file or directory" 错误
检查 `--dataset-root` 路径是否正确，默认是 `./datasets`

### Q: 模型 API 调用失败
确认 API Key 已设置：`echo $DEEPSEEK_API_KEY`

### Q: 变异测试太慢
- 减少 `--max-samples`
- 减小 `--mutation-timeout`
- 或不加 `--mutation-enabled`

### Q: 想跳过某些场景
使用 `--scenario` 指定，如 `--scenario boundary,simple_function`

---

## 七、完整示例：3模型 × 4语言 × 4场景 × 5样本

```bash
export DEEPSEEK_API_KEY="sk-xxx"
export MINIMAX_API_KEY="xxx"
export VOLCENGINE_API_KEY="xxx"

./utbench run \
  --models deepseek,minimax,doubao \
  --langs python,java,go,cpp \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --class self_contained \
  --scenario boundary,simple_function,complex_dependency,interface_mock \
  --max-samples 5 \
  --mutation-enabled \
  --mutation-timeout 1800
```

**预期产出：**
- 3 模型 × 4 语言 × 4 场景 × 5 样本 = 240 个评测结果
- 每个结果包含：编译通过率、测试通过率、覆盖率、变异分数
