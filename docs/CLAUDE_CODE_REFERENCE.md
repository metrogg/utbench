# Claude Code 官方文档整理 — UT-Bench 适配参考

> 整理时间: 2026-05-05 | 来源: https://code.claude.com/docs/
> 本文档聚焦 UT-Bench 沙箱评测场景所需的全部适配信息

---

## 目录

1. [Headless 模式 (Print Mode)](#1-headless-模式-print-mode)
2. [CLI 标志完整参考](#2-cli-标志完整参考)
3. [输出格式 (text / json / stream-json)](#3-输出格式)
4. [环境变量完整列表](#4-环境变量完整列表)
5. [权限模式](#5-权限模式)
6. [LLM Gateway 与第三方模型接入](#6-llm-gateway-与第三方模型接入)
7. [模型配置](#7-模型配置)
8. [各厂商 Anthropic 兼容接口汇总](#8-各厂商-anthropic-兼容接口汇总)
9. [Skills 系统](#9-skills-系统)
10. [Hooks 系统](#10-hooks-系统)
11. [CLAUDE.md 项目指令](#11-claudemd-项目指令)
12. [UT-Bench 适配要点](#12-ut-bench-适配要点)

---

## 1. Headless 模式 (Print Mode)

Claude Code 的非交互模式通过 `-p` / `--print` 标志启用，适用于 CI/CD、脚本和自动化评测。

### 基本用法

```bash
# 单次执行
claude -p "你的prompt"

# 流式JSON输出 + 最大轮次限制
claude -p "你的prompt" --output-format stream-json --max-turns 50 --verbose

# 最简模式(跳过hooks/skills/plugins/MCP发现，启动更快)
claude --bare -p "你的prompt" --output-format stream-json --verbose
```

### 与 UT-Bench 相关的关键标志组合

```bash
# 推荐的评测命令模板
claude -p "{{.ContainerPrompt}}" \
  --bare \
  --output-format stream-json \
  --verbose \
  --max-turns 50 \
  --permission-mode bypassPermissions \
  --system-prompt-file /path/to/system-prompt.txt \
  --append-system-prompt "You are running in an isolated Docker sandbox..."
```

| 标志 | 作用 | 评测必要性 |
|------|------|-----------|
| `-p` | 非交互模式 | **必需** |
| `--bare` | 跳过自动发现(hooks/skills/plugins/MCP)，加速启动 | **推荐** |
| `--output-format stream-json` | 流式JSON输出，便于程序解析 | **必需** |
| `--verbose` | 输出thinking/tool_use/tool_result等详细过程 | **推荐** |
| `--max-turns N` | 限制agentic轮次 | **必需** |
| `--permission-mode bypassPermissions` | 跳过所有权限提示 | **必需**(沙箱环境) |
| `--model <model>` | 指定模型(可用别名或完整名) | 可选(已有env) |
| `--system-prompt-file` | 替换系统提示 | 可选 |
| `--append-system-prompt` | 追加系统提示 | 可选 |
| `--max-budget-usd N` | 限制花费(美元) | 可选(防止失控) |

---

## 2. CLI 标志完整参考

### 2.1 核心命令

| 命令 | 描述 |
|------|------|
| `claude` | 启动交互式会话 |
| `claude "query"` | 带初始提示启动交互式会话 |
| `claude -p "query"` | 非交互模式，执行后退出 |
| `cat file \| claude -p "query"` | 处理管道输入 |
| `claude -c` | 继续最近对话 |
| `claude -r "<session>" "query"` | 按ID或名称恢复会话 |
| `claude update` | 更新到最新版本 |
| `claude install [version]` | 安装/重装原生二进制 |
| `claude auth status --text` | 查看认证状态 |

### 2.2 全部标志

| 标志 | 描述 |
|------|------|
| `--add-dir` | 添加额外工作目录供读写 |
| `--agent` | 指定当前会话的子代理 |
| `--allowedTools` | 无需权限即可使用的工具列表 |
| `--append-system-prompt` | 追加文本到默认系统提示末尾 |
| `--append-system-prompt-file` | 从文件加载并追加系统提示 |
| **`--bare`** | **最小模式: 跳过hooks/skills/plugins/MCP/自动内存/CLAUDE.md发现** |
| `--betas` | 包含额外anthropic-beta头 |
| `--continue`, `-c` | 加载当前目录最近的对话 |
| **`--dangerously-skip-permissions`** | **跳过权限提示(等同`--permission-mode bypassPermissions`)** |
| `--debug` | 启用调试模式 |
| `--debug-file <path>` | 调试日志写入指定文件 |
| `--disallowedTools` | 移除且不可用的工具 |
| `--effort` | 努力级别: `low`/`medium`/`high`/`xhigh`/`max` |
| `--exclude-dynamic-system-prompt-sections` | 从系统提示移除每机器部分，改善缓存复用 |
| `--fallback-model` | 默认模型过载时回退模型(仅print模式) |
| `--fork-session` | 恢复时创建新会话ID |
| `--from-pr` | 恢复链接到特定PR的会话 |
| `--include-hook-events` | 在输出中包含hook生命周期事件(需stream-json) |
| **`--include-partial-messages`** | **包含部分流事件(逐token)，需配合-p和stream-json** |
| `--init` | 会话前运行init hooks |
| `--init-only` | 仅运行Setup和SessionStart hooks后退出 |
| `--input-format` | print模式输入格式: `text`/`stream-json` |
| `--json-schema` | 获取匹配JSON Schema的输出(仅print模式) |
| `--maintenance` | 会话前运行maintenance hooks |
| `--max-budget-usd` | API花费上限(美元，仅print模式) |
| **`--max-turns`** | **限制代理轮次(仅print模式)，到达限制时以错误退出** |
| `--mcp-config` | 加载MCP服务器配置 |
| `--model` | 设置模型(别名或完整名) |
| `--name`, `-n` | 设置会话显示名称 |
| `--no-session-persistence` | 禁用会话持久化(仅print模式) |
| **`--output-format`** | **输出格式: `text`/`json`/`stream-json`** |
| **`--permission-mode`** | **权限模式: `default`/`acceptEdits`/`plan`/`auto`/`dontAsk`/`bypassPermissions`** |
| `--permission-prompt-tool` | MCP工具处理权限提示(非交互模式) |
| **`--print`, `-p`** | **非交互模式** |
| `--resume`, `-r` | 按ID或名称恢复会话 |
| `--session-id` | 使用特定会话ID |
| `--settings` | 加载额外设置(JSON文件或字符串) |
| `--system-prompt` | 替换整个系统提示 |
| `--system-prompt-file` | 从文件替换系统提示 |
| `--tools` | 限制可用的内置工具(`""`禁用全部,`"default"`全部,`"Bash,Edit,Read"`) |
| `--verbose` | 启用详细日志，显示完整轮次输出 |
| `--version`, `-v` | 输出版本号 |
| `--worktree`, `-w` | 在隔离git worktree中启动 |

### 2.3 系统提示标志

| 标志 | 行为 | 互斥 |
|------|------|------|
| `--system-prompt` | 替换整个默认提示 | 与`--system-prompt-file`互斥 |
| `--system-prompt-file` | 用文件内容替换默认提示 | 与`--system-prompt`互斥 |
| `--append-system-prompt` | 追加到默认提示末尾 | — |
| `--append-system-prompt-file` | 追加文件内容到默认提示 | — |

---

## 3. 输出格式

`--output-format` **只在 `-p` (print) 模式下生效**。

### 3.1 `text` (默认)

纯文本输出，无元数据。适用于终端直接查看。

### 3.2 `json` (一次性完整结果)

等 Claude 执行完毕后，一次性返回一个 JSON 对象。

```json
{
  "type": "result",
  "subtype": "success",
  "session_id": "550e8400-e29b-41d4-a716-446655440000",
  "result": "递归就是函数调用自身来解决问题...",
  "cost_usd": 0.003,
  "duration_ms": 1200,
  "duration_api_ms": 950,
  "is_error": false,
  "num_turns": 1,
  "usage": {
    "input_tokens": 42,
    "output_tokens": 38
  }
}
```

**字段说明:**

| 字段 | 类型 | 说明 |
|------|------|------|
| `type` | string | 固定 `"result"` |
| `subtype` | string | `"success"` / `"error"` 等 |
| `session_id` | string | 会话ID，可用于`--resume` |
| `result` | string | Claude的回答内容(自由文本) |
| `cost_usd` | number | 本次调用花费(美元) |
| `duration_ms` | number | 总执行时长(毫秒) |
| `duration_api_ms` | number | API调用时长(毫秒) |
| `is_error` | boolean | 是否发生错误 |
| `num_turns` | integer | agentic循环轮次 |
| `usage.input_tokens` | integer | 输入token数 |
| `usage.output_tokens` | integer | 输出token数 |

搭配 `--json-schema` 时额外返回 `structured_output` 字段。

### 3.3 `stream-json` (NDJSON 实时流)

Newline-delimited JSON (NDJSON)，每行一个独立JSON对象。

#### 事件类型完整列表

| type | subtype / content[].type | 说明 | 需要参数 |
|------|--------------------------|------|---------|
| `system` | `init` | 会话开始，包含可用工具列表、模型、权限模式等 | — |
| `system` | `start` | 请求处理开始 | — |
| `system` | `end` | 请求处理结束(在result之前) | — |
| `system` | `api_retry` | API重试通知 | — |
| `assistant` | `text` | Claude文本输出(同一message.id多次发射) | — |
| `assistant` | `tool_use` | 工具调用请求(同一message.id多次发射) | `--verbose` |
| `assistant` | `thinking` | Claude内部思考过程 | `--verbose` |
| `user` | `text` | 用户输入(程序化使用时通过stdin) | — |
| `user` | `tool_result` | **工具执行结果(在user事件中，不是assistant!)** | `--verbose` |
| `stream_event` | `text_delta` | 逐token文本流 | `--include-partial-messages` |
| `result` | `success` | 最终结果(与json格式结构一致) | — |
| `control_request` | `interrupt` | SDK→Claude方向的中断请求(通过stdin) | — |
| `control_response` | `success`/`error` | Claude→SDK方向的中断响应 | — |

#### 事件结构详解

**system/init:**
```json
{"type":"system","subtype":"init","session_id":"abc-123","tools":["Read","Edit","Bash"]}
```

**system/api_retry:**
```json
{
  "type": "system",
  "subtype": "api_retry",
  "attempt": 1,
  "max_retries": 5,
  "retry_delay_ms": 2000,
  "error_status": 429,
  "error": "rate_limit",
  "uuid": "...",
  "session_id": "..."
}
```

错误类别: `rate_limit` / `server_error` / `authentication_failed` / `billing_error` / `invalid_request` / `max_output_tokens` / `unknown`

**assistant (thinking):**
```json
{
  "type": "assistant",
  "message": {
    "role": "assistant",
    "content": [{"type": "thinking", "thinking": "让我先看看项目结构..."}]
  }
}
```

**assistant (text):**
```json
{
  "type": "assistant",
  "message": {
    "role": "assistant",
    "content": [{"type": "text", "text": "发现 auth.py 第42行..."}]
  }
}
```

**assistant (tool_use):**
```json
{
  "type": "assistant",
  "message": {
    "role": "assistant",
    "content": [{"type": "tool_use", "name": "Bash", "input": {"command": "pytest test_auth.py"}}]
  }
}
```

**user (tool_result) — 工具执行结果在user事件中，不在assistant中:**
```json
{
  "type": "user",
  "message": {
    "role": "user",
    "content": [
      {"type": "tool_result", "tool_use_id": "toolu_01Kp...", "content": "PASSED 5 tests", "is_error": false}
    ]
  },
  "parent_tool_use_id": null,
  "session_id": "..."
}
```

**stream_event (text_delta):**
```json
{
  "type": "stream_event",
  "event": {
    "delta": {"type": "text_delta", "text": "这是增量文本..."}
  }
}
```

**result (最终):**
```json
{
  "type": "result",
  "subtype": "success",
  "session_id": "abc-123",
  "result": "已修复3个bug...",
  "total_cost_usd": 0.05,
  "duration_ms": 15000,
  "num_turns": 8,
  "is_error": false,
  "usage": {"input_tokens": 5000, "output_tokens": 2000, "cache_creation_input_tokens": 0, "cache_read_input_tokens": 0}
}
```

> **注意**: stream-json的result事件中成本字段为 `total_cost_usd`，而json格式的同一字段为 `cost_usd`。usage对象可能包含缓存相关字段(`cache_creation_input_tokens`, `cache_read_input_tokens`)。

#### stream-json 消费注意事项

1. **行缓冲是最大陷阱**: 事件可能跨TCP chunk，必须用变量保留未完成行
2. **空行跳过**: `if (!line.trim()) continue;`
3. **进程退出时处理残留**: `buf`中可能还有最后一个未换行事件
4. **NDJSON不是标准JSON数组**: 逐行解析，不要整体parse
5. **⚠️ assistant事件同一message.id多次发射**: 同一条assistant消息的text和tool_use会分多次发射，每次追加新content块。**必须按message.id去重，保留最后一次的完整content**
6. **⚠️ tool_result在user事件中**: 不要在assistant事件里找tool_result，它在`type:"user"`事件中，字段为`content[].type:"tool_result"`，通过`tool_use_id`关联对应的tool_use

---

## 4. 环境变量完整列表

### 4.1 API 认证与密钥

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `ANTHROPIC_API_KEY` | 作为 `X-Api-Key` 头发送 | 无 |
| `ANTHROPIC_AUTH_TOKEN` | 自定义 `Authorization` 头(自动加`Bearer`前缀) | 无 |
| **`ANTHROPIC_BASE_URL`** | **覆盖API端点(核心!用于第三方模型接入)** | 无 |
| `ANTHROPIC_BEDROCK_BASE_URL` | Bedrock端点URL | 无 |
| `ANTHROPIC_BEDROCK_SERVICE_TIER` | Bedrock服务层级(`default`/`flex`/`priority`) | 无 |
| `ANTHROPIC_BETAS` | 额外`anthropic-beta`头值(逗号分隔) | 无 |
| `ANTHROPIC_CUSTOM_HEADERS` | 自定义请求头(`Name: Value`格式，换行分隔) | 无 |
| `ANTHROPIC_VERTEX_PROJECT_ID` | Vertex AI GCP项目ID | 无 |
| `ANTHROPIC_VERTEX_BASE_URL` | Vertex AI端点URL | 无 |
| `AWS_BEARER_TOKEN_BEDROCK` | Bedrock API密钥认证 | 无 |

### 4.2 模型配置 (核心!)

| 变量 | 说明 | 默认值 |
|------|------|--------|
| **`ANTHROPIC_MODEL`** | **默认使用的模型** | 无 |
| **`ANTHROPIC_DEFAULT_OPUS_MODEL`** | **映射Opus层级(复杂推理)** | 无 |
| **`ANTHROPIC_DEFAULT_SONNET_MODEL`** | **映射Sonnet层级(日常编程)** | 无 |
| **`ANTHROPIC_DEFAULT_HAIKU_MODEL`** | **映射Haiku层级(轻量任务)** | 无 |
| `ANTHROPIC_SMALL_FAST_MODEL` | ~~后台任务模型~~ [已弃用] | 无 |
| `CLAUDE_CODE_SUBAGENT_MODEL` | 子代理使用的模型 | 无 |
| `ANTHROPIC_CUSTOM_MODEL_OPTION` | 添加自定义模型ID到选择器 | 无 |
| `ANTHROPIC_CUSTOM_MODEL_OPTION_NAME` | 自定义模型显示名称 | 默认模型ID |

### 4.3 请求与超时

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `API_TIMEOUT_MS` | API请求超时(毫秒)，最大2147483647 | `600000`(10分钟) |
| `BASH_DEFAULT_TIMEOUT_MS` | bash命令默认超时 | `120000`(2分钟) |
| `BASH_MAX_OUTPUT_LENGTH` | bash输出最大字符数 | 无 |
| `BASH_MAX_TIMEOUT_MS` | 模型可设的最大bash超时 | `600000`(10分钟) |
| `CLAUDE_CODE_MAX_RETRIES` | API请求失败重试次数 | `10` |
| `CLAUDE_CODE_MAX_OUTPUT_TOKENS` | 大多数请求最大输出token | 因模型而异 |
| `CLAUDE_CODE_MAX_CONTEXT_TOKENS` | 覆盖上下文窗口大小(需设DISABLE_COMPACT) | 无 |

### 4.4 上下文与压缩

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `CLAUDE_AUTOCOMPACT_PCT_OVERRIDE` | 自动压缩触发的上下文容量百分比(1-100) | ~`95%` |
| `CLAUDE_CODE_AUTO_COMPACT_WINDOW` | 自动压缩计算的上下文容量(token) | 标准模型200K, 扩展1M |

### 4.5 思考与推理

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `CLAUDE_CODE_DISABLE_THINKING` | 设`1`强制禁用扩展思考 | 无 |
| `CLAUDE_CODE_DISABLE_ADAPTIVE_THINKING` | 设`1`禁用自适应推理 | 无 |
| `CLAUDE_CODE_EFFORT_LEVEL` | 努力级别: `low`/`medium`/`high`/`xhigh`/`max`/`auto` | 无 |

### 4.6 后台任务与并发

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `CLAUDE_AUTO_BACKGROUND_TASKS` | 设`1`强制启用长时间任务自动后台化 | 无 |
| `CLAUDE_CODE_DISABLE_BACKGROUND_TASKS` | 设`1`禁用所有后台任务 | 无 |
| `CLAUDE_CODE_MAX_TOOL_USE_CONCURRENCY` | 只读工具和子代理最大并行数 | `10` |

### 4.7 网络、代理与TLS

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `HTTP_PROXY` / `HTTPS_PROXY` | HTTP(S)代理 | 无 |
| `NO_PROXY` | 不走代理的地址 | 无 |
| `CLAUDE_CODE_CERT_STORE` | TLS CA证书源(逗号分隔: `bundled`,`system`) | `bundled,system` |
| `CLAUDE_CODE_PROXY_RESOLVES_HOSTS` | 设`1`允许代理执行DNS解析 | 无 |

### 4.8 认证跳过 (LLM Gateway用)

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `CLAUDE_CODE_SKIP_BEDROCK_AUTH` | 跳过Bedrock认证 | 无 |
| `CLAUDE_CODE_SKIP_VERTEX_AUTH` | 跳过Vertex认证 | 无 |

### 4.9 遥测与流量控制

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC` | 同时禁用自动更新/反馈/错误报告/遥测 | 无 |
| `CLAUDE_CODE_DISABLE_FEEDBACK_SURVEY` | 禁用会话质量调查 | 无 |

### 4.10 Shell与终端

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `CLAUDE_CODE_SHELL` | 覆盖自动检测的shell | 自动检测 |
| `CLAUDECODE` | Claude Code生成的shell环境中设为`1` | 无 |

### 4.11 简化模式

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `CLAUDE_CODE_SIMPLE` | 设`1`最小化系统提示，仅Bash/读取/编辑工具，禁用自动发现 | 无 |
| `CLAUDE_CODE_DISABLE_EXPERIMENTAL_BETAS` | 设`1`剥离beta头和字段(第三方模型推荐) | 无 |

### 4.12 其他重要

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC` | 禁用非必要网络流量 | 无 |
| `CLAUDE_CODE_SKIP_PROMPT_HISTORY` | 跳过写入会话记录到磁盘 | 无 |
| `CLAUDE_CODE_DISABLE_1M_CONTEXT` | 禁用1M上下文窗口 | 无 |
| `CLAUDE_CODE_EXTRA_BODY` | 合并到每个API请求体的JSON对象 | 无 |
| `CLAUDE_CODE_DISABLE_NONSTREAMING_FALLBACK` | 禁用非流式回退 | 无 |
| `API_KEY_HELPER` | 动态API Key脚本路径(设置文件中) | 无 |
| `CLAUDE_CODE_API_KEY_HELPER_TTL_MS` | 动态Key刷新间隔(毫秒) | 无 |

---

## 5. 权限模式

### 5.1 全部6种模式

| 模式 | 描述 | 适用场景 |
|------|------|---------|
| `default` | 标准行为，首次使用每个工具时提示 | 交互式使用 |
| `acceptEdits` | 自动接受文件编辑和常见文件系统命令 | 轻度自动化 |
| `plan` | Claude可分析但**不能修改文件或执行命令** | 只读审查 |
| `auto` | 自动批准工具调用，后台安全检查 | 研究预览 |
| `dontAsk` | 未预先批准的工具**自动拒绝** | 受限自动化 |
| **`bypassPermissions`** | **跳过所有权限提示** | **沙箱/容器/Docker(评测场景)** |

### 5.2 评测场景推荐

```bash
# UT-Bench 沙箱评测: 使用 bypassPermissions
claude -p "prompt" --permission-mode bypassPermissions
# 或等价的:
claude -p "prompt" --dangerously-skip-permissions
```

> `bypassPermissions` 仍保留 `rm -rf /` 和 `rm -rf ~` 的断路器提示。

### 5.3 权限规则语法

```json
{
  "permissions": {
    "allow": ["Bash(npm run *)", "Bash(git *)"],
    "deny": ["Bash(rm -rf *)"]
  }
}
```

---

## 6. LLM Gateway 与第三方模型接入

### 6.1 核心原理

Claude Code 通过 `ANTHROPIC_BASE_URL` 将请求路由到 LLM Gateway/代理。Gateway 必须实现以下三种 API 格式之一:

| API 格式 | 端点路径 | 说明 |
|---------|---------|------|
| **Anthropic Messages** | `/v1/messages`, `/v1/messages/count_tokens` | **最通用，国产厂商均用此格式** |
| Amazon Bedrock | `/invoke`, `/invoke-with-response-stream` | AWS专属 |
| Google Vertex | `:rawPredict`, `:streamRawPredict` | GCP专属 |

### 6.2 认证方式

| 方式 | 配置 | 说明 |
|------|------|------|
| 静态API Key | `ANTHROPIC_AUTH_TOKEN=sk-xxx` 或 `ANTHROPIC_API_KEY=sk-xxx` | 最简单 |
| 动态API Key | `apiKeyHelper`脚本 + `CLAUDE_CODE_API_KEY_HELPER_TTL_MS` | 支持密钥轮换 |
| 自定义请求头 | `ANTHROPIC_CUSTOM_HEADERS` | 附加HTTP头 |

优先级: `apiKeyHelper` < `ANTHROPIC_AUTH_TOKEN` / `ANTHROPIC_API_KEY`

### 6.3 模型自动发现 (v2.1.126+)

仅当 `ANTHROPIC_BASE_URL` 已设置且**不指向** `api.anthropic.com` 时:
- 自动查询 `/v1/models` 端点
- **仅添加 ID 以 `claude` 或 `anthropic` 开头的模型**
- 缓存至 `~/.claude/cache/gateway-models.json`

### 6.4 设置文件配置方式

`~/.claude/settings.json` 或 `.claude/settings.json`:

```json
{
  "env": {
    "ANTHROPIC_BASE_URL": "https://your-gateway.example.com",
    "ANTHROPIC_AUTH_TOKEN": "sk-xxx",
    "ANTHROPIC_MODEL": "your-model-name",
    "ANTHROPIC_DEFAULT_OPUS_MODEL": "your-model-name",
    "ANTHROPIC_DEFAULT_SONNET_MODEL": "your-model-name",
    "ANTHROPIC_DEFAULT_HAIKU_MODEL": "your-fast-model",
    "CLAUDE_CODE_SUBAGENT_MODEL": "your-fast-model"
  }
}
```

> UT-Bench 不使用 settings.json，直接通过命令行环境变量注入。

---

## 7. 模型配置

### 7.1 模型别名

| 别名 | 解析结果 |
|------|---------|
| `sonnet` | Anthropic API → Sonnet 4.6; Bedrock/Vertex → Sonnet 4.5 |
| `opus` | Anthropic API → Opus 4.7; Bedrock/Vertex → Opus 4.6 |
| `haiku` | 最新Haiku模型 |
| `opusplan` | 计划模式用`opus`，执行模式用`sonnet` |
| `sonnet[1m]` / `opus[1m]` | 带100万token上下文窗口的版本 |

### 7.2 环境变量覆盖别名

```bash
export ANTHROPIC_DEFAULT_OPUS_MODEL='claude-opus-4-7'
export ANTHROPIC_DEFAULT_SONNET_MODEL='claude-sonnet-4-5'
```

支持 `[1m]` 后缀启用扩展上下文。

### 7.3 手动添加自定义模型

```bash
export ANTHROPIC_CUSTOM_MODEL_OPTION="my-model-id"
export ANTHROPIC_CUSTOM_MODEL_OPTION_NAME="My Model Display Name"
export ANTHROPIC_CUSTOM_MODEL_OPTION_DESCRIPTION="Description..."
export ANTHROPIC_CUSTOM_MODEL_OPTION_SUPPORTED_CAPABILITIES="effort,thinking"
```

### 7.4 可声明的能力值

| 能力值 | 启用的功能 |
|--------|-----------|
| `effort` | 效力等级和`/effort`命令 |
| `xhigh_effort` | `xhigh`效力等级 |
| `max_effort` | `max`效力等级 |
| `thinking` | 扩展思考 |
| `adaptive_thinking` | 自适应推理(按任务复杂度动态分配) |
| `interleaved_thinking` | 工具调用间的思考 |

### 7.5 modelOverrides (精细映射)

```json
{
  "modelOverrides": {
    "claude-opus-4-7": "arn:aws:bedrock:us-east-2:123:app-inference-profile/opus",
    "claude-sonnet-4-5": "arn:aws:bedrock:us-east-2:123:app-inference-profile/sonnet"
  }
}
```

### 7.6 优先级

```
--model / ANTHROPIC_MODEL (原样传递)        → 最高
modelOverrides (替换/model选择器)           → 次之
ANTHROPIC_DEFAULT_*_MODEL (覆盖别名)        → 再次
ANTHROPIC_CUSTOM_MODEL_OPTION (添加条目)    → 最低
```

---

## 8. 各厂商 Anthropic 兼容接口汇总

> 所有国产模型厂商均提供了 `/anthropic` 端点，实现 Anthropic Messages API 协议兼容。

### 8.1 DeepSeek

```bash
export ANTHROPIC_BASE_URL=https://api.deepseek.com/anthropic
export ANTHROPIC_AUTH_TOKEN=<your-deepseek-api-key>
export ANTHROPIC_MODEL=deepseek-v4-pro[1m]
export ANTHROPIC_DEFAULT_OPUS_MODEL=deepseek-v4-pro[1m]
export ANTHROPIC_DEFAULT_SONNET_MODEL=deepseek-v4-pro[1m]
export ANTHROPIC_DEFAULT_HAIKU_MODEL=deepseek-v4-flash
export CLAUDE_CODE_SUBAGENT_MODEL=deepseek-v4-flash
```

**可用模型:** `deepseek-v4-pro`, `deepseek-v4-flash`

### 8.2 阿里百炼 (Qwen)

```bash
export ANTHROPIC_BASE_URL=https://dashscope.aliyuncs.com/apps/anthropic
export ANTHROPIC_AUTH_TOKEN=<your-bailian-api-key>
export ANTHROPIC_MODEL=qwen3.6-plus
export ANTHROPIC_DEFAULT_OPUS_MODEL=qwen3.6-plus
export ANTHROPIC_DEFAULT_SONNET_MODEL=qwen3.6-plus
export ANTHROPIC_DEFAULT_HAIKU_MODEL=qwen3.6-flash
export CLAUDE_CODE_SUBAGENT_MODEL=qwen3.6-flash
```

**可用模型:** `qwen3.6-max-preview`, `qwen3.6-plus`, `qwen3.6-flash`, `qwen3-coder-plus`, `qwen3-coder-next`, `qwen3.5-plus`, `qwen3.5-flash`, `qwen-turbo` 等

**注意:** Base URL与API Key必须归属同一地域。新加坡区域: `https://dashscope-intl.aliyuncs.com/apps/anthropic`

### 8.3 智谱 GLM

```bash
export ANTHROPIC_BASE_URL=https://open.bigmodel.cn/api/anthropic
export ANTHROPIC_AUTH_TOKEN=<your-zhipu-api-key>
export ANTHROPIC_MODEL=glm-5.1
export ANTHROPIC_DEFAULT_OPUS_MODEL=glm-5.1
export ANTHROPIC_DEFAULT_SONNET_MODEL=glm-5.1
export ANTHROPIC_DEFAULT_HAIKU_MODEL=glm-4.7
export CLAUDE_CODE_SUBAGENT_MODEL=glm-4.7
```

**可用模型:** `glm-5.1`, `glm-5`, `glm-4.7`, `glm-4.6`

### 8.4 MiniMax

> **注意**: 国内用户用 `api.minimaxi.com`（多一个i），国际用户用 `api.minimax.io`。

```bash
# 国内用户 (推荐)
export ANTHROPIC_BASE_URL=https://api.minimaxi.com/anthropic
# 国际用户
# export ANTHROPIC_BASE_URL=https://api.minimax.io/anthropic
export ANTHROPIC_AUTH_TOKEN=<your-minimax-api-key>
export ANTHROPIC_MODEL=MiniMax-M2.7
export ANTHROPIC_DEFAULT_OPUS_MODEL=MiniMax-M2.7
export ANTHROPIC_DEFAULT_SONNET_MODEL=MiniMax-M2.7
export ANTHROPIC_DEFAULT_HAIKU_MODEL=MiniMax-M2.7
export CLAUDE_CODE_SUBAGENT_MODEL=MiniMax-M2.7
export API_TIMEOUT_MS=3000000
export CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1
```

**可用模型:** `MiniMax-M2.7`, `MiniMax-M2.7-highspeed`, `MiniMax-M2.5`, `MiniMax-M2.5-highspeed`, `MiniMax-M2.1`, `MiniMax-M2.1-highspeed`, `MiniMax-M2`

**⚠️ 重要**: 配置前需先清除已有 Anthropic 环境变量: `unset ANTHROPIC_AUTH_TOKEN && unset ANTHROPIC_BASE_URL`

### 8.5 豆包 (火山引擎)

> **注意**: 端点路径因套餐不同而异。Coding Plan 套餐用 `/api/coding`，普通兼容接口用 `/api/compatible`。请确认你的套餐类型后选择对应端点。

```bash
# Coding Plan 套餐
export ANTHROPIC_BASE_URL=https://ark.cn-beijing.volces.com/api/coding
# 普通 Anthropic 兼容接口
# export ANTHROPIC_BASE_URL=https://ark.cn-beijing.volces.com/api/compatible
export ANTHROPIC_AUTH_TOKEN=<your-volcengine-api-key>
export ANTHROPIC_MODEL=doubao-seed-2.0-code
export ANTHROPIC_DEFAULT_OPUS_MODEL=doubao-seed-2.0-code
export ANTHROPIC_DEFAULT_SONNET_MODEL=doubao-seed-2.0-code
export ANTHROPIC_DEFAULT_HAIKU_MODEL=doubao-seed-2.0-code
export CLAUDE_CODE_SUBAGENT_MODEL=doubao-seed-2.0-code
```

**注意:** Coding Plan 套餐支持通过同一端点调用第三方模型(Kimi-K2.5、GLM系列等)，具体可用模型以火山方舟控制台为准。

### 8.6 一站式多模型平台 (推荐)

| 平台 | Base URL | 首月价格 | 可用模型 |
|------|----------|---------|---------|
| **阿里百炼** | `https://dashscope.aliyuncs.com/apps/anthropic` | 7.5元 | Qwen全系列、DeepSeek、智谱、MiniMax、Kimi、GLM **(北京地域独有第三方模型)** |
| **字节方舟** | `https://ark.cn-beijing.volces.com/api/coding` | 9.9元 | 豆包、Kimi-K2.5、GLM系列等(Coding Plan套餐) |

> **百炼平台优势**: 北京地域下可通过同一端点直接调用DeepSeek/智谱/MiniMax/Kimi/GLM等第三方模型（见§8.2的完整模型列表），免去多厂商分别接入的麻烦。新加坡地域仅支持Qwen系列。

---

## 9. Skills 系统

### 9.1 概念

Skills 扩展了 Claude 能做的事情。通过创建 `SKILL.md` 文件（包含指令），Claude 会将其添加到工具包中。**Claude 在相关时自动使用 skills**，或用户可通过 `/skill-name` 手动调用。

与 CLAUDE.md 的关键区别：Skill 的正文**仅在使用时加载**，长参考资料在需要之前几乎不花费任何 token。

### 9.2 内置捆绑 Skills

| 捆绑 Skill | 说明 |
|:--|:--|
| `/simplify` | 简化代码或内容 |
| `/batch` | 批量处理 |
| `/debug` | 调试 |
| `/loop` | 循环执行 |
| `/claude-api` | Claude API 相关操作 |

内置命令（`/help`、`/compact`等）直接执行固定逻辑；捆绑 skills 是**基于提示的**，为 Claude 提供详细剧本，让它使用工具编排工作。

### 9.3 SKILL.md 格式

```yaml
---
name: my-skill                        # 显示名称和斜杠命令(省略用目录名)
description: What this skill does     # Claude据其判断何时使用(最关键字段!)
when_to_use: Additional trigger hints # 额外触发上下文(description+when_to_use限1536字符)
argument-hint: [filename] [format]    # 自动补全提示
arguments: file format                # 命名位置参数(空格分隔或YAML列表)
disable-model-invocation: true        # true=禁止自动触发,仅手动 /skill-name
user-invocable: true                  # false=从/菜单隐藏,仅作背景知识
allowed-tools: Read, Grep, Glob      # 激活时免权限使用的工具
model: opus                           # 激活时覆盖使用的模型
effort: high                          # 工作量级别: low/medium/high/xhigh/max
context: fork                         # fork=在隔离的subagent中运行
agent: Explore                        # context:fork时的subagent类型
paths: "src/**/*.ts"                  # Glob模式,限制何时自动激活
shell: bash                           # bash(默认)或powershell,用于shell注入
hooks:                                # 限定于skill生命周期的hooks
---

## Skill 内容 (Markdown)

当解释代码时，始终包含:
1. **类比**: 将代码与日常事物比较
2. **图示**: 使用ASCII art展示流程
3. **逐步讲解**: 解释每一步
4. **陷阱**: 常见错误或误解
```

### 9.4 字符串替换

| 变量 | 描述 |
|------|------|
| `$ARGUMENTS` | 传递的所有参数(完整字符串) |
| `$ARGUMENTS[N]` / `$N` | 按索引访问特定参数(0基) |
| `$name` | `arguments` 列表中的命名参数 |
| `${CLAUDE_SESSION_ID}` | 当前会话 ID |
| `${CLAUDE_EFFORT}` | 当前工作量级别 |
| `${CLAUDE_SKILL_DIR}` | skill 目录路径 |

### 9.5 动态上下文注入

`` !`command` `` 语法在发送给 Claude 前执行 shell 命令，输出替换占位符：

```markdown
## Pull request context
- PR diff: !`gh pr diff`
- PR comments: !`gh pr view --comments`
- Changed files: !`gh pr diff --name-only`
```

多行命令用 ` ```! ` 围栏代码块。

### 9.6 目录结构

```
my-skill/
├── SKILL.md           # 核心指令(必需, 建议500行以下)
├── reference.md       # 详细API文档(按需加载)
├── examples/
│   └── sample.md      # 示例输出
└── scripts/
    └── validate.sh    # 辅助脚本(Claude可执行)
```

在 SKILL.md 中引用附属文件：
```markdown
For complete API details, see [reference.md](reference.md).
For usage examples, see [examples/sample.md](examples/sample.md).
```

### 9.7 存储位置

| 级别 | 路径 | 适用范围 |
|:--|:--|:--|
| **企业** | 托管设置 | 组织内所有用户 |
| **个人** | `~/.claude/skills/<skill-name>/SKILL.md` | 用户所有项目 |
| **项目** | `.claude/skills/<skill-name>/SKILL.md` | 仅当前项目 |
| **插件** | `<plugin>/skills/<skill-name>/SKILL.md` | 启用插件的位置 |

**优先级**: 企业 > 个人 > 项目。插件使用 `plugin-name:skill-name` 命名空间。

### 9.8 加载机制

1. **启动时**: 仅加载 frontmatter（名称+描述），每个 Skill 约 **100 token**
2. **对话中**: Claude 读取描述判断哪些 Skills 与当前内容相关
3. **激活时**: Claude 判断需要某个 Skill 时，才加载其完整内容

> 可安装几十个 Skills 而不会撑爆上下文窗口。Claude 只在需要时才加载所需内容。

### 9.9 调用控制矩阵

| Frontmatter | 用户可调用 | Claude可自动调用 | 上下文加载 |
|:--|:--|:--|:--|
| （默认） | ✅ | ✅ | 描述始终在上下文中，调用时加载完整skill |
| `disable-model-invocation: true` | ✅ | ❌ | 描述不在上下文中，用户调用时加载 |
| `user-invocable: false` | ❌ | ✅ | 描述始终在上下文中，调用时加载 |

### 9.10 Skills vs CLAUDE.md vs Hooks

| 系统 | 加载时机 | AI决定? | 保证执行? | 用途 |
|:--|:--|:--|:--|:--|
| **CLAUDE.md** | 每次会话 | 否 | 是 | 项目上下文和规范 |
| **Skills** | 按需/自动检测 | 是 | 否 | 可复用工作流 |
| **Hooks** | 匹配工具事件时 | 否 | 是 | 必须执行的操作 |

### 9.11 Skills vs 旧版 Commands

| 维度 | Commands(旧) | Skills(新) |
|:--|:--|:--|
| **位置** | `.claude/commands/deploy.md` | `.claude/skills/deploy/SKILL.md` |
| **调用方式** | `/deploy` | `/deploy`(相同) |
| **目录支持** | ❌ 单文件 | ✅ 支持目录(模板/示例/脚本) |
| **自动加载** | ❌ | ✅ Claude自动检测 |
| **优先级** | 较低 | 同名时skill优先 |

### 9.12 安装方式

```bash
# 插件市场
claude plugin install document-skills@anthropic-agent-skills

# 手动安装(个人级)
mkdir -p ~/.claude/skills/my-new-skill
cp SKILL.md ~/.claude/skills/my-new-skill/SKILL.md

# 手动安装(项目级, 提交git与团队共享)
mkdir -p .claude/skills/my-new-skill
cp SKILL.md .claude/skills/my-new-skill/SKILL.md
git add .claude/skills/
```

### 9.13 对 UT-Bench Skill 注入的意义

UT-Bench 的 `injectAgentNativeSkill()` 需要将 skill 复制到 `.claude/skills/<name>/SKILL.md`：

```go
// subjects.go
func injectAgentNativeSkill(workdir, skillName, skillDir string, framework string) error {
    var targetDir string
    switch framework {
    case "codebuddy":
        targetDir = filepath.Join(workdir, ".codebuddy", "skills", skillName)
    case "opencode":
        targetDir = filepath.Join(workdir, ".opencode", "skills", skillName)
    case "claudecode":
        targetDir = filepath.Join(workdir, ".claude", "skills", skillName) // 新增
    }
    // ... 复制 SKILL.md 到 targetDir
}
```

Claude Code 的 skill 系统比 CodeBuddy 更成熟：
- 支持 `allowed-tools` 预授权工具（评测中可免权限）
- 支持 `context: fork` 在子代理中隔离执行
- 支持动态上下文注入 `!`command``
- 描述预算管理（1536字符上限）

---

## 10. Hooks 系统

### 10.1 概念

Hooks 是用户定义的 shell 命令、HTTP 端点或 LLM 提示，在 Claude Code 生命周期中的特定点自动执行。

### 10.2 事件类型

| 事件 | 触发时机 | 用途 |
|:--|:--|:--|
| `PreToolUse` | 工具调用**之前** | 拦截危险操作、修改参数、注入额外上下文 |
| `PostToolUse` | 工具调用**之后** | 验证输出、自动格式化、审计日志 |
| `Notification` | Claude 需要用户注意时 | 桌面通知、声音提醒 |
| `Stop` | Claude 停止生成时（等待用户输入前） | 自动运行lint/格式化、汇总 |
| `SessionStart` | 会话开始时 | 环境检查、初始化 |
| `SessionEnd` | 会话结束时 | 清理、报告 |

### 10.3 Hook 类型

| 类型 | 说明 |
|:--|:--|
| **command** (默认) | 执行 shell 命令 |
| **http** | 调用 HTTP 端点 |
| **llm** | 调用 LLM 评估（用小模型判断是否继续） |

### 10.4 配置方式

在 `settings.json` 中配置：

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "command": "echo 'About to run bash command'"
      }
    ],
    "PostToolUse": [
      {
        "matcher": "Edit",
        "command": "prettier --write $FILEPATH"
      }
    ]
  }
}
```

### 10.5 输入输出协议

Claude Code 通过 **stdin** 传入 JSON（包含工具名、输入、会话ID等），Hook 通过 **stdout** 返回 JSON。

**PreToolUse exit code 语义:**
- `0` — 继续执行工具
- `2` — 阻止工具调用
- 其他 — 提示用户确认

**stdout 输出可修改工具输入:**

```json
{
  "decision": "allow",       // "allow" | "deny" | "prompt"
  "reason": "Auto-approved", // 拒绝/提示时显示的理由
  "tool_input": {            // 修改后的工具输入(可选)
    "command": "echo 'safe command'"
  }
}
```

### 10.6 环境变量

Hook 运行时可用的环境变量:

| 变量 | 说明 |
|:--|:--|
| `$CLAUDE_TOOL_NAME` | 工具名称 |
| `$FILEPATH` | 文件路径(Edit/Read时) |
| `$CLAUDE_SESSION_ID` | 会话 ID |
| `$CLAUDE_PROJECT_DIR` | 项目目录 |

### 10.7 对 UT-Bench 的意义

评测场景中 `--bare` 标志会跳过 hooks 发现，所以 hooks 在 UT-Bench 中**不涉及**。但如果需要自定义 Agent 行为（如每次 Bash 命令前审计），可以考虑不用 `--bare`。

---

## 11. CLAUDE.md 项目指令

### 11.1 概念

`CLAUDE.md` 是项目级指令文件，定义"这个项目是谁、怎么工作"。Claude 在每次会话开始时**始终加载**。

### 11.2 文件位置与优先级

| 文件 | 位置 | 是否提交git |
|:--|:--|:--|
| `CLAUDE.md` | 项目根目录 | ✅ 提交 |
| `.claude/CLAUDE.md` | `.claude/`目录 | ✅ 提交 |
| `CLAUDE.local.md` | 项目根目录 | ❌ 不提交(个人) |
| `.claude/CLAUDE.local.md` | `.claude/`目录 | ❌ 不提交(个人) |

按顺序加载，后加载的追加到前者的上下文中。

### 11.3 用途

- 技术栈说明（语言、框架、版本）
- 架构和代码组织约定
- 编码规范（命名、格式、测试要求）
- 常用命令（构建、测试、部署）
- 项目特有约定

### 11.4 与 `--bare` 的关系

`--bare` 标志会跳过 CLAUDE.md 加载。UT-Bench 评测中：
- 如果需要给 Agent 项目上下文 → 不用 `--bare`，提供 CLAUDE.md
- 如果只需纯任务指令 → 用 `--bare` + `--system-prompt` 或 `--append-system-prompt`

---

## 12. UT-Bench 适配要点

### 12.1 agents.yaml 框架配置 (草案)

```yaml
frameworks:
  claudecode:
    kind: agent
    description: "Anthropic Claude Code"
    command: >
      claude -p "{{.ContainerPrompt}}"
      --bare
      --output-format stream-json
      --verbose
      --max-turns 50
      --permission-mode bypassPermissions
    env:
      # --- 认证 (按模型动态注入) ---
      ANTHROPIC_BASE_URL: "{{.ModelEndpoint}}"
      ANTHROPIC_AUTH_TOKEN: "${{{.ModelAPIKeyEnv}}}"
      ANTHROPIC_MODEL: "{{.ModelID}}"
      # --- 子任务模型映射 ---
      ANTHROPIC_DEFAULT_SONNET_MODEL: "{{.ModelID}}"
      ANTHROPIC_DEFAULT_OPUS_MODEL: "{{.ModelID}}"
      ANTHROPIC_DEFAULT_HAIKU_MODEL: "{{.ModelID}}"
      CLAUDE_CODE_SUBAGENT_MODEL: "{{.ModelID}}"
      # --- 评测环境优化 ---
      API_TIMEOUT_MS: "600000"
      BASH_DEFAULT_TIMEOUT_MS: "120000"
      CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC: "1"
      CLAUDE_CODE_DISABLE_EXPERIMENTAL_BETAS: "1"
      CLAUDE_CODE_DISABLE_1M_CONTEXT: "1"
      CLAUDE_CODE_SKIP_PROMPT_HISTORY: "1"
    sandbox:
      image: "utbench:latest"
    output_globs:
      - "*.json"
      - "*.txt"
      - "*.py"
      - "*.go"
      - "*.java"
      - "*.cpp"
    preflight:
      - "claude --version"
```

### 12.2 models.yaml 模型配置 (草案)

Claude Code 接入第三方模型需要 `endpoint` 为 Anthropic 兼容接口:

```yaml
models:
  deepseek:
    id: "deepseek-v4-pro"
    endpoint: "https://api.deepseek.com/anthropic"
    api_key_env: "DEEPSEEK_API_KEY"

  qwen:
    id: "qwen3.6-plus"
    endpoint: "https://dashscope.aliyuncs.com/apps/anthropic"
    api_key_env: "DASHSCOPE_API_KEY"

  glm:
    id: "glm-5.1"
    endpoint: "https://open.bigmodel.cn/api/anthropic"
    api_key_env: "BIGMODEL_API_KEY"

  minimax:
    id: "MiniMax-M2.7"
    endpoint: "https://api.minimaxi.com/anthropic"  # 国内用minimaxi.com, 国际用minimax.io
    api_key_env: "MINIMAX_API_KEY"

  doubao:
    id: "doubao-seed-2.0-code"
    endpoint: "https://ark.cn-beijing.volces.com/api/coding"
    api_key_env: "ARK_API_KEY"
```

### 12.3 Dockerfile 改动 (Stage 1)

```dockerfile
# Stage 1: 安装 Agent CLI
FROM node:22-bookworm-slim AS agent-cli

RUN npm install -g @anthropic-ai/claude-code
# claude 命令现在可用
```

### 12.4 Go 代码适配点

| 文件 | 改动 | 说明 |
|------|------|------|
| `injectAgentNativeSkill()` | 新增 `case "claudecode"` | Skill目录: `.claude/skills/` |
| `parseAgentOutput()` | 新增Claude Code stream-json解析 | 解析`result`事件提取usage/token |
| `collectSessionExport()` | 按需新增Claude Code session导出 | Claude Code的session存储在`~/.claude/` |

### 12.5 stream-json 输出解析 (Go 伪代码)

```go
// 逐行解析 NDJSON
scanner := bufio.NewScanner(output)
lastAssistantMsg := make(map[string]*ClaudeMessage) // 按message.id去重

for scanner.Scan() {
    line := strings.TrimSpace(scanner.Text())
    if line == "" { continue }

    var event ClaudeEvent
    if err := json.Unmarshal([]byte(line), &event); err != nil { continue }

    switch event.Type {
    case "result":
        // 提取最终结果
        // event.TotalCostUSD, event.NumTurns, event.Usage.InputTokens, event.Usage.OutputTokens
    case "assistant":
        // ⚠️ 同一message.id会多次发射，需去重保留最后一次
        if event.Message != nil && event.Message.ID != "" {
            lastAssistantMsg[event.Message.ID] = event.Message
        }
        // 提取工具调用和思考过程(需--verbose)
        for _, content := range event.Message.Content {
            switch content.Type {
            case "tool_use":
                // content.Name, content.Input
            case "text":
                // content.Text
            case "thinking":
                // content.Thinking
            }
        }
    case "user":
        // ⚠️ tool_result 在这里，不在assistant中!
        for _, content := range event.Message.Content {
            if content.Type == "tool_result" {
                // content.ToolUseID, content.Content, content.IsError
            }
        }
    case "system":
        // init/start/end/api_retry 系统事件
    case "stream_event":
        // text_delta 逐token流
    }
}

type ClaudeEvent struct {
    Type          string         `json:"type"`
    Subtype       string         `json:"subtype,omitempty"`
    SessionID     string         `json:"session_id,omitempty"`
    Result        string         `json:"result,omitempty"`
    TotalCostUSD  float64        `json:"total_cost_usd,omitempty"`
    CostUSD       float64        `json:"cost_usd,omitempty"` // json格式用此字段
    NumTurns      int            `json:"num_turns,omitempty"`
    DurationMs    int64          `json:"duration_ms,omitempty"`
    DurationAPIMs int64          `json:"duration_api_ms,omitempty"`
    IsError       bool           `json:"is_error,omitempty"`
    Usage         ClaudeUsage    `json:"usage,omitempty"`
    Message       *ClaudeMessage `json:"message,omitempty"`
}

type ClaudeUsage struct {
    InputTokens              int `json:"input_tokens"`
    OutputTokens             int `json:"output_tokens"`
    CacheCreationInputTokens int `json:"cache_creation_input_tokens,omitempty"`
    CacheReadInputTokens     int `json:"cache_read_input_tokens,omitempty"`
}

type ClaudeMessage struct {
    ID      string          `json:"id"`
    Role    string          `json:"role"`
    Content []ClaudeContent `json:"content"`
}

type ClaudeContent struct {
    Type       string          `json:"type"`
    Text       string          `json:"text,omitempty"`
    Name       string          `json:"name,omitempty"`
    ToolUseID  string          `json:"tool_use_id,omitempty"`
    Input      json.RawMessage `json:"input,omitempty"`
    Content    string          `json:"content,omitempty"`      // tool_result的内容
    IsError    bool            `json:"is_error,omitempty"`     // tool_result是否有错
    Thinking   string          `json:"thinking,omitempty"`
}
```

### 12.6 Skill 注入适配 (injectAgentNativeSkill)

Claude Code 的 Skill 目录为 `.claude/skills/<name>/SKILL.md`：

```go
// subjects.go - injectAgentNativeSkill 新增 case
case "claudecode":
    targetDir = filepath.Join(workdir, ".claude", "skills", skillName)
```

**Claude Code Skill 的优势（vs CodeBuddy/OpenCode）：**

| 能力 | Claude Code | CodeBuddy | OpenCode |
|------|-------------|-----------|----------|
| `allowed-tools` 预授权 | ✅ frontmatter字段 | ❌ | ❌ |
| `context: fork` 隔离执行 | ✅ 子代理 | ❌ | ❌ |
| 动态上下文 `!`command`` | ✅ | ❌ | ❌ |
| 附属文件(模板/脚本) | ✅ 目录结构 | ✅ | ❌ |
| 描述预算管理 | ✅ 1536字符上限 | ❌ | ❌ |

### 12.7 --bare 与 Skill/CLAUDE.md 的取舍

| 场景 | 推荐标志组合 |
|------|-------------|
| 纯任务评测(只需prompt) | `--bare --append-system-prompt "..."` |
| 需要 Skill 辅助 | **不用** `--bare`，但准备 CLAUDE.md |
| 需要 Skill + 自定义prompt | 不用 `--bare` + `--append-system-prompt` + 放置 skill 到 `.claude/skills/` |
| 评测 Agent 自身能力 | `--bare` 最干净 |

| 维度 | CodeBuddy | Claude Code |
|------|-----------|-------------|
| 安装 | `npm install -g @codebuddyai/cli` | `npm install -g @anthropic-ai/claude-code` |
| 多模型配置 | 需写 `models.json` 文件 + sed替换 | 环境变量直接注入，零配置文件 |
| Headless命令 | `codebuddy -p -y --output-format json --model <id>` | `claude -p --bare --output-format stream-json` |
| 输出格式 | 单个JSON对象 | NDJSON流(更丰富的事件) |
| Token追踪 | 正则提取或JSON解析 | `usage.input_tokens`/`output_tokens` 原生字段 |
| 成本追踪 | 无原生支持 | `cost_usd` 原生字段 |
| 轮次限制 | `--max-turns` | `--max-turns`(到达时以错误退出) |
| 权限 | `-y` 自动确认 | `--permission-mode bypassPermissions` |
| Skill目录 | `.codebuddy/skills/` | `.claude/skills/` |

### 12.8 与 CodeBuddy / OpenCode 的全面对比

| 维度 | Claude Code | CodeBuddy | OpenCode |
|------|-------------|-----------|----------|
| **安装** | `npm i -g @anthropic-ai/claude-code` | `npm i -g @codebuddyai/cli` | `go install` |
| **多模型配置** | 纯env变量(零文件) | 需写`models.json`+sed | OPENCODE_CONFIG_CONTENT注入 |
| **Headless命令** | `claude -p --bare --output-format stream-json` | `codebuddy -p -y --output-format json` | `opencode --agent-mode --approval auto` |
| **输出格式** | NDJSON流(丰富事件) | 单个JSON对象 | 文本+session export |
| **Token追踪** | `usage.input/output_tokens` 原生 | 正则提取/JSON解析 | 正则提取 |
| **成本追踪** | `cost_usd` 原生 | ❌ | ❌ |
| **轮次限制** | `--max-turns`(超限错误退出) | `--max-turns` | 交互轮次参数 |
| **权限** | `--permission-mode bypassPermissions` | `-y` 自动确认 | `--approval auto` |
| **Skill目录** | `.claude/skills/` | `.codebuddy/skills/` | `.opencode/skills/` |
| **Skill能力** | allowed-tools/fork/动态注入 | 基础 | 基础 |
| **轻量启动** | `--bare` 跳过所有发现 | 无等价 | 无等价 |
| **子代理模型** | `CLAUDE_CODE_SUBAGENT_MODEL` | 单一模型 | 单一模型 |
| **模型分层** | Opus/Sonnet/Haiku三级映射 | 单一模型 | 单一模型 |

---

## 附录: 配置文件位置

| 文件 | 说明 |
|------|------|
| `~/.claude/settings.json` | 全局用户设置 |
| `.claude/settings.json` | 项目级设置 |
| `.claude/settings.local.json` | 项目级本地设置(不提交git) |
| `.claude/CLAUDE.md` | 项目级自定义指令 |
| `.claude/rules/` | 项目级规则目录 |
| `~/.claude/cache/gateway-models.json` | Gateway模型发现缓存 |
| `~/.claude/debug/<session-id>.txt` | 调试日志 |

---

> 本文档基于 Claude Code 官方文档 (https://code.claude.com/docs/) 及各厂商公开文档整理。
> 最后更新: 2026-05-05
