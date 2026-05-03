# UT-Bench Agent 评测升级 —— 开源生态调研与优化建议

> 调研时间：2026-05-01
> 调研范围：Agent 沙箱、Agent 框架、评测基准、Skill 系统、评测指标体系

---

## 一、Agent 沙箱方案调研

### 1.1 沙箱技术全景对比

| 技术 | 隔离机制 | 启动时间 | 内存开销 | 安全等级 | 适用场景 |
|:---|:---|:---|:---|:---|:---|
| **Docker (runc)** | Linux namespaces + cgroups + seccomp | ~10ms | ~10MB | ⭐ 弱 | 可信代码，开发调试 |
| **gVisor** | 用户空间内核（Sentry）拦截所有 syscall | ~100ms | ~20MB | ⭐⭐⭐ 强 | 计算密集型 Agent 任务 |
| **Firecracker** | 硬件虚拟化（KVM），独立内核 | ~125ms | ~5MB | ⭐⭐⭐⭐⭐ 最强 | 不可信代码，生产环境 |
| **Kata Containers** | 多 VMM 后端 + OCI 容器 API | ~200ms | ~30MB | ⭐⭐⭐⭐⭐ 最强 | K8s 环境的 Agent 部署 |
| **WebAssembly (WASI)** | 能力模型，默认拒绝 | <1ms | <1MB | ⭐⭐⭐⭐ 强 | 有界计算，浏览器执行 |
| **Cloudflare V8 Isolates** | V8 沙箱 | <1ms | ~1MB | ⭐⭐⭐⭐ 强 | 高频短生命周期调用 |

> **2026 行业共识**：Docker/runc 对于 AI 生成代码根本不足（fundamentally insufficient）。

### 1.2 各方案详细分析

#### Docker（当前 UT-Bench 第一版方案）

**现状**：UT-Bench 当前使用 Docker 内层容器作为每样本沙箱，通过 DOOD 方式控制。

**优点**：
- 启动快，生态成熟
- 与现有 CI/CD 流程兼容
- 配置简单

**缺点**：
- 共享主机内核，容器逃逸风险
- seccomp 配置容易出错
- 不适合不可信代码（AI 生成的代码就是不可信的）

**结论**：可作为第一版过渡方案，但不应作为最终方案。

#### gVisor

**原理**：Google 开发的用户空间内核，拦截所有系统调用，不让其直接到达主机内核。

**优点**：
- 系统调用从不直接到达主机内核
- 启动时间适中（~100ms）
- 计算密集型任务开销极小
- 支持 Systrap 和 KVM 两种模式

**缺点**：
- I/O 密集型有 10-30% 性能开销
- 需要额外的 Sentry 进程

**对 UT-Bench 的价值**：
- 可直接作为 Docker runtime 替换：`docker run --runtime=runsc`
- **最小改动即可升级安全等级**
- 适合编译、测试这类 I/O 中等的任务

**实施路径**：
```bash
# 宿主机安装 gVisor
wget https://gvisor.dev/releases/release/latest/x86_64/runsc
sudo mv runsc /usr/local/bin/

# Docker 配置 /etc/docker/daemon.json
{
  "runtimes": {
    "runsc": {
      "path": "/usr/local/bin/runsc"
    }
  }
}

# UT-Bench 内层容器使用 gVisor
sandbox_mode: docker
docker_runtime: runsc  # 新增配置
```

#### Firecracker MicroVM

**原理**：AWS 开源的轻量级虚拟化，每个工作负载运行自己的 Linux 内核。

**优点**：
- 硬件级隔离，独立内核
- 启动极快：~125ms
- 内存开销极低：<5 MiB/VM
- 高密度：每秒 150 个 microVM/主机
- Firecracker 进程本身仅允许 24 个系统调用

**缺点**：
- 需要 KVM 支持
- 云环境（如某些容器服务）可能受限

**对 UT-Bench 的价值**：
- 最高安全等级
- 每个样本一个 microVM，完全隔离
- 适合最终生产环境

**实施路径**：
- 使用 `firecracker-containerd` 或 `weaveworks/ignite`
- 或直接用 `firectl` 管理 microVM

#### 腾讯云 Cube Sandbox（2026-04 刚开源）

**关键信息**：
- 2026 年 4 月 21 日腾讯云开源
- 基于 KVM + RustVMM 技术栈
- 硬件级隔离
- **Drop-in 兼容 E2B 接口**
- 冷启动 <60ms，单实例内存 <5MB
- 单机支持 2000+ 实例
- Apache 2.0 协议

**对 UT-Bench 的价值**：
- **国内部署首选**，网络延迟低
- 兼容 E2B 意味着我们的 `SandboxRunner` 接口可以平滑接入
- 如果最终要上生产，这是比 E2B 更可控的选择

**建议**：持续关注其 GitHub 仓库开源进展，作为后续沙箱后端候选。

#### E2B

**原理**：托管沙箱服务，为 AI Agent 提供安全代码执行环境。

**核心 API**：
```python
from e2b import Sandbox
sandbox = Sandbox.create()
result = sandbox.commands.run('pytest test_file.py')
print(result.stdout)
```

**优点**：
- API 简洁，SDK 完善（JS/TS + Python）
- 无需自己维护基础设施
- 支持模板（Template）预配置环境

**缺点**：
- 托管服务，有费用
- 需要网络访问 E2B 服务
- 数据隐私考虑

**对 UT-Bench 的价值**：
- 作为 `SandboxRunner` 的一个后端实现
- 适合快速验证，不适合大规模批量评测（成本）

### 1.3 沙箱方案演进建议

```
当前（第一版）        短期（2-4 周）         中期（1-2 月）          长期（3-6 月）
    │                    │                    │                     │
    ▼                    ▼                    ▼                     ▼
┌─────────┐        ┌─────────┐         ┌─────────────┐        ┌─────────────┐
│ Docker  │   →    │ Docker  │    →    │ gVisor      │   →    │ Firecracker │
│ (runc)  │        │ + gVisor│         │ (默认)       │        │ / Cube      │
│ DOOD    │        │ (可选)  │         │             │        │ / E2B       │
└─────────┘        └─────────┘         └─────────────┘        └─────────────┘
```

**具体建议**：
1. **现在**：保持 Docker，但把 `docker_runtime` 做成可配置项
2. **两周内**：支持 `--runtime=runsc`，让用户一键启用 gVisor
3. **一月内**：调研 Firecracker 在本地部署的可行性
4. **长期**：Cube Sandbox 开源后，作为国内首选后端

---

## 二、Agent 框架调研

### 2.1 开源 Agent 框架对比

| 框架 | 类型 | 模型切换 | 非交互式 | 轨迹输出 | Skill 系统 | 沙箱集成 |
|:---|:---|:---|:---|:---|:---|:---|
| **OpenCode** | CLI | ✅ 75+ 模型 | ✅ `run` 模式 | ⚠️ 日志 | ✅ Skills | ⚠️ 需配置 |
| **Claude Code** | CLI | ⚠️ Anthropic 为主 | ✅ `-p` 模式 | ⚠️ 日志 | ✅ Skills | ⚠️ 需配置 |
| **Aider** | CLI | ✅ 多模型 | ✅ 非交互 | ✅ benchmark | ⚠️ 配置化 | ⚠️ 需配置 |
| **SWE-agent** | 框架 | ✅ 可配置 | ✅ 自动化 | ✅ trajectory | ⚠️ 配置 | ✅ 内置 Docker |
| **OpenHands** | 框架 | ✅ 可配置 | ✅ 自动化 | ✅ EventStream | ✅ 插件 | ✅ 内置多后端 |

### 2.2 各框架详细分析

#### OpenCode（UT-Bench 已接入）

**状态**：已作为第一个真实 framework 接入。

**验证结果**：
- ✅ 非交互 `run` 模式可用
- ✅ `--model provider/model` 支持任意底层模型
- ✅ `OPENCODE_CONFIG_CONTENT` 可动态配置 provider
- ✅ 能生成测试文件到 workspace
- ⚠️ Token 消耗无法直接获取（stderr 只有日志，没有结构化数据）
- ⚠️ Agent 过程轨迹不够详细（缺少工具调用序列）

**优化建议**：
1. 尝试用 `--output-last-message` 或 `--format json` 获取结构化输出
2. 如果 OpenCode 支持 MCP，可以通过 MCP 协议获取更详细的工具调用记录
3. 考虑用 wrapper 脚本包装 `opencode run`，拦截 stdout/stderr 做更详细的解析

#### Claude Code

**特点**：
- `-p` 参数支持非交互式单轮执行
- Skills 系统完善：`.claude/skills/` 目录 + `SKILL.md` 格式
- 系统提示追加能力强
- 但模型切换主要面向 Anthropic 自身模型

**接入 UT-Bench 的路径**：
```bash
# 非交互执行
claude -p "$(cat prompt.md)" --allowed-tools "Edit,Bash"
```

**挑战**：
- Claude Code 不一定支持非 Claude 模型（需要验证）
- 如果只能跑 Claude 模型，那只能做"原生方案对比"，不能做"同模型框架对比"

**建议**：
1. 先验证 Claude Code 是否能挂 OpenRouter 或其他 provider
2. 如果不能，把它定位为"Claude 原生方案"被测对象
3. 如果能，它就可以进入"同模型框架对比组"

#### Aider

**特点**：
- 非常成熟的 benchmark 体系（`aider/benchmark`）
- 基于 Exercism 编程练习题
- 完全非交互式，支持 `--threads` 并行
- 强制 Docker 运行（安全考虑）
- 详细的评分指标：pass_rate、语法错误、缩进错误、上下文窗口耗尽等

**对 UT-Bench 的价值**：
- Aider 的 benchmark 设计是**最佳实践参考**
- 它的评分指标非常细，可以借鉴
- 它的"多轮通过率"（pass_rate_1 / pass_rate_2）概念值得引入

**Aider Benchmark 核心指标**：
```yaml
pass_rate_1: 0.75        # 一次尝试通过率
pass_rate_2: 0.82        # 两次尝试通过率
percent_cases_well_formed: 0.95
error_outputs: 3
num_malformed_responses: 2
syntax_errors: 1
indentation_errors: 0
exhausted_context_windows: 0
test_timeouts: 0
seconds_per_case: 12.5
total_cost: 0.42
```

**建议**：把 Aider 作为第二个接入的 CLI Agent 框架，同时学习它的 benchmark 指标设计。

#### SWE-agent

**特点**：
- 专门用于代码修复（issue resolution）
- **SWE-ReX**：独立的沙箱抽象层，支持 Docker/本地/远程
- 原生产出 trajectory（完整的工具调用序列）
- 与 SWE-bench 深度集成

**SWE-ReX 设计**（对 UT-Bench 非常有价值）：
```python
# SWE-ReX 的核心抽象
class Sandbox:
    def run_command(self, cmd: str) -> CommandResult
    def read_file(self, path: str) -> str
    def write_file(self, path: str, content: str)
    def close(self)
```

**对 UT-Bench 的价值**：
- SWE-ReX 的抽象和我们的 `SandboxRunner` 非常相似
- 可以借鉴它的"沙箱作为独立服务"的设计理念
- trajectory 格式可以作为我们 Agent 过程采集的参考

**建议**：
1. 研究 SWE-ReX 的接口设计，丰富我们的 `SandboxRunner`
2. 学习 SWE-agent 的 trajectory 格式，改进我们的 `trace.jsonl`

#### OpenHands

**特点**：
- 最完整的 Agent 框架之一
- Runtime 抽象最完善：Docker / Local / Remote / E2B
- Action/Observation 事件流模型
- 插件系统：Jupyter、VSCode、Agent Skills

**Runtime 架构**（对 UT-Bench 极具参考价值）：
```
Agent ←→ EventStream ←→ ActionExecutor (REST API)
                         └─ Docker Container
                             └─ ActionExecutor Server
                                 ├─ Browser
                                 ├─ Bash Shell
                                 ├─ Jupyter
                                 └─ Plugins
```

**对 UT-Bench 的价值**：
- **Runtime 多后端设计**是我们 `SandboxRunner` 的最佳参考
- EventStream 模型可以用来改进我们的 trace 采集
- 插件系统可以作为 Skill 系统的长期参考

**建议**：
1. 参考 OpenHands 的 Runtime 接口，丰富 `SandboxRunner`
2. 学习它的 EventStream 设计，把 Agent 过程采集做得更完整

### 2.3 Agent 框架接入优先级

```
已接入：     OpenCode
下一个：     Aider（benchmark 体系成熟，值得学习）
再下一个：   Claude Code（如果模型切换可行）
长期：       SWE-agent（trajectory 丰富）
            OpenHands（Runtime 抽象最完善）
```

---

## 三、评测基准与数据集调研

### 3.1 现有基准对比

| 基准 | 任务类型 | 语言 | 样本量 | Agent 友好 | 对 UT-Bench 的价值 |
|:---|:---|:---|:---|:---|:---|
| **SWE-bench** | Issue 修复 | Python/JS | ~2K | ✅ 专为 Agent 设计 | 仓库级样本来源 |
| **BigCodeBench** | 代码生成 | Python/JS/Java/Go | 1,140 | ⚠️ 偏模型 | 高质量编程任务 |
| **Exercism** | 编程练习 | 多语言 | ~5K | ✅ 有测试验证 | Aider 使用，可借鉴 |
| **HumanEval** | 函数生成 | Python | 164 | ❌ 太简单 | 历史基准，不适用 |
| **MBPP** | 函数生成 | Python | ~1K | ❌ 太简单 | 历史基准，不适用 |

### 3.2 SWE-bench

**核心设计**：
- 从真实 GitHub Issues 提取任务
- 每个任务 = 问题描述 + 代码仓库 + 测试验证
- Docker-based evaluation harness
- 自动应用 patch 并运行测试验证

**对 UT-Bench 的价值**：
- **仓库级样本的最佳来源**
- 但 SWE-bench 的任务是"修复 bug"，不是"生成单元测试"
- 需要适配：把"修复任务"改成"为这段代码生成单元测试"

**适配思路**：
```
SWE-bench 原始任务：
  "修复这个 issue：当输入为空列表时程序崩溃"

UT-Bench 适配后：
  "为这段代码（修复前的版本）生成单元测试，
   测试应该能发现'输入为空列表时程序崩溃'这个 bug"
```

### 3.3 BigCodeBench

**核心设计**：
- 1,140 个实际编程任务
- 覆盖 Python、JavaScript、Java、Go
- 每个任务有完整的函数签名、docstring、测试用例
- 任务难度分布均匀

**对 UT-Bench 的价值**：
- 任务质量高，比 HumanEval/MBPP 更实用
- 多语言覆盖，与 UT-Bench 的语言策略一致
- 但它是"从头写函数"，不是"为已有代码写测试"

**适配思路**：
- 把 BigCodeBench 的"参考实现"作为被测代码
- 把"测试用例"作为期望测试的参考（但不直接暴露给 Agent）
- 评测 Agent 生成的测试是否能覆盖参考实现的行为

### 3.4 自建仓库级样本

**建议方向**：

```
repos/
├── from_swe_bench/           # 从 SWE-bench 适配
│   └── python_flask_001/
│       ├── src/               # 被测源码（问题版本）
│       ├── task.md            # 任务描述
│       ├── setup.py           # 环境安装
│       └── expected_tests/    # 参考测试（隐藏）
│
├── from_bigcodebench/        # 从 BigCodeBench 适配
│   └── python_data_process_001/
│       ├── src/
│       ├── task.md
│       └── expected_tests/
│
└── from_github/              # 从 GitHub 热门项目提取
    └── python_requests_001/
        ├── src/
        ├── task.md
        └── expected_tests/
```

**每个样本的元数据**：
```json
{
  "sample_id": "python_flask_001",
  "language": "python",
  "source": "swe_bench",
  "task_type": "generate_tests_for_module",
  "target_module": "flask/app.py",
  "task_description": "为 Flask 的 app.py 模块生成单元测试...",
  "setup_script": "setup.sh",
  "expected_test_count": 15,
  "difficulty": "medium",
  "lines_of_code": 450
}
```

---

## 四、Skill 系统设计调研

### 4.1 各框架 Skill 系统对比

| 框架 | Skill 格式 | 安装方式 | 注入方式 | 触发方式 |
|:---|:---|:---|:---|:---|
| **OpenCode** | `SKILL.md` (YAML + Markdown) | GitHub/手动/市场 | 自动注入 prompt | 自动识别 / `@` 手动 |
| **Claude Code** | `SKILL.md` (YAML + Markdown) | `.claude/skills/` 目录 | 自动注入 system prompt | 自动识别 |
| **UT-Bench 当前** | `*.md` 文件 | 配置文件引用 | `prompt_append` / `workspace_mount` | 配置绑定 |

### 4.2 OpenCode / Claude Code Skill 设计

**统一格式**：
```markdown
---
name: unit-test-generator
version: "1.0"
description: Generate comprehensive unit tests for Python modules
---

# Unit Test Generator Skill

## Instructions
When generating unit tests, follow these principles:

1. Use pytest framework
2. Cover normal, boundary, and error paths
3. Mock external dependencies
4. Keep tests deterministic
5. Use descriptive test names

## Examples

### Example 1: Simple function
```python
def test_add_basic():
    assert add(2, 3) == 5

def test_add_boundary():
    assert add(0, 0) == 0
    assert add(-1, 1) == 0
```

## Constraints
- Do not modify the source code being tested
- Import all dependencies explicitly
- Use `pytest.raises` for exception testing
```

**对 UT-Bench 的 Skill 系统的启示**：
1. **标准化格式**：采用 `SKILL.md`（YAML frontmatter + Markdown body）
2. **版本管理**：Skill 应该有 version 字段
3. **示例驱动**：Skill 里应该包含具体示例，不只是抽象指令
4. **约束明确**：明确告诉 Agent 什么不能做

### 4.3 UT-Bench Skill 系统优化建议

**当前问题**：
- `unit_test_skill.md` 只有抽象指令，缺少具体示例
- 没有版本管理
- 注入方式只有 `prompt_append` 和 `workspace_mount`

**优化方向**：

```yaml
# skills/unit_test_skill_v2.md
---
name: unit_test_expert
version: "2.0"
description: Expert-level unit test generation with coverage optimization
author: UT-Bench Team
compatible_frameworks: [opencode, claude_code, aider]
compatible_languages: [python, go, java, cpp]
inject_mode: prompt_append
---

# Unit Test Expert

## Core Principles
1. **Coverage-first**: Aim for >90% line coverage
2. **Mutation-aware**: Write tests that catch common mutation patterns
3. **Deterministic**: Never use random data without fixed seeds

## Language-Specific Guidelines

### Python
- Use `pytest` with `parametrize` for boundary cases
- Mock at the import point: `@patch('module_under_test.requests.get')`
- Use `tmp_path` fixture for file I/O tests

### Go
- Use table-driven tests
- Test both success and error paths with `errors.Is`
- Use `testify/assert` for complex assertions

## Examples
[具体示例...]

## Anti-patterns to Avoid
- ❌ Don't test implementation details (private methods)
- ❌ Don't use `time.sleep()` in tests
- ❌ Don't call real external APIs
```

**新增注入模式**：
```yaml
inject_modes:
  - prompt_append      # 当前已有
  - workspace_mount    # 当前已有
  - agent_native       # 预留
  - system_prompt      # 追加到 Agent system prompt（如果 Agent 支持）
  - pre_task_script    # 在执行前运行脚本（安装依赖、准备环境）
  - post_task_script   # 在执行后运行脚本（格式化、lint）
```

---

## 五、评测指标体系调研

### 5.1 现有评测维度

UT-Bench 当前已有：
- 编译通过率
- 测试通过率
- 行覆盖率
- 分支覆盖率
- 变异得分
- 耗时
- Token 消耗（仅 model_api）

### 5.2 开源项目评测指标参考

#### Aider Benchmark 指标

| 指标 | 说明 |
|:---|:---|
| `pass_rate_1` | 一次尝试通过率 |
| `pass_rate_2` | 两次尝试通过率 |
| `percent_cases_well_formed` | 格式良好的案例比例 |
| `syntax_errors` | 语法错误数 |
| `indentation_errors` | 缩进错误数 |
| `exhausted_context_windows` | 上下文窗口耗尽次数 |
| `seconds_per_case` | 每案例平均耗时 |
| `total_cost` | 总 API 费用 |

#### SWE-eval 指标（Trajectory-Enhanced）

| 维度 | 指标 | 说明 |
|:---|:---|:---|
| **效率** | 资源消耗 | 完成任务所需的 token/时间/费用 |
| **逻辑一致性** | Intra-turns | 单轮对话内的逻辑一致性 |
| | Inter-turns | 多轮对话间的逻辑一致性 |
| **工具使用** | Info-gain | 工具为解决问题提供的新信息量 |

### 5.3 UT-Bench 应新增的 Agent 过程指标

**第一层：基础过程指标**

| 指标 | 采集方式 | 用途 |
|:---|:---|:---|
| `agent_duration_ms` | runner 计时 | Agent 执行总耗时 |
| `agent_command_count` | 解析 trace | Agent 执行的命令数量 |
| `agent_file_reads` | 解析 trace / diff | 读取文件次数 |
| `agent_file_writes` | 解析 diff | 写入/修改文件次数 |
| `agent_rounds` | 框架回传 | 多轮交互轮次 |

**第二层：Token & 成本指标**

| 指标 | 采集方式 | 用途 |
|:---|:---|:---|
| `agent_input_tokens` | 框架回传 / 估算 | Agent 输入 token |
| `agent_output_tokens` | 框架回传 / 估算 | Agent 输出 token |
| `agent_total_tokens` | 框架回传 / 估算 | Agent 总 token |
| `agent_estimated_cost` | 按模型定价计算 | 估算费用 |
| `tokens_per_line_of_test` | 计算 | Token 效率 |
| `cost_per_coverage_point` | 计算 | 成本效率 |

**第三层：工具使用精准度**

| 指标 | 采集方式 | 用途 |
|:---|:---|:---|
| `tool_read_file_count` | 解析 trace | 读文件工具调用次数 |
| `tool_write_file_count` | 解析 trace | 写文件工具调用次数 |
| `tool_execute_command_count` | 解析 trace | 执行命令次数 |
| `tool_search_count` | 解析 trace | 搜索/查找次数 |
| `tool_usage_accuracy` | 人工标注 / LLM 评估 | 工具使用是否合理 |

**第四层：质量与失败分析**

| 指标 | 采集方式 | 用途 |
|:---|:---|:---|
| `failure_stage` | runner 判断 | 失败发生在哪一阶段 |
| `failure_category` | runner 判断 | 失败类型分类 |
| `retry_count` | runner 记录 | Agent 重试次数 |
| `abnormal_termination` | runner 记录 | 是否异常终止 |

### 5.4 失败阶段分类

```
generation_pipeline:
  1. prompt_render      → fail: prompt_error
  2. workspace_prepare  → fail: workspace_error
  3. skill_inject       → fail: skill_error
  4. agent_execute      → fail: agent_timeout / agent_crash / agent_output_error
  5. test_file_collect  → fail: no_test_file / malformed_test
  6. compile            → fail: compile_error
  7. test_run           → fail: test_failure
  8. coverage           → fail: coverage_error
  9. mutation           → fail: mutation_error
```

---

## 六、对 UT-Bench 的全面优化建议

### 6.1 架构层面优化

#### 1. SandboxRunner 后端扩展

**当前**：`local` / `docker`

**建议扩展**：
```go
type SandboxBackend string

const (
    BackendLocal          SandboxBackend = "local"
    BackendDocker         SandboxBackend = "docker"
    BackendDockerGVisor   SandboxBackend = "docker+gvisor"
    BackendFirecracker    SandboxBackend = "firecracker"
    BackendE2B            SandboxBackend = "e2b"
    BackendCube           SandboxBackend = "cube"
)
```

**实施优先级**：
1. `docker+gvisor`：最小改动，最大安全提升
2. `e2b`：验证可行性，适合小规模测试
3. `firecracker`：长期目标，最高安全等级
4. `cube`：等腾讯云开源后接入

#### 2. Agent 过程采集增强

**当前**：只有 `trace.jsonl`（命令、stdout/stderr、exit code、duration）

**建议增强**：
```json
{
  "trace_version": "v2",
  "events": [
    {
      "timestamp": "2026-05-01T08:15:00Z",
      "type": "command_start",
      "command": "opencode run ...",
      "cwd": "/workspace"
    },
    {
      "timestamp": "2026-05-01T08:15:05Z",
      "type": "tool_call",
      "tool": "read_file",
      "args": {"path": "boundary_000.py"},
      "result_summary": "read 45 lines"
    },
    {
      "timestamp": "2026-05-01T08:15:10Z",
      "type": "tool_call",
      "tool": "write_file",
      "args": {"path": "generated_test.py"},
      "result_summary": "wrote 120 lines"
    },
    {
      "timestamp": "2026-05-01T08:15:30Z",
      "type": "llm_request",
      "model": "deepseek-v4-flash",
      "input_tokens": 2048,
      "output_tokens": 1024,
      "cost_usd": 0.0015
    },
    {
      "timestamp": "2026-05-01T08:15:45Z",
      "type": "command_end",
      "exit_code": 0,
      "duration_ms": 45000
    }
  ],
  "summary": {
    "total_commands": 1,
    "total_tool_calls": 5,
    "total_llm_requests": 3,
    "total_input_tokens": 6144,
    "total_output_tokens": 3072,
    "total_cost_usd": 0.0045
  }
}
```

**实现路径**：
- 短期：解析 Agent 的 stdout，用正则提取工具调用
- 中期：如果 Agent 支持 MCP，通过 MCP 协议获取结构化事件
- 长期：让 Agent 框架原生支持结构化事件输出

#### 3. 配置系统插件化

**当前**：framework 配置在 `agents.yaml` 里写死

**建议**：
```yaml
# configs/agents.yaml
frameworks:
  opencode:
    plugin: cli_agent  # 引用内置插件
    config:
      command: "opencode run ..."
      docker_image: "utbench-agent-opencode:latest"

  # 未来用户可以自己添加，不改动代码
  my_custom_agent:
    plugin: cli_agent
    config:
      command: "my-agent --task {{.PromptFile}}"
      docker_image: "my-agent:latest"
```

### 6.2 评测层面优化

#### 1. 新增 Agent 过程维度

**立即可以做的**（不依赖 Agent 框架配合）：
- `agent_duration_ms`：runner 已有数据
- `agent_file_writes`：从 diff.json 统计
- `failure_stage`：runner 已有判断逻辑

**需要 Agent 框架配合的**：
- `agent_input_tokens` / `agent_output_tokens`
- `agent_command_count`
- `tool_call` 明细

**建议**：先实现"不依赖 Agent 配合"的指标，再逐步推进"需要配合"的指标。

#### 2. 引入 Aider 的 pass_rate 概念

**当前**：一次生成，成功/失败

**建议**：支持多次尝试
```yaml
# agents.yaml
frameworks:
  opencode:
    retry_policy:
      max_retries: 2
      retry_on: [test_failure, compile_error]
```

报告增加：
```json
{
  "pass_rate_1": 0.65,
  "pass_rate_2": 0.78,
  "pass_rate_3": 0.82
}
```

#### 3. 引入 SWE-eval 的 Trajectory 评估

**建议**：用 LLM 评估 Agent 的 trajectory 质量

```python
# 伪代码
def evaluate_trajectory(trace_path):
    trace = load_trace(trace_path)
    
    # 效率评估
    efficiency_score = evaluate_efficiency(trace)
    
    # 逻辑一致性评估
    logical_consistency = evaluate_logical_consistency(trace)
    
    # 工具使用评估
    tool_utilization = evaluate_tool_utilization(trace)
    
    return {
        "efficiency_score": efficiency_score,
        "logical_consistency": logical_consistency,
        "tool_utilization": tool_utilization,
        "info_gain": calculate_info_gain(trace)
    }
```

### 6.3 数据集层面优化

#### 1. 仓库级样本接入路线图

```
Phase 1（现在-2周）：    完善单文件样本包装
  - 统一最小 repo 结构
  - 每种语言预配置测试框架

Phase 2（2-4周）：        接入 SWE-bench 子集
  - 选取 50-100 个 Python 样本
  - 适配"生成测试"任务格式
  - 验证可行性

Phase 3（1-2月）：         接入 BigCodeBench 子集
  - 选取多语言样本
  - 适配测试生成任务

Phase 4（2-3月）：         自建 GitHub 样本
  - 从热门项目提取模块
  - 人工标注任务描述
```

#### 2. 最小 repo 结构标准化

```
workspace/
├── src/
│   └── {sample_id}.py          # 被测源码
├── tests/
│   └── __init__.py             # 预创建
├── conftest.py                 # pytest 配置（如有需要）
├── pyproject.toml              # 项目配置
├── requirements.txt            # 依赖
├── .utbench/
│   └── skills/
│       └── {skill_name}.md     # skill 文件
└── utbench_agent_prompt.md     # 执行契约
```

### 6.4 Skill 层面优化

#### 1. Skill 格式升级

采用 `SKILL.md` 标准格式（与 OpenCode / Claude Code 兼容）：

```markdown
---
name: {skill_id}
version: "{version}"
description: {description}
author: {author}
compatible_frameworks: [opencode, claude_code, aider]
compatible_languages: [python, go, java, cpp]
inject_mode: prompt_append
---

# {Skill Name}

## Instructions
{具体指令}

## Examples
{具体示例}

## Constraints
{明确约束}
```

#### 2. Skill 市场/仓库

```
skills/
├── official/
│   ├── unit_test_basic/
│   │   └── SKILL.md
│   ├── unit_test_advanced/
│   │   └── SKILL.md
│   ├── mutation_test_aware/
│   │   └── SKILL.md
│   └── mocking_best_practices/
│       └── SKILL.md
└── community/
    └── ...
```

### 6.5 报告层面优化

#### 1. 增强 Agent vs Model API 对比报告

```json
{
  "agent_comparisons": [
    {
      "subject_id": "opencode__deepseek-v4-flash__no_skill",
      "baseline_subject_id": "model_api__deepseek-v4-flash__no_skill",
      "model": "deepseek-v4-flash",
      "language": "python",
      "sample_count": 100,
      
      "compile_pass": {
        "baseline": 0.95,
        "agent": 0.97,
        "delta": 0.02
      },
      "test_pass": {
        "baseline": 0.78,
        "agent": 0.85,
        "delta": 0.07
      },
      "line_coverage": {
        "baseline": 0.72,
        "agent": 0.84,
        "delta": 0.12
      },
      "mutation_score": {
        "baseline": 0.55,
        "agent": 0.68,
        "delta": 0.13
      },
      
      "efficiency": {
        "avg_latency_ms": {
          "baseline": 15000,
          "agent": 120000,
          "delta": 105000
        },
        "avg_tokens": {
          "baseline": 2500,
          "agent": 8000,
          "delta": 5500
        },
        "avg_cost_usd": {
          "baseline": 0.005,
          "agent": 0.015,
          "delta": 0.01
        }
      },
      
      "process": {
        "avg_file_reads": {"agent": 3.5},
        "avg_file_writes": {"agent": 1.2},
        "avg_tool_calls": {"agent": 8.3}
      }
    }
  ]
}
```

#### 2. 新增 Trajectory 可视化

**建议输出**：
- `report/trajectories/{subject}/{sample}.html`：单个样本的 trajectory 时间线
- `report/comparison.html`：两个 subject 的并排对比
- `report/dashboard.html`：总览仪表盘

### 6.6 实施优先级矩阵

| 优化项 | 影响 | 工作量 | 优先级 | 建议时间 |
|:---|:---|:---|:---|:---|
| **gVisor 支持** | 高（安全提升） | 低 | **P0** | 1 周内 |
| **Agent 基础过程指标**（duration, file writes） | 高 | 低 | **P0** | 1 周内 |
| **Skill 格式标准化** | 中 | 低 | **P1** | 2 周内 |
| **Aider framework 接入** | 高（多框架对比） | 中 | **P1** | 2-3 周 |
| **Token 消耗采集** | 高 | 中 | **P1** | 2-3 周 |
| **失败阶段细化** | 中 | 低 | **P1** | 1-2 周 |
| **仓库级样本（SWE-bench 子集）** | 高 | 高 | **P2** | 1 月 |
| **Trajectory 评估（LLM-based）** | 高 | 高 | **P2** | 1-2 月 |
| **Firecracker 支持** | 高（安全最高） | 高 | **P2** | 2-3 月 |
| **报告可视化增强** | 中 | 中 | **P3** | 2-3 月 |
| **Cube Sandbox 接入** | 高（国内生产） | 中 | **P3** | 等开源 |

---

## 七、可直接使用的开源项目/代码

### 7.1 可直接集成

| 项目 | 用途 | 集成方式 |
|:---|:---|:---|
| **gVisor** | 沙箱安全升级 | Docker runtime 替换 |
| **SWE-ReX** | 沙箱抽象层参考 | 学习接口设计 |
| **E2B SDK** | 托管沙箱后端 | `SandboxRunner` 后端实现 |
| **Aider Benchmark** | 评测指标参考 | 学习指标设计 |

### 7.2 可参考学习

| 项目 | 学习点 |
|:---|:---|
| **OpenHands Runtime** | 多后端 Runtime 架构、EventStream 模型 |
| **SWE-eval** | Trajectory 评估框架、三维评测指标 |
| **Claude Code Skills** | Skill 格式设计、系统提示注入 |
| **OpenCode Skills** | Skill 市场设计、插件安装机制 |

### 7.3 值得关注

| 项目 | 关注原因 |
|:---|:---|
| **腾讯云 Cube Sandbox** | 2026-04 刚开源，国内部署首选，兼容 E2B |
| **kubernetes-sigs/agent-sandbox** | K8s 官方 Agent 沙箱工作组 |
| **AI Fortress (Flatcar + gVisor)** | 多层沙箱最佳实践 |

---

## 八、总结

UT-Bench 的第一版 Agent 评测升级已经**成功跑通端到端链路**，当前状态：

- ✅ `subject = framework + model + skill` 架构
- ✅ OpenCode 首个真实 framework 接入
- ✅ 每样本独立 Docker 沙箱
- ✅ Agent trace / diff 采集
- ✅ 报告包含 Agent vs Model API 对比

**下一步最关键的 3 件事**：

1. **安全升级（1 周内）**：支持 gVisor 作为 Docker runtime，一行配置即可提升安全等级
2. **过程指标补齐（1-2 周内）**：采集 Agent duration、file writes、failure stage 等基础指标
3. **多框架接入（2-3 周内）**：接入 Aider，实现真正的"框架横评"

**长期愿景**：
- 沙箱：从 Docker → gVisor → Firecracker/Cube
- 框架：从 OpenCode → Aider → Claude Code → SWE-agent → OpenHands
- 样本：从单文件 → 仓库级（SWE-bench + BigCodeBench + GitHub）
- 指标：从结果指标 → 过程指标 → Trajectory 评估
- Skill：从静态 prompt → 标准化 Skill 市场
