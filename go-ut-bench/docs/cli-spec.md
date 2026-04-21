# CLI 命令文档

## 命令概览

```
utbench run          完整流程 (generate -> evaluate -> report)
utbench generate     仅生成单元测试
utbench evaluate     评测已生成的单元测试
utbench report       生成评测报告
utbench ingest       结果导入 SQLite
utbench dataset      数据集管理 (index, manifest, stats)
```

---

## 1. run

一键编排完整流程。

**示例：**
```bash
utbench run \
  --models deepseek,minimax \
  --langs python,go \
  --dataset-manifest ./configs/dataset_index.json \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --class self_contained \
  --scenario boundary \
  --max-samples 10 \
  --mutation-enabled
```

**参数：**

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--config` | 模型配置文件路径 | `../benchmark/config/models.yaml` |
| `--models` | 模型列表（逗号分隔） | - |
| `--langs` | 语言列表（逗号分隔） | - |
| `--dataset-root` | 数据集根目录 | `./datasets` |
| `--dataset-manifest` | 数据集清单路径 | `./configs/dataset_index.json` |
| `--class` | 数据集大类 | - |
| `--scenario` | 数据集场景 | - |
| `--level` | 数据集级别 | - |
| `--max-samples` | 样本数量上限 | 0（不限制） |
| `--output-root` | 输出根目录 | `./artifacts` |
| `--mode` | 执行模式 (`full`/`incremental`) | `full` |
| `--mutation-enabled` | 启用变异测试 | false |
| `--mutation-timeout` | 变异超时（秒） | 1800 |
| `--mutation-policy` | 变异异常策略 (`warn`/`fail`) | `warn` |
| `--dry-run` | 跳过 API 调用 | false |
| `--reset-checkpoint` | 重置 checkpoint | false |
| `--ingest` | 完成后导入 SQLite | false |
| `--db-path` | SQLite 数据库路径 | `./storage/utbench.db` |
| `--run-id` | 运行 ID | 自动生成 |

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
  --class self_contained \
  --scenario boundary \
  --max-samples 5 \
  --dry-run
```

**参数：**

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--config` | 模型配置文件路径 | `../benchmark/config/models.yaml` |
| `--models` | 模型列表 | - |
| `--langs` | 语言列表 | - |
| `--dataset-root` | 数据集根目录 | `./datasets` |
| `--dataset-manifest` | 数据集清单路径 | `./configs/dataset_index.json` |
| `--class` | 数据集大类 | - |
| `--scenario` | 数据集场景 | - |
| `--level` | 数据集级别 | - |
| `--max-samples` | 样本数量上限 | 0（不限制） |
| `--output-root` | 输出根目录 | `./artifacts` |
| `--mode` | 执行模式 | `full` |
| `--dry-run` | 跳过 API 调用 | false |
| `--reset-checkpoint` | 重置 checkpoint | false |
| `--run-id` | 运行 ID | 自动生成 |

---

## 3. evaluate

评测已生成的单元测试。

**示例：**
```bash
utbench evaluate \
  --manifest ./artifacts/runs/run_xxx/generated/generated_manifest.json \
  --output-root ./artifacts \
  --run-id eval_001 \
  --mutation-enabled
```

**参数：**

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--manifest` | **必需** - generated_manifest.json 路径 | - |
| `--output-root` | 输出根目录 | `./artifacts` |
| `--run-id` | 运行 ID | 自动生成 |
| `--mutation-enabled` | 启用变异测试 | false |
| `--mutation-timeout` | 变异超时（秒） | 1800 |
| `--mutation-policy` | 变异异常策略 | `warn` |
| `-v` | 详细输出 | false |

---

## 4. report

生成评测报告。

**示例：**
```bash
utbench report \
  --input ./artifacts/runs/run_xxx/evaluation/evaluation_result.json \
  --output-root ./artifacts
```

---

## 5. ingest

导入评测结果到 SQLite。

**示例：**
```bash
utbench ingest \
  --input ./artifacts/runs/run_xxx/evaluation/evaluation_result.json \
  --db-path ./storage/utbench.db
```

---

## 6. dataset

数据集管理子命令。

### 6.1 dataset index

扫描数据集生成索引文件。

```bash
utbench dataset index --dataset-root ./datasets --output ./configs/dataset_index.json
```

### 6.2 dataset manifest

基于索引生成过滤后的清单。

```bash
utbench dataset manifest \
  --index ./configs/dataset_index.json \
  --langs python,go,java,cpp \
  --class self_contained \
  --level l1 \
  --limit-per-scenario 20 \
  --output ./configs/manifest_l1.json
```

### 6.3 dataset stats

查看数据集统计。

```bash
utbench dataset stats --dataset-root ./datasets
```

---

## 数据集结构

```
datasets/
  python/
    python_code_files_self_contained/
      boundary/
      simple_function/
      complex_dependency/
      interface_mock/
    python_code_files_module_level/
      ...
  java/
    java_code_files_self_contained/
      ...
  go/
    go_code_files_self_contained/
      ...
  cpp/
    cpp_code_files_self_contained/
      ...
```

---

## 输出结构

```
artifacts/
  runs/
    <run-id>/
      generated/
        generated_manifest.json
        tests/
          <model>/
            <lang>/
              *.test.<ext>
      evaluation/
        evaluation_result.json
      reports/
        reporter_*.{json,csv,html}
  runner_summary_<timestamp>.json
  evaluator_summary_<timestamp>.json
```
