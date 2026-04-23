# 本地 CLI 快速上手

本文档面向直接在仓库根目录运行 `utbench` 的场景，命令已按当前代码实现核对。

## 1. 环境准备

最小要求：

- Go 1.21+
- 仓库根目录下有 `datasets/`、`configs/models.yaml`

如果要跑真实评测而不是 dry-run，还需要安装对应语言工具链。

Python:

```bash
pip install pytest coverage mutmut
```

Go 变异测试：

```bash
go install github.com/go-gremlins/gremlins/cmd/gremlins@latest
```

Java:

- JDK 17+
- Maven 3+

C++:

- CMake
- GoogleTest
- gcov
- mull

## 2. 构建可执行文件

Linux / macOS:

```bash
go build -o utbench ./cmd/utbench/
```

Windows PowerShell:

```powershell
go build -o utbench.exe ./cmd/utbench/
```

查看帮助：

```bash
./utbench --help
```

```powershell
.\utbench.exe --help
```

提示：当前仓库根目录运行时，请优先显式传入 `--config ./configs/models.yaml`。CLI 自带的默认路径是 `../benchmark/config/models.yaml`。

## 3. 先跑一个 dry-run

```bash
./utbench run \
  --models deepseek \
  --langs python \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --config ./configs/models.yaml \
  --class self_contained \
  --scenario simple_function \
  --max-samples 1 \
  --dry-run
```

这条命令会依次执行：

1. `generate`
2. `evaluate`
3. `report`

并在 `artifacts/runs/<run-id>/` 下生成完整产物。

## 4. 真实运行示例

```bash
export DEEPSEEK_API_KEY="..."
export DASHSCOPE_API_KEY="..."

./utbench run \
  --models deepseek,qwen \
  --langs python,go \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --config ./configs/models.yaml \
  --class self_contained \
  --max-samples 5
```

`--models` 的可选值请以 `configs/models.yaml` 中的键名为准，例如：

- `deepseek`
- `qwen`
- `minimax`
- `doubao-seed`
- `glm-4.7`

## 5. 分步执行

如果希望 `generate / evaluate / report` 的结果都落在同一个目录，必须手动复用同一个 `--run-id`。

```bash
RUN_ID=local_step_001

./utbench generate \
  --run-id "$RUN_ID" \
  --models deepseek \
  --langs python \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --config ./configs/models.yaml \
  --max-samples 2 \
  --dry-run

./utbench evaluate \
  --run-id "$RUN_ID" \
  --manifest ./artifacts/runs/$RUN_ID/generated/generated_manifest.json \
  --output-root ./artifacts

./utbench report \
  --run-id "$RUN_ID" \
  --input ./artifacts/runs/$RUN_ID/evaluation/evaluation_result.json \
  --output-root ./artifacts

./utbench ingest \
  --input ./artifacts/runs/$RUN_ID/evaluation/evaluation_result.json \
  --db-path ./storage/utbench.db
```

不要写错：

- `evaluate --manifest`
- `report --input`
- `ingest --db-path`

## 6. 数据集相关命令

索引全部样本：

```bash
./utbench dataset index \
  --dataset-root ./datasets \
  --output ./configs/dataset_index.json
```

从索引构建清单：

```bash
./utbench dataset manifest \
  --index ./configs/dataset_index.json \
  --langs python,go \
  --class self_contained \
  --level l1 \
  --limit-per-scenario 20 \
  --output ./configs/dataset_l1.json
```

校验布局并做依赖扫描：

```bash
./utbench dataset stats --dataset-root ./datasets
```

说明：

- 当前没有 `dataset validate` 子命令
- `dataset stats` 在 Windows 下会调用 `python3`；如果系统只有 `python`，建议使用 WSL / Docker 或自行建立别名

## 7. 结果目录

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
```

如果是单独执行 `report` 且没有传 `--run-id`，CLI 会创建一个新的 run 目录，而不是回写到原来的目录。

## 8. 常见坑

- Windows 下直接执行仓库里已有的 `utbench` 文件可能会失败，因为那通常不是当前平台重新构建出的可执行文件；请先 `go build`
- 报告目录名是 `report/`，不是 `reports/`
- `--dataset-manifest` 和 `--level` 任意一个生效时，样本发现会优先走 manifest，而不是直接扫目录
- 只做 dry-run 时不需要 API Key，但依然需要可读取的 `configs/models.yaml`
