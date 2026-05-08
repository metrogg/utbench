# CodeBuddy Code 官方文档

> 来源: https://www.codebuddy.ai/docs/
> 抓取时间: 2026-05-04

---

## 目录

1. [快速入门](#1-快速入门)
2. [CLI 完整参考](#2-cli-完整参考)
3. [模型配置 (models.json)](#3-模型配置-modelsjson)
4. [设置配置 (settings.json)](#4-设置配置-settingsjson)
5. [环境变量](#5-环境变量)

---

## 1. 快速入门

### 系统要求
- **Node.js**: 版本 18.20 或更高
- **操作系统**: macOS、Linux 或 Windows

### 验证环境
```bash
node --version  # 应显示 v18.0.0 或更高
npm --version
```

### 安装

**方式一：npm 全局安装**
```bash
npm install -g @tencent-ai/codebuddy-code
```

**方式二：原生安装器（Beta）**
```bash
# macOS/Linux
curl -fsSL https://copilot.tencent.com/cli/install.sh | bash

# Windows
irm https://copilot.tencent.com/cli/install.ps1 | iex
```

**验证安装**
```bash
codebuddy --version
```

### 登录认证

首次使用需完成登录认证，启动后选择登录方式：
- **中文站点 (Chinese Site)**: 通过腾讯云中国站认证 (copilot.tencent.com)
- **国际站点 (International Site)**: 通过腾讯云国际站认证 (codebuddy.ai)
- **企业域名 (Enterprise Domain)**: 连接企业专用或自托管 CodeBuddy 服务
- **iOA**: 腾讯内部员工通过 iOA 零信任系统认证

### 初次体验

```bash
cd /path/to/your/project
codebuddy

# 建议先初始化项目上下文
> /init

# 开始对话
> Help me analyze the structure of this project
```

### 核心使用模式

**交互式对话模式：**
```bash
codebuddy
```

**单命令模式（脚本/自动化）：**
```bash
# 直接查询
codebuddy -p "Optimize the performance of this SQL query"

# 管道输入
cat error.log | codebuddy -p "Analyze these error logs"

# 文件分析（需要授权时必须加 -y）
codebuddy -p "Review the code quality of src/utils.js" -y
```

> **重要提示**: 使用 `-p` 时，如果操作需要授权（如文件访问或命令执行），必须添加 `-y`。

### 快捷键

| 快捷键 | 功能 |
|--------|------|
| `↑/↓` | 浏览命令历史 |
| `Tab` | 命令自动补全 |
| `Esc` | 清除输入（按两次）/ 返回上一级 |
| `Ctrl+C` | 退出程序 |
| `Ctrl+D` | 退出程序 |
| `Enter` | 发送消息 |
| `Shift+Enter` | 换行（多行输入） |
| `Ctrl+G` | 在外部编辑器中编辑 prompt |
| `Ctrl+R` | 展开/折叠详细输出 |
| `Shift+Tab` (macOS) / `Alt+M` (Win) | 切换权限模式 |
| `Ctrl+O` | 查看/关闭思考详情面板 |
| `k` | 终止选中的后台任务 |

---

## 2. CLI 完整参考

### 核心命令

| 命令 | 描述 | 示例 |
|------|------|------|
| `codebuddy` | 启动交互式 REPL | `codebuddy` |
| `codebuddy "query"` | 带初始提示启动 REPL | `codebuddy "explain this project"` |
| `codebuddy -p "query"` | 通过 SDK 查询后退出 | `codebuddy -p "explain this function"` |
| `cat file \| codebuddy -p "query"` | 处理管道输入 | `cat logs.txt \| codebuddy -p "analyze logs"` |
| `codebuddy -c` | 继续最近的对话 | `codebuddy -c` |
| `codebuddy -c -p "query"` | 通过 SDK 继续对话 | `codebuddy -c -p "check for type errors"` |
| `codebuddy -r "<session-id>" "query"` | 按 ID 恢复会话 | `codebuddy -r "abc123" "complete this MR"` |
| `codebuddy update` | 更新到最新版本 | `codebuddy update` |
| `codebuddy mcp` | 配置 MCP 服务器 | 见 MCP 文档 |

### Daemon 命令

| 命令 | 描述 | 示例 |
|------|------|------|
| `codebuddy daemon start` | 启动 Daemon 进程 | `codebuddy daemon start --port 8080` |
| `codebuddy daemon stop` | 停止 Daemon | `codebuddy daemon stop` |
| `codebuddy daemon status` | 查看 Daemon 状态 | `codebuddy daemon status` |
| `codebuddy daemon restart` | 重启 Daemon | `codebuddy daemon restart` |
| `codebuddy ps` | 列出所有活跃 Worker 进程 | `codebuddy ps` |
| `codebuddy logs <pid\|name>` | 查看 Worker 日志 | `codebuddy logs feature-x` |
| `codebuddy attach <pid\|name>` | 附着到后台 Worker | `codebuddy attach feature-x` |
| `codebuddy kill <pid\|name>` | 终止 Worker 进程 | `codebuddy kill feature-x` |

### 所有命令行参数

#### 输入/输出控制

| 参数 | 描述 | 示例 |
|------|------|------|
| `-p` / `--print` | 打印响应后退出，不进入交互模式 | `codebuddy -p "query"` |
| `--output-format` | 输出格式：`text`、`json`、`stream-json` | `codebuddy -p "query" --output-format json` |
| `--input-format` | 输入格式：`text`、`stream-json` | `codebuddy -p --output-format json --input-format stream-json` |
| `--json-schema` | JSON Schema 验证结构化输出 | `codebuddy -p --output-format json --json-schema '{...}' "query"` |
| `--include-partial-messages` | 包含部分流事件（需 `--print` + `stream-json`） | `codebuddy -p --output-format stream-json --include-partial-messages "query"` |

#### 模型控制

| 参数 | 描述 | 示例 |
|------|------|------|
| `--model` | 设置模型（别名或完整模型名） | `codebuddy --model gpt-5` |
| `--text-to-image-model` | 文本到图像生成模型 ID | `codebuddy --text-to-image-model your-image-model` |
| `--image-to-image-model` | 图像到图像编辑模型 ID | `codebuddy --image-to-image-model your-edit-model` |

#### 会话控制

| 参数 | 描述 | 示例 |
|------|------|------|
| `-c` / `--continue` | 加载当前目录中最近的对话 | `codebuddy --continue` |
| `--resume` | 按 ID 恢复特定会话 | `codebuddy --resume abc123 "query"` |
| `--max-turns` | 限制非交互模式中的 agent 轮次 | `codebuddy -p --max-turns 3 "query"` |

#### 权限控制

| 参数 | 描述 | 示例 |
|------|------|------|
| `-y` / `--dangerously-skip-permissions` | 跳过权限提示（谨慎使用） | `codebuddy -y` |
| `--permission-mode` | 以指定权限模式启动（如 `plan`） | `codebuddy --permission-mode plan` |
| `--subagent-permission-mode` | 设置 subagent/team 成员默认权限模式 | `codebuddy --subagent-permission-mode bypassPermissions` |
| `--permission-prompt-tool` | MCP 工具处理非交互模式的权限提示 | `codebuddy -p --permission-prompt-tool mcp_auth_tool "query"` |

#### 工具控制

| 参数 | 描述 | 示例 |
|------|------|------|
| `--allowedTools` | 允许无需提示的工具列表 | `"Bash(git log:*)" "Read"` |
| `--disallowedTools` | 禁止的工具列表 | `"Bash(git log:*)" "Edit"` |
| `--tools` | 限制内置工具集（白名单） | `codebuddy --tools "Bash,Read,Edit"` |

#### 系统提示词控制（3 种方式）

| 参数 | 行为 | 可用模式 |
|------|------|----------|
| `--system-prompt` | **替换**整个默认系统提示 | 交互 + Print |
| `--system-prompt-file` | **替换**为文件内容 | 仅 Print |
| `--append-system-prompt` | **追加**到默认提示末尾 | 交互 + Print |

> `--system-prompt` 和 `--system-prompt-file` 互斥。推荐 `--append-system-prompt`。

#### Agent 控制

| 参数 | 描述 | 示例 |
|------|------|------|
| `--agent` | 指定 agent 名称 | `codebuddy --agent my-reviewer` |
| `--agents` | 通过 JSON 动态定义自定义 Sub-Agents | 见下方格式 |

**Sub-Agent JSON 格式：**
```json
{
  "code-reviewer": {
    "description": "何时调用此 Sub-Agent 的自然语言描述",
    "prompt": "引导 Sub-Agent 行为的系统提示",
    "tools": ["Read", "Edit", "Bash"],
    "model": "sonnet"
  }
}
```

#### 配置控制

| 参数 | 描述 | 示例 |
|------|------|------|
| `--settings` | 从 JSON 文件或字符串加载额外设置 | `codebuddy --settings '{"model":"gpt-5"}' "query"` |
| `--setting-sources` | 指定设置源：`user`、`project`、`local` | `codebuddy --setting-sources project,local "query"` |

#### 目录控制

| 参数 | 描述 | 示例 |
|------|------|------|
| `--add-dir` | 添加额外工作目录 | `codebuddy --add-dir ../apps ../lib` |

#### 调试

| 参数 | 描述 |
|------|------|
| `--verbose` | 详细日志，显示完整轮次输出 |
| `--debug` | 调试模式，可选类别过滤 |

#### IDE 集成

| 参数 | 描述 |
|------|------|
| `--ide` | 自动连接 IDE |

#### 沙盒模式（Beta）

| 参数 | 描述 | 示例 |
|------|------|------|
| `--sandbox [url]` | 在沙盒运行：无参数用容器；提供 URL 用 E2B | `codebuddy --sandbox "analyze project"` |
| `--sandbox-upload-dir` | 上传当前目录到沙盒（仅 E2B） | |
| `--sandbox-new` | 强制创建新沙盒 | `codebuddy --sandbox --sandbox-new "start"` |
| `--sandbox-id <id>` | 连接到指定沙盒 | `codebuddy --sandbox --sandbox-id sb_abc123` |
| `--sandbox-kill` | 退出时终止沙盒 | `codebuddy --sandbox --sandbox-kill` |

#### Git Worktree

| 参数 | 描述 |
|------|------|
| `--worktree [name]` | 在隔离的 git worktree 中运行 |
| `--tmux` | 在 tmux 会话中运行（配合 `--worktree`） |

#### 插件

| 参数 | 描述 |
|------|------|
| `--plugin-dir <dirs...>` | 从本地目录加载插件，支持多路径 |

#### 后台/服务模式

| 参数 | 描述 |
|------|------|
| `--bg` | 后台运行，日志输出到 `~/.codebuddy/logs/` |
| `--name <name>` | 后台会话名称（配合 `--bg`） |
| `--serve` | 启动 HTTP 服务器（Web UI/REST API/ACP） |

### Headless/Pipeline 典型用法

```bash
# 基本查询
codebuddy -p "explain this function"

# 自动化必需：跳过权限
codebuddy -p -y "refactor this code"

# 限制轮次（控制成本/时间）
codebuddy -p --max-turns 3 "fix all linting errors"

# JSON 输出供脚本解析
codebuddy -p "analyze code" --output-format json

# 结构化输出 + Schema 验证
codebuddy -p --output-format json --json-schema '{"type":"object",...}' "find all bugs"

# 自定义模型
codebuddy -p --model gpt-5 "explain this"

# 后台运行 + 管理
codebuddy --bg --name feature-x "implement feature"
codebuddy logs feature-x
codebuddy kill feature-x
```

---

## 3. 模型配置 (models.json)

### 配置文件位置

| 级别 | 路径 | 优先级 |
|------|------|--------|
| 用户级 | `~/.codebuddy/models.json` | 低 |
| 项目级 | `<project-root>/.codebuddy/models.json` | 高 |

> 项目级配置会覆盖用户级同名模型配置（基于 `id` 字段匹配）。
> `availableModels` 项目级完全覆盖用户级，不进行合并。

### 完整字段说明

#### `models` — `Array<LanguageModel>`

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `id` | string | ✅ | 模型唯一标识符 |
| `name` | string | - | 模型显示名称 |
| `vendor` | string | - | 模型供应商 |
| `apiKey` | string | - | API 密钥，支持环境变量引用 `${VAR_NAME}` |
| `maxInputTokens` | number | - | 最大输入 token 数 |
| `maxOutputTokens` | number | - | 最大输出 token 数 |
| `url` | string | - | API 端点 URL，**必须以 `/chat/completions` 结尾** |
| `temperature` | number | - | 采样温度，范围 0-2 |
| `supportsToolCall` | boolean | - | 是否支持工具调用 |
| `supportsImages` | boolean | - | 是否支持图片输入 |
| `supportsReasoning` | boolean | - | 是否支持推理模式 |
| `relatedModels` | object | - | 相关模型配置（`lite`/`reasoning`） |

#### `availableModels` — `Array<string>`

控制模型下拉列表中显示哪些模型。未配置或空数组 → 显示所有模型。

### 自定义模型端点 URL 配置

**关键规则：** 仅支持 OpenAI 接口格式 API，`url` 必须以 `/chat/completions` 结尾。

✅ 正确：`https://api.openai.com/v1/chat/completions`
❌ 错误：`https://api.openai.com/v1`

### API Key 配置方式

**推荐：环境变量引用**
```json
{ "apiKey": "${OPENAI_API_KEY}" }
```

设置环境变量：
```bash
export OPENAI_API_KEY="sk-your-actual-api-key"
# 或
OPENAI_API_KEY="sk-xxx" codebuddy
```

### relatedModels 相关模型配置

| 场景 | 用途 | 状态 |
|------|------|------|
| `lite` | 轻量快速模型（后台提取、摘要） | 已激活 |
| `reasoning` | 推理增强模型（复杂推理） | 已激活 |
| `subagent` | 子代理默认模型 | 保留未激活 |
| `vision` | 视觉理解模型 | 保留未激活 |
| `longContext` | 长上下文模型 | 保留未激活 |

**解析优先级：**
1. 环境变量（`CODEBUDDY_SMALL_FAST_MODEL` → lite，`CODEBUDDY_BIG_SLOW_MODEL` → reasoning）
2. 主模型条目的 `relatedModels[variant]`
3. 内置 `defaultRelatedModels`（仅内置模型生效）
4. 回退到主模型自身

### 完整配置示例

**基础配置（云端 + 本地）：**
```json
{
  "models": [
    {
      "id": "gpt-4o",
      "name": "GPT-4o",
      "vendor": "OpenAI",
      "apiKey": "sk-your-openai-key",
      "maxInputTokens": 128000,
      "maxOutputTokens": 16384,
      "supportsToolCall": true,
      "supportsImages": true
    },
    {
      "id": "my-local-llm",
      "name": "My Local LLM",
      "vendor": "Ollama",
      "url": "http://localhost:11434/v1/chat/completions",
      "apiKey": "ollama",
      "maxInputTokens": 8192,
      "maxOutputTokens": 2048,
      "supportsToolCall": true
    }
  ],
  "availableModels": ["gpt-4o", "my-local-llm"]
}
```

**环境变量 + relatedModels（推荐生产配置）：**
```json
{
  "models": [
    {
      "id": "deepseek-v4-pro",
      "name": "DeepSeek V4 Pro",
      "vendor": "DeepSeek",
      "url": "https://api.deepseek.com/v1/chat/completions",
      "apiKey": "${DEEPSEEK_API_KEY}",
      "maxInputTokens": 128000,
      "maxOutputTokens": 8192,
      "supportsToolCall": true,
      "relatedModels": {
        "lite": "deepseek-v4-flash",
        "reasoning": "deepseek-v4-pro"
      }
    },
    {
      "id": "deepseek-v4-flash",
      "name": "DeepSeek V4 Flash",
      "vendor": "DeepSeek",
      "url": "https://api.deepseek.com/v1/chat/completions",
      "apiKey": "${DEEPSEEK_API_KEY}",
      "maxInputTokens": 128000,
      "maxOutputTokens": 8192,
      "supportsToolCall": true
    }
  ],
  "availableModels": ["deepseek-v4-pro", "deepseek-v4-flash"]
}
```

**OpenRouter 配置：**
```json
{
  "models": [
    {
      "id": "openai/gpt-4o",
      "name": "open-router-model",
      "url": "https://openrouter.ai/api/v1/chat/completions",
      "apiKey": "sk-or-v1-your-openrouter-api-key",
      "maxInputTokens": 128000,
      "maxOutputTokens": 4096,
      "supportsToolCall": true,
      "supportsImages": false
    }
  ]
}
```

**热重载：** 配置文件修改后自动检测（1秒防抖延迟），无需重启。

---

## 4. 设置配置 (settings.json)

### 配置文件层级（优先级从高到低）

| 优先级 | 层级 | 文件路径 | 用途 |
|--------|------|----------|------|
| 1 (最高) | 命令行参数 | `codebuddy config set ...` | 临时会话覆盖 |
| 2 | 本地项目设置 | `.codebuddy/settings.local.json` | 个人项目偏好（不提交 Git） |
| 3 | 共享项目设置 | `.codebuddy/settings.json` | 团队共享（提交到源码管理） |
| 4 (最低) | 用户全局设置 | `~/.codebuddy/settings.json` | 所有项目通用的个人偏好 |

**合并规则：** 高层配置覆盖低层配置；按 key 合并。

### 完整 JSON Schema

```json
{
  // ===== 核心设置 =====
  "language": "English",
  "model": "gpt-5",
  "agent": "my-reviewer",
  "textToImageModel": "your-image-model",
  "imageToImageModel": "your-edit-model",
  "endpoint": "https://api.example.com",
  "envRouteMode": "production",
  "reasoningEffort": "high",
  "apiKeyHelper": "/bin/generate_temp_api_key.sh",

  // ===== 环境变量 =====
  "env": {
    "NODE_ENV": "development",
    "DEBUG": "codebuddy:*"
  },

  // ===== Git =====
  "includeCoAuthoredBy": false,

  // ===== 聊天历史 =====
  "cleanupPeriodDays": 30,

  // ===== UI =====
  "showTokensCounter": false,
  "statusLine": {
    "type": "command",
    "command": "~/.codebuddy/statusline.sh"
  },

  // ===== 行为开关 =====
  "autoCompactEnabled": true,
  "autoUpdates": false,
  "alwaysThinkingEnabled": true,
  "promptSuggestionEnabled": false,
  "disableAllHooks": true,

  // ===== 目录信任 =====
  "trustedDirectories": ["~/workspace/myproj"],
  "trustAll": true,

  // ===== MCP 服务器 =====
  "enableAllProjectMcpServers": false,
  "enabledMcpjsonServers": ["memory", "github"],
  "disabledMcpjsonServers": ["filesystem"],

  // ===== 权限 =====
  "permissions": {
    "allow": ["Bash(npm run lint)", "Bash(npm run test:*)", "Read(~/.zshrc)"],
    "ask": ["Bash(git push:*)"],
    "deny": ["Bash(curl:*)", "Read(./.env)", "Read(./.env.*)", "Read(./secrets/**)"],
    "additionalDirectories": ["../docs/"],
    "defaultMode": "acceptEdits",
    "disableBypassPermissionsMode": "disable",
    "subagentPermissionMode": "bypassPermissions"
  },

  // ===== Hooks =====
  "hooks": {
    "PreToolUse": { "Bash": "echo 'Running command...'" }
  },

  // ===== 沙箱 =====
  "sandbox": {
    "enabled": true,
    "autoAllowBashIfSandboxed": true,
    "excludedCommands": ["git", "docker"],
    "allowUnsandboxedCommands": true,
    "enableWeakerNestedSandbox": true,
    "network": {
      "allowUnixSockets": ["~/.ssh/agent-socket"],
      "allowLocalBinding": true,
      "httpProxyPort": 8080,
      "socksProxyPort": 8081
    }
  },

  // ===== 内存（实验性）=====
  "memory": {
    "autoMemoryEnabled": true,
    "typedMemory": true,
    "relevanceSelection": true,
    "memoryExtraction": false,
    "teamMemory": { "enabled": true, "userId": "username" }
  },

  // ===== 插件 =====
  "enabledPlugins": {
    "formatter@company-tools": true,
    "deployer@company-tools": true
  },
  "extraKnownMarketplaces": {
    "company-tools": {
      "source": { "source": "github", "repo": "company/codebuddy-plugins" }
    }
  }
}
```

### 核心设置字段说明

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `language` | string | 空(自动检测) | 首选响应语言 |
| `model` | string | 产品默认 | 覆盖默认模型 |
| `agent` | string | 产品默认 | 覆盖主线程 agent |
| `endpoint` | string | - | 自定义服务器端点地址 |
| `reasoningEffort` | string | 产品默认 | 推理深度：`low`/`medium`/`high`/`xhigh` |
| `apiKeyHelper` | string | - | 自定义认证脚本路径 |

### 权限设置说明

| 字段 | 说明 |
|------|------|
| `allow` | 允许的工具使用规则（Bash 使用前缀匹配） |
| `ask` | 需确认的工具使用规则 |
| `deny` | 拒绝的工具使用规则 |
| `defaultMode` | 默认权限模式：`"default"` / `"acceptEdits"` / `"bypassPermissions"` |
| `disableBypassPermissionsMode` | 设 `"disable"` 可阻止 `-y` 和 `--dangerously-skip-permissions` |
| `subagentPermissionMode` | 子代理权限模式覆盖 |

### 可用工具列表

| 工具 | 说明 | 需要权限 |
|------|------|----------|
| AskUserQuestion | 向用户提问（多选题） | 否 |
| Bash | 执行 shell 命令 | **是** |
| Edit | 目标化文件编辑 | **是** |
| MultiEdit | 单文件多编辑操作 | **是** |
| Write | 创建或覆盖文件 | **是** |
| NotebookEdit | 修改 Jupyter notebook | **是** |
| WebFetch | 获取 URL 内容 | **是** |
| WebSearch | 带域名过滤的网络搜索 | **是** |
| Read | 读取文件内容 | 否 |
| Glob | 模式匹配查找文件 | 否 |
| Grep | 文件内容搜索 | 否 |
| LSP | LSP 代码智能 | 否 |
| Task | 运行子代理 | 否 |

### CLI 配置管理命令

```bash
codebuddy config get <key>              # 获取配置值
codebuddy config set <key> <value>      # 设置项目级配置
codebuddy config set -g <key> <value>   # 设置全局配置
codebuddy config list                   # 列出所有配置
codebuddy config add <key> <values...>  # 向数组配置添加项
codebuddy config remove <key> [values...] # 删除配置或数组项
```

---

## 5. 环境变量

### Authentication（认证）

| 环境变量 | 描述 |
|----------|------|
| `CODEBUDDY_API_KEY` | API 密钥，非交互模式（`-p`）下始终使用 |
| `CODEBUDDY_AUTH_TOKEN` | CodeBuddy 平台认证令牌 |
| `CODEBUDDY_CUSTOM_HEADERS` | 自定义 HTTP 请求头（格式：`Name: Value`，多个用换行符分隔） |

### API Endpoints and Proxy（API 端点与代理）

| 环境变量 | 描述 |
|----------|------|
| `CODEBUDDY_BASE_URL` | 覆盖 API 端点地址 |
| `CODEBUDDY_INTERNET_ENVIRONMENT` | 网络环境（`internal` 中国版 / `ioa` 企业版） |
| `HTTP_PROXY` / `http_proxy` | HTTP 代理地址 |
| `HTTPS_PROXY` / `https_proxy` | HTTPS 代理地址 |
| `NO_PROXY` / `no_proxy` | 绕过代理的域名列表 |

### Model Configuration（模型配置）

| 环境变量 | 描述 |
|----------|------|
| `CODEBUDDY_MODEL` | 覆盖默认 Agent 模型 |
| `CODEBUDDY_SMALL_FAST_MODEL` | 后台任务的小型快速模型 |
| `CODEBUDDY_BIG_SLOW_MODEL` | 复杂推理任务的大型模型 |
| `CODEBUDDY_CODE_SUBAGENT_MODEL` | 子 Agent 使用的模型 |
| `MAX_THINKING_TOKENS` | 扩展思维的 token 预算（默认禁用） |

### Bash Tool Configuration

| 环境变量 | 默认值 | 描述 |
|----------|--------|------|
| `BASH_DEFAULT_TIMEOUT_MS` | `120000` | bash 命令默认超时（ms） |
| `BASH_MAX_OUTPUT_LENGTH` | `30000` | bash 输出最大字符数 |
| `BASH_MAX_TIMEOUT_MS` | `600000` | 模型可设置的最大超时（ms） |

### Tools and Feature Toggles

| 环境变量 | 描述 |
|----------|------|
| `CODEBUDDY_DISABLE_HOT_RELOAD` | 设 `1` 禁用热重载 |
| `CODEBUDDY_SKIP_BUILTIN_MARKETPLACE` | 设 `1` 跳过内置插件市场 |
| `CODEBUDDY_PLUGIN_DIRS` | 本地插件目录（冒号分隔） |
| `CODEBUDDY_IMAGE_GEN_ENABLED` | 设 `false`/`0` 禁用图像生成 |
| `CODEBUDDY_DEFER_TOOL_LOADING` | 设 `false`/`0` 禁用 MCP 工具延迟加载 |
| `CODEBUDDY_DISABLE_CRON` | 设 `1` 禁用定时任务 |
| `CODEBUDDY_TOOL_RESULT_THRESHOLD_KB` | 非 bash 工具结果外部化阈值（默认 `50` KB） |

### Context and Memory

| 环境变量 | 默认值 | 描述 |
|----------|--------|------|
| `CODEBUDDY_AUTOCOMPACT_PCT_OVERRIDE` | 产品配置 | 触发自动压缩的上下文百分比 |
| `CODEBUDDY_PRE_MESSAGE_COMPACT_PCT` | `10` | 消息前压缩检查百分比 |
| `CODEBUDDY_DISABLE_AUTO_MEMORY` | 启用 | 设 `1` 禁用自动内存 |
| `CODEBUDDY_MEMORY_ENABLED` | 禁用 | 设 `true`/`1` 启用内存 |
| `CODEBUDDY_TYPED_MEMORY_ENABLED` | 禁用 | 设 `true`/`1` 启用类型化内存 |
| `CODEBUDDY_TEAM_MEMORY_ENABLED` | 禁用 | 设 `true`/`1` 启用团队内存 |

### MCP

| 环境变量 | 默认值 | 描述 |
|----------|--------|------|
| `MCP_TIMEOUT` | 无 | MCP 服务器连接超时（ms） |
| `MCP_TOOL_TIMEOUT` | 无 | MCP 工具执行超时（ms） |
| `MAX_MCP_OUTPUT_TOKENS` | `20000` | MCP 工具响应最大 token 数 |

### Performance and Output

| 环境变量 | 默认值 | 描述 |
|----------|--------|------|
| `CODEBUDDY_CODE_MAX_OUTPUT_TOKENS` | 无 | 请求最大输出 token 数 |
| `CODEBUDDY_CODE_FILE_READ_MAX_OUTPUT_TOKENS` | `20000` | 文件读取 token 限制 |
| `CODEBUDDY_STREAM_TIMEOUT_MS` | `120000` | 流式响应最大静默时间（ms） |
| `CODEBUDDY_FIRST_TOKEN_TIMEOUT_MS` | `120000` | 等待首个输出的最大时间（ms） |
| `CODEBUDDY_SESSION_MAX_ITEMS` | `1000` | 会话回放最大历史消息数 |

### Filesystem and Configuration

| 环境变量 | 描述 |
|----------|------|
| `CODEBUDDY_CONFIG_DIR` | 配置和数据文件自定义位置 |
| `CODEBUDDY_CODE_DEBUG_LOGS_DIR` | 调试日志目录 |
| `CODEBUDDY_SANDBOX_IMAGE` | 容器沙箱镜像（默认 `node:20-alpine`） |

### Shell Configuration

| 环境变量 | 描述 |
|----------|------|
| `CODEBUDDY_CODE_SHELL` | 覆盖 Shell 检测（`bash`/`zsh`/`sh`/`powershell`） |
| `CODEBUDDY_CODE_SHELL_PREFIX` | 包装所有 Shell 命令的前缀 |
| `CODEBUDDY_CODE_GIT_BASH_PATH` | Windows 上 Git Bash 路径 |
| `CODEBUDDY_POWERSHELL_PATH` | PowerShell 可执行文件路径 |
| `CODEBUDDY_ENV_FILE` | 每个 Shell 命令前自动 source 的环境文件 |

### Agent Execution Control

| 环境变量 | 默认值 | 描述 |
|----------|--------|------|
| `CODEBUDDY_CODE_MAX_TURNS` | `500` | 主 Agent 最大执行轮次 |
| `CODEBUDDY_CODE_SUBAGENT_MAX_TURNS` | `500` | 子 Agent 最大执行轮次 |
| `CODEBUDDY_SUBAGENT_PERMISSION_MODE` | 映射表默认 | 子 Agent 默认权限模式 |

### Telemetry and Debugging

| 环境变量 | 描述 |
|----------|------|
| `DISABLE_TELEMETRY` | 设 `1` 禁用遥测 |
| `DISABLE_ERROR_REPORTING` | 设 `1` 禁用错误报告 |
| `DISABLE_AUTOUPDATER` | 设 `1` 禁用自动更新 |
| `CODEBUDDY_DEBUG` | 设 `1` 启用调试模式 |
| `CODEBUDDY_DEBUG_REQUEST` | 设 `1` 启用请求调试 |

### Gateway and Remote Access

| 环境变量 | 描述 |
|----------|------|
| `CODEBUDDY_GATEWAY_AUTH` | 网关认证模式（`password`/`none`） |
| `CODEBUDDY_GATEWAY_PASSWORD` | 网关访问密码 |
| `CODEBUDDY_GATEWAY_FORCE_TUNNEL` | 设 `1` 强制隧道模式 |
| `SERVER__HOST` | `--serve` 监听地址（默认 `127.0.0.1`） |
| `SERVER__PORT` | `--serve` 监听端口 |

### E2E Testing - Record/Replay

| 环境变量 | 描述 |
|----------|------|
| `CODEBUDDY_RECORD_DIR` | 录制模式：响应保存目录 |
| `CODEBUDDY_REPLAY_DIR` | 回放模式：录制读取目录 |
| `CODEBUDDY_REPLAY_SPEED` | 回放速度倍率（`0`=即时） |
| `CODEBUDDY_REPLAY_STRICT` | 设 `1` 录制用尽时抛出错误 |

> **注意：** `RECORD_DIR` 和 `REPLAY_DIR` 互斥。
