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
```
