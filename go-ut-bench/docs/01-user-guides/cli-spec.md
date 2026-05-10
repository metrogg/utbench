# CLI 命令文档

## 命令概览

```
utbench run          完整流程 (generate -> evaluate -> report)
utbench generate     仅生成单元测试
utbench evaluate     评测已生成的单元测试
utbench report       生成评测报告
utbench db           管理 SQLite 评测数据库
utbench assets       查询可复用生成/评测资产
utbench dataset      数据集管理 (index, manifest, stats, validate)
utbench doctor       评测工具链自检
utbench web          启动 Web 管理界面
```

---

## 1. run

一键编排完整流程。

**示例：**
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

**参数：**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--config` | 历史默认值 `../benchmark/config/models.yaml` | 当前仓库实际应显式传 `./configs/models.yaml` |
| `--models` | 空 | 模型列表（逗号分隔）；为空时加载 `models.yaml` 中所有启用模型 |
| `--subjects` | 空 | subject 列表（逗号分隔）；为空时至少包含默认 `model_api__<model>__no_skill` |
| `--langs` | 空 | 语言列表（逗号分隔）；为空时使用全部支持语言 |
| `--dataset-root` | `./datasets` | 数据集根目录 |
| `--dataset-manifest` | 空 | 数据集清单路径 |
| `--output-root` | `./artifacts` | 输出根目录 |
| `--class` | `self_contained` | 数据集大类：`self_contained`/`repo_level` |
| `--scenario` | 全部 | 数据集场景：`boundary`/`simple_function`/`complex_dependency`/`interface_mock` |
| `--level` | 空 | 数据集级别 |
| `--max-samples` | `0`（不限制） | 样本数量上限 |
| `--mode` | `full` | 执行模式：`full`/`incremental` |
| `--reuse-generated` | `true` | 允许优先复用历史生成结果 |
| `--reuse-evaluation` | `false` | 允许在环境指纹匹配时复用历史评测结果 |
| `--mutation-enabled` | `true` | 启用变异测试 |
| `--mutation-timeout` | `600` | 变异超时（秒） |
| `--mutation-policy` | `warn` | 变异异常策略：`warn`/`fail` |
| `--test-timeout` | `180` | 测试执行超时（秒） |
| `--workers` | `16` | 并发 worker 数量 |
| `--eval-backend` | `local` | 评测后端：`local`/`docker` |
| `--eval-docker-image` | `utbench:latest` | docker 评测镜像 |
| `--dry-run` | `false` | 跳过 API 调用 |
| `--reset-checkpoint` | `false` | 重置 checkpoint |
| `--ingest` | `false` | 完成后写入 v2 SQLite 数据库 |
| `--db-path` | `./storage/utbench.db` | SQLite 数据库路径 |
| `-v` | `false` | 打开详细日志 |

---

## 2. generate

仅生成单元测试，不进行评测。

**示例：**
```bash
utbench generate \
  --models deepseek \
  --langs python \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --max-samples 5 \
  --dry-run
```

**参数：**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--config` | 历史默认值 `../benchmark/config/models.yaml` | 当前仓库实际应显式传 `./configs/models.yaml` |
| `--models` | 空 | 模型列表；为空时加载所有启用模型 |
| `--subjects` | 空 | subject 列表 |
| `--langs` | 空 | 语言列表；为空时使用全部支持语言 |
| `--dataset-root` | `./datasets` | 数据集根目录 |
| `--dataset-manifest` | 空 | 数据集清单路径 |
| `--output-root` | `./artifacts` | 输出根目录 |
| `--class` | `self_contained` | 数据集大类 |
| `--scenario` | 全部 | 数据集场景 |
| `--level` | 空 | 数据集级别 |
| `--max-samples` | `0` | 样本数量上限 |
| `--mode` | `full` | 执行模式 |
| `--reuse-generated` | `true` | 允许优先复用历史生成结果 |
| `--db-path` | `./storage/utbench.db` | SQLite 数据库路径 |
| `--dry-run` | `false` | 跳过 API 调用 |
| `--reset-checkpoint` | `false` | 重置 checkpoint |
| `-v` | `false` | 打开详细日志 |

---

## 3. evaluate

评测已生成的单元测试。

**示例：**
```bash
utbench evaluate \
  --manifest ./artifacts/runs/<run-id>/generated/generated_manifest.json \
  --mutation-enabled \
  --mutation-timeout 360
```

说明：`evaluate` 会把结构化日志写到 `artifacts/runs/<run-id>/logs/`，便于排查 compile/test/coverage/mutation 卡点。所有语言都必须样本级测试通过后才运行变异测试；基线测试失败的样本变异分记 0，并归入模型问题。

**参数：**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--manifest` | 必填 | 生成的 manifest 文件路径 |
| `--output-root` | `./artifacts` | 输出根目录 |
| `--mutation-enabled` | `true` | 启用变异测试 |
| `--mutation-timeout` | `600` | 变异超时（秒） |
| `--mutation-policy` | `warn` | 变异异常策略 |
| `--test-timeout` | `180` | 测试超时（秒） |
| `--db-path` | `./storage/utbench.db` | SQLite 数据库路径 |
| `--reuse-evaluation` | `false` | 允许复用历史评测结果 |
| `--eval-backend` | `local` | 评测后端：`local`/`docker` |
| `--eval-docker-image` | `utbench:latest` | docker 评测镜像 |
| `-v` | `false` | 打开详细日志 |

---

## 4. report

生成评测报告。

**示例：**
```bash
utbench report \
  --evaluation ./artifacts/runs/<run-id>/evaluation/evaluation_result.json
```

说明：报告中的 `test_pass_rate` / `avg_test_pass_rate` 为样本级口径；`test_case_pass_rate` / `avg_test_case_pass_rate` 为测试用例级口径。

报告的 `failures` 会把变异测试失败拆成更细的类型，避免全部混在 `mutation_error` 中：

| 类型 | 含义 |
|------|------|
| `mutation_skipped_baseline_failed` | 基线测试未全部通过，跳过变异测试 |
| `mutation_target_not_exercised` | 生成测试没有导入/执行被测模块，变异体无法被测试关联 |
| `mutation_no_results` | 变异工具没有产出可解析结果 |
| `mutation_no_coverage` | 工具产出结果但没有 killed/survived 有效计分项 |
| `mutation_no_effective_mutants` | 变异体为 0 或没有执行任何变异体 |
| `mutation_timeout` | 变异测试超时 |
| `mutation_tool_error` | 工具解析、报告或插件问题 |
| `mutation_error` | 未归类的变异测试错误 |

**参数：**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--evaluation` | 必填 | 评测结果 JSON 文件路径 |
| `--output-root` | `./artifacts` | 输出根目录 |
| `--run-id` | 空 | 输出报告的 run ID；为空时自动生成 |

---

## 5. db

管理 SQLite 评测数据库。v2 schema 会保存 manifest、生成测试代码、prompt、模型响应、评测结果、报告和 artifact 索引；旧 `utbench ingest` 两表结构不再作为入口。

**示例：**
```bash
utbench db init --db-path ./storage/utbench.db

utbench db ingest-run \
  --run-id <run-id> \
  --output-root ./artifacts \
  --db-path ./storage/utbench.db

utbench db overview --db-path ./storage/utbench.db
utbench db list-results --run-id <run-id> --db-path ./storage/utbench.db

utbench db report \
  --run-ids <run-a>,<run-b> \
  --models deepseek,qwen \
  --langs python,go \
  --output-root ./artifacts \
  --db-path ./storage/utbench.db
```

**参数：**

| 子命令 | 说明 |
|------|------|
| `init` | 初始化数据库 schema |
| `ingest-manifest --manifest <path>` | 入库 `generated_manifest.json` 和关联生成产物 |
| `ingest-evaluation --evaluation <path>` | 入库 `evaluation_result.json`，并尽量追溯 manifest |
| `ingest-report --report <path>` | 入库 `report_summary.json` 和 HTML 报告 |
| `ingest-run --run-id <id>` | 入库 `artifacts/runs/<run-id>` 下所有已知产物 |
| `overview` | 输出数据库计数和最近运行 |
| `list-runs` | 列出数据库中的运行 |
| `list-results` | 按 run/model/lang 查询样本级评测结果 |
| `report` | 基于数据库筛选结果生成横向对比报告 |

常用参数：

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--db-path` | `./storage/utbench.db` | SQLite 数据库路径 |
| `--run-id` | 空 | `ingest-run` / `list-results` 的运行 ID |
| `--run-ids` | 空 | `report` 的来源 run ID 列表，逗号分隔 |
| `--evaluation-run-ids` | 空 | `report` 的来源 evaluation run ID 列表，逗号分隔 |
| `--models` | 空 | `report` 的模型过滤，逗号分隔 |
| `--langs` | 空 | `report` 的语言过滤，逗号分隔 |
| `--run-dir` | 空 | 直接指定 run 目录 |
| `--output-root` | `./artifacts` | `ingest-run --run-id` 推导 run 目录时使用 |
| `--json` | `false` | `overview` / `list-runs` / `list-results` 输出 JSON |

---

## 6. dataset

数据集管理操作。

### dataset index

生成数据集索引文件。

```bash
utbench dataset index \
  --dataset-root ./datasets \
  --output ./configs/dataset_index.json
```

**参数：**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--dataset-root` | `./datasets` | 数据集根目录 |
| `--output` | `./configs/dataset_index.json` | 输出索引文件路径 |

### dataset stats

查看数据集统计信息。

```bash
utbench dataset stats \
  --dataset-root ./datasets
```

**参数：**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--dataset-root` | `./datasets` | 数据集根目录 |

### dataset validate

检查数据集 readiness，统计 language/class/scenario 分布并标记高风险样本。

```bash
utbench dataset validate \
  --dataset-root ./datasets \
  --langs python,go,java,cpp \
  --class self_contained \
  --strict
```

**参数：**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--dataset-root` | `./datasets` | 数据集根目录 |
| `--langs` | 全部 | 语言列表 |
| `--class` | 全部 | 数据集类别过滤 |
| `--scenario` | 全部 | 场景过滤 |
| `--strict` | `false` | 存在错误时返回非 0 |
| `--json` | 空 | 写出 JSON 报告 |

---

## 7. doctor

检查评测环境是否能正常编译、运行测试、收集覆盖率和执行变异测试。

```bash
utbench doctor \
  --langs python,go,java,cpp \
  --mutation-enabled \
  --mutation-timeout 120
```

**参数：**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--langs` | `python,go,java,cpp` | 要检查的语言 |
| `--mutation-enabled` | `true` | 是否运行变异测试 canary |
| `--mutation-timeout` | `120` | 变异测试超时（秒） |
| `--test-timeout` | `60` | canary 测试超时（秒） |
| `--json` | 空 | 写出 JSON 报告 |

---

## 8. assets

查询 SQLite 中可复用的 subject / generation / evaluation 资产，以及复用命中解释。

```bash
utbench assets explain-reuse \
  --db-path ./storage/utbench.db \
  --subject model_api__deepseek__no_skill \
  --lang python \
  --sample sample_001
```

常用子命令：

| 子命令 | 说明 |
|------|------|
| `subjects` | 查看 subject 级资产汇总 |
| `generations` | 查看生成资产 |
| `evaluations` | 查看评测资产 |
| `explain-reuse` | 解释某个 subject/sample 当前是否能命中复用 |

---

## 数据集类别说明

| 类别 | 说明 | 适用语言 |
|------|------|---------|
| `self_contained` | 自包含代码，无外部依赖 | Python, Go, Java, C++ |
| `repo_level` | 仓库级别，有外部依赖 | 预留/旧数据 |

**注意**：当前仓库内置数据集实际为 `self_contained`，正式运行建议显式指定：
```bash
--class self_contained
```

---

## 模型列表

支持的模型（逗号分隔）：

| 模型名 | Provider |
|--------|----------|
| `deepseek` | deepseek |
| `qwen` | dashscope |
| `minimax` | minimax |
| `doubao-seed` | volcengine |
| `doubao-seed-2.0-lite` | volcengine |
| `doubao-seed-1.6` | volcengine |
| `doubao-seed-2.0-pro-v2` | volcengine |

---

## 输出文件

每次运行生成以下文件：

```
artifacts/runs/<run-id>/
  generated/
    tests/                    # 生成的测试文件
    metadata/                 # 元数据（响应、统计）
    generated_manifest.json   # 生成清单
  evaluation/
    evaluation_result.json    # 评测结果
  report/
    report.html               # HTML 报告
    report_summary.json       # 报告摘要
  run.log                     # 运行日志
  api.log                     # API 调用日志
```

## 说明

- 当前代码里 `run/generate` 的 `--config` 默认值仍保留历史路径，因此日常运行建议始终显式传 `--config ./configs/models.yaml`。
- 当前项目的统一被测对象不是单纯 model，而是 `subject = framework + model + optional skill`。

