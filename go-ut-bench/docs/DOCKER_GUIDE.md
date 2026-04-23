# Docker 使用说明

当前仓库的 Docker 工作流基于 [Dockerfile](../Dockerfile)，不再依赖额外的 `docker.sh` 或 `docker-compose.yml`。

## 镜像包含什么

Docker 镜像会构建 `utbench` CLI，并预装常用评测工具链，包括：

- Python: `pytest`, `coverage`, `mutmut`
- Go: `gremlins`
- Java: JDK 21, Maven
- C++: clang / llvm / gcov / GoogleTest / mull

## 1. 构建镜像

```bash
docker build -t utbench:latest .
```

验证镜像入口：

```bash
docker run --rm utbench:latest --help
```

## 2. 准备挂载目录

建议在仓库根目录运行以下命令，并挂载这些路径：

| 主机目录 | 容器目录 | 用途 |
|----------|----------|------|
| `./datasets` | `/app/datasets` | 数据集 |
| `./artifacts` | `/app/artifacts` | 运行产物 |
| `./configs` | `/app/configs` | 模型配置 |
| `./storage` | `/app/storage` | SQLite 数据库 |

## 3. 准备环境变量

如果要调用真实模型，创建 `.env`：

```dotenv
DEEPSEEK_API_KEY=...
DASHSCOPE_API_KEY=...
MINIMAX_API_KEY=...
VOLCENGINE_API_KEY=...
ARK_API_KEY=...
```

使用哪些变量，取决于你在 `configs/models.yaml` 里启用了哪些模型。

## 4. dry-run 示例

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
    --run-id docker_doc_001 \
    --dry-run
```

## 5. 真实运行示例

```bash
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  -v "$(pwd)/storage:/app/storage" \
  utbench:latest run \
    --models deepseek,qwen \
    --langs python,go \
    --dataset-root /app/datasets \
    --output-root /app/artifacts \
    --config /app/configs/models.yaml \
    --max-samples 5 \
    --run-id docker_real_001 \
    --ingest \
    --db-path /app/storage/utbench.db
```

## 6. 分阶段执行

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

docker run --rm \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/storage:/app/storage" \
  utbench:latest ingest \
    --input /app/artifacts/runs/$RUN_ID/evaluation/evaluation_result.json \
    --db-path /app/storage/utbench.db
```

要点：

- `evaluate` 使用 `--manifest`
- `report` 使用 `--input`
- `ingest` 使用 `--db-path`
- 如果不手动指定同一个 `--run-id`，每一步都会落到不同目录

## 7. 使用仓库现成脚本

PowerShell:

```powershell
.\run_bench.ps1 -Models deepseek -Langs python -MaxSamples 2 -Mutation $false
```

Bash / WSL:

```bash
chmod +x ./run_bench.sh
./run_bench.sh deepseek python 2 false
```

## 8. 输出位置

```text
artifacts/runs/<run-id>/
  generated/
    generated_manifest.json
  evaluation/
    evaluation_result.json
  report/
    report_summary.json
    report.html
  run_summary.json
```

## 9. 常见问题

### 为什么我以前的 `evaluate --input ...generated_manifest.json` 现在不行？

因为当前 CLI 的 `evaluate` 子命令只接受 `--manifest`。

### 为什么镜像里命令是 `utbench:latest`？

因为文档和仓库自带脚本现在统一使用这个镜像标签。

### 为什么本地说明写的是 `./configs/models.yaml`，CLI 默认却不是这里？

CLI 默认仍是 `../benchmark/config/models.yaml`，这是代码里的默认值；但在当前仓库根目录运行时，显式传 `--config ./configs/models.yaml` 更可靠。
