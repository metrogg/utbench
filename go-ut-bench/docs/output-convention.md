# 输出目录规范

## 1. 目标

统一 run 目录、checkpoint、数据库和 Agent 过程产物的落盘方式，避免历史 `results_*` 一类临时目录继续扩散。

## 2. 顶层目录

默认约定：

- 输出根目录：`./artifacts`
- 数据库目录：`./storage`

完整结构：

```text
artifacts/
  checkpoints/
    runner_<scope_hash>.checkpoint.json
  runs/
    <run-id>/
      generated/
        generated_manifest.json
        tests/
          <subject-id>/
            <language>/
              <sample-id>.test<ext>
        prompts/
          rendered/
            <subject-id>/
              <language>/
                <sample-id>.prompt.txt
        metadata/
          <subject-id>_<language>_<sample-id>.response.json
          <subject-id>_<language>_<sample-id>.metadata.json
          agent_traces/
            <subject-id>/
              <language>/
                <sample-id>.trace.jsonl
                <sample-id>.diff.json
      evaluation/
        evaluation_result.json
      report/
        report_summary.json
        report.html
      agent_workspaces/
        <subject-id>/
          <language>/
            <sample-id>/
      logs/
      run_summary.json

storage/
  utbench.db
```

## 3. `generated/` 目录

### `generated_manifest.json`

生成阶段主清单，包含：

- `spec`
- `prompt_strategy`
- `prompt_version_id`
- `cases[]`

每个 case 现在不仅有旧的 `model` 字段，还会带：

- `subject_id`
- `subject_kind`
- `agent_framework`
- `agent_model`
- `skill_name`
- `skill_version`
- `sample_uid`
- `subject_version_id`
- `generation_key`
- `dependency_fingerprint`
- `generation_env_fingerprint`
- `trace_path`
- `workspace_diff_path`
- `sandbox_fingerprint`
- `reused`
- `reuse_stage`
- `reuse_key`
- `reuse_reason`
- `reused_from_run_id`
- `reused_from_case_id`

注意：为兼容旧 evaluator / reporter，`model` 字段当前仍填 `subject_id`。

复用命中时，当前 run 仍会写出完整 case 记录，并把 `reused=true`、`reuse_stage=generation` 写入 manifest。生成测试文件会进入当前 run 的 `generated/tests/...` 目录，DB 中同时保留来源 run/case，便于追溯历史。

### `evaluation_result.json`

每条 evaluation result 会额外带：

- `sample_uid`
- `evaluation_key`
- `evaluation_env_fingerprint`
- `evaluator_version`
- `mutation_config_sha256`
- `reused`
- `reuse_stage`
- `reuse_key`
- `reuse_reason`
- `reused_from_run_id`
- `reused_from_result_id`

第一阶段 `--reuse-evaluation` 默认关闭，因此这些字段主要用于记录和解释，暂不默认跳过重评测。

### `tests/`

最终生成出的测试代码，路径按 subject 和语言分层：

```text
generated/tests/<subject-id>/<language>/<sample-id>.test<ext>
```

### `prompts/rendered/`

实际渲染后发给模型或 Agent 的 prompt 快照：

```text
generated/prompts/rendered/<subject-id>/<language>/<sample-id>.prompt.txt
```

### `metadata/*.response.json`

保留原始模型响应或 CLI Agent 的执行摘要，例如：

- adapter 类型
- command
- exit_code
- stdout / stderr 摘要
- trace_path
- workspace_diff_path

### `metadata/*.metadata.json`

保存归一化元数据，例如：

- prompt 版本
- prompt 路径
- response 路径
- token 使用
- 生成延迟
- subject / framework / skill 信息
- `subject_version_id`
- `framework_config_sha256`
- `skill_sha256`
- `agent_command_sha256`
- `docker_image`
- `docker_image_digest`
- `env_contract_sha256`

### `metadata/agent_traces/`

Agent 特有过程产物：

- `*.trace.jsonl`：一条或多条执行轨迹记录
- `*.diff.json`：工作区变更摘要

当前 `trace.jsonl` 主要记录：

- `subject_id`
- `framework`
- `model`
- `skill`
- `sample_id`
- `language`
- `command`
- `exit_code`
- `duration_ms`
- `stdout`
- `stderr`

## 4. `agent_workspaces/` 目录

这是 CLI Agent 的临时工作区快照根目录：

```text
agent_workspaces/<subject-id>/<language>/<sample-id>/
```

它体现的是当前真正的 Agent 隔离粒度：

```text
一个 subject 执行一个 sample，对应一个独立 workspace
```

当 framework 使用 `sandbox_mode: docker` 时，这个 workspace 会被单独挂载进内层容器的 `/workspace`。

## 5. `evaluation/` 目录

`evaluation_result.json` 保留原有评测结果，并新增 Agent 相关字段：

- `subject_id`
- `subject_kind`
- `agent_framework`
- `agent_model`
- `skill_name`
- `skill_version`
- `trace_path`
- `workspace_diff_path`
- `sandbox_fingerprint`
- `evaluation_env_fingerprint`

## 6. `report/` 目录

### `report_summary.json`

除了原有聚合统计，现在还包含：

- `agent_comparisons`
- `skill_uplifts`

### `report.html`

HTML 报告继续沿用现有页面结构，但数据源已经支持 subject 维度聚合。

## 7. checkpoint 作用域

runner checkpoint 现在不再只看模型名，而是按 subject 作用域计算。也就是：

```text
subject_id + language + sample
```

这样可以避免：

- `agent + no_skill` 误复用 `agent + skill`
- `model_api` 误复用 `cli_agent`

## 8. 数据库

SQLite 已扩展 subject 维度字段，`generated_cases` 和 `evaluation_results` 都会保存：

- `subject_id`
- `subject_kind`
- `agent_framework`
- `agent_model`
- `skill_name`
- `skill_version`
- `sandbox_fingerprint`

旧字段 `model` 仍保留兼容。

## 9. 当前边界

当前输出规范已经覆盖：

- 纯模型 baseline
- CLI Agent
- skill 注入
- trace / diff / sandbox fingerprint

但还没有做到：

- 多步 tool 调用的细粒度结构化 trace
- syscall / 网络 / 文件系统级别的完整审计
- 真正不可逃逸的 microVM 级隔离证明
