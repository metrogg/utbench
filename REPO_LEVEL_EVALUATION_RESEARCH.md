# 仓库级(Repo-Level)单元测试生成评测：技术调研与UT-Bench建议

> 调研时间：2026-05-04
> 目标：为UT-Bench添加仓库级评测能力，梳理行业技术方案并提出具体建议

---

## 一、行业项目全景概览

| 项目 | 核心定位 | 仓库级方式 | 语言覆盖 | 环境管理 | 评测指标 |
|------|----------|-----------|---------|---------|---------|
| **SWE-bench** | Issue修复评测 | PR→任务实例,FAIL/PASS契约 | Python | Conda环境映射 | Resolved Rate |
| **SWE-Bench++** | SWE-bench自动化扩展 | 4阶段全自动pipeline | 11语言 | 自动合成Docker | Resolved Rate + AutoQA |
| **SWE-rebench** | 防污染持续评测 | 自动化持续采集 | Python | GitHub真实环境 | Resolved Rate |
| **TestGenEval** | 单元测试生成评测 | 基于SWE-bench转换 | Python | Docker(预构建镜像) | pass@k, 覆盖率, 变异得分 |
| **FEA-Bench** | 特性开发评测 | PR过滤→新组件识别 | Python | Docker(基于SWE-bench) | Resolved Ratio |
| **ExecRepoBench** | 可执行代码补全 | AST掩码+单元测试验证 | Python | 隐式(执行环境) | Pass@k |
| **CrossCodeEval** | 跨文件代码补全 | 静态分析提取跨文件依赖 | 4语言 | 无执行 | EM, ES, Identifier Match |
| **RepoBench** | 仓库级代码补全 | 跨文件上下文检索 | Python, Java | 无执行 | EM, ES |
| **RepoST** | 可扩展环境构建 | 沙盒隔离函数+依赖 | Python | 沙盒测试(非全仓库) | Pass@1 |
| **TestPilot** | LLM测试生成 | 函数签名+文档+示例 | JS/TS | npm本地执行 | 通过率(隐式) |
| **EvalPlus** | 函数级代码评测 | 无仓库级(仅函数级) | Python | 进程隔离 | pass@k, DPS |
| **Aider** | 代码编辑评测 | Exercism练习题 | Python | 本地执行 | 通过率(2次机会) |
| **LLM变异测试研究** | 变异测试效果 | 851个真实Bug | Java | 多服务器并行 | 检测率, 耦合率, 可编译率 |

---

## 二、各项目详细技术方案

### 2.1 SWE-bench：仓库级评测的奠基者

#### 仓库选择标准
- **12个知名Python开源项目**：Django, Flask, Scikit-learn, Matplotlib, Sympy, Requests等
- 选择标准：知名度高、活跃开发、测试覆盖完善、Issue/PR历史丰富
- 领域多样性：Web框架、数据科学、符号计算等

#### 任务提取流程（3步Pipeline）
```
Step 1: print_pulls.py → 从GitHub API收集PR元数据
Step 2: build_dataset.py → PR→结构化任务实例(含验证过滤)
Step 3: get_tasks_pipeline.py → 多仓库自动化+并行Token
```

**关键过滤规则**：
- PR必须包含closing keywords（如`fixes #123`）
- 必须产生有效的代码diff（空patch被过滤）
- Lite版额外过滤：patch≤25行、≤1文件、≤3 hunk、无超链接/图片

#### 任务实例Schema
```json
{
  "instance_id": "django__django-12345",
  "repo": "django/django",
  "pull_number": 12345,
  "issue_numbers": ["67890"],
  "base_commit": "419b7d...",
  "patch": "黄金修复patch(不含测试文件)",
  "test_patch": "测试patch(仅测试文件)",
  "problem_statement": "Issue标题+正文",
  "hints_text": "首次commit前的讨论评论",
  "created_at": "2023-01-15T10:30:00Z",
  "FAIL_TO_PASS": ["test_foo", "test_bar"],
  "PASS_TO_PASS": ["test_baz", "test_qux"],
  "version": "4.2",
  "environment_setup_commit": "abc123..."
}
```

**核心设计**：`FAIL_TO_PASS` + `PASS_TO_PASS` 双重验证契约——修复必须让失败测试通过，且不能破坏已有通过的测试。

#### 环境管理
- **Conda环境映射**：每个任务实例指定`environment_setup_commit`，映射到对应的Conda环境规格（Python版本、依赖版本）
- 每次评测前创建全新的Conda环境
- 仓库checkout到`base_commit`

#### 评测指标
- **Resolved Rate**：同时满足FAIL_TO_PASS全通过 + PASS_TO_PASS全通过的实例比例

---

### 2.2 SWE-Bench++：全自动多语言扩展

#### 核心创新（vs SWE-bench）
| 维度 | SWE-bench | SWE-Bench++ |
|------|-----------|-------------|
| 数据构建 | 手动 | **全自动** |
| 规模 | 12仓库 | **3,971仓库, 11,133实例** |
| 语言 | Python | **11语言** |
| 任务类型 | Bug修复 | Bug修复+功能需求(38.5%) |
| 环境策略 | 预构建Conda | **自动合成Docker** |
| 日志解析 | 静态regex | **自适应解析器(regex+神经回退)** |
| 发布方式 | 静态 | **持续"活"数据集(防污染)** |

#### 4阶段自动化Pipeline
```
Stage 1: 程序化源采集
  ↓ 扫描GitHub firehose，可扩展过滤器
Stage 2: 环境合成
  ↓ 模板脚手架 + LLM迭代修正(最多5轮) → Docker环境
Stage 3: 状态差分测试预言机提取
  ↓ 三状态执行(Base/Before/After) + 自适应日志解析
Stage 4: 自动质量保证
  ↓ 4层AutoQA + 人工验证子集
```

#### 仓库选择5项标准
1. **活跃维护**：近期有提交活动
2. **社区采用**：>100 stars + 可识别的测试框架
3. **足够复杂**：代码量>10,000 LOC
4. **合并PR+Issue关联**：PR必须合并且明确解决关联Issue
5. **测试文件变更**：PR必须包含测试文件的修改或新增

#### 环境合成（最重要创新）
- **混合架构**：语言特定Dockerfile模板 + LLM迭代修正
- 模板包含语义占位符，LLM填充 → 尝试构建 → 失败时捕获stderr → LLM生成修正方案 → 重试(最多5次)
- 构建成功后进入**Test-Run反馈循环**验证测试执行器正常工作
- Python产出率41%，C++产出率9.5%

#### 三状态差分分类
- **状态A(Bug修复)**：Before状态构建成功 → F2P(修复的测试) + P2P(回归检查)
- **状态B(功能需求)**：Before状态构建失败(缺少符号) → F2P(新增测试在After通过) → 构建失败被视为语义信号而非错误

#### 4层AutoQA
| 层 | 目的 | 方法 |
|----|------|------|
| Layer 1 | 环境确定性 | 构建3次，全部成功才保留 |
| Layer 2 | 预言机一致性 | 运行黄金方案3次，结果一致才保留 |
| Layer 3 | 语义对齐 | LLM-Judge评估问题-测试对齐度 |
| Layer 4 | 假阴性过滤 | 检查SOTA模型失败是否由基础设施导致 |

---

### 2.3 SWE-rebench：防污染持续评测

#### 核心价值
- **21,000+交互式Python任务**，全自动化持续采集
- 解决SWE-bench的**数据污染问题**：静态benchmark会被模型训练数据覆盖
- 可持续生成**全新任务**，保证评测有效性
- 适合RL训练(不仅评测)

#### 防污染策略
- 持续采集新任务，可通过PR日期过滤控制污染
- 实验表明：部分模型在SWE-bench Verified上的分数可能因污染被高估

---

### 2.4 TestGenEval：最接近UT-Bench的评测项目

#### 任务设计（2种任务）
| 任务 | 描述 | 场景 |
|------|------|------|
| **UnitTest Completion** | 在已有测试文件中补全测试 | first/last/extra test |
| **Full File Generation** | 从零生成完整测试文件 | full |

#### 数据集构建
- **基于SWE-bench转换**：`transform_swebench.py`将SWE-bench数据转为测试生成任务
- **覆盖率过滤**：`filter_unittests.py`移除黄金测试对目标代码无覆盖的实例
- **1,210个代码-测试文件对**，来自11个仓库(3,523-78,287 stars)
- 代码文件平均782 LOC，测试文件平均677 LOC
- 黄金测试中位覆盖率60.4%

#### 环境管理
- **Docker预构建镜像**：每个仓库版本对应一个Docker镜像
- 预装coverage和变异测试依赖
- 支持本地构建(Makefile)或从DockerHub拉取
- Conda管理评测框架本身

#### 评测指标（最值得借鉴）
| 任务 | 指标 |
|------|------|
| 测试补全 | pass@k, 覆盖率提升, 覆盖率提升@pass, avg pass@5 |
| 全文件生成 | pass@1, all pass@1, 覆盖率, 覆盖率@pass, **变异得分**, 变异得分@pass |

**关键设计**：`@pass`变体——仅在通过的测试上计算覆盖率和变异得分，防止"总是通过的断言"刷分。

#### Pipeline流程
```
1. run_pipeline.py → 模型预测(4种prompt设置)
2. run_evaluation.py → Docker内执行测试 + 计算指标
3. generate_report.py → 聚合报告
```

---

### 2.5 FEA-Bench：仓库级特性开发评测

#### 任务提取流程
```
1. 仓库采集: Top PyPI 8000包 → license+>1000 PR → 快速验证(20 PR试跑pytest) → 119仓库
2. PR爬取: 只保留修改了测试文件的合并PR
3. 新组件提取: 解析gold patch，比较前后状态，识别新增的类和函数
4. 规则过滤: ≥1新组件 + 新组件占>25%修改行
5. 意图过滤: GPT-4o分类PR意图，仅保留"新功能"
6. 长度过滤: patch ≤ 8K tokens
7. 验证: 应用test_patch→部分失败; 应用gold patch→全部通过
   ⚠️ 允许gold前的AttributeError/ImportError(新组件不存在)
```

#### 最终数据集
- **1,401实例**，83个仓库
- 平均编辑2.62个文件，128.5行
- 平均新增4.49个函数，0.78个类

#### 评测方法
- **执行评测**：模型patch应用后，全部单元测试通过=resolved
- 4维度提示配置：新组件提示(简略/详细) × 检索方法(Oracle/BM25) × 输出格式(Natural/Patch) × 上下文长度(27K/40K)
- **关键发现**：Natural格式显著优于Patch格式（LLM难以遵循严格patch格式）

#### FEA-Bench Lite
- 200任务，48仓库
- 过滤：低质量(描述<40词、级联Issue、含图片)、高难度(文件删除、>3代码文件、>10 hunk、>4K tokens变更、含新类、>10新增函数)

---

### 2.6 ExecRepoBench：可执行的仓库级补全

#### 核心创新
- **可执行验证**：将生成的代码填回仓库，运行单元测试验证
- 区别于RepoBench/CrossCodeEval仅用EM/ES静态指标

#### 数据集构建
- 25个活跃Python仓库，1.0-1.5K样本
- **AST感知掩码**：基于tree-sitter解析，按Expression/Statement/Function/Class四级粒度掩码
- 严格时间限制：<2分钟通过所有测试
- 输入截断：8K tokens

#### 评测指标
- **Pass@k**：基于执行结果的功能正确性

---

### 2.7 CrossCodeEval：跨文件依赖评测

#### 核心方法
- **静态分析提取跨文件依赖**：
  1. 用空类替换import语句
  2. 运行静态分析找undefined names
  3. 这些undefined names定位跨文件使用

#### 评测配置
- 3种上下文设置：仅文件内/检索跨文件/带参考检索(上界估计)
- 2类指标：代码匹配(EM, ES) + 标识符匹配

---

### 2.8 RepoST：可扩展沙盒环境

#### 核心创新
- **沙盒测试**：不构建整个仓库，而是**隔离目标函数+其依赖到独立脚本**
- 避免全仓库构建的复杂性（外部服务、配置、互联依赖）
- 大幅提升可扩展性：7,415函数，832仓库

#### 关键洞察
> 全仓库执行是不必要的复杂——只需提取函数及其直接依赖即可获得有意义的执行反馈。

---

### 2.9 TestPilot：LLM测试生成工具

#### 架构
```
1. 识别目标函数(包导出)
2. 收集上下文: 签名 + 函数体 + 文档示例(自动挖掘)
3. 用测试骨架提示LLM
4. 解析LLM响应为可运行测试
5. 执行测试(可选)
6. 失败→反馈给LLM→重新提示(迭代精化)
```

#### 特点
- 不需要训练/微调/少样本示例
- 上下文嵌入在代码注释中
- 已归档，后继为TestPilot2

---

### 2.10 LLM变异测试研究

#### 关键发现
| 指标 | 规则方法平均 | LLM平均 | 提升 |
|------|-------------|---------|------|
| 真实Bug检测率 | 41.64% | 79.07% | +37.43pp |
| 耦合率 | 24.38% | 46.21% | +21.83pp |
| 可编译率 | ~99.9% | 63.8% | -36.1pp |
| 重复率 | 0% | 13.1% | +13.1pp |
| 等价变异率 | ~1.1% | 5.3% | +4.2pp |

**对UT-Bench的启示**：
- LLM生成的变异体检测能力强但可编译率低，变异测试评测需要**过滤不可编译变异体**
- GPT-4o-Mini检测率最高(90.8%)，但编译率仅~76%
- 变异测试应作为**测试质量**的指标，而非仅看覆盖率

---

### 2.11 EvalPlus

- **不支持仓库级评测**：仅函数级(HumanEval+ 164任务, MBPP+ 378任务)
- 测试增强：将原始测试扩充80倍(HumanEval)和35倍(MBPP)
- 进程隔离执行，4GB内存限制
- **可借鉴**：测试增强方法、进程隔离策略、sanitization管道

---

### 2.12 Aider Benchmark

- 基于Exercism Python的133道练习题
- **2次机会**评测：第1次失败后可看测试错误输出再试
- GPT永远不看到测试源码，只看到失败输出
- **纯文本编辑格式优于函数调用格式**
- 变差控制：temperature=0 + SHA哈希记录 + 去除时间信息

---

## 三、技术方案横向对比

### 3.1 仓库选择策略对比

| 项目 | 来源 | 核心筛选标准 | 最终仓库数 |
|------|------|-------------|-----------|
| SWE-bench | 手工选择知名项目 | 知名度+测试覆盖+活跃度 | 12 |
| SWE-Bench++ | GitHub firehose | >100 stars + >10K LOC + 测试框架 + 活跃 | 3,971 |
| TestGenEval | SWE-bench衍生 | SWE-bench仓库 + 覆盖率过滤 | 11 |
| FEA-Bench | Top PyPI | license + >1000 PR + pip install可工作 | 83 |
| ExecRepoBench | GitHub搜索 | 活跃更新 + 可创建测试 + <2min通过 | 25 |
| RepoST | GitHub | 函数可隔离 + 依赖可提取 | 832 |

**UT-Bench建议**：采用FEA-Bench的"快速验证"策略——先粗筛(license/stars/PR数)，再对每个仓库抽样PR试跑测试，大幅降低无效仓库比例。

### 3.2 任务提取策略对比

| 项目 | 提取方式 | 过滤层级 | 产出率 |
|------|---------|---------|--------|
| SWE-bench | PR→Issue关联+测试变更 | 2层(有效PR+有效实例) | 低 |
| SWE-Bench++ | PR→全自动4阶段 | 4层AutoQA | 8.1% |
| TestGenEval | SWE-bench→测试配对转换 | 覆盖率>0过滤 | 中 |
| FEA-Bench | PR→新组件识别+GPT-4o意图分类 | 6层 | 中 |
| ExecRepoBench | AST掩码 | 语法正确性 | 高 |
| RepoST | 函数提取+依赖隔离 | 可执行性 | 41%(Python) |

**UT-Bench建议**：UT-Bench的任务不是"修复Bug"而是"生成测试"，因此提取逻辑应改为——
1. 从仓库中识别**测试覆盖不足的模块/函数**
2. 提取函数签名+跨文件依赖上下文
3. 过滤：确保目标函数有可运行的import链

### 3.3 环境管理方案对比

| 项目 | 环境方案 | 隔离粒度 | 构建方式 | 可扩展性 |
|------|---------|---------|---------|---------|
| SWE-bench | Conda环境映射 | 每实例 | 手动映射表 | 低 |
| SWE-Bench++ | 自动合成Docker | 每仓库版本 | 模板+LLM修正 | 高 |
| TestGenEval | 预构建Docker镜像 | 每仓库版本 | Makefile/DockerHub | 中 |
| FEA-Bench | Docker(基于SWE-bench) | 每仓库版本 | pip install -e . + pytest | 中 |
| RepoST | 沙盒测试(函数级隔离) | 每函数 | 自动依赖提取 | 高 |
| EvalPlus | 进程隔离 | 每样本 | 无需构建 | 极高 |
| UT-Bench(当前) | Docker(全仓库) | 每任务 | 手动Dockerfile | 低 |

**UT-Bench建议**：采用**SWE-Bench++的模板+LLM修正**方案，但降级为**RepoST的沙盒隔离**思路——
- Phase 1：为每个仓库构建Dockerfile模板(pip install -e . + pytest)
- Phase 2：LLM自动修正构建失败的环境
- Phase 3(可选)：对复杂仓库降级为RepoST沙盒模式

### 3.4 评测指标对比

| 指标 | SWE-bench | TestGenEval | FEA-Bench | ExecRepoBench | UT-Bench(当前) |
|------|-----------|-------------|-----------|---------------|---------------|
| Resolved/Pass Rate | ✅ | ✅(pass@k) | ✅ | ✅(Pass@k) | ✅(TestPass) |
| 代码覆盖率 | ❌ | ✅(含@pass变体) | ❌ | ❌ | ✅ |
| 变异测试 | ❌ | ✅(含@pass变体) | ❌ | ❌ | ✅ |
| FAIL_TO_PASS/PASS_TO_PASS | ✅ | ❌ | ❌ | ❌ | ❌ |
| Existing Tests Still Pass | 隐含 | ❌ | 隐含 | ❌ | ❌ |
| 覆盖率提升 | ❌ | ✅ | ❌ | ❌ | ❌ |
| Flaky Test Rate | ❌ | ❌ | ❌ | ❌ | ❌ |
| Cost Efficiency | ❌ | ❌ | ❌ | ❌ | ✅(部分) |

---

## 四、对UT-Bench仓库级评测的具体建议

### 4.1 数据集Schema建议

综合SWE-bench、FEA-Bench、TestGenEval的最佳实践，建议UT-Bench仓库级任务实例Schema如下：

```json
{
  "instance_id": "django__query_001",
  "source": {
    "repo_url": "https://github.com/django/django",
    "base_commit": "419b7d...",
    "license": "BSD-3",
    "stars": 78000
  },
  "target": {
    "target_files": ["django/db/models/query.py"],
    "target_functions": ["QuerySet.annotate", "QuerySet.filter"],
    "target_classes": ["QuerySet"],
    "test_directory": "tests/query/"
  },
  "context": {
    "context_scope": "module",
    "dependency_level": 3,
    "relevant_files": [
      "django/db/models/sql/query.py",
      "django/db/models/manager.py"
    ],
    "cross_file_dependencies": {
      "imports_from_other_modules": ["django.db.models.sql"],
      "type_references": ["Query"]
    }
  },
  "environment": {
    "language": "python",
    "language_version": "3.11",
    "setup_commands": ["pip install -e ."],
    "test_runner": "pytest",
    "pinned_deps": {
      "django": "5.0.1",
      "pytest": "8.0.0",
      "coverage": "7.4.0",
      "mutmut": "2.5.0"
    },
    "dockerfile_template": "python-generic"
  },
  "evaluation": {
    "existing_tests_command": "pytest tests/query/ --tb=short",
    "existing_tests_expected_pass": ["test_annotate", "test_filter_basic", "test_filter_chain"],
    "coverage_target": "django/db/models/query.py",
    "mutation_target": "django/db/models/query.py",
    "timeout_seconds": 120,
    "must_pass_existing": true
  },
  "difficulty": {
    "dependency_level": 3,
    "scenario": "complex_dependency",
    "estimated_loc": 150,
    "num_target_functions": 2,
    "tags": ["orm", "database", "query-builder"]
  },
  "golden_test": {
    "test_patch": "--- a/tests/query/test_query.py\n+++ b/tests/query/test_query.py\n...",
    "coverage_achieved": 0.72,
    "mutation_score_achieved": 0.65
  }
}
```

**新增关键字段说明**：

| 字段 | 借鉴来源 | 作用 |
|------|---------|------|
| `source.base_commit` | SWE-bench | 锁定仓库快照，保证可复现 |
| `evaluation.existing_tests_expected_pass` | SWE-bench PASS_TO_PASS | 回归检查，防止生成测试破坏已有功能 |
| `evaluation.must_pass_existing` | SWE-bench | 强制验证已有测试仍通过 |
| `context.relevant_files` | FEA-Bench | 指定检索范围，减少无关上下文 |
| `context.cross_file_dependencies` | CrossCodeEval | 明确跨文件依赖关系 |
| `golden_test` | TestGenEval | 提供黄金标准对比 |
| `difficulty` | SWE-Bench++ | 难度分级，支持分层评测 |
| `environment.dockerfile_template` | SWE-Bench++ | 模板化环境构建 |
| `environment.pinned_deps` | SWE-bench | 依赖版本锁定 |

### 4.2 数据集构建Pipeline建议

```
Phase 1: 仓库采集与验证
├── 1.1 来源: GitHub热门仓库(Go/Python/Java/C++)
│   ├── 筛选: license + >500 stars + 测试目录存在
│   └── 快速验证: 抽样10个commit, 检查test命令可运行
├── 1.2 仓库镜像(可选, 借鉴SWE-bench)
│   ├── 创建mirror repo到组织下(防止上游变动)
│   └── 清理CI/CD配置(.github/workflows等)
└── 1.3 目标: 50-100个仓库(每语言10-25个)

Phase 2: 目标函数提取
├── 2.1 静态分析(借鉴CrossCodeEval)
│   ├── 提取所有公开API(函数/类/方法)
│   ├── 识别跨文件依赖(import链分析)
│   └── 计算dependency_level
├── 2.2 覆盖率基线测量
│   ├── 运行现有测试套件, 收集覆盖率
│   ├── 识别覆盖不足的函数(行覆盖<50%)
│   └── 标记无测试的函数(覆盖=0%)
├── 2.3 过滤
│   ├── 排除纯getter/setter
│   ├── 排除private/内部函数
│   ├── 排除需要外部服务的函数(DB/网络/文件系统)
│   └── 确保目标函数有可运行的import链
└── 2.4 目标: 每仓库5-20个目标函数

Phase 3: 环境构建(借鉴SWE-Bench++)
├── 3.1 模板化Dockerfile
│   ├── python-generic: pip install -e . + pytest + coverage + mutmut
│   ├── go-generic: go test + go tool cover + go-mutesting
│   ├── java-generic: mvn test + jacoco + pitest
│   └── cpp-generic: cmake + gtest + gcov + mull
├── 3.2 LLM辅助环境修正(借鉴SWE-Bench++)
│   ├── 构建失败 → 捕获stderr → LLM生成修正
│   ├── 最多5轮重试
│   └── 记录成功/失败率
├── 3.3 环境验证(3次确定性检查)
│   ├── 构建3次，全部成功才保留
│   ├── 运行测试3次，结果一致才保留
│   └── 排除flaky仓库
└── 3.4 预构建镜像推送到Registry

Phase 4: 黄金测试收集(借鉴TestGenEval + SWE-bench)
├── 4.1 从PR中提取测试patch
│   ├── 识别修改了测试文件的PR
│   ├── 提取新增/修改的测试用例
│   └── 验证测试在gold patch后通过
├── 4.2 或: 手写/LLM辅助生成黄金测试
│   ├── 用强模型(GPT-4o)为每个目标函数生成参考测试
│   ├── 人工审核过滤
│   └── 记录黄金覆盖率和变异得分
└── 4.3 记录existing_tests_expected_pass列表

Phase 5: 质量保证(借鉴SWE-Bench++ AutoQA)
├── 5.1 环境确定性: 3次构建一致性
├── 5.2 测试稳定性: 3次运行一致性
├── 5.3 语义对齐: LLM-Judge验证函数-测试对齐
├── 5.4 假阴性: 检查强模型失败是否因基础设施问题
└── 5.5 防污染: 检查代码是否出现在公开训练数据中(时间过滤)
```

### 4.3 环境管理架构建议

```
┌─────────────────────────────────────────────────┐
│               UT-Bench Orchestrator              │
├─────────────────────────────────────────────────┤
│  Environment Manager                            │
│  ├── Dockerfile Template Registry               │
│  │   ├── python-generic/Dockerfile              │
│  │   ├── go-generic/Dockerfile                  │
│  │   ├── java-maven/Dockerfile                  │
│  │   └── cpp-cmake/Dockerfile                   │
│  ├── Image Build Queue                          │
│  │   ├── Build → Test → Push to Registry        │
│  │   └── LLM修正(最多5轮)                       │
│  └── Image Registry (DockerHub/内网)            │
├─────────────────────────────────────────────────┤
│  Sandbox Runner                                 │
│  ├── 全仓库模式(Docker)                         │
│  │   ├── 源码只读挂载: -v src:/workspace/src:ro │
│  │   ├── 测试写入: -v output:/workspace/tests   │
│  │   └── 资源限制: --memory=4g --cpus=2         │
│  ├── 沙盒模式(RepoST风格)                       │
│  │   ├── 函数+依赖提取到独立脚本                │
│  │   ├── 进程级隔离(借鉴EvalPlus)               │
│  │   └── 适用于环境构建失败的仓库                │
│  └── 回退: 本地模式(已有)                       │
├─────────────────────────────────────────────────┤
│  Safety Guard                                   │
│  ├── 交互轮次限制: max 50轮                     │
│  ├── 超时: 120s/test, 600s/instance             │
│  ├── 自动回滚: git init + commit → git checkout │
│  ├── 孤儿容器清理: 定时扫描 + 退出钩子          │
│  └── 资源监控: docker stats采样                  │
└─────────────────────────────────────────────────┘
```

**双层环境策略**（借鉴RepoST思路）：
1. **首选全仓库Docker模式**：完整环境，更真实
2. **回退到沙盒模式**：仅当Docker构建失败时，提取函数+依赖到独立脚本运行

### 4.4 评测指标体系建议

```yaml
metrics:
  correctness:
    - name: test_pass_rate
      description: "生成的测试是否通过(pass@1, pass@5)"
      weight: 0.15
    - name: existing_tests_still_pass
      description: "已有测试是否仍然通过(PASS_TO_PASS契约)"
      weight: 0.10
    - name: flaky_test_rate
      description: "同一测试跑3次的不一致率"
      weight: 0.05

  coverage:
    - name: line_coverage
      description: "行覆盖率"
      weight: 0.10
    - name: branch_coverage
      description: "分支覆盖率"
      weight: 0.10
    - name: coverage_improvement
      description: "相比现有测试的覆盖率提升(借鉴TestGenEval)"
      weight: 0.05

  mutation:
    - name: mutation_score
      description: "变异测试得分"
      weight: 0.15
    - name: mutation_score_at_pass
      description: "仅计算通过测试的变异得分(借鉴TestGenEval)"
      weight: 0.10

  quality:
    - name: syntax_error_rate
      description: "语法错误占比"
      weight: 0.05
    - name: assertion_quality
      description: "断言质量(非平凡断言占比)"
      weight: 0.05
    - name: test_readability
      description: "LLM-as-Judge评估可读性"
      weight: 0.05

  efficiency:
    - name: cost_per_coverage_point
      description: "$ / line_coverage(借鉴Aider)"
      weight: 0.03
    - name: tokens_per_test_case
      description: "总token数 / 测试用例数"
      weight: 0.02
```

**核心创新指标**（UT-Bench独有）：
- **`mutation_score_at_pass`**：借鉴TestGenEval的@pass设计，防止"总是通过的断言"刷变异得分
- **`coverage_improvement`**：测试生成不仅是覆盖，更要**提升**覆盖
- **`existing_tests_still_pass`**：SWE-bench的PASS_TO_PASS契约在测试生成场景的映射

### 4.5 评测流程建议

```
1. 环境准备
   ├── 拉取/构建Docker镜像
   ├── checkout到base_commit
   ├── 安装pinned依赖
   └── 运行existing_tests, 记录PASS_TO_PASS基线

2. 测试生成
   ├── 模型接收prompt(目标函数+上下文+依赖信息)
   ├── 生成测试代码
   └── (可选)CLI Agent: 执行→看输出→修正(最多N轮)

3. 评测执行
   ├── 3a. 编译检查(syntax_error_rate)
   ├── 3b. 运行已有测试(existing_tests_still_pass)
   │   └── 若已有测试失败 → 标记为"破窗", 终止
   ├── 3c. 运行生成测试(test_pass_rate)
   │   └── 重复3次计算flaky_test_rate
   ├── 3d. 覆盖率分析(line_coverage, branch_coverage, coverage_improvement)
   └── 3e. 变异测试(mutation_score, mutation_score_at_pass)

4. 质量评估
   ├── 断言质量分析(assertion_quality)
   ├── LLM-as-Judge可读性评估(test_readability)
   └── 效率指标计算(cost_per_coverage_point, tokens_per_test_case)

5. 报告输出
   ├── JSON详细结果(per-instance)
   ├── HTML可视化报告
   ├── 横向对比图表(多模型)
   └── 结构化trajectory(兼容SWE-agent格式)
```

### 4.6 分阶段实施路线图

#### Phase 1: 最小可用版(1-2周)
- [ ] 扩展repo_level meta.json schema(5维度)
- [ ] 3-5个Python仓库手动构建Docker镜像+meta.json
- [ ] 实现must_pass_existing验证
- [ ] 源码只读挂载
- [ ] 交互轮次限制(50轮)

**交付物**：5个Python仓库的repo_level数据集 + 可运行的评测pipeline

#### Phase 2: 自动化扩展(2-4周)
- [ ] Dockerfile模板注册表(4语言)
- [ ] 仓库自动验证脚本(快速验证)
- [ ] 目标函数自动提取(静态分析)
- [ ] LLM辅助环境修正(SWE-Bench++方案)
- [ ] 环境确定性验证(3次构建)
- [ ] 扩展到30-50个仓库

**交付物**：半自动化数据集构建pipeline

#### Phase 3: 完整评测体系(1-2月)
- [ ] 覆盖率提升指标
- [ ] mutation_score_at_pass指标
- [ ] Flaky Test Rate指标
- [ ] 沙盒回退模式(RepoST风格)
- [ ] 结构化trajectory输出
- [ ] pass_rate_n多轮通过率
- [ ] 扩展到100+仓库, 4语言
- [ ] 防污染机制(时间过滤)

**交付物**：完整的仓库级UT生成评测benchmark

#### Phase 4: 差异化壁垒(2-3月)
- [ ] LLM-as-Judge轨迹质量评估
- [ ] 公开Leaderboard
- [ ] 黄金测试自动收集(PR mining)
- [ ] SWE-Bench++式AutoQA(4层)
- [ ] 持续采集pipeline(SWE-rebench风格)
- [ ] 论文发表

---

## 五、与现有UT-Bench差距的映射

根据`UTBENCH_GAP_ANALYSIS.md`中识别的差距，以下是本次调研结果对应的解决方案：

| 差距ID | 差距描述 | 行业最佳实践 | 具体实施方案 |
|--------|---------|-------------|-------------|
| D1 | 无FAIL/PASS契约 | SWE-bench FAIL_TO_PASS/PASS_TO_PASS | 实现`existing_tests_still_pass`指标 |
| D2 | 无base_commit锁定 | SWE-bench base_commit | meta.json增加source.base_commit |
| D3 | 无环境锁定 | SWE-bench environment_setup_commit | Dockerfile模板+pinned_deps |
| D4 | repo_level meta太薄 | FEA-Bench/SWE-bench schema | 扩展5维度meta.json |
| D5 | 样本量太少 | SWE-Bench++ 11K实例 | 自动化pipeline扩展到100+仓库 |
| D6 | 无难度分级 | SWE-Bench++难度分布 | 增加difficulty维度 |
| D7 | 无污染防护 | SWE-rebench持续采集 | 时间过滤+持续更新机制 |
| D8 | 源码可被修改 | SWE-agent只读挂载 | `-v src:/workspace/src:ro` |
| D18 | 无must_pass_existing | SWE-bench PASS_TO_PASS | Phase 1实现 |
| D23 | 评测无重试 | 行业通用 | 增加重试(1次) |

---

## 六、关键参考资料

1. **SWE-bench** - https://www.swebench.com/ | https://github.com/SWE-bench/SWE-bench
2. **SWE-Bench++** - https://arxiv.org/html/2512.17419v1 (自动合成Docker+多语言+AutoQA)
3. **SWE-rebench** - https://arxiv.org/abs/2505.20411 (21K+任务, 防污染)
4. **TestGenEval** - https://github.com/facebookresearch/testgeneval (最接近UT-Bench, 变异得分@pass)
5. **FEA-Bench** - https://github.com/microsoft/FEA-Bench (仓库选择+新组件识别)
6. **ExecRepoBench** - https://execrepobench.github.io/ (可执行验证)
7. **CrossCodeEval** - https://crosscodeeval.github.io/ (跨文件依赖提取)
8. **RepoBench** - https://github.com/Leolty/repobench (跨文件上下文检索)
9. **RepoST** - https://arxiv.org/abs/2503.07358 (沙盒隔离, 可扩展)
10. **TestPilot** - https://github.com/githubnext/testpilot (迭代精化循环)
11. **LLM变异测试** - https://arxiv.org/html/2406.09843v3 (LLM vs 规则变异测试)
12. **EvalPlus** - https://github.com/evalplus/evalplus (测试增强+进程隔离)
13. **Aider** - https://aider.chat/docs/benchmarks.html (2次机会评测)
