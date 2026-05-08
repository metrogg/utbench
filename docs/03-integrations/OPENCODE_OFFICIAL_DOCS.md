# OpenCode 官方文档

> 来源: https://opencode.ai/docs/zh-cn/
> 抓取时间: 2026-05-04

---

## 目录

1. [CLI 参考](#1-cli-参考)
2. [配置文件 (opencode.json)](#2-配置文件-opencodejson)
3. [Provider 配置](#3-provider-配置)
4. [环境变量](#4-环境变量)

---

## 1. CLI 参考

### 默认行为

运行 `opencode` 不带参数启动 TUI（Terminal User Interface）。

```bash
opencode                        # 启动 TUI
opencode run "Explain closures" # 非交互模式
```

### TUI 默认启动 (`opencode [project]`)

| Flag | Short | 描述 |
|------|-------|------|
| `--continue` | `-c` | 继续最近的会话 |
| `--session` | `-s` | 指定会话 ID 继续 |
| `--fork` | — | 继续时 fork 会话 |
| `--prompt` | — | 使用指定提示 |
| `--model` | `-m` | 模型，格式：`provider/model` |
| `--agent` | — | 指定 agent |
| `--port` | — | 监听端口 |
| `--hostname` | — | 监听主机名 |

### `opencode run [message..]` — 非交互模式

**主要用途：** 脚本、自动化、无需完整 TUI 的快速问答。

| Flag | Short | 描述 |
|------|-------|------|
| `--command` | — | 要执行的命令 |
| `--continue` | `-c` | 继续最近的会话 |
| `--session` | `-s` | 会话 ID 继续 |
| `--fork` | — | 继续时 fork |
| `--share` | — | 分享会话 |
| `--model` | `-m` | 模型，格式：`provider/model` |
| `--agent` | — | 指定 agent |
| `--file` | `-f` | 附加到消息的文件 |
| `--format` | — | 输出格式：`default`（格式化）或 `json`（原始 JSON 事件） |
| `--title` | — | 会话标题 |
| `--attach` | — | 连接到运行中的 opencode 服务器 |
| `--port` | — | 本地服务器端口（默认随机） |

**Pipeline/Headless 示例：**
```bash
# 非交互运行
opencode run Explain the use of context in Go

# 指定模型
opencode run -m anthropic/claude-3-opus "Refactor this code"

# JSON 输出
opencode run --format json "Explain async/await in JavaScript"

# 连接到运行中的服务器（避免 MCP 冷启动）
opencode serve                          # 终端 1
opencode run --attach http://localhost:4096 "Explain closures"  # 终端 2

# 继续会话
opencode run --continue "Follow up question"
opencode run --session abc123 "More questions"
```

### `opencode serve` — HTTP API 服务器

启动 HTTP 服务器提供 API 访问。设置 `OPENCODE_SERVER_PASSWORD` 启用 HTTP basic auth（默认用户名：`opencode`）。

| Flag | 描述 |
|------|------|
| `--port` | 监听端口 |
| `--hostname` | 监听主机名 |
| `--mdns` | 启用 mDNS 发现 |
| `--cors` | 额外 CORS 来源 |

### `opencode attach [url]` — 远程 TUI

| Flag | Short | 描述 |
|------|-------|------|
| `--dir` | — | 启动 TUI 的工作目录 |
| `--session` | `-s` | 会话 ID 继续 |

```bash
opencode web --port 4096 --hostname 0.0.0.0   # 终端 1
opencode attach http://10.20.30.40:4096        # 终端 2
```

### `opencode web` — Web UI 服务器

启动 HTTP 服务器 + 浏览器 Web 界面。

| Flag | 描述 |
|------|------|
| `--port` | 监听端口 |
| `--hostname` | 监听主机名 |
| `--mdns` | 启用 mDNS 发现 |
| `--cors` | 额外 CORS 来源 |

### `opencode agent` — Agent 管理

| 子命令 | 描述 |
|--------|------|
| `opencode agent create` | 创建新 agent（交互式向导） |
| `opencode agent list` | 列出所有可用 agent |

### `opencode auth` — Provider 认证

| 子命令 | 描述 |
|--------|------|
| `opencode auth login` | 配置 API 密钥（基于 Models.dev） |
| `opencode auth list` / `ls` | 列出已认证的 provider |
| `opencode auth logout` | 清除凭据 |

> 凭据存储：`~/.local/share/opencode/auth.json`
> 启动时加载 auth 文件 + 环境变量 + `.env` 文件的密钥。

### `opencode models [provider]` — 列出模型

| Flag | 描述 |
|------|------|
| `--refresh` | 从 models.dev 刷新模型缓存 |
| `--verbose` | 详细输出（含成本元数据） |

```bash
opencode models              # 所有模型
opencode models anthropic    # 按 provider 过滤
opencode models --refresh    # 更新缓存
```

### `opencode mcp` — MCP 服务器管理

| 子命令 | 描述 |
|--------|------|
| `opencode mcp add` | 添加本地或远程 MCP 服务器 |
| `opencode mcp list` / `ls` | 列出 MCP 服务器 + 连接状态 |
| `opencode mcp auth [name]` | OAuth MCP 服务器认证 |
| `opencode mcp auth list` / `ls` | 列出 OAuth MCP 服务器 |
| `opencode mcp logout [name]` | 移除 OAuth 凭据 |
| `opencode mcp debug <name>` | 调试 OAuth 连接 |

### `opencode session` — 会话管理

| 子命令 | 描述 |
|--------|------|
| `opencode session list` | 列出所有会话 |

`session list` 的 Flags：
| Flag | Short | 描述 |
|------|-------|------|
| `--max-count` | `-n` | 限制显示数量 |
| `--format` | — | 输出格式：`table`（默认）或 `json` |

### `opencode stats` — Token 使用统计

| Flag | 描述 |
|------|------|
| `--days` | 显示最近 N 天的统计 |
| `--tools` | 显示的工具数量 |
| `--models` | 模型使用分布（传入数字显示 Top N） |
| `--project` | 按项目过滤（空字符串 = 当前项目） |

### `opencode export [sessionID]` — 导出会话

无 ID 时交互式选择。

### `opencode import <file>` — 导入会话

支持本地文件或 OpenCode 分享链接：
```bash
opencode import session.json
opencode import https://opncd.ai/s/abc123
```

### `opencode github` — GitHub Agent

| 子命令 | Flag | 描述 |
|--------|------|------|
| `opencode github install` | — | 在仓库安装 GitHub Agent |
| `opencode github run` | `--event`, `--token` | 运行 GitHub Agent（通常在 Actions 中） |

### `opencode acp` — Agent Client Protocol 服务器

通过 stdin/stdout 使用 nd-JSON 通信。

| Flag | 描述 |
|------|------|
| `--cwd` | 工作目录 |
| `--port` | 监听端口 |
| `--hostname` | 监听主机名 |

### `opencode uninstall`

| Flag | Short | 描述 |
|------|-------|------|
| `--keep-config` | `-c` | 保留配置文件 |
| `--keep-data` | `-d` | 保留会话数据和快照 |
| `--dry-run` | — | 预览删除而不实际删除 |
| `--force` | `-f` | 跳过确认提示 |

### `opencode upgrade [target]`

| Flag | Short | 描述 |
|------|-------|------|
| `--method` | `-m` | 安装方式：`curl`, `npm`, `pnpm`, `bun`, `brew` |

### 全局 Flags

| Flag | Short | 描述 |
|------|-------|------|
| `--help` | `-h` | 显示帮助 |
| `--version` | `-v` | 打印版本 |
| `--print-logs` | — | 输出日志到 stderr |
| `--log-level` | — | 日志级别：`DEBUG`, `INFO`, `WARN`, `ERROR` |

### Headless 模式总结

| 模式 | 命令 | 用途 |
|------|------|------|
| 单次运行 | `opencode run "prompt"` | 脚本集成 |
| 持久服务器 | `opencode serve` | API 服务器 |
| 服务器 + 附加运行 | `opencode serve` → `opencode run --attach URL` | 避免 MCP 冷启动 |
| Web UI | `opencode web` | 浏览器界面 |
| 远程 TUI | `opencode attach URL` | 远程终端 |
| ACP 协议 | `opencode acp` | stdin/stdout nd-JSON |
| GitHub Actions | `opencode github run` | CI/CD 自动化 |

---

## 2. 配置文件 (opencode.json)

### 配置文件格式

- **支持格式：** JSON 和 JSONC（带注释的 JSON）
- **Schema：** `https://opencode.ai/config.json`

### 配置文件位置与优先级（后者覆盖前者）

| 优先级 | 来源 | 路径 | 用途 |
|--------|------|------|------|
| 1 (最低) | 远程配置 | `.well-known/opencode` 端点 | 组织默认值 |
| 2 | 全局配置 | `~/.config/opencode/opencode.json` | 用户偏好 |
| 3 | 自定义配置 | `OPENCODE_CONFIG` 环境变量 | 自定义覆盖 |
| 4 | 项目配置 | 项目根目录 `opencode.json` | 项目设置 |
| 5 | `.opencode` 目录 | `.opencode/` | 代理、命令、插件 |
| 6 (最高) | 内联配置 | `OPENCODE_CONFIG_CONTENT` 环境变量 | 运行时覆盖 |

> 配置文件是**合并**的，非冲突设置会保留。

### 完整 JSON Schema

```jsonc
{
  "$schema": "https://opencode.ai/config.json",

  // ═══ 模型配置 ═══
  "model": "anthropic/claude-sonnet-4-5",
  "small_model": "anthropic/claude-haiku-4-5",

  // ═══ Provider 配置 ═══
  "provider": {
    "<provider-id>": {
      "npm": "@ai-sdk/openai-compatible",  // AI SDK 包名
      "name": "Provider Display Name",
      "models": {},                        // Provider 特定模型
      "options": {
        "timeout": 600000,                 // 请求超时(ms)，默认 300000
        "apiKey": "{env:API_KEY}",         // API 密钥
        "baseURL": "https://...",          // 自定义 Base URL
        "headers": {}                      // 自定义请求头
      }
    }
  },

  // ═══ Provider 白名单/黑名单 ═══
  "disabled_providers": ["openai", "gemini"],
  "enabled_providers": ["anthropic", "openai"],

  // ═══ TUI 配置 ═══
  "tui": {
    "scroll_speed": 3,
    "scroll_acceleration": { "enabled": true },
    "diff_style": "auto"
  },

  // ═══ 服务器配置 ═══
  "server": {
    "port": 4096,
    "hostname": "0.0.0.0",
    "mdns": true,
    "mdnsDomain": "myproject.local",
    "cors": ["http://localhost:5173"]
  },

  // ═══ 工具配置 ═══
  "tools": {
    "write": false,
    "bash": false,
    "edit": false
  },

  // ═══ 代理配置 ═══
  "agent": {
    "<agent-name>": {
      "description": "代理描述",
      "model": "anthropic/claude-sonnet-4-5",
      "prompt": "系统提示词",
      "tools": {
        "write": false,
        "edit": false
      }
    }
  },
  "default_agent": "plan",

  // ═══ 自定义命令 ═══
  "command": {
    "<command-name>": {
      "template": "提示模板，可用 $ARGUMENTS",
      "description": "命令描述",
      "agent": "build",
      "model": "anthropic/claude-haiku-4-5"
    }
  },

  // ═══ 快捷键 ═══
  "keybinds": {},

  // ═══ 分享 ═══
  "share": "manual",  // "manual" | "auto" | "disabled"

  // ═══ 自动更新 ═══
  "autoupdate": true,  // true | false | "notify"

  // ═══ 格式化程序 ═══
  "formatter": {
    "<formatter-name>": {
      "disabled": true,
      // 或自定义格式化器:
      "command": ["npx", "prettier", "--write", "$FILE"],
      "environment": { "NODE_ENV": "development" },
      "extensions": [".js", ".ts", ".jsx", ".tsx"]
    }
  },

  // ═══ 权限 ═══
  "permission": {
    "edit": "ask",   // "allow" | "ask" | "deny"
    "bash": "ask"
  },

  // ═══ 上下文压缩 ═══
  "compaction": {
    "auto": true,
    "prune": true,
    "reserved": 10000
  },

  // ═══ 文件监视器 ═══
  "watcher": {
    "ignore": ["node_modules/**", "dist/**", ".git/**"]
  },

  // ═══ MCP 服务器 ═══
  "mcp": {
    "<server-name>": {
      "type": "remote",
      "url": "https://...",
      "enabled": false
    }
  },

  // ═══ 插件 ═══
  "plugin": ["opencode-helicone-session"],

  // ═══ 指令 ═══
  "instructions": [
    "CONTRIBUTING.md",
    "docs/guidelines.md",
    ".cursor/rules/*.md"
  ],

  // ═══ 实验性功能 ═══
  "experimental": {}
}
```

### 变量替换

```json
// 环境变量替换
{ "model": "{env:OPENCODE_MODEL}" }

// 文件内容替换
{ "apiKey": "{file:~/.secrets/openai-key}" }
```
路径支持相对路径（相对于配置文件目录）、绝对路径、`~` 开头的路径。

### 目录结构配置

| 目录路径 | 用途 |
|----------|------|
| `.opencode/agents/` 或 `~/.config/opencode/agents/` | 代理定义（Markdown 文件） |
| `.opencode/commands/` 或 `~/.config/opencode/commands/` | 命令定义（Markdown 文件） |
| `.opencode/modes/` 或 `~/.config/opencode/modes/` | 模式定义 |
| `.opencode/plugins/` 或 `~/.config/opencode/plugins/` | 插件文件 |
| `.opencode/skills/` 或 `~/.config/opencode/skills/` | 代理技能 |
| `.opencode/tools/` 或 `~/.config/opencode/tools/` | 自定义工具 |
| `.opencode/themes/` 或 `~/.config/opencode/themes/` | 自定义主题 |

### 权限配置

```jsonc
{
  "permission": {
    "edit": "ask",    // 编辑文件前需要确认
    "bash": "ask",    // 执行 bash 命令前需要确认
    "write": "allow", // 允许无需确认
    "read": "deny"    // 完全拒绝
  }
}
```

| 值 | 行为 |
|----|------|
| `"allow"` | 始终允许（默认行为） |
| `"ask"` | 每次操作前询问 |
| `"deny"` | 始终拒绝 |

### 自定义 Provider 配置示例（UT-Bench 使用）

```json
{
  "$schema": "https://opencode.ai/config.json",
  "provider": {
    "utbench": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "UT-Bench Provider",
      "options": {
        "baseURL": "{env:MODEL_ENDPOINT}",
        "apiKey": "{env:DEEPSEEK_API_KEY}"
      },
      "models": {
        "deepseek-v4-flash": {
          "name": "deepseek-v4-flash"
        }
      }
    }
  }
}
```

---

## 3. Provider 配置

### 内置 Provider 列表（41 个）

| 序号 | Provider 名称 | 类型 |
|------|--------------|------|
| 1 | 302.AI | 云端 API |
| 2 | Amazon Bedrock | 云端（AWS） |
| 3 | Anthropic | 云端 API |
| 4 | Atomic Chat | 本地模型 |
| 5 | Azure OpenAI | 云端（Azure） |
| 6 | Azure Cognitive Services | 云端（Azure） |
| 7 | Baseten | 云端 API |
| 8 | Cerebras | 云端 API |
| 9 | Cloudflare AI Gateway | AI 网关/代理 |
| 10 | Cortecs | 云端 API |
| 11 | DeepSeek | 云端 API |
| 12 | Deep Infra | 云端 API |
| 13 | Firmware | 云端 API |
| 14 | Fireworks AI | 云端 API |
| 15 | GitLab Duo | 云端 API |
| 16 | GitHub Copilot | 云端（OAuth） |
| 17 | Google Vertex AI | 云端（GCP） |
| 18 | Groq | 云端 API |
| 19 | Hugging Face | 云端 API |
| 20 | Helicone | AI 网关/可观测性 |
| 21 | llama.cpp | 本地模型 |
| 22 | IO.NET | 云端 API |
| 23 | LM Studio | 本地模型 |
| 24 | Moonshot AI | 云端 API |
| 25 | MiniMax | 云端 API |
| 26 | Nebius Token Factory | 云端 API |
| 27 | Ollama | 本地模型 |
| 28 | Ollama Cloud | 云端 API |
| 29 | OpenAI | 云端 API |
| 30 | OpenCode Zen | 云端（官方推荐） |
| 31 | OpenRouter | AI 网关/代理 |
| 32 | SAP AI Core | 云端（企业级） |
| 33 | STACKIT | 云端（欧洲主权托管） |
| 34 | OVHcloud AI Endpoints | 云端 API |
| 35 | Scaleway | 云端 API |
| 36 | Together AI | 云端 API |
| 37 | Venice AI | 云端 API |
| 38 | Vercel AI Gateway | AI 网关/代理 |
| 39 | xAI | 云端 API |
| 40 | Z.AI | 云端 API |
| 41 | ZenMux | AI 网关/代理 |

### 自定义 Provider 配置方式

**通过配置文件（opencode.json）：**

```json
{
  "$schema": "https://opencode.ai/config.json",
  "provider": {
    "myprovider": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "My AI Provider",
      "options": {
        "baseURL": "https://api.myprovider.com/v1",
        "apiKey": "{env:MY_API_KEY}",
        "headers": {
          "Authorization": "Bearer custom-token"
        }
      },
      "models": {
        "my-model-name": {
          "name": "My Model",
          "limit": {
            "context": 200000,
            "output": 65536
          }
        }
      }
    }
  }
}
```

**完整配置字段：**

| 字段 | 必填 | 说明 |
|------|------|------|
| `npm` | ✅ | AI SDK 包名 |
| `name` | ✅ | UI 显示名称 |
| `options.baseURL` | ✅ | API 端点 URL |
| `options.apiKey` | ❌ | API 密钥，支持 `{env:ENV_VAR}` |
| `options.headers` | ❌ | 自定义请求头 |
| `options.timeout` | ❌ | 请求超时(ms)，默认 300000 |
| `models` | ✅ | 可用模型映射 |
| `models.<id>.name` | ❌ | 模型显示名称 |
| `models.<id>.limit.context` | ❌ | 最大输入 Token |
| `models.<id>.limit.output` | ❌ | 最大输出 Token |
| `models.<id>.options.provider` | ❌ | 路由提供商（网关类） |

**npm 包选择规则：**
- OpenAI 兼容（`/v1/chat/completions`）→ `@ai-sdk/openai-compatible`
- 走 `/v1/responses` → `@ai-sdk/openai`
- 有专用包的用专用包（如 `@ai-sdk/cerebras`）

### Amazon Bedrock 特殊配置

```json
{
  "provider": {
    "amazon-bedrock": {
      "options": {
        "region": "us-east-1",
        "profile": "my-aws-profile",
        "endpoint": "https://bedrock-runtime.us-east-1.vpce-xxxxx.amazonaws.com"
      }
    }
  }
}
```

> `endpoint` 是 `baseURL` 的 AWS 专用别名，如果同时指定，`endpoint` 优先。

### 本地模型配置

**Ollama（端口 11434）：**
```json
{
  "provider": {
    "ollama": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "Ollama (local)",
      "options": { "baseURL": "http://localhost:11434/v1" },
      "models": { "llama2": { "name": "Llama 2" } }
    }
  }
}
```

**llama.cpp（端口 8080）：**
```json
{
  "provider": {
    "llama.cpp": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "llama-server (local)",
      "options": { "baseURL": "http://127.0.0.1:8080/v1" },
      "models": {
        "qwen3-coder:a3b": {
          "name": "Qwen3-Coder: a3b-30b (local)",
          "limit": { "context": 128000, "output": 65536 }
        }
      }
    }
  }
}
```

**LM Studio（端口 1234）：**
```json
{
  "provider": {
    "lmstudio": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "LM Studio (local)",
      "options": { "baseURL": "http://127.0.0.1:1234/v1" },
      "models": {
        "google/gemma-3n-e4b": { "name": "Gemma 3n-e4b (local)" }
      }
    }
  }
}
```

### 网关类 Provider 路由

**OpenRouter（指定提供商路由）：**
```json
{
  "provider": {
    "openrouter": {
      "models": {
        "moonshotai/kimi-k2": {
          "options": {
            "provider": {
              "order": ["baseten"],
              "allow_fallbacks": false
            }
          }
        }
      }
    }
  }
}
```

**Vercel AI Gateway：**
```json
{
  "provider": {
    "vercel": {
      "models": {
        "anthropic/claude-sonnet-4": {
          "options": {
            "order": ["anthropic", "vertex"]
          }
        }
      }
    }
  }
}
```

### 特殊认证方式

- **Amazon Bedrock:** Bearer Token > AWS 凭证链
- **Anthropic:** 支持 Claude Pro/Max 订阅 OAuth
- **OpenAI:** 支持 ChatGPT Plus/Pro 订阅 OAuth
- **GitHub Copilot:** 设备码认证（`github.com/login/device`）
- **GitLab Duo:** OAuth 或 Personal Access Token

---

## 4. 环境变量

### 核心配置

| 变量 | 描述 |
|------|------|
| `OPENCODE_AUTO_SHARE` | 自动分享会话 |
| `OPENCODE_CONFIG` | 配置文件路径 |
| `OPENCODE_TUI_CONFIG` | TUI 配置文件路径 |
| `OPENCODE_CONFIG_DIR` | 配置目录路径 |
| `OPENCODE_CONFIG_CONTENT` | 内联 JSON 配置（优先级最高） |
| `OPENCODE_PERMISSION` | 内联 JSON 权限配置 |
| `OPENCODE_CLIENT` | 客户端标识（默认 `cli`） |
| `OPENCODE_GIT_BASH_PATH` | Windows 上 Git Bash 路径 |

### 服务器/认证

| 变量 | 描述 |
|------|------|
| `OPENCODE_SERVER_PASSWORD` | 启用 `serve`/`web` 的 basic auth |
| `OPENCODE_SERVER_USERNAME` | 覆盖 basic auth 用户名（默认 `opencode`） |

### 功能开关

| 变量 | 描述 |
|------|------|
| `OPENCODE_DISABLE_AUTOUPDATE` | 禁用自动更新检查 |
| `OPENCODE_DISABLE_PRUNE` | 禁用旧数据清理 |
| `OPENCODE_DISABLE_TERMINAL_TITLE` | 禁用终端标题更新 |
| `OPENCODE_DISABLE_DEFAULT_PLUGINS` | 禁用默认插件 |
| `OPENCODE_DISABLE_LSP_DOWNLOAD` | 禁用 LSP 自动下载 |
| `OPENCODE_DISABLE_AUTOCOMPACT` | 禁用自动上下文压缩 |
| `OPENCODE_DISABLE_MODELS_FETCH` | 禁用远程模型获取 |

### Claude Code 兼容性

| 变量 | 描述 |
|------|------|
| `OPENCODE_DISABLE_CLAUDE_CODE` | 禁止读取 `.claude` 配置 |
| `OPENCODE_DISABLE_CLAUDE_CODE_PROMPT` | 禁止读取 `~/.claude/CLAUDE.md` |
| `OPENCODE_DISABLE_CLAUDE_CODE_SKILLS` | 禁止加载 `.claude/skills` |

### 模型与搜索

| 变量 | 描述 |
|------|------|
| `OPENCODE_ENABLE_EXPERIMENTAL_MODELS` | 启用实验模型 |
| `OPENCODE_ENABLE_EXA` | 启用 Exa 网络搜索工具 |
| `OPENCODE_MODELS_URL` | 自定义模型配置 URL |

### 实验性功能

| 变量 | 描述 |
|------|------|
| `OPENCODE_EXPERIMENTAL` | 启用所有实验功能 |
| `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | Bash 命令默认超时(ms) |
| `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX` | LLM 响应最大输出 token |
| `OPENCODE_EXPERIMENTAL_PLAN_MODE` | 启用计划模式 |
| `OPENCODE_EXPERIMENTAL_MARKDOWN` | 启用实验 Markdown 功能 |
| `OPENCODE_EXPERIMENTAL_FILEWATCHER` | 启用完整目录文件监视 |
| `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER` | 禁用文件监视 |

### 测试

| 变量 | 描述 |
|------|------|
| `OPENCODE_FAKE_VCS` | 测试用的假 VCS provider |

### 故障排除

1. 运行 `opencode auth list` 查看凭据
2. 确认 provider ID 在 `/connect` 和 `opencode.json` 中一致
3. 选择正确的 npm 包：
   - 有专用包用专用包
   - OpenAI 兼容用 `@ai-sdk/openai-compatible`
   - `/v1/responses` 用 `@ai-sdk/openai`
4. 验证 `options.baseURL` 地址正确
