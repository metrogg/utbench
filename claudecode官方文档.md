# Claude Code 官方文档梳理与 UT-Bench 接入适配说明

更新日期：2026-05-05

本文目标不是复述 Claude Code 文档，而是把 **接入 UT-Bench 时真正需要的信息** 提炼出来，并明确：

1. Claude Code 官方到底支持哪些接入方式；
2. 它和你现在已经接入的 OpenCode / CodeBuddy 有哪些本质区别；
3. UT-Bench 为了接入 `claudecode` framework 需要改哪些地方；
4. “Claude Code 接各家模型”这件事，哪些是官方支持的，哪些不能想当然。

---

## 1. 先说结论

### 1.1 Claude Code 可以接入，但不要把它理解成“任意 OpenAI 兼容模型前端”

从官方文档看，Claude Code 当前明确支持的模型接入路径主要是：

1. **Anthropic 官方 API**
2. **Amazon Bedrock**
3. **Google Vertex AI**
4. **Microsoft Azure / Foundry（官方文档已有独立章节）**
5. **LLM Gateway / Proxy**  
   通过 `ANTHROPIC_BASE_URL`、`ANTHROPIC_AUTH_TOKEN` / `ANTHROPIC_API_KEY`、`apiKeyHelper` 等方式接到网关

但这里的重点是：

- Claude Code 的主世界观仍然是 **Claude 模型**；
- 即便通过网关接入，官方对模型发现、模型名、能力假设，也更偏向 Claude 体系；
- 不能简单宣传成“Claude Code = 任意模型通用 CLI 外壳”。

对 UT-Bench 的含义是：

- 你可以把 `claudecode` 当成一个 **CLI Agent framework** 接进来；
- 但“借 Claude Code 跑 DeepSeek / Qwen / Minimax / 豆包”等第三方模型，必须走 **网关 / 映射 / 代理** 路线；
- 而且这条路线是否稳定、是否完全受官方支持，要按 **官方文档明确列出的配置方式** 来做，不能按 OpenCode / CodeBuddy 的思路直接拼一个 OpenAI Compatible 配置就假定能跑。

### 1.2 对 UT-Bench 来说，Claude Code 接入本质上是第三个 CLI Agent framework

你现在项目里已经有：

- `model_api`
- `opencode`
- `codebuddy`

Claude Code 应该按同一层级接入：

- `framework = claudecode`
- `kind = cli_agent`
- subject 形式：
  - `claudecode__claude-sonnet-4__no_skill`
  - `claudecode__claude-sonnet-4__unit_test_skill`
  - 如果以后做网关映射，也可能出现：
    - `claudecode__deepseek-v4-flash__no_skill`
    - 但这类 subject 是否官方支持、结果是否可解释，需要单独标记

### 1.3 接入前最重要的判断

你接下来其实有两条路线：

#### 路线 A：先把 Claude Code 作为“官方 Claude 工作流基线”接入

只支持官方明确支持的模型后端：

- Anthropic API
- Bedrock
- Vertex
- Foundry

这是最稳的。

#### 路线 B：再研究 Claude Code 通过 LLM Gateway 跑“各家模型”

这条路可以做，但要提高技术标准：

- 不能只看“能不能跑起来”
- 要确认：
  - 模型发现机制
  - 模型名映射机制
  - 输出 / usage / session 数据是否完整
  - 是否影响 Claude Code 的权限 / 工具 / prompt 工作流假设

建议顺序是 **先 A 后 B**。  
不要一上来就把 Claude Code 当“统一任意模型代理器”，那样后面评测解释会很混乱。

---

## 2. 当前 UT-Bench 与 Claude Code 接入点

基于当前仓库状态，Claude Code 接入不会从零开始。

### 2.1 已有基础

你现在的核心结构已经具备：

- `subject = framework + model + skill`
- `framework.kind = cli_agent`
- `framework.sandbox.*`
- `framework.command`
- `framework.preflight`
- `framework.env`
- `framework.env_from_host`
- `framework.forbidden_command_patterns`

关键文件：

- [go-ut-bench/configs/agents.example.yaml](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/configs/agents.example.yaml)
- [go-ut-bench/internal/runner/adapter_cli.go](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/internal/runner/adapter_cli.go)
- [go-ut-bench/internal/agentconfig/config.go](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/internal/agentconfig/config.go)
- [go-ut-bench/internal/runner/subjects.go](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/internal/runner/subjects.go)

### 2.2 当前直接缺的不是框架能力，而是 Claude Code 专属适配

还没补的主要是：

1. `agents.yaml / agents.example.yaml` 里新增 `claudecode` framework
2. Agent 沙箱镜像里安装 `claude` CLI
3. `agent_native` skill 注入支持 Claude Code 原生目录
4. Claude Code 非交互模式命令模板
5. Claude Code 的 usage / session / JSON 输出解析
6. 认证与 provider 选择策略

---

## 3. 官方文档里对我们最重要的信息

下面只保留对 UT-Bench 接入有直接价值的部分。

### 3.1 安装与运行入口

官方快速开始文档说明了 Claude Code 的安装、登录和基本 CLI 用法。  
入口文档：

- [Claude Code overview](https://docs.anthropic.com/en/docs/claude-code/overview)
- [Quickstart](https://docs.anthropic.com/en/docs/claude-code/quickstart)
- [CLI reference](https://docs.anthropic.com/en/docs/claude-code/cli-reference)

对 UT-Bench 的直接意义：

- 我们要的不是交互式 TUI；
- 我们要的是 **非交互 / 可脚本化 / 有结构化输出** 的命令调用方式。

### 3.2 非交互模式是官方支持的

官方 CLI 文档明确提供了适合自动化的参数：

- `-p` / `--print`
- `--output-format text|json|stream-json`
- `--input-format stream-json`
- `--resume`
- `--continue`
- `--max-turns`
- `--allowedTools`
- `--disallowedTools`
- `--permission-mode`
- `--json-schema`
- `--mcp-config`
- `--model`

这对 UT-Bench 很关键，因为我们现在的 `opencode` / `codebuddy` 就是走 CLI agent 的非交互模式。

**推荐方向：**

- Claude Code 接入时优先用 `-p`
- 输出优先尝试 `--output-format stream-json`
- 如 stream-json 难以稳定解析，再退回 `json`

原因：

- `stream-json` 更适合后续做工具调用、轮次、usage、事件流采集；
- 但真正的数据结构，仍然需要你跑一个最小样例做实测，不要只靠文档想象。

### 3.3 权限模式是必须考虑的

Claude Code 官方有明确的权限 / 安全模式文档：

- [Permission modes](https://docs.anthropic.com/en/docs/claude-code/permission-modes)
- [Settings](https://docs.anthropic.com/en/docs/claude-code/settings)

对 UT-Bench 最关键的是：

- `bypassPermissions`
- `dontAsk`
- 默认交互权限模式

官方对 `bypassPermissions` 的态度很明确：  
**只能在受信任、隔离的环境中使用。**

这和你现在的 UT-Bench 设计是匹配的：

- Agent 在独立 sandbox 容器内运行；
- 对源码有只读保护；
- 未来还在做更强的 runtime / sandbox 重构。

所以，对 UT-Bench 来说，Claude Code 进入沙箱后使用：

- `--permission-mode bypassPermissions`

是合理候选。

但有两个前提：

1. 沙箱镜像里不要注入无关宿主机秘密；
2. 网络能力、挂载范围、可执行命令需要继续收紧。

### 3.4 配置文件、CLAUDE.md、skills / commands 体系

官方相关文档：

- [Memory and CLAUDE.md files](https://docs.anthropic.com/en/docs/claude-code/memory)
- [Settings](https://docs.anthropic.com/en/docs/claude-code/settings)
- [Commands](https://docs.anthropic.com/en/docs/claude-code/commands)
- [Skills](https://docs.anthropic.com/en/docs/claude-code/skills)

这部分对你非常重要，因为你现在已经有 `agent_native` skill 注入机制。

对 UT-Bench 的直接映射建议：

#### 1. 项目级说明文件

Claude Code 认 `CLAUDE.md`。  
这和你现有的 “prompt + skill + workspace 注入” 机制可以共存。

但不要一开始就把 UT-Bench 的主 prompt 全塞到 `CLAUDE.md`。

建议：

- 主任务说明仍然通过 `{{.ContainerPrompt}}` 注入
- `CLAUDE.md` 作为长期规则 / 团队规范载体

#### 2. 原生 skill 注入

Claude Code 官方支持 skills / commands 目录体系。  
对于 UT-Bench 的 `agent_native`，更合适的目标是：

- `.claude/skills/<skill_name>/SKILL.md`

如果 skill 附带 references / scripts，也按官方技能目录方式落进去。

这意味着你后续需要在：

- [go-ut-bench/internal/runner/subjects.go](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/internal/runner/subjects.go)

里给 `claudecode` 增加类似：

- `copySkillToNativeDir(workRoot, skill, ".claude", "skills")`

### 3.5 认证与 provider 选择

官方认证 / provider 相关文档：

- [IAM / authentication](https://docs.anthropic.com/en/docs/claude-code/iam)
- [Amazon Bedrock](https://code.claude.com/docs/en/amazon-bedrock)
- [Google Vertex AI](https://code.claude.com/docs/en/google-vertex-ai)
- [Microsoft Azure / Foundry](https://code.claude.com/docs/en/microsoft-azure)
- [Third-party integrations](https://code.claude.com/docs/en/third-party-integrations)
- [LLM gateway](https://code.claude.com/docs/en/llm-gateway)

对 UT-Bench 最值得记录的结论：

#### 1. 官方明确支持多种后端，但不是无约束的“各家模型”

Claude Code 可以接：

- Anthropic 直连
- Bedrock
- Vertex
- Foundry
- Gateway / Proxy

但文档里的 Gateway 接法，仍然围绕：

- Anthropic 风格 API
- Claude 模型发现
- provider passthrough

#### 2. 网关模式的关键变量

官方文档里和我们最相关的是：

- `ANTHROPIC_BASE_URL`
- `ANTHROPIC_AUTH_TOKEN`
- `ANTHROPIC_API_KEY`
- `apiKeyHelper`

这意味着，如果你要让 UT-Bench 的 `claudecode` framework 动态切换后端：

- 最现实的方式不是给 Claude Code 写 OpenAI 风格 `models.json`
- 而是通过环境变量给它注入：
  - 网关地址
  - 鉴权 token / API key
  - 必要时的取 key helper

#### 3. 模型发现机制有约束

官方对第三方集成写得很明确：

- 如果要让 Claude Code 自动从网关发现模型，通常依赖 `/v1/models`
- 返回的模型 ID 需要满足它的识别规则
- 文档明确提到只会自动加入 **以 `claude` 或 `anthropic` 开头** 的模型

这对你“接各家模型”的计划有一个非常重要的约束：

> 如果你打算通过 Claude Code 跑 DeepSeek / Qwen / Minimax / 豆包，不要默认认为直接把这些模型名暴露给 Claude Code 就能正常工作。

更稳的方式可能是：

- 在网关层做 **Claude 风格模型别名**
- 再把真实下游模型映射到这些别名

但这已经不是“Claude Code 原生支持各家模型”，而是 **你自己做模型映射层**。

#### 4. 自动登录 / CI token 也要考虑

Claude Code 官方支持非交互认证方式。  
对 UT-Bench 的 sandbox 容器来说，优先级建议是：

1. **环境变量型认证**
2. **网关 token**
3. **apiKeyHelper**
4. 避免依赖浏览器 OAuth 弹登录

因为 benchmark 容器必须是可重放、可脚本化的。

### 3.6 Dev container / sandbox 建议

官方 dev container 文档：

- [Dev containers](https://code.claude.com/docs/en/devcontainer)

虽然你不是直接做 VS Code devcontainer 集成，但这个文档对沙箱设计有现实意义：

- Claude Code 官方承认容器化开发环境是常见用法；
- 但也明确提醒：容器并不自动等于完全安全；
- 如果容器内有网络、有 secrets、有过宽挂载，仍然可能被恶意仓库利用。

这和你现在 UT-Bench 的 runtime 重构方向是一致的：

- control plane / eval plane / agent sandbox plane 分层；
- sandbox provider、sandbox image、只读挂载、env 注入需要继续收紧；
- 不要因为“已经在 Docker 里”就放弃权限控制。

---

## 4. “Claude Code 接各家模型”这件事，应该怎么理解

这是这次文档里最需要讲清楚的部分。

### 4.1 可以做，但要分三种等级

#### 等级 1：官方原生支持

- Anthropic API
- Bedrock
- Vertex
- Foundry

这是最稳的。

#### 等级 2：官方承认的 Gateway / Proxy

通过：

- `ANTHROPIC_BASE_URL`
- `ANTHROPIC_AUTH_TOKEN`
- `ANTHROPIC_API_KEY`
- `apiKeyHelper`

接到你自己的代理层。

这也合理，但你要承担模型映射和兼容性解释。

#### 等级 3：用 Gateway 把任意第三方模型“伪装成 Claude”

理论上有机会，但风险最大：

- 模型名发现未必通过
- 能力假设可能不一致
- 工具调用、JSON 输出、长上下文、会话恢复可能和官方 Claude 行为不同
- 评测出来的结果到底在衡量“Claude Code”还是“某个第三方模型 + Claude Code 外壳”，解释会变复杂

### 4.2 对 benchmark 报告的要求

如果你真的要让 Claude Code 跑第三方模型，报告层一定要区分：

- `framework = claudecode`
- `provider_mode = anthropic|bedrock|vertex|foundry|gateway`
- `gateway_vendor = litellm|自建代理|其他`
- `underlying_model = deepseek-v4-flash / qwen / ...`
- `claude_visible_model_id = 传给 Claude Code 的模型名`

否则最终结果会混淆：

- 是 Claude Code 框架能力？
- 还是第三方模型能力？
- 还是网关兼容层的副作用？

---

## 5. 对 UT-Bench 的接入建议

## 5.1 framework 命名

建议统一使用：

- `claudecode`

不要混：

- `claude_code`
- `claude-code`
- `claude`

原因：

- 你现在 subject id 已经是 `opencode__...` / `codebuddy__...`
- `claudecode__...` 最一致

### 5.2 第一版只做官方直连模式

建议第一版先做：

- Anthropic API

可选第二步：

- Bedrock
- Vertex

先不要第一版就做“Claude Code + 任意第三方模型”。

### 5.3 第一版命令模板建议

可先按下面方向设计：

```bash
PROMPT="$(cat {{.ContainerPrompt}})" &&
claude -p \
  --output-format stream-json \
  --permission-mode bypassPermissions \
  --max-turns 50 \
  --model "{{.ModelID}}" \
  "$PROMPT"
```

说明：

- `-p`：非交互模式
- `--output-format stream-json`：便于后续采 event / usage / tool calls
- `--permission-mode bypassPermissions`：适合隔离容器内 benchmark
- `--max-turns 50`：与现有 CodeBuddy 设计保持量级一致
- `--model`：显式固定模型

但这里先不要写死进仓库，接入前应先做一次最小实测，验证：

1. stream-json 实际 schema
2. session / usage 字段是否稳定存在
3. 工具调用事件如何表达
4. stderr / stdout 分流行为

### 5.4 第一版环境变量建议

如果走 Anthropic 原生：

- `ANTHROPIC_API_KEY`

如果走网关：

- `ANTHROPIC_BASE_URL`
- `ANTHROPIC_AUTH_TOKEN` 或 `ANTHROPIC_API_KEY`

如果后续做更复杂动态切换，再考虑：

- `apiKeyHelper`

不建议第一版就把所有 provider 逻辑揉成一个大命令模板。

### 5.5 skill 注入建议

对于 `agent_native`，建议优先用：

- `.claude/skills/<skill_name>/SKILL.md`

不要先走 `.claude/commands`，因为你这里的 skill 更像长期规则，不是 slash command。

### 5.6 usage 采集策略建议

Claude Code 的使用量采集不要照搬 OpenCode / CodeBuddy。

推荐顺序：

1. **优先解析官方 JSON / stream-json 输出里的 usage**
2. 如果存在 session export / transcript export，再补采
3. 解析不到时，明确标记：
   - `token_source = missing`
   - 不要估一个看起来好看的低值

这点你现在在 [go-ut-bench/internal/runner/adapter_cli.go](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/internal/runner/adapter_cli.go) 里对 CLI agent 已经有正确方向了，Claude Code 应该保持同样标准。

---

## 6. 接入 Claude Code 时具体要改哪些地方

下面按代码改动面列清楚。

### 6.1 配置层

文件：

- [go-ut-bench/configs/agents.example.yaml](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/configs/agents.example.yaml)
- [go-ut-bench/configs/agents.yaml](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/configs/agents.yaml)

新增 `frameworks.claudecode`：

- `enabled`
- `kind: cli_agent`
- `sandbox.*`
- `preflight`
- `forbidden_command_patterns`
- `env`
- `env_from_host`
- `command`
- `compatible_languages`
- `output_globs`

### 6.2 skill 原生注入

文件：

- [go-ut-bench/internal/runner/subjects.go](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/internal/runner/subjects.go)

当前已经支持：

- CodeBuddy → `.codebuddy/skills/...`
- OpenCode → `.opencode/skills/...`

需要补：

- Claude Code → `.claude/skills/...`

### 6.3 CLI 命令执行与解析

文件：

- [go-ut-bench/internal/runner/adapter_cli.go](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/internal/runner/adapter_cli.go)

需要补三块：

1. Claude Code 的命令模板渲染
2. Claude Code 的 JSON / stream-json usage 解析
3. Claude Code 的工具调用 / 轮次 / session 解析

注意：

- 代码注释里已经写了“支持 OpenCode、Claude Code、CodeBuddy 等 CLI Agent”
- 但目前真正有针对性解析的是 OpenCode / CodeBuddy
- Claude Code 现在还没有实质适配

### 6.4 Agent 沙箱镜像

你现在已有统一沙箱镜像方向。  
需要决定：

#### 方案 A：继续统一镜像

在统一 agent sandbox image 里安装：

- `claude`
- `opencode`
- `codebuddy`

优点：

- 管理简单

缺点：

- 镜像更胖
- 版本冲突和升级节奏更复杂

#### 方案 B：按 framework 拆 agent image

例如：

- `utbench-agent-opencode:latest`
- `utbench-agent-codebuddy:latest`
- `utbench-agent-claudecode:latest`

从长期架构看，这更合理。

你现在在 runtime / sandbox 重构，建议认真考虑走 **方案 B**。

### 6.5 Web / report / registry

后续接入后要同步更新：

- framework 注册表
- subject 注册表
- runtime topology 展示
- report 中的 framework/provider/model visibility

尤其如果后面做 “Claude Code + Gateway + 第三方模型”，必须把：

- Claude 可见模型名
- 底层真实模型名
- provider 模式

分开显示。

---

## 7. 推荐接入顺序

### 第一步：官方直连最小闭环

目标：

- `claudecode + Anthropic API + no_skill`

只验证：

1. sandbox 内 `claude` 能启动
2. 非交互 `-p` 能执行
3. 能生成测试文件
4. 能拿到最小 usage / stdout / stderr / session 数据

### 第二步：加 skill

目标：

- `claudecode + Anthropic API + unit_test_skill`

验证：

1. `.claude/skills/...` 是否被 Claude Code 自动识别
2. 是否需要额外 project-level `CLAUDE.md`

### 第三步：加 provider 变体

按顺序：

1. Bedrock
2. Vertex
3. Foundry
4. Gateway

### 第四步：最后才尝试“各家模型”

如果要做：

- 不要直接让模型名暴露成 `deepseek-v4-flash`
- 应先设计网关映射和报告字段

---

## 8. 明确的风险与不要踩的坑

### 8.1 不要把 Claude Code 当成 OpenCode 那样的自由 provider 容器

OpenCode 当前更像“Agent 外壳 + 自定义 provider config”。  
Claude Code 不是这个产品心智。

如果你照 OpenCode 的思路去强行拼：

- 非 Claude 风格模型名
- 任意 OpenAI 兼容地址
- 任意响应 schema

最后很容易得到一个“能跑但不稳定、能评但不可解释”的框架。

### 8.2 不要第一版就混第三方模型

先把：

- `claudecode + 官方 Claude`

跑通，才能知道后面的问题到底来自：

- Claude Code 本身
- 你的 sandbox
- 还是第三方 provider 网关

### 8.3 不要先假设 usage schema

Claude Code 文档告诉你有 JSON / stream-json 输出，但不会替你保证你关心的每个字段都稳定。  
接入时必须做最小实测，把真实输出样本存下来。

### 8.4 不要把 OAuth 登录当 benchmark 主路径

benchmark 容器最需要的是：

- 可脚本化
- 可重放
- 无交互

所以主路径应该是：

- API key
- gateway token
- apiKeyHelper

而不是浏览器登录。

---

## 9. 官方文档索引

下面是这次整理过程中最关键的官方页面，后续真正接入时优先看这些。

### 总览 / 快速开始

- [Claude Code overview](https://docs.anthropic.com/en/docs/claude-code/overview)
- [Quickstart](https://docs.anthropic.com/en/docs/claude-code/quickstart)
- [CLI reference](https://docs.anthropic.com/en/docs/claude-code/cli-reference)

### 权限 / 配置 / 记忆

- [Permission modes](https://docs.anthropic.com/en/docs/claude-code/permission-modes)
- [Settings](https://docs.anthropic.com/en/docs/claude-code/settings)
- [Memory and CLAUDE.md files](https://docs.anthropic.com/en/docs/claude-code/memory)
- [Commands](https://docs.anthropic.com/en/docs/claude-code/commands)
- [Skills](https://docs.anthropic.com/en/docs/claude-code/skills)

### 认证 / provider / 网关

- [IAM / authentication](https://docs.anthropic.com/en/docs/claude-code/iam)
- [Amazon Bedrock](https://code.claude.com/docs/en/amazon-bedrock)
- [Google Vertex AI](https://code.claude.com/docs/en/google-vertex-ai)
- [Microsoft Azure / Foundry](https://code.claude.com/docs/en/microsoft-azure)
- [Third-party integrations](https://code.claude.com/docs/en/third-party-integrations)
- [LLM gateway](https://code.claude.com/docs/en/llm-gateway)

### 容器 / 开发环境

- [Dev containers](https://code.claude.com/docs/en/devcontainer)

---

## 10. 对下一步实施的建议

如果下一步开始真正接代码，建议按下面顺序推进：

1. 新增 `claudecode` framework 配置样例
2. 先做 Anthropic 官方 API 直连
3. 在统一 sandbox image 或独立 agent image 中安装 `claude`
4. 增加 `.claude/skills/...` 原生 skill 注入
5. 用最小样例采集真实 `stream-json` 输出，补 usage / tool / session parser
6. 跑 smoke benchmark
7. 再评估 Bedrock / Vertex / Gateway
8. 最后再讨论“Claude Code 跑各家模型”

如果只用一句话概括：

> Claude Code 可以接进 UT-Bench，但应先按“官方 Claude 工作流基线”接入，再把 Gateway / 第三方模型当成二阶段问题处理；不要一开始把它做成“任意模型通用前端”。
