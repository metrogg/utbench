# ut-bench

AI 单元测试生成横向评测工具（当前重点是 Runner 生成阶段）。

## 项目现状

当前仓库已落地：

- 数据集：`dataset/`（`java` / `python` / `go` / `cpp`）
- 模型配置：`benchmark/config/models.yaml`
- 生成模块：`benchmark/runner/`
  - Prompt 构造（含中英双语说明）
  - 模型 API 调用
  - 响应解析与测试代码落盘
  - 并发执行
  - 断点续跑（checkpoint）
  - 按模型实时进度日志
- 评估模块：`benchmark/evaluator/`
  - 编译验证 / 测试执行 / 覆盖率 / 变异测试（Python 已接入）
  - 输出 `evaluator_summary_*.json`
- 报告模块：`benchmark/reporter/`
  - 基于 evaluator 结果进行多维聚合
  - 输出 JSON / CSV / HTML 报告

## 目录说明

```text
ut-bench/
  README.md
  dataset/
  benchmark/
    config/models.yaml
    runner/
      runner.py
      prompt_builder.py
      __main__.py
    evaluator/
    reporter/
  results/
```

详细使用手册：`docs/使用说明.md`

## 环境准备

- Python 3.10+
- 依赖：`pyyaml`

安装依赖：

```bash
pip install pyyaml
```

或使用脚本一键初始化：

```bash
bash scripts/setup.sh
```

## 模型配置

配置文件：`benchmark/config/models.yaml`

关键字段：

- `enabled`：是否参与运行
- `provider`：模型提供方
- `config.api_endpoint`：接口基地址
- `config.model`：模型 ID
- `config.api_key_env`：读取 API Key 的环境变量名

### 当前主要模型（按你当前配置）

- `qwen` -> `qwen3.6-plus`
- `deepseek` -> `deepseek-chat`
- `minimax` -> `MiniMax-M2.7`
- `doubao-seed` -> `doubao-seed-2.0-pro`
- `glm-4.7` -> `glm-4.7`
- `kimi-k2.5` 当前已禁用
- `glm`（旧配置）当前已禁用

## 设置环境变量（PowerShell）

```powershell
$env:DASHSCOPE_API_KEY="你的 key"
$env:DEEPSEEK_API_KEY="你的 key"
$env:MINIMAX_API_KEY="你的 key"
$env:VOLCENGINE_API_KEY="你的 key"
```

## 运行方式

### 1) 全量运行（所有启用模型 × 全语言数据集）

```bash
python -m benchmark.runner
```

也可使用一键流水线脚本（runner -> evaluator -> reporter）：

```bash
bash scripts/run_benchmark.sh
```

### 2) 指定模型/语言

```bash
python -m benchmark.runner --model doubao-seed,glm-4.7 --lang python,java
```

### 2.1) 指定样本子集（按 glob）

例如仅跑当前 Python 基线集中 1 条：

```bash
python -m benchmark.runner \
  --lang python \
  --sample-glob "*_python_*_0.py" \
  --max-samples 1
```

### 3) 并发 + 断点续跑

```bash
# 默认开启断点续跑，按配置并发
python -m benchmark.runner

# 指定并发线程数
python -m benchmark.runner --max-workers 8

# 清空当前范围 checkpoint 后重跑
python -m benchmark.runner --reset-checkpoint

# 关闭续跑（强制全跑）
python -m benchmark.runner --no-resume
```

### 4) 调试模式（不调用真实 API）

```bash
python -m benchmark.runner --dry-run --model doubao-seed --lang python --max-samples 1
```

## 运行输出

输出目录：`results/`

- `results/<model>/tests/`：生成的单测代码
- `results/<model>/reports/`：每条样本 metadata
- `results/<model>/artifacts/`：原始响应/失败产物
- `results/checkpoints/`：断点续跑状态
- `results/runner_summary_*.json`：本轮汇总
- `results/evaluator_summary_*.json`：评测汇总
- `results/reports/reporter_*.{json,csv,html}`：报告产物

## Evaluator 与 Reporter

运行 evaluator：

```bash
python -m benchmark.evaluator --results-root results
```

只看自包含白名单样本（推荐用于主对比口径）：

```bash
python -m benchmark.evaluator --results-root results --lang python --only-self-contained
```

说明：Python 样本会按 `benchmark/config/python_self_contained_allowlist.txt` 自动分层，
汇总中包含 `sample_bucket_counts`、`self_contained_*` 和 `non_self_contained_*` 指标。

运行 reporter（默认读取最新 evaluator summary）：

```bash
python -m benchmark.reporter --results-root results
```

中文可视化报告（摘要/详细分析/附录 + 图表）：

- 柱状图：模型关键指标横向对比
- 折线图：复杂度维度趋势
- 雷达图：五维能力分布
- 热力图：模型 × 语言表现

指定 evaluator summary：

```bash
python -m benchmark.reporter --results-root results --evaluator-summary results/evaluator_summary_20260406T064045963256Z.json
```

可配置参数示例：

```bash
python -m benchmark.reporter \
  --results-root results \
  --formats json,csv,html \
  --chart-style teal \
  --threshold-test 0.7 \
  --threshold-line 0.7 \
  --threshold-branch 0.6 \
  --threshold-mutation 0.85
```

也可通过脚本一键生成：

```bash
bash scripts/gen_report.sh --results-root results
```

## 进度日志示例

```text
[INFO] Runner start: models=[...], total_samples=15, pending_tasks=75, ...
[INFO] OK model=deepseek lang=python sample=boundary_python_boundary_0
[INFO] PROGRESS model=deepseek done=3/15 success=3 failed=0
[ERROR] FAILED model=qwen lang=go sample=... err={...}
```

## 常见问题

1. 报 `Missing API key env var`
   - 说明 `api_key_env` 对应的环境变量未设置。

2. 某模型大量 `timeout`
   - 先降低并发：`--max-workers 4`
   - 先小样本验证：`--max-samples 1`
   - 必要时先禁用该模型再跑全量。

3. 需要从中断点继续
   - 直接重跑同一命令即可（默认 `resume=true`）。

