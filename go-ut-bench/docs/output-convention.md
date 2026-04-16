# 输出目录规范

## 1. 目标

统一产物落盘路径，避免 `results_real*` 这类临时目录散落。

## 2. 目录约定

- 运行产物根目录：`./artifacts`
- 数据库存储目录：`./storage`

完整结构：

```text
artifacts/
  checkpoints/
    runner_<scope_hash>.checkpoint.json
  runs/
    <run_id>/
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

storage/
  utbench.db
```

## 3. 命令约束

- 默认 `--output-root` 为 `./artifacts`
- 默认 `--db` 为 `./storage/utbench.db`
- 调试临时目录建议用 `--output-root ./artifacts_debug`，不要创建多个 `results_*`
