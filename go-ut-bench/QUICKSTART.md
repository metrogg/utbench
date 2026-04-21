# go-ut-bench 快速使用指南

## 前置要求

- Docker Desktop 已安装并运行
- 约 3GB 磁盘空间（镜像大小）
- API 密钥（DeepSeek/Qwen/Minimax 等）

## 步骤一：获取镜像

### 方式 A：从 tar 文件加载（推荐）

```bash
# Windows PowerShell
docker load -i utbench.tar

# Linux/Mac
docker load -i utbench.tar
```

### 方式 B：自己构建

```bash
cd go-ut-bench
docker build -t utbench:latest .
```

## 步骤二：配置 API 密钥

编辑 `.env` 文件：

```bash
# API Keys
DEEPSEEK_API_KEY=sk-your-deepseek-key
DASHSCOPE_API_KEY=sk-your-qwen-key      # 通义千问
MINIMAX_API_KEY=your-minimax-key
VOLCENGINE_API_KEY=your-doubao-key       # 豆包
```

## 步骤三：运行测试

### Dry-run（无需密钥，验证环境）

```powershell
docker run --rm `
  -v "${PWD}/datasets:/app/datasets" `
  -v "${PWD}/artifacts:/app/artifacts" `
  -v "${PWD}/../benchmark/config:/app/config" `
  utbench:latest run --models deepseek --langs python --dataset-root /app/datasets --output-root /app/artifacts --config /app/config/models.yaml --max-samples 2 --dry-run
```

### 实际运行（需要密钥）

```powershell
docker run --rm --env-file .env `
  -v "${PWD}/datasets:/app/datasets" `
  -v "${PWD}/artifacts:/app/artifacts" `
  -v "${PWD}/../benchmark/config:/app/config" `
  utbench:latest run --models deepseek,qwen --langs python,go --dataset-root /app/datasets --output-root /app/artifacts --config /app/config/models.yaml --max-samples 10
```

## 步骤四：查看结果

结果在 `artifacts/runs/<run-id>/` 目录下：

```
artifacts/runs/<run-id>/
  generated/           # 生成的测试文件
  evaluation/          # 评测结果 JSON
  report/              # HTML 报告
```

## 常用命令

```bash
# 查看帮助
docker run --rm utbench:latest --help

# 只生成测试（不评测）
docker run --rm --env-file .env -v ... utbench:latest generate --models deepseek --langs python --max-samples 5

# 只评测（需要已有生成的测试）
docker run --rm -v ... utbench:latest evaluate --input /app/artifacts/runs/<run-id>/generated/generated_manifest.json

# 生成报告
docker run --rm -v ... utbench:latest report --input /app/artifacts/runs/<run-id>/evaluation/evaluation_result.json
```

## Windows 简化脚本

创建 `run.ps1`：

```powershell
$env:PWD = Get-Location

docker run --rm --env-file .env `
  -v "$env:PWD/datasets:/app/datasets" `
  -v "$env:PWD/artifacts:/app/artifacts" `
  -v "$env:PWD/../benchmark/config:/app/config" `
  utbench:latest run --models deepseek --langs python --dataset-root /app/datasets --output-root /app/artifacts --config /app/config/models.yaml --max-samples 10
```

运行：

```powershell
.\run.ps1
```

## Linux/Mac 简化脚本

创建 `run.sh`：

```bash
#!/bash/bin
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/../benchmark/config:/app/config" \
  utbench:latest run --models deepseek --langs python --dataset-root /app/datasets --output-root /app/artifacts --config /app/config/models.yaml --max-samples 10
```

运行：

```bash
chmod +x run.sh
./run.sh
```

## 导出镜像给他人

```bash
docker save utbench:latest -o utbench.tar
```

然后将以下文件打包给别人：
- `utbench.tar`（镜像文件）
- `datasets/`（数据集目录）
- `benchmark/config/models.yaml`（模型配置）
- `.env`（API密钥模板，不含真实密钥）
- `QUICKSTART.md`（本指南）