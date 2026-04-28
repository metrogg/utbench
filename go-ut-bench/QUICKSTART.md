# go-ut-bench Docker 快速开始

这份文档只覆盖当前仓库真实存在的 Docker 用法，不再引用已经不存在的 `docker.sh` / `docker-compose` 包装脚本。

## 1. 构建镜像

```bash
docker build -t utbench:latest .
```

检查 CLI 是否正常：

```bash
docker run --rm utbench:latest --help
```

## 2. 准备目录和配置

仓库根目录下需要至少有这些内容：

- `datasets/`
- `artifacts/`
- `configs/models.yaml`

如果要调用真实模型，再准备一个 `.env` 文件，例如：

```dotenv
DEEPSEEK_API_KEY=sk-your-deepseek-key
DASHSCOPE_API_KEY=sk-your-qwen-key
MINIMAX_API_KEY=your-minimax-key
VOLCENGINE_API_KEY=your-volcengine-key
ARK_API_KEY=your-ark-key
```

## 3. 先跑一个 dry-run

Windows PowerShell:

```powershell
docker run --rm `
  -v "${PWD}/datasets:/app/datasets" `
  -v "${PWD}/artifacts:/app/artifacts" `
  -v "${PWD}/configs:/app/configs" `
  utbench:latest run `
    --models deepseek `
    --langs python `
    --dataset-root /app/datasets `
    --output-root /app/artifacts `
    --config /app/configs/models.yaml `
    --class self_contained `
    --scenario simple_function `
    --max-samples 1 `
    --run-id docker_dryrun_001 `
    --dry-run
```

Linux / macOS:

```bash
docker run --rm \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  utbench:latest run \
    --models deepseek \
    --langs python \
    --dataset-root /app/datasets \
    --output-root /app/artifacts \
    --config /app/configs/models.yaml \
    --class self_contained \
    --scenario simple_function \
    --max-samples 1 \
    --run-id docker_dryrun_001 \
    --dry-run
```

## 4. 跑真实模型

```bash
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  utbench:latest run \
    --models deepseek,qwen \
    --langs python,go \
    --dataset-root /app/datasets \
    --output-root /app/artifacts \
    --config /app/configs/models.yaml \
    --class self_contained \
    --max-samples 5 \
    --run-id docker_real_001
```

## 5. 分步执行

如果你不想直接跑 `run`，可以手动分阶段执行。注意要复用同一个 `--run-id`，否则输出会分散到不同目录。

```bash
RUN_ID=docker_step_001

docker run --rm \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  utbench:latest generate \
    --run-id "$RUN_ID" \
    --models deepseek \
    --langs python \
    --dataset-root /app/datasets \
    --output-root /app/artifacts \
    --config /app/configs/models.yaml \
    --max-samples 2 \
    --dry-run

docker run --rm \
  -v "$(pwd)/artifacts:/app/artifacts" \
  utbench:latest evaluate \
    --run-id "$RUN_ID" \
    --manifest /app/artifacts/runs/$RUN_ID/generated/generated_manifest.json \
    --output-root /app/artifacts

docker run --rm \
  -v "$(pwd)/artifacts:/app/artifacts" \
  utbench:latest report \
    --run-id "$RUN_ID" \
    --input /app/artifacts/runs/$RUN_ID/evaluation/evaluation_result.json \
    --output-root /app/artifacts
```

关键点：

- `evaluate` 用 `--manifest`
- `report` 用 `--input`
- 报告目录是 `/app/artifacts/runs/<run-id>/report/`

## 6. 使用仓库自带脚本

仓库里保留了两个可直接使用的包装脚本：

Windows PowerShell:

```powershell
.\run_bench.ps1 -Models deepseek -Langs python -MaxSamples 2 -Mutation $false
```

Linux / macOS / WSL:

```bash
chmod +x ./run_bench.sh
./run_bench.sh deepseek python 2 false
```

这两个脚本默认都使用：

- 镜像名 `utbench:latest`
- 配置文件 `/app/configs/models.yaml`
- 挂载 `datasets/`、`artifacts/`、`configs/`

## 7. 结果目录

```text
artifacts/runs/<run-id>/
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

## 8. 常见问题

### `docker run ... evaluate --input ...` 为什么报错？

因为当前 CLI 的 `evaluate` 子命令只接受 `--manifest`。

### 为什么文档里的模型名和我本地不一样？

以 [configs/models.yaml](configs/models.yaml) 为准，命令里的 `--models` 必须使用其中的键名。

### 为什么不建议依赖默认配置路径？

因为 CLI 默认路径是 `../benchmark/config/models.yaml`。在当前仓库根目录运行时，更稳妥的做法是显式传入 `--config /app/configs/models.yaml` 或 `--config ./configs/models.yaml`。
