# UTBench Database Design

## 目标

数据库用于管理评测全流程产生的可复用资产和可解释结果，而不是只保存一次运行的汇总分数。它需要支持：

- 持久化保存数据集样本、prompt、模型生成代码、模型响应、评测结果、报告和日志索引。
- 判断历史生成结果是否可复用，避免重复调用模型。
- 判断历史评测结果是否可复用，避免重复执行 compile/test/coverage/mutation。
- 从历史结果中按模型、语言、场景、数据集、prompt、评测环境筛选并生成横向对比报告。
- 保留完整溯源信息，能解释某个模型在某个样本上的分数来自哪份源码、哪版 prompt、哪份生成代码、哪个工具环境。

旧的 `runs` / `sample_results` 简化表不再作为长期设计基础，也不再提供兼容入口。新实现直接按本文 v2 schema 落地。

## 已确认决策

- artifact 路径统一保存相对路径，例如 `artifacts/runs/...`，避免换机器或换挂载目录后绝对路径失效。
- 成功和失败的模型响应都保存，失败响应同样要可溯源。
- ingest 时做基础脱敏，至少移除 API key、Authorization、Cookie、Set-Cookie 等敏感字段。
- P0 不做多用户，但 schema 不绑定本机绝对路径，主键和外键设计保留未来迁移服务端的空间。
- 删除策略先采用软删除，不物理删除 artifact；未来再加引用计数和垃圾回收。

## 当前产物与入库覆盖

当前运行已经会产生这些文件或字段，v2 schema 需要全部接住：

| 来源 | 当前数据 | 入库位置 |
|------|----------|----------|
| `generated_manifest.json` | run spec、prompt strategy/version、prompt snapshot dir、generated cases | `generation_runs`、`generated_cases` |
| `GeneratedCase` | model/language/sample_id/sample_path/prompt_path/generated_test_path/response_path/metadata_path/tokens/latency/truncated/error | `generated_cases`、`artifacts` |
| prompt snapshot | 每个模型/语言/样本的 rendered prompt txt | `prompt_renderings`、`artifacts` |
| response json | 模型原始响应或错误响应 | `artifacts`，并关联到 `generated_cases.response_artifact_id` |
| metadata json | prompt 信息、数据集 class、source md5、token、latency | `artifacts`，并关联到 `generated_cases.metadata_artifact_id` |
| `evaluation_result.json` | 每条样本 compile/test/coverage/mutation/assertion/failure_origin/score_eligible | `evaluation_runs`、`evaluation_results` |
| report summary/html | 聚合结果和 HTML 报告 | `report_snapshots`、`artifacts` |
| logs | api/errors/evaluator/run/runner 日志 | `run_artifacts` 或 `artifacts` |
| doctor / dataset validate | 环境自检和数据集 readiness 报告 | `quality_checks`、`artifacts` |

当前还不完整、建议补齐的产物：

- 每条评测各阶段的完整 stdout/stderr。现在 `EvaluationResult` 只保存错误摘要，成功阶段的原始输出主要在日志里，不利于逐样本复查。
- 每条评测实际执行的命令、工作目录、退出码、开始/结束时间。建议后续新增结构化 stage trace。
- Docker image digest。仅有 image tag 不够判断环境完全一致。
- dataset validation report 与 dataset fingerprint 的绑定。否则只知道样本 hash，不知道当时数据集 readiness 是否有 warning。

## 设计原则

1. SQLite 先行  
   当前先服务单机本地使用，默认数据库为 `storage/utbench.db`。schema 尽量避免 SQLite 特有技巧，方便未来迁移到 PostgreSQL 或服务端数据库。

2. 大文件不进数据库  
   生成测试代码、模型原始响应、prompt 快照、日志、报告 HTML/JSON 继续存放在文件系统。数据库只存 artifact 的路径、hash、大小、类型和关联关系。

3. 内容可寻址  
   关键资产都计算 `sha256`。是否可复用主要靠内容 hash 和环境 fingerprint，而不是靠 run_id。

4. 生成与评测解耦  
   同一份 generated test 可以在不同评测环境下评测多次；同一批评测结果可以被不同报告组合复用。

5. 结果可解释  
   每条 evaluation result 保留失败归因、剔除原因、原始 JSON、artifact 关联。报告排名不能只给分数，必须能追到原始记录。

6. 全信息优先，聚合后置  
   原始事实记录优先入库，报告中的排名和聚合数据可以重算。数据库中可以保存报告快照，但不能只保存报告快照。

## 核心实体

### schema_migrations

记录数据库 schema 版本。

```sql
CREATE TABLE schema_migrations (
  version INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  applied_at TEXT NOT NULL
);
```

### experiments

实验分组，用于把多次 generation/evaluation/report 归到同一个研究目标下。P0 可以在 CLI 中自动创建默认实验，Web UI 后续再暴露。

```sql
CREATE TABLE experiments (
  experiment_id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  tags_json TEXT,
  created_at TEXT NOT NULL,
  archived_at TEXT
);
```

### artifacts

保存所有外部文件的索引。

```sql
CREATE TABLE artifacts (
  artifact_id TEXT PRIMARY KEY,
  kind TEXT NOT NULL,
  sha256 TEXT NOT NULL,
  path TEXT NOT NULL,
  size_bytes INTEGER,
  mime TEXT,
  redacted INTEGER NOT NULL DEFAULT 0,
  deleted_at TEXT,
  created_at TEXT NOT NULL,
  UNIQUE(kind, sha256)
);
```

`kind` 建议值：

- `source_code`
- `generated_test`
- `rendered_prompt`
- `model_response`
- `metadata`
- `generated_manifest`
- `evaluation_result`
- `report_json`
- `report_html`
- `log`
- `doctor_json`
- `dataset_validation_json`
- `stage_stdout`
- `stage_stderr`
- `stage_trace_json`

### run_artifacts

把 artifact 关联到生成、评测或报告 run。

```sql
CREATE TABLE run_artifacts (
  run_artifact_id TEXT PRIMARY KEY,
  run_kind TEXT NOT NULL,
  run_id TEXT NOT NULL,
  artifact_id TEXT NOT NULL,
  role TEXT NOT NULL,
  created_at TEXT NOT NULL,
  UNIQUE(run_kind, run_id, artifact_id, role)
);
```

`run_kind` 建议值：`generation`、`evaluation`、`report`、`doctor`、`dataset_validate`。

### dataset_samples

保存数据集样本身份和源码 fingerprint。

```sql
CREATE TABLE dataset_samples (
  sample_uid TEXT PRIMARY KEY,
  language TEXT NOT NULL,
  dataset_class TEXT NOT NULL,
  scenario TEXT NOT NULL,
  sample_id TEXT NOT NULL,
  source_path TEXT NOT NULL,
  source_md5 TEXT,
  source_sha256 TEXT NOT NULL,
  source_artifact_id TEXT,
  risk_flags_json TEXT,
  created_at TEXT NOT NULL,
  UNIQUE(language, dataset_class, scenario, sample_id, source_sha256)
);
```

`sample_uid` 不能只用 `sample_id`，应由 `language/class/scenario/sample_id/source_sha256` 派生，避免不同语言或数据集版本冲突。

### dataset_snapshots

保存一次运行实际使用的数据集样本集合。

```sql
CREATE TABLE dataset_snapshots (
  dataset_snapshot_id TEXT PRIMARY KEY,
  dataset_root TEXT NOT NULL,
  dataset_fingerprint TEXT NOT NULL UNIQUE,
  validation_artifact_id TEXT,
  created_at TEXT NOT NULL
);
```

### dataset_snapshot_members

```sql
CREATE TABLE dataset_snapshot_members (
  dataset_snapshot_id TEXT NOT NULL,
  sample_uid TEXT NOT NULL,
  ordinal INTEGER,
  PRIMARY KEY(dataset_snapshot_id, sample_uid)
);
```

### model_configs

保存模型别名、真实模型 ID、provider 和生成参数。

```sql
CREATE TABLE model_configs (
  model_config_id TEXT PRIMARY KEY,
  model_alias TEXT NOT NULL,
  provider TEXT NOT NULL,
  model_id TEXT NOT NULL,
  config_sha256 TEXT NOT NULL,
  temperature REAL,
  max_tokens INTEGER,
  extra_json TEXT,
  created_at TEXT NOT NULL,
  UNIQUE(provider, model_id, config_sha256)
);
```

### prompt_profiles

保存 prompt 模板版本和公平规则版本。

```sql
CREATE TABLE prompt_profiles (
  prompt_profile_id TEXT PRIMARY KEY,
  prompt_version_id TEXT,
  prompt_strategy TEXT,
  template_sha256 TEXT NOT NULL,
  rules_sha256 TEXT NOT NULL,
  created_at TEXT NOT NULL,
  UNIQUE(prompt_version_id, prompt_strategy, template_sha256, rules_sha256)
);
```

### prompt_renderings

保存每个样本实际渲染出来的 prompt。

```sql
CREATE TABLE prompt_renderings (
  prompt_rendering_id TEXT PRIMARY KEY,
  prompt_profile_id TEXT NOT NULL,
  sample_uid TEXT NOT NULL,
  rendered_prompt_sha256 TEXT NOT NULL,
  artifact_id TEXT NOT NULL,
  created_at TEXT NOT NULL,
  UNIQUE(prompt_profile_id, sample_uid, rendered_prompt_sha256)
);
```

### generation_runs

保存一次生成批次的配置和范围。

```sql
CREATE TABLE generation_runs (
  generation_run_id TEXT PRIMARY KEY,
  experiment_id TEXT,
  run_id TEXT NOT NULL UNIQUE,
  started_at TEXT,
  finished_at TEXT,
  status TEXT NOT NULL,
  dataset_snapshot_id TEXT,
  dataset_fingerprint TEXT NOT NULL,
  prompt_profile_id TEXT,
  model_config_set_sha256 TEXT,
  spec_json TEXT NOT NULL,
  git_sha TEXT,
  manifest_artifact_id TEXT
);
```

### generated_cases

保存模型对某样本生成的测试代码记录。

```sql
CREATE TABLE generated_cases (
  generated_case_id TEXT PRIMARY KEY,
  generation_run_id TEXT NOT NULL,
  sample_uid TEXT NOT NULL,
  model_config_id TEXT NOT NULL,
  prompt_rendering_id TEXT,
  success INTEGER NOT NULL,
  truncated INTEGER NOT NULL DEFAULT 0,
  latency_ms INTEGER,
  prompt_tokens INTEGER,
  completion_tokens INTEGER,
  total_tokens INTEGER,
  generated_test_sha256 TEXT,
  generated_test_artifact_id TEXT,
  response_artifact_id TEXT,
  metadata_artifact_id TEXT,
  error_json TEXT,
  created_at TEXT NOT NULL,
  UNIQUE(sample_uid, model_config_id, prompt_rendering_id, generated_test_sha256)
);
```

生成复用的核心查询条件：

```text
sample_uid + model_config_id + prompt_rendering_id
```

如果已有成功生成结果，且用户允许复用，就可以跳过模型 API 调用。

### evaluation_envs

保存评测环境 fingerprint。

```sql
CREATE TABLE evaluation_envs (
  eval_env_id TEXT PRIMARY KEY,
  env_fingerprint TEXT NOT NULL UNIQUE,
  docker_image TEXT,
  docker_image_digest TEXT,
  utbench_git_sha TEXT,
  utbench_version TEXT,
  os TEXT,
  python_version TEXT,
  go_version TEXT,
  java_version TEXT,
  cpp_compiler_version TEXT,
  pytest_version TEXT,
  coverage_version TEXT,
  mutmut_version TEXT,
  go_mutesting_version TEXT,
  pitest_version TEXT,
  mull_version TEXT,
  tool_versions_json TEXT,
  created_at TEXT NOT NULL
);
```

`env_fingerprint` 至少由以下内容组成：

- Docker image digest 或本地 OS/tool versions
- utbench git sha
- evaluator 版本
- 各语言工具版本
- mutation gate 策略
- score policy

### score_policies

保存排名公式和剔除规则版本。

```sql
CREATE TABLE score_policies (
  score_policy_id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  policy_sha256 TEXT NOT NULL,
  compile_weight REAL NOT NULL,
  test_weight REAL NOT NULL,
  coverage_weight REAL NOT NULL,
  mutation_weight REAL NOT NULL,
  rules_json TEXT NOT NULL,
  created_at TEXT NOT NULL,
  UNIQUE(name, policy_sha256)
);
```

### evaluation_runs

保存一次评测批次。

```sql
CREATE TABLE evaluation_runs (
  evaluation_run_id TEXT PRIMARY KEY,
  experiment_id TEXT,
  run_id TEXT NOT NULL UNIQUE,
  generation_run_id TEXT,
  manifest_artifact_id TEXT,
  eval_env_id TEXT NOT NULL,
  score_policy_id TEXT NOT NULL,
  mutation_enabled INTEGER NOT NULL,
  mutation_timeout INTEGER,
  test_timeout INTEGER,
  started_at TEXT,
  finished_at TEXT,
  status TEXT NOT NULL,
  spec_json TEXT NOT NULL,
  evaluation_artifact_id TEXT
);
```

### evaluation_results

保存单条样本评测结果。

```sql
CREATE TABLE evaluation_results (
  evaluation_result_id TEXT PRIMARY KEY,
  evaluation_run_id TEXT NOT NULL,
  generated_case_id TEXT,
  sample_uid TEXT NOT NULL,
  model_config_id TEXT NOT NULL,
  language TEXT NOT NULL,
  sample_id TEXT NOT NULL,

  compile_pass INTEGER NOT NULL,
  test_pass INTEGER,
  test_pass_count INTEGER,
  test_total_count INTEGER,
  test_pass_rate REAL,

  line_coverage REAL,
  branch_coverage REAL,

  mutation_tool TEXT,
  mutation_score REAL,
  mutation_total INTEGER,
  mutation_killed INTEGER,
  mutation_survived INTEGER,
  mutation_no_tests INTEGER,
  mutation_timeouts INTEGER,
  mutation_skipped INTEGER,
  mutation_suspicious INTEGER,

  assertion_count INTEGER,
  test_case_count INTEGER,
  assertion_density REAL,
  runtime_ms INTEGER,

  failure_origin TEXT,
  score_eligible INTEGER,
  score_exclusion_reason TEXT,

  compile_error TEXT,
  test_error TEXT,
  coverage_error TEXT,
  mutation_error TEXT,

  raw_result_json TEXT NOT NULL,

  UNIQUE(evaluation_run_id, model_config_id, sample_uid)
);
```

### evaluation_stage_results

保存单条样本各阶段的结构化执行记录。P0 可以先只在发生错误时写 stdout/stderr artifact；后续逐步做到每阶段全量记录。

```sql
CREATE TABLE evaluation_stage_results (
  stage_result_id TEXT PRIMARY KEY,
  evaluation_result_id TEXT NOT NULL,
  stage TEXT NOT NULL,
  tool TEXT,
  command_json TEXT,
  workdir TEXT,
  started_at TEXT,
  finished_at TEXT,
  runtime_ms INTEGER,
  exit_code INTEGER,
  timed_out INTEGER NOT NULL DEFAULT 0,
  stdout_artifact_id TEXT,
  stderr_artifact_id TEXT,
  trace_artifact_id TEXT,
  error_summary TEXT
);
```

`stage` 建议值：`prepare`、`compile`、`test`、`coverage`、`mutation`、`cleanup`。

评测复用的核心查询条件：

```text
generated_case_id + eval_env_id + score_policy_id + mutation_enabled + mutation_timeout + test_timeout
```

如果完全匹配，就可以跳过评测。

### report_snapshots

保存一次报告选择范围和输出 artifact。

```sql
CREATE TABLE report_snapshots (
  report_id TEXT PRIMARY KEY,
  created_at TEXT NOT NULL,
  title TEXT,
  filters_json TEXT NOT NULL,
  score_policy_id TEXT NOT NULL,
  result_set_sha256 TEXT NOT NULL,
  report_json_artifact_id TEXT,
  report_html_artifact_id TEXT
);
```

### report_result_members

保存报告由哪些 evaluation result 组成。

```sql
CREATE TABLE report_result_members (
  report_id TEXT NOT NULL,
  evaluation_result_id TEXT NOT NULL,
  PRIMARY KEY(report_id, evaluation_result_id)
);
```

## 报告查询方式

报告不再只能基于单个 `evaluation_result.json`。数据库报告应支持条件筛选：

- 模型列表
- 语言列表
- 场景列表
- 数据集版本 / dataset fingerprint
- prompt profile
- model config
- evaluation env
- score policy
- 时间范围

示例：

```text
选择 prompt_profile=A、eval_env=B、score_policy=C 下，
deepseek-v4-flash / minimax2.7 / doubao-seed-2.0-code
在 python/go/java/cpp 的 self_contained 样本结果，
生成横向排名报告。
```

## 日志策略

P0 不把完整日志写入数据库。

做法：

- 日志文件继续写入 `artifacts/runs/<run-id>/logs/`
- `artifacts` 表记录日志路径、hash、大小、类型
- `evaluation_runs` 关联日志 artifact

后续如果需要日志检索，再增加：

```sql
CREATE TABLE log_events (
  log_event_id TEXT PRIMARY KEY,
  run_id TEXT NOT NULL,
  stage TEXT,
  level TEXT,
  message TEXT NOT NULL,
  fields_json TEXT,
  created_at TEXT
);
```

## CLI 建议

新增或扩展以下命令：

```bash
utbench db init --db ./storage/utbench.db
utbench db ingest-manifest --manifest artifacts/runs/<run>/generated/generated_manifest.json
utbench db ingest-evaluation --evaluation artifacts/runs/<run>/evaluation/evaluation_result.json
utbench db ingest-report --report artifacts/runs/<run>/report/report_summary.json
utbench db report --models deepseek-v4-flash,minimax2.7 --langs python,go --env <fingerprint> --output-root artifacts
utbench db list-runs
utbench db list-results --model deepseek-v4-flash --language python
```

不保留旧 `utbench ingest` 兼容语义。实现时可以删除或重定向为 v2 ingest，避免用户误以为旧表仍是权威数据源。

## 复用规则

### 生成复用

允许复用的必要条件：

- `sample_uid` 相同。
- `model_config_id` 相同。
- `prompt_rendering_id` 相同。
- 已有 `generated_cases.success = 1`。
- 生成代码 artifact 仍存在且 sha256 校验通过。

可选条件：

- 如果用户指定只复用同一 `generation_run_id`，则限制在该 run 内。
- 如果用户允许跨 run 复用，则按内容 hash 选最新成功记录。

### 评测复用

允许复用的必要条件：

- `generated_case_id` 相同。
- `eval_env_id` 相同。
- `score_policy_id` 相同。
- `mutation_enabled`、`mutation_timeout`、`test_timeout` 相同。
- 相关 artifact 仍存在且 sha256 校验通过。

不允许复用的情况：

- Docker digest 或工具版本变化。
- evaluator git sha 变化且该变化影响评测逻辑。
- score policy 或 failure attribution 规则变化。
- 数据集源码 hash 变化。

## 数据完整性约束

ingest 阶段必须做这些校验：

- manifest cases 与 generated artifact 数量一致。
- evaluation results 的 `(model, language, sample_id)` 必须能在 manifest 中找到。
- 同一个 evaluation run 中不能有重复 `(model_config_id, sample_uid)`。
- source 和 generated test artifact 的 sha256 必须可计算。
- `score_eligible=false` 时 `failure_origin` 必须是 `environment`、`dataset` 或 `tool`。
- `compile_pass=false` 时 `test_pass` 应为 NULL。
- `test_pass=false` 且 mutation enabled 时 mutation error 应为 baseline skip，不能真实运行 mutation。
- report snapshot 的 members 必须全部来自同一个 score policy；如允许跨环境对比，报告必须显式展示环境差异。

## 需要补充到运行产物的内容

为了达到“全信息可溯源”，后续实现时建议补这些字段或文件：

- 每条 generated case 的 `model_config_id` 或 model config hash。
- 每条 generated case 的 rendered prompt sha256，目前可以从 prompt file 计算，但 manifest 中最好直接带上。
- 每条 evaluation result 的 `generated_case_id` 或 generated test sha256，目前只能通过 path/model/language/sample_id 反推。
- 每条 evaluation result 的 stage trace artifact，至少包含命令、退出码、stdout/stderr 摘要。
- evaluation run 的 env fingerprint artifact，即 doctor 输出或 tool version snapshot。
- report snapshot 的 filter JSON 和 result membership。

## P0 实施计划

1. 新建 v2 schema 和迁移版本表。
2. 实现 artifact indexing：保存 path、sha256、kind。
3. 实现 manifest ingest：写入 dataset samples、prompt profiles、prompt renderings、model configs、generation runs、generated cases。
4. 实现 evaluation ingest：写入 evaluation env、score policy、evaluation run、evaluation results。
5. 实现 DB 查询生成报告：先复用现有 reporter 聚合逻辑，把 DB rows 转成 `EvaluationResult`。
6. Web UI 增加历史结果选择入口。

## 当前对齐结论

- 新库直接按 v2 schema 实现，不兼容旧 `runs/sample_results` 入口。
- artifact 保存相对路径，文件内容用 sha256 校验。
- 成功和失败响应都保存。
- ingest 做基础脱敏。
- P0 不做多用户，但 schema 保留 `experiments`，未来可自然扩展到 `projects/workspaces/users`。
- 删除先软删除 artifact，不做物理清理。

## 当前落地状态

代码已落到 `internal/store/sqlite.go` 和 Web 的“数据管理”页：

- `utbench db init`
- `utbench db ingest-manifest`
- `utbench db ingest-evaluation`
- `utbench db ingest-report`
- `utbench db ingest-run`
- `utbench db overview`
- `utbench db list-runs`
- `utbench db list-results`
- `utbench db report`

`run --ingest --db-path ...` 已改为写入 v2 schema，会补录当前 run 目录中的 manifest、evaluation、report 和关联 artifact。Web 当前支持数据库概览、运行列表、样本级评测结果、artifact 索引、已有 run 补录，以及按当前筛选生成数据库报告。`utbench db report` 支持从数据库选择 run/model/language 后复用现有 reporter 生成对比报告；后续还需要在 UI 上补更完整的多选控件和环境一致性提示。

## 仍需最终确认

1. 是否需要记录每个样本每个阶段的完整 stdout/stderr，还是只在失败时保存？  
   推荐：P0 失败时保存完整 stdout/stderr，成功时保存摘要和命令；后续如果磁盘可接受，再改全量保存。

2. 是否允许报告跨评测环境对比？  
   推荐：允许，但默认报告必须按同一 `eval_env_id` 过滤；跨环境报告需要明显标注，不参与正式排名。

3. 是否需要实验分组命名？  
   推荐：CLI 默认用 run_id 创建实验；Web UI 后续允许用户给实验命名，如“P0 prompt baseline”、“go-mutesting migration”。
