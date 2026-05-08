# Prompt 模块设计参考

> 适用模块：`internal/runner/prompt.go`
> 灵感来源：[facebookresearch/testgeneval](https://github.com/facebookresearch/testgeneval) 的 `inference/configs/*` 提示词设计
> 目标：生成**可执行、与源码语义严格一致**的单元测试，最小化后处理成本

---

## 1. 设计原则

| 原则 | 做法 |
|------|------|
| **Transport / Prompt 解耦** | `api.go` 只做 HTTP；`prompt.go` 只负责构造提示词 |
| **多模式派发** | `PromptMode` 枚举：`fullfile` / `completion` / `modulelevel` |
| **结构化上下文注入** | 从源码自动提取 dependencies / critical conditions / module symbols |
| **严格输出契约** | 明确要求 "raw code only, no fences"，降低下游解析复杂度 |
| **语言策略模式** | 每种语言一组规则（Python / Go / Java / C++ / JS），由 `buildLanguageSpecificRules` 分派 |
| **双语清单** | 关键约束中英双语呈现，降低模型"理解偏差" |

---

## 2. 模块结构

```
internal/runner/
├── api.go           # HTTP 调用 + 输出清洗 + 源码静态分析
└── prompt.go        # 提示词构造（本文档重点）
```

### `prompt.go` 关键导出

```go
type PromptMode string
const (
    PromptModeFullFile    PromptMode = "fullfile"
    PromptModeCompletion  PromptMode = "completion"
    PromptModeRepoLevel PromptMode = "modulelevel"
)

type PromptRequest struct {
    Mode            PromptMode
    Language        string
    SamplePath      string
    SourceCode      string
    ExistingTestSrc string                    // 仅 completion 模式
    ModuleMeta      *repoLevelMetaForRunner // 仅 modulelevel 模式
}

func BuildPrompt(req PromptRequest) string  // 公共入口
func buildPrompt(lang, samplePath, src string) string // 兼容旧调用
```

---

## 3. 三种 Prompt 模式

### 3.1 `fullfile`（默认）— 整文件生成

**用途**：给定一段被测源码，从零生成完整测试文件。

**Prompt 骨架**：

```
You are an expert unit testing engineer.
你是一名资深单元测试工程师。

## Role & Objective（角色与目标）
## Language & Framework（语言与框架）
## Step-by-Step Workflow（分步流程）   <- 5 步思考链
## Test Requirements（测试要求）        <- 通用 + 语言特定
## Semantic Alignment Hard Rules        <- 语义对齐硬约束
## Critical Conditions Extracted        <- 从源码提取的分支条件
## Error Prevention Checklist           <- 内部自检清单
## Context Information                  <- sample_id / scenario / deps
## Output Format                        <- 严格输出契约
## Source Code Under Test               <- 被测源码（围栏包裹）
```

**亮点**：
- **Step-by-Step Workflow**：显式 CoT 链（识别符号 → 构造测试矩阵 → 导出期望值 → 编写 → 自检）
- **Critical Conditions**：通过正则扫描源码抽取所有含比较符（`==`, `<=`, `>=`…）的 `if / while / for / return` 行，注入到 prompt，引导模型覆盖关键分支
- **Module-Level Symbols**：提取 `func / type / const / var` 定义，明确告知模型"哪些符号真实存在"，降低幻觉

### 3.2 `completion`（新增）— 增量追加测试

**灵感**：testgeneval 的 `PROMPT_COMPLETION`（first / last / extra 三变体）

**用途**：在已有测试文件末尾追加一个新 test，专注提升覆盖率。

**Prompt 骨架**：

```
## Objective（目标）：只写下一个测试函数
## Language & Framework
## Rules                     <- 不重复已有 import/setup/tests；选未覆盖路径
## Source Code Under Test
## Existing Test File        <- 已有测试内容
## Output Format             <- 只输出下一个函数
```

**对比 testgeneval 的 4 变体**：

| testgeneval | go-ut-bench | 说明 |
|-------------|-------------|------|
| `full` | `PromptModeFullFile` | 完整文件生成 |
| `first` | `PromptModeCompletion` + `ExistingTestSrc=preamble` | 有 import/setup，写第一个 test |
| `last` | `PromptModeCompletion` + `ExistingTestSrc=file_minus_last_test` | 补最后一个 test |
| `extra` | `PromptModeCompletion` + `ExistingTestSrc=full_file` | 追加额外 test |

> **后续扩展点**：需要 AST 切片（类似 testgeneval 的 `get_prompt_contexts.py`）才能真正区分 first / last / extra。当前先提供 `ExistingTestSrc` 字段，调用方自己切。

### 3.3 `modulelevel` — 多文件包生成

**用途**：被测代码位于一个多文件 package 中，配套 `meta.json` 描述 `module_import / package_name / target_file`。

**触发条件**：与样本同目录存在 `meta.json` 或 `<name>.meta.json`，并包含 `module_import` 字段。`loadRepoLevelMetaForRunner()` 自动识别。

**关键差异**：
- 显式告知完整模块导入路径（`from pkg.submod import ...`）
- 禁止相对导入（`from ._common import ...`）
- 提示测试文件会被放到 `tests/` 子目录

---

## 4. 核心子模块

### 4.1 `buildLanguageSpecificRules(lang, moduleName)`

按语言分派规则片段。**已覆盖**：Python / Go / Java / C++ / JavaScript。

| 语言 | 关键约束 |
|------|----------|
| Python | `def test_*`、function-based、`pytest.raises`、本地模块导入 |
| Go | `func TestXxx(t *testing.T)`、**package 必须与源一致**、表格驱动、仅 stdlib、可测未导出函数 |
| Java | `ClassNameTest`、`@Test`、`public void`、JUnit4 导入、大括号平衡自检 |
| C++ | GoogleTest 宏 `TEST() / EXPECT_EQ / ASSERT_EQ` |
| JS | Jest `test() / it()`、`expect().toBe()` |

### 4.2 `buildSemanticAlignmentRules(lang)`

**语义对齐硬约束**（防幻觉核心）：

通用：
- 期望值只能源自展示的源码
- 不明确时测实现分支，不测"理想分支"
- 禁止测试未出现的符号

语言特化：
- Python：异常类型精确匹配（TypeError vs ValueError）
- Go：优先 `errors.Is` / 源码中的 sentinel；尊重值 / 指针接收器
- Java：匹配声明的 checked 异常，不泛化

### 4.3 `buildErrorPreventionChecklist(lang)`

**模型内部自检清单**（Checklist 风格），例如：

```
- [ ] All imports resolve to symbols that exist in the provided source.
- [ ] No placeholder tests like `assert True`.
- [ ] No reliance on real network / files / env variables.
- [ ] Every test has at least one meaningful assertion.
- [ ] `func TestXxx(t *testing.T)` signature is correct; package matches source.   (Go)
```

### 4.4 源码静态分析（`api.go`）

| 函数 | 产出 | 用途 |
|------|------|------|
| `extractDependencies` | import / include 列表 | 提示模型可用依赖 |
| `extractCriticalConditions` | 含比较符的关键分支行 | 引导覆盖 |
| `extractRepoLevelSymbols` | 顶层 func / type / const / var 符号 | 告知真实存在的 API |
| `mockRequirement` | 是否需要 mock | 根据源码关键词启发式判断（http/file/db/...） |
| `parseSampleMeta` | scenario / complexity | 从 sampleID 推断元信息 |

---

## 5. Chat Payload 组装（`buildPayload`）

统一 OpenAI-Chat 格式：

```json
{
  "model": "<model>",
  "messages": [
    {"role": "system", "content": "<systemMessage>"},
    {"role": "user",   "content": "<BuildPrompt(...)>"}
  ],
  "stream": false,
  "<params from models.yaml>": "..."
}
```

**DashScope 原生模式**（非 compatible-mode）自动转换为 `input.messages` + `parameters` 结构，由 `resolveEndpoint` + `buildPayload` 处理。

---

## 6. 输出后处理链（`api.go`）

模型原始响应 → 测试代码，经过 4 步清洗：

```
LLM raw output
 │
 ├─► sanitizeModelOutput       去除 <think>/<analysis> 块、Python 行前噪声
 ├─► extractCode               匹配 ```lang``` 围栏，取对应语言代码块
 ├─► stripMarkdownFence        兜底去除孤立围栏
 └─► validateGeneratedTest     校验非空、无泄漏推理标签、Python 必含 test_/pytest
```

失败分类（`contracts.ErrorInfo.Kind`）：
- `auth_config_error` / `payload_error` / `request_build_error` —— 不可重试
- `network_error` / `timeout` / `http_error`(5xx/429/408) —— 可重试，指数退避 + 抖动
- `response_parse_error` / `response_extract_error` / `quality_error` —— 不可重试

---

## 7. 使用示例

### 7.1 默认调用（runner 内部）

```go
prompt := buildPrompt("go", "benchmark/samples/go/boundary_001.go", src)
// 自动探测 meta.json → 派发到 fullfile 或 modulelevel
```

### 7.2 主动选择 completion 模式

```go
prompt := BuildPrompt(PromptRequest{
    Mode:            PromptModeCompletion,
    Language:        "python",
    SamplePath:      "/path/to/sample.py",
    SourceCode:      srcCode,
    ExistingTestSrc: existingTests, // 已有测试片段
})
```

### 7.3 显式 repo-level

```go
meta := loadRepoLevelMetaForRunner(samplePath)
prompt := BuildPrompt(PromptRequest{
    Mode:       PromptModeRepoLevel,
    Language:   "python",
    SamplePath: samplePath,
    SourceCode: entrySrc,
    ModuleMeta: meta,
})
```

---

## 8. 与 testgeneval 的对照总结

| 维度 | testgeneval | go-ut-bench |
|------|-------------|-------------|
| 组织方式 | 基类 `InstructPrompt` + 模型子类（Llama3/Codestral/Gemma2） | 单文件 `prompt.go` + 语言策略函数 |
| Prompt 变体 | 4 个（full/first/last/extra） | 3 个（fullfile/completion/modulelevel） |
| 上下文切片 | AST-based（`get_prompt_contexts.py`） | 正则抽取（dependencies/conditions/symbols） |
| 输出格式 | ` ```python ``` ` 围栏 | **raw code 无围栏**（更严格） |
| 双语 | 仅英文 | **中英双语**（降低中文模型理解偏差） |
| 硬约束 | 基础 system message + 用户提示 | System + Semantic Alignment + Error Checklist **三层约束** |

---

## 9. 未来扩展路线

1. **AST 切片**（参考 testgeneval）：为 `completion` 模式提供 `preamble / first / last_minus_one / last` 四种切分方式，真实对齐 testgeneval 的 4 变体。
2. **Chat template 适配**：针对 local HuggingFace 模型增加 `tokenizer.apply_chat_template` 对应的包装层。
3. **Few-shot**：在 `Context Information` 段注入同类样本的人写测试作为 demo。
4. **Prompt A/B**：把 `PromptMode` 提升为 CLI 参数（`utbench generate --prompt-mode=...`），便于对比评测。
5. **多 system message**：像 testgeneval 那样区分 `SYSTEM_MESSAGE`（完成型）和 `SYSTEM_MESSAGE_FULL`（生成型）。

---

## 10. 参考文件（代码位置）

- `@/f:/Code/ut-bench/ut-bench-feat_go/ut-bench/go-ut-bench/internal/runner/prompt.go`
- `@/f:/Code/ut-bench/ut-bench-feat_go/ut-bench/go-ut-bench/internal/runner/api.go`
- `@/f:/Code/ut-bench/ut-bench-feat_go/ut-bench/go-ut-bench/internal/runner/api_test.go`

