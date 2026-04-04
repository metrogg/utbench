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
  results/
```

## 环境准备

- Python 3.10+
- 依赖：`pyyaml`

安装依赖：

```bash
pip install pyyaml
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

### 2) 指定模型/语言

```bash
python -m benchmark.runner --model doubao-seed,glm-4.7 --lang python,java
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

---

如果你后续要接 evaluator/reporter，我可以继续把编译、执行、覆盖率、变异测试和报告汇总补齐到同一条流水线命令。
