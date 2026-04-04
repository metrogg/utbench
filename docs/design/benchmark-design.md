# Benchmark 评测方案总体设计

## 1. 目标与范围

ut-bench 的目标是建立一套可复现、可扩展、可横向对比的 AI 单测生成评测流水线。

当前已实现范围：

- Runner 阶段：数据集读取、Prompt 构造、模型 API 调用、响应解析、测试代码落盘
- 并发执行：支持多任务并发
- 断点续跑：基于 checkpoint 的任务跳过与恢复
- 结果追踪：按模型分目录输出，保留可追溯元数据

规划中范围：

- Evaluator 阶段：编译验证、测试执行、覆盖率统计、变异测试
- Reporter 阶段：多维聚合与可视化报告

## 2. 设计原则

- 标准化：统一输入（dataset）、统一配置（models.yaml）、统一输出目录（results）
- 可重复：相同输入 + 相同配置应得到可比结果
- 容错性：单样本失败不阻断全局任务
- 可扩展：新增模型、语言、指标时不破坏现有流程
- 可观测：输出进度日志、失败原因、摘要文件

## 3. 总体架构

```text
dataset/*
  -> Runner
       - sample reader
       - prompt builder
       - model client
       - response parser
       - output writer
  -> results/<model>/{tests,reports,artifacts}
  -> results/runner_summary_*.json
```

## 4. 流水线阶段定义

### 阶段一：Runner（已实现）

输入：

- `dataset/<lang>/*` 样本代码
- `benchmark/config/models.yaml` 模型配置
- 环境变量中的 API Key

输出：

- `results/<model>/tests/*.test.<ext>`
- `results/<model>/reports/*.metadata.json`
- `results/<model>/artifacts/*.response.json`（可按需清理）
- `results/runner_summary_*.json`

关键能力：

- Prompt 模板化构造（含语言框架、覆盖率目标、上下文信息）
- provider 级协议适配（OpenAI 兼容与 DashScope 差异）
- 重试与退避（针对可恢复错误）
- 并发执行 + checkpoint 续跑

### 阶段二到五（规划）

- 编译验证：统计编译通过率
- 测试执行：统计执行通过率
- 覆盖率分析：行/分支/函数覆盖率
- 变异测试：Mutation Score

### 阶段六（规划）

- Reporter：汇总模型、语言、场景维度结果并生成报告

## 5. 配置设计

核心配置文件：`benchmark/config/models.yaml`

- `models.<name>.enabled`：是否参与评测
- `models.<name>.provider`：模型供应方
- `models.<name>.config.api_endpoint`：接口地址
- `models.<name>.config.model`：模型 ID
- `models.<name>.config.api_key_env`：API Key 环境变量名
- `models.<name>.config.parameters`：推理参数
- `benchmark.parallel`：并发参数（模型并发、样本并发）
- `benchmark.timeouts.api_call`：API 超时阈值

## 6. 执行模型

默认入口：

```bash
python -m benchmark.runner
```

常用参数：

- `--model`：指定模型列表
- `--lang`：指定语言列表
- `--max-samples`：限制每语言样本数
- `--max-workers`：覆盖总并发线程数
- `--no-resume`：关闭断点续跑
- `--reset-checkpoint`：重置当前范围 checkpoint
- `--dry-run`：不调用真实 API

## 7. 异常处理策略

- `auth_config_error`：Key 缺失或配置错误，立即失败
- `http_non_retryable`：不可恢复请求错误（如参数错误）
- `http_retryable`：可恢复错误（如 429/5xx），自动重试
- `timeout`：请求超时，可重试
- `network_error`：网络异常，可重试

## 8. 结果目录规范

```text
results/
  <model>/
    tests/
    reports/
    artifacts/
  checkpoints/
  runner_summary_*.json
```

命名规则（Runner）：

- 测试代码：`{model}_{lang}_{sample}_{timestamp}.test.{ext}`
- 元数据：`{model}_{lang}_{sample}_{timestamp}.metadata.json`
- 原始响应：`{model}_{lang}_{sample}_{timestamp}.response.json`

## 9. 扩展点

- 新模型：在 `models.yaml` 新增配置项
- 新语言：扩展 dataset 子目录与语言映射
- 新指标：在 Evaluator 增加计算逻辑
- 新报告：在 Reporter 增加维度聚合与展示

## 10. 当前已知限制

- 目前仅完成生成阶段，不代表最终单测质量评分
- 长响应模型可能出现超时，需要结合并发与 token 参数调优
- 不同 provider 协议差异需要持续维护适配器
