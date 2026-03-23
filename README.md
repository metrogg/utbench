# ut-bench

> **多模型单测生成效果横向评测平台**
>
> AI Unit Test Generation — Multi-Model Benchmark

---

## 项目简介

`ut-bench` 是一个针对 **AI 辅助单元测试生成能力** 的横向评测平台，旨在系统性地评估和对比不同大语言模型（LLM）在自动生成单元测试方面的效果。

### 解决的问题

- 不同模型（GPT-4o、Claude、DeepSeek、混元等）生成单测的质量差异难以量化
- 缺乏统一的评测标准和基准数据集
- 生成的测试用例在编译通过率、执行通过率、覆盖率等维度缺乏横向对比
- 难以指导工程实践中的模型选型决策

### 核心目标

1. **建立标准化 Benchmark 数据集**：覆盖 Java、Python、Go 多语言，涵盖不同复杂度的代码样本
2. **定义统一评测指标体系**：正确性、覆盖率、有效性、效率、工程质量五大维度
3. **自动化评测流水线**：端到端自动完成生成 → 编译 → 执行 → 覆盖率 → 变异测试全链路
4. **可视化评测报告**：多模型横向对比，支持按语言、复杂度、场景多维度分析

---

## 目录结构

```
ut-bench/
├── README.md                     # 项目总览（本文件）
├── docs/                         # 技术与方案文档
│   ├── design/                   # 方案设计文档
│   │   ├── benchmark-design.md   # 评测方案总体设计
│   │   ├── metrics-definition.md # 评测指标定义
│   │   └── dataset-design.md     # 数据集设计
│   ├── research/                 # 调研报告
│   │   ├── model-survey.md       # 模型能力调研
│   │   └── tools-survey.md       # 相关工具调研
│   └── reports/                  # 评测报告（历次评测结果文档）
├── dataset/                      # 评测基准数据集
│   ├── README.md                 # 数据集说明
│   ├── java/                     # Java 语言测试样本
│   ├── python/                   # Python 语言测试样本
│   └── go/                       # Go 语言测试样本
├── benchmark/                    # 评测核心代码
│   ├── README.md                 # 评测模块说明
│   ├── runner/                   # 评测执行器（调用模型 API 生成测试）
│   ├── evaluator/                # 评测指标计算（覆盖率、变异测试等）
│   ├── reporter/                 # 评测报告生成
│   └── config/                   # 评测配置
│       └── models.yaml           # 模型列表与 API 配置
├── scripts/                      # 工具脚本
│   ├── setup.sh                  # 环境初始化
│   ├── run_benchmark.sh          # 运行完整评测
│   └── gen_report.sh             # 生成评测报告
└── results/                      # 评测结果存档
```

---

## 评测指标体系

| 维度 | 指标 | 说明 |
|------|------|------|
| **正确性** | 编译通过率 | 生成的测试代码能否成功编译 |
| **正确性** | 测试执行通过率 | 编译后测试用例能否全部通过 |
| **覆盖率** | 行覆盖率 | 测试覆盖的代码行比例 |
| **覆盖率** | 分支覆盖率 | 测试覆盖的分支比例 |
| **覆盖率** | 函数覆盖率 | 测试覆盖的函数比例 |
| **有效性** | 变异测试得分 | Mutation Score，衡量测试的缺陷检测能力（目标 ≥ 85%） |
| **效率** | 生成耗时 | 完成单个文件测试生成的时间 |
| **效率** | Token 消耗 | 每次生成消耗的 Token 数量 |
| **工程质量** | Mock 准确率 | 依赖 Mock 的准确性和完整性 |
| **工程质量** | 断言充分性 | 断言覆盖核心逻辑的比例 |

---

## 参评模型

| 模型 | 提供方 | 备注 |
|------|--------|------|
| GPT-4o | OpenAI | 基准对比 |
| Claude 3.5 Sonnet | Anthropic | 基准对比 |
| DeepSeek-V3 | DeepSeek | 国产模型代表 |
| 混元 | 腾讯 | 内部模型 |

> 持续扩充中，欢迎提交 PR 添加新模型。

---

## 数据集说明

评测数据集覆盖以下维度：

- **编程语言**：Java、Python、Go
- **代码复杂度**：简单（纯函数）、中等（含依赖）、复杂（含 DB/网络/文件 IO）
- **业务场景**：工具类、服务类、数据处理类

详见 [dataset/README.md](dataset/README.md)

---

## 快速开始

### 环境准备

```bash
# 初始化环境
bash scripts/setup.sh
```

### 运行评测

```bash
# 运行完整评测（所有模型 × 所有数据集）
bash scripts/run_benchmark.sh

# 运行单个模型评测
bash scripts/run_benchmark.sh --model deepseek --lang java
```

### 生成报告

```bash
bash scripts/gen_report.sh
```

---

## 相关文档

- [评测方案总体设计](docs/design/benchmark-design.md)
- [评测指标定义](docs/design/metrics-definition.md)
- [数据集设计](docs/design/dataset-design.md)
- [模型能力调研](docs/research/model-survey.md)

---

## 贡献指南

1. Fork 本仓库
2. 新建 `feature/xxx` 分支
3. 提交代码并发起 MR
4. 等待 Review 通过后合并

---

## 项目背景

本项目是**校企联合项目「多模型单测自动生成横向评测平台」**的核心代码仓库，旨在通过系统性评测，为 CSIG 单测质量防护体系的模型选型和能力优化提供数据支撑。

---

*Last updated: 2026-03*
