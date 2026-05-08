# go-ut-bench 快速开始

这份文档只写当前仓库里已经能跑通的路径，并补上 Agent + Skill 的最小用法。

## 1. 构建镜像

在 [go-ut-bench](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench) 目录执行：

```bash
docker build -t utbench:latest .
```

检查 CLI：

```bash
docker run --rm utbench:latest --help
```

## 2. 准备配置

至少要有：

- `datasets/`
- `artifacts/`
- `configs/models.yaml`

如果要调真实模型，再准备 `.env`：

```dotenv
DEEPSEEK_API_KEY=...
DASHSCOPE_API_KEY=...
MINIMAX_API_KEY=...
VOLCENGINE_API_KEY=...
ARK_API_KEY=...
BIGMODEL_API_KEY=...
```

如果要跑 CLI Agent，再准备一个 Agent 配置文件。仓库里已有示例：

- [configs/agents.example.yaml](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/configs/agents.example.yaml)

如果你要跑当前示例里的 `OpenCode` framework，还需要先构建内层 Agent 镜像：

```bash
docker build -t utbench-agent-opencode:latest -f ./docker/agents/opencode/Dockerfile .
```

## 3. 先跑纯模型 baseline

Windows PowerShell：

```powershell
docker run --rm --env-file .env `
  -v "${PWD}/datasets:/app/datasets" `
  -v "${PWD}/artifacts:/app/artifacts" `
  -v "${PWD}/configs:/app/configs" `
  utbench:latest run `
    --models deepseek-v4-flash `
    --langs python `
    --dataset-root /app/datasets `
    --output-root /app/artifacts `
    --config /app/configs/models.yaml `
    --class self_contained `
    --max-samples 1 `
    --run-id quick_model_001
```

Linux / macOS：

```bash
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  utbench:latest run \
    --models deepseek-v4-flash \
    --langs python \
    --dataset-root /app/datasets \
    --output-root /app/artifacts \
    --config /app/configs/models.yaml \
    --class self_contained \
    --max-samples 1 \
    --run-id quick_model_001
```

这会自动包含纯模型 subject：

```text
model_api__deepseek-v4-flash__no_skill
```

## 4. 跑 Agent + Skill

当前 CLI Agent 的入口不是写死在代码里，而是从 `agents.yaml` 读：

- `frameworks.<name>.command`
- `frameworks.<name>.sandbox_mode`
- `frameworks.<name>.docker_image`
- `frameworks.<name>.env`
- `frameworks.<name>.env_from_host`
- `skills.<name>.*`

最小命令：

```bash
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  utbench:latest run \
    --models deepseek-v4-flash \
    --langs python \
    --config /app/configs/models.yaml \
    --agents-config /app/configs/agents.example.yaml \
    --subjects model_api__deepseek-v4-flash__no_skill,opencode__deepseek-v4-flash__no_skill,opencode__deepseek-v4-flash__unit_test_skill \
    --dataset-root /app/datasets \
    --output-root /app/artifacts \
    --class self_contained \
    --max-samples 1 \
    --run-id quick_agent_001
```

如果不传 `--subjects`，会自动展开 `framework × model × skill` 的全部合法组合。

## 5. CLI Agent 实际怎么被调用

当前实现是通用模板调用，不是为某一个 Agent 写死参数。

`cli_agent` 的实际执行流程：

1. 为每个 `subject × sample` 复制一个独立工作区
2. 生成 `utbench_agent_prompt.md`
3. 注入 skill 文件或 skill prompt
4. 渲染 `framework.command`
5. 执行命令
6. 查找测试文件输出
7. 记录 trace 和 workspace diff

当 `sandbox_mode: docker` 时，UT-Bench 会实际执行一条内层 Docker 命令，逻辑等价于：

```text
docker run --rm \
  [--network none] \
  [--cpus ...] \
  [--memory ...] \
  -v <workspace>:/workspace \
  -w /workspace \
  <docker_image> \
  /bin/sh -lc "<rendered command>"
```

所以，CLI Agent 不是共享整个 benchmark 运行目录，而是只看到当前样本的那个工作区。

## 6. 只跑 generate / evaluate / report

```bash
RUN_ID=quick_step_001

docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  utbench:latest generate \
    --run-id "$RUN_ID" \
    --models deepseek-v4-flash \
    --langs python \
    --config /app/configs/models.yaml \
    --agents-config /app/configs/agents.example.yaml \
    --dataset-root /app/datasets \
    --output-root /app/artifacts \
    --class self_contained \
    --max-samples 1

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
    --evaluation /app/artifacts/runs/$RUN_ID/evaluation/evaluation_result.json \
    --output-root /app/artifacts
```

注意：

- `evaluate` 用 `--manifest`
- `report` 用 `--evaluation`

## 7. 结果目录

```text
artifacts/runs/<run-id>/
  generated/
    generated_manifest.json
    tests/
    prompts/
    metadata/
  evaluation/
    evaluation_result.json
  report/
    report_summary.json
    report.html
  agent_workspaces/
  logs/
  run_summary.json
```

## 8. 对 CLI Agent 的选择建议

如果现在就要接一个真实 Agent，建议先接 `OpenCode`。

原因：

- 更贴合 `framework × model` 的组合目标
- 更适合作为“Agent 壳 + 多模型 provider”入口
- 自动化 benchmark 场景里更容易做非交互式调用

`Claude Code` 适合第二个接入，用来评测 Claude 自身工作流，而不是作为第一优先的多模型框架入口。
