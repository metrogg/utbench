# ut-bench 系统设计（v0.1）

> 状态：Draft（迭代式完善）
> 技术路线：Electron + React + TypeScript（桌面端优先）
> 对齐文档：`docs/design/requirements-spec.md`

## 1. 设计目标与原则

### 1.1 设计目标

- 做成可交互的软件，而非一次性脚本。
- 支持模型多选、数据集多选、提示词工程、可视化评测报告。
- v0.1 跑通 Python 端到端；其他语言保留扩展位。

### 1.2 核心原则

- 分层解耦：UI、应用服务、执行引擎、适配器分离。
- 配置驱动：模型与任务配置优先通过配置表达。
- 可追溯：任何指标都能追溯到样本级原始数据与日志。
- 可扩展：新增模型/语言/图表尽量不改核心流程。

## 2. 整体架构

## 2.1 逻辑分层

1. UI 层（Renderer）
   - 页面交互、参数选择、任务控制、报告展示。
2. 应用服务层（Main Process Service）
   - IPC 路由、权限边界、任务生命周期管理。
3. 评测执行引擎（Core Engine）
   - 任务编排、模型调用、评测计算、结果聚合。
4. 外部工具适配层（Adapter）
   - pytest/coverage/mutation/docker 等调用封装。
5. 存储层（Storage）
   - 配置存储（SQLite）+ 产物存储（文件系统 results/）。

### 2.2 端到端执行链路

1. 用户在 UI 选择模型、数据集、提示词模板。
2. UI 发起 `createRun` + `startRun`。
3. Orchestrator 生成任务矩阵（模型 x 样本）。
4. Prompt Engine 生成最终提示词。
5. Model Adapter 调用 LLM 生成测试代码。
6. Evaluator 计算 G1/G2/R1/C1/C2/C3 指标。
7. Reporter 产出 JSON 与 HTML（矩阵、雷达图、排名）。
8. UI 实时展示进度并支持结果回放。

## 3. 模块划分

### 3.1 Model Manager

职责：模型配置管理与连接校验。

- 新增/编辑/启用/禁用模型。
- 支持多 Provider（OpenAI compatible 优先）。
- 密钥通过环境变量引用，不明文入库。

### 3.2 Dataset Manager

职责：数据集发现、过滤、样本清单构建。

- 支持按语言、难度、场景多选。
- 生成标准样本描述（sample manifest）。

### 3.3 Prompt Studio

职责：提示词模板、变量替换、版本管理。

- 支持 system/user 模板分离。
- 支持变量：语言、风格、约束、样本上下文。
- 支持 prompt 预览与版本回溯。

### 3.4 Run Orchestrator

职责：任务调度与生命周期管理。

- 创建 Run、启动、暂停、取消、重试。
- 控制并发度、超时、失败隔离。
- 汇总进度事件供 UI 展示。

### 3.5 Evaluator

职责：执行评测并输出标准指标。

- G1 编译通过率。
- G2 运行通过率。
- R1 分支覆盖率。
- C1 变异测试得分。
- C2 Token 消耗。
- C3 生成耗时。

### 3.6 Report Engine

职责：聚合数据并生成可视化结果。

- 指标矩阵（模型 x 指标）。
- 雷达图（能力轮廓）。
- 排名视图（默认按主指标排序，暂不固定加权总分）。

### 3.7 Artifact Store

职责：结果、日志、快照存储与检索。

- Run 元数据、样本明细、错误日志。
- 配置快照（模型、提示词模板版本、运行参数）。

## 4. 技术选型

### 4.1 桌面端

- 框架：Electron。
- UI：React + TypeScript。
- 构建：Vite。
- 图表：ECharts（矩阵热力图/雷达图/排名图实现友好）。

### 4.2 执行与存储

- 运行时：Node.js (TypeScript)。
- 数据库存储：SQLite（轻量、本地部署友好）。
- 产物存储：文件系统 `results/`。
- 外部工具：通过命令适配层封装调用（pytest、coverage、mutmut、docker）。

### 4.3 选型理由（简要）

- Electron 生态成熟，工程交付速度快。
- TS 统一前后逻辑，降低心智切换成本。
- SQLite 易于本地安装与导出调试。

## 5. 详细设计

### 5.1 任务编排逻辑

- 输入：模型集合 M、样本集合 S。
- 任务矩阵：`Tasks = M x S`。
- 每个 Task 独立执行，失败互不影响。
- 支持按模型维度或样本维度并发（默认按模型并发）。

### 5.2 状态机

Run 状态：

- `created` -> `running` -> `completed`
- `created/running` -> `cancelled`
- `running` -> `failed`（系统级故障）

Task 状态：

- `pending` -> `generating` -> `evaluating` -> `done`
- 任一阶段失败 -> `error`（记录错误类型和上下文）

### 5.3 失败处理策略

- 模型调用失败：按策略重试（指数退避，最多 N 次）。
- 评测超时：标记当前 Task 失败，不中断整批 Run。
- 外部工具缺失：Run 前预检并给出可操作提示。

## 6. 数据结构设计

### 6.1 核心实体

- `ModelConfig`
  - `id`, `name`, `provider`, `model`, `baseUrl`, `apiKeyEnv`, `enabled`
- `DatasetSample`
  - `id`, `language`, `path`, `tags`, `metadata`
- `PromptTemplate`
  - `id`, `name`, `version`, `systemTemplate`, `userTemplate`, `variables`
- `Run`
  - `id`, `status`, `createdAt`, `startedAt`, `endedAt`, `configSnapshot`
- `TaskResult`
  - `runId`, `modelId`, `sampleId`, `status`, `generatedTestPath`, `error`
- `MetricResult`
  - `compilePass`, `runPass`, `branchCoverage`, `mutationScore`, `tokenUsage`, `latencyMs`

### 6.2 存储布局建议

`results/<run_id>/`

- `raw.json`（样本级完整结果）
- `summary.json`（聚合结果）
- `report.html`（可视化报告）
- `logs/*.log`（结构化执行日志）
- `snapshots/*.json`（配置快照）

## 7. API / IPC 设计

桌面端采用 IPC（Renderer <-> Main）方式，建议接口如下：

### 7.1 Run 相关

- `run.create(payload)` -> `{ runId }`
- `run.start(runId)` -> `{ accepted: true }`
- `run.cancel(runId)` -> `{ accepted: true }`
- `run.get(runId)` -> `RunDetail`
- `run.list(query)` -> `RunSummary[]`

### 7.2 配置相关

- `model.list()` / `model.create()` / `model.update()` / `model.delete()`
- `model.testConnection(modelId)`
- `dataset.scan(filters)`
- `prompt.list()` / `prompt.saveVersion()` / `prompt.preview()`

### 7.3 报告相关

- `report.get(runId)` -> 报告 JSON 结构
- `report.exportHtml(runId, outPath)`

## 8. 算法与排名设计

- v0.1 不做固定综合加权总分。
- 默认排序策略：
  1. 先按 G1/G2 通过率过滤与分组。
  2. 再按主指标排序（默认 C1 变异测试得分）。
  3. 同分时按 R1、C3、C2 作为次级排序键。
- 排序策略可配置，但需记录到 run 配置快照。

## 9. 扩展性设计

### 9.1 模型扩展

- 统一 `ModelAdapter` 接口：`generateTest(input) -> output`。
- 新 Provider 仅实现适配器，不改 Orchestrator。

### 9.2 语言扩展

- 统一 `LanguageEvaluator` 接口：`compile/run/coverage/mutation`。
- v0.1 实现 Python；Java/Go/C++ 提供接口占位。

### 9.3 图表扩展

- 报告数据与渲染层分离。
- 新图表通过注册机制接入，不影响核心计算逻辑。

## 10. 可靠性设计

- 执行前预检：工具链、目录权限、模型配置完整性。
- 超时机制：模型调用、测试执行、变异测试分别配置。
- 幂等与可恢复：Run 失败后可局部重试未完成 Task。
- 数据一致性：先写 raw，再写 summary，再生成 report。

## 11. 监控与运维设计

- 日志规范：JSON 行日志，包含 `runId/modelId/sampleId/stage`。
- 运行事件：开始、进度、失败、完成均可订阅。
- 诊断包导出：自动收集配置快照、日志、错误摘要。
- 本地运维：提供“环境检查”按钮（工具可用性检测）。

## 12. 文档规范与协作

- 使用统一术语：Run、Task、Sample、Metric、Artifact。
- 关键决策采用 ADR 记录（如 Electron vs Tauri）。
- 需求变更需回写 `requirements-spec.md` 与本设计文档。
- 未实现部分必须标注 TODO，不得伪造实现状态。

## 13. 迭代计划（设计演进）

### Iteration 1（本周）

- Python 单语言评测闭环。
- 模型多选、数据集多选。
- HTML 报告（矩阵 + 雷达图 + 排名）。

### Iteration 2

- 强化提示词工程（模板市场、A/B 对比）。
- 引入更多报告维度（趋势、失败画像、成本分析）。

### Iteration 3

- Java/Go/C++ 实现落地。
- 可靠性增强（断点续跑、增量执行、任务缓存）。

## 14. 当前待决策项

- Electron 打包与发布策略（内部安装包/自动更新方案）。
- 数据库存储迁移策略（本地 SQLite -> 可选远端）。
- 报告导出格式（HTML 必选，PDF/图片是否纳入 v0.1）。

## 15. 模块实现细则（v0.1 冻结草案）

### 15.1 Model Manager

输入：用户在 UI 填写模型配置。

处理逻辑：

1. 字段校验（`name/provider/model/baseUrl/apiKeyEnv` 必填）。
2. 配置落库（密钥只保存环境变量名，不保存明文）。
3. `testConnection` 发送最小请求（短 prompt + 5s 超时）。
4. 记录诊断（成功/失败原因、耗时）。

输出：模型可用状态与错误信息。

默认参数：

- 连接超时：5s
- 请求超时：30s
- 重试次数：2（仅网络错误和 5xx）

### 15.2 Dataset Manager

输入：数据集目录与筛选条件（语言、标签、难度）。

处理逻辑：

1. 扫描目录生成样本清单。
2. 为每个样本计算 `checksum`（SHA-256）。
3. 验证样本结构（源文件存在、扩展名合法）。
4. 产出 `manifest` 并缓存到本地数据库。

输出：可执行样本列表与无效样本报告。

默认参数：

- 扫描缓存 TTL：300s
- 允许最大样本数（单次 Run）：100（可配置）

### 15.3 Prompt Studio

输入：模板版本 + 变量字典。

处理逻辑：

1. 加载模板（system/user）。
2. 校验变量完整性（缺失变量直接报错）。
3. 渲染最终 prompt（Mustache）。
4. 存储 prompt 快照用于回溯。

输出：最终 prompt 文本与 `prompt_version_id`。

默认参数：

- 模板引擎：Mustache
- 模板渲染长度上限：64KB

### 15.4 Run Orchestrator

输入：`models[]`、`samples[]`、`promptTemplate`、运行参数。

处理逻辑：

1. 创建 Run 与配置快照。
2. 展开任务矩阵（`models x samples`）。
3. 按并发池调度任务（默认并发 2）。
4. 任务执行顺序：生成 -> 评测 -> 落盘。
5. 按事件流推送 UI 进度。

输出：Run 状态、Task 状态、聚合统计。

默认参数：

- 并发度：2
- 单 Task 总超时：600s
- 可重试阶段：生成阶段

### 15.5 Generator Adapter

输入：模型配置 + prompt + 样本上下文。

处理逻辑：

1. 组装统一请求体。
2. 调用 Provider 接口。
3. 标准化响应（代码、usage、耗时、raw）。
4. 返回给 Evaluator。

输出：`GenerateResult`。

错误分级：

- `E_MODEL_AUTH`（鉴权失败）
- `E_MODEL_RATE_LIMIT`（限流）
- `E_MODEL_TIMEOUT`（超时）
- `E_MODEL_RESPONSE`（响应结构异常）

### 15.6 Python Evaluator

输入：源代码文件 + 生成测试代码 + 运行配置。

处理逻辑：

1. 生成临时工作目录，写入测试文件。
2. 执行 G1：`python -m py_compile`。
3. 执行 G2：`pytest -q`。
4. 执行 R1：`pytest --cov --cov-branch --cov-report=json`。
5. 执行 C1：`mutmut run`（建议容器执行）。
6. 合并 C2/C3（来自生成阶段 usage/latency）。
7. 产出样本级 `MetricResult`。

输出：样本级指标与错误摘要。

默认参数：

- `pytest` 超时：120s
- `mutmut` 超时：600s
- 失败保留：stdout/stderr 各前 200 行

### 15.7 Report Engine

输入：`raw.json` + `summary.json`。

处理逻辑：

1. 构建指标矩阵（模型 x 指标）。
2. 构建雷达图数据（模型能力轮廓）。
3. 构建排名列表（按策略排序）。
4. 渲染 HTML 模板并内嵌图表数据。

输出：`report.html`。

## 16. SQLite 表结构草案（DDL）

```sql
CREATE TABLE IF NOT EXISTS model_configs (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  provider TEXT NOT NULL,
  model TEXT NOT NULL,
  base_url TEXT NOT NULL,
  api_key_env TEXT NOT NULL,
  enabled INTEGER NOT NULL DEFAULT 1,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS prompt_templates (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  version INTEGER NOT NULL,
  system_template TEXT NOT NULL,
  user_template TEXT NOT NULL,
  variables_json TEXT NOT NULL,
  created_at TEXT NOT NULL,
  UNIQUE(name, version)
);

CREATE TABLE IF NOT EXISTS runs (
  id TEXT PRIMARY KEY,
  status TEXT NOT NULL,
  language TEXT NOT NULL,
  config_snapshot_json TEXT NOT NULL,
  created_at TEXT NOT NULL,
  started_at TEXT,
  ended_at TEXT
);

CREATE TABLE IF NOT EXISTS run_tasks (
  id TEXT PRIMARY KEY,
  run_id TEXT NOT NULL,
  model_id TEXT NOT NULL,
  sample_id TEXT NOT NULL,
  status TEXT NOT NULL,
  retry_count INTEGER NOT NULL DEFAULT 0,
  generated_test_path TEXT,
  error_code TEXT,
  error_message TEXT,
  started_at TEXT,
  ended_at TEXT,
  FOREIGN KEY(run_id) REFERENCES runs(id)
);

CREATE TABLE IF NOT EXISTS metric_results (
  id TEXT PRIMARY KEY,
  run_id TEXT NOT NULL,
  task_id TEXT NOT NULL,
  compile_pass INTEGER,
  run_pass INTEGER,
  branch_coverage REAL,
  mutation_score REAL,
  prompt_tokens INTEGER,
  completion_tokens INTEGER,
  latency_ms INTEGER,
  raw_json TEXT NOT NULL,
  created_at TEXT NOT NULL,
  FOREIGN KEY(run_id) REFERENCES runs(id),
  FOREIGN KEY(task_id) REFERENCES run_tasks(id)
);
```

## 17. IPC 接口契约（TypeScript）

```ts
type ApiResult<T> =
  | { ok: true; data: T }
  | { ok: false; error: { code: string; message: string; detail?: unknown } }

interface RunCreatePayload {
  language: "python"
  modelIds: string[]
  datasetFilters: Record<string, unknown>
  promptTemplateId: string
  concurrency?: number
}

interface GenerateResult {
  code: string
  promptTokens: number
  completionTokens: number
  latencyMs: number
  rawResponse: unknown
}
```

约束：

- 所有 IPC 返回统一 `ApiResult<T>`。
- 错误码稳定可枚举，禁止仅返回自由文本。
- 每个接口必须携带 `runId` 或实体 `id` 用于追踪。

## 18. 排名与聚合算法细则

### 18.1 聚合

- G1/G2：按样本通过率计算百分比。
- R1/C1：按样本平均值与中位数双输出。
- C2/C3：输出均值、P50、P95。

### 18.2 排名

不使用固定加权总分，采用字典序排序：

1. G1 是否达阈值（高优先）
2. G2 是否达阈值（高优先）
3. C1（高优先）
4. R1（高优先）
5. C3（低优先）
6. C2（低优先）

同分处理：

- 若所有排序键相同，按模型名升序稳定排序。

## 19. 可靠性与运维细则

### 19.1 预检项

- Python 版本可用。
- `pytest`、`pytest-cov` 可用。
- `mutmut` 可用（或容器可用）。
- `results/` 目录可写。

### 19.2 日志规范

统一字段：

- `ts`, `level`, `runId`, `taskId`, `modelId`, `sampleId`, `stage`, `message`

日志级别：

- `INFO`：状态变化、关键里程碑
- `WARN`：可恢复异常
- `ERROR`：任务失败或系统故障

### 19.3 诊断包

导出内容：

- `summary.json`
- `logs/*.log`
- `snapshots/*.json`
- 工具链检查结果

## 20. 文档与变更管理

- 每次实现变更需同步更新本文件对应章节。
- 关键取舍以 ADR 文档记录（`docs/design/adr/`，后续创建）。
- 需求变更优先修改 `requirements-spec.md`，再更新系统设计。
