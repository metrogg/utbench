# Reporter 设计文档

## 1. 目标

`reporter` 是 `ut-bench` 的结果展示层，负责把 `evaluator` 产出的逐样本评测结果，整理成适合对比分析和汇报的聚合报告。

它解决的核心问题有三个：

1. 把原始逐样本结果转成稳定的聚合指标。
2. 把 runner 和 evaluator 分散在不同文件里的信息补齐到同一份视图里。
3. 生成可直接消费的报告产物，包括 JSON、CSV、HTML。

## 2. 设计目标与非目标

### 2.1 设计目标

- 输入以 `evaluator_summary_*.json` 为主，兼容从 `results/` 自动发现最新文件。
- 对每条样本结果做标准化补全，得到统一的 `sample_rows`。
- 聚合维度固定且可扩展，至少支持：
  - 按模型
  - 按语言
  - 按复杂度
  - 按场景
  - 按模型 x 语言
- 保持指标口径稳定，缺失值保留为 `None`，不伪造 0。
- 支持导出多种产物：
  - 机器可消费的 JSON
  - 表格分析友好的 CSV
  - 面向汇报和人工阅读的 HTML

### 2.2 非目标

- 不重新执行评测。
- 不修改 runner/evaluator 原始结果。
- 不承担长时存储、数据库建模或在线服务职责。

## 3. 输入与输出

### 3.1 输入

`reporter` 依赖三类输入源：

1. `evaluator_summary_*.json`
   - 主数据源
   - 提供 `results`、`summary`、`filters`
2. `runner_summary_*.json`
   - 作为生成耗时、token 等信息的兜底来源
3. `results/<model>/reports/*.metadata.json`
   - 单样本级别的生成元数据，优先级高于 `runner_summary`

### 3.2 输出

默认输出到 `results/reports/`：

- `reporter_summary_<run_id>.json`
- `reporter_by_model_<run_id>.csv`
- `reporter_by_language_<run_id>.csv`
- `reporter_samples_<run_id>.csv`
- `reporter_report_<run_id>.html`

## 4. 总体架构

当前推荐的逻辑分层如下：

```text
CLI(__main__.py)
  -> Reporter.run()
    -> 解析输入与参数
    -> 读取 evaluator summary
    -> 加载 runner metadata 兜底信息
    -> enrich sample rows
    -> 聚合 metrics
    -> 构建 failure breakdown
    -> 组装 report payload
    -> 按 format 导出 JSON / CSV / HTML
```

建议的职责边界：

- `__main__.py`
  - 只负责参数解析和调用
- `pipeline.py`
  - 只负责编排，不承载大段业务细节
- `loader.py`
  - 读取 JSON
  - enrich 单样本行
  - 解析 sample_id、metadata、runner fallback
- `aggregator.py`
  - 计算总体指标
  - 计算分组指标
- `diagnostics.py`
  - 失败样本提取
  - 错误分类和示例摘要
- `exporters/`
  - 各格式导出器
- `contracts.py`
  - 常量、字段名、展示标签
- `common.py`
  - 数值转换、格式化、通用小函数

## 5. 核心数据流

### 5.1 原始输入行

`evaluator` 的 `results` 是 reporter 的主输入，每条记录至少包含：

- `model`
- `language`
- `sample_id`
- `generated_test_path`
- `source_path`
- `compile_pass`
- `test_pass`
- `line_coverage`
- `branch_coverage`
- `function_coverage`
- `mutation_score`
- `runtime_ms`
- `compile_error`
- `test_error`
- `coverage_error`
- `mutation_error`
- `sample_bucket`
- `sample_bucket_reason`
- mutation 统计字段

### 5.2 标准化后的 sample row

reporter 会把每条 evaluator 结果补全成统一结构，核心新增字段包括：

- `complexity`
- `scenario`
- `generation_latency_ms`
- `prompt_tokens`
- `completion_tokens`
- `total_tokens`
- `generation_metrics_source`
- `generation_metrics_reason`

补全优先级：

1. `results/<model>/reports/*.metadata.json`
2. `runner_summary_*.json`
3. 缺失，保留 `None`

这样设计的原因是：

- metadata 粒度最细，可信度最高
- runner summary 能兼容历史数据
- 缺失时必须显式保留，不污染均值和占比

## 6. 指标设计

### 6.1 总体指标

总体指标来自 `build_metrics(rows)`，建议保持以下口径稳定：

- `total_samples`
- `compile_pass_count`
- `compile_pass_rate = compile_pass_count / total_samples`
- `test_pass_count`
- `test_pass_rate = test_pass_count / compile_pass_count`
- `stage2_executable_samples`
- `avg_line_coverage`
- `avg_branch_coverage`
- `avg_function_coverage`
- `avg_mutation_score`
- `avg_eval_runtime_ms`
- `avg_generation_latency_ms`
- `avg_prompt_tokens`
- `avg_completion_tokens`
- `avg_total_tokens`
- mutation 体量统计
- `mutation_kill_rate = killed / total_mutants`
- `mutation_effective_kill_rate = killed / (killed + survived)`

其中几个设计点要固定：

- `test_pass_rate` 的分母是 `compile_pass_count`，不是 `total_samples`
- 覆盖率、token、耗时等平均值只对非空值求均值
- mutation 相关统计要区分：
  - 总 mutant
  - 有效 mutant
  - `no_tests`
  - `not_checked`
  - `timeout`
  - `skipped`
  - `suspicious`

### 6.2 分组指标

分组逻辑统一由 `group_metrics(rows, group_keys)` 负责，避免每种维度写一套新逻辑。

当前建议的固定维度：

- `by_model`
- `by_language`
- `by_complexity`
- `by_scenario`
- `by_model_language`

这样做的价值是：

- 便于 CSV 导出
- 便于 HTML 图表直接消费
- 后续新增维度时只需要增加 group key

## 7. 故障分析设计

`diagnostics.py` 负责把错误信息从“原始字符串”收敛成“可统计类别”。

处理流程：

1. 从 `compile_error`、`test_error`、`coverage_error`、`mutation_error` 中按阶段顺序选首个错误。
2. 基于错误文本做分类：
   - `module_not_found`
   - `name_error`
   - `assertion_failure`
   - `import_error`
   - `syntax_error`
   - `timeout`
   - `stop_iteration`
   - `recursion_error`
   - `other`
3. 生成带示例的 Top N 失败统计：
   - `stage`
   - `error_type`
   - `count`
   - `example_model`
   - `example_sample`
   - `example_message`

这样设计的目的不是替代完整日志，而是给横向对比提供“高频失败画像”。

## 8. HTML 报告设计

HTML 不是简单导出表格，而是给人工阅读做一层可视化包装。

页面结构建议固定为三段：

1. 摘要
   - 总样本数
   - 编译通过率
   - 测试通过率
   - 平均覆盖率
   - 平均变异得分
   - 平均耗时 / token
2. 详细分析
   - 模型柱状图
   - 复杂度趋势图
   - 多维雷达图
   - 模型 x 语言热力图
   - 高频失败类型表
   - 数据质量说明
3. 附录
   - 按模型统计表
   - 按语言统计表
   - 样本明细表
   - 前端筛选器

阈值设计：

- 通过 CLI 注入阈值，不写死在模板里
- HTML 只负责展示阈值，不负责重新计算业务结论

筛选器设计：

- 模型
- 语言
- 测试状态
- 覆盖率状态
- 变异状态
- `sample_id` 关键词

## 9. 稳定性与兼容性原则

为了和 evaluator 契约对齐，reporter 应遵守下面几条：

- 不修改 evaluator summary 的原始语义
- 所有缺失指标保持 `None`
- 历史数据缺字段时允许降级，但要标明来源
- 文件自动发现逻辑要优先使用最新 evaluator summary
- 输出 JSON 的字段命名应保持稳定，便于后续脚本消费

## 10. 当前实现现状

结合现有代码，当前 reporter 已经具备可运行能力，但存在两个明显的结构问题。

### 10.1 单文件实现和拆分模块并存

当前仓库里既有：

- `pipeline.py` 中的大量内联逻辑

也有：

- `loader.py`
- `aggregator.py`
- `diagnostics.py`
- `exporters/`

这说明 reporter 处于“从单文件向模块化迁移”的中间状态。继续放任下去的风险是：

- 同一逻辑维护两份
- 修 bug 时只改到一边
- 文档、测试、实现逐步漂移

### 10.2 中文标签存在编码异常

`contracts.py` 和 HTML 模板里能看到明显的中文乱码字符串。这说明当前文件编码或历史复制过程存在问题。这个问题不会直接影响聚合逻辑，但会影响：

- HTML 展示质量
- CSV/表头可读性
- 后续维护成本

## 11. 推荐落地方案

建议把 reporter 收敛成下面这种结构：

```text
benchmark/reporter/
  __main__.py
  __init__.py
  pipeline.py          # 只保留 orchestration
  loader.py            # 输入读取与样本 enrich
  aggregator.py        # 指标聚合
  diagnostics.py       # 错误归类
  contracts.py         # 常量、字段、标签
  common.py            # 通用函数
  exporters/
    json_exporter.py
    csv_exporter.py
    html_exporter.py
  DESIGN.md
```

建议的演进顺序：

1. 先把 `pipeline.py` 改成只调用拆分模块。
2. 清理重复实现，保证每类逻辑只有一个真实来源。
3. 修复中文编码问题。
4. 再补 reporter 级测试，重点覆盖：
   - metadata 优先级
   - runner summary fallback
   - 空值处理
   - 分组聚合口径
   - failure breakdown 分类

## 12. 最小测试建议

如果后续补测试，至少应覆盖这些场景：

- 指定 `--evaluator-summary` 时能正确读取目标文件
- 未指定时能自动找到最新 summary
- metadata 存在时优先使用 metadata 中的耗时和 token
- metadata 缺失时回退到 `runner_summary`
- 二者都缺失时输出 `None` 且 `generation_metrics_source=missing`
- `test_pass_rate` 分母使用 `compile_pass_count`
- mutation 统计字段为空时不应被误算成 0 分样本
- 错误分类能把典型 Python 错误映射到稳定类别

## 13. 结论

`reporter` 的正确定位不是“再做一次评测”，而是“把评测结果变成稳定、可解释、可汇报的分析视图”。

所以它的设计重点应该是：

- 契约稳定
- 空值语义正确
- 聚合口径统一
- 导出层与计算层解耦
- HTML 可读但不侵入核心逻辑

按照这个思路，reporter 后续既可以继续输出静态报告，也可以很自然地扩展成 API 或 dashboard 的后端数据层。
