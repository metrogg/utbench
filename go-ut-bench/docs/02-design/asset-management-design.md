# UT-Bench 资产管理、复用与扩展设计

## 目标

UT-Bench 的长期管理单位不是一次 run，而是可复用资产：

- `subject`: 用户视角的被测对象，例如 `model_api__deepseek-v4-flash__no_skill`、`opencode__deepseek-v4-flash__unit_test_skill`。
- `subject_version`: 复用判断使用的精确版本，包含模型配置、agent/framework 配置、skill 内容、sandbox 和镜像环境。
- `sample_uid`: 数据集样本的稳定身份，兼容当前单文件样本和未来仓库级样本。
- `generation_key`: 判断生成测试是否可复用。
- `evaluation_key`: 判断评测结果是否可复用。

文件系统继续保存大文件，SQLite 只保存索引、指纹、关联关系和查询入口。

## Subject 模型

`subject_id` 继续沿用现有格式：

```text
model_api__deepseek-v4-flash__no_skill
opencode__deepseek-v4-flash__no_skill
opencode__deepseek-v4-flash__unit_test_skill
claudecode__deepseek-v4-flash__no_skill
```

`subjects` 保存用户视角身份：

```text
subject_id
subject_kind
framework
model
skill
display_name
enabled
tags_json
created_at_utc
```

`subject_versions` 保存复用判断所需的精确配置：

```text
subject_version_id
subject_id
model_config_id
framework_config_sha256
skill_sha256
agent_command_sha256
docker_image
docker_image_digest
sandbox_fingerprint
env_contract_sha256
created_at_utc
```

复用必须使用 `subject_version_id` 和资产 key，不能只看 `subject_id`。

## Dataset Sample 模型

### 单文件样本

单文件样本指纹由这些字段组成：

```text
language
dataset_class = self_contained
scenario
sample_id
source_sha256
dependency_fingerprint
```

### 仓库级样本

仓库级样本从第一版资产层开始纳入设计：

```text
language
dataset_class = repo_level
scenario
sample_id
workspace_sha256
target_files_sha256
dependency_fingerprint
repo_manifest_sha256
```

仓库级样本需要支持这些元数据：

```text
workspace_root
target_files
test_command
build_command
dependency_files
lock_files
setup_commands
forbidden_paths
```

`sample_uid` 由 `language`、`dataset_class`、`scenario`、`sample_id`、源码或 workspace 指纹、repo manifest 指纹稳定生成。这样后续接入真实仓库、最小仓库、SWE-style 样本时，不需要重写资产层。

## Artifact 存储

文件仍然落在：

```text
artifacts/runs/<run-id>/
```

SQLite 保存：

```text
artifact_id
kind
path
sha256
size_bytes
created_at_utc
deleted_at_utc
```

关键 artifact 类型：

```text
source_code
repo_workspace_snapshot
rendered_prompt
generated_test
model_response
agent_trace
agent_session_export
workspace_diff
generated_manifest
evaluation_result
report_json
report_html
```

历史结果默认不物理删除，后续通过软删除或归档状态管理。

## Generation Reuse

`generation_key` 由以下信息组成：

```text
subject_version_id
sample_uid
prompt_rendering_sha256
prompt_version_id
language
dataset_class
dependency_fingerprint
generation_env_fingerprint
```

复用规则：

- 只复用 `success = true`。
- 默认取 `generated_at_utc` 最新的一条。
- 生成测试 artifact 必须存在且 sha256 匹配。
- 不跨 subject、skill、framework、docker image、sandbox fingerprint 复用。
- 仓库级样本 workspace 指纹不同则不复用。

当前实现入口为 `FindReusableGeneratedAsset(generation_key)`，适用于 `model_api`、`cli_agent` 以及未来的 agent adapter。

## Evaluation Reuse

`evaluation_key` 由以下信息组成：

```text
generated_test_sha256
sample_uid
evaluation_env_id
score_policy_id
evaluator_version
mutation_config_sha256
test_command_sha256
dependency_fingerprint
```

当前实现已经把 `evaluation_env_fingerprint` 纳入 `evaluation_key`，fingerprint 会记录本地工具链版本或 Docker digest，并与 `evaluator_version`、`mutation_config_sha256` 一起决定是否可复用。`--reuse-evaluation=true` 仍默认关闭，原因是正式横评前还需要继续观察不同语言 evaluator 指纹是否还有遗漏项。

评测复用的规则：

- 同 generated test、同 evaluator env、同 mutation config 才可命中。
- evaluator 或 mutation 配置变化必须重评测。
- 网络失败、环境失败、工具失败不作为可复用能力结果。
- 编译失败、测试失败可保存历史，但默认不作为跳过重评测的结果。

## Run 流程

新 run 的逻辑：

```text
resolve subjects
resolve samples
compute subject_version_id
compute sample_uid
compute prompt_rendering_sha256
compute generation_key
lookup reusable generation
  hit -> copy/link generated artifact into current run, mark reused
  miss -> execute generation, persist asset
compute evaluation_key
lookup reusable evaluation
  hit -> write reused evaluation row into current run
  miss -> evaluate, persist result
generate report from current run rows
ingest run artifacts
```

当前 run 的 manifest 保持完整，并为复用 case 写入：

```json
{
  "reused": true,
  "reused_from_run_id": "...",
  "reused_from_case_id": "...",
  "reuse_key": "...",
  "reuse_stage": "generation",
  "reuse_reason": "generation_key_match"
}
```

## 第一阶段落地范围

已纳入第一阶段：

- DB schema 增加 `subjects`、`subject_versions`、generation/evaluation key 字段和索引。
- 生成阶段默认开启 `--reuse-generated=true`。
- `FindReusableGeneratedAsset(generation_key)` 校验 artifact 存在和 sha256。
- `subject_versions` 已落库存储 `framework_config_sha256`、`skill_sha256`、`agent_command_sha256`、`docker_image_digest`、`env_contract_sha256`。
- manifest、metadata、DB 中保存复用来源和 key。
- 评测阶段记录 `evaluation_key`，并支持 `--reuse-evaluation=true` 的保守复用；默认仍为 `false`。
- generation/evaluation 都会在阶段开始前批量预判可复用资产，优先命中最新成功结果，再下发 worker。
- CLI 查询入口：`utbench assets subjects|generations|evaluations|explain-reuse`。
- Web “数据管理 → 资产审计” 已增加 subject 列表、subject version、生成资产、评测资产、复用解释。

后置增强：

- 更完整的 Web 树形管理页和历史展开视图。
- evaluation reuse 默认开启前的失败类型白名单。
- artifact 软删除、归档、引用计数和垃圾回收。
