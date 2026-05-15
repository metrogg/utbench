# **概述**

**腾讯云智能编程助手 - 让代码开发更智能、更高效**

CodeBuddy Code 是基于腾讯云 AI 技术的智能编程工具，深度集成腾讯云生态，提供从代码编写到项目部署的全链路 AI 辅助。

## **为什么选择 CodeBuddy Code？**

### **🚀 用自然语言驱动整个开发运维生命周期**

CodeBuddy Code 让您能够用自然语言描述需求，自动化完成从代码编写、测试、调试到部署的全链路开发任务，实现极致的自动化效率提升。无论是简单的代码修改还是复杂的架构重构，都能通过对话式交互轻松完成。

### **🔧 终端原生，无缝集成**

- **熟悉的环境**：直接在您熟悉的命令行环境中获得 AI 辅助，无需切换开发工具或学习新界面
- **原生体验**：完美融入现有的开发工作流，支持所有主流操作系统和终端
- **零学习成本**：保持原有的开发习惯，AI 助手静默工作在后台

### **⚡ 开箱即用的强大能力**

- **内置工具链**：集成文件编辑、命令运行、Git 操作、测试执行等核心开发工具
- **智能提交**：自动生成规范的提交信息，支持代码审查和变更管理
- **灵活扩展**：通过 MCP (模型上下文协议) 轻松集成第三方工具和服务
- **自定义开发工具**：根据项目需求定制专属的开发助手

### **🛠️ Unix 哲学的 AI 集成**

- **管道友好**：像 `grep` 和 `awk` 一样，原生支持管道输入进行智能分析
- **脚本集成**：完美融入 shell 脚本和自动化工具链
- **组合能力**：与现有 Unix 工具无缝组合，构建强大的 AI 驱动工作流
- **标准输入输出**：遵循 Unix 标准，支持重定向和管道操作
  ```
  # 管道集成示例

  git log --oneline | codebuddy "分析这些提交，找出可能的问题"
  cat error.log | codebuddy "帮我分析这些错误日志"
  ```

## **快速体验**

### **环境要求**

- Node.js 18.0+

### **一键安装**

```
npm install -g @tencent-ai/codebuddy-code
```

### **开始使用**

```
# 进入项目目录
cd my-project

# 启动 CodeBuddy
codebuddy
cbc

# 或直接提问
codebuddy "帮我优化这个函数的性能"
cbc "帮我优化这个函数的性能"
```

## **下一步操作**

### **📚 深入了解**

- **[快速入门指南](https://www.codebuddy.ai/docs/zh/cli/quickstart)** - 详细的设置和使用指南
- **[常见工作流](https://www.codebuddy.ai/docs/zh/cli/common-workflows)** - 扩展思考、图片粘贴、--resume 等功能
- **[IDE 集成](https://www.codebuddy.ai/docs/zh/cli/ide-integrations)** - 在你喜欢的编辑器中使用

### **🔧 配置和扩展**

- **[MCP 集成](https://www.codebuddy.ai/docs/zh/cli/mcp)** - 模型上下文协议集成
- **[斜杠命令](https://www.codebuddy.ai/docs/zh/cli/slash-commands)** - 内置命令参考
- **[设置配置](https://www.codebuddy.ai/docs/zh/cli/settings)** - 配置文件、环境变量、工具设置

### **🚀 高级用法**

- **GitHub Actions** - CI/CD 集成
- **SDK 开发** - 扩展和自定义
- **企业部署** - 容器化部署

### **🆘 获取帮助**

- **[故障排除](https://www.codebuddy.ai/docs/zh/cli/troubleshooting)** - 常见问题解决
- **[CLI 参考](https://www.codebuddy.ai/docs/zh/cli/reference)** - 完整命令行参考
- **交互模式** - 键盘快捷键和技巧

## **反馈和支持**

- 📧 技术支持：codebuddy\@tencent.com
- 🌐 **[国内官方网站](https://copilot.tencent.com/cli)**
- 🌐 **[海外官方网站](https://www.codebuddy.ai/cli)**

**最后更新: 2025/12/10 14:19**

<br />

# **快速入门指南**

欢迎使用 CodeBuddy Code！这份指南将帮助您在 5 分钟内上手，体验自然语言驱动的编程助手。

## **🎯 开始之前**

### **系统要求**

- **Node.js**：版本 18.20 或更高
- **操作系统**: macOS、Linux 或 Windows

### **验证环境**

```
# 检查 Node.js 版本
node --version  # 应显示 v18.0.0 或更高

# 检查 npm 版本
npm --version
```

## **⚡ 极速安装**

### **npm 全局安装**

```
npm install -g @tencent-ai/codebuddy-code
```

### **原生安装器（Beta）**

> ⚠️ **Beta 功能**：原生安装器目前处于 Beta 阶段。我们推荐您尝试使用，以获得更快速、独立的安装体验。

原生安装器提供独立的 CodeBuddy 安装，无需 Node.js 环境。

**下载并安装：**

```
# macOS/Linux
curl -fsSL https://copilot.tencent.com/cli/install.sh | bash
```

```
# Windows
irm https://copilot.tencent.com/cli/install.ps1 | iex
```

**优势：**

- ✅ 无需 Node.js 依赖
- ✅ 安装和启动速度更快

### **验证安装**

```
codebuddy --version
```

## **🔐 登录认证**

首次使用 CodeBuddy Code 时，您需要完成登录认证。启动后会显示登录方式选择界面：

```
Select login method:
› Log in via Chinese Site
  Log in via International Site
  Log in via Enterprise Domain
  Log in via iOA (Tencent only)
```

### **登录方式说明**

**登录方式**

**适用场景**

**说明**

**Chinese Site**

国内用户

通过腾讯云国内站 (copilot.tencent.com) 进行认证，支持国内主流模型

**International Site**

海外用户

通过腾讯云国际站 (codebuddy.ai) 进行认证，支持海外主流模型

**Enterprise Domain**

专享版/私有化部署

连接企业专享版或自建的 CodeBuddy 服务，需要输入企业提供的服务地址

**iOA**

腾讯内部员工

通过腾讯 iOA 零信任系统进行认证，仅限腾讯内部员工使用

使用 `↑↓` 键选择登录方式，按 `Enter` 确认后会自动打开浏览器完成认证。

## **🚀 首次体验**

### **1. 进入项目目录**

```
cd /path/to/your/project
```

### **2. 启动交互模式**

```
codebuddy
```

您将看到欢迎界面：

> **提示**：如果您希望 CodeBuddy Code 始终使用特定语言回复（如简体中文），可以在启动后运行 `/config` 命令设置 Language 选项。

```
🤖 CodeBuddy Code v1.0.0
💡 输入 /help 查看可用命令
📝 开始对话，让 AI 成为您的编程伙伴

>
```

### **3. 初始化项目上下文（强烈推荐）**

在正式开始对话之前，**强烈建议**先使用 `/init` 命令初始化项目上下文：

```
> /init
```

**为什么 /init 如此重要？**

**📊 效果提升：**

- ✅ **理解更准确**：预先构建项目知识图谱，AI 能更准确理解代码结构和业务逻辑
- ✅ **响应更快速**：避免重复扫描文件，后续对话响应速度显著提升
- ✅ **建议更精准**：基于全局视图提供更符合项目架构的建议
- ✅ **减少误判**：了解项目依赖关系，避免提出不兼容的修改方案

**💰 成本优化：**

- ✅ **Token 消耗更少**：一次性建立上下文，避免每次对话都重新分析
- ✅ **减少重复请求**：预加载关键信息，减少 30-50% 的上下文 Token 开销
- ✅ **更高效的对话**：每轮对话携带更少的冗余信息，整体成本更低

**最佳实践：**

```
# 第一次使用项目时
> /init

# 项目结构发生重大变化时（如添加新模块、重构等）
> /clear  # 开启全新对话
> /init   # 重新初始化
```

### **4. 尝试第一个对话**

```
> 帮我分析这个项目的结构
```

CodeBuddy 会自动扫描您的项目文件，并提供详细的结构分析。

## **💡 核心使用模式**

### **交互式对话模式**

最自然的使用方式，适合探索性开发：

```
codebuddy
```

**典型对话示例：**

```
> 我想给这个 React 组件添加一个加载状态
> 帮我重构这个函数，让它更易读
> 这段代码有什么潜在的性能问题？
> 为这个 API 接口写单元测试
```

### **单次命令模式**

适合脚本化和自动化场景：

```
# 直接提问
codebuddy -p "优化这个 SQL 查询的性能"

# 管道输入
cat error.log | codebuddy -p "分析这些错误日志"

# 文件分析（需要授权时必须添加 -y 或 --dangerously-skip-permissions）
codebuddy -p "审查 src/utils.js 的代码质量" -y
```

> **重要提示**：使用 `-p/--print` 参数进行单次执行时,如果操作需要访问文件、执行命令等授权操作,必须添加 `-y` （或 `--dangerously-skip-permissions`) 参数。

### **项目级操作**

处理复杂的跨文件任务：

```
# 项目重构（需要文件操作授权）
codebuddy -p "将所有组件从 class 组件迁移到函数组件" -y

# 代码规范（需要文件读取授权）
codebuddy -p "检查整个项目的 TypeScript 类型定义" -y

# 测试覆盖（需要文件操作授权）
codebuddy -p "为 services 目录下的所有文件添加单元测试" -y
```

### **快捷键**

#### **基础导航**

**快捷键**

**功能**

`↑/↓`

浏览命令历史

`↓`

查看后台任务（当有运行中任务时）

`Tab`

命令自动补全

`Esc`

清空输入（按两次） / 返回上级菜单

`Ctrl+C`

退出程序

`Ctrl+D`

退出程序（输入框为空且无进行中对话时，连按两次）

#### **权限和模式**

**快捷键**

**功能**

`Shift+Tab` (macOS/Linux) 或 `Alt+M` (Windows)

切换权限模式（default → bypass → accept → plan）

#### **编辑功能**

**快捷键**

**功能**

`Ctrl+R`

展开/收起详细输出（在对话中）

`Ctrl+G`

使用外部编辑器编辑提示词

`Enter`

发送消息

`Shift+Enter`

换行（多行输入）

`\Enter`

换行（反斜杠转义）

`Ctrl+J`

插入换行（多行输入）

#### **面板操作**

**快捷键**

**功能**

`↑/↓`

在面板中导航选项

`Enter`

选择当前项

`Space`

切换选择（多选面板）

`k`

终止选中的后台任务

#### **专用功能**

**快捷键**

**功能**

`j/k`

Vim风格的上下导航（部分面板支持）

`Ctrl+O`

查看思考详情面板

当终端显示“Thinking”指示时，可用 `Ctrl+O` 打开完整推理内容,再按 `Ctrl+O` 退出。

## **🎓 进阶学习**

恭喜您完成快速入门！接下来推荐阅读：

- **[常见工作流](https://www.codebuddy.ai/docs/zh/cli/common-workflows)** - 学习扩展思考、图片分析等高级功能
- **[IDE 集成](https://www.codebuddy.ai/docs/zh/cli/ide-integrations)** - 在 VS Code、JetBrains 中使用
- **[斜杠命令](https://www.codebuddy.ai/docs/zh/cli/slash-commands)** - 掌握所有内置命令
- **[MCP 集成](https://www.codebuddy.ai/docs/zh/cli/mcp)** - 扩展自定义工具

## **💬 获取帮助**

遇到问题？我们随时为您提供支持：

- 🐛 **[提交 Bug](https://cnb.cool/codebuddy/codebuddy-code/-/issues)**
- 📧 技术支持：codebuddy\@tencent.com
- 📚 **[完整文档](https://www.codebuddy.ai/docs/zh/cli/README)**
- 🌐 **[官方网站](https://copilot.tencent.com/cli)**

***

*现在开始，让 AI 成为您的编程伙伴！🚀*

**最后更新: 2026/4/7 15:00**

<br />

# **CodeBuddy Code 安装指南**

## **安装方式**

### **📦 使用包管理器安装**

#### **Node.js 包管理器**

**前置要求:** Node.js 18.20 或更高版本

选择你喜欢的包管理器执行以下命令:

**npmpnpmyarnbun**

```
npm install -g @tencent-ai/codebuddy-code
```

#### **Homebrew (macOS/Linux)**

无需 Node.js，直接安装:

**两步安装单命令安装Brewfile**

```
# 添加 tap
brew tap Tencent-CodeBuddy/tap

# 安装工具
brew install codebuddy-code
```

#### **验证安装**

安装完成后,运行以下命令验证是否安装成功:

```
codebuddy --version
```

### **⚙️ 使用原生二进制安装（Beta）**

> ⚠️ **Beta 测试阶段**
>
> 原生二进制安装目前处于 Beta 测试阶段，功能仍在完善中。 如遇到任何问题，请在 **[Issues 页面](https://cnb.cool/codebuddy/codebuddy-code/-/issues)** 提交问题报告，或联系技术支持（codebuddy\@tencent.com）。

#### **特性说明**

原生二进制安装相比 npm 版本提供以下特性：

- 单一可执行文件，无需额外依赖
- 无需 Node.js 运行时
- 改进的自动更新机制

#### **支持平台**

- macOS (Apple Silicon M1/M2/M3 或 Intel x86\_64)
- Linux (arm64 或 x86\_64)
- Windows (x86\_64)

#### **从 npm 版本迁移**

如果你已经通过 npm 安装了 CodeBuddy Code，可以使用以下命令迁移到原生二进制版本：

```
codebuddy install
```

#### **全新安装**

**macOS / LinuxWindows**

```
curl -fsSL https://www.codebuddy.cn/cli/install.sh | bash
```

#### **验证安装**

安装脚本会自动下载最新版本并配置环境变量。安装完成后，运行以下命令验证：

```
codebuddy --version
```

**如果命令不可用**,请手动将安装路径添加到环境变量 `PATH`:

**macOS / LinuxWindows**

```
export PATH="$HOME/.local/bin:$PATH"

# 为了永久生效,建议添加到 shell 配置文件:
# Bash: ~/.bashrc 或 ~/.bash_profile
# Zsh: ~/.zshrc
```

## **📁 配置目录**

CodeBuddy Code 默认将配置文件存储在以下目录:

**平台**

**默认配置目录**

macOS / Linux

`~/.codebuddy`

Windows

`%USERPROFILE%\.codebuddy`

### **配置目录内容**

```
~/.codebuddy/
├── settings.json      # 用户设置
├── mcp.json           # MCP 服务器配置
├── .mcp.json          # MCP 服务器配置（备选位置）
└── skills/            # 用户自定义 Skills
```

### **自定义配置目录**

通过设置环境变量 `CODEBUDDY_CONFIG_DIR` 可以自定义配置目录位置:

```
export CODEBUDDY_CONFIG_DIR="$HOME/.my-codebuddy-config"
```

这在以下场景中非常有用:

- 多个 CodeBuddy 实例需要独立配置
- 企业环境中需要统一管理配置位置
- 与其他使用 CodeBuddy 引擎的应用（如 WorkBuddy）共存时避免配置冲突

## **🔄 更新**

### **自动更新**

CodeBuddy Code 默认会自动保持最新状态，以确保你拥有最新的功能和安全修复。

#### **关闭自动更新**

如需关闭自动更新，可设置环境变量：

```
export DISABLE_AUTOUPDATER=1
```

### **手动更新**

使用以下命令手动更新到最新版本：

```
codebuddy update
```

`update` 命令会自动检测你的安装方式并执行相应的更新操作。

#### **使用包管理器更新**

如果 `codebuddy update` 命令未能成功更新，你也可以使用包管理器重新安装：

```
npm install -g @tencent-ai/codebuddy-code
```

或使用其他包管理器（pnpm、yarn、bun）执行相应的安装命令。

## **🔧 故障排查**

### **命令不可用**

**问题：** 安装后提示 `codebuddy: command not found`

**解决方案：**

1. 检查安装路径是否在 `PATH` 环境变量中：
   ```
   echo $PATH
   ```
2. 将 CodeBuddy 安装路径添加到 `PATH`（参考上方\[验证安装]\(#验证安装）部分）
3. 重启终端或重新加载配置文件：
   ```
   source ~/.bashrc  # 或 ~/.zshrc
   ```

### **网络问题**

**问题：** 安装或更新时网络连接失败

**解决方案：**

1. 检查网络连接
2. 配置 npm 镜像源（如果使用 npm 安装）：
   ```
   npm config set registry https://registry.npmmirror.com
   ```

## **🗑️ 卸载**

### **包管理器版本卸载**

**Homebrewnpmpnpmyarnbun**

```
# 卸载工具
brew uninstall codebuddy-code

# 移除 tap (可选)
brew untap Tencent-CodeBuddy/tap
```

### **原生二进制版本卸载**

#### **macOS / Linux**

删除可执行文件:

```
rm -f ~/.local/bin/codebuddy
```

### **清理配置文件(可选)**

如需完全清理,可删除配置目录:

**macOS / Linux:**

```
rm -rf ~/.codebuddy
rm -rf ~/.local/share/codebuddy
```

> 💡 **提示:** 如果你使用了 `CODEBUDDY_CONFIG_DIR` 环境变量自定义了配置目录，请删除对应的目录。

# **无头模式 （Headless Mode)**

> 以编程方式运行 CodeBuddy Code,无需交互式 UI

## **概述**

无头模式允许您通过命令行脚本和自动化工具以编程方式运行 CodeBuddy Code,无需任何交互式 UI。

无头模式也支持定时任务相关能力。在脚本、SDK 或服务端集成场景中，可以使用 `CronCreate`、`CronList`、`CronDelete` 等工具来创建、查看和取消定时任务。

> **⚠️ 重要提示：** `-y` （或 `--dangerously-skip-permissions`) 是非交互模式的必需参数。在使用 `-p/--print` 参数进行非交互式执行时，必须添加此参数才能执行需要授权的操作（文件读写、命令执行、网络请求等）,否则这些操作会被阻止。仅在受信任的环境和明确的任务场景下使用此参数。详见 **[CLI 参考](https://www.codebuddy.ai/docs/zh/cli/cli-reference)**。

## **基本用法**

CodeBuddy Code 的主要命令行接口是 `codebuddy` （或 `cbc`) 命令。使用 `--print` （或 `-p`) 标志在非交互模式下运行并打印最终结果：

```
codebuddy -p "暂存我的更改并为它们编写一组提交" \
  --allowedTools "Bash,Read" \
  --permission-mode acceptEdits
```

## **配置选项**

无头模式利用 CodeBuddy Code 中所有可用的 CLI 选项。以下是用于自动化和脚本编写的关键选项：

**标志**

**描述**

**示例**

`--print`, `-p`

在非交互模式下运行

`codebuddy -p "查询"`

`--output-format`

指定输出格式 （`text`, `json`, `stream-json`)

`codebuddy -p --output-format json`

`--resume`, `-r`

通过会话 ID 恢复对话

`codebuddy --resume abc123`

`--continue`, `-c`

继续最近的对话

`codebuddy --continue`

`--verbose`

启用详细日志记录

`codebuddy --verbose`

`--append-system-prompt`

追加到系统提示词 （仅与 `--print` 配合使用）

`codebuddy --append-system-prompt "自定义指令"`

`--allowedTools`

允许的工具列表,空格分隔或\
\
逗号分隔的字符串

`codebuddy --allowedTools mcp__slack mcp__filesystem`\
\
`codebuddy --allowedTools "Bash(npm install),mcp__filesystem"`

`--disallowedTools`

拒绝的工具列表,空格分隔或\
\
逗号分隔的字符串

`codebuddy --disallowedTools mcp__splunk mcp__github`\
\
`codebuddy --disallowedTools "Bash(git commit),mcp__github"`

`--settings`

从 JSON 文件或 JSON 字符串加载额外的设置配置

`codebuddy -p --settings '{"model":"gpt-5"}' "查询"`

`--setting-sources`

指定要加载的设置源（可选值: `user`, `project`, `local`）

`codebuddy -p --setting-sources project,local "查询"`

`--mcp-config`

从 JSON 文件加载 MCP 服务器

`codebuddy --mcp-config servers.json`

`--permission-prompt-tool`

用于处理权限提示的 MCP 工具 （仅与 `--print` 配合使用）

❌ 不支持

> **说明：** `--permission-prompt-tool` 功能当前不支持。

有关 CLI 选项和功能的完整列表，请参阅 **[CLI 参考](https://www.codebuddy.ai/docs/zh/cli/cli-reference)** 文档。

## **多轮对话**

对于多轮对话，您可以恢复对话或从最近的会话继续：

```
# 继续最近的对话
codebuddy --continue "现在重构以提高性能"

# 通过会话 ID 恢复特定对话
codebuddy --resume 550e8400-e29b-41d4-a716-446655440000 "更新测试"

# 在非交互模式下恢复
codebuddy --resume 550e8400-e29b-41d4-a716-446655440000 "修复所有 linting 问题" -p
```

## **输出格式**

### **文本输出 （默认）**

```
codebuddy -p "解释文件 src/components/Header.tsx"
# 输出： 这是一个 React 组件，显示...
```

### **JSON 输出**

返回包含元数据的结构化数据：

```
codebuddy -p "数据层是如何工作的?" --output-format json
```

响应格式：

```
{
 ...
}
```

### **流式 JSON 输出**

在收到每条消息时流式传输：

```
codebuddy -p "构建一个应用程序" --output-format stream-json
```

每个对话都以初始 `init` 系统消息开始,然后是用户和助手消息列表,最后是包含统计信息的最终 `result` 系统消息。每条消息都作为单独的 JSON 对象发出。

### **结构化 JSON 输出**

要获得符合特定架构的输出，请使用 `--output-format json` 和 `--json-schema` 以及 **[JSON Schema](https://json-schema.org/)** 定义。响应包括关于请求的元数据（会话 ID、使用情况等），结构化输出在 `structured_output` 字段中。

此示例从 auth.py 中提取函数名称并将其作为字符串数组返回：

```
codebuddy -p "提取 auth.py 中的主要函数名称" \
  --output-format json \
  --json-schema '{"type":"object","properties":{"functions":{"type":"array","items":{"type":"string"}}},"required":["functions"]}'
```

> **提示**：使用 **[jq](https://jqlang.github.io/jq/)** 之类的工具来解析响应并提取特定字段：
>
> ```
> # 提取文本结果
> codebuddy -p "总结这个项目" --output-format json | jq -r '.result'
>
> # 提取结构化输出
> codebuddy -p "提取 auth.py 中的函数名称" \
>   --output-format json \
>   --json-schema '{"type":"object","properties":{"functions":{"type":"array","items":{"type":"string"}}},"required":["functions"]}' \
>   | jq '.structured_output'
> ```

## **输入格式**

### **文本输入 （默认）**

```
# 直接参数
codebuddy -p "解释这段代码"

# 从 stdin
echo "解释这段代码" | codebuddy -p
```

### **流式 JSON 输入**

通过 `stdin` 提供的消息流,其中每条消息代表一个用户轮次。这允许在不重新启动 `codebuddy` 二进制文件的情况下进行多轮对话，并允许在模型处理请求时向模型提供指导。

每条消息都是一个 JSON "用户消息" 对象,遵循与输出消息模式相同的格式。消息使用 **[jsonl](https://jsonlines.org/)** 格式进行格式化,其中每行输入都是一个完整的 JSON 对象。流式 JSON 输入需要 `-p` 和 `--output-format stream-json`。

```
echo '{"type":"user","message":{"role":"user","content":[{"type":"text","text":"解释这段代码"}]}}' | \
  codebuddy -p --output-format=stream-json --input-format=stream-json --verbose

# 单条消息（包含图片）
echo '{"type":"user","message":{"role":"user","content":[{"type":"text","text":"文本提示词，如参考以下图片的文案"},{"type":"image","source":{"type":"base64","media_type":"image/png","data":"原始base64（不含协议前缀）"}}]}}' \
  | codebuddy -p --input-format stream-json --output-format stream-json

# 多轮对话（多行 JSON，对同一进程持续发送）
printf '%s\n' \
  '{"type":"user","message":{"role":"user","content":[{"type":"text","text":"第一问"}]}}' \
  '{"type":"user","message":{"role":"user","content":[{"type":"text","text":"第二问"}]}}' \
  | codebuddy -p --input-format stream-json --output-format stream-json --verbose
```

## **Agent 集成示例**

### **SRE 事件响应机器人**

```
#!/bin/bash

# 自动化事件响应 agent
investigate_incident() {
    local incident_description="$1"
    local severity="${2:-medium}"

    codebuddy -p "事件: $incident_description (严重性: $severity)" \
      --append-system-prompt "你是一名 SRE 专家。诊断问题，评估影响，并提供即时行动项。" \
      --output-format json \
      --allowedTools "Bash,Read,WebSearch,mcp__datadog" \
      --mcp-config monitoring-tools.json
}

# 使用方式
investigate_incident "支付 API 返回 500 错误" "high"
```

### **自动化安全审查**

```
# PR 的安全审计 agent
audit_pr() {
    local pr_number="$1"

    gh pr diff "$pr_number" | codebuddy -p \
      --append-system-prompt "你是一名安全工程师。审查此 PR 的漏洞、不安全模式和合规问题。" \
      --output-format json \
      --allowedTools "Read,Grep,WebSearch"
}

# 使用并保存到文件
audit_pr 123 > security-report.json
```

### **多轮法律助手**

```
# 具有会话持久性的法律文档审查
session_id=$(codebuddy -p "开始法律审查会话" --output-format json | jq -r '.session_id')

# 分多个步骤审查合同
codebuddy -p --resume "$session_id" "审查 contract.pdf 的责任条款"
codebuddy -p --resume "$session_id" "检查 GDPR 要求的合规性"
codebuddy -p --resume "$session_id" "生成风险执行摘要"
```

## **最佳实践**

- **使用 JSON 输出格式** 进行程序化解析响应：
  ```
  # 使用 jq 解析 JSON 响应
  result=$(codebuddy -p "生成代码" --output-format json)
  code=$(echo "$result" | jq -r '.result')
  cost=$(echo "$result" | jq -r '.total_cost_usd')
  ```
- **优雅地处理错误** - 检查退出代码和 stderr:
  ```
  if ! codebuddy -p "$prompt" 2>error.log; then
      echo "发生错误:" >&2
      cat error.log >&2
      exit 1
  fi
  ```
- **使用会话管理** 在多轮对话中维护上下文
- **考虑超时** 对于长时间运行的操作：
  ```
  timeout 300 codebuddy -p "$complex_prompt" || echo "5 分钟后超时"
  ```
- **遵守速率限制** 在进行多个请求时，通过在调用之间添加延迟
- **使用 `-y`** 在非交互模式下执行需要授权的操作：
  ```
  # 非交互模式下的完整示例
  codebuddy -p "分析代码并运行测试" \
    --output-format json \
    -y \
    --allowedTools "Bash,Read,Grep"
  ```
  > **⚠️ 重要提示：** `-y` （或 `--dangerously-skip-permissions`) 是非交互模式的必需参数。在使用 `-p/--print` 参数进行非交互式执行时，必须添加此参数才能执行需要授权的操作（文件读写、命令执行、网络请求等）,否则这些操作会被阻止。仅在受信任的环境和明确的任务场景下使用此参数。详见 **[CLI 参考](https://www.codebuddy.ai/docs/zh/cli/cli-reference)**。

## **相关资源**

- **[CLI 参考](https://www.codebuddy.ai/docs/zh/cli/cli-reference)** - 完整的 CLI 文档
- **[常见工作流](https://www.codebuddy.ai/docs/zh/cli/common-workflows)** - 常见用例的分步指南
- **[交互模式](https://www.codebuddy.ai/docs/zh/cli/interactive-mode)** - 交互式会话功能
- **[IAM 权限](https://www.codebuddy.ai/docs/zh/cli/iam)** - 工具权限和访问控制

***

> **提示**：无头模式非常适合 CI/CD 管道、自动化脚本和 agent 集成。将其与 MCP 服务器结合使用以扩展功能。

**最后更新: 2026/3/11 18:03**

<br />

# **故障排查与最佳实践**

本文档涵盖常见问题解决方案和使用优化建议。

***

## **安装问题**

### **Node.js 版本要求**

CodeBuddy Code 需要 Node.js v18.20 或更高版本。

```
node -v          # 检查版本
```

升级地址：**<https://nodejs.org/en/download/>**

### **Windows 平台**

#### **Git Bash 依赖**

Windows 平台推荐安装 Git Bash。缺失时 CodeBuddy Code 会启动时打印一次性提示并自动降级到 PowerShell 执行 shell 命令，`/enter-worktree`、`/leave-worktree` 等 Git Bash 特有功能将不可用。

- 下载：**<https://git-scm.com/downloads/win>**
- 安装时勾选 "Git Bash" 组件

**自定义路径**（非标准安装位置）：

```
# CMD
set CODEBUDDY_CODE_GIT_BASH_PATH=C:\Program Files\Git\bin\bash.exe

# PowerShell
$env:CODEBUDDY_CODE_GIT_BASH_PATH="C:\Program Files\Git\bin\bash.exe"
```

**抑制启动提示**（例如上游进程已管理 shell）：

```
# PowerShell
$env:CODEBUDDY_SKIP_GIT_BASH_CHECK="1"
```

#### **"codebuddy 不是内部或外部命令"**

npm 全局目录未加入 PATH。

```
npm config get prefix    # 查找安装路径
```

将路径添加到系统 PATH（默认：`%USERPROFILE%\AppData\Roaming\npm`），重启终端。

### **搜索工具 (Ripgrep)**

CodeBuddy 自动处理 ripgrep 依赖，无需手动安装。如需最佳性能：

```
# macOS
brew install ripgrep

# Windows
choco install ripgrep

# Ubuntu/Debian
sudo apt install ripgrep
```

***

## **常见问题**

### **额度共享**

CLI、CodeBuddy IDE 和 CodeBuddy Plugin 共享同一账号的资源配额。

### **JetBrains IDE 中 ESC 键不生效**

JetBrains 终端对 ESC 键处理不同，改用 `Ctrl+ESC` 或 `Shift+ESC`。

**操作**

**标准终端**

**JetBrains 终端**

退出/取消

`ESC`

`Ctrl+ESC` 或 `Shift+ESC`

### **模型切换**

```
/model              # 交互式选择
/model [模型名称]    # 直接切换
/status             # 查看当前模型
```

### **--serve 模式网络访问问题**

#### **局域网 IP 访问报 `ERR_EMPTY_RESPONSE`**

**症状**：使用 `codebuddy --serve` 启动后，通过局域网 IP（如 `http://10.31.110.26:52477`）访问时，浏览器报 `ERR_EMPTY_RESPONSE`（未发送任何数据）。

**原因**：HTTP 服务默认监听 `127.0.0.1`（本地回环地址），只接受本机连接。使用局域网 IP 访问时，请求到达了机器但被服务拒绝。

**解决方案**：启动时添加 `--host 0.0.0.0` 参数，使服务监听所有网络接口：

```
codebuddy --serve --host 0.0.0.0 --port 8080
```

> 非回环地址启动时会自动启用密码认证，密码在控制台输出中可见。

#### **确认服务监听状态**

可以通过以下命令检查服务是否正确监听：

```
# macOS / Linux
lsof -i :PORT_NUMBER

# 或
netstat -an | grep PORT_NUMBER
```

- 若显示 `127.0.0.1:PORT` — 只接受本机连接
- 若显示 `0.0.0.0:PORT` 或 `*:PORT` — 接受所有网络连接

***

## **更新**

### **自动更新**

默认开启，下次启动时自动应用新版本。通过 `/config` 管理开关。

### **手动更新**

```
codebuddy update                                    # 内置命令（推荐）
npm install -g @tencent-ai/codebuddy-code@latest   # npm 更新
```

### **版本检查**

```
codebuddy --version                                 # 当前版本
npm view @tencent-ai/codebuddy-code version        # 最新版本
```

***

## **从 Claude Code 迁移**

### **迁移内容**

**目录/文件**

**说明**

`agents/`

自定义 agents 配置

`commands/`

斜杠命令定义

`skills/`

专业技能定义

`CLAUDE.md` → `CODEBUDDY.md`

AI 指令和记忆文档

### **方案一：符号链接（推荐）**

共享配置，修改一处两边生效。

```
# macOS/Linux
cd ~/.codebuddy
ln -s ~/.claude/agents agents
ln -s ~/.claude/commands commands
ln -s ~/.claude/skills skills
ln -s ~/.claude/CLAUDE.md CODEBUDDY.md
```

```
# Windows (需管理员权限)
cd $env:USERPROFILE\.codebuddy
New-Item -ItemType SymbolicLink -Path agents -Target $env:USERPROFILE\.claude\agents
New-Item -ItemType SymbolicLink -Path commands -Target $env:USERPROFILE\.claude\commands
New-Item -ItemType SymbolicLink -Path skills -Target $env:USERPROFILE\.claude\skills
New-Item -ItemType SymbolicLink -Path CODEBUDDY.md -Target $env:USERPROFILE\.claude\CLAUDE.md
```

### **方案二：复制文件**

独立配置，互不影响。

```
# macOS/Linux
cp -r ~/.claude/agents ~/.codebuddy/agents
cp -r ~/.claude/commands ~/.codebuddy/commands
cp -r ~/.claude/skills ~/.codebuddy/skills
cp ~/.claude/CLAUDE.md ~/.codebuddy/CODEBUDDY.md
```

```
# Windows
Copy-Item -Recurse $env:USERPROFILE\.claude\agents $env:USERPROFILE\.codebuddy\agents
Copy-Item -Recurse $env:USERPROFILE\.claude\commands $env:USERPROFILE\.codebuddy\commands
Copy-Item -Recurse $env:USERPROFILE\.claude\skills $env:USERPROFILE\.codebuddy\skills
Copy-Item $env:USERPROFILE\.claude\CLAUDE.md $env:USERPROFILE\.codebuddy\CODEBUDDY.md
```

### **插件 Skills 一键安装**

Claude Code 插件中的 Skills 支持一键安装，安装后自动加载。

### **验证迁移**

```
codebuddy         # 启动
/skills           # 检查 Skills
/config           # 查看配置
```

***

## **成本优化**

### **核心原则**

- 新任务用 `/clear` 开启新会话
- 长对话用 `/compact` 压缩历史
- 用 `@filename` 引用文件，避免粘贴代码

### **会话管理命令**

**命令**

**功能**

`/cost`

查看 Token 消耗

`/clear`

开启新会话

`/compact`

压缩历史

`/resume`

恢复旧对话

### **成本对比**

**方式**

**输入 Token**

**相对成本**

单会话连续 10 个任务

\~50,000

高

每个任务新会话

\~15,000

低

定期 `/compact`

\~25,000

中

### **推荐做法**

- ✓ 新任务开新会话
- ✓ 每 20-30 轮用 `/compact`
- ✓ 用 `@filename` 引用文件
- ✓ 精简提问

### **避免做法**

- ✗ 同一会话处理多个无关任务
- ✗ 对话超过 30 轮不清理
- ✗ 重复粘贴已知代码

***

**更多帮助**：**[快速入门指南](https://www.codebuddy.ai/docs/zh/cli/quickstart)**

**最后更新: 2026/4/23 14:36**

<br />

# **设置配置**

CodeBuddy Code 使用分层配置系统，让您能够在不同级别进行个性化定制，从个人偏好到团队标准，再到项目特定需求。

## **配置文件**

`settings.json` 文件是配置 CodeBuddy Code 的官方机制，支持分层设置：

- **用户设置** 定义在 `~/.codebuddy/settings.json`，应用于所有项目
- **项目设置** 保存在项目目录中：
  - `.codebuddy/settings.json` 用于检入源代码控制并与团队共享的设置
  - `.codebuddy/settings.local.json` 用于不检入的设置，适合个人偏好和实验。CodeBuddy Code 会自动配置 git 忽略此文件

### **完整配置示例**

```
{
  "language": "简体中文",
  "permissions": {
    "allow": [
      "Bash(npm run lint)",
      "Bash(npm run test:*)",
      "Read(~/.zshrc)"
    ],
    "ask": [
      "Bash(git push:*)"
    ],
    "deny": [
      "Bash(curl:*)",
      "Read(./.env)",
      "Read(./.env.*)",
      "Read(./secrets/**)"
    ]
  },
  "env": {
    "NODE_ENV": "development",
    "DEBUG": "codebuddy:*"
  },
  "model": "gpt-5",
  "cleanupPeriodDays": 30,
  "includeCoAuthoredBy": false,
  "statusLine": {
    "type": "command",
    "command": "~/.codebuddy/statusline.sh"
  }
}
```

## **可用设置**

`settings.json` 支持以下选项：

**配置键**

**描述**

**示例**

`language`

首选响应语言，设置后 CodeBuddy Code 将使用指定语言进行回复。留空则自动根据用户输入判断语言

`"简体中文"`

`apiKeyHelper`

自定义脚本，在 `/bin/sh` 中执行，生成认证值。此值将作为模型请求的 `X-Api-Key` 和 `Authorization: Bearer` 头发送

`/bin/generate_temp_api_key.sh`

`textToImageModel`

文生图功能使用的模型 ID

`"your-image-model"`

`imageToImageModel`

图生图功能使用的模型 ID

`"your-edit-model"`

`cleanupPeriodDays`

根据最后活动日期本地保留聊天记录的时长(默认:30 天)

`20`

`env`

应用于每个会话的环境变量

`{"FOO": "bar"}`

`includeCoAuthoredBy`

是否在 git 提交和拉取请求中包含 `co-authored-by CodeBuddy` 署名(默认:`true`）

`false`

`permissions`

权限配置，见下表

<br />

`hooks`

配置在工具执行前后运行的自定义命令。见 **[hooks 文档](https://www.codebuddy.ai/docs/zh/cli/hooks)**

`{"PreToolUse": {"Bash": "echo 'Running command...'"}}`

`disableAllHooks`

禁用所有 **[hooks](https://www.codebuddy.ai/docs/zh/cli/hooks)**

`true`

`model`

覆盖 CodeBuddy Code 使用的默认模型

`"gpt-5"`

`agent`

覆盖主线程使用的 agent 名称（内置或自定义 agent），应用该 agent 的 system prompt、工具限制和模型配置。优先级：`product.json default` → `plugin agent` → `settings.json agent` → `CLI --agent`

`"my-reviewer"`

`statusLine`

配置自定义状态行以显示上下文。见 \[statusLine 文档]\(#状态行配置）

`{"type": "command", "command": "~/.codebuddy/statusline.sh"}`

`enableAllProjectMcpServers`

自动批准项目 `.mcp.json` 文件中定义的所有 MCP 服务器

`false`

`enabledMcpjsonServers`

从 `.mcp.json` 文件批准的特定 MCP 服务器列表

`["memory", "github"]`

`disabledMcpjsonServers`

从 `.mcp.json` 文件拒绝的特定 MCP 服务器列表

`["filesystem"]`

`autoCompactEnabled`

开启自动压缩功能

`true`

`autoUpdates`

自动更新设置

`false`

`alwaysThinkingEnabled`

始终启用思考模式

`true`

`showTokensCounter`

是否在界面中显示 Tokens 计数器

`false`

`endpoint`

自定义服务端点地址

`"https://api.example.com"`

`envRouteMode`

环境路由模式配置

`"production"`

`sandbox`

Bash 沙箱配置,见**[Bash沙箱设置](https://www.codebuddy.ai/docs/zh/cli/settings#bash%E6%B2%99%E7%AE%B1%E8%AE%BE%E7%BD%AE)**

`{"enabled": true}`

`promptSuggestionEnabled`

启用 Prompt 建议功能，在 Agent 完成对话后自动预测下一步操作（默认：`true`）

`false`

`reasoningEffort`

Reasoning effort 级别配置，控制模型推理的深度。可选值：`low`、`medium`、`high`、`xhigh`。留空时使用产品配置默认值。可通过 `/config` 面板切换，选择 `auto` 等效于清除此设置

`"high"`

`memory`

\[Experimental] 记忆功能配置，见**[记忆功能配置](https://www.codebuddy.ai/docs/zh/cli/settings#%E8%AE%B0%E5%BF%86%E5%8A%9F%E8%83%BD%E9%85%8D%E7%BD%AEexperimental)**

`{"enabled": true}`

`trustedDirectories`

已经信任过的工作目录列表。命中的目录启动时不会再弹"是否信任此目录"的授权提示。通常由首次启动时的弹窗自动写入，也可手动编辑

`["~/workspace/myproj"]`

`trustAll`

信任所有工作目录，启动时不再弹"是否信任此目录"的授权提示。**仅免除目录信任授权，不会跳过工具执行权限**——是否弹工具审批仍由 `permissions.defaultMode` / `bypassPermissions` 模式决定，与本字段相互独立

`true`

### **权限设置**

**配置键**

**描述**

**示例**

`allow`

**[权限规则](https://www.codebuddy.ai/docs/zh/cli/iam#%E9%85%8D%E7%BD%AE%E6%9D%83%E9%99%90)**数组,允许工具使用。**注意:** Bash 规则使用前缀匹配,不是正则表达式

`[ "Bash(git diff:*)" ]`

`ask`

**[权限规则](https://www.codebuddy.ai/docs/zh/cli/iam#%E9%85%8D%E7%BD%AE%E6%9D%83%E9%99%90)**数组,在工具使用时询问确认

`[ "Bash(git push:*)" ]`

`deny`

**[权限规则](https://www.codebuddy.ai/docs/zh/cli/iam#%E9%85%8D%E7%BD%AE%E6%9D%83%E9%99%90)**数组,拒绝工具使用。用于排除 CodeBuddy Code 访问敏感文件。**注意:** Bash 模式是前缀匹配,可以被绕过(参见 **[Bash 权限限制](https://www.codebuddy.ai/docs/zh/cli/iam#%E5%B7%A5%E5%85%B7%E7%89%B9%E5%AE%9A%E7%9A%84%E6%9D%83%E9%99%90%E8%A7%84%E5%88%99)**)

`[ "WebFetch", "Bash(curl:*)", "Read(./.env)", "Read(./secrets/**)" ]`

`additionalDirectories`

CodeBuddy 可以访问的额外**[工作目录](https://www.codebuddy.ai/docs/zh/cli/iam#%E5%B7%A5%E4%BD%9C%E7%9B%AE%E5%BD%95)**

`[ "../docs/" ]`

`defaultMode`

打开 CodeBuddy Code 时的默认**[权限模式](https://www.codebuddy.ai/docs/zh/cli/iam#%E6%9D%83%E9%99%90%E6%A8%A1%E5%BC%8F)**

`"acceptEdits"`

`disableBypassPermissionsMode`

设置为 `"disable"` 以防止激活 `bypassPermissions` 模式。这会禁用 `-y` 和 `--dangerously-skip-permissions` 命令行标志

`"disable"`

`subagentPermissionMode`

覆盖 subagent/团队成员的默认权限模式。设置后所有 subagent 使用此模式，而非继承主 session 的模式。Agent 工具的 `mode` 参数优先级更高

`"bypassPermissions"`

### **记忆功能配置**

记忆功能允许 CodeBuddy Code 在会话之间保持持久化记忆，自动管理项目上下文和学习历史。

**配置键**

**描述**

**示例**

`autoMemoryEnabled`

是否启用 Auto Memory 功能（默认：`true`）。Auto Memory 允许 CodeBuddy 自动管理跨会话的持久化记忆，存储在 `~/.codebuddy/memories/` 目录

`true`

`typedMemory`

是否启用 Typed Memory 模式（默认：`true`）。启用后使用 4 种记忆类型（user/feedback/project/reference）+ YAML frontmatter 格式管理记忆

`true`

`relevanceSelection`

是否启用记忆相关性选择（默认：`true`）。启用后根据用户查询自动选择最多 5 个相关记忆注入上下文

`true`

`memoryExtraction`

是否启用后台记忆提取（默认：`false`）。启用后在对话结束时自动从对话中提取值得记住的信息

`true`

`teamMemory.enabled`

是否启用团队记忆模式（默认：`false`）。启用后，项目记忆存储在项目目录下，便于团队共享

`true`

`teamMemory.userId`

团队用户 ID，用于隔离不同用户的记忆。默认自动获取（git user.name > 系统用户名）

`"yangsubo"`

**配置示例：**

```
{
  "memory": {
    "autoMemoryEnabled": true,
    "typedMemory": true,
    "relevanceSelection": true,
    "memoryExtraction": false,
    "teamMemory": {
      "enabled": true,
      "userId": "yangsubo"
    }
  }
}
```

**记忆存储位置：**

- **个人模式**（默认）：`~/.codebuddy/memories/{project-id}/`
- **团队模式**：`{project}/.codebuddy/memories/@{user-id}/`
- **全局记忆**：`~/.codebuddy/memories/global/`

也可以通过 `/config` 命令在设置面板中启用此功能。

### **Bash沙箱设置**

配置高级沙箱行为。沙箱将 bash 命令与您的文件系统和网络隔离。详见 **[Bash 沙箱文档](https://www.codebuddy.ai/docs/zh/cli/bash-sandboxing)**。

**文件系统和网络限制**通过 Read、Edit 和 WebFetch 权限规则配置，而非通过这些沙箱设置。

**配置键**

**描述**

**示例**

`enabled`

启用 bash 沙箱(仅限 macOS/Linux)。默认:false

`true`

`autoAllowBashIfSandboxed`

在沙箱环境中自动批准 bash 命令。默认:true

`true`

`excludedCommands`

应在沙箱外运行的命令

`["git", "docker"]`

`allowUnsandboxedCommands`

允许通过 `dangerouslyDisableSandbox` 参数在沙箱外运行命令。设置为 `false` 时，完全禁用

<br />

`network.allowUnixSockets`

沙箱中可访问的 Unix 套接字路径（用于 SSH 代理等）

`["~/.ssh/agent-socket"]`

`network.allowLocalBinding`

允许绑定到 localhost 端口(仅限 macOS)。默认: false

`true`

`network.httpProxyPort`

如果您希望使用自己的代理,使用的 HTTP 代理端口。如果未指定,CodeBuddy 将运行自己的代理

`8080`

`network.socksProxyPort`

如果您希望使用自己的代理,使用的 SOCKS5 代理端口。如果未指定,CodeBuddy 将运行自己的代理

`8081`

`enableWeakerNestedSandbox`

为无特权的 Docker 环境启用较弱的沙箱(仅限 Linux)。**降低安全性。** 默认:false

`true`

**配置示例：**

```
{
  "sandbox": {
    "enabled": true,
    "autoAllowBashIfSandboxed": true,
    "excludedCommands": ["docker"],
    "network": {
      "allowUnixSockets": [
        "/var/run/docker.sock"
      ],
      "allowLocalBinding": true
    }
  },
  "permissions": {
    "deny": [
      "Read(.envrc)",
      "Read(~/.aws/**)"
    ]
  }
}
```

**文件系统访问**通过 Read/Edit 权限控制：

- Read deny 规则阻止沙箱中的文件读取
- Edit allow 规则允许文件写入（除默认值外，如当前工作目录）
- Edit deny 规则阻止允许路径内的写入

> **注意**：沙箱默认将 CodeBuddy 配置文件（`settings.json`、`settings.local.json`）加入写保护列表，防止沙箱内的命令或工具篡改配置。详见 **[Bash 沙箱 - 配置文件保护](https://www.codebuddy.ai/docs/zh/cli/bash-sandboxing#%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6%E4%BF%9D%E6%8A%A4)**。

**网络访问**通过 WebFetch 权限控制：

- WebFetch allow 规则允许网络域
- WebFetch deny 规则阻止网络域

### **设置优先级**

设置按优先级顺序应用（从高到低）:

1. **命令行参数**
   - 特定会话的临时覆盖
2. **本地项目设置** (`.codebuddy/settings.local.json`)
   - 个人项目特定设置
3. **共享项目设置** (`.codebuddy/settings.json`)
   - 源代码控制中的团队共享项目设置
4. **用户设置** (`~/.codebuddy/settings.json`)
   - 个人全局设置

此层次结构确保团队可以建立共享标准，同时仍允许个人自定义体验。

### **配置系统要点**

- **内存文件 （CODEBUDDY.md)**：包含 CodeBuddy 在启动时加载的指令和上下文
- **设置文件 （JSON)**：配置权限、环境变量和工具行为
- **斜杠命令**：可在会话期间使用 `/command-name` 调用的自定义命令
- **MCP 服务器**：使用额外工具和集成扩展 CodeBuddy Code
- **优先级**：更高级别的配置覆盖更低级别的配置
- **继承**：设置被合并，更具体的设置添加或覆盖更广泛的设置

### **排除敏感文件**

为防止 CodeBuddy Code 访问包含敏感信息的文件（如 API 密钥、秘密、环境文件），在 `.codebuddy/settings.json` 文件中使用 `permissions.deny` 设置：

```
{
  "permissions": {
    "deny": [
      "Read(./.env)",
      "Read(./.env.*)",
      "Read(./secrets/**)",
      "Read(./config/credentials.json)",
      "Read(./build)"
    ]
  }
}
```

匹配这些模式的文件将对 CodeBuddy Code 完全不可见，防止任何敏感数据的意外泄露。

## **子代理配置**

CodeBuddy Code 支持可在用户和项目级别配置的自定义 AI 子代理。这些子代理存储为带有 YAML frontmatter 的 Markdown 文件：

- **用户子代理**：`~/.codebuddy/agents/` - 在所有项目中可用
- **项目子代理**：`.codebuddy/agents/` - 特定于项目，可与团队共享

子代理文件定义具有自定义提示和工具权限的专用 AI 助手。详见 **[子代理文档](https://www.codebuddy.ai/docs/zh/cli/sub-agents)**。

## **插件配置**

CodeBuddy Code 支持插件系统，允许您使用自定义命令、代理、hooks 和 MCP 服务器扩展功能。插件通过市场分发，可在用户和项目级别配置。

### **插件设置**

`settings.json` 中的插件相关设置：

```
{
  "enabledPlugins": {
    "formatter@company-tools": true,
    "deployer@company-tools": true,
    "analyzer@security-plugins": false
  },
  "extraKnownMarketplaces": {
    "company-tools": {
      "source": {
        "source": "github",
        "repo": "company/codebuddy-plugins"
      }
    }
  }
}
```

#### **`enabledPlugins`**

控制启用哪些插件。格式：`"plugin-name@marketplace-name": true/false`

**作用域**：

- **用户设置** (`~/.codebuddy/settings.json`)：个人插件偏好
- **项目设置** (`.codebuddy/settings.json`)：与团队共享的项目特定插件
- **本地设置** (`.codebuddy/settings.local.json`)：每台机器的覆盖（不提交）

**示例**:

```
{
  "enabledPlugins": {
    "code-formatter@team-tools": true,
    "deployment-tools@team-tools": true,
    "experimental-features@personal": false
  }
}
```

#### **`extraKnownMarketplaces`**

定义应为项目提供的额外市场。通常在项目级设置中使用，以确保团队成员可以访问所需的插件源。

**当项目包含 `extraKnownMarketplaces` 时**:

1. 团队成员在信任文件夹时被提示安装市场
2. 然后团队成员被提示从该市场安装插件
3. 用户可以跳过不需要的市场或插件（存储在用户设置中）
4. 安装遵守信任边界并需要明确同意

**示例**:

```
{
  "extraKnownMarketplaces": {
    "company-tools": {
      "source": {
        "source": "github",
        "repo": "company-org/codebuddy-plugins"
      }
    },
    "security-plugins": {
      "source": {
        "source": "git",
        "url": "https://git.company.com/security/plugins.git"
      }
    }
  }
}
```

**市场源类型**:

- `github`: GitHub 仓库（使用 `repo`)
- `git`:任何 git URL(使用 `url`)
- `directory`:本地文件系统路径(使用 `path`,仅用于开发）

### **管理插件**

使用 `/plugin` 命令交互式管理插件：

- 浏览市场中的可用插件
- 安装/卸载插件
- 启用/禁用插件
- 查看插件详细信息（提供的命令、代理、hooks)
- 添加/删除市场

详见 **[插件文档](https://www.codebuddy.ai/docs/zh/cli/plugins)**。

## **环境变量**

CodeBuddy Code 支持通过环境变量来控制其行为。所有环境变量也可以在 **[`settings.json`](https://www.codebuddy.ai/docs/zh/cli/settings#%E5%8F%AF%E7%94%A8%E8%AE%BE%E7%BD%AE)** 的 `env` 字段中配置，这样可以自动为每个会话应用，或为整个团队推出配置。

完整的环境变量参考文档请参见 **[环境变量参考](https://www.codebuddy.ai/docs/zh/cli/env-vars)**。

### **快速入门**

**基础认证配置**：

```
# 使用 API 密钥
export CODEBUDDY_API_KEY="your-api-key"
codebuddy

# 或使用授权令牌
export CODEBUDDY_AUTH_TOKEN="your-token"
codebuddy
```

**设置代理**：

```
export HTTPS_PROXY="https://proxy.example.com:8080"
export NO_PROXY="localhost,127.0.0.1"
codebuddy
```

**启用高级功能**：

```
# 扩展思考
export MAX_THINKING_TOKENS="10000"

# 自动内存
export CODEBUDDY_DISABLE_AUTO_MEMORY="0"

codebuddy -p "你的查询"
```

### **在 settings.json 中配置**

环境变量也可以在 `settings.json` 的 `env` 字段中设置：

```
{
  "env": {
    "CODEBUDDY_API_KEY": "your-api-key",
    "HTTPS_PROXY": "https://proxy.example.com:8080",
    "MAX_THINKING_TOKENS": "10000"
  }
}
```

更多配置示例和高级用法，请参见 **[环境变量参考](https://www.codebuddy.ai/docs/zh/cli/env-vars)** 和 **[使用示例](https://www.codebuddy.ai/docs/zh/cli/env-vars#%E4%BD%BF%E7%94%A8%E7%A4%BA%E4%BE%8B)**。

## **状态行配置**

配置终端底部显示的状态行，可以显示当前会话、模型、成本等信息：

**配置键**

**类型**

**描述**

`statusLine.type`

string

状态行类型，目前支持 "command"

`statusLine.command`

string

执行的命令路径，支持 \~ 路径扩展

```
{
  "statusLine": {
    "type": "command",
    "command": "~/.codebuddy/statusline-script.sh"
  }
}
```

状态行命令会接收包含会话信息的 JSON 数据作为 stdin 输入，包括：

- `session_id`：会话 ID
- `model`：当前模型信息
- `workspace`：工作空间路径信息
- `cost`：成本统计信息
- `version`：应用版本

使用 `/statusline` 命令可以快速配置状态行。

## **配置管理命令**

使用 `codebuddy config` 命令管理配置：

### **基本语法**

```
codebuddy config [command] [options]
```

### **可用命令**

**命令**

**语法**

**描述**

`get`

`codebuddy config get <key>`

获取配置值

`set`

`codebuddy config set [options] <key> <value>`

设置配置值

`list`

`codebuddy config list`(别名:`ls`）

列出所有配置

`add`

`codebuddy config add <key> <values...>`

向数组配置添加项目

`remove`

`codebuddy config remove <key> [values...]`(别名:`rm`）

移除配置或数组项

### **选项**

**选项**

**描述**

**适用命令**

`-g, --global`

设置全局配置

`set`

### **使用示例**

#### **查看配置**

```
# 列出所有配置
codebuddy config list

# 获取特定配置值
codebuddy config get model
codebuddy config get permissions
```

#### **设置配置**

```
# 设置项目级模型（不需要 -g 标志）
codebuddy config set model gpt-5

# 设置全局模型（需要 -g 标志）
codebuddy config set -g model gpt-4

# 设置项目级权限配置（不需要 -g 标志）
codebuddy config set permissions '{"allow": ["Read", "Edit"], "deny": ["Bash(rm:*)"]}'

# 设置项目级环境变量（不需要 -g 标志）
codebuddy config set env '{"NODE_ENV": "development", "DEBUG": "true"}'

# 设置全局专用配置（需要 -g 标志）
codebuddy config set -g cleanupPeriodDays 30
codebuddy config set -g includeCoAuthoredBy false
```

## **CodeBuddy 可用的工具**

CodeBuddy Code 可以访问一组强大的工具，帮助它理解和修改您的代码库：

**工具**

**描述**

**需要权限**

**AskUserQuestion**

向用户询问多选问题以收集信息或澄清歧义

否

**Bash**

在您的环境中执行 shell 命令

是

**TaskOutput**

从正在运行或已完成的后台任务检索输出

否

**Edit**

对特定文件进行有针对性的编辑

是

**MultiEdit**

在单个操作中对单个文件进行多次编辑

是

**ExitPlanMode**

提示用户退出计划模式并开始编码

是

**Glob**

基于模式匹配查找文件

否

**Grep**

在文件内容中搜索模式

否

**TaskStop**

通过 ID 终止正在运行的后台任务

否

**LSP**

与 LSP 服务器交互获取代码智能功能（跳转定义、查找引用、悬停信息等）

否

**NotebookEdit**

修改 Jupyter notebook 单元格

是

**Read**

读取文件内容

否

**Skill**

在主对话中执行技能

是

**SlashCommand**

运行**[自定义斜杠命令](https://www.codebuddy.ai/docs/zh/cli/slash-commands#slashcommand-%E5%B7%A5%E5%85%B7)**

是

**Task**

运行子代理以处理复杂的多步骤任务

否

**TaskOutput**

从正在运行或已完成的后台任务检索输出

否

**TaskCreate**

创建任务以跟踪工作进度

否

**TaskUpdate**

更新任务状态（pending/in\_progress/completed）

否

**TaskList**

列出当前任务

否

**TaskGet**

获取特定任务详情

否

**WebFetch**

从指定 URL 获取内容

是

**WebSearch**

执行带域过滤的网络搜索

是

**Write**

创建或覆盖文件

是

权限规则可以使用 `/permissions` 或在**[权限设置](https://www.codebuddy.ai/docs/zh/cli/settings#%E6%9D%83%E9%99%90%E8%AE%BE%E7%BD%AE)**中配置。另见**[工具特定的权限规则](https://www.codebuddy.ai/docs/zh/cli/iam#%E5%B7%A5%E5%85%B7%E7%89%B9%E5%AE%9A%E7%9A%84%E6%9D%83%E9%99%90%E8%A7%84%E5%88%99)**。

### **使用 hooks 扩展工具**

您可以使用 **[CodeBuddy Code hooks](https://www.codebuddy.ai/docs/zh/cli/hooks)** 在任何工具执行前后运行自定义命令。

例如，您可以在 CodeBuddy 修改 Python 文件后自动运行 Python 格式化程序，或通过阻止对某些路径的 Write 操作来防止修改生产配置文件。

## **常见配置场景**

### **团队协作配置**

**项目共享配置**（`.codebuddy/settings.json`）：

```
{
  "model": "gpt-5",
  "permissions": {
    "allow": ["Read", "Edit", "Bash(git:*)", "Bash(npm:*)"],
    "ask": ["WebFetch", "Bash(docker:*)"],
    "deny": ["Bash(rm:*)", "Bash(sudo:*)"]
  },
  "env": {
    "NODE_ENV": "development"
  }
}
```

**个人本地配置**（`.codebuddy/settings.local.json`）：

```
{
  "model": "gpt-4",
  "env": {
    "DEBUG": "myapp:*"
  }
}
```

### **安全配置**

限制敏感操作和文件访问：

```
{
  "permissions": {
    "allow": ["Read", "Edit(src/**)", "Bash(git:status,git:diff)"],
    "ask": ["WebFetch", "Bash(curl:*)"],
    "deny": [
      "Edit(**/*.env)",
      "Edit(**/*.key)",
      "Edit(**/*.pem)",
      "Bash(wget:*)",
      "Read(/etc/**)",
      "Read(~/.ssh/**)"
    ],
    "defaultMode": "default"
  }
}
```

### **沙箱安全配置**

启用沙箱并配置文件系统和网络访问：

```
{
  "sandbox": {
    "enabled": true,
    "autoAllowBashIfSandboxed": true,
    "excludedCommands": ["docker", "git"],
    "network": {
      "allowUnixSockets": ["/var/run/docker.sock"],
      "allowLocalBinding": true
    }
  },
  "permissions": {
    "allow": [
      "Edit(src/**)",
      "WebFetch(https://api.github.com/**)"
    ],
    "deny": [
      "Read(.envrc)",
      "Read(~/.aws/**)",
      "Edit(**/*.env)"
    ]
  }
}
```

## **另见**

- **[身份和访问管理](https://www.codebuddy.ai/docs/zh/cli/iam#%E9%85%8D%E7%BD%AE%E6%9D%83%E9%99%90)** - 了解 CodeBuddy Code 的权限系统
- **[Bash 沙箱](https://www.codebuddy.ai/docs/zh/cli/bash-sandboxing)** - 了解沙箱隔离功能
- **[故障排除](https://www.codebuddy.ai/docs/zh/cli/troubleshooting)** - 常见配置问题的解决方案

***

*合适的配置让 CodeBuddy Code 更懂您的需求 ⚙️*

**最后更新: 2026/4/27 08:00**

<br />

# **models.json 配置指南**

## **概述**

`models.json` 是一个配置文件，用于自定义模型列表和控制模型下拉列表的显示。该配置支持两个级别：

- **用户级**: `~/.codebuddy/models.json` - 全局配置，适用于所有项目
- **项目级**: `<workspace>/.codebuddy/models.json` - 项目特定配置，优先级高于用户级

## **配置文件位置**

### **用户级配置**

```
~/.codebuddy/models.json
```

### **项目级配置**

```
<project-root>/.codebuddy/models.json
```

## **配置优先级**

配置合并优先级从高到低：

1. 项目级 models.json
2. 用户级 models.json
3. 内置默认配置

项目级配置会覆盖用户级配置中的相同模型定义（基于 `id` 字段匹配）。`availableModels` 字段：项目级完全覆盖用户级，不进行合并。

## **配置结构**

```
{
  "models": [
    {
      "id": "model-id",
      "name": "Model Display Name",
      "vendor": "vendor-name",
      "apiKey": "sk-actual-api-key-value",
      "maxInputTokens": 200000,
      "maxOutputTokens": 8192,
      "url": "https://api.example.com/v1/chat/completions",
      "temperature": 0.7,
      "supportsToolCall": true,
      "supportsImages": true
    }
  ],
  "availableModels": ["model-id-1", "model-id-2"]
}
```

## **配置字段说明**

### **models**

类型： `Array<LanguageModel>`

定义自定义模型列表。可以添加新模型或覆盖内置模型配置。

#### **LanguageModel 字段**

**字段**

**类型**

**必填**

**说明**

`id`

string

✓

模型唯一标识符

`name`

string

\-

模型显示名称

`vendor`

string

\-

模型供应商 （如 OpenAI, Google)

`apiKey`

string

\-

API 密钥，支持环境变量引用（见下方安全配置说明）

`maxInputTokens`

number

\-

最大输入 token 数

`maxOutputTokens`

number

\-

最大输出 token 数

`url`

string

\-

API 端点 URL，支持环境变量引用 (必须是接口完整路径,一般以 `/chat/completions` 结尾）

`temperature`

number

\-

采样温度，范围 0-2，值越高输出越随机，值越低输出越确定

`supportsToolCall`

boolean

\-

是否支持工具调用

`supportsImages`

boolean

\-

是否支持图片输入

`supportsReasoning`

boolean

\-

是否支持推理模式

`relatedModels`

object

\-

关联模型配置，指定在不同场景（`lite`/`reasoning`/`vision`/`longContext`/`subagent`）下使用哪个模型 id。详见**[配置关联模型](https://www.codebuddy.ai/docs/zh/cli/models#%E9%85%8D%E7%BD%AE%E5%85%B3%E8%81%94%E6%A8%A1%E5%9E%8B)**

**重要说明：**

- 目前仅支持 OpenAI 接口格式的 API
- `url` 字段必须是接口完整路径,一般以 `/chat/completions` 结尾
- 例如: `https://api.openai.com/v1/chat/completions` 或 `http://localhost:11434/v1/chat/completions`

### **安全配置：使用环境变量引用**

为避免 API 密钥明文存储在配置文件中，`apiKey` 和 `url` 字段支持环境变量引用语法 `${VAR_NAME}`。

**语法格式：**

```
${环境变量名}
```

**配置示例：**

```
{
  "models": [
    {
      "id": "gpt-4o",
      "name": "GPT-4o",
      "vendor": "OpenAI",
      "apiKey": "${OPENAI_API_KEY}",
      "url": "https://api.openai.com/v1/chat/completions"
    }
  ]
}
```

**设置环境变量：**

```
# 在 ~/.zshrc 或 ~/.bashrc 中添加
export OPENAI_API_KEY="your-openai-api-key"

# 或者在启动时临时设置
OPENAI_API_KEY="your-openai-api-key" codebuddy
```

**使用系统 Keychain（macOS）：**

```
# 存储密钥到 Keychain
security add-generic-password -a "$USER" -s "openai-api-key" -w "your-openai-api-key"

# 在 ~/.zshrc 中配置自动导出
export OPENAI_API_KEY=$(security find-generic-password -s "openai-api-key" -w 2>/dev/null)
```

**注意事项：**

- 环境变量在 CLI 启动时解析
- 如果环境变量不存在，将保留原始占位符（会导致 API 调用失败）
- 建议将 `models.json` 文件权限设置为 `600`（仅所有者可读写）
- 不要将包含实际密钥的配置文件提交到版本控制系统

### **availableModels**

类型： `Array<string>`

控制模型下拉列表中显示哪些模型。只有在此数组中列出的模型 ID 才会在 UI 中显示。

- 如果未配置或为空数组，则显示所有模型
- 配置后，只显示列出的模型 ID
- 可以同时包含内置模型和自定义模型的 ID

## **使用场景**

### **1. 添加自定义模型**

在用户级或项目级添加新的模型配置：

```
{
  "models": [
    {
      "id": "my-custom-model",
      "name": "My Custom Model",
      "vendor": "OpenAI",
      "apiKey": "sk-custom-key-here",
      "maxInputTokens": 128000,
      "maxOutputTokens": 4096,
      "url": "https://api.myservice.com/v1/chat/completions",
      "supportsToolCall": true
    }
  ]
}
```

### **2. 覆盖内置模型配置**

修改内置模型的默认参数：

```
{
  "models": [
    {
      "id": "gpt-4-turbo",
      "name": "GPT-4 Turbo (Custom Endpoint)",
      "vendor": "OpenAI",
      "url": "https://my-proxy.example.com/v1/chat/completions",
      "apiKey": "sk-your-key-here"
    }
  ]
}
```

### **3. 限制可用模型列表**

只在下拉列表中显示特定模型：

```
{
  "availableModels": [
    "gpt-4-turbo",
    "gpt-4o",
    "my-custom-model"
  ]
}
```

### **4. 项目特定配置**

为特定项目使用不同的模型或 API 端点：

**项目 A** (`.codebuddy/models.json`):

```
{
  "models": [
    {
      "id": "project-a-model",
      "name": "Project A Model",
      "vendor": "OpenAI",
      "url": "https://project-a-api.example.com/v1/chat/completions",
      "apiKey": "sk-project-a-key",
      "maxInputTokens": 100000,
      "maxOutputTokens": 4096
    }
  ],
  "availableModels": ["project-a-model", "gpt-4-turbo"]
}
```

## **热重载**

配置文件支持热重载：

- 文件变更会被自动检测
- 使用 1 秒防抖延迟避免频繁重载
- 配置更新后会自动同步到应用

监听的文件：

- `~/.codebuddy/models.json` （用户级）
- `<workspace>/.codebuddy/models.json` （项目级）

## **标签系统**

通过 `models.json` 添加的模型会自动标记 `custom` 标签，便于在 UI 中识别和过滤。

## **合并策略**

配置使用 `SmartMerge` 策略：

- 相同 ID 的模型配置会被覆盖
- 不同 ID 的模型会被追加
- 项目级配置优先于用户级配置
- `availableModels` 过滤在所有合并完成后执行

## **示例配置**

### **API 端点 URL 格式说明**

**必须使用完整路径：** 所有自定义模型的 `url` 字段一般以 `/chat/completions` 结尾。

✅ **正确示例：**

```
https://api.openai.com/v1/chat/completions
https://api.myservice.com/v1/chat/completions
http://localhost:11434/v1/chat/completions
https://my-proxy.example.com/v1/chat/completions
```

❌ **错误示例：**

```
https://api.openai.com/v1
https://api.myservice.com
http://localhost:11434
```

### **OpenRouter 平台配置示例**

使用 OpenRouter 访问多种模型：

```
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

### **DeepSeek 平台配置示例**

使用 DeepSeek 模型（配置 `url` 后，即使与云端同 id 也会按"完全替换"语义生效，不会被云端默认项合并覆盖）：

```
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
      "supportsImages": false
    },
    {
      "id": "deepseek-v4-flash",
      "name": "DeepSeek V4 Flash",
      "vendor": "DeepSeek",
      "url": "https://api.deepseek.com/v1/chat/completions",
      "apiKey": "${DEEPSEEK_API_KEY}",
      "maxInputTokens": 128000,
      "maxOutputTokens": 8192,
      "supportsToolCall": true,
      "supportsImages": false
    }
  ],
  "availableModels": [
    "deepseek-v4-pro",
    "deepseek-v4-flash"
  ]
}
```

设置 API 密钥环境变量后启动：

```
export DEEPSEEK_API_KEY="<your-deepseek-api-key>"
codebuddy --model deepseek-v4-pro
```

> **提示**：如果不希望维护 `models.json`，也可以完全通过环境变量对接 DeepSeek，见 **[env-vars.md 对接 DeepSeek 示例](https://www.codebuddy.ai/docs/zh/cli/env-vars#%E5%AF%B9%E6%8E%A5-deepseek-%E7%A4%BA%E4%BE%8B)**。

### **配置关联模型**

CodeBuddy Code 在一次会话中会根据场景切换模型，避免用大模型处理简单任务、或用通用模型处理需要推理 / 视觉 / 长上下文的请求。这些场景通过模型条目的 `relatedModels` 字段声明。

**支持的场景（variant type）：**

**场景**

**用途**

**当前状态**

`lite`

轻量快速模型，用于后台提取、摘要等低价值请求；也是 Agent 工具 `model: "lite"` 参数对应的模型

**已生效**

`reasoning`

推理增强模型，用于需要深度思考的复杂推理；Agent 工具 `model: "reasoning"` 参数对应的模型

**已生效**

`subagent`

子 Agent / 团队成员默认使用的模型

**预留未启用**——当前子代理模型通过 `CODEBUDDY_CODE_SUBAGENT_MODEL` 环境变量或 agent 配置的 `models[0]` 决定，不读取本字段

`vision`

视觉理解模型，用于需要处理图片的请求

**预留未启用**——类型已定义，尚无调用点消费此 variant

`longContext`

长上下文模型，用于上下文超长的请求

**预留未启用**——类型已定义，尚无调用点消费此 variant

> **当前实际可用**：只有 `lite` 和 `reasoning` 两个 variant 在 agent-manager 中被消费并映射到模型切换逻辑。`subagent` / `vision` / `longContext` 三项仅保留在类型定义中，为后续迭代预留，现在写进 `relatedModels` 不会报错但也不会生效。

**关键规则（自定义模型必读）：**

> 通过 `models.json` 添加的自定义模型**不会继承**产品内置的 `defaultRelatedModels`。 如果没有在自身条目里显式声明 `relatedModels`，所有场景都会回退到主模型自己——也就是说子代理、lite、reasoning 全部用同一个大模型跑，成本和速度都不划算。

**配置示例（DeepSeek 主模型 + flash 作为 lite / reasoning）：**

```
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
  "availableModels": [
    "deepseek-v4-pro",
    "deepseek-v4-flash"
  ]
}
```

**解析优先级（从高到低，当前仅 `lite` / `reasoning` 会走到这个解析链）：**

1. 环境变量显式指定（`CODEBUDDY_SMALL_FAST_MODEL` 对应 `lite`、`CODEBUDDY_BIG_SLOW_MODEL` 对应 `reasoning`）
2. 当前主模型条目里的 `relatedModels[variant]`
3. 产品内置的 `defaultRelatedModels[variant]`（**仅对内置模型生效，自定义模型跳过这步**）
4. 回落到主模型自身

> **子代理模型的取值规则不走 `relatedModels.subagent`**，而是独立的链路：Agent 工具 `model` 参数 > `CODEBUDDY_CODE_SUBAGENT_MODEL` 环境变量 > agent 配置的 `models[0]` > 主模型。

**与环境变量方式的取舍：**

- 想让某主模型在 `lite` / `reasoning` 场景自动切到另一个同厂小模型：在主模型条目里声明 `relatedModels` 最直观，绑定关系跟着模型走。
- 想让子代理用小模型：目前必须走环境变量 `CODEBUDDY_CODE_SUBAGENT_MODEL`，或在 agent 配置里把小模型写到 `models[0]`——**不是** `relatedModels.subagent`。
- 想在不同项目里用不同的大小模型组合：用环境变量（`CODEBUDDY_MODEL` / `CODEBUDDY_BIG_SLOW_MODEL` / `CODEBUDDY_SMALL_FAST_MODEL` / `CODEBUDDY_CODE_SUBAGENT_MODEL`），配合项目级 `.env` 切换。
- 两种方式同时生效时，环境变量优先。

### **完整示例**

```
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
  "availableModels": [
    "gpt-4o",
    "my-local-llm"
  ]
}
```

## **故障排查**

### **配置未生效**

1. 检查 JSON 格式是否正确
2. 确认文件路径是否正确
3. 查看日志输出确认配置是否被加载
4. 确认环境变量中的 API 密钥是否已设置

### **模型未在列表中显示**

1. 检查模型 ID 是否在 `availableModels` 中列出
2. 确认 `models` 配置是否正确
3. 验证必填字段 （`id`, `name`, `provider`) 是否都已提供

### **热重载未触发**

- 配置文件变更有 1 秒防抖延迟
- 确保文件确实被保存到磁盘
- 检查文件监听是否正常启动 （查看调试日志）

**最后更新: 2026/4/27 08:00**

<br />

# **优化终端配置**

> CodeBuddy Code 在合适的终端配置下能发挥最佳性能。按照以下指南优化您的使用体验。

## **🎨 主题和外观**

CodeBuddy Code 无法控制您终端的主题，这由您的终端应用程序处理。您可以随时通过 `/config` 命令将 CodeBuddy Code 的主题与终端主题匹配。

如需进一步自定义 CodeBuddy Code 界面本身，您可以配置**[自定义状态行](https://www.codebuddy.ai/docs/zh/cli/statusline)**在终端底部显示上下文信息，如当前模型、工作目录或 git 分支等。

## **⌨️ 换行输入**

在 CodeBuddy Code 中输入换行符有以下几种方式：

- **Ctrl+J**:跨平台快捷键，在多行输入时按 Ctrl+J 即可插入换行
- **快速转义**:输入 `\` 后按 Enter 即可创建换行
- **Shift+Enter**:通过 `/terminal-setup` 命令自动配置（推荐）

### **配置 Shift+Enter（推荐）:**

在 CodeBuddy Code 中运行 `/terminal-setup` 命令，自动配置 Shift+Enter 快捷键。

支持的终端：

- **macOS**: iTerm2, Terminal.app
- **IDE**: VSCode, Cursor, Windsurf, Zed, CodeBuddy
- **JetBrains**: PyCharm, IntelliJ IDEA, WebStorm, PhpStorm, GoLand, Rider, CLion, RubyMine, AppCode, DataGrip
- **终端模拟器**: Ghostty, WezTerm, Kitty, Alacritty, Hyper, Tabby, Warp
- **Windows**: Windows Terminal

### **配置 Option+Enter(VS Code、iTerm2 或 macOS Terminal.app):**

**对于 Mac Terminal.app:**

1. 打开 设置 → 描述文件 → 键盘
2. 勾选"将 Option 用作 Meta 键"

**对于 iTerm2 和 VS Code 终端：**

1. 打开 Settings → Profiles → Keys
2. 在"General"下，将 Left/Right Option 键设置为"Esc+"

## **🔔 通知设置**

通过合适的通知配置，让您在 CodeBuddy 完成任务时不会错过：

### **iTerm 2 系统通知**

为 iTerm 2 配置任务完成时的提醒：

1. 打开 iTerm 2 偏好设置
2. 导航到 Profiles → Terminal
3. 启用"Silence bell"和 Filter Alerts → “Send escape sequence-generated alerts”
4. 设置您偏好的通知延迟

注意：这些通知功能仅适用于 iTerm 2,在 macOS 默认终端中不可用。

### **自定义通知 Hook**

如需高级通知处理，您可以创建**[通知 Hook](https://www.codebuddy.ai/docs/zh/cli/hooks#notification)**来运行自己的逻辑。

## **📝 处理大量输入**

在处理大量代码或长指令时：

- **避免直接粘贴**:CodeBuddy Code 可能无法很好地处理非常长的粘贴内容
- **使用基于文件的工作流**:将内容写入文件并要求 CodeBuddy 读取它
- **注意 VS Code 的限制**:VS Code 终端特别容易截断长粘贴内容

## **⌨️ Vim 模式**

CodeBuddy Code 支持部分 Vim 快捷键，可以通过 `/vim` 命令启用,或通过 `/config` 配置。

支持的功能子集包括：

- 模式切换：`Esc`（切换到 NORMAL 模式）、`i`/`I`、`a`/`A`、`o`/`O`（切换到 INSERT 模式）
- 导航：`h`/`j`/`k`/`l`、`w`/`e`/`b`、`0`/`$`/`^`、`gg`/`G`
- 编辑：`x`、`dw`/`de`/`db`/`dd`/`D`、`cw`/`ce`/`cb`/`cc`/`C`、`.`（重复）

# **环境变量参考**

CodeBuddy Code 支持通过环境变量来控制其行为。这些变量可以在启动前设置，也可以在 **[`settings.json`](https://www.codebuddy.ai/docs/zh/cli/settings#%E5%8F%AF%E7%94%A8%E8%AE%BE%E7%BD%AE)** 的 `env` 字段中配置以应用到每个会话。

> **提示**：所有环境变量也可以在 `settings.json` 的 `env` 字段中设置，这样可以自动为每个会话应用，或为整个团队推出配置。

## **认证相关**

**环境变量**

**说明**

`CODEBUDDY_API_KEY`

API 密钥。设置此密钥用于模型接口调用。在非交互模式 (`-p`) 下始终使用此密钥

`CODEBUDDY_AUTH_TOKEN`

CodeBuddy 平台认证令牌，用于所有平台接口调用

`CODEBUDDY_CUSTOM_HEADERS`

自定义 HTTP 请求头。格式：`Name: Value`，多个请求头用换行符或 `\n` 分隔

## **API 端点和代理**

**环境变量**

**说明**

`CODEBUDDY_BASE_URL`

覆盖 API 端点地址，通常与 `CODEBUDDY_API_KEY` 配合使用

`CODEBUDDY_INTERNET_ENVIRONMENT`

网络环境配置（`internal` 用于中国版，`ioa` 用于 iOA 企业版）

`HTTP_PROXY` / `http_proxy`

HTTP 代理服务器地址

`HTTPS_PROXY` / `https_proxy`

HTTPS 代理服务器地址

`NO_PROXY` / `no_proxy`

绕过代理的域和 IP 列表（逗号分隔，如 `localhost,.example.com`）

## **模型配置**

**环境变量**

**说明**

`CODEBUDDY_MODEL`

覆盖默认代理模型

`CODEBUDDY_SMALL_FAST_MODEL`

后台任务使用的小型快速模型

`CODEBUDDY_BIG_SLOW_MODEL`

复杂推理任务使用的大型模型

`CODEBUDDY_CODE_SUBAGENT_MODEL`

子代理使用的模型

`MAX_THINKING_TOKENS`

启用扩展思考并设置思考过程的 token 预算。默认禁用

## **Bash 工具配置**

**环境变量**

**说明**

`BASH_DEFAULT_TIMEOUT_MS`

长时间运行 bash 命令的默认超时（默认：120000）

`BASH_MAX_OUTPUT_LENGTH`

bash 输出在内存中保留的最大字符数（默认：30000，上限：150000）。超出部分会被中间截断（保留 head 20% + tail 80%），完整输出会自动保存到磁盘

`BASH_MAX_TIMEOUT_MS`

模型可为长时间运行的 bash 命令设置的最大超时（默认：600000）

## **工具输出外部化**

**环境变量**

**说明**

`CODEBUDDY_TOOL_RESULT_THRESHOLD_KB`

工具结果外部化的大小阈值（KB），超过此阈值的非 bash 工具结果会被保存到磁盘并替换为占位符（默认：50）

> **说明**：Bash 工具的输出外部化由 `BASH_MAX_OUTPUT_LENGTH` 控制，当输出超过该值发生截断时，完整输出自动流式写入磁盘。`CODEBUDDY_TOOL_RESULT_THRESHOLD_KB` 主要影响其他工具（如 MCP 工具）的大输出处理。详见**[工具输出外部化](https://www.codebuddy.ai/docs/zh/cli/env-vars#%E5%B7%A5%E5%85%B7%E8%BE%93%E5%87%BA%E5%A4%96%E9%83%A8%E5%8C%96%E6%9C%BA%E5%88%B6)**章节。

## **工具和功能开关**

**环境变量**

**说明**

`CODEBUDDY_DISABLE_HOT_RELOAD`

设置为 `1` 禁用热更新系统

`CODEBUDDY_SKIP_BUILTIN_MARKETPLACE`

设置为 `1` 跳过内置插件市场加载

`CODEBUDDY_AUTO_UPDATE_THIRD_PARTY_MARKETPLACES`

设置为 `true` 或 `1` 启用第三方插件市场自动更新（默认：禁用）

`CODEBUDDY_PLUGIN_DIRS`

冒号分隔的本地插件目录路径列表（等同于 `--plugin-dir`），插件的 `bin/` 目录会自动注入到 `PATH`

`CODEBUDDY_IMAGE_GEN_ENABLED`

设置为 `false` 或 `0` 禁用图片生成功能

`CODEBUDDY_IMAGE_EDIT_ENABLED`

设置为 `false` 或 `0` 禁用图片编辑功能

`CODEBUDDY_DEFER_TOOL_LOADING`

设置为 `false` 或 `0` 禁用 MCP 工具延迟加载

`CODEBUDDY_SHOW_ALL_DEFERRED_TOOLS`

设置为 `true` 或 `1` 显示所有延迟工具的完整描述

`CODEBUDDY_DISABLE_CRON`

设置为 `1` 禁用计划任务

`CODEBUDDY_REHYDRATE_IMAGE_BLOB_REFS`

设置为 `true` 在 `-p` 模式流式输出中将图片 blob 引用还原为完整 base64 数据。适用于需要直接获取图片数据的下游集成场景

## **上下文和内存**

**环境变量**

**说明**

`CODEBUDDY_AUTOCOMPACT_PCT_OVERRIDE`

设置自动压缩触发的上下文容量百分比（1-100）。默认由产品配置决定（通常 70-92%）。使用更低的值（如 `50`）来更早压缩

`CODEBUDDY_PRE_MESSAGE_COMPACT_PCT`

设置用户发消息前预检压缩的上下文容量百分比（1-100）。默认 10%。当上下文超过此阈值时，在处理用户新消息前自动压缩

`CODEBUDDY_DISABLE_AUTO_MEMORY`

设置为 `1` 禁用自动内存，设置为 `0` 启用

`CODEBUDDY_MEMORY_ENABLED`

设置为 `true` 或 `1` 启用记忆功能

`CODEBUDDY_TYPED_MEMORY_ENABLED`

设置为 `true` 或 `1` 启用分类记忆模式

`CODEBUDDY_TEAM_MEMORY_ENABLED`

设置为 `true` 或 `1` 启用团队记忆模式

`CODEBUDDY_USER_ID`

团队记忆模式下的用户 ID

## **MCP (Model Context Protocol)**

**环境变量**

**说明**

`MCP_TIMEOUT`

MCP 服务器连接的超时时间（毫秒）

`MCP_TOOL_TIMEOUT`

MCP 工具执行的超时时间（毫秒）

`MAX_MCP_OUTPUT_TOKENS`

MCP 工具响应中允许的最大 token 数（默认：20000）

## **性能和输出**

**环境变量**

**说明**

`CODEBUDDY_CODE_MAX_OUTPUT_TOKENS`

设置大多数请求的最大输出 token 数

`CODEBUDDY_CODE_FILE_READ_MAX_OUTPUT_TOKENS`

覆盖文件读取的默认 token 限制（默认：20000）

`CODEBUDDY_STREAM_TIMEOUT_MS`

流式响应中两个数据块之间允许的最大静默时间（毫秒）（默认：120000）

`CODEBUDDY_FIRST_TOKEN_TIMEOUT_MS`

等待第一个模型输出的最大时间（毫秒）（默认：120000）

`CODEBUDDY_SESSION_MAX_ITEMS`

`session/load` 回放时历史消息的最大条数（默认：1000）。达到阈值且遇到 user message 时停止逆序读取 JSONL。需要支持超长会话（如沙箱场景）时可调大（例如 2000 或更多）；零/负数/非数字会回退到默认值

## **文件系统和配置**

**环境变量**

**说明**

`CODEBUDDY_CONFIG_DIR`

自定义 CodeBuddy Code 存储配置和数据文件的位置

`CODEBUDDY_CODE_DEBUG_LOGS_DIR`

调试日志目录

`CODEBUDDY_SANDBOX_IMAGE`

容器沙箱镜像（默认：`node:20-alpine`）

`USE_BUILTIN_RIPGREP`

设置为 `0` 使用系统安装的 `rg` 而不是 CodeBuddy Code 附带的 `rg`

## **Shell 配置**

**环境变量**

**说明**

`CODEBUDDY_CODE_SHELL`

覆盖自动 shell 检测。支持的值：`bash`、`zsh`、`sh`、`powershell`

`CODEBUDDY_CODE_SHELL_PREFIX`

包装所有 shell 命令的命令前缀（如用于日志或审计）

`CODEBUDDY_CODE_GIT_BASH_PATH`

Windows 下显式指定 Git Bash 路径；若指定的路径无效则启动失败

`CODEBUDDY_SKIP_GIT_BASH_CHECK`

设置为 `1` 跳过启动时的 Windows Git Bash 检测和提示（适用于上游已管理 shell 的场景）

`CODEBUDDY_POWERSHELL_PATH`

显式指定 PowerShell 可执行文件路径（优先于自动检测）

`CODEBUDDY_USE_POWERSHELL_TOOL`

控制 PowerShell 工具启用状态。Windows 上默认启用，设为 `0` 可禁用

`CODEBUDDY_ENV_FILE`

在执行每个 shell 命令前自动 source 的环境文件路径

`CODEBUDDY_DISABLE_SHELL_SNAPSHOT`

设置为 `1` 全平台禁用 shell 环境快照（在 bash profile 加载慢时有效）

`CODEBUDDY_ENABLE_SHELL_SNAPSHOT`

Windows 无 Git Bash 时默认跳过 snapshot；设为 `1` 可强制启用（通常不需要）

## **UI 和交互**

**环境变量**

**说明**

`CODEBUDDY_CODE_DISABLE_TERMINAL_TITLE`

设置为 `1` 禁用自动终端标题更新

`CODEBUDDY_PROMPT_SUGGESTION_DISABLED`

设置为 `true` 禁用提示建议

`IS_DEMO`

设置为 `true` 启用演示模式：隐藏邮箱和组织

## **安全和认证**

**环境变量**

**说明**

`CODEBUDDY_CODE_CLIENT_CERT`

mTLS 客户端证书文件路径 ⚠️ *暂未支持*

`CODEBUDDY_CODE_CLIENT_KEY`

mTLS 客户端私钥文件路径 ⚠️ *暂未支持*

`CODEBUDDY_CODE_CLIENT_KEY_PASSPHRASE`

mTLS 加密私钥的密码（可选）⚠️ *暂未支持*

## **遥测和报告**

**环境变量**

**说明**

`DISABLE_TELEMETRY`

设置为 `1` 禁用遥测

`DISABLE_ERROR_REPORTING`

设置为 `1` 禁用错误报告

`DISABLE_AUTOUPDATER`

设置为 `1` 禁用自动更新

`DISABLE_FEEDBACK_COMMAND`

设置为 `1` 禁用 `/feedback` 命令

## **任务和后台工作**

**环境变量**

**说明**

`CODEBUDDY_DISABLE_BACKGROUND_TASKS`

设置为 `1` 禁用所有后台任务功能

## **Daemon 模式**

**环境变量**

**说明**

`CODEBUDDY_DAEMON_ALLOW_SLEEP`

设置为 `1` 或 `true` 禁用 daemon 的防休眠功能（允许系统正常进入 idle sleep）。唤醒后的自动重连不受影响

`CODEBUDDY_DAEMON_AUTO_CONNECT_CHANNELS`

设置为 `0` 或 `false` 禁用 daemon 启动时自动连接微信/企微 channel。也可通过 `settings.json` 的 `daemonAutoConnectChannels: false` 配置。单个实例可在 `instances.json` 中设置 `autoConnect: false`

`CODEBUDDY_DAEMON_RESTORE_CHANNELS`

daemon 自动重启时由系统设置，包含需要恢复的 channel 列表（逗号分隔，如 `wechat:abc,wecom:default`）。用户通常无需手动设置

## **Agent 执行控制**

**环境变量**

**说明**

`CODEBUDDY_CODE_MAX_TURNS`

主 Agent 的最大执行轮次。优先级：CLI `--max-turns` > 此环境变量 > 默认值 (500)

`CODEBUDDY_CODE_SUBAGENT_MAX_TURNS`

子 Agent 的最大执行轮次。优先级：CLI `--max-turns` > 此环境变量 > 模型动态传入的 `max_turns` > 默认值 (500)

`CODEBUDDY_SUBAGENT_PERMISSION_MODE`

子 Agent/团队成员的默认权限模式（如 `bypassPermissions`、`acceptEdits`、`default`、`plan`）。优先级：Agent 工具 `mode` 参数 > CLI `--subagent-permission-mode` > 此环境变量 > Settings `permissions.subagentPermissionMode` > 映射表默认值

## **Gateway 和远程访问**

**环境变量**

**说明**

`CODEBUDDY_GATEWAY_AUTH`

Gateway 认证模式（`password` 或 `none`）

`CODEBUDDY_GATEWAY_PASSWORD`

Gateway 访问密码

`CODEBUDDY_GATEWAY_FORCE_TUNNEL`

设置为 `1` 强制使用 tunnel 模式

`CODEBUDDY_DISABLE_REQUEST_VALIDATION`

设置为 `1` 关闭 Gateway 自定义请求头校验（`X-CodeBuddy-Request`）。详见 **[HTTP API 安全](https://www.codebuddy.ai/docs/zh/cli/http-api#%E5%AE%89%E5%85%A8)**

`CODEBUDDY_CODE_CORS_ORIGINS`

额外的 CORS 允许来源（逗号分隔）。支持精确 origin、`*.domain` 子域通配和 `*` 全开。如 `https://*.example.com,https://specific.com`

`SERVER__HOST`

`--serve` 模式监听地址（默认：`127.0.0.1`）

`SERVER__PORT`

`--serve` 模式监听端口

## **企业微信集成**

**环境变量**

**说明**

`CODEBUDDY_GATEWAY_WECHAT_KF_TOKEN`

企业微信客服 Token

`CODEBUDDY_GATEWAY_WECHAT_KF_ENCODING_AES_KEY`

企业微信客服加密密钥

`CODEBUDDY_GATEWAY_WECHAT_KF_CORP_ID`

企业微信客服企业 ID

`CODEBUDDY_GATEWAY_WECHAT_KF_CORP_SECRET`

企业微信客服企业密钥

`CODEBUDDY_GATEWAY_WECHAT_KF_ACCOUNT_NAME`

企业微信客服账户名

`CODEBUDDY_GATEWAY_WECOM_TOKEN`

企业微信 Token

`CODEBUDDY_GATEWAY_WECOM_ENCODING_AES_KEY`

企业微信加密密钥

`CODEBUDDY_GATEWAY_WECOM_CORP_ID`

企业微信企业 ID

`CODEBUDDY_GATEWAY_WECOM_CORP_SECRET`

企业微信企业密钥

`CODEBUDDY_GATEWAY_WECOM_AGENT_ID`

企业微信应用 ID

## **调试和诊断**

**环境变量**

**说明**

`CODEBUDDY_DEBUG`

设置为 `1` 启用调试模式

`CODEBUDDY_DEBUG_REQUEST`

设置为 `1` 启用请求调试

`CODEBUDDY_STARTUP_PROFILE`

设置为 `1` 启用启动性能分析

## **E2E 测试 (Record/Replay)**

用于 E2E 测试的模型响应录制/回放功能。录制模式下捕获真实模型响应，回放模式下使用录制文件替代真实 API 调用。

**环境变量**

**说明**

`CODEBUDDY_RECORD_DIR`

录制模式：指定录制文件保存目录。设置后模型响应会保存到该目录的 `recording.jsonl` 文件

`CODEBUDDY_REPLAY_DIR`

回放模式：指定录制文件读取目录。设置后从 `recording.jsonl` 回放响应，无需真实 API

`CODEBUDDY_REPLAY_SPEED`

回放速度倍率。`0` 表示无延迟立即返回，`1` 表示按原始时间间隔回放（默认：`1`）

`CODEBUDDY_REPLAY_STRICT`

设置为 `1` 启用严格模式：录制耗尽时抛出错误而非透传到真实 API

> **注意**：`CODEBUDDY_RECORD_DIR` 和 `CODEBUDDY_REPLAY_DIR` 互斥，不能同时设置。

## **其他**

**环境变量**

**说明**

`SLASH_COMMAND_TOOL_CHAR_BUDGET`

斜杠命令工具元数据的最大字符数（默认：15000）

`CODEBUDDY_CODE_API_KEY_HELPER_TTL_MS`

刷新凭证的间隔（毫秒）（默认：300000）

## **使用示例**

### **基础认证配置**

```
# 设置 API 密钥
export CODEBUDDY_API_KEY="your-api-key"

# 设置代理服务器
export HTTP_PROXY="http://proxy.example.com:8080"
export HTTPS_PROXY="https://proxy.example.com:8080"

# 设置代理绕过列表
export NO_PROXY="localhost,127.0.0.1,.internal.example.com"

# 设置自定义请求头（多个 header 使用 \n 分隔）
export CODEBUDDY_CUSTOM_HEADERS="X-Custom-Header: value1\nX-Another-Header: value2"

# 启动 CodeBuddy
codebuddy
```

### **使用第三方模型服务**

```
# 使用自定义 API 端点
export CODEBUDDY_API_KEY="your-api-key"
export CODEBUDDY_BASE_URL="https://api.example.com/v1"
codebuddy --model your-model-name
```

#### **对接 DeepSeek 示例**

对接任意兼容 Anthropic 协议的第三方模型服务（如 DeepSeek），只需配置 Base URL、API Key 和模型变量，无需额外修改 `models.json`：

```
# 端点与密钥
export CODEBUDDY_BASE_URL="https://api.deepseek.com"
export CODEBUDDY_API_KEY="<your-deepseek-api-key>"

# 主 Agent 默认模型
export CODEBUDDY_MODEL="deepseek-v4-pro"

# 复杂推理使用的大模型
export CODEBUDDY_BIG_SLOW_MODEL="deepseek-v4-pro"

# 后台/轻量任务使用的小模型
export CODEBUDDY_SMALL_FAST_MODEL="deepseek-v4-flash"

# 子代理使用的模型（不设置则继承主 Agent）
export CODEBUDDY_CODE_SUBAGENT_MODEL="deepseek-v4-flash"

# 启动时可通过 --model 显式指定主模型
codebuddy --model deepseek-v4-pro
```

> **提示**：以上变量也可以写入 `settings.json` 的 `env` 字段，每个会话自动应用，便于团队统一配置。

### **中国版配置**

```
# 设置中国版环境标识
export CODEBUDDY_INTERNET_ENVIRONMENT=internal

# 设置 API 密钥
export CODEBUDDY_API_KEY="your-api-key"

# 启动 CodeBuddy
codebuddy
```

### **启用高级功能**

```
# 启用自动内存
export CODEBUDDY_DISABLE_AUTO_MEMORY="0"

# 启用扩展思考
export MAX_THINKING_TOKENS="10000"

# 非交互模式运行
codebuddy -p -y "你的查询"
```

### **调试和性能分析**

```
# 启用调试模式
export CODEBUDDY_DEBUG="1"

# 启用启动性能分析
export CODEBUDDY_STARTUP_PROFILE="1"

# 启动 CodeBuddy
codebuddy
```

## **在 settings.json 中配置**

环境变量也可以在 `settings.json` 的 `env` 字段中设置：

```
{
  "env": {
    "CODEBUDDY_API_KEY": "your-api-key",
    "HTTPS_PROXY": "https://proxy.example.com:8080",
    "MAX_THINKING_TOKENS": "10000",
    "CODEBUDDY_DISABLE_AUTO_MEMORY": "0"
  }
}
```

## **工具输出外部化机制**

当工具执行产生的输出超过阈值时，CodeBuddy Code 会自动将完整输出保存到磁盘，给模型只发送截断后的内容和文件路径指针，模型可按需读取完整内容。

### **数据流**

```
Shell 输出
  ├─→ OutputSpiller（完整输出流式写入磁盘）
  └─→ TruncateBuffer（内存中保留 head + tail，约 30KB）
                          ↓
                  检测到截断 → 生成 placeholder（~2KB preview + 文件路径）
                          ↓
                  模型收到 placeholder，可通过 Read 工具按需读取完整文件
```

### **各阶段数据大小（以 1.3MB 输出为例）**

**阶段**

**内容**

**大小**

磁盘文件（OutputSpiller）

完整原始输出

1,355,099 bytes

内存缓冲（TruncateBuffer）

head 6KB + tail 24KB

\~30KB

发给模型（placeholder）

文件路径 + preview

\~2KB

### **存储目录**

工具输出文件存储在项目数据目录中：

```
~/.codebuddy/projects/{projectDir}/
  └── {sessionId}/
      ├── tool-results/                          ← 主 session 的工具结果
      │   ├── {callId}.txt
      │   └── ...
      └── subagents/                             ← 子代理数据
          ├── agent-{agentId}.jsonl              ← 子代理对话历史
          └── agent-{agentId}/
              └── tool-results/                  ← 子代理的工具结果
                  ├── {callId}.txt
                  └── ...
```

### **相关环境变量**

**环境变量**

**影响范围**

**默认值**

`BASH_MAX_OUTPUT_LENGTH`

Bash 工具内存保留量，超出则截断并触发磁盘外部化

30000

`CODEBUDDY_TOOL_RESULT_THRESHOLD_KB`

非 bash 工具（如 MCP）在 session 层的外部化阈值

50

## **另见**

- **[Settings](https://www.codebuddy.ai/docs/zh/cli/settings)** - 在 `settings.json` 中配置环境变量和其他设置
- **[CLI Reference](https://www.codebuddy.ai/docs/zh/cli/cli-reference)** - 命令行参数完整列表
- **[MCP Setup](https://www.codebuddy.ai/docs/zh/cli/mcp)** - MCP 服务器配置
- **[子代理](https://www.codebuddy.ai/docs/zh/cli/sub-agents)** - 子代理存储目录说明

**最后更新: 2026/4/27 08:00**

<br />

# **ACP 协议集成**

> ACP (Agent Client Protocol) 是 Zed 编辑器推出的一种通用智能体协议，使智能体的核心功能（服务端）和用户界面（客户端）解耦，允许用户自由选择不同的智能体服务端和客户端进行搭配使用。

CodeBuddy Code 原生支持 ACP 协议，可以作为智能体服务端与支持 ACP 的编辑器无缝集成。

## **快速开始**

### **启动 ACP 模式**

使用 `--acp` 参数启动 CodeBuddy Code 的 ACP 服务器：

```
codebuddy --acp
```

## **Zed 编辑器集成**

### **配置步骤**

打开 Zed 配置文件（`~/.config/zed/settings.json`），添加以下配置：

```
{
  "agent_servers": {
    "CodeBuddy Code": {
      "command": "codebuddy",
      "args": ["--acp"],
      "env": {}
    }
  }
}
```

随后即可在 Zed 侧边栏创建 CodeBuddy Code Thread，开始使用。

### **配置说明**

- **command**：指定 CodeBuddy Code 的命令路径（确保 `codebuddy` 在 PATH 中可用）
- **args**：使用 `["--acp"]` 启用 ACP 协议模式
- **env**：可选的环境变量配置，例如：
  ```
  {
    "env": {
      "CODEBUDDY_API_KEY": "your-api-key",
      "CODEBUDDY_INTERNET_ENVIRONMENT": "internal"
    }
  }
  ```
  > **注意**：使用 `CODEBUDDY_API_KEY` 时，必须根据版本正确配置 `CODEBUDDY_INTERNET_ENVIRONMENT`：
  >
  > - 海外版：不设置（默认）
  > - 中国版：`internal`
  > - iOA 版：`ioa`
  >
  > 详见 **[身份和访问管理文档](https://www.codebuddy.ai/docs/zh/cli/iam#%E4%B8%AA%E4%BA%BA%E7%94%A8%E6%88%B7%E8%8E%B7%E5%8F%96-api-key)**。

## **ACP 协议特性**

### **认证信息扩展**

CodeBuddy Code 在 `authenticate` 响应的 `_meta` 字段中返回用户信息：

```
{
  "_meta": {
    "codebuddy.ai/userinfo": {
      "userId": "用户 ID",
      "userName": "用户名",
      "userNickname": "用户昵称"
    }
  }
}
```

客户端可以利用这些信息提供更好的用户体验，例如显示当前登录用户、个性化界面等。

### **工具代理机制**

ACP 协议支持客户端代理部分工具操作，提升性能和安全性：

- **文件操作代理**：基于客户端的 `fs.readTextFile` 和 `fs.writeTextFile` 能力
- **终端操作代理**：基于客户端的 `terminal` 能力

当客户端声明支持这些能力时，CodeBuddy Code 会自动将相关工具调用代理给客户端执行。

### **命令列表推送**

CodeBuddy Code 会在创建新会话时自动向客户端推送可用的 Slash 命令列表（`available_commands_update`），让客户端能够：

- 提供命令自动补全功能
- 显示命令提示和帮助信息
- 动态更新可用命令

命令列表会自动过滤掉本地命令（如 `/clear`、`/exit`）和客户端专属命令（如 `/theme`、`/config`），只推送适用于 ACP 模式的命令。

### **Agent Teams 协议扩展**

CodeBuddy Code 通过 `session_info_update` 的 `_meta` 字段扩展 ACP 协议，支持 Agent Teams 多智能体协作的实时状态推送。

#### **Team 状态事件**

通过 `_meta['codebuddy.ai/teamUpdate']` 推送以下事件类型：

**成员状态变化** (`member_status_change`)：

```
{
  "sessionUpdate": "session_info_update",
  "_meta": {
    "codebuddy.ai/teamUpdate": {
      "type": "member_status_change",
      "teamName": "my-team",
      "isAutoTeam": false,
      "members": [
        {
          "name": "ux-designer",
          "color": "blue",
          "description": "用户体验设计分析",
          "status": "running",
          "taskId": "agent-abc123",
          "sessionId": "session-xyz",
          "tokenUsage": { "inputTokens": 1000, "outputTokens": 500, "lastContextWindow": 42000 },
          "toolCallCount": 5
        }
      ]
    }
  }
}
```

**Team 创建** (`team_created`) / **删除** (`team_deleted`)：

```
{
  "sessionUpdate": "session_info_update",
  "_meta": {
    "codebuddy.ai/teamUpdate": {
      "type": "team_created",
      "teamName": "my-team"
    }
  }
}
```

#### **成员流式消息**

成员的实时消息（文本、工具调用）通过标准 ACP 事件推送，附加 `_meta['codebuddy.ai/memberEvent']` 标记来标识消息来源：

```
{
  "sessionUpdate": "agent_message_chunk",
  "content": { "type": "text", "text": "正在分析架构方案..." },
  "_meta": {
    "codebuddy.ai/memberEvent": "tech-architect"
  }
}
```

客户端收到带 `memberEvent` 标记的事件后，应将其路由到对应成员的对话时间线，而非主对话区。

#### **页面刷新恢复**

页面刷新后，`loadSession` 的 `replayHistory` 完成后会自动推送当前 Team 状态（`member_status_change` 事件），客户端无需单独请求。`AcpTeamBridge` 在订阅成员 session 时会自动重放其完整历史，因此成员的对话数据也通过 ACP SSE 完整恢复，无需额外 HTTP API。

## **其他编辑器支持**

ACP 是开放协议，理论上任何支持 ACP 的编辑器都可以集成 CodeBuddy Code。配置方式与 Zed 类似：

```
{
  "agent_servers": {
    "CodeBuddy": {
      "command": "codebuddy",
      "args": ["--acp"]
    }
  }
}
```

## **故障排除**

### **连接失败**

**问题**: Zed 无法连接到 CodeBuddy

**解决方法**:

1. 确认 `codebuddy` 命令可用：
   ```
   which codebuddy
   ```
2. 测试 ACP 模式启动：
   ```
   codebuddy --acp
   ```
3. 检查配置文件 JSON 格式是否正确

### **工具调用失败**

**问题**：文件操作或命令执行报错

**解决方法**:

1. 检查工作目录权限
2. 查看 CodeBuddy 日志

## **相关链接**

- **[CLI 参考手册](https://www.codebuddy.ai/docs/zh/cli/cli-reference)** - 查看所有命令行参数
- **[IDE 集成说明](https://www.codebuddy.ai/docs/zh/cli/ide-integrations)** - 更多编辑器集成方式
- **[ACP 协议规范](https://github.com/agentclientprotocol/agent-client-protocol)** - 协议详细文档

***

*通过 ACP 协议，让 CodeBuddy Code 融入您喜爱的编辑器 🚀*

<br />

# **CodeBuddy Code 开发容器**

CodeBuddy Code 开发容器提供了一个预配置、安全的开发环境，适合需要一致性和隔离性工作空间的团队使用。它与 Visual Studio Code 的 Dev Containers 扩展以及类似工具无缝集成。

***

## **核心特性**

1. **生产就绪的 Node.js 环境**：基于 Node.js 20，包含必要的开发依赖
2. **安全设计**：自定义防火墙，限制网络访问到必要的服务
3. **开发者友好工具**：包含 git、ZSH 和生产力增强工具、fzf 等
4. **VS Code 无缝集成**：预配置扩展和优化设置
5. **会话持久化**：命令历史和配置在容器重启后保留
6. **跨平台支持**：兼容 macOS、Windows 和 Linux 开发环境

***

## **四步快速开始**

1. **安装 VS Code 和 Remote - Containers 扩展**
2. **参考下方配置详解在工作区创建 `.devcontainer` 目录及相关文件**
3. **在 VS Code 中打开仓库**
4. **当提示时，点击 "Reopen in Container"**（或使用命令面板：Cmd+Shift+P → "Remote-Containers: Reopen in Container"）

***

## **配置详解**

**目录结构**

```
your-project/
├── .devcontainer/
│   ├── devcontainer.json
│   ├── Dockerfile
│   └── init-firewall.sh
└── ...
```

开发容器设置由**三个主要组件**构成：

***

### **🚀 安装 CodeBuddy Code（推荐方案）**

#### **方案 1：使用 Dev Containers Feature（推荐）**

**Dev Containers Feature** 是在 Dev Container 中安装 CodeBuddy Code 的**推荐方式**。它提供自动版本管理、与其他特性配置一致、简化维护等优势。

**优势：**

- ✅ **自动版本管理** - 轻松升级到最新版本或固定特定版本
- ✅ **配置一致性** - 与其他 Dev Containers 特性采用相同配置方式
- ✅ **简化维护** - 无需在 Dockerfile 中管理复杂的安装逻辑
- ✅ **团队共享** - 易于在团队中统一配置

**配置方法：**

在 `devcontainer.json` 中的 `features` 字段添加配置：

```
{
    "name": "CodeBuddy Code Sandbox",
    "features": {
        "ghcr.io/devcontainers-contrib/features/codebuddy-code:1": {
            "version": "latest"
        }
    },
    // ... 其他配置
}
```

**固定特定版本：**

```
{
    "features": {
        "ghcr.io/devcontainers-contrib/features/codebuddy-code:1": {
            "version": "2.16.0"
        }
    }
}
```

**使用默认最新版本：**

```
{
    "features": {
        "ghcr.io/devcontainers-contrib/features/codebuddy-code:1": {}
    }
}
```

#### **方案 2：在 Dockerfile 中手动安装**

如果需要更多控制或特殊场景，可在 Dockerfile 中手动安装。详见下方 **Dockerfile** 部分。

***

### **1. devcontainer.json**

控制容器设置、管理扩展、配置卷挂载。

**包含 Dev Containers Feature 的完整示例：**

```
{
    "name": "CodeBuddy Code Sandbox",
    "build": {
        "dockerfile": "Dockerfile",
        "args": {
            "TZ": "${localEnv:TZ:America/Los_Angeles}",
            "GIT_DELTA_VERSION": "0.18.2",
            "ZSH_IN_DOCKER_VERSION": "1.2.0"
        }
    },
    "features": {
        "ghcr.io/devcontainers-contrib/features/codebuddy-code:1": {
            "version": "latest"
        }
    },
    "runArgs": [
        "--cap-add=NET_ADMIN",
        "--cap-add=NET_RAW"
    ],
    "customizations": {
        "vscode": {
            "extensions": [
                "dbaeumer.vscode-eslint",
                "esbenp.prettier-vscode",
                "eamodio.gitlens"
            ],
            "settings": {
                "editor.formatOnSave": true,
                "editor.defaultFormatter": "esbenp.prettier-vscode",
                "editor.codeActionsOnSave": {
                    "source.fixAll.eslint": "explicit"
                },
                "terminal.integrated.defaultProfile.linux": "zsh",
                "terminal.integrated.profiles.linux": {
                    "bash": {
                        "path": "bash",
                        "icon": "terminal-bash"
                    },
                    "zsh": {
                        "path": "zsh"
                    }
                }
            }
        }
    },
    "remoteUser": "node",
    "mounts": [
        "source=codebuddy-code-bashhistory-${devcontainerId},target=/commandhistory,type=volume",
        "source=codebuddy-code-config-${devcontainerId},target=/home/node/.codebuddy,type=volume"
    ],
    "containerEnv": {
        "NODE_OPTIONS": "--max-old-space-size=4096",
        "CODEBUDDY_CONFIG_DIR": "/home/node/.codebuddy",
        "POWERLEVEL9K_DISABLE_GITSTATUS": "true"
    },
    "workspaceMount": "source=${localWorkspaceFolder},target=/workspace,type=bind,consistency=delegated",
    "workspaceFolder": "/workspace",
    "postStartCommand": "sudo /usr/local/bin/init-firewall.sh",
    "waitFor": "postStartCommand"
}
```

**关键配置说明：**

- `features` - 声明使用 CodeBuddy Code Dev Containers Feature，支持版本管理
- `version: "latest"` - 使用最新版本（可替换为具体版本号如 "2.16.0"）
- 注意：使用 Feature 方式时，Dockerfile 中无需 `CODEBUDDY_CODE_VERSION` 参数

### **2. Dockerfile**

定义容器镜像、指定安装的工具。

#### **使用 Feature 时的简化 Dockerfile**

如使用上方 Dev Containers Feature 安装 CodeBuddy Code，Dockerfile 可以更简洁（无需 `CODEBUDDY_CODE_VERSION` 参数和手动安装命令）：

> 如果你选择使用 Dev Containers Feature 方式（推荐），可以参考下方"简化版 Dockerfile"。如需在 Dockerfile 中手动安装，可参考"完整版 Dockerfile"。

**简化版 Dockerfile（推荐配合 Feature 使用）：**

```
FROM node:20

ARG TZ
ENV TZ="$TZ"

# Install basic development tools and iptables/ipset
RUN apt-get update && apt-get install -y --no-install-recommends \
  less \
  git \
  procps \
  sudo \
  fzf \
  zsh \
  man-db \
  unzip \
  gnupg2 \
  gh \
  iptables \
  ipset \
  iproute2 \
  dnsutils \
  aggregate \
  jq \
  nano \
  vim \
  && apt-get clean && rm -rf /var/lib/apt/lists/*

# Ensure default node user has access to /usr/local/share
RUN mkdir -p /usr/local/share/npm-global && \
  chown -R node:node /usr/local/share

ARG USERNAME=node

# Persist bash history.
RUN SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  && mkdir /commandhistory \
  && touch /commandhistory/.bash_history \
  && chown -R $USERNAME /commandhistory

# Set `DEVCONTAINER` environment variable to help with orientation
ENV DEVCONTAINER=true

# Create workspace and config directories and set permissions
RUN mkdir -p /workspace /home/node/.codebuddy && \
  chown -R node:node /workspace /home/node/.codebuddy

WORKDIR /workspace

ARG GIT_DELTA_VERSION=0.18.2
RUN ARCH=$(dpkg --print-architecture) && \
  wget "https://github.com/dandavison/delta/releases/download/${GIT_DELTA_VERSION}/git-delta_${GIT_DELTA_VERSION}_${ARCH}.deb" && \
  sudo dpkg -i "git-delta_${GIT_DELTA_VERSION}_${ARCH}.deb" && \
  rm "git-delta_${GIT_DELTA_VERSION}_${ARCH}.deb"

# Set up non-root user
USER node

# Install global packages
ENV NPM_CONFIG_PREFIX=/usr/local/share/npm-global
ENV PATH=$PATH:/usr/local/share/npm-global/bin

# Set the default shell to zsh rather than sh
ENV SHELL=/bin/zsh

# Set the default editor and visual
ENV EDITOR=nano
ENV VISUAL=nano

# Default powerline10k theme
ARG ZSH_IN_DOCKER_VERSION=1.2.0
RUN sh -c "$(wget -O- https://github.com/deluan/zsh-in-docker/releases/download/v${ZSH_IN_DOCKER_VERSION}/zsh-in-docker.sh)" -- \
  -p git \
  -p fzf \
  -a "source /usr/share/doc/fzf/examples/key-bindings.zsh" \
  -a "source /usr/share/doc/fzf/examples/completion.zsh" \
  -a "export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  -x

# CodeBuddy Code will be installed via Dev Containers Feature

# Copy and set up firewall script
COPY init-firewall.sh /usr/local/bin/
USER root
RUN chmod +x /usr/local/bin/init-firewall.sh && \
  echo "node ALL=(root) NOPASSWD: /usr/local/bin/init-firewall.sh" > /etc/sudoers.d/node-firewall && \
  chmod 0440 /etc/sudoers.d/node-firewall
USER node
```

#### **手动安装方式的完整 Dockerfile**

如需完全控制安装过程，可在 Dockerfile 中手动安装 CodeBuddy Code（此方案不使用 Dev Containers Feature）：

```
FROM node:20

ARG TZ
ENV TZ="$TZ"

ARG CODEBUDDY_CODE_VERSION=latest

# Install basic development tools and iptables/ipset
RUN apt-get update && apt-get install -y --no-install-recommends \
  less \
  git \
  procps \
  sudo \
  fzf \
  zsh \
  man-db \
  unzip \
  gnupg2 \
  gh \
  iptables \
  ipset \
  iproute2 \
  dnsutils \
  aggregate \
  jq \
  nano \
  vim \
  && apt-get clean && rm -rf /var/lib/apt/lists/*

# Ensure default node user has access to /usr/local/share
RUN mkdir -p /usr/local/share/npm-global && \
  chown -R node:node /usr/local/share

ARG USERNAME=node

# Persist bash history.
RUN SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  && mkdir /commandhistory \
  && touch /commandhistory/.bash_history \
  && chown -R $USERNAME /commandhistory

# Set `DEVCONTAINER` environment variable to help with orientation
ENV DEVCONTAINER=true

# Create workspace and config directories and set permissions
RUN mkdir -p /workspace /home/node/.codebuddy && \
  chown -R node:node /workspace /home/node/.codebuddy

WORKDIR /workspace

ARG GIT_DELTA_VERSION=0.18.2
RUN ARCH=$(dpkg --print-architecture) && \
  wget "https://github.com/dandavison/delta/releases/download/${GIT_DELTA_VERSION}/git-delta_${GIT_DELTA_VERSION}_${ARCH}.deb" && \
  sudo dpkg -i "git-delta_${GIT_DELTA_VERSION}_${ARCH}.deb" && \
  rm "git-delta_${GIT_DELTA_VERSION}_${ARCH}.deb"

# Set up non-root user
USER node

# Install global packages
ENV NPM_CONFIG_PREFIX=/usr/local/share/npm-global
ENV PATH=$PATH:/usr/local/share/npm-global/bin

# Set the default shell to zsh rather than sh
ENV SHELL=/bin/zsh

# Set the default editor and visual
ENV EDITOR=nano
ENV VISUAL=nano

# Default powerline10k theme
ARG ZSH_IN_DOCKER_VERSION=1.2.0
RUN sh -c "$(wget -O- https://github.com/deluan/zsh-in-docker/releases/download/v${ZSH_IN_DOCKER_VERSION}/zsh-in-docker.sh)" -- \
  -p git \
  -p fzf \
  -a "source /usr/share/doc/fzf/examples/key-bindings.zsh" \
  -a "source /usr/share/doc/fzf/examples/completion.zsh" \
  -a "export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  -x

# Install CodeBuddy Code (manual installation method - only if not using Dev Containers Feature)
RUN npm install -g @tencent-ai/codebuddy-code@${CODEBUDDY_CODE_VERSION}

# Copy and set up firewall script
COPY init-firewall.sh /usr/local/bin/
USER root
RUN chmod +x /usr/local/bin/init-firewall.sh && \
  echo "node ALL=(root) NOPASSWD: /usr/local/bin/init-firewall.sh" > /etc/sudoers.d/node-firewall && \
  chmod 0440 /etc/sudoers.d/node-firewall
USER node
```

> **注意：** 若使用此方式，devcontainer.json 中的 `build.args` 需包含 `"CODEBUDDY_CODE_VERSION"` 参数,且 `features` 字段中不应包含 CodeBuddy Code Feature。

### **3. init-firewall.sh**

建立网络安全规则。

```
#!/bin/bash
set -euo pipefail  # Exit on error, undefined vars, and pipeline failures
IFS=$'\n\t'       # Stricter word splitting

# 1. Extract Docker DNS info BEFORE any flushing
DOCKER_DNS_RULES=$(iptables-save -t nat | grep "127\.0\.0\.11" || true)

# Flush existing rules and delete existing ipsets
iptables -F
iptables -X
iptables -t nat -F
iptables -t nat -X
iptables -t mangle -F
iptables -t mangle -X
ipset destroy allowed-domains 2>/dev/null || true

# 2. Selectively restore ONLY internal Docker DNS resolution
if [ -n "$DOCKER_DNS_RULES" ]; then
    echo "Restoring Docker DNS rules..."
    iptables -t nat -N DOCKER_OUTPUT 2>/dev/null || true
    iptables -t nat -N DOCKER_POSTROUTING 2>/dev/null || true
    echo "$DOCKER_DNS_RULES" | xargs -L 1 iptables -t nat
else
    echo "No Docker DNS rules to restore"
fi

# First allow DNS and localhost before any restrictions
# Allow outbound DNS
iptables -A OUTPUT -p udp --dport 53 -j ACCEPT
# Allow inbound DNS responses
iptables -A INPUT -p udp --sport 53 -j ACCEPT
# Allow outbound SSH
iptables -A OUTPUT -p tcp --dport 22 -j ACCEPT
# Allow inbound SSH responses
iptables -A INPUT -p tcp --sport 22 -m state --state ESTABLISHED -j ACCEPT
# Allow localhost
iptables -A INPUT -i lo -j ACCEPT
iptables -A OUTPUT -o lo -j ACCEPT

# Create ipset with CIDR support
ipset create allowed-domains hash:net

# Fetch GitHub meta information and aggregate + add their IP ranges
echo "Fetching GitHub IP ranges..."
gh_ranges=$(curl -s https://api.github.com/meta)
if [ -z "$gh_ranges" ]; then
    echo "ERROR: Failed to fetch GitHub IP ranges"
    exit 1
fi

if ! echo "$gh_ranges" | jq -e '.web and .api and .git' >/dev/null; then
    echo "ERROR: GitHub API response missing required fields"
    exit 1
fi

echo "Processing GitHub IPs..."
while read -r cidr; do
    if [[ ! "$cidr" =~ ^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}/[0-9]{1,2}$ ]]; then
        echo "ERROR: Invalid CIDR range from GitHub meta: $cidr"
        exit 1
    fi
    echo "Adding GitHub range $cidr"
    ipset add allowed-domains "$cidr"
done < <(echo "$gh_ranges" | jq -r '(.web + .api + .git)[]' | aggregate -q)

# Resolve and add other allowed domains
for domain in \
    "registry.npmjs.org" \
    "copilot.tencent.com" \
    "sentry.io" \
    "marketplace.visualstudio.com" \
    "vscode.blob.core.windows.net" \
    "update.code.visualstudio.com"; do
    echo "Resolving $domain..."
    ips=$(dig +noall +answer A "$domain" | awk '$4 == "A" {print $5}')
    if [ -z "$ips" ]; then
        echo "ERROR: Failed to resolve $domain"
        exit 1
    fi
    
    while read -r ip; do
        if [[ ! "$ip" =~ ^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$ ]]; then
            echo "ERROR: Invalid IP from DNS for $domain: $ip"
            exit 1
        fi
        echo "Adding $ip for $domain"
        ipset add allowed-domains "$ip"
    done < <(echo "$ips")
done

# Get host IP from default route
HOST_IP=$(ip route | grep default | cut -d" " -f3)
if [ -z "$HOST_IP" ]; then
    echo "ERROR: Failed to detect host IP"
    exit 1
fi

HOST_NETWORK=$(echo "$HOST_IP" | sed "s/\.[0-9]*$/.0\/24/")
echo "Host network detected as: $HOST_NETWORK"

# Set up remaining iptables rules
iptables -A INPUT -s "$HOST_NETWORK" -j ACCEPT
iptables -A OUTPUT -d "$HOST_NETWORK" -j ACCEPT

# Set default policies to DROP first
iptables -P INPUT DROP
iptables -P FORWARD DROP
iptables -P OUTPUT DROP

# First allow established connections for already approved traffic
iptables -A INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT
iptables -A OUTPUT -m state --state ESTABLISHED,RELATED -j ACCEPT

# Then allow only specific outbound traffic to allowed domains
iptables -A OUTPUT -m set --match-set allowed-domains dst -j ACCEPT

# Explicitly REJECT all other outbound traffic for immediate feedback
iptables -A OUTPUT -j REJECT --reject-with icmp-admin-prohibited

echo "Firewall configuration complete"
echo "Verifying firewall rules..."
if curl --connect-timeout 5 https://example.com >/dev/null 2>&1; then
    echo "ERROR: Firewall verification failed - was able to reach https://example.com"
    exit 1
else
    echo "Firewall verification passed - unable to reach https://example.com as expected"
fi

# Verify GitHub API access
if ! curl --connect-timeout 5 https://api.github.com/zen >/dev/null 2>&1; then
    echo "ERROR: Firewall verification failed - unable to reach https://api.github.com"
    exit 1
else
    echo "Firewall verification passed - able to reach https://api.github.com as expected"
fi
```

***

## **安全特性**

容器实现了**多层安全防护**：

- **精确访问控制**：限制出站连接到白名单域名（npm registry、GitHub、CodeBuddy API 等）
- **允许的出站连接**：防火墙允许出站 DNS 和 SSH 连接
- **默认拒绝策略**：阻止所有其他外部网络访问
- **启动验证**：容器初始化时验证防火墙规则
- **隔离**：创建与主系统分离的安全开发环境

### **重要安全提示**

> 虽然开发容器提供了实质性的保护，但没有系统能够完全免疫所有攻击。当使用 `-y` （或 `--dangerously-skip-permissions`) 执行时，开发容器无法阻止恶意项目窃取容器内可访问的任何内容，包括 CodeBuddy Code 凭证。**我们建议仅在处理可信仓库时使用开发容器。** 始终保持良好的安全实践并监控 CodeBuddy 的活动。

### **无人值守操作**

容器增强的安全措施（隔离和防火墙规则）允许您运行 `codebuddy -y` （或 `codebuddy --dangerously-skip-permissions`) 来绕过权限提示，实现无人值守操作。

***

## **自定义选项**

开发容器配置设计灵活，可根据需求调整：

- 根据工作流添加或删除 VS Code 扩展
- 为不同硬件环境调整资源分配
- 调整网络访问权限
- 自定义 shell 配置和开发工具

***

## **使用场景**

### **1. 安全的客户项目开发**

使用开发容器隔离不同客户的项目，确保代码和凭证永不混合。

### **2. 团队快速入职**

新团队成员可在几分钟内获得完全配置的开发环境，所有必要的工具和设置都已预装。

### **3. 一致的 CI/CD 环境**

在 CI/CD 流水线中镜像您的开发容器配置，确保开发和生产环境匹配。

***

## **相关资源**

- **[VS Code 开发容器文档](https://code.visualstudio.com/docs/devcontainers/containers)**
- **[Bash 沙箱](https://www.codebuddy.ai/docs/zh/cli/bash-sandboxing)**
- **[设置配置](https://www.codebuddy.ai/docs/zh/cli/settings)**
- **[GitLab CI/CD 集成](https://www.codebuddy.ai/docs/zh/cli/gitlab-ci-cd)**
- 📖 **[Official Dev Containers Docs](https://containers.dev/)**
- 🔗 **[devcontainers-contrib/features](https://github.com/devcontainers-contrib/features)**

**最后更新: 2025/12/29 23:34**

<br />

# **MCP (Model Context Protocol) 使用文档**

## **概述**

MCP (Model Context Protocol) 是一个开放标准，允许 CodeBuddy 与外部工具和数据源进行集成。通过 MCP，您可以扩展 CodeBuddy 的功能，连接到各种外部服务、数据库、API 等。

## **核心概念**

### **MCP 服务器**

MCP 服务器是提供工具、资源和提示的独立进程，CodeBuddy 通过不同的传输协议与这些服务器通信。

### **MCP Prompts 集成**

MCP 服务器可以提供 Prompts（提示模板）,这些 Prompts 会自动转换为 CodeBuddy 的斜杠命令。当 MCP 服务器连接后：

- 服务器提供的 Prompts 会自动注册为斜杠命令
- 命令名称格式为： `/服务器名:prompt名称`
- 支持动态参数，通过交互式界面收集用户输入
- 命令执行时会调用 MCP 服务器的 `prompts/get` 接口获取完整内容
- 支持实时监听配置变更，自动更新可用命令列表

### **传输类型**

- **STDIO**：通过标准输入输出与本地进程通信
- **SSE**：通过 Server-Sent Events 与远程服务通信
- **HTTP**：通过 HTTP 流式传输与远程服务通信

### **配置作用域**

- **user**：全局用户配置，应用于所有项目
- **project**：项目级配置，应用于特定项目
- **local**：本地配置，仅应用于当前会话或工作区

对于同名服务（即在多个作用域有同名配置），生效的优先级为：`local > project > user`

### **安全审批机制**

项目作用域的 MCP 服务器在首次连接时需要用户审批，以确保安全性。系统会显示服务器详细信息，用户可以选择批准或拒绝连接。

#### **非交互模式（-p/--print）下的审批**

在非交互模式（如使用 `-p/--print` 参数)下,由于无法通过 UI 进行审批,需要通过 `--settings` 参数预先配置允许的 MCP 服务器：

```
# 方式 1：允许所有项目 MCP 服务器
codebuddy --settings '{"enableAllProjectMcpServers": true}' -p "your prompt"

# 方式 2：允许特定的 MCP 服务器
codebuddy --settings '{"enabledMcpjsonServers": ["server-name-1", "server-name-2"]}' -p "your prompt"
```

### **工具权限管理**

MCP 工具支持完整的权限管理系统，可以精确控制哪些工具可以被使用：

#### **权限规则类型**

权限系统支持三种规则类型（按优先级排序）：

1. **拒绝规则 （deny)** - 阻止使用指定工具（最高优先级）
2. **询问规则 （ask)** - 使用工具前需要用户确认（覆盖允许规则）
3. **允许规则 （allow)** - 允许工具使用而无需手动批准

#### **MCP 权限规则格式**

**重要**：MCP 权限不支持通配符 （\*)

##### **服务器级权限**

```
mcp__服务器名
```

- 匹配指定服务器提供的任何工具
- 服务器名是在 CodeBuddy 中配置的名称

##### **工具级权限**

```
mcp__服务器名__工具名
```

- 匹配指定服务器的特定工具

#### **配置示例**

##### **批准服务器的所有工具**

```
{
  "permissions": {
    "allow": [
      "mcp__github"
    ]
  }
}
```

##### **仅批准特定工具**

```
{
  "permissions": {
    "allow": [
      "mcp__github__get_issue",
      "mcp__github__list_issues"
    ]
  }
}
```

##### **拒绝特定工具**

```
{
  "permissions": {
    "deny": [
      "mcp__dangerous_server__delete_file"
    ]
  }
}
```

## **配置文件**

### **配置文件位置**

配置文件使用优先级机制，系统会按优先级顺序查找第一个存在的文件进行读取。写入时，如果文件已存在，会写入第一个存在的文件；如果都不存在，会创建最高优先级的文件。

#### **USER 作用域**

优先级顺序（从高到低）：

1. `~/.codebuddy/.mcp.json`（推荐）
2. `~/.codebuddy/mcp.json`（已废弃）
3. `~/.codebuddy.json`（旧版配置文件）

**读取规则**：系统会按上述顺序查找第一个存在的文件并读取其内容。

**写入规则**：

- 如果上述文件中存在任意一个，写入第一个存在的文件
- 如果都不存在，创建 `~/.codebuddy/.mcp.json`（最高优先级）

#### **PROJECT 作用域**

优先级顺序（从高到低）：

1. `<项目根目录>/.mcp.json`（推荐）
2. `<项目根目录>/mcp.json`（已废弃）

**读取规则**：系统会按上述顺序查找第一个存在的文件并读取其内容。

**写入规则**：

- 如果上述文件中存在任意一个，写入第一个存在的文件
- 如果都不存在，创建 `<项目根目录>/.mcp.json`（最高优先级）

#### **LOCAL 作用域**

local 作用域的配置实际上保存在 user 作用域的配置文件中，通过 `projects` 字段来区分不同项目的 local 配置。

文件路径：`~/.codebuddy.json#/projects/<workspace_path>`

`#/projects/<workspace_path>` 使用的是 JSON Pointer 语法，用于指向 JSON 文档中的特定位置。关于 JSON Pointer 的详细说明，请参考：**<https://datatracker.ietf.org/doc/html/rfc6901>**

**注意**：

- 系统不会合并同一个作用域的多个配置文件内容，只会使用第一个存在的文件

示例见下方配置文件格式说明。

### **配置文件格式**

MCP 配置文件支持 **JSONC (JSON with Comments)** 格式，允许在配置中添加注释，提升可读性和可维护性。

#### **JSONC 支持的特性**

- **单行注释**：使用 `//` 添加行内或行尾注释
- **多行注释**：使用 `/* */` 添加块注释
- **尾随逗号**：数组和对象的最后一个元素后可以添加逗号

#### **基础配置格式**

```
{
  // MCP 服务器配置
  "mcpServers": {
    "server-name": {
      "type": "stdio|sse|http",
      "command": "命令路径",
      "args": ["参数1", "参数2"],
      "env": {
        "ENV_VAR": "value"
      },
      "url": "http://example.com/mcp",
      "headers": {
        "Authorization": "Bearer token"
      },
      "description": "服务器描述"
    }
  },
  // projects 字段仅在 user 作用域的文件里有效，用于识别 local 作用域的配置
  "projects": {
    "/path/to/project": {
      "mcpServers": {
        "local-server": {
          "type": "stdio",
          "command": "./local-tool"
        }
      }
    }
  }
}
```

#### **带注释的完整示例**

```
{
  // MCP Server Configuration for CodeBuddy
  // 这个文件配置了项目使用的 MCP 服务器
  
  "mcpServers": {
    /*
     * Filesystem Server
     * 提供文件系统访问能力
     * 文档: https://github.com/modelcontextprotocol/servers
     */
    "filesystem": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@modelcontextprotocol/server-filesystem",
        "/path/to/workspace",  // 工作目录路径
      ],
      "env": {
        "DEBUG": "true",  // 启用调试模式
      },
    },
    
    // HTTP API 服务器示例
    "api-server": {
      "type": "http",
      "url": "http://localhost:3000/mcp",  // 本地开发服务器
      "headers": {
        "Authorization": "Bearer your-token",
      },
    },
  },
  
  // 已禁用的服务器列表（供参考）
  "disabledMcpServers": [
    "deprecated-server",
  ],
}
```

**注意**：

- 标准 JSON 格式文件仍然完全兼容
- 解析错误时会提供清晰的错误提示

````

**注意**：`type` 字段是可选的。如果未指定，系统会根据配置内容自动推断：
- 包含 `command` 字段时,推断为 `stdio` 类型
- 包含 `url` 字段时,推断为 `http` 类型

建议显式指定 `type` 字段以确保配置的准确性。

### 环境变量扩展

MCP 配置支持环境变量扩展，允许您在配置中引用系统环境变量。这对于在团队间共享配置、管理敏感信息（如 API 密钥、令牌）以及支持环境特定配置（开发、测试、生产）非常有用。

#### 支持的语法

- **`${VAR_NAME}`** - 展开为环境变量 VAR_NAME 的值
- **`${VAR_NAME:-default_value}`** - 如果 VAR_NAME 未设置，使用默认值

#### 变量命名规则

- 变量名必须以大写字母或下划线 `[A-Z_]` 开头
- 后续字符只能是大写字母、数字或下划线 `[A-Z0-9_]*`
- 小写字母、混合大小写以及以数字开头的变量不会被展开

#### 支持的配置字段

环境变量可以在以下配置字段中展开：

**STDIO 类型配置**：
- `command` - 可执行文件路径或命令
- `args` - 命令行参数列表中的每个参数
- `env` - 环境变量值（键不会被展开）

**SSE/HTTP/Remote 类型配置**：
- `url` - 服务端点 URL
- `headers` - HTTP 请求头值（键不会被展开）

#### 错误处理

**环境变量未设置的行为**：
- 如果环境变量未设置且**有默认值**，使用默认值
- 如果环境变量未设置且**无默认值**，保留原始占位符（`${VAR}`），并在诊断中报告 WARNING 消息

这意味着配置不会因缺失的环境变量而失败，而是保留占位符并发出警告。

#### 示例配置

**示例 1：STDIO 类型服务器，使用环境变量**
```json
{
  "mcpServers": {
    "python-tools": {
      "type": "stdio",
      "command": "${PYTHON_PATH:-python}",
      "args": [
        "-m",
        "my_mcp_server",
        "--config",
        "${CONFIG_DIR:-/etc/config}"
      ],
      "env": {
        "PYTHONPATH": "${PYTHON_LIB_PATH}",
        "DEBUG": "${DEBUG_MODE:-false}",
        "API_KEY": "${API_KEY}"
      }
    }
  }
}
````

**示例 2：HTTP 类型服务器，使用环境变量和默认值**

```
{
  "mcpServers": {
    "api-server": {
      "type": "http",
      "url": "${API_BASE_URL:-https://api.example.com}/mcp",
      "headers": {
        "Authorization": "Bearer ${API_TOKEN}",
        "X-API-Version": "${API_VERSION:-v1}",
        "User-Agent": "CodeBuddy/${CODEBUDDY_VERSION:-1.0}"
      }
    }
  }
}
```

#### **常见用例**

1. **团队共享配置**
   ```
   # 在 .mcp.json 中使用环境变量
   # 每个团队成员在本地设置环境变量
   export API_TOKEN="their-personal-token"
   export LOCAL_TOOL_PATH="/home/user/tools"
   ```
2. **环境特定配置**
   ```
   # 开发环境
   export API_BASE_URL="http://localhost:3000"

   # 生产环境
   export API_BASE_URL="https://api.production.com"
   ```
3. **管理敏感信息**
   ```
   {
     "headers": {
       "Authorization": "Bearer ${MY_API_KEY}"
     }
   }
   ```
   将 API 密钥存储在环境变量中，不要直接写在配置文件中。

#### **诊断和调试**

当配置中的环境变量展开时，您可以通过以下方式了解展开结果：

1. 如果某个环境变量未设置且无默认值，系统会发出 WARNING 诊断
2. 可以通过 `/mcp` 命令查看 MCP 服务器配置和诊断信息
3. 诊断消息会列出所有缺失的环境变量

**示例诊断消息**：

```
Missing environment variables: API_TOKEN, DATABASE_URL
```

### **配置结构详解**

根据不同的传输类型，MCP 服务器配置具有不同的结构：

#### **STDIO 类型配置**

通过标准输入输出与本地进程通信。

**字段**

**类型**

**必填**

**说明**

`type`

string

是

固定值 `"stdio"`

`command`

string

是

可执行文件路径或命令

`args`

Array\<string>

否

命令行参数列表

`env`

Object

否

环境变量键值对

`defer_loading`

boolean

否

是否延迟加载工具（默认 false）

`tools`

Object

否

工具级别配置，可覆盖服务器级别设置

**示例**：

```
{
  "type": "stdio",
  "command": "python",
  "args": ["-m", "my_mcp_server"],
  "env": {
    "PYTHONPATH": "/path/to/tools",
    "DEBUG": "true"
  }
}
```

#### **SSE 类型配置**

通过 Server-Sent Events 与远程服务通信。

**字段**

**类型**

**必填**

**说明**

`type`

string

是

固定值 `"sse"`

`url`

string

是

SSE 服务端点 URL

`headers`

Object

否

HTTP 请求头键值对

`defer_loading`

boolean

否

是否延迟加载工具（默认 false）

`tools`

Object

否

工具级别配置，可覆盖服务器级别设置

**示例**：

```
{
  "type": "sse",
  "url": "https://api.example.com/mcp/sse",
  "headers": {
    "Authorization": "Bearer your-api-token",
    "X-API-Version": "v1"
  }
}
```

#### **HTTP 类型配置**

通过 HTTP 流式传输与远程服务通信。

**字段**

**类型**

**必填**

**说明**

`type`

string

是

固定值 `"http"`

`url`

string

是

HTTP 服务端点 URL

`headers`

Object

否

HTTP 请求头键值对

`defer_loading`

boolean

否

是否延迟加载工具（默认 false）

`tools`

Object

否

工具级别配置，可覆盖服务器级别设置

**示例**：

```
{
  "type": "http",
  "url": "https://mcp.example.com/api/v1",
  "headers": {
    "Authorization": "Bearer secret-token",
    "Content-Type": "application/json"
  }
}
```

### **延迟加载 (defer\_loading)**

当 MCP 服务器提供大量工具时，可以使用 `defer_loading` 配置来延迟加载工具，减少上下文消耗并提高模型工具选择的准确性。

#### **工作原理**

- 设置 `defer_loading: true` 的工具不会在初始请求时加载到模型上下文
- 模型可以通过 `ToolSearch` 工具搜索这些延迟加载的工具
- 搜索到的工具会被激活，并在后续请求中可用
- 激活状态在当前会话中保持

#### **服务器级别配置**

将服务器的所有工具设为延迟加载：

```
{
  "mcpServers": {
    "my-server": {
      "type": "stdio",
      "command": "my-mcp-server",
      "defer_loading": true
    }
  }
}
```

#### **工具级别配置**

可以为单个工具覆盖服务器级别的设置：

```
{
  "mcpServers": {
    "my-server": {
      "type": "stdio",
      "command": "my-mcp-server",
      "defer_loading": true,
      "tools": {
        "frequently_used_tool": {
          "defer_loading": false
        }
      }
    }
  }
}
```

#### **继承规则**

**服务器 defer\_loading**

**工具 defer\_loading**

**最终结果**

true

未设置

true（继承）

true

false

false（覆盖）

false/未设置

未设置

false

false/未设置

true

true（覆盖）

#### **使用场景**

- **工具数量多**：当 MCP 服务器提供超过 30 个工具时
- **减少成本**：减少每次请求的 token 消耗
- **提高准确性**：让模型在更少的工具中做出更准确的选择

## **命令行使用**

### **添加 MCP 服务器**

#### **STDIO 服务器**

```
# 添加本地可执行文件
codebuddy mcp add --scope user my-tool -- /path/to/tool arg1 arg2

# 添加 Python 脚本
codebuddy mcp add --scope project python-tool -- python /path/to/script.py
```

#### **SSE 服务器**

```
# 添加 SSE 服务器
codebuddy mcp add --scope user --transport sse sse-server https://example.com/mcp/sse
```

#### **HTTP 服务器**

```
# 添加 HTTP 流式服务器
codebuddy mcp add --scope project --transport http http-server https://example.com/mcp/http
```

### **使用 JSON 配置添加服务器**

```
# 添加 STDIO 类型服务器
codebuddy mcp add-json --scope user my-server '{"type":"stdio","command":"/usr/local/bin/tool","args":["--verbose"]}'

# 添加 HTTP 类型服务器
codebuddy mcp add-json --scope user http-server '{"type":"http","url":"https://example.com/mcp","headers":{"Authorization":"Bearer token"}}'

# 添加 SSE 类型服务器
codebuddy mcp add-json --scope project sse-server '{"type":"sse","url":"https://api.example.com/mcp/sse","headers":{"X-API-Key":"your-api-key"}}'

# 添加带环境变量的 STDIO 服务器
codebuddy mcp add-json --scope user python-tool '{"type":"stdio","command":"python","args":["-m","my_mcp_server"],"env":{"PYTHONPATH":"/path/to/tools"}}'
```

### **管理 MCP 服务器**

#### **列出所有服务器**

```
# 列出所有作用域的服务器
codebuddy mcp list
```

#### **查看服务器详情**

```
# 查看特定服务器信息
codebuddy mcp get my-server
```

#### **移除服务器**

```
# 移除特定服务器
codebuddy mcp remove my-server

# 移除特定作用域的服务器
codebuddy mcp remove my-server --scope user
```

## **最佳实践**

### **1. 作用域选择**

- 使用 **user** 作用域存储个人工具和全局服务
- 使用 **project** 作用域存储项目特定的工具
- 使用 **local** 作用域存储临时或实验性工具

### **2. 安全考虑**

- 避免在配置文件中存储敏感信息
- **使用环境变量传递认证信息**：利用 MCP 的环境变量扩展功能（`${API_TOKEN}` 或 `${API_TOKEN:-default}`）来管理 API 密钥、令牌等敏感数据
- 定期审查和更新服务器配置
- 项目作用域的 MCP 服务器需要用户审批后才能连接，确保安全性
- OAuth 授权 URL 会在打开前进行安全验证，仅支持 http/https 协议
- 将包含环境变量引用的配置文件提交到版本控制系统，但在 `.gitignore` 中排除实际的环境变量文件

### **3. 性能优化**

- 合理配置服务器超时时间
- 避免同时运行过多的 STDIO 服务器
- 使用缓存机制减少重复连接

### **4. 错误处理**

- 监控服务器连接状态
- 实现重连机制
- 记录和分析错误日志

## **故障排除**

### **常见问题**

#### **服务器连接失败**

1. 检查命令路径是否正确
2. 验证参数和环境变量
3. 确认网络连接（对于远程服务器）
4. 查看服务器日志输出

#### **工具不可用**

1. 确认服务器已成功连接
2. 检查工具权限设置
3. 验证工具兼容性

#### **配置不生效**

1. 检查配置文件语法
2. 确认作用域优先级
3. 重启 CodeBuddy 应用

## **示例配置**

### **Python 工具服务器**

```
{
  "mcpServers": {
    "python-tools": {
      "type": "stdio",
      "command": "python",
      "args": ["-m", "my_mcp_server"],
      "env": {
        "PYTHONPATH": "/path/to/tools"
      },
      "description": "Python 工具集合"
    }
  }
}
```

### **远程 API 服务器**

```
{
  "mcpServers": {
    "api-server": {
      "type": "sse",
      "url": "https://api.example.com/mcp/sse",
      "headers": {
        "Authorization": "Bearer your-token",
        "X-API-Version": "v1"
      },
      "description": "远程 API 服务"
    }
  }
}
```

### **Node.js 本地服务器**

```
{
  "mcpServers": {
    "node-server": {
      "type": "stdio", 
      "command": "node",
      "args": ["./mcp-server.js"],
      "env": {
        "NODE_ENV": "production"
      },
      "description": "Node.js MCP 服务器"
    }
  }
}
```

## **扩展开发**

### **创建自定义 MCP 服务器**

1. **选择实现语言**: Python、Node.js、Go 等
2. **实现 MCP 协议**：使用官方 SDK 或自行实现
3. **定义工具接口**：描述工具功能和参数
4. **处理请求**：接收和处理来自 CodeBuddy 的请求
5. **返回结果**：按 MCP 格式返回执行结果

### **SDK 和库**

- **Python**: `FastMCP`
- **TypeScript/JavaScript**: `@modelcontextprotocol/sdk`
- **其他语言**：参考官方文档实现

## **配置示例**

### **TAPD**

```
codebuddy mcp add --scope user --transport http --header "X-Tapd-Access-Token: TAPD_ACCESS_TOKEN" -- tapd_mcp_http https://mcp-oa.tapd.woa.com/mcp
```

### **Chrome Devtools**

```
codebuddy mcp add --scope user chrome-devtools -- npx -y chrome-devtools-mcp@latest
```

### **iWiki**

```
codebuddy mcp add --scope user iwiki -- npx -y mcp-remote@latest https://prod.mcp.it.woa.com/app_iwiki_mcp/mcp3
```

## **相关链接**

- **[MCP 官方文档](https://modelcontextprotocol.io/)**
- **[MCP GitHub 仓库](https://github.com/modelcontextprotocol)**
- **[CodeBuddy 官方文档](https://cnb.cool/codebuddy/codebuddy-code/-/blob/main/docs)**

**最后更新: 2026/3/7 22:06**

<br />

# **CodeBuddy Code Skills （技能系统）**

Skills 是 CodeBuddy Code 的扩展能力系统，允许您创建专业的领域知识和工作流模板，让 AI 助手能够更专业地处理特定类型的任务。

## **什么是 Skills**

Skills 类似于为 AI 助手提供的"专业培训"。通过 Skill，您可以：

- **封装专业知识**：将特定领域的最佳实践和操作流程封装成可复用的技能
- **提供工作流模板**：定义标准化的任务处理流程，提高工作效率
- **扩展 AI 能力**：让 AI 助手能够处理更专业、更复杂的任务
- **团队协作共享**：项目级 Skills 可以在团队成员间共享专业知识

## **Skills vs Slash Commands**

**特性**

**Skills**

**Slash Commands**

**触发方式**

AI 模型自动识别并调用

用户手动输入命令

**使用场景**

专业领域任务处理

快捷操作和工作流

**权限控制**

支持工具白名单限制

无特殊权限控制

**工作目录**

支持自定义基础目录

使用当前工作目录

**可见性**

对用户透明，AI 自动决策

用户主动发起

**简单来说**：

- **Slash Commands** 是用户主动调用的快捷方式
- **Skills** 是 AI 根据任务需求自动选择的专业能力

## **创建 Skills**

### **目录结构**

Skills 通过在特定目录中创建 `SKILL.md` 文件来定义：

1. **项目级 Skills**：`.codebuddy/skills/`（项目根目录下）
2. **用户级 Skills**：`~/.codebuddy/skills/`（用户主目录下）

每个 Skill 一个独立的目录，包含 `SKILL.md` 文件：

```
.codebuddy/skills/
├── pdf/
│   └── SKILL.md
├── data-analysis/
│   └── SKILL.md
└── code-review/
    └── SKILL.md
```

### **SKILL.md 格式**

Skill 文件使用 Markdown 格式，支持 YAML Frontmatter 定义元数据：

```
---
name: pdf
description: PDF 文档处理专家
allowed-tools: Read, Write, Bash, WebFetch
---

你是一个 PDF 文档处理专家，擅长：
- 解析和提取 PDF 内容
- 转换 PDF 为其他格式
- 生成 PDF 报告

当用户需要处理 PDF 相关任务时，请使用以下工作流：
1. 首先检查 PDF 文件是否存在
2. 使用适当的工具提取内容
3. 根据需求进行处理
4. 生成结果报告

可用工具：
- pdftotext：提取文本内容
- pdfinfo：获取 PDF 信息
```

### **Frontmatter 字段**

**字段**

**必填**

**说明**

**示例**

`name`

否

Skill 名称，未指定时使用目录名

`pdf`

`description`

否

Skill 描述，帮助 AI 理解何时使用

`PDF 文档处理专家 （project)`

`allowed-tools`

否

允许使用的工具白名单，逗号分隔

`Read, Write, Bash`

`disable-model-invocation`

否

设置为 `true` 时，Skill 不会出现在 Skill 工具中，只能通过 `/skill-name` 手动触发

`true`

`user-invocable`

否

设置为 `false` 时，Skill 从 `/` 菜单中隐藏，仅供 AI 内部调用或其他 Skill 引用，默认 `true`

`false`

`context`

否

设置为 `fork` 时，Skill 在独立的 subagent 上下文中执行

`fork`

`agent`

否

指定 subagent 类型，仅在 `context: fork` 时有效

`Explore`

## **执行 Shell 命令**

与**[斜杠命令](https://www.codebuddy.ai/docs/zh/cli/slash-commands)**一样，Skills 也支持在 SKILL.md 中使用 `!`command\`\` 语法内联执行 Shell 命令。当 Skill 被触发时（无论是 AI 自动调用还是用户通过 `/skill-name` 手动触发），这些命令会被执行，输出结果会替换到 Skill 内容中，供 AI 后续分析。

**示例**：

```
---
description: 项目状态分析
---

### 当前工作目录

!`echo "CWD=$(pwd)"`

### Git 状态

!`git status --short`

### 最近提交

!`git log --oneline -5`

请基于以上信息分析项目当前状态。
```

### **支持的特性**

- **`$ARGUMENTS` 参数替换**：在 Shell 命令执行前，`$ARGUMENTS` 会被替换为用户传入的参数
- **`@file` 文件引用**：Shell 命令执行后，`@file` 引用会被处理并注入文件内容
- **错误隔离**：单个命令执行失败不会影响其他命令，失败的命令会被替换为空字符串

> 💡 **处理管道**：`$ARGUMENTS` 替换 → `!`command\`\` 执行 → `@file` 引用处理，与斜杠命令的处理顺序一致。

## **Context Fork**

`context: fork` 使 Skill 在隔离的子代理上下文中运行，不访问对话历史。

```
---
name: deep-research
description: 深入研究某个主题
context: fork
agent: Explore
---

研究 $ARGUMENTS：
1. 使用 Glob 和 Grep 查找相关文件
2. 读取并分析代码
3. 总结发现并附加具体文件引用
```

### **可用 Agent 类型**

**类型**

**说明**

`general-purpose`

通用（默认）

`Explore`

只读工具，优化代码库探索

`Plan`

规划和分析

自定义

`.codebuddy/agents/` 中定义的 agent

### **隐藏 Skill（user-invocable）**

`user-invocable: false` 使 Skill 从 `/` 菜单中隐藏，适用于：

- 背景知识类 Skill（如项目规范、编码标准）
- 仅供其他 Skill 或 AI 内部引用的辅助 Skill

```
---
name: project-guidelines
description: 项目编码规范和最佳实践
user-invocable: false
---

# 项目编码规范

本项目遵循以下编码标准：
- 使用 TypeScript 严格模式
- 函数命名使用 camelCase
- 组件命名使用 PascalCase
...
```

这类 Skill 会被加载到 AI 的上下文中，但用户无法通过 `/` 菜单直接调用。

### **执行流程**

1. 创建新的隔离上下文
2. 子代理接收 Skill 内容作为提示
3. `agent` 字段决定执行环境
4. 结果返回主对话

> **注意**：`context: fork` 只适用于包含明确任务的 Skill。仅有指导方针没有具体任务时，不会产生有意义的输出。

## **使用示例**

### **示例 1：PDF 处理 Skill**

**文件**：`.codebuddy/skills/pdf/SKILL.md`

```
---
name: pdf
description: PDF 文档处理和转换专家
allowed-tools: Read, Write, Bash, WebFetch
---

# PDF 处理专家

你是一个专业的 PDF 文档处理专家。

## 核心能力
- 提取 PDF 文本内容
- 转换 PDF 为 Markdown、HTML 等格式
- 合并和拆分 PDF 文件
- 提取 PDF 元数据和书签

## 工作流程
1. 检查 PDF 文件是否存在并可访问
2. 使用 pdftotext 或 pdfinfo 获取基本信息
3. 根据任务类型选择合适的处理工具
4. 验证输出结果的完整性

## 可用工具
- pdftotext：提取纯文本
- pdfinfo：获取文档信息
- pdftk：合并拆分操作
```

**使用**：当用户询问 "帮我提取这个 PDF 的内容" 时，AI 会自动识别需要 PDF 处理能力并调用该 Skill。

### **示例 2：数据分析 Skill**

**文件**：`~/.codebuddy/skills/data-analysis/SKILL.md`

```
---
name: data-analysis
description：数据分析和可视化专家
allowed-tools: Read, Write, Bash, WebFetch, NotebookEdit
---

# 数据分析专家

你是一个专业的数据分析师，擅长使用 Python 和相关工具进行数据分析。

## 核心能力
- 数据清洗和预处理
- 统计分析和建模
- 数据可视化
- 生成分析报告

## 分析流程
1. 理解数据结构和质量
2. 清洗和预处理数据
3. 执行统计分析
4. 创建可视化图表
5. 生成分析结论

## 工具库
- pandas：数据处理
- numpy：数值计算
- matplotlib/seaborn：可视化
- scikit-learn：机器学习

## 最佳实践
- 始终先探索数据质量
- 使用 Jupyter Notebook 进行交互式分析
- 保存中间结果避免重复计算
```

### **示例 3：代码审查 Skill**

**文件**：`.codebuddy/skills/code-review/SKILL.md`

```
---
name: code-review
description：代码审查和质量检查专家
allowed-tools: Read, Grep, Bash, Edit
---

# 代码审查专家

你是一个经验丰富的代码审查者，遵循业界最佳实践。

## 审查重点
1. **代码质量**
   - 命名规范
   - 代码复杂度
   - 重复代码

2. **安全性**
   - SQL 注入风险
   - XSS 漏洞
   - 认证授权问题

3. **性能**
   - 算法效率
   - 资源使用
   - 缓存策略

4. **可维护性**
   - 代码注释
   - 模块化设计
   - 测试覆盖

## 审查流程
1. 理解代码变更的目的
2. 检查代码风格和规范
3. 分析潜在的 Bug 和性能问题
4. 验证安全性
5. 提供建设性的改进建议

## 输出格式
- ✅ 优点：列出做得好的地方
- ⚠️ 问题：指出需要改进的地方
- 💡 建议：提供具体的改进方案
```

## **AI 如何选择 Skills**

AI 根据以下因素决定是否调用 Skill：

1. **任务匹配度**：任务描述与 Skill description 的相关性
2. **工具需求**：任务所需工具是否在 allowed-tools 范围内
3. **上下文相关性**：当前对话上下文是否适合使用该 Skill
4. **Skill 来源**：项目级 Skills 优先于用户级 Skills

## **权限控制**

### **allowed-tools 白名单**

通过 `allowed-tools` 字段限制 Skill 可以使用的工具：

```
allowed-tools: Read, Write, Bash(git:*), Grep
```

支持的工具模式匹配：

- `Bash(git:*)` - 只允许 git 相关命令
- `Edit(src/**/*.ts)` - 只允许编辑特定路径文件

### **工作目录限制**

每个 Skill 都有自己的 `baseDirectory`（SKILL.md 所在目录），可以在 Skill 指令中引用：

```
当处理文件时，优先在 {baseDirectory} 目录下查找相关资源。
```

## **最佳实践**

### **1. 清晰的 Skill 描述**

```
# ❌ 不好
description：处理文件

# ✅ 好
description: PDF 文档解析和转换专家，支持文本提取和格式转换 （project)
```

### **2. 详细的指令内容**

提供详细的：

- 核心能力说明
- 标准工作流程
- 可用工具列表
- 常见场景处理方法
- 输出格式要求

### **3. 合理的工具权限**

只授予必需的工具权限：

```
# ❌ 权限过大
allowed-tools: Bash

# ✅ 精确控制
allowed-tools: Read, Write, Bash(git:status,git:diff), Grep
```

### **4. 组织 Skill 目录**

按功能领域组织 Skills：

```
.codebuddy/skills/
├── document/
│   ├── pdf/SKILL.md
│   └── markdown/SKILL.md
├── data/
│   ├── analysis/SKILL.md
│   └── visualization/SKILL.md
└── code/
    ├── review/SKILL.md
    └── refactor/SKILL.md
```

## **调试 Skills**

### **查看已加载的 Skills**

使用 `/skills` 命令查看当前已加载的所有 Skills：

```
/skills
```

Skills 面板会显示：

- **User skills**：用户级 Skills（`~/.codebuddy/skills/`）
- **Project skills**：项目级 Skills（`.codebuddy/skills/`）
- **Plugin skills**：插件提供的 Skills

每个 Skill 会显示名称和预估的 token 数量。

### **常见问题**

**Q: Skill 没有被触发？**

- 检查 description 是否清晰描述了 Skill 的功能
- 确认任务描述与 Skill 能力匹配
- 验证 allowed-tools 是否包含所需工具

**Q: Skill 权限不足？**

- 检查 allowed-tools 配置
- 确认工具名称拼写正确
- 使用模式匹配精确控制权限

**Q：项目级和用户级 Skill 冲突？**

- 项目级 Skills 优先级更高
- 使用不同的 name 避免冲突

## **与其他功能的配合**

### **Skills + Memory**

Skills 可以访问 Memory 系统存储的信息：

```
在执行数据分析时，参考 Memory 中保存的数据模式和业务规则。
```

### **Skills + Slash Commands**

Slash Commands 可以引用 Skills：

```
<!-- .codebuddy/commands/analyze-data.md -->
请使用 data-analysis skill 分析文件：$1
```

### **Skills + MCP**

Skills 可以调用 MCP 提供的外部工具（如果在 allowed-tools 中）。

## **下一步**

- **[斜杠命令](https://www.codebuddy.ai/docs/zh/cli/slash-commands)** - 了解用户主动命令
- **[设置配置](https://www.codebuddy.ai/docs/zh/cli/settings)** - 配置工具权限
- **[MCP 集成](https://www.codebuddy.ai/docs/zh/cli/mcp)** - 扩展外部工具能力

***

*Skills - 让 AI 成为领域专家*

**最后更新: 2026/2/24 22:11**

<br />

# **Web UI**

CodeBuddy Code 提供内置的 Web UI，在浏览器中提供完整的交互界面。当您以 serve 模式启动或开启远程控制时，Web UI 自动可用。

## **概述**

Web UI 提供与终端界面相同的核心能力，并针对浏览器进行了可视化布局优化：

- **对话**：发送消息、查看对话、实时监控工具执行
- **终端**：内嵌终端，支持分屏布局（最多 4 个面板）
- **Workers**：管理 CLI Worker 进程和 Daemon 守护进程
- **日志**：独立日志查看器，支持多种日志类型和关键词搜索
- **远程控制**：连接微信和企业微信渠道
- **监控**：系统资源指标和各 Worker 进程级内存/运行时间指标
- **任务**：浏览任务模版并创建定时任务
- **插件**：管理插件安装和插件市场
- **设置**：配置主题、语言、模型和权限模式
- **文档**：浏览 CLI 文档，支持全文搜索
- **API 文档**：查看交互式 Swagger UI，方便 HTTP API 探索

## **访问 Web UI**

### **方式一：Serve 模式**

使用 `--serve` 参数启动 CodeBuddy Code：

```
codebuddy --serve --port 7890
```

然后在浏览器中打开：

```
http://127.0.0.1:7890
```

### **方式二：远程控制**

在已有的 CodeBuddy Code 会话中启动 Gateway：

```
/gateway
```

终端会显示二维码和 URL。用手机扫码或在浏览器中打开 URL。详见**[远程控制](https://www.codebuddy.ai/docs/zh/cli/remote-control)**文档。

## **认证方式**

Web UI 支持两种认证模式：

**模式**

**设置**

**说明**

免认证（默认）

`CODEBUDDY_GATEWAY_AUTH=none`

无需密码

密码认证

`CODEBUDDY_GATEWAY_AUTH=password`

启动时终端显示密码

认证方式（以下任一方式均可）：

- **URL 参数**：`?password=xxx` — 通过 URL 自动登录
- **登录页面**：输入终端显示的密码
- **Bearer Token**：`Authorization: Bearer <password>`，用于 API 访问

在 `~/.codebuddy/settings.json` 中配置：

```
{
  "gateway.auth": "none"
}
```

## **功能详解**

### **对话视图**

默认视图，用于与 Agent 交互。核心功能：

- **富文本消息渲染**：Markdown、语法高亮代码块、表格、图片
- **工具执行展示**：内联查看工具调用、参数和结果
- **权限管理**：在浏览器中直接批准或拒绝工具权限
- **问答面板**：回答 Agent 的多选问题
- **任务进度**：实时监控后台任务和 Team 进度
- **会话管理**：新建对话、浏览历史、切换会话

### **终端视图**

基于 xterm.js 的内嵌终端：

- **分屏布局**：支持水平和垂直分屏，最多 4 个面板
- **独立会话**：每个面板有独立的 PTY 会话
- **持久连接**：终端会话在页面刷新后保持
- **自适应调整**：窗口大小变化时面板自动调整

### **文档视图**

在 Web UI 中直接浏览 CLI 文档：

- **全文搜索**：基于 MiniSearch 搜索所有文档
- **多语言**：自动跟随 UI 语言设置（中文/英文）
- **目录导航**：根据文档标题自动生成，带滚动追踪
- **内部跳转**：文档链接在查看器内导航（SPA 模式）
- **API 文档**：快速链接到 `/api/docs` 的交互式 Swagger UI

### **实例管理**

管理多个 CodeBuddy Code 实例：

- **实例列表**：查看所有运行中的实例及其工作目录和状态
- **快速切换**：一键切换实例
- **手动添加**：通过 URL 添加远程实例
- **隧道支持**：通过 Cloudflare Tunnel 访问实例

### **设置**

- **主题**：浅色、深色或跟随系统（自动检测）
- **语言**：中文、英文或跟随系统（自动检测）
- **模型**：从可用选项中选择 AI 模型
- **权限模式**：默认、接受编辑、跳过权限或规划模式

## **API 文档**

HTTP 服务运行时，可在以下地址访问交互式 API 探索器：

```
http://127.0.0.1:{PORT}/api/docs
```

Swagger UI 提供以下功能：

- 浏览所有可用的 REST API 端点
- 查看请求/响应 Schema
- 直接在浏览器中测试 API 调用
- 在 `/api/openapi.json` 下载 OpenAPI 3.1 规范

完整的 API 参考请见 **[HTTP API 文档](https://www.codebuddy.ai/docs/zh/cli/http-api)**。

## **移动端支持**

Web UI 完全响应式，支持移动设备：

- **侧边栏**：小屏幕上折叠为滑出式抽屉
- **PWA 支持**：添加到主屏幕获得类原生应用体验
- **触控优化**：所有交互针对触摸操作优化
- **扫码访问**：终端扫码即可在手机上打开

## **快捷键**

**快捷键**

**操作**

`Enter`

发送消息

`Shift+Enter`

输入换行

`Escape`

停止运行中的 Agent

## **技术细节**

- **框架**：React 18 + Zustand 状态管理
- **通信**：ACP 协议，基于 HTTP/SSE（非 WebSocket）
- **样式**：Tailwind CSS + CSS 变量主题
- **终端**：xterm.js + fit addon
- **搜索**：MiniSearch 客户端全文搜索
- **Markdown**：react-markdown + remark-gfm + 语法高亮

## **相关文档**

- **[远程控制](https://www.codebuddy.ai/docs/zh/cli/remote-control)** — 通过 Gateway 和 Tunnel 启动 Web UI
- **[HTTP API](https://www.codebuddy.ai/docs/zh/cli/http-api)** — 完整的 REST API 文档
- **[ACP 协议](https://www.codebuddy.ai/docs/zh/cli/acp)** — IDE 集成的 Agent Client Protocol
- **[设置](https://www.codebuddy.ai/docs/zh/cli/settings)** — 配置选项

**最后更新: 2026/4/7 15:00**

<br />

# **Bash 沙箱**

> 了解 CodeBuddy Code 的沙箱化 Bash 工具如何提供文件系统和网络隔离，实现更安全、更自主的智能体执行。

## **概述**

CodeBuddy Code 具有原生沙箱功能，为智能体执行提供更安全的环境，同时减少持续的权限提示。沙箱不需要为每个 bash 命令请求权限，而是预先创建定义好的边界，让 CodeBuddy Code 可以在降低风险的情况下更自由地工作。

沙箱化 Bash 工具使用操作系统级原语来强制执行文件系统和网络隔离。

## **为什么沙箱很重要**

传统的基于权限的安全性需要持续的用户批准才能执行 bash 命令。虽然这提供了控制，但可能导致：

- **审批疲劳**：重复点击"批准"可能导致用户对批准的内容关注度降低
- **降低生产力**：持续的中断会减慢开发工作流程
- **有限的自主性**：在等待批准时，CodeBuddy Code 无法高效工作

沙箱通过以下方式解决这些挑战：

1. **定义清晰的边界**：明确指定 CodeBuddy Code 可以访问哪些目录和网络主机
2. **减少权限提示**：沙箱内的安全命令不需要批准
3. **保持安全性**：访问沙箱外资源的尝试会立即触发通知
4. **实现自主性**: CodeBuddy Code 可以在定义的限制内更独立地运行

**WARNING**

有效的沙箱需要\*\*同时\*\*具备文件系统和网络隔离。没有网络隔离，受损的智能体可能会窃取敏感文件（如 SSH 密钥）。没有文件系统隔离，受损的智能体可能会后门系统资源以获得网络访问。配置沙箱时，确保你的配置设置不会在这些系统中创建绕过方式非常重要。

## **工作原理**

### **文件系统隔离**

沙箱化 Bash 工具将文件系统访问限制在特定目录：

- **默认写入行为**：对当前工作目录及其子目录具有读写访问权限
- **默认读取行为**：对整个计算机具有读取访问权限，除了某些被拒绝的目录
- **阻止访问**：未经明确许可，无法修改当前工作目录之外的文件
- **可配置**：通过设置定义自定义的允许和拒绝路径

### **网络隔离**

网络访问通过在沙箱外运行的代理服务器控制：

- **域限制**：只能访问批准的域
- **用户确认**：新的域请求会触发权限提示
- **自定义代理支持**：高级用户可以对出站流量实施自定义规则
- **全面覆盖**：限制适用于命令生成的所有脚本、程序和子进程

### **操作系统级强制执行**

沙箱化 Bash 工具利用操作系统安全原语：

- **Linux**: 使用 **[bubblewrap](https://github.com/containers/bubblewrap)** 进行隔离
- **macOS**：使用 Seatbelt 进行沙箱强制执行

这些操作系统级限制确保 CodeBuddy Code 命令生成的所有子进程都继承相同的安全边界。

## **入门**

### **启用沙箱**

你可以通过运行 `/sandbox` 斜杠命令来启用沙箱：

```
> /sandbox
```

这将使用默认设置激活沙箱化 Bash 工具，允许访问当前工作目录，同时阻止访问敏感的系统位置。

### **配置沙箱**

通过 `settings.json` 文件自定义沙箱行为。完整的配置参考请参阅**[设置](https://www.codebuddy.ai/docs/zh/cli/settings#bash%E6%B2%99%E7%AE%B1%E8%AE%BE%E7%BD%AE)**。

并非所有命令都能开箱即用地与沙箱兼容。以下一些注意事项可能帮助你充分利用沙箱：

- 许多 CLI 工具需要访问某些主机。当你使用这些工具时，它们会请求访问某些主机的权限。授予权限将允许它们现在和将来访问这些主机，使它们能够在沙箱内安全执行。
- `watchman` 与在沙箱中运行不兼容。如果你正在运行 `jest`,请考虑使用 `jest --no-watchman`
- `docker` 与在沙箱中运行不兼容。考虑在 `excludedCommands` 中指定 `docker` 以强制其在沙箱外运行。

**NOTE**

CodeBuddy Code 包含一个有意设置的逃生舱机制，允许在必要时在沙箱外运行命令。当命令因沙箱限制（如网络连接问题或不兼容的工具）而失败时，CodeBuddy 会被提示分析失败，并可能使用 \`dangerouslyDisableSandbox\` 参数重试命令。使用此参数的命令会经过正常的 CodeBuddy Code 权限流程，需要用户权限才能执行。这允许 CodeBuddy Code 处理某些工具或网络操作无法在沙箱约束内运行的边缘情况。

你可以通过在**[沙箱设置](https://www.codebuddy.ai/docs/zh/cli/settings#bash%E6%B2%99%E7%AE%B1%E8%AE%BE%E7%BD%AE)**中设置 `"allowUnsandboxedCommands": false` 来禁用此逃生舱。禁用后,`dangerouslyDisableSandbox` 参数会被完全忽略,所有命令必须在沙箱中运行或明确列在 `excludedCommands` 中。

## **沙箱自动批准**

在 AcceptEdits 权限模式下，CodeBuddy Code 提供了沙箱自动批准功能，进一步减少权限提示：

### **工作原理**

当启用沙箱自动批准时，满足以下所有条件的 Bash 命令将自动批准执行，无需用户确认：

1. **当前处于 AcceptEdits 权限模式**
2. **沙箱已启用** (`sandbox.enabled = true`)
3. **自动批准配置已开启** (`sandbox.autoAllowBashIfSandboxed = true`)
4. **命令将在沙箱中运行** （不在 `excludedCommands` 列表中）
5. **未使用 dangerouslyDisableSandbox 参数**

### **配置示例**

在 `.codebuddy/settings.json` 中启用沙箱自动批准：

```
{
  "sandbox": {
    "enabled": true,
    "autoAllowBashIfSandboxed": true,
    "excludedCommands": ["git", "docker"]
  }
}
```

### **使用场景**

**自动批准** （无需用户确认）:

```
# 在 AcceptEdits 模式下，这些命令会自动批准
npm test
npm run build
ls -la
grep "pattern" file.txt
```

**需要批准** （以下情况仍需用户确认）:

```
# 1. 排除的命令 （在 excludedCommands 中）
git push

# 2. 明确禁用沙箱的命令
# （使用 dangerouslyDisableSandbox 参数）

# 3. 在非 AcceptEdits 模式下的所有命令
```

### **安全保障**

沙箱自动批准在保持安全的同时提升效率：

- ✅ **只在沙箱环境内**：只有沙箱化的命令才能自动批准
- ✅ **文件系统隔离**：命令只能访问允许的目录
- ✅ **网络隔离**：命令只能访问批准的域名
- ✅ **可配置边界**：通过 `excludedCommands` 排除高风险命令
- ✅ **安全优先**：任何配置错误或异常都会回退到需要用户批准

### **最佳实践**

1. **谨慎使用**：只在理解沙箱限制的情况下启用自动批准
2. **排除高风险命令**：将 `git push`、`rm -rf` 等危险操作加入 `excludedCommands`
3. **定期审查**：检查日志了解哪些命令被自动批准
4. **组合使用**：配合 IAM 权限规则使用，提供深度防御

## **安全优势**

### **防范提示注入**

即使攻击者成功通过提示注入操纵 CodeBuddy Code 的行为，沙箱也能确保你的系统保持安全：

**文件系统保护：**

- 无法修改关键配置文件，如 `~/.bashrc`
- 无法修改 `/bin/` 中的系统级文件
- 无法修改 CodeBuddy 配置文件（`settings.json`、`settings.local.json`），防止通过注入 hooks 实现沙箱逃逸
- 无法读取在 **[CodeBuddy 权限设置](https://www.codebuddy.ai/docs/zh/cli/iam#%E9%85%8D%E7%BD%AE%E6%9D%83%E9%99%90)**中被拒绝的文件

**网络保护：**

- 无法将数据窃取到攻击者控制的服务器
- 无法从未授权的域下载恶意脚本
- 无法对未批准的服务进行意外的 API 调用
- 无法联系任何未明确允许的域

**监控和控制：**

- 所有访问沙箱外的尝试都在操作系统级别被阻止
- 当边界被测试时，你会立即收到通知
- 你可以选择拒绝、允许一次或永久更新配置

### **配置文件保护**

沙箱默认阻止对 CodeBuddy 配置文件的写入，以防止沙箱逃逸攻击。攻击者可能通过提示注入让 Agent 写入 `settings.json`，注入恶意 `SessionStart` hooks，在用户下次启动 CodeBuddy 时以主机权限执行恶意命令。

以下配置文件在沙箱中受到写保护：

- `~/.codebuddy/settings.json` - 用户全局设置
- `~/.codebuddy/settings.local.json` - 用户本地设置
- `.codebuddy/settings.json` - 项目共享设置
- `.codebuddy/settings.local.json` - 项目本地设置

此保护同时应用于 Bash 命令和文件编辑工具（Write、Edit、MultiEdit），确保沙箱的 `denyWrite` 规则在所有工具中统一生效。

### **减少攻击面**

沙箱限制了以下潜在损害：

- **恶意依赖**：具有有害代码的 NPM 包或其他依赖项
- **受损脚本**：具有安全漏洞的构建脚本或工具
- **社会工程**：诱骗用户运行危险命令的攻击
- **提示注入**：诱骗 CodeBuddy 运行危险命令的攻击

### **透明操作**

当 CodeBuddy Code 尝试访问沙箱外的网络资源时：

1. 操作在操作系统级别被阻止
2. 你立即收到通知
3. 你可以选择：
   - 拒绝请求
   - 允许一次
   - 更新沙箱配置以永久允许

## **安全限制**

- **网络沙箱限制**：网络过滤系统通过限制进程允许连接的域来运作。它不会检查通过代理的流量，用户有责任确保他们只允许其策略中的受信任域。

**WARNING**

用户应该意识到允许广泛域名(如 \`github.com\`)可能带来的风险,这可能允许数据窃取。此外,在某些情况下,可能通过\[域前置]\(https\://en.wikipedia.org/wiki/Domain\_fronting)绕过网络过滤。

- **通过 Unix 套接字的权限提升**: `allowUnixSockets` 配置可能会无意中授予对强大系统服务的访问权限,这可能导致沙箱绕过。例如,如果使用它允许访问 `/var/run/docker.sock`,这将通过利用 docker 套接字有效地授予对主机系统的访问权限。鼓励用户仔细考虑他们允许通过沙箱的任何 unix 套接字。
- **文件系统权限提升**：过于宽泛的文件系统写入权限可能会启用权限提升攻击。允许对包含 `$PATH` 中可执行文件的目录、系统配置目录或用户 shell 配置文件(`.bashrc`、`.zshrc`)进行写入，当其他用户或系统进程访问这些文件时，可能导致在不同安全上下文中执行代码。
- **Linux 沙箱强度**: Linux 实现提供了强大的文件系统和网络隔离，但包括一个 `enableWeakerNestedSandbox` 模式，使其能够在没有特权命名空间的 Docker 环境中工作。此选项会大大削弱安全性，只应在其他方式强制执行额外隔离的情况下使用。

## **高级用法**

### **自定义代理配置**

对于需要高级网络安全的组织，你可以实现自定义代理来：

- 解密和检查 HTTPS 流量
- 应用自定义过滤规则
- 记录所有网络请求
- 与现有安全基础设施集成

```
{
  "sandbox": {
    "network": {
      "httpProxyPort": 8080,
      "socksProxyPort": 8081
    }
  }
}
```

### **与现有安全工具集成**

沙箱化 Bash 工具与以下工具协同工作：

- **IAM 策略**：与**[权限设置](https://www.codebuddy.ai/docs/zh/cli/iam)**结合使用以实现深度防御
- **开发容器**：与 devcontainers 一起使用以实现额外隔离
- **企业策略**：通过**[管理设置](https://www.codebuddy.ai/docs/zh/cli/settings#%E8%AE%BE%E7%BD%AE%E4%BC%98%E5%85%88%E7%BA%A7)**强制执行沙箱配置

## **最佳实践**

1. **从限制性开始**：从最小权限开始，根据需要扩展
2. **监控日志**：审查沙箱违规尝试以了解 CodeBuddy Code 的需求
3. **使用特定于环境的配置**：为开发和生产环境使用不同的沙箱规则
4. **与权限结合**：将沙箱与 IAM 策略一起使用以实现全面安全
5. **测试配置**：验证你的沙箱设置不会阻止合法工作流程
6. **谨慎使用自动批准**：只在充分理解沙箱限制的情况下启用 `autoAllowBashIfSandboxed`

## **开源**

沙箱运行时可作为开源 npm 包用于你自己的智能体项目。这使更广泛的 AI 智能体社区能够构建更安全、更可靠的自主系统。这也可以用于沙箱化你可能希望运行的其他程序。例如，要沙箱化 MCP 服务器，你可以运行：

```
npx @anthropic-ai/sandbox-runtime <command-to-sandbox>
```

有关实现细节和源代码,请访问 **[GitHub 仓库](https://github.com/anthropic-experimental/sandbox-runtime)**。

## **限制**

- **性能开销**：最小，但某些文件系统操作可能略慢
- **兼容性**：一些需要特定系统访问模式的工具可能需要配置调整，甚至可能需要在沙箱外运行
- **平台支持**：目前支持 Linux 和 macOS;计划支持 Windows

## **相关文档**

- **[安全](https://www.codebuddy.ai/docs/zh/cli/security)** - 全面的安全功能和最佳实践
- **[IAM](https://www.codebuddy.ai/docs/zh/cli/iam)** - 权限配置和访问控制
- **[设置](https://www.codebuddy.ai/docs/zh/cli/settings)** - 完整的配置参考

***

*让 CodeBuddy Code 在安全的沙箱中更自主地工作 🛡️*

**最后更新: 2026/3/30 17:49**

<br />

<br />

# **CLI 参考**

> CodeBuddy Code 命令行工具完整参考手册，包含所有命令和参数说明。

## **CLI 命令**

**命令**

**说明**

**示例**

`codebuddy`

启动交互式 REPL

`codebuddy`

`codebuddy "查询"`

带初始提示词启动 REPL

`codebuddy "解释这个项目"`

`codebuddy -p "查询"`

通过 SDK 查询后退出

`codebuddy -p "解释这个函数"`

`cat 文件 | codebuddy -p "查询"`

处理管道内容

`cat logs.txt | codebuddy -p "分析日志"`

`codebuddy -c`

继续最近的对话

`codebuddy -c`

`codebuddy -c -p "查询"`

通过 SDK 继续对话

`codebuddy -c -p "检查类型错误"`

`codebuddy -r "<session-id>" "查询"`

通过 ID 恢复会话

`codebuddy -r "abc123" "完成这个 MR"`

`codebuddy update`

更新到最新版本

`codebuddy update`

`codebuddy mcp`

配置 Model Context Protocol (MCP) 服务器

参见 **[CodeBuddy Code MCP 文档](https://www.codebuddy.ai/docs/zh/cli/mcp)**

`codebuddy daemon start`

启动 Daemon 守护进程

`codebuddy daemon start --port 8080`

`codebuddy daemon stop`

停止 Daemon

`codebuddy daemon stop`

`codebuddy daemon status`

查看 Daemon 状态

`codebuddy daemon status`

`codebuddy daemon restart`

重启 Daemon

`codebuddy daemon restart`

`codebuddy ps`

列出所有活跃 Worker 进程

`codebuddy ps`

`codebuddy logs <pid|name>`

查看 Worker 日志

`codebuddy logs feature-x`

`codebuddy attach <pid|name>`

附加到后台 Worker

`codebuddy attach feature-x`

`codebuddy kill <pid|name>`

终止 Worker 进程

`codebuddy kill feature-x`

## **CLI 参数**

自定义 CodeBuddy Code 行为的命令行参数：

**参数**

**说明**

**示例**

`--add-dir`

添加额外的工作目录供 CodeBuddy 访问（验证每个路径是否存在）

`codebuddy --add-dir ../apps ../lib`

`--agent`

指定当前会话使用的 agent 名称（内置或自定义 agent），优先级高于 settings.json 的 `agent` 配置

`codebuddy --agent my-reviewer`

`--agents`

通过 JSON 动态定义自定义**[子代理](https://www.codebuddy.ai/docs/zh/cli/sub-agents)**（格式见下文）

`codebuddy --agents '{"reviewer":{"description":"审查代码","prompt":"你是代码审查员"}}'`

`--allowedTools`

除了**[settings.json 文件](https://www.codebuddy.ai/docs/zh/cli/settings)**外,无需提示用户即可允许的工具列表

`"Bash(git log:*)" "Bash(git diff:*)" "Read"`

`--disallowedTools`

除了**[settings.json 文件](https://www.codebuddy.ai/docs/zh/cli/settings)**外,应禁止使用的工具列表

`"Bash(git log:*)" "Bash(git diff:*)" "Edit"`

`--tools`

限制可用的内置工具集（白名单）。空字符串 `""` 禁用所有内置工具，`"default"` 使用全部工具，或指定逗号分隔的工具名

`codebuddy --tools "Bash,Read,Edit"`

`--print`, `-p`

打印响应后退出,不进入交互模式

`codebuddy -p "查询"`

`--settings`

从 JSON 文件或 JSON 字符串加载额外的设置配置

`codebuddy --settings '{"model":"gpt-5"}' "查询"`

`--setting-sources`

指定要加载的设置源,逗号分隔（可选值: `user`, `project`, `local`）。默认: `user,project,local`

`codebuddy --setting-sources project,local "查询"`

`--system-prompt`

用自定义文本替换整个系统提示词（在交互和打印模式下都可用）

`codebuddy --system-prompt "你是 Python 专家"`

`--system-prompt-file`

从文件加载系统提示词,替换默认提示词(仅打印模式)

`codebuddy -p --system-prompt-file ./custom-prompt.txt "查询"`

`--append-system-prompt`

在默认系统提示词末尾追加自定义文本（在交互和打印模式下都可用）

`codebuddy --append-system-prompt "始终使用 TypeScript"`

`--output-format`

指定打印模式的输出格式(选项: `text`, `json`, `stream-json`)

`codebuddy -p "查询" --output-format json`

`--input-format`

指定打印模式的输入格式(选项: `text`, `stream-json`)

`codebuddy -p --output-format json --input-format stream-json`

`--json-schema`

使用 JSON Schema 验证结构化输出。示例: `'{"type":"object","properties":{"name":{"type":"string"}},"required":["name"]}'`

`codebuddy -p --output-format json --json-schema '{"type":"object","properties":{...}}' "查询"`

`--include-partial-messages`

在输出中包含部分流式事件（需要 `--print` 和 `--output-format=stream-json`)

`codebuddy -p --output-format stream-json --include-partial-messages "查询"`

`--verbose`

启用详细日志记录,显示完整的轮次输出(在打印和交互模式下都有助于调试)

`codebuddy --verbose`

`--max-turns`

限制非交互模式下的代理轮次数

`codebuddy -p --max-turns 3 "查询"`

`--model`

使用别名设置当前会话的模型，如最新模型的别名（`sonnet` 或 `opus`）或模型全名

`codebuddy --model gpt-5`

`--text-to-image-model`

设置文生图功能使用的模型 ID

`codebuddy --text-to-image-model your-image-model`

`--image-to-image-model`

设置图生图功能使用的模型 ID

`codebuddy --image-to-image-model your-edit-model`

`--permission-mode`

以指定的**[权限模式](https://www.codebuddy.ai/docs/zh/cli/iam#permission-modes)**开始

`codebuddy --permission-mode plan`

`--subagent-permission-mode`

设置 subagent/团队成员的默认权限模式，覆盖从主 session 继承的模式

`codebuddy --subagent-permission-mode bypassPermissions`

`--permission-prompt-tool`

指定在非交互模式下处理权限提示的 MCP 工具

`codebuddy -p --permission-prompt-tool mcp_auth_tool "查询"`

`--resume`

通过 ID 恢复特定会话,或在交互模式下选择

`codebuddy --resume abc123 "查询"`

`--continue`

加载当前目录中最近的对话

`codebuddy --continue`

`-y` / `--dangerously-skip-permissions`

跳过权限提示（谨慎使用）

`codebuddy -y` 或 `codebuddy --dangerously-skip-permissions`

`--ide`

启动时自动连接到 IDE（如果恰好有一个有效的 IDE 可用且打开了当前工作目录）

`codebuddy --ide`

`--sandbox`

在沙箱中运行 CodeBuddy（详见下方**[沙箱模式](https://www.codebuddy.ai/docs/zh/cli/cli-reference#%E6%B2%99%E7%AE%B1%E6%A8%A1%E5%BC%8F-beta)**)

`codebuddy --sandbox "分析项目"`

`--debug`

启用调试模式,支持可选的类别过滤

`codebuddy --debug`

`--worktree [name]`

在独立的 git worktree 中运行（详见 **[Worktree 文档](https://www.codebuddy.ai/docs/zh/cli/worktree)**）

`codebuddy --worktree` 或 `codebuddy --worktree my-feature`

`--tmux`

在 tmux 会话中运行（与 `--worktree` 配合使用）

`codebuddy --worktree --tmux`

`--plugin-dir <dirs...>`

从本地目录加载插件（用于开发/测试），可指定多个路径。详见 **[插件文档](https://www.codebuddy.ai/docs/zh/cli/plugins)**

`codebuddy --plugin-dir ./my-plugin ../other-plugin`

`--bg`

后台运行会话（detached 模式），日志输出到 `~/.codebuddy/logs/`。详见 **[Daemon 文档](https://www.codebuddy.ai/docs/zh/cli/daemon)**

`codebuddy --bg "实现登录页面"`

`--name <name>`

后台会话名称（与 `--bg` 配合使用，便于通过 `ps`/`logs`/`kill` 查找）

`codebuddy --bg --name feature-x "实现功能"`

`--serve`

启动 HTTP 服务（Web UI、REST API、ACP 协议）

`codebuddy --serve --port 8080`

> **重要提示**：在使用 `-p/--print` 参数进行非交互式执行时,涉及文件读写、命令执行、网络请求等操作通常需要添加 `-y` （或 `--dangerously-skip-permissions`)参数才能执行，否则操作会被阻止。

\`--output-format json\` 参数特别适用于脚本和自动化，允许您以编程方式解析 CodeBuddy 的响应。

### **Agents 参数格式**

`--agents` 参数接受定义一个或多个自定义子代理的 JSON 对象。每个子代理需要一个唯一的名称（作为键）和一个包含以下字段的定义对象：

**字段**

**必需**

**说明**

`description`

是

何时应调用子代理的自然语言描述

`prompt`

是

指导子代理行为的系统提示词

`tools`

否

子代理可以使用的特定工具数组（如 `["Read", "Edit", "Bash"]`)。省略则继承所有工具

`model`

否

要使用的模型别名: `sonnet`、`opus` 或 `haiku`。省略则使用默认子代理模型

示例：

```
codebuddy --agents '{
  "code-reviewer": {
    "description": "专业代码审查员。代码更改后主动使用。",
    "prompt": "你是高级代码审查员。专注于代码质量、安全性和最佳实践。",
    "tools": ["Read", "Grep", "Glob", "Bash"],
    "model": "sonnet"
  },
  "debugger": {
    "description": "错误和测试失败的调试专家。",
    "prompt": "你是专业调试人员。分析错误，识别根本原因并提供修复方案。"
  }
}'
```

有关创建和使用子代理的更多详细信息，请参见**[子代理文档](https://www.codebuddy.ai/docs/zh/cli/sub-agents)**。

### **系统提示词参数**

CodeBuddy Code 提供三个自定义系统提示词的参数，每个参数用途不同：

**参数**

**行为**

**模式**

**使用场景**

`--system-prompt`

**替换**整个默认提示词

交互 + 打印模式

完全控制 CodeBuddy 的行为和指令

`--system-prompt-file`

用文件内容**替换**

仅打印模式

从文件加载提示词以确保可重现性和版本控制

`--append-system-prompt`

**追加**到默认提示词

交互 + 打印模式

添加特定指令同时保留默认 CodeBuddy Code 行为

**何时使用：**

- **`--system-prompt`**：当您需要完全控制 CodeBuddy 的系统提示词时使用。这会移除所有默认 CodeBuddy Code 指令，给您一个空白画布。
  ```
  codebuddy --system-prompt "你是只编写带类型注解代码的 Python 专家"
  ```
- **`--system-prompt-file`**：当您想从文件加载自定义提示词时使用，适用于团队一致性或版本控制的提示词模板。
  ```
  codebuddy -p --system-prompt-file ./prompts/code-review.txt "审查这个 MR"
  ```
- **`--append-system-prompt`**：当您想添加特定指令同时保留 CodeBuddy Code 的默认功能时使用。这是大多数用例的最安全选项。
  ```
  codebuddy --append-system-prompt "始终使用 TypeScript 并包含 JSDoc 注释"
  ```

**NOTE**

\`--system-prompt\` 和 \`--system-prompt-file\` 互斥。不能同时使用这两个参数。

对于大多数用例，建议使用 \`--append-system-prompt\`,因为它在添加自定义需求的同时保留了 CodeBuddy Code 的内置功能。仅当需要完全控制系统提示词时才使用 \`--system-prompt\` 或 \`--system-prompt-file\`。

## **沙箱模式 （Beta)**

> **Beta 功能**: Sandbox 功能目前处于 Beta 阶段。
>
> **详细文档**：查看 **[Bash 沙箱](https://www.codebuddy.ai/docs/zh/cli/bash-sandboxing)** 获取沙箱隔离功能说明。

### **沙箱参数**

```
--sandbox [url]                       在沙箱中运行 CodeBuddy:
                                      - 不带参数或 "container"：使用容器 （Docker/Podman)
                                      - 提供完整的 E2B API URL：使用云端沙箱
--sandbox-upload-dir                  上传当前工作目录到沙箱 （仅 E2B)
--sandbox-new                         强制创建新沙箱 （忽略缓存的沙箱）
--sandbox-id <id>                     连接到指定的沙箱 ID 或别名
--sandbox-kill                        退出时终止沙箱 （默认： 保持运行以便复用）
--teleport <value>                    Teleport 模式： 连接到远程创建的沙箱
```

### **沙箱使用示例**

```
# 容器沙箱 （Docker/Podman，自动挂载当前目录）
codebuddy --sandbox "分析这个项目"

# E2B 云端沙箱 （自动复用）
codebuddy --sandbox https://api.e2b.dev "创建 Python web 应用"

# 强制创建新沙箱
codebuddy --sandbox --sandbox-new "从头开始"

# 连接到指定沙箱
codebuddy --sandbox --sandbox-id sb_abc123 "继续工作"

# 退出时清理沙箱
codebuddy --sandbox --sandbox-kill "临时测试"

# Teleport 模式 - 连接到远程创建的沙箱
codebuddy --teleport session_abc123XYZ4567890 "连接到远程沙箱"
```

### **沙箱环境变量**

```
E2B_API_KEY                          E2B API 密钥 （E2B 沙箱必需）
E2B_TEMPLATE                         E2B 模板 ID （默认： base)
CODEBUDDY_SANDBOX_IMAGE              自定义 Docker 镜像 （容器沙箱）
```

## **下一步**

掌握 CLI 命令后，您可以：

- **[学习交互模式](https://www.codebuddy.ai/docs/zh/cli/interactive-mode)** - 掌握键盘快捷键和技巧
- **[探索斜杠命令](https://www.codebuddy.ai/docs/zh/cli/slash-commands)** - 了解内置命令
- **[Skills 技能系统](https://www.codebuddy.ai/docs/zh/cli/skills)** - 扩展 AI 专业能力
- **[学习 MCP 集成](https://www.codebuddy.ai/docs/zh/cli/mcp)** - 扩展工具能力

***

*精确的命令行操作是高效开发的基础*

**最后更新: 2026/4/23 14:36**

<br />

# **CodeBuddy 插件参考文档**

> 完整的 CodeBuddy 插件系统技术参考，包括组件规范、CLI 命令和开发工具。

**插件**是一个自包含的组件目录，用于扩展 CodeBuddy 的自定义功能。插件组件包括技能 (Skills)、代理 (Agents)、钩子 (Hooks)、MCP 服务器和 LSP 服务器。

***

## **一、插件组件参考**

### **1. Skills（技能）**

插件通过添加技能来扩展 CodeBuddy，创建 `/name` 快捷方式供用户或 AI 助手调用。

**位置**：插件根目录的 `skills/` 或 `commands/` 目录

**文件格式**：技能是包含 `SKILL.md` 的目录；命令是简单的 Markdown 文件

**目录结构**:

```
skills/
├── pdf-processor/
│   ├── SKILL.md
│   ├── reference.md （可选）
│   └── scripts/ （可选）
└── code-reviewer/
    └── SKILL.md
```

**集成行为**:

- 安装插件时自动发现技能和命令
- AI 助手可根据任务上下文自动调用
- 技能可以包含辅助文件和脚本

详见 **[Skills](https://www.codebuddy.ai/docs/zh/cli/skills)**。

### **2. Agents（代理）**

插件可以提供专门的子代理用于特定任务，AI 助手可以在适当时自动调用。

**位置**：插件根目录的 `agents/` 目录

**格式**：描述代理能力的 Markdown 文件

**Frontmatter 配置**:

```
---
name: agent-name
description: 代理的专长和调用时机
model: sonnet
effort: medium
maxTurns: 20
disallowedTools: Write, Edit
---

代理的详细系统提示词，描述其角色、专长和行为。
```

插件代理支持以下 frontmatter 字段：`name`、`description`、`model`、`effort`、`maxTurns`、`tools`、`disallowedTools`、`skills`、`memory`、`background` 和 `isolation`。`isolation` 的唯一有效值为 `"worktree"`。出于安全原因，插件代理不支持 `hooks`、`mcpServers` 和 `permissionMode` 字段。

**集成方式**:

- 代理出现在 `/agents` 界面中
- AI 助手可根据任务上下文自动调用代理
- 用户也可手动调用代理
- 插件代理与内置代理协同工作

详见 **[Subagents](https://www.codebuddy.ai/docs/zh/cli/sub-agents)**。

### **3. Hooks（钩子）**

插件可以提供事件处理器，自动响应 CodeBuddy 事件。

**位置**：插件根目录的 `hooks/hooks.json`，或在 `plugin.json` 中内联配置

**格式**: JSON 配置，包含事件匹配器和操作

**配置示例**:

```
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Write|Edit",
        "hooks": [
          {
            "type": "command",
            "command": "${CODEBUDDY_PLUGIN_ROOT}/scripts/format-code.sh"
          }
        ]
      }
    ]
  }
}
```

**可用事件**:

插件钩子响应与用户定义钩子相同的生命周期事件：

**事件**

**触发时机**

`SessionStart`

会话开始或恢复时

`UserPromptSubmit`

用户提交提示时，在 AI 处理之前

`PreToolUse`

工具调用执行前，可以阻断

`PermissionRequest`

权限对话框出现时

`PermissionDenied`

工具调用被自动模式分类器拒绝时。返回 `{retry: true}` 告知模型可重试

`PostToolUse`

工具调用成功后

`PostToolUseFailure`

工具调用失败后

`Notification`

CodeBuddy 发送通知时

`SubagentStart`

子代理启动时

`SubagentStop`

子代理完成时

`TaskCreated`

通过 `TaskCreate` 创建任务时

`TaskCompleted`

任务被标记为已完成时

`Stop`

AI 完成响应时

`StopFailure`

轮次因 API 错误结束时。输出和退出码被忽略

`TeammateIdle`

团队成员即将进入空闲时

`InstructionsLoaded`

CODEBUDDY.md 或 `.codebuddy/rules/*.md` 文件加载到上下文时

`ConfigChange`

会话期间配置文件变更时

`CwdChanged`

工作目录变更时（例如 AI 执行 `cd` 命令）

`FileChanged`

监视的文件在磁盘上变更时。`matcher` 字段指定要监视的文件名

`WorktreeCreate`

通过 `--worktree` 或 `isolation: "worktree"` 创建工作树时

`WorktreeRemove`

工作树被移除时（会话退出或子代理完成时）

`PreCompact`

上下文压缩之前

`PostCompact`

上下文压缩完成后

`Elicitation`

MCP 服务器在工具调用期间请求用户输入时

`ElicitationResult`

用户响应 MCP elicitation 后，响应发回服务器之前

`SessionEnd`

会话终止时

**钩子类型**:

- `command`：执行 shell 命令或脚本
- `http`：将事件 JSON 作为 POST 请求发送到 URL
- `prompt`：使用 LLM 评估提示词（使用 `$ARGUMENTS` 占位符获取上下文）
- `agent`：运行带工具的代理验证器，用于复杂验证任务

### **4. MCP Servers（MCP 服务器）**

插件可以捆绑 Model Context Protocol (MCP) 服务器，将 CodeBuddy 与外部工具和服务连接。

**位置**：插件根目录的 `.mcp.json`，或在 `plugin.json` 中内联配置

**格式**: 标准 MCP 服务器配置

**配置示例**:

```
{
  "mcpServers": {
    "plugin-database": {
      "command": "${CODEBUDDY_PLUGIN_ROOT}/servers/db-server",
      "args": ["--config", "${CODEBUDDY_PLUGIN_ROOT}/config.json"],
      "env": {
        "DB_PATH": "${CODEBUDDY_PLUGIN_ROOT}/data"
      }
    },
    "plugin-api-client": {
      "command": "npx",
      "args": ["@company/mcp-server", "--plugin-mode"],
      "cwd": "${CODEBUDDY_PLUGIN_ROOT}"
    }
  }
}
```

**集成行为**:

- 插件启用时自动启动 MCP 服务器
- 服务器作为标准 MCP 工具出现在工具包中
- 服务器功能与现有工具无缝集成
- 插件服务器可独立于用户 MCP 服务器进行配置

### **5. LSP Servers（LSP 服务器）**

> **提示**: 需要使用 LSP 插件? 可从官方市场安装——在 `/plugin` Discover 标签中搜索 "lsp"。本节介绍如何为官方市场未涵盖的语言创建 LSP 插件。

插件可以提供 **[Language Server Protocol](https://microsoft.github.io/language-server-protocol/)** (LSP) 服务器，为 AI 助手在代码库上工作时提供实时代码智能支持。

LSP 集成提供:

- **即时诊断**: AI 助手在每次编辑后立即看到错误和警告
- **代码导航**: 跳转到定义、查找引用和悬停信息
- **语言感知**: 代码符号的类型信息和文档

**位置**: 插件根目录的 `.lsp.json` 文件，或在 `plugin.json` 中内联配置

**格式**: JSON 配置，将语言服务器名称映射到其配置

**`.lsp.json` 文件格式**:

```
{
  "go": {
    "command": "gopls",
    "args": ["serve"],
    "extensionToLanguage": {
      ".go": "go"
    }
  }
}
```

**在 `plugin.json` 中内联配置**:

```
{
  "name": "my-plugin",
  "lspServers": {
    "go": {
      "command": "gopls",
      "args": ["serve"],
      "extensionToLanguage": {
        ".go": "go"
      }
    }
  }
}
```

**必需字段**:

**字段**

**描述**

`command`

要执行的 LSP 二进制文件（必须在 PATH 中）

`extensionToLanguage`

将文件扩展名映射到语言标识符

**可选字段**:

**字段**

**描述**

`args`

LSP 服务器的命令行参数

`transport`

通信传输方式: `stdio`（默认）或 `socket`

`env`

启动服务器时设置的环境变量

`initializationOptions`

在初始化期间传递给服务器的选项

`settings`

通过 `workspace/didChangeConfiguration` 传递的设置

`workspaceFolder`

服务器的工作区文件夹路径

`startupTimeout`

等待服务器启动的最大时间（毫秒）

`shutdownTimeout`

等待优雅关闭的最大时间（毫秒）

`restartOnCrash`

服务器崩溃时是否自动重启

`maxRestarts`

放弃前的最大重启尝试次数

> **警告**: **您必须单独安装语言服务器二进制文件。** LSP 插件配置 CodeBuddy 如何连接到语言服务器，但不包含服务器本身。如果在 `/plugin` Errors 标签中看到 `Executable not found in $PATH` 错误，请为您的语言安装所需的二进制文件。

**可用的 LSP 插件**:

**插件**

**语言服务器**

**安装命令**

`pyright-lsp`

Pyright (Python)

`pip install pyright` 或 `npm install -g pyright`

`typescript-lsp`

TypeScript Language Server

`npm install -g typescript-language-server typescript`

`rust-lsp`

rust-analyzer

**[参见 rust-analyzer 安装](https://rust-analyzer.github.io/manual.html#installation)**

先安装语言服务器，然后从市场安装插件。

***

## **二、插件安装作用域**

安装插件时，选择一个**作用域**来决定插件的可用范围：

**作用域**

**设置文件**

**使用场景**

`user`

`~/.codebuddy/settings.json`

个人插件，所有项目可用（默认）

`project`

`.codebuddy/settings.json`

团队插件，通过版本控制共享

`local`

`.codebuddy/settings.local.json`

项目特定插件，被 gitignore

`managed`

托管设置

托管插件（只读，仅支持更新）

插件使用与其他 CodeBuddy 配置相同的作用域系统。详见 **[Settings](https://www.codebuddy.ai/docs/zh/cli/settings)**。

***

## **三、插件清单架构（plugin.json）**

`.codebuddy-plugin/plugin.json`（或 `.claude-plugin/plugin.json`）文件定义插件的元数据和配置。本节记录所有支持的字段和选项。

清单是可选的。如果省略，CodeBuddy 会自动发现**[默认位置](https://www.codebuddy.ai/docs/zh/cli/plugins-reference#%E6%96%87%E4%BB%B6%E4%BD%8D%E7%BD%AE%E5%8F%82%E8%80%83)**中的组件，并从目录名派生插件名称。当需要提供元数据或自定义组件路径时使用清单。

### **完整架构示例**

```
{
  "name": "plugin-name",
  "version": "1.2.0",
  "description": "插件简短描述",
  "author": {
    "name": "作者名称",
    "email": "[email protected]",
    "url": "https://github.com/author"
  },
  "homepage": "https://docs.example.com/plugin",
  "repository": "https://github.com/author/plugin",
  "license": "MIT",
  "keywords": ["keyword1", "keyword2"],
  "commands": ["./custom/commands/special.md"],
  "agents": "./custom/agents/",
  "skills": "./custom/skills/",
  "hooks": "./config/hooks.json",
  "mcpServers": "./mcp-config.json",
  "outputStyles": "./styles/",
  "lspServers": "./.lsp.json"
}
```

### **必需字段**

如果包含清单，`name` 是唯一必需字段。

**字段**

**类型**

**描述**

**示例**

`name`

string

唯一标识符（kebab-case，无空格）

`"deployment-tools"`

此名称用于组件命名空间。例如在 UI 中，插件 `plugin-dev` 的代理 `agent-creator` 将显示为 `plugin-dev:agent-creator`。

### **元数据字段**

**字段**

**类型**

**描述**

**示例**

`version`

string

语义化版本。如果在 marketplace 条目中也设置了，`plugin.json` 优先。只需在一处设置。

`"2.1.0"`

`description`

string

插件用途简述

`"部署自动化工具"`

`author`

object

作者信息

`{"name": "开发团队", "email": "[email protected]"}`

`homepage`

string

文档 URL

`"https://docs.example.com"`

`repository`

string

源代码 URL

`"https://github.com/user/plugin"`

`license`

string

许可证标识符

`"MIT"`, `"Apache-2.0"`

`keywords`

array

发现标签

`["deployment", "ci-cd"]`

### **组件路径字段**

**字段**

**类型**

**描述**

**示例**

`commands`

string | array

自定义命令文件/目录（替换默认 `commands/`）

`"./custom/cmd.md"` 或 `["./cmd1.md"]`

`agents`

string | array

自定义代理文件（替换默认 `agents/`）

`"./custom/agents/reviewer.md"`

`skills`

string | array

自定义技能目录（替换默认 `skills/`）

`"./custom/skills/"`

`hooks`

string | array | object

钩子配置路径或内联配置

`"./my-extra-hooks.json"`

`mcpServers`

string | array | object

MCP 配置路径或内联配置

`"./my-extra-mcp-config.json"`

`outputStyles`

string | array

自定义输出样式文件/目录（替换默认 `output-styles/`）

`"./styles/"`

`lspServers`

string | array | object

LSP 配置路径或内联配置

`"./.lsp.json"`

`userConfig`

object

启用时提示用户配置的值。详见**[用户配置](https://www.codebuddy.ai/docs/zh/cli/plugins-reference#%E7%94%A8%E6%88%B7%E9%85%8D%E7%BD%AE)**

见下方

`channels`

array

消息注入的频道声明。详见**[频道](https://www.codebuddy.ai/docs/zh/cli/plugins-reference#%E9%A2%91%E9%81%93)**

见下方

### **用户配置**

`userConfig` 字段声明在插件启用时 CodeBuddy 提示用户输入的值。用此替代要求用户手动编辑 `settings.json`。

```
{
  "userConfig": {
    "api_endpoint": {
      "description": "您团队的 API 端点",
      "sensitive": false
    },
    "api_token": {
      "description": "API 认证令牌",
      "sensitive": true
    }
  }
}
```

键必须是有效标识符。每个值可作为 `${user_config.KEY}` 在 MCP 和 LSP 服务器配置、钩子命令中替换，以及（仅非敏感值）在技能和代理内容中替换。值也作为 `CODEBUDDY_PLUGIN_OPTION_<KEY>` 环境变量导出到插件子进程。

非敏感值存储在 `settings.json` 的 `pluginConfigs[<plugin-id>].options` 中。敏感值存储到系统密钥链（或在密钥链不可用时存储到 `~/.codebuddy/.credentials.json`）。密钥链存储与 OAuth 令牌共享，总限制约 2 KB，因此敏感值应保持较小。

### **频道**

`channels` 字段允许插件声明一个或多个消息频道，用于向对话中注入内容。每个频道绑定到插件提供的一个 MCP 服务器。

```
{
  "channels": [
    {
      "server": "telegram",
      "userConfig": {
        "bot_token": { "description": "Telegram bot token", "sensitive": true },
        "owner_id": { "description": "您的 Telegram 用户 ID", "sensitive": false }
      }
    }
  ]
}
```

`server` 字段是必需的，必须匹配插件 `mcpServers` 中的一个键。可选的按频道 `userConfig` 使用与顶层相同的 schema，允许在启用插件时提示输入 bot token 或 owner ID。

### **路径行为规则**

对于 `commands`、`agents`、`skills` 和 `outputStyles`，自定义路径会**替换**默认目录。如果清单指定了 `commands`，则不会扫描默认的 `commands/` 目录。Hooks、MCP servers 和 LSP servers 对处理多个来源有不同的语义。

- 所有路径必须相对于插件根目录并以 `./` 开头
- 自定义路径中的组件使用相同的命名和命名空间规则
- 可以将多个路径指定为数组
- 要保留默认目录并添加更多路径，在数组中包含默认目录：`"commands": ["./commands/", "./extras/deploy.md"]`

**路径示例**:

```
{
  "commands": [
    "./specialized/deploy.md",
    "./utilities/batch-process.md"
  ],
  "agents": [
    "./custom-agents/reviewer.md",
    "./custom-agents/tester.md"
  ]
}
```

### **环境变量**

CodeBuddy 提供两个变量用于引用插件路径。两者都在技能内容、代理内容、钩子命令、MCP 或 LSP 服务器配置中的任何位置进行内联替换。两者也作为环境变量导出到钩子进程和 MCP 或 LSP 服务器子进程。

**`${CODEBUDDY_PLUGIN_ROOT}`**：插件安装目录的绝对路径。用于引用插件捆绑的脚本、二进制文件和配置文件。此路径在插件更新时会变化，因此写入此处的文件不会在更新后保留。

**兼容性**：同时支持 `${CLAUDE_PLUGIN_ROOT}` 变量名以兼容 Claude Code 插件。

**`${CODEBUDDY_PLUGIN_DATA}`**：用于插件状态的持久化目录，在更新后保留。用于已安装的依赖项（如 `node_modules` 或 Python 虚拟环境）、生成的代码、缓存以及任何需要跨插件版本持久化的文件。首次引用时自动创建该目录。

**兼容性**：同时支持 `${CLAUDE_PLUGIN_DATA}` 变量名。

```
{
  "hooks": {
    "PostToolUse": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "${CODEBUDDY_PLUGIN_ROOT}/scripts/process.sh"
          }
        ]
      }
    ]
  }
}
```

#### **持久化数据目录**

`${CODEBUDDY_PLUGIN_DATA}` 目录解析为 `~/.codebuddy/plugins/data/{id}/`，其中 `{id}` 是插件标识符，`a-z`、`A-Z`、`0-9`、`_` 和 `-` 之外的字符替换为 `-`。例如安装为 `formatter@my-marketplace` 的插件，目录为 `~/.codebuddy/plugins/data/formatter-my-marketplace/`。

常见用法是安装语言依赖一次并跨会话和插件更新重用。由于数据目录的生命周期超越任何单一插件版本，仅检查目录存在性无法检测更新何时更改了插件的依赖清单。推荐的模式是将捆绑的清单与数据目录中的副本进行比较，不同时重新安装。

以下 `SessionStart` 钩子在首次运行时安装 `node_modules`，并在插件更新包含更改的 `package.json` 时再次安装：

```
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "diff -q \"${CODEBUDDY_PLUGIN_ROOT}/package.json\" \"${CODEBUDDY_PLUGIN_DATA}/package.json\" >/dev/null 2>&1 || (cd \"${CODEBUDDY_PLUGIN_DATA}\" && cp \"${CODEBUDDY_PLUGIN_ROOT}/package.json\" . && npm install) || rm -f \"${CODEBUDDY_PLUGIN_DATA}/package.json\""
          }
        ]
      }
    ]
  }
}
```

`diff` 在存储副本缺失或与捆绑副本不同时以非零退出，涵盖首次运行和依赖变更更新。如果 `npm install` 失败，尾部的 `rm` 删除已复制的清单，以便下次会话重试。

捆绑在 `${CODEBUDDY_PLUGIN_ROOT}` 中的脚本可以针对持久化的 `node_modules` 运行：

```
{
  "mcpServers": {
    "routines": {
      "command": "node",
      "args": ["${CODEBUDDY_PLUGIN_ROOT}/server.js"],
      "env": {
        "NODE_PATH": "${CODEBUDDY_PLUGIN_DATA}/node_modules"
      }
    }
  }
}
```

卸载插件（从最后一个安装作用域）时数据目录会自动删除。`/plugin` 界面显示目录大小并在删除前提示。CLI 默认删除；传递 `--keep-data` 可保留。

***

## **四、插件缓存和文件解析**

插件有两种指定方式：

- 通过 `codebuddy --plugin-dir`，仅在会话期间有效
- 通过市场安装，适用于未来会话

出于安全和验证目的，CodeBuddy 将**市场**插件复制到用户的本地**插件缓存**（`~/.codebuddy/plugins/cache`），而不是原地使用。理解此行为对于开发引用外部文件的插件很重要。

### **路径遍历限制**

已安装的插件无法引用其目录之外的文件。遍历到插件根目录之外的路径（如 `../shared-utils`）在安装后不会工作，因为那些外部文件不会被复制到缓存中。

### **使用外部依赖**

如果插件需要访问其目录之外的文件，可以在插件目录中创建指向外部文件的符号链接。符号链接在复制过程中会被保留：

```
# 在插件目录内
ln -s /path/to/shared-utils ./shared-utils
```

符号链接的内容将被复制到插件缓存中。这在保持缓存系统安全优势的同时提供了灵活性。

***

## **五、插件目录结构**

### **标准插件布局**

完整的插件遵循以下结构：

```
enterprise-plugin/
├── .codebuddy-plugin/        # 元数据目录（可选）
│   └── plugin.json             # 插件清单
├── commands/                 # 默认命令位置
│   ├── status.md
│   └── logs.md
├── agents/                   # 默认代理位置
│   ├── security-reviewer.md
│   ├── performance-tester.md
│   └── compliance-checker.md
├── skills/                   # 代理技能
│   ├── code-reviewer/
│   │   └── SKILL.md
│   └── pdf-processor/
│       ├── SKILL.md
│       └── scripts/
├── output-styles/            # 输出样式定义
│   └── terse.md
├── hooks/                    # 钩子配置
│   ├── hooks.json            # 主钩子配置
│   └── security-hooks.json   # 额外钩子
├── bin/                      # 插件可执行文件，添加到 PATH
│   └── my-tool               # 可在 Bash 工具中作为裸命令调用
├── settings.json             # 插件默认设置
├── .mcp.json                 # MCP 服务器定义
├── .lsp.json                 # LSP 服务器配置
├── scripts/                  # 钩子和实用脚本
│   ├── security-scan.sh
│   ├── format-code.py
│   └── deploy.js
├── LICENSE                   # 许可证文件
└── CHANGELOG.md              # 版本历史
```

> **重要**: `.codebuddy-plugin/` 目录包含 `plugin.json` 文件。所有其他目录（`commands/`、`agents/`、`skills/`、`output-styles/`、`hooks/`）必须位于插件根目录，而不是 `.codebuddy-plugin/` 内部。同时兼容 `.claude-plugin/` 目录。

### **文件位置参考**

**组件**

**默认位置**

**用途**

**Manifest**

`.codebuddy-plugin/plugin.json`

插件元数据和配置（可选）

**Commands**

`commands/`

技能 Markdown 文件（旧版；新技能使用 `skills/`）

**Agents**

`agents/`

子代理 Markdown 文件

**Skills**

`skills/`

带有 `<name>/SKILL.md` 结构的技能

**Output styles**

`output-styles/`

输出样式定义

**Hooks**

`hooks/hooks.json`

钩子配置

**MCP servers**

`.mcp.json`

MCP 服务器定义

**LSP servers**

`.lsp.json`

语言服务器配置

**Executables**

`bin/`

添加到 Bash 工具 PATH 的可执行文件。启用插件时此目录中的文件可在任何 Bash 工具调用中作为裸命令调用

**Settings**

`settings.json`

启用插件时应用的默认配置。目前仅支持 **[agent](https://www.codebuddy.ai/docs/zh/cli/sub-agents)** 设置

***

## **六、CLI 命令参考**

CodeBuddy 提供 CLI 命令用于非交互式插件管理，适用于脚本和自动化。

### **plugin install**

从可用市场安装插件。

```
codebuddy plugin install <plugin> [options]
```

**参数**:

- `<plugin>`: 插件名称或 `plugin-name@marketplace-name`（指定特定市场）

**选项**:

**选项**

**描述**

**默认**

`-s, --scope <scope>`

安装作用域：`user`、`project` 或 `local`

`user`

`-h, --help`

显示命令帮助

<br />

作用域决定已安装插件添加到哪个设置文件。例如 `--scope project` 写入 `.codebuddy/settings.json` 中的 `enabledPlugins`，使插件对所有克隆该项目仓库的人可用。

**示例**:

```
# 安装到用户作用域（默认）
codebuddy plugin install formatter@my-marketplace

# 安装到项目作用域（与团队共享）
codebuddy plugin install formatter@my-marketplace --scope project

# 安装到本地作用域（被 gitignore）
codebuddy plugin install formatter@my-marketplace --scope local
```

### **plugin uninstall**

移除已安装的插件。

```
codebuddy plugin uninstall <plugin> [options]
```

**参数**:

- `<plugin>`: 插件名称或 `plugin-name@marketplace-name`

**选项**:

**选项**

**描述**

**默认**

`-s, --scope <scope>`

从指定作用域卸载：`user`、`project` 或 `local`

`user`

`--keep-data`

保留插件的**[持久化数据目录](https://www.codebuddy.ai/docs/zh/cli/plugins-reference#%E6%8C%81%E4%B9%85%E5%8C%96%E6%95%B0%E6%8D%AE%E7%9B%AE%E5%BD%95)**

<br />

`-h, --help`

显示命令帮助

<br />

**别名**: `remove`, `rm`

默认情况下，从最后一个剩余作用域卸载时也会删除插件的 `${CODEBUDDY_PLUGIN_DATA}` 目录。使用 `--keep-data` 可保留它，例如在测试新版本后重新安装时。

### **plugin enable**

启用已禁用的插件。

```
codebuddy plugin enable <plugin> [options]
```

**参数**:

- `<plugin>`: 插件名称或 `plugin-name@marketplace-name`

**选项**:

**选项**

**描述**

**默认**

`-s, --scope <scope>`

启用作用域：`user`、`project` 或 `local`

`user`

`-h, --help`

显示命令帮助

<br />

### **plugin disable**

禁用插件但不卸载。

```
codebuddy plugin disable <plugin> [options]
```

**参数**:

- `<plugin>`: 插件名称或 `plugin-name@marketplace-name`

**选项**:

**选项**

**描述**

**默认**

`-s, --scope <scope>`

禁用作用域：`user`、`project` 或 `local`

`user`

`-h, --help`

显示命令帮助

<br />

### **plugin update**

更新插件到最新版本。

```
codebuddy plugin update <plugin> [options]
```

**参数**:

- `<plugin>`: 插件名称或 `plugin-name@marketplace-name`

**选项**:

**选项**

**描述**

**默认**

`-s, --scope <scope>`

更新作用域：`user`、`project`、`local` 或 `managed`

`user`

`-h, --help`

显示命令帮助

<br />

### **市场管理**

```
# 添加市场
codebuddy plugin marketplace add <source> [--name <name>]

# 列出市场
codebuddy plugin marketplace list

# 更新市场
codebuddy plugin marketplace update <name>

# 删除市场
codebuddy plugin marketplace remove <name>
```

**市场源格式**:

```
# 本地目录
codebuddy plugin marketplace add /path/to/marketplace

# GitHub 简写
codebuddy plugin marketplace add owner/repo

# Git URL
codebuddy plugin marketplace add https://github.com/owner/repo.git

# HTTP URL (marketplace.json)
codebuddy plugin marketplace add https://example.com/marketplace.json
```

***

## **七、调试和开发工具**

### **调试命令**

使用 `codebuddy --debug` 查看插件加载详情：

**显示内容**:

- 正在加载哪些插件
- 插件清单中的任何错误
- 命令、代理和钩子注册
- MCP 服务器初始化

### **常见问题排查**

**问题**

**原因**

**解决方案**

插件未加载

无效的 `plugin.json`

运行 `codebuddy plugin validate` 或 `/plugin validate` 检查 `plugin.json`、skill/agent/command frontmatter 和 `hooks/hooks.json` 的语法和 schema 错误

命令未出现

错误的目录结构

确保 `commands/` 在插件根目录，而不是 `.codebuddy-plugin/` 内部

钩子未触发

脚本不可执行

运行 `chmod +x script.sh`

MCP 服务器失败

缺少 `${CODEBUDDY_PLUGIN_ROOT}`

所有插件路径使用此变量

路径错误

使用了绝对路径

所有路径必须是相对路径并以 `./` 开头

LSP `Executable not found in $PATH`

语言服务器未安装

安装二进制文件（例如：`npm install -g typescript-language-server typescript`）

### **常见错误信息**

**清单验证错误**:

- `Invalid JSON syntax: Unexpected token } in JSON at position 142`: 检查缺少的逗号、多余的逗号或未引用的字符串
- `Plugin has an invalid manifest file at .codebuddy-plugin/plugin.json. Validation errors: name: Required`: 缺少必需字段
- `Plugin has a corrupt manifest file at .codebuddy-plugin/plugin.json. JSON parse error: ...`: JSON 语法错误

**插件加载错误**:

- `Warning: No commands found in plugin my-plugin custom directory: ./cmds. Expected .md files or SKILL.md in subdirectories.`: 命令路径存在但不包含有效的命令文件
- `Plugin directory not found at path: ./plugins/my-plugin. Check that the marketplace entry has the correct path.`: marketplace.json 中的 `source` 路径指向不存在的目录
- `Plugin my-plugin has conflicting manifests: both plugin.json and marketplace entry specify components.`: 删除重复的组件定义或在 marketplace 条目中移除 `strict: false`

### **钩子排查**

**钩子脚本未执行**:

1. 检查脚本是否可执行：`chmod +x ./scripts/your-script.sh`
2. 验证 shebang 行：第一行应为 `#!/bin/bash` 或 `#!/usr/bin/env bash`
3. 检查路径是否使用 `${CODEBUDDY_PLUGIN_ROOT}`：`"command": "${CODEBUDDY_PLUGIN_ROOT}/scripts/your-script.sh"`
4. 手动测试脚本：`./scripts/your-script.sh`

**钩子未在预期事件上触发**:

1. 验证事件名称正确（区分大小写）：`PostToolUse`，而不是 `postToolUse`
2. 检查 matcher 模式是否匹配目标工具：`"matcher": "Write|Edit"` 用于文件操作
3. 确认钩子类型有效：`command`、`http`、`prompt` 或 `agent`

### **MCP 服务器排查**

**服务器未启动**:

1. 检查命令是否存在且可执行
2. 验证所有路径使用 `${CODEBUDDY_PLUGIN_ROOT}` 变量
3. 检查 MCP 服务器日志：`codebuddy --debug` 显示初始化错误
4. 在 CodeBuddy 之外手动测试服务器

**服务器工具未出现**:

1. 确保服务器在 `.mcp.json` 或 `plugin.json` 中正确配置
2. 验证服务器正确实现了 MCP 协议
3. 检查调试输出中的连接超时

### **目录结构错误**

**症状**: 插件加载但组件（命令、代理、钩子）缺失。

**正确结构**: 组件必须在插件根目录，而不是 `.codebuddy-plugin/` 内部。只有 `plugin.json` 属于 `.codebuddy-plugin/`。

```
my-plugin/
├── .codebuddy-plugin/
│   └── plugin.json      <- 只有清单在这里
├── commands/             <- 在根级别
├── agents/               <- 在根级别
└── hooks/                <- 在根级别
```

如果组件在 `.codebuddy-plugin/` 内部，将它们移到插件根目录。

**调试检查清单**:

1. 运行 `codebuddy --debug` 并查找 "loading plugin" 消息
2. 检查每个组件目录是否在调试输出中列出
3. 验证文件权限允许读取插件文件

***

## **八、版本管理参考**

### **语义化版本控制**

遵循语义化版本控制进行插件发布：

```
{
  "name": "my-plugin",
  "version": "2.1.0"
}
```

**版本格式**: `MAJOR.MINOR.PATCH`

- **MAJOR**: 不兼容的 API 更改
- **MINOR**: 向后兼容的功能添加
- **PATCH**: 向后兼容的错误修复

**最佳实践**:

- 第一个稳定版本从 `1.0.0` 开始
- 分发更改前更新 `plugin.json` 中的版本
- 在 `CHANGELOG.md` 文件中记录变更
- 使用预发布版本（如 `2.0.0-beta.1`）进行测试

> **注意**: CodeBuddy 使用版本来判断是否需要更新插件。如果更改了插件代码但未更新 `plugin.json` 中的版本，由于缓存机制，现有用户不会看到变更。
>
> 如果插件在**[市场](https://www.codebuddy.ai/docs/zh/cli/plugin-marketplaces)**目录中，可以通过 `marketplace.json` 管理版本，并从 `plugin.json` 中省略 `version` 字段。

***

## **九、与 Claude Code 的兼容性**

CodeBuddy 插件系统在设计上兼容 Claude Code 插件规范，但存在以下差异：

### **命名差异**

**概念**

**Claude Code**

**CodeBuddy**

元数据目录

`.claude-plugin/`

`.codebuddy-plugin/`（优先）或 `.claude-plugin/`（兼容）

环境变量

`${CLAUDE_PLUGIN_ROOT}`

`${CODEBUDDY_PLUGIN_ROOT}`（优先）或 `${CLAUDE_PLUGIN_ROOT}`（兼容）

数据目录变量

`${CLAUDE_PLUGIN_DATA}`

`${CODEBUDDY_PLUGIN_DATA}`（优先）或 `${CLAUDE_PLUGIN_DATA}`（兼容）

### **迁移指南**

从 Claude Code 迁移到 CodeBuddy：

1. 可选择将 `.claude-plugin/` 重命名为 `.codebuddy-plugin/`
2. 可选择将脚本中的 `${CLAUDE_PLUGIN_ROOT}` 替换为 `${CODEBUDDY_PLUGIN_ROOT}`
3. 可选择将 `${CLAUDE_PLUGIN_DATA}` 替换为 `${CODEBUDDY_PLUGIN_DATA}`

**注意**：保持原有命名也完全兼容，CodeBuddy 会自动识别。

***

## **相关资源**

- **[插件](https://www.codebuddy.ai/docs/zh/cli/plugins)** - 教程和实用指南
- **[插件市场](https://www.codebuddy.ai/docs/zh/cli/plugin-marketplaces)** - 创建和管理市场
- **[Skills](https://www.codebuddy.ai/docs/zh/cli/skills)** - 技能开发详情
- **[Subagents](https://www.codebuddy.ai/docs/zh/cli/sub-agents)** - 代理配置和能力
- **[Hooks](https://www.codebuddy.ai/docs/zh/cli/hooks)** - 事件处理和自动化
- **[MCP](https://www.codebuddy.ai/docs/zh/cli/mcp)** - 外部工具集成
- **[Settings](https://www.codebuddy.ai/docs/zh/cli/settings)** - 插件配置选项

**最后更新: 2026/4/3 23:53**

# **成本管理**

> 通过多场景模型选择和异步压缩策略优化成本，在保证效果的同时实现更快的响应速度和更低的使用成本。

CodeBuddy Code 每次交互都会消耗 Token。成本因代码库大小、查询复杂度和对话长度而异。本文档介绍如何追踪成本、多场景模型机制和降低 Token 消耗。

## **追踪成本**

### **使用 /cost 命令**

`/cost` 命令提供当前会话的详细 Token 使用统计：

```
/cost
  ⎿ Total duration (API):  9m 35.6s
    Total duration (wall): 22m 14.9s
    Total code changes:    0 lines added, 0 lines removed
    Usage by model:
         claude-sonnet-4:  875.5k input, 11.7k output, 714.3k cache read, 0 cache write
```

### **使用 /context 命令**

`/context` 命令可以分析当前上下文的占用情况，查看不同类型上下文的大小分布：

```
> /context 
  ⎿  Context Usage
     ⛁ ⛁ ⛁ ⛁ ⛁ ⛁ ⛁ ⛁ ⛁ ⛁    glm-4.7 · 38.1k/200.0k tokens (19.1%)
     ⛁ ⛁ ⛁ ⛁ ⛁ ⛁ ⛁ ⛁ ⛁ ⛁
     ⛁ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶    ⛁ System prompt: 2.1k tokens (1.1%)
     ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶    ⛁ System tools: 16.4k tokens (8.2%)
     ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶    ⛁ Memory files: 3.7k tokens (1.9%)
     ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶    ⛁ Messages: 15.9k tokens (7.9%)
     ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶    ⛶ Free space: 145.9k (72.9%)
     ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶    ⛝ Autocompact buffer: 16.0k tokens (8.0%)
     ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶
     ⛶ ⛶ ⛶ ⛶ ⛝ ⛝ ⛝ ⛝ ⛝ ⛝

     Memory files · /memory
     └ /Users/yangsubo/.codebuddy/CODEBUDDY.md (User): 18 tokens
     └ /Users/yangsubo/CODEBUDDY.md (Project): 15 tokens
     └ /Users/yangsubo/workspace/genie/CODEBUDDY.md (Project): 1.2k tokens
     └ /Users/yangsubo/workspace/genie/packages/agent-cli/CODEBUDDY.md (Project): 2.5k tokens

     Skills and slash commands · /skills

     Project
     └ release: 1.1k tokens
     └ gen-drawio: 846 tokens
     └ task-manager: 815 tokens
     └ task-add: 730 tokens
     └ task-done: 525 tokens
     └ mr: 519 tokens
     └ task-start: 493 tokens
     └ my-task: 238 tokens
     └ task-list: 233 tokens
     └ security-review: 30 tokens
```

通过 `/context` 可以快速识别哪些内容占用了大量上下文空间，从而有针对性地优化。

## **多场景模型机制**

不同的任务场景对模型能力的要求不同。简单的文件搜索、快速查询等任务使用轻量级模型即可完成；复杂的架构设计、多步骤推理等任务需要更强大的推理模型。

通过为不同场景自动选择不同的模型，可以实现：

- **效果最优**：复杂任务使用高能力模型，确保质量
- **速度更快**：简单任务使用轻量模型，响应更迅速
- **成本更低**：避免在简单任务上使用昂贵的高端模型

### **场景类型**

**场景类型**

**说明**

**典型用例**

`default`

默认模型，平衡性能与成本

一般编程任务、代码编写

`lite`

轻量快速模型，低成本高速度

文件搜索、简单查询、快速操作

`reasoning`

推理增强模型，强大分析能力

复杂分析、架构决策、多步骤推理

### **自动模型选择**

CodeBuddy Code 会根据任务类型自动选择合适的场景模型。当子代理执行时，系统会根据用户当前选择的主模型，自动解析出对应的场景模型。

例如，`contentAnalyzer` 等轻量级子代理会自动使用 `lite` 模型，在保证功能的同时降低成本、提升速度。

Agent 工具支持通过 `model` 参数指定场景类型：

- `default`：继承父模型，适用于一般任务
- `lite`：快速且低成本，适用于简单搜索、快速文件操作
- `reasoning`：增强推理能力，适用于复杂分析、架构决策

## **降低 Token 消耗**

Token 成本随上下文大小增长：CodeBuddy Code 处理的上下文越大，消耗的 Token 越多。CodeBuddy Code 通过 Prompt 缓存（减少重复内容如系统提示的成本）和自动压缩（在接近上下文限制时压缩对话历史）自动优化成本。

以下策略帮助你保持较小的上下文，降低每条消息的成本。

### **主动管理上下文**

使用 `/cost` 检查当前 Token 使用情况。

- **任务间清理**：切换到不相关的工作时使用 `/clear` 重新开始。过时的上下文会在后续每条消息中浪费 Token。清理前使用 `/rename` 以便之后通过 `/resume` 返回。
- **添加自定义压缩指令**：`/compact Focus on code samples and API usage` 告诉 CodeBuddy Code 在压缩时保留什么内容。

你也可以在 CODEBUDDY.md 中自定义压缩行为：

```
# Compact instructions

When you are using compact, please focus on test output and code changes.
```

### **异步压缩策略**

当对话历史接近上下文限制时，系统会自动进行压缩：

- **自动触发**：当上下文接近限制时自动启动压缩
- **后台执行**：压缩过程在后台异步执行，不阻塞用户操作
- **智能摘要**：保留关键信息，压缩冗余内容
- **无缝衔接**：用户无感知，获得"无限上下文"体验

压缩时会保留以下关键信息：代码变更记录、重要决策点、用户明确的偏好和指示、当前任务的关键上下文。

### **选择合适的模型**

根据任务复杂度选择模型。使用 `/model` 在会话中切换模型，或在 `/config` 中设置默认值。

- **简单任务使用 lite**：文件搜索、快速查询、代码格式化
- **复杂任务使用 reasoning**：架构设计、性能优化、复杂调试
- **一般任务使用 default**：日常编码、功能实现

### **减少 MCP 服务器开销**

每个 MCP 服务器会将工具定义添加到上下文中，即使处于空闲状态。运行 `/mcp` 查看已配置的服务器。

- **优先使用 CLI 工具**：`gh`、`aws`、`gcloud` 等工具比 MCP 服务器更节省上下文，因为它们不会添加持久的工具定义。CodeBuddy Code 可以直接运行 CLI 命令，无需额外开销。
- **禁用未使用的服务器**：运行 `/mcp` 查看并禁用未使用的服务器。

### **将详细操作委托给子代理**

运行测试、获取文档或处理日志文件可能消耗大量上下文。将这些操作委托给子代理，详细输出保留在子代理的上下文中，只有摘要返回主对话。

### **编写精确的提示**

模糊的请求如"改进这个代码库"会触发大范围扫描。精确的请求如"为 auth.ts 中的 login 函数添加输入验证"让 CodeBuddy Code 能够以最少的文件读取高效工作。

### **复杂任务的高效工作方式**

对于较长或更复杂的工作，以下习惯有助于避免因走错方向而浪费 Token：

- **复杂任务使用计划模式**：按 Shift+Tab (macOS/Linux) 或 Alt+M (Windows) 进入计划模式。CodeBuddy Code 会探索代码库并提出方案供你批准，避免初始方向错误时的昂贵返工。
- **尽早纠正方向**：如果 CodeBuddy Code 开始走错方向，按 Escape 立即停止。使用 `/rewind` 或双击 Escape 将对话和代码恢复到之前的检查点。
- **提供验证目标**：在提示中包含测试用例、截图或预期输出。当 CodeBuddy Code 能够自我验证工作时，它可以在你需要请求修复之前发现问题。
- **增量测试**：写一个文件，测试它，然后继续。这样可以在问题还容易修复时尽早发现。

## **后台 Token 消耗**

CodeBuddy Code 在空闲时也会为某些后台功能消耗 Token：

- **对话摘要**：后台任务会为 `--resume` 功能摘要之前的对话
- **Prompt 预测**：根据历史对话信息推测下一条最有可能的 Prompt 输入

这些后台进程即使没有活跃交互也会消耗少量 Token。

## **相关文档**

- **[子代理](https://www.codebuddy.ai/docs/zh/cli/sub-agents)** - 使用子代理隔离高消耗操作
- **[MCP](https://www.codebuddy.ai/docs/zh/cli/mcp)** - 管理 MCP 服务器开销
- **[模型](https://www.codebuddy.ai/docs/zh/cli/models)** - 了解可用的模型选项

***

*本文档帮助您了解如何有效管理 CodeBuddy Code 的使用成本。*

**最后更新: 2026/3/23 03:17**

<br />

# **工具参考**

> CodeBuddy Code 可用工具的完整参考，包括权限要求。

CodeBuddy Code 内置一系列工具来帮助理解和修改代码库。下表中的工具名称即为**[权限规则](https://www.codebuddy.ai/docs/zh/cli/iam#%E5%B7%A5%E5%85%B7%E7%89%B9%E5%AE%9A%E7%9A%84%E6%9D%83%E9%99%90%E8%A7%84%E5%88%99)**、**[子代理工具列表](https://www.codebuddy.ai/docs/zh/cli/sub-agents)**和 **[Hook 匹配器](https://www.codebuddy.ai/docs/zh/cli/hooks)**中使用的标识符。

**工具**

**说明**

**需要权限**

`Agent`

生成具有独立上下文窗口的**[子代理](https://www.codebuddy.ai/docs/zh/cli/sub-agents)**来处理任务

否

`AskUserQuestion`

向用户提出多选问题，收集需求或澄清歧义

是

`Bash`

在你的环境中执行 Shell 命令。参见 **[Bash 工具行为](https://www.codebuddy.ai/docs/zh/cli/tools-reference#bash-%E5%B7%A5%E5%85%B7%E8%A1%8C%E4%B8%BA)**

是

`CronCreate`

在当前会话内调度定时或一次性任务（退出后失效）。参见**[定时任务](https://www.codebuddy.ai/docs/zh/cli/scheduled-tasks)**

否

`CronDelete`

按 ID 取消定时任务

否

`CronList`

列出当前会话中所有定时任务

否

`DeferExecuteTool`

执行通过 `ToolSearch` 发现的延迟加载工具

否

`Edit`

对文件进行精确的字符串替换编辑

是

`EnterPlanMode`

切换到计划模式，在编码前设计实现方案

否

`EnterWorktree`

创建隔离的 **[git worktree](https://www.codebuddy.ai/docs/zh/cli/worktree)** 并切换到其中

是

`ExitPlanMode`

提交计划供用户审批并退出计划模式

是

`Glob`

基于 glob 模式查找文件

否

`Grep`

在文件内容中搜索正则表达式模式

否

`ImageEdit`

基于文本指令编辑或修改已有图片

是

`ImageGen`

根据文本描述生成图片

是

`LeaveWorktree`

退出 worktree 会话并返回原始目录

是

`LSP`

通过语言服务器提供代码智能。文件编辑后自动报告类型错误和警告。还支持跳转定义、查找引用、获取类型信息、列出符号、查找实现、追踪调用层级等导航操作。需要**[代码智能插件](https://www.codebuddy.ai/docs/zh/cli/plugins)**及其语言服务器二进制文件

否

`MultiEdit`

在单个原子操作中对同一文件执行多步编辑

是

`NotebookEdit`

修改 Jupyter notebook 单元格内容

是

`PowerShell`

在 Windows 上执行 PowerShell 命令。仅 Windows 可用，参见 **[PowerShell 工具行为](https://www.codebuddy.ai/docs/zh/cli/tools-reference#powershell-%E5%B7%A5%E5%85%B7%E8%A1%8C%E4%B8%BA)**

是

`Read`

读取文件内容，支持图片、PDF 和 Jupyter notebook

否

`SendMessage`

在 **[Agent 团队](https://www.codebuddy.ai/docs/zh/cli/agent-teams)**中向队友发送消息

否

`Skill`

在主对话中执行 **[Skill 技能](https://www.codebuddy.ai/docs/zh/cli/skills)**

否

`SlashCommand`

执行自定义**[斜杠命令](https://www.codebuddy.ai/docs/zh/cli/slash-commands)**

是

`StructuredOutput`

返回符合 JSON Schema 的结构化输出

否

`TaskCreate`

创建新任务到任务列表

否

`TaskGet`

获取特定任务的完整详情

否

`TaskList`

列出所有任务及其当前状态

否

`TaskOutput`

获取后台任务或子代理的输出

否

`TaskStop`

按 ID 终止正在运行的后台任务

否

`TaskUpdate`

更新任务状态、依赖、详情或删除任务

否

`TeamCreate`

创建 **[Agent 团队](https://www.codebuddy.ai/docs/zh/cli/agent-teams)**以协调多个代理协作

否

`TeamDelete`

删除 Agent 团队及其任务目录

否

`ToolSearch`

搜索并加载延迟加载的工具，支持内置工具和 **[MCP 工具](https://www.codebuddy.ai/docs/zh/cli/mcp#%E5%BB%B6%E8%BF%9F%E5%8A%A0%E8%BD%BD-defer_loading)**

否

`WebFetch`

获取指定 URL 的内容并进行 AI 分析

是

`WebSearch`

执行网络搜索

是

`Write`

创建或覆盖文件

是

权限规则可通过 `/permissions` 命令或在**[权限设置](https://www.codebuddy.ai/docs/zh/cli/settings#%E6%9D%83%E9%99%90%E8%AE%BE%E7%BD%AE)**中配置。另请参阅**[工具特定的权限规则](https://www.codebuddy.ai/docs/zh/cli/iam#%E5%B7%A5%E5%85%B7%E7%89%B9%E5%AE%9A%E7%9A%84%E6%9D%83%E9%99%90%E8%A7%84%E5%88%99)**。

## **工具别名**

部分工具拥有别名，可在权限规则中互换使用：

**工具**

**别名**

`TaskOutput`

`BashOutput`

`TaskStop`

`KillShell`

`PowerShell`

`pwsh`、`ps`

## **Bash 工具行为**

Bash 工具在独立进程中执行每条命令，具有以下持久化特性：

- **工作目录**在命令间保持不变。设置 `CODEBUDDY_BASH_MAINTAIN_PROJECT_WORKING_DIR=1` 可在每次命令后重置到项目目录。
- **环境变量不会持久化**。一条命令中的 `export` 不会在下一条命令中生效。

在启动 CodeBuddy Code 前激活你的 virtualenv 或 conda 环境。要使环境变量在 Bash 命令间持久化，启动前设置 **[`CODEBUDDY_ENV_FILE`](https://www.codebuddy.ai/docs/zh/cli/env-vars)** 指向一个 shell 脚本，或使用 **[SessionStart Hook](https://www.codebuddy.ai/docs/zh/cli/hooks)** 动态填充。

### **沙箱模式**

Bash 工具支持**[沙箱隔离](https://www.codebuddy.ai/docs/zh/cli/bash-sandboxing)**，可限制文件系统和网络访问。启用沙箱时，命令在受限环境中执行，防止未授权的系统访问。

通过 `dangerouslyDisableSandbox` 参数可逐条命令绕过沙箱（需要用户审批）。

### **后台执行**

通过 `run_in_background` 参数可将命令在后台运行，使用 `TaskOutput` 工具读取输出。适用于长时间运行的构建、测试等场景。

## **PowerShell 工具行为**

PowerShell 工具仅在 Windows 上可用，提供原生 PowerShell 命令执行能力。

### **与 Bash 工具的关系**

- **有 Git Bash 时**：Bash 工具和 PowerShell 工具同时可用，模型根据场景选择合适的工具
- **无 Git Bash 时**：Bash 工具自动禁用，PowerShell 工具成为唯一的 shell 工具
- macOS/Linux 上 PowerShell 工具不可用

### **版本适配**

PowerShell 工具自动检测 PowerShell 版本，优先使用 PowerShell 7+（pwsh），其次使用 Windows PowerShell 5.1。prompt 中的语法指导会根据版本差异自动调整（如 `&&` 操作符仅 7+ 支持）。

### **安全检查**

PowerShell 工具内置安全检查器，覆盖代码注入、下载执行、提权操作、系统破坏等危险模式。危险命令（如 `Invoke-Expression`、`Add-Type`）会被阻止，系统修改类命令需要用户确认。

### **环境变量**

**环境变量**

**说明**

`CODEBUDDY_POWERSHELL_PATH`

显式指定 PowerShell 路径（优先于自动检测）

`CODEBUDDY_USE_POWERSHELL_TOOL`

设为 `0` 禁用 PowerShell 工具

## **延迟加载工具**

部分工具（如通过 **[MCP 服务器](https://www.codebuddy.ai/docs/zh/cli/mcp)**提供的工具）采用延迟加载机制。这些工具不会在初始工具列表中出现，需要通过 `ToolSearch` 发现和激活。一旦激活，工具在会话剩余时间内保持可用。

通过产品配置可以将内置工具设置为延迟加载：

```
{
  "tools": {
    "ToolName": {
      "deferLoading": true
    }
  }
}
```

## **另请参阅**

- **[身份和访问管理](https://www.codebuddy.ai/docs/zh/cli/iam)**：权限系统、规则语法和工具特定规则
- **[子代理](https://www.codebuddy.ai/docs/zh/cli/sub-agents)**：为子代理配置工具访问权限
- **[Hooks 钩子系统](https://www.codebuddy.ai/docs/zh/cli/hooks-guide)**：在工具执行前后运行自定义命令
- **[MCP 集成](https://www.codebuddy.ai/docs/zh/cli/mcp)**：通过 MCP 服务器扩展可用工具
- **[定时任务](https://www.codebuddy.ai/docs/zh/cli/scheduled-tasks)**：使用 CronCreate 调度自动化任务
- **[Agent 团队](https://www.codebuddy.ai/docs/zh/cli/agent-teams)**：多代理协作的团队系统

**最后更新: 2026/4/8 20:42**
身份和访问管理
了解如何为组织中的 CodeBuddy Code 配置用户身份验证、授权和访问控制。

认证方法
快速开始
根据你的场景选择认证方式：

场景	推荐方式	如何获取凭据
个人开发者	CODEBUDDY_API_KEY	从平台获取 API Key
企业/团队（OAuth 集成）	apiKeyHelper	创建应用获取 Client ID/Secret
CI/CD 已有 OAuth token	CODEBUDDY_AUTH_TOKEN	直接使用已有 token
第三方模型服务	CODEBUDDY_API_KEY + BASE_URL	从第三方服务商获取
当多个认证方式同时配置时，按以下优先级生效：CODEBUDDY_AUTH_TOKEN > apiKeyHelper > CODEBUDDY_API_KEY

个人用户：获取 API Key
适用于个人开发者快速上手，使用 CodeBuddy 平台提供的模型服务。

第 1 步：获取 API Key

访问对应平台获取你的 API Key：

版本	获取地址
海外版	https://www.codebuddy.ai/profile/keys
中国版	https://copilot.tencent.com/profile/
iOA 版	https://tencent.sso.copilot.tencent.com/profile/keys
第 2 步：配置环境变量

⚠️ 重要提示：必须正确配置 CODEBUDDY_INTERNET_ENVIRONMENT 环境变量！

这是用户最常遗漏的配置项。如果不设置或设置错误，将导致认证失败或连接到错误的服务端点。

根据你使用的版本，配置对应的环境变量：

海外版：


export CODEBUDDY_API_KEY="your-api-key"
# 海外版无需设置 CODEBUDDY_INTERNET_ENVIRONMENT（默认值）
中国版：


export CODEBUDDY_API_KEY="your-api-key"
export CODEBUDDY_INTERNET_ENVIRONMENT=internal
iOA 版：


export CODEBUDDY_API_KEY="your-api-key"
export CODEBUDDY_INTERNET_ENVIRONMENT=ioa
版本	CODEBUDDY_INTERNET_ENVIRONMENT 值	说明
海外版	不设置（或 public）	默认值，连接海外服务
中国版	internal	连接中国区服务
iOA 版	ioa	连接腾讯 iOA 内网服务
💡 持久化配置建议： 将环境变量添加到 ~/.bashrc、~/.zshrc 或 shell 配置文件中，避免每次手动设置。


# 添加到 ~/.zshrc 或 ~/.bashrc
echo 'export CODEBUDDY_API_KEY="your-api-key"' >> ~/.zshrc
echo 'export CODEBUDDY_INTERNET_ENVIRONMENT=internal' >> ~/.zshrc  # 中国版
source ~/.zshrc
第 3 步：开始使用


codebuddy
企业用户：OAuth 认证
适用于企业/团队集成，通过 OAuth 2.0 获取 token。

目前仅介绍 Client Credentials 授权方式，适用于服务端应用和 CI/CD 场景。

前提条件：企业用户需要先购买 CodeBuddy 旗舰版才能使用 OAuth 认证。详见 CodeBuddy 旗舰版购买指南。

第 1 步：创建应用，获取 Client ID 和 Secret

参考 企业开发者快速入门 完成：

创建企业应用
获取 Client ID 和 Client Secret
第 2 步：创建获取 token 的脚本


#!/bin/bash
# get-oauth-token.sh - OAuth 2.0 Client Credentials 流程

CLIENT_ID="${OAUTH_CLIENT_ID}"
CLIENT_SECRET="${OAUTH_CLIENT_SECRET}"
TOKEN_URL="https://copilot.tencent.com/oauth2/token"

response=$(curl -s -X POST "$TOKEN_URL" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials" \
  -d "client_id=$CLIENT_ID" \
  -d "client_secret=$CLIENT_SECRET")

echo "$response" | jq -r '.access_token'
第 3 步：配置 apiKeyHelper

在 ~/.codebuddy/settings.json 或项目 .codebuddy/settings.json 中配置：


{
  "apiKeyHelper": "/path/to/get-oauth-token.sh"
}
配置完成后即可使用 codebuddy 命令。

apiKeyHelper 获取的 token 默认缓存 5 分钟，可通过 CODEBUDDY_CODE_API_KEY_HELPER_TTL_MS 环境变量调整（单位：毫秒）。

第三方模型服务
适用于使用 OpenRouter 等第三方模型服务。

注意： 使用第三方模型服务时，无需设置 CODEBUDDY_INTERNET_ENVIRONMENT 环境变量，因为请求直接发送到第三方服务端点。


export CODEBUDDY_API_KEY="sk-or-v1-xxx"
export CODEBUDDY_BASE_URL="https://openrouter.ai/api/v1"
codebuddy --model openai/gpt-4
认证方式详解
CODEBUDDY_API_KEY
静态 API Key，适用于大部分场景。

特性	说明
环境变量	CODEBUDDY_API_KEY
认证类型	API Key（X-Api-Key 请求头）
适用场景	个人开发、第三方模型服务
⚠️ 注意：使用 API Key 时，必须同时配置 CODEBUDDY_INTERNET_ENVIRONMENT 环境变量！

版本	配置值
海外版	不设置（默认）
中国版	internal
iOA 版	ioa

# 海外版
export CODEBUDDY_API_KEY="your-api-key"

# 中国版
export CODEBUDDY_API_KEY="your-api-key"
export CODEBUDDY_INTERNET_ENVIRONMENT=internal

# iOA 版
export CODEBUDDY_API_KEY="your-api-key"
export CODEBUDDY_INTERNET_ENVIRONMENT=ioa
CODEBUDDY_AUTH_TOKEN
已获取的 OAuth Bearer Token，适用于 CI/CD 或已有 token 的场景。

特性	说明
环境变量	CODEBUDDY_AUTH_TOKEN
认证类型	Bearer Token（Authorization 请求头）
适用场景	CI/CD 自动化、已有 OAuth token

export CODEBUDDY_AUTH_TOKEN="your-oauth-bearer-token"
也可以在 settings.json 中配置：


{
  "env": {
    "CODEBUDDY_AUTH_TOKEN": "your-oauth-token"
  }
}
apiKeyHelper
通过脚本动态获取 token，适用于 OAuth 集成或需要定期刷新 token 的场景。

特性	说明
配置方式	settings.json 的 apiKeyHelper 字段
认证类型	Bearer Token（由脚本返回）
缓存机制	默认 5 分钟，可通过 CODEBUDDY_CODE_API_KEY_HELPER_TTL_MS 配置
适用场景	OAuth Client Credentials、Vault 集成、token 自动刷新
脚本要求：

将 token 输出到标准输出（stdout）
以退出码 0 表示成功
脚本执行超时时间为 30 秒
配置示例：


{
  "apiKeyHelper": "/path/to/get-token.sh"
}
从 Vault 获取 token 示例：


#!/bin/bash
vault read -field=api_key secret/codebuddy/api-key
详细的认证配置请参考设置文档 - 环境变量。

访问控制和权限
我们支持细粒度权限，以便您可以精确指定代理允许执行的操作（例如运行测试、运行 linter)和不允许执行的操作（例如更新云基础设施）。这些权限设置可以检入版本控制并分发给组织中的所有开发人员，也可以由各个开发人员自定义。

权限系统
CodeBuddy Code 使用分层权限系统来平衡功能和安全性：

工具类型	示例	需要批准	"是，不再询问"行为
只读	文件读取、LS、Grep	否	N/A
Bash 命令	Shell 执行	是	按项目目录和命令永久记住
文件修改	编辑/写入文件	是	会话结束前有效
配置权限
您可以使用 /permissions 查看和管理 CodeBuddy Code 的工具权限。此 UI 列出所有权限规则及其来源的 settings.json 文件。

Allow 规则将允许 CodeBuddy Code 使用指定的工具，无需进一步手动批准。
Ask 规则将在 CodeBuddy Code 尝试使用指定工具时询问用户确认。Ask 规则优先于 allow 规则。
Deny 规则将阻止 CodeBuddy Code 使用指定工具。Deny 规则优先于 allow 和 ask 规则。
Additional directories 将 CodeBuddy 的文件访问扩展到初始工作目录之外的目录。
Default mode 控制 CodeBuddy 在遇到新请求时的权限行为。
权限规则使用格式： Tool 或 Tool(optional-specifier)

仅工具名称的规则匹配该工具的任何使用。例如，将 Bash 添加到 allow 规则列表将允许 CodeBuddy Code 使用 Bash 工具而无需用户批准。

权限模式
CodeBuddy Code 支持几种权限模式，可以在设置文件中设置为 defaultMode:

模式	描述
default	标准行为 - 首次使用每个工具时提示权限
acceptEdits	自动接受会话的文件编辑权限
plan	计划模式 - CodeBuddy 可以分析但不能修改文件或执行命令
bypassPermissions	跳过所有权限提示（需要安全环境 - 见下方警告）
WARNING

`bypassPermissions` 模式应仅在安全、隔离的环境中使用，例如 Docker 容器或 VM。在生产环境或包含敏感数据的系统上使用此模式可能会带来安全风险。
NOTE

**`trustAll` / `trustedDirectories` 不是权限模式替代项。** 这两个字段只影响启动时的**目录信任授权提示**（即"是否信任此目录并允许在其中运行 CodeBuddy"的一次性弹窗），和工具执行时是否弹审批无关。
免除工具审批请设置 permissions.defaultMode 或使用 --permission-mode bypassPermissions / -y / --dangerously-skip-permissions
免除目录信任弹窗才用 trustAll: true 或把目录加入 trustedDirectories
两个开关独立；bypassPermissions 开着时目录信任弹窗仍会正常出现，反之亦然
所以"开了 bypassPermissions + trustAll，其他就不用配"这个说法是不准确的——前者控制工具审批，后者控制目录信任，两者解决的是不同的确认入口。

工作目录
默认情况下，CodeBuddy 可以访问其启动目录中的文件。您可以扩展此访问权限：

启动时：使用 --add-dir <path> CLI 参数
会话期间：使用 /add-dir 斜杠命令
持久配置：添加到设置文件的 additionalDirectories
附加目录中的文件遵循与原始工作目录相同的权限规则 - 它们可以无提示读取，文件编辑权限遵循当前权限模式。

工具特定的权限规则
某些工具支持更细粒度的权限控制：

Bash

Bash(npm run build) 精确匹配 Bash 命令 npm run build
Bash(npm run test:*) 匹配以 npm run test 开头的 Bash 命令
Bash(curl http://site.com/:*) 匹配以 curl http://site.com/ 开头的 curl 命令
CodeBuddy Code 能识别 shell 操作符（如 `&&`),因此前缀匹配规则如 `Bash(safe-cmd:*)` 不会授予它运行命令 `safe-cmd && other-cmd` 的权限
WARNING

Bash 权限模式的重要限制：
此工具使用前缀匹配,而非正则表达式或 glob 模式
通配符 :* 仅在模式末尾有效，用于匹配任何后续内容
像 Bash(curl http://github.com/:*) 这样的模式可以通过多种方式绕过:
URL 前的选项: curl -X GET http://github.com/... 不匹配
不同协议: curl https://github.com/... 不匹配
重定向: curl -L http://bit.ly/xyz (重定向到 github)
变量: URL=http://github.com && curl $URL 不匹配
额外空格: curl http://github.com 不匹配
要更可靠地过滤 URL,请考虑：

使用带 WebFetch(domain:github.com) 权限的 WebFetch 工具
通过 CODEBUDDY.md 指示 CodeBuddy Code 您允许的 curl 模式
使用 hooks 进行自定义权限验证
Read & Edit

Edit 规则适用于所有编辑文件的内置工具。CodeBuddy 将尽力将 Read 规则应用于所有读取文件的内置工具，如 Grep、Glob 和 LS。

Read 和 Edit 规则都遵循 gitignore 规范,有四种不同的模式类型:

模式	含义	示例	匹配
//path	从文件系统根目录的绝对路径	Read(//Users/alice/secrets/**)	/Users/alice/secrets/**
~/path	从家目录的路径	Read(~/Documents/*.pdf)	/Users/alice/Documents/*.pdf
/path	相对于设置文件的路径	Edit(/src/**/*.ts)	<设置文件路径>/src/**/*.ts
path 或 ./path	相对于当前目录的路径	Read(*.env)	<cwd>/*.env
WARNING

像 `/Users/alice/file` 这样的模式不是绝对路径 - 它相对于您的设置文件！使用 `//Users/alice/file` 表示绝对路径。
示例：

Edit(/docs/**) - 在 <项目>/docs/ 中编辑（不是 /docs/!)
Read(~/.zshrc) - 读取家目录的 .zshrc
Edit(//tmp/scratch.txt) - 编辑绝对路径 /tmp/scratch.txt
Read(src/**) - 从 <当前目录>/src/ 读取
WebFetch

WebFetch(domain:example.com) 匹配对 example.com 的获取请求
MCP

mcp__puppeteer 匹配 puppeteer 服务器提供的任何工具（在 CodeBuddy Code 中配置的名称）
mcp__puppeteer__* 通配符语法,也匹配 puppeteer 服务器的所有工具
mcp__puppeteer__puppeteer_navigate 匹配 puppeteer 服务器提供的 puppeteer_navigate 工具
要批准 MCP 服务器的所有工具，可以使用以下任一格式：
✅ 使用： mcp__github （批准所有 GitHub 工具）
✅ 使用： mcp__github__* （批准所有 GitHub 工具，与上一条等价）
要仅批准特定工具，列出每一个：

✅ 使用： mcp__github__get_issue
✅ 使用： mcp__github__list_issues
权限配置示例
基础权限配置：


{
  "permissions": {
    "allow": [
      "Read",
      "Edit",
      "Bash(git:*)",
      "Bash(npm:*)"
    ],
    "ask": [
      "WebFetch",
      "Bash(docker:*)"
    ],
    "deny": [
      "Bash(rm:*)",
      "Bash(sudo:*)",
      "Edit(**/*.env)",
      "Read(~/.ssh/**)"
    ]
  }
}
安全限制配置：


{
  "permissions": {
    "allow": [
      "Read",
      "Edit(src/**)",
      "Bash(git:status,git:diff)"
    ],
    "deny": [
      "Edit(**/*.env)",
      "Edit(**/*.key)",
      "Edit(**/*.pem)",
      "Bash(wget:*)",
      "Bash(curl:*)",
      "Read(/etc/**)",
      "Read(~/.ssh/**)",
      "Read(~/.aws/**)"
    ],
    "defaultMode": "default"
  }
}
使用 hooks 进行额外的权限控制
CodeBuddy Code hooks 提供了一种注册自定义 shell 命令在运行时执行权限评估的方法。当 CodeBuddy Code 进行工具调用时，PreToolUse hooks 在权限系统运行之前运行，hook 输出可以决定是批准还是拒绝工具调用，以代替权限系统。

详见 Hooks 文档。

设置优先级
当存在多个设置源时，它们按以下顺序应用（从高到低优先级）:

命令行参数
本地项目设置(.codebuddy/settings.local.json)
共享项目设置(.codebuddy/settings.json)
用户设置(~/.codebuddy/settings.json)
此层次结构确保在适当的地方仍允许项目和用户级别的灵活性。

凭据管理
CodeBuddy Code 安全地管理您的身份验证凭据：

平台	存储位置
macOS	加密的 macOS Keychain
Linux	系统密钥环（GNOME Keyring、KWallet）
Windows	Windows 凭据管理器
动态凭据获取： 使用 apiKeyHelper 配置自定义脚本动态获取 token，详见 apiKeyHelper 配置。

安全最佳实践
1. 最小权限原则
仅授予 CodeBuddy Code 完成任务所需的最小权限：


{
  "permissions": {
    "allow": [
      "Read",
      "Edit(src/**/*.ts)",
      "Bash(npm:test,npm:build)"
    ],
    "deny": [
      "Edit(**/*.env)",
      "Bash(rm:*)",
      "Bash(sudo:*)"
    ]
  }
}
2. 保护敏感文件
始终拒绝访问包含敏感信息的文件：


{
  "permissions": {
    "deny": [
      "Read(.env)",
      "Read(.env.*)",
      "Read(secrets/**)",
      "Read(~/.ssh/**)",
      "Read(~/.aws/**)",
      "Edit(**/*.key)",
      "Edit(**/*.pem)"
    ]
  }
}
3. 使用 Bash 沙箱
在支持的平台上启用沙箱以隔离 bash 命令：


{
  "sandbox": {
    "enabled": true,
    "autoAllowBashIfSandboxed": true,
    "excludedCommands": ["docker"]
  }
}
详见Bash 沙箱文档。

4. 审查权限日志
定期检查 CodeBuddy Code 的权限使用情况，确保符合安全策略。

5. 团队配置共享
将团队级别的权限配置检入版本控制：


# 创建团队共享配置
.codebuddy/settings.json

# 添加到 .gitignore
.codebuddy/settings.local.json
常见问题
使用 API Key 认证失败怎么办？
这是最常见的问题，通常是因为 未配置或错误配置了 CODEBUDDY_INTERNET_ENVIRONMENT 环境变量。

排查步骤：

确认你使用的版本（海外版/中国版/iOA 版）
检查环境变量是否正确设置：

# 检查当前配置
echo $CODEBUDDY_API_KEY
echo $CODEBUDDY_INTERNET_ENVIRONMENT
根据版本设置正确的环境变量：
版本	CODEBUDDY_INTERNET_ENVIRONMENT
海外版	不设置或 public
中国版	internal
iOA 版	ioa
常见错误：

❌ 中国版用户忘记设置 CODEBUDDY_INTERNET_ENVIRONMENT=internal
❌ iOA 版用户设置成了 internal 而不是 ioa
❌ 环境变量只在当前终端生效，新开终端后失效（建议添加到 shell 配置文件）
如何临时绕过权限？
使用 --permission-mode bypassPermissions 启动 CodeBuddy Code:


codebuddy --permission-mode bypassPermissions
WARNING

仅在安全、隔离的环境中使用此选项！
如何为特定项目设置不同的权限？
在项目根目录创建 .codebuddy/settings.json:


{
  "permissions": {
    "allow": ["项目特定的权限"]
  }
}
如何查看当前权限配置？
使用 /permissions 命令查看所有生效的权限规则及其来源。

另见
设置配置 - 了解完整的配置选项
Hooks 文档 - 使用 hooks 进行高级权限控制
Bash 沙箱 - 了解沙箱隔离功能
MCP 集成 - 配置 MCP 服务器权限
通过适当的权限配置，确保 CodeBuddy Code 在安全边界内工作

最后更新: 2026/4/27 08:00
