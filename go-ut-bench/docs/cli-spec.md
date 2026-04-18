# CLI 命令设计（MVP）

## 1. 命令列表

- `utbench run`
- `utbench generate`
- `utbench evaluate`
- `utbench report`
- `utbench ingest`
- `utbench dataset list`
- `utbench dataset validate`
- `utbench dataset index`
- `utbench dataset manifest`

## 2. 命令说明

### 2.1 run

一键编排：`generate -> evaluate -> report`，可选 `--ingest`。

示例：

```bash
utbench run --models deepseek,minimax --langs python,go --class self_contained --dataset-manifest ./configs/dataset_index.json --dataset-root ./datasets --output-root ./artifacts --mode full --max-samples 20 --ingest --db ./storage/utbench.db
```

支持参数：

- `--config`：模型配置文件（默认 `../benchmark/config/models.yaml`）
- `--dry-run`：跳过模型 API 调用，使用占位测试代码

### 2.2 generate

仅生成单测和 manifest。

示例：

```bash
utbench generate --models deepseek --langs python --class module_level --dataset-manifest ./configs/dataset_index.json --dataset-root ./datasets --output-root ./artifacts --config ../benchmark/config/models.yaml
```

### 2.3 evaluate

输入已生成好的 `generated_manifest.json`，仅做评测。

示例：

```bash
utbench evaluate --input ./artifacts/runs/<run-id>/generated/generated_manifest.json --output-root ./artifacts
```

说明：

- `evaluate` 与 `generate` 完全解耦，可以直接评测历史生成产物

### 2.4 report

输入评测结果，输出报告。

示例：

```bash
utbench report --input ./artifacts/runs/<run-id>/evaluation/evaluation_result.json --output-root ./artifacts
```

### 2.5 ingest

输入评测结果，入库 SQLite。

示例：

```bash
utbench ingest --input ./artifacts/runs/<run-id>/evaluation/evaluation_result.json --db ./storage/utbench.db
```

### 2.6 dataset list

查看可被扫描到的样本。

示例：

```bash
utbench dataset list --langs python,java,go,cpp --class self_contained --dataset-root ./datasets --max-samples 20
```

### 2.7 dataset validate

校验数据集基础布局。

示例：

```bash
utbench dataset validate --dataset-root ./datasets
```

### 2.8 dataset index

扫描目录并生成标准索引。

示例：

```bash
utbench dataset index --dataset-root ./datasets --output ./configs/dataset_index.json
```

### 2.9 dataset manifest

基于索引构建预过滤清单（可选，方便分享）。

示例：

```bash
utbench dataset manifest --index ./configs/dataset_index.json --level l1 --output ./configs/dataset_l1.json --langs python,java,go,cpp --limit-per-scenario 20
```

> 直接跑不需要这一层，直接用 `dataset_index.json` + CLI 过滤参数即可。

## 3. 参数语义

- `--models`：模型名称列表（逗号分隔）
- `--langs`：语言列表（逗号分隔）
- `--class`：样本大类（`self_contained`/`module_level`）
- `--scenario`：样本子类（`boundary`/`simple_function`/`complex_dependency`/`interface_mock`）
- `--dataset-manifest`：样本清单（默认 `./configs/dataset_index.json`，用 CLI 参数过滤）
- `--mode`：执行模式（`full`/`incremental`）
- `--reset-checkpoint`：重置当前作用域 checkpoint
- `--mutation-enabled`：是否启用变异阶段
- `--mutation-timeout`：变异阶段超时（秒）
- `--mutation-policy`：变异异常策略（`warn`/`fail`）

说明：

- `warn`：记录 `mutation_error` 并继续
- `fail`：评测阶段返回非零退出码
- `--max-samples`：样本上限，0 表示不限制
- `--output-root`：输出根目录
- `--db`：SQLite 文件路径
