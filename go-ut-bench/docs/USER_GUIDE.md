# UT-Bench 用户指南

多语言单元测试生成效果横向评测工具，支持 Python、Go、Java、C++ 四种语言。

---

## 目录

1. [快速开始](#快速开始)
2. [环境配置](#环境配置)
3. [命令详解](#命令详解)
4. [模型配置](#模型配置)
5. [数据集说明](#数据集说明)
6. [输出结果](#输出结果)
7. [常见问题](#常见问题)

---

## 快速开始

### 方式一：Docker 运行（推荐）

```bash
# 1. 加载镜像（如有 tar 文件）
docker load -i utbench.tar

# 或构建镜像
docker build -t utbench:latest .

# 2. 配置 API 密钥
cp .env.example .env
# 编辑 .env 文件填入密钥

# 3. 运行测试
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  -w /app \
  utbench:latest run \
    --models deepseek \
    --langs python,java \
    --config /app/configs/models.yaml \
    --dataset-root /app/datasets \
    --max-samples 5 \
    --class self_contained
```

### 方式二：本地运行

```bash
# 1. 构建
go build -o utbench ./cmd/utbench/

# 2. 设置环境变量
export DEEPSEEK_API_KEY="sk-xxx"

# 3. 运行
./utbench run --models deepseek --langs python --max-samples 5
```

---

## 环境配置

### API 密钥

在 `.env` 文件中配置：

```bash
DEEPSEEK_API_KEY=sk-xxx        # DeepSeek
DASHSCOPE_API_KEY=xxx          # 通义千问 (Qwen)
MINIMAX_API_KEY=xxx            # Minimax
VOLCENGINE_API_KEY=xxx         # 豆包 (Doubao)
ARK_API_KEY=xxx                # 豆包 ARK 版本
```

### 语言工具链（本地运行需要）

| 语言 | 测试框架 | 覆盖率工具 | 变异测试 |
|------|---------|-----------|---------|
| Python | pytest | coverage | mutmut |
| Go | go test | go test -cover | go-mutesting |
| Java | Maven/JUnit | JaCoCo | pitest |
| C++ | GoogleTest | gcov | mull |

安装命令：

```bash
# Python
pip install pytest coverage mutmut

# Go
go install github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest

# Java - Maven 自动下载依赖

# C++ (Ubuntu)
sudo apt-get install cmake clang-15 libgtest-dev g++-15 mull-15
```

---

## 命令详解

### utbench run（完整流程）

一键执行：生成 → 评测 → 报告

```bash
utbench run \
  --models deepseek,qwen \
  --langs python,java,go,cpp \
  --config ./configs/models.yaml \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --class self_contained \
  --max-samples 10 \
  --mutation-enabled \
  --mutation-timeout 360
```

**参数说明：**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--config` | `../benchmark/config/models.yaml` | 模型配置文件 |
| `--models` | `deepseek` | 模型列表（逗号分隔） |
| `--langs` | `python` | 语言列表（逗号分隔） |
| `--dataset-root` | `./datasets` | 数据集目录 |
| `--dataset-manifest` | `./configs/dataset_index.json` | 数据集索引 |
| `--output-root` | `./artifacts` | 输出目录 |
| `--class` | `self_contained` | 数据集类别 |
| `--scenario` | 全部 | 场景过滤 |
| `--level` | `l1` | 数据集级别 |
| `--max-samples` | `0`（全部） | 样本数量上限 |
| `--mode` | `full` | 运行模式：`full`/`incremental` |
| `--mutation-enabled` | `true` | 启用变异测试 |
| `--mutation-timeout` | `360` | 变异超时（秒） |
| `--mutation-policy` | `warn` | 变异错误策略 |
| `--total-timeout` | `0`（无限制） | 总超时（分钟） |
| `--dry-run` | `false` | 试运行（不调用API） |
| `--reset-checkpoint` | `false` | 重置断点 |
| `--ingest` | `false` | 导入SQLite |
| `--verbose` | `true` | 详细日志 |

### utbench generate（仅生成）

```bash
utbench generate \
  --models deepseek \
  --langs python \
  --max-samples 5 \
  --dry-run
```

### utbench evaluate（仅评测）

```bash
utbench evaluate \
  --manifest ./artifacts/runs/<run-id>/generated/generated_manifest.json \
  --mutation-enabled
```

### utbench report（仅报告）

```bash
utbench report \
  --input ./artifacts/runs/<run-id>/evaluation/evaluation_result.json
```

### utbench dataset（数据集管理）

```bash
# 生成数据集索引
utbench dataset index --root ./datasets --output ./configs/dataset_index.json

# 查看统计
utbench dataset stats --manifest ./configs/dataset_index.json
```

---

## 模型配置

配置文件：`configs/models.yaml`

### 支持的模型

| 模型名 | Provider | 说明 |
|--------|----------|------|
| `qwen` | dashscope | 通义千问 3.6-plus |
| `deepseek` | deepseek | DeepSeek Chat |
| `minimax` | minimax | MiniMax M2.7 |
| `doubao-seed` | volcengine | 豆包 Seed 2.0 Pro |
| `doubao-seed-2.0-lite` | volcengine | 豆包 Seed 2.0 Lite |
| `doubao-seed-1.6` | volcengine | 豆包 Seed 1.6 |
| `doubao-seed-2.0-pro-v2` | volcengine | 豆包 Seed 2.0 Pro V2 |

### 配置示例

```yaml
models:
  deepseek:
    enabled: true
    provider: deepseek
    config:
      api_endpoint: "https://api.deepseek.com/v1"
      model: "deepseek-chat"
      api_key_env: "DEEPSEEK_API_KEY"
      parameters:
        temperature: 0.7
        top_p: 0.9
        max_tokens: 4096
```

### 参数说明

- `max_tokens`: 最大输出长度（建议 4096-8192）
- `temperature`: 生成随机性（0.7 推荐）
- `top_p`: 采样范围（0.9 推荐）

---

## 数据集说明

### 目录结构

```
datasets/
  python/
    python_code_files_self_contained/
      boundary/           # 边界条件测试
      simple_function/    # 简单函数测试
      complex_dependency/ # 复杂依赖测试
      interface_mock/     # 接口模拟测试
  java/
  go/
  cpp/
```

### 数据集类别

| 类别 | 说明 | 适用语言 |
|------|------|---------|
| `self_contained` | 自包含代码 | Python, Go, Java, C++ |
| `module_level` | 模块级别(需要workspace上下文) | (待扩展) |

**当前数据集均为 `self_contained`**，默认 `--class self_contained` 即可正常工作。若后续新增真正的 module_level 样本（需要 meta.json 和 workspace_root），则需指定 `--class module_level`。

### 场景类型

- `boundary`: 边界条件测试
- `simple_function`: 简单函数
- `complex_dependency`: 复杂依赖
- `interface_mock`: 接口模拟

---

## 输出结果

### 目录结构

```
artifacts/
  runs/
    <run-id>/
      generated/
        tests/           # 生成的测试文件
        metadata/        # 元数据（响应、统计）
        generated_manifest.json
      evaluation/
        evaluation_result.json
      report/
        report.html      # HTML 报告
        report_summary.json
      run.log            # 运行日志
      api.log            # API 调用日志
```

### 报告内容

HTML 报告包含：

1. **模型排名**：综合得分排序
2. **按语言统计**：各语言表现
3. **按场景统计**：各场景表现
4. **截断分析**：截断率、续写统计、调优建议
5. **错误分析**：失败案例详情
6. **图表分析**：可视化对比

### 关键指标

- **编译通过率**：生成的测试能否编译
- **测试通过率**：测试是否正确
- **行覆盖率**：代码覆盖程度
- **变异得分**：测试质量（杀死变异比例）

---

## 常见问题

### Q: 截断问题如何处理？

系统已内置自动续写功能：
- 检测到截断时自动发送续写请求
- 最多续写 3 次
- 报告中显示截断统计和调优建议

如截断率过高，建议：
1. 增加 `max_tokens`（如 8192）
2. 简化提示词
3. 选择更大上下文的模型

### Q: Docker 命令在 Windows 如何运行？

PowerShell 格式：
```powershell
docker run --rm --env-file .env `
  -v "${PWD}/datasets:/app/datasets" `
  -v "${PWD}/artifacts:/app/artifacts" `
  -v "${PWD}/configs:/app/configs" `
  -w /app `
  utbench:latest run --models deepseek --langs python
```

### Q: 如何增量运行？

使用 `--mode incremental`：
```bash
utbench run --mode incremental --models deepseek --langs python
```

系统会跳过已完成的样本（基于 checkpoint）。

### Q: 变异测试太慢怎么办？

1. 减少样本数：`--max-samples 5`
2. 缩短超时：`--mutation-timeout 180`
3. 关闭变异：去掉 `--mutation-enabled`

### Q: 如何查看详细日志？

日志文件位置：
- `artifacts/runs/<run-id>/run.log`
- `artifacts/runs/<run-id>/api.log`

---

## 示例命令

### 快速测试（2个样本）

```bash
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  -w /app \
  utbench:latest run \
    --models deepseek \
    --langs python \
    --config /app/configs/models.yaml \
    --dataset-root /app/datasets \
    --max-samples 2 \
    --class module_level
```

### 全量评测（所有模型+语言）

```bash
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  -w /app \
  utbench:latest run \
    --models deepseek,qwen,minimax,doubao-seed \
    --langs python,java,go,cpp \
    --config /app/configs/models.yaml \
    --dataset-root /app/datasets \
    --mutation-enabled \
    --mutation-timeout 360
```

### 仅生成（调试用）

```bash
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  -w /app \
  utbench:latest generate \
    --models deepseek \
    --langs cpp \
    --config /app/configs/models.yaml \
    --dataset-root /app/datasets \
    --max-samples 1 \
    --class self_contained \
    --dry-run
```