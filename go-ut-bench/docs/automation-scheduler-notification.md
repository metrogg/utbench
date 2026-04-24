# 自动化评测与通知方案

## 1. 背景

当前 Web 端已经可以手动创建评测任务、查看任务状态、重跑任务，并通过模型管理配置连接信息。下一步可以增加自动化能力，让系统按计划定时触发评测任务，并在任务结束后通过消息通道网关发送结果通知，例如邮箱、企业微信、飞书、钉钉或 Webhook。

目标是把“人工点一次运行”扩展为“按规则自动运行 + 自动汇报结果”，适合长期监控模型能力变化、回归检查评测链路，以及周期性产出报告。

## 2. 核心需求

- 支持定时触发评测任务。
- 支持配置评测任务参数，包括模型、语言、样本类型、场景、样本数量、是否 Docker 运行、是否开启 mutation、是否入库。
- 支持任务执行完成后读取评测结果。
- 支持通过消息通道发送通知，第一阶段优先支持邮箱。
- 通知内容需要包含任务状态、成功率、失败摘要、报告链接或报告路径。
- 不在日志、报告、通知中泄露 API Key、SMTP 密码等敏感信息。

## 3. 方案 A：轻量内置调度器

### 设计

在 Web Server 进程内增加一个后台 goroutine，启动时读取 YAML 配置中的自动化任务列表，并使用 ticker 周期性检查是否有任务需要触发。

示例配置：

```yaml
automation:
  enabled: true
  schedules:
    - name: nightly-python-smoke
      enabled: true
      cron: "0 2 * * *"
      use_docker: true
      notify_channels: ["mail-default"]
      run:
        models: ["deepseek", "qwen"]
        langs: ["python"]
        class: "self_contained"
        scenarios: ["simple_function"]
        max_samples: 5
        mutation_enabled: false

notifications:
  channels:
    - name: mail-default
      type: smtp
      smtp:
        host: "smtp.example.com"
        port: 587
        username_env: "UTBENCH_SMTP_USER"
        password_env: "UTBENCH_SMTP_PASSWORD"
        from: "utbench@example.com"
        to: ["team@example.com"]
```

### 优点

- 实现成本最低。
- 不需要额外数据库表或复杂 UI。
- 适合先验证自动化评测和邮件通知链路。

### 缺点

- 任务配置主要依赖配置文件，Web 端编辑能力弱。
- 进程重启、并发保护、错过触发时间等场景需要额外补偿。
- 调度历史不够完整，不利于审计和排查。

### 适用场景

适合第一版快速落地，或者仅需要少量固定周期任务的本地/单机部署。

## 4. 方案 B：SQLite 持久化调度中心（推荐）

### 设计

把自动化任务作为一类持久化资源保存在 SQLite 中，由 Web Server 内置调度器负责扫描、触发和记录执行历史。Web UI 提供自动化任务管理页面，可以新增、编辑、暂停、手动触发和查看历史。

建议新增表：

```sql
CREATE TABLE automation_schedules (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  enabled INTEGER NOT NULL DEFAULT 1,
  cron TEXT NOT NULL,
  use_docker INTEGER NOT NULL DEFAULT 1,
  run_spec_json TEXT NOT NULL,
  notify_channel_ids_json TEXT NOT NULL,
  last_run_at TEXT,
  next_run_at TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE notification_channels (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  type TEXT NOT NULL,
  config_json TEXT NOT NULL,
  enabled INTEGER NOT NULL DEFAULT 1,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE automation_runs (
  id TEXT PRIMARY KEY,
  schedule_id TEXT NOT NULL,
  run_id TEXT NOT NULL,
  status TEXT NOT NULL,
  started_at TEXT,
  completed_at TEXT,
  summary_json TEXT,
  notify_status_json TEXT,
  created_at TEXT NOT NULL
);
```

### Web API

建议增加以下接口：

```text
GET    /api/automations
POST   /api/automations
PUT    /api/automations/{id}
DELETE /api/automations/{id}
POST   /api/automations/{id}/trigger
GET    /api/automations/{id}/runs

GET    /api/notification-channels
POST   /api/notification-channels
PUT    /api/notification-channels/{id}
DELETE /api/notification-channels/{id}
POST   /api/notification-channels/{id}/test
```

### 通知内容

邮件标题示例：

```text
[UTBench] nightly-python-smoke completed: 12/14 passed
```

邮件正文建议包含：

- 自动化任务名称。
- 本次运行 ID。
- 触发时间、完成时间、耗时。
- 模型、语言、场景、样本数量。
- 总样本数、通过数、失败数、编译失败数、测试失败数、超时数。
- Top N 失败原因摘要。
- `report_summary.json` 路径。
- `report.html` 路径或 Web 详情页链接。

### 优点

- 调度配置、历史、通知状态都有持久化记录。
- 可以自然接入现有 Web 管理后台。
- 支持暂停、立即触发、失败重试、通知测试等完整运维能力。
- 后续扩展飞书、企业微信、钉钉、Webhook 比较顺滑。

### 缺点

- 实现成本高于方案 A。
- 需要设计数据库表、迁移、后端接口和 Web 页面。
- 需要处理调度器并发、进程重启恢复、重复触发保护。

### 推荐原因

这个项目已经有 Web UI、SQLite、任务管理和报告产物，自动化任务本质上是“可持久化的运行模板 + 触发规则 + 通知策略”。把它放进 SQLite 和 Web 管理后台，长期维护成本最低，也最符合当前架构方向。

## 5. 方案 C：外部调度系统

### 设计

不在 go-ut-bench 内部实现调度器，而是把 CLI 或 Web API 暴露给外部系统调用。例如：

- Windows Task Scheduler
- Linux cron
- Jenkins
- GitHub Actions
- GitLab CI
- Airflow
- Kubernetes CronJob

外部任务定时执行：

```bash
./utbench run \
  --models deepseek,qwen \
  --langs python,go \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --config ./configs/models.yaml \
  --class self_contained \
  --scenario simple_function \
  --max-samples 10 \
  --mutation-enabled
```

通知也可以由外部系统完成，或者调用 go-ut-bench 后续提供的通知脚本/API。

### 优点

- go-ut-bench 内部改动最少。
- 适合已经有 CI/CD 或运维调度平台的团队。
- 调度可靠性、权限、日志、失败重试可以交给成熟平台。

### 缺点

- Web UI 无法统一管理自动化任务。
- 用户需要理解外部调度系统。
- 本地单机使用体验不如内置调度器。

### 适用场景

适合团队已经有 Jenkins、GitHub Actions、Kubernetes 等基础设施，并希望把 go-ut-bench 作为评测命令接入现有流水线。

## 6. 推荐实现路径

建议分两阶段推进：

### 阶段 1：先做方案 B 的最小闭环

- 新增 SQLite 表：`automation_schedules`、`notification_channels`、`automation_runs`。
- 后端实现自动化任务 CRUD。
- 后端实现通知通道 CRUD。
- 支持 SMTP 邮箱通知。
- 支持手动触发自动化任务。
- 支持定时触发自动化任务。
- 调度触发默认使用 Docker 运行。
- 任务完成后发送简版邮件通知。

### 阶段 2：增强可运维性

- Web UI 增加自动化任务页面。
- Web UI 增加通知通道测试按钮。
- 支持失败重试策略。
- 支持只在失败时通知、每次都通知、状态变化时通知。
- 支持 Webhook、飞书、企业微信、钉钉。
- 支持查看最近 N 次自动化运行记录。
- 支持从历史任务一键保存为自动化模板。

## 7. 第一版边界建议

第一版不要做太复杂，建议先限定：

- 只支持单机 Web Server 内置调度。
- cron 表达式可以先支持常见五段式，或者先支持 daily/weekly/interval 三种简单模式。
- 通知先只支持 SMTP 邮箱。
- 邮箱密码只允许从环境变量读取，不写入 YAML 或 SQLite 明文。
- 自动化任务触发时默认 Docker 运行。
- 同一 schedule 上一次还在运行时，默认跳过本次触发。

## 8. 风险点

- 进程重启时可能错过触发时间，需要定义是否补跑。
- 多个 Web Server 实例同时运行时可能重复触发，需要数据库锁或部署约束。
- 邮件发送失败不能影响评测结果落盘，需要单独记录通知状态。
- 通知内容必须避免泄露 API Key、SMTP 密码、模型返回中的敏感内容。
- 大任务执行时间可能跨过下一个周期，需要明确并发策略。

## 9. 结论

如果目标是快速验证，可以先做方案 A。但从当前项目已经具备 Web UI、SQLite 和任务历史的情况看，更推荐直接做方案 B 的最小闭环：SQLite 持久化调度中心 + Web 管理 + SMTP 邮件通知。这样既能满足定时评测和结果通知，也能为后续企业微信、飞书、钉钉、Webhook 等消息通道扩展留下清晰接口。
