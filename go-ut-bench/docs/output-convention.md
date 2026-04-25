<<<<<<< HEAD
# 输出目录约定

## 根目录

- 运行产物：`./artifacts`
- SQLite：`./storage`

## 标准结构
=======
# 输出目录规范

## 1. 目标

统一产物落盘路径，避免 `results_real*` 这类临时目录散落。

## 2. 目录约定

- 运行产物根目录：`./artifacts`
- 数据库存储目录：`./storage`

完整结构：
>>>>>>> origin/feat/go

```text
artifacts/
  checkpoints/
    runner_<scope_hash>.checkpoint.json
  runs/
<<<<<<< HEAD
    <run-id>/
      generated/
        generated_manifest.json
        tests/
          <model>/
            <lang>/
              *.test.<ext>
        metadata/
          <model>_<lang>_<sample>.response.json
          <model>_<lang>_<sample>.metadata.json
=======
    <run_id>/
      generated/
        generated_manifest.json
        tests/
        metadata/
>>>>>>> origin/feat/go
      evaluation/
        evaluation_result.json
      report/
        report_summary.json
        report.html
      run_summary.json

storage/
  utbench.db
```

<<<<<<< HEAD
## 当前默认值

- `--output-root` 默认 `./artifacts`
- `--db-path` 默认 `./storage/utbench.db`

## 重要说明

- 报告目录是 `report/`，不是 `reports/`
- `run_summary.json` 只会在 `utbench run` 里自动生成
- 如果分阶段执行且不复用同一个 `--run-id`，`generate`、`evaluate`、`report` 会各自生成新的 `runs/<run-id>/` 目录
- `evaluate` 的输入文件是 `generated/generated_manifest.json`
- `report` 和 `ingest` 的输入文件是 `evaluation/evaluation_result.json`
=======
## 3. 命令约束

- 默认 `--output-root` 为 `./artifacts`
- 默认 `--db` 为 `./storage/utbench.db`
- 调试临时目录建议用 `--output-root ./artifacts_debug`，不要创建多个 `results_*`
>>>>>>> origin/feat/go
