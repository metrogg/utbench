<<<<<<< HEAD
# CLI 命令说明

以下内容以当前 `cmd/utbench/main.go` 为准。

## 命令总览

```text
utbench run
utbench generate
utbench evaluate
utbench report
utbench ingest
utbench dataset index
utbench dataset manifest
utbench dataset stats
```

## 约定说明

- `run` 和 `generate` 共享同一组样本选择 / 输出 / 模型配置参数
- 如果单独执行 `generate`、`evaluate`、`report`，想把结果写进同一个目录，必须复用同一个 `--run-id`
- CLI 默认模型配置路径是 `../benchmark/config/models.yaml`
- 在当前仓库根目录运行时，建议显式传入 `--config ./configs/models.yaml`

## 1. `utbench run`

完整流水线：`generate -> evaluate -> report`，可选 `--ingest`。

示例：

```bash
utbench run \
  --run-id demo_run_001 \
  --models deepseek,qwen \
  --langs python,go \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --config ./configs/models.yaml \
  --class self_contained \
  --scenario boundary \
  --max-samples 5 \
  --dry-run
```

参数：

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--config` | 模型配置文件路径 | `../benchmark/config/models.yaml` |
| `--models` | 模型键名，逗号分隔 | 空 |
| `--langs` | 语言列表，逗号分隔 | 空 |
| `--dataset-root` | 数据集根目录 | `./datasets` |
| `--dataset-manifest` | 数据集 manifest 路径 | 空 |
| `--class` | 数据集类别，逗号分隔，可选 `self_contained,module_level` | 空 |
| `--scenario` | 场景过滤，逗号分隔 | 空 |
| `--level` | manifest level 过滤 | 空 |
| `--max-samples` | 每个语言/场景的最大样本数，`0` 表示不限制 | `0` |
| `--output-root` | 输出根目录 | `./artifacts` |
| `--mode` | 运行模式，`full` 或 `incremental` | `full` |
| `--mutation-enabled` | 是否开启变异测试 | `false` |
| `--mutation-timeout` | 变异测试超时秒数 | `1800` |
| `--mutation-policy` | 变异错误策略，`warn` 或 `fail` | `warn` |
| `--dry-run` | 跳过真实 API 调用 | `false` |
| `--reset-checkpoint` | 清空 runner checkpoint | `false` |
| `--ingest` | 流程结束后写入 SQLite | `false` |
| `--db-path` | SQLite 文件路径 | `./storage/utbench.db` |
| `--run-id` | 本次运行 ID | 自动生成 |
| `--workers` | worker 数量，`0` 为自动 | `0` |
| `-v` | 详细日志 | `false` |

## 2. `utbench generate`

只生成测试，不做评测和报告。

示例：

```bash
utbench generate \
  --run-id demo_gen_001 \
=======
# CLI 命令文档

## 命令概览

```
utbench run          完整流程 (generate -> evaluate -> report)
utbench generate     仅生成单元测试
utbench evaluate     评测已生成的单元测试
utbench report       生成评测报告
utbench ingest       结果导入 SQLite
utbench dataset      数据集管理 (index, manifest, stats, validate)
utbench doctor       评测工具链自检
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
| `--config` | `../benchmark/config/models.yaml` | 模型配置文件路径 |
| `--models` | `deepseek` | 模型列表（逗号分隔） |
| `--langs` | `python` | 语言列表（逗号分隔） |
| `--dataset-root` | `./datasets` | 数据集根目录 |
| `--dataset-manifest` | `./configs/dataset_index.json` | 数据集清单路径 |
| `--output-root` | `./artifacts` | 输出根目录 |
| `--class` | `self_contained` | 数据集大类：`self_contained`/`module_level` |
| `--scenario` | 全部 | 数据集场景：`boundary`/`simple_function`/`complex_dependency`/`interface_mock` |
| `--level` | `l1` | 数据集级别 |
| `--max-samples` | `0`（不限制） | 样本数量上限 |
| `--mode` | `full` | 执行模式：`full`/`incremental` |
| `--mutation-enabled` | `true` | 启用变异测试 |
| `--mutation-timeout` | `360` | 变异超时（秒） |
| `--mutation-policy` | `warn` | 变异异常策略：`warn`/`fail` |
| `--total-timeout` | `0`（不限制） | 总超时（分钟） |
| `--dry-run` | `false` | 跳过 API 调用 |
| `--reset-checkpoint` | `false` | 重置 checkpoint |
| `--ingest` | `false` | 完成后导入 SQLite |
| `--db` | `./storage/utbench.db` | SQLite 数据库路径 |
| `--verbose` | `true` | 详细日志 |

---

## 2. generate

仅生成单元测试，不进行评测。

**示例：**
```bash
utbench generate \
>>>>>>> origin/feat/go
  --models deepseek \
  --langs python \
  --dataset-root ./datasets \
  --output-root ./artifacts \
<<<<<<< HEAD
  --config ./configs/models.yaml \
  --max-samples 2 \
  --dry-run
```

说明：

- 接受与 `run` 相同的共享参数
- `--ingest` / `--db-path` 虽然会被解析，但对 `generate` 本身没有效果
- 输出主文件是 `artifacts/runs/<run-id>/generated/generated_manifest.json`

## 3. `utbench evaluate`

读取生成结果并做编译 / 测试 / 覆盖率 / 变异评测。

示例：

```bash
utbench evaluate \
  --run-id demo_gen_001 \
  --manifest ./artifacts/runs/demo_gen_001/generated/generated_manifest.json \
  --output-root ./artifacts \
  --mutation-enabled
```

参数：

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--manifest` | 必填，`generated_manifest.json` 路径 | 无 |
| `--output-root` | 输出根目录 | `./artifacts` |
| `--run-id` | 运行 ID | 自动生成 |
| `--mutation-enabled` | 是否开启变异测试 | `false` |
| `--mutation-timeout` | 变异测试超时秒数 | `1800` |
| `--mutation-policy` | 变异错误策略 | `warn` |
| `-v` | 详细日志 | `false` |

注意：当前命令使用 `--manifest`，不是 `--input`。

## 4. `utbench report`

根据 `evaluation_result.json` 生成汇总 JSON 和 HTML 报告。

示例：

```bash
utbench report \
  --run-id demo_gen_001 \
  --input ./artifacts/runs/demo_gen_001/evaluation/evaluation_result.json \
  --output-root ./artifacts
```

参数：

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--input` | 必填，`evaluation_result.json` 路径 | 无 |
| `--output-root` | 输出根目录 | `./artifacts` |
| `--run-id` | 运行 ID | 自动生成 |
| `-v` | 详细日志 | `false` |

默认输出：

```text
artifacts/runs/<run-id>/report/report_summary.json
artifacts/runs/<run-id>/report/report.html
```

## 5. `utbench ingest`

把评测结果写入 SQLite。

示例：

```bash
utbench ingest \
  --input ./artifacts/runs/demo_gen_001/evaluation/evaluation_result.json \
  --db-path ./storage/utbench.db
```

参数：

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--input` | 必填，`evaluation_result.json` 路径 | 无 |
| `--db-path` | SQLite 路径 | `./storage/utbench.db` |

## 6. `utbench dataset index`

扫描 `datasets/` 并输出索引文件。

示例：

```bash
utbench dataset index \
  --dataset-root ./datasets \
  --output ./configs/dataset_index.json
```

参数：

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--dataset-root` | 数据集根目录 | `./datasets` |
| `--output` | 索引输出路径 | `./configs/dataset_index.json` |

## 7. `utbench dataset manifest`

根据索引生成筛选后的 manifest。

示例：

```bash
utbench dataset manifest \
  --index ./configs/dataset_index.json \
  --langs python,go \
  --class self_contained \
  --scenario boundary \
  --level l1 \
  --limit-per-scenario 20 \
  --output ./configs/dataset_l1.json
```

参数：

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--index` | 必填，索引文件路径 | 无 |
| `--langs` | 语言列表，逗号分隔 | 空 |
| `--class` | 类别过滤 | 空 |
| `--scenario` | 场景过滤 | 空 |
| `--level` | level 过滤 | 空 |
| `--limit-per-scenario` | 每个场景保留的样本数 | `20` |
| `--output` | 必填，manifest 输出路径 | 无 |

## 8. `utbench dataset stats`

校验数据集布局并做依赖扫描。

示例：

```bash
utbench dataset stats --dataset-root ./datasets
```

参数：

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--dataset-root` | 数据集根目录 | `./datasets` |

说明：

- 这是当前唯一的数据集校验命令
- 当前实现会调用 `python3` 做依赖扫描；Windows 下如无 `python3` 命令，建议用 WSL / Docker

## 9. 产物目录

```text
artifacts/
  checkpoints/
    runner_<hash>.checkpoint.json
  runs/
    <run-id>/
      generated/
        generated_manifest.json
        tests/
        metadata/
      evaluation/
        evaluation_result.json
      report/
        report_summary.json
        report.html
      run_summary.json
=======
  --max-samples 5 \
  --dry-run
```

**参数：**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--config` | `../benchmark/config/models.yaml` | 模型配置文件路径 |
| `--models` | `deepseek` | 模型列表 |
| `--langs` | `python` | 语言列表 |
| `--dataset-root` | `./datasets` | 数据集根目录 |
| `--dataset-manifest` | `./configs/dataset_index.json` | 数据集清单路径 |
| `--output-root` | `./artifacts` | 输出根目录 |
| `--class` | `self_contained` | 数据集大类 |
| `--scenario` | 全部 | 数据集场景 |
| `--level` | `l1` | 数据集级别 |
| `--max-samples` | `0` | 样本数量上限 |
| `--mode` | `full` | 执行模式 |
| `--dry-run` | `false` | 跳过 API 调用 |
| `--reset-checkpoint` | `false` | 重置 checkpoint |
| `--verbose` | `true` | 详细日志 |

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

**参数：**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--manifest` | 必填 | 生成的 manifest 文件路径 |
| `--output-root` | `./artifacts` | 输出根目录 |
| `--mutation-enabled` | `true` | 启用变异测试 |
| `--mutation-timeout` | `360` | 变异超时（秒） |
| `--mutation-policy` | `warn` | 变异异常策略 |
| `--test-timeout` | `0` | 测试超时（秒） |
| `--verbose` | `true` | 详细日志 |

---

## 4. report

生成评测报告。

**示例：**
```bash
utbench report \
  --input ./artifacts/runs/<run-id>/evaluation/evaluation_result.json
```

**参数：**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--input` | 必填 | 评测结果 JSON 文件路径 |
| `--output-root` | `./artifacts` | 输出根目录 |

---

## 5. ingest

将评测结果导入 SQLite 数据库。

**示例：**
```bash
utbench ingest \
  --input ./artifacts/runs/<run-id>/evaluation/evaluation_result.json \
  --db ./storage/utbench.db
```

**参数：**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--input` | 必填 | 评测结果 JSON 文件路径 |
| `--db` | `./storage/utbench.db` | SQLite 数据库路径 |

---

## 6. dataset

数据集管理操作。

### dataset index

生成数据集索引文件。

```bash
utbench dataset index \
  --root ./datasets \
  --output ./configs/dataset_index.json
```

**参数：**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--root` | `./datasets` | 数据集根目录 |
| `--output` | `./configs/dataset_index.json` | 输出索引文件路径 |

### dataset stats

查看数据集统计信息。

```bash
utbench dataset stats \
  --manifest ./configs/dataset_index.json
```

**参数：**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--manifest` | `./configs/dataset_index.json` | 数据集索引文件路径 |

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

## 数据集类别说明

| 类别 | 说明 | 适用语言 |
|------|------|---------|
| `self_contained` | 自包含代码，无外部依赖 | Python, Go, Java, C++ |
| `module_level` | 模块级别，有外部依赖 | 预留/旧数据 |

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
>>>>>>> origin/feat/go
```
