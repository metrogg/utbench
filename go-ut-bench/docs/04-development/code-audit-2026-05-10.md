# 基于代码的偏差审计（2026-05-10）

本文记录一次以当前代码为准的仓库审计。重点不是历史方案，而是列出仍然会影响开发、使用或理解成本的真实偏差。

## 审计范围

- `go-ut-bench/cmd/utbench/main.go`
- `go-ut-bench/internal/**`
- `go-ut-bench/docs/**`
- 父仓库根目录 `README.md`、`docs/**`、`go.mod`

## 结论摘要

当前仓库的主要问题不在“核心架构不可用”，而在：

1. 代码和帮助文案之间仍有历史用词残留。
2. `go-ut-bench/docs/` 与父仓库 `docs/` 并存，容易让新成员进入错误文档树。
3. 一部分默认值和文档口径没有跟上当前主实现。
4. 测试在当前沙箱环境下无法直接验证，存在环境权限噪音。

## P1：代码与 CLI 帮助文案不一致

### 1. `--class` 帮助文案仍写 `module_level`

当前代码契约只接受 `self_contained` 和 `repo_level`：

- `internal/contracts/constants.go`
- `internal/dataset/service.go`

但 `run` 子命令帮助文案仍然写成：

```text
Dataset class(es), comma-separated (self_contained, module_level)
```

定位：

- `cmd/utbench/main.go:570`
- `internal/dataset/service.go:64`
- `internal/dataset/service.go:68`
- `internal/contracts/constants.go:18`

影响：

- 用户按帮助文案传 `module_level` 会直接触发参数校验失败。
- 新成员会误以为 `module_level` 仍是当前有效契约。

建议：

- 统一改为 `repo_level`。
- 补一个兼容层或更明确的错误提示也可以，但至少帮助文案不能继续误导。

### 2. CLI Quick Start 存在命令拼写错误

帮助文案里的 Quick Start 第 2 步写成了：

```text
.\tbench doctor --langs python
```

应为：

```text
.\utbench doctor --langs python
```

定位：

- `cmd/utbench/main.go:97`

影响：

- 新用户复制即失败。
- 会直接降低 CLI 首次使用成功率。

建议：

- 修正文案。

## P1：代码默认值与当前仓库布局不一致

### 3. `run/generate` 默认 `--config` 仍指向历史路径

`run` 和 `generate` 默认模型配置路径仍是：

```text
../benchmark/config/models.yaml
```

定位：

- `cmd/utbench/main.go:564`
- `cmd/utbench/main.go:676`

而当前仓库内的真实配置路径是：

- `go-ut-bench/configs/models.yaml`

并且 `web` 子命令默认值已经是当前路径：

- `cmd/utbench/main.go:424`

影响：

- 用户如果不显式传 `--config`，CLI 行为和 Web 行为不一致。
- 在当前仓库里直接跑 `utbench run`，默认值很可能找不到配置。

建议：

- 把 CLI 默认值统一改到 `./configs/models.yaml`。
- 如果必须保留历史兼容，可以做“当前路径不存在时回退旧路径”。

## P1：文档主入口与代码真实定位存在偏差

### 4. `project-introduction.md` 仍把项目描述成 LLM benchmark，弱化了 Agent/subject 平台化事实

`docs/00-overview/project-introduction.md` 当前开头仍然把项目定义为：

```text
LLM 单元测试生成能力横向评测基准工具
```

定位：

- `go-ut-bench/docs/00-overview/project-introduction.md:7`

但从代码看，系统已经明确支持：

- `subject`
- `framework`
- `skill`
- `cli_agent`
- trace / workspace diff / sandbox fingerprint
- Web 里的 agent/skill 管理 API

对应代码：

- `internal/agentconfig/config.go`
- `internal/runner/subjects.go`
- `internal/contracts/spec.go`
- `internal/contracts/results.go`
- `internal/web/server.go`

影响：

- 文档会让新成员低估 `subject` 抽象的重要性。
- 架构理解容易停留在“模型跑基准”而不是“实验对象平台化”。

建议：

- 把 overview 主文档改成“模型 API + Agent + skill 的统一评测平台”。

### 5. 文件职责文档出现过旧 prompt 模式名，本次已修正

`project-file-map.md` 之前把 prompt 模式写成了：

```text
full_file、completion、module_level
```

而代码中的真实模式是：

- `full_file`
- `completion`
- `repo_level`

对应代码：

- `internal/runner/prompt.go:24`
- `internal/runner/prompt.go:27`

处理结果：

- 本次整理已将 `project-file-map.md` 同步为 `repo_level`。

剩余建议：

- 继续排查其他 overview / user guide 文档中是否还有零散旧词。

## P2：双文档树并存，认知成本高

### 6. 父仓库 `docs/` 和 `go-ut-bench/docs/` 同时存在

目前仓库中有两套文档树：

- 父仓库根目录 `docs/`
- 主模块 `go-ut-bench/docs/`

当前根 README 已经说明以 `go-ut-bench/docs` 为准：

- `README.md:3`

但父仓库 `docs/` 里仍有大量“项目分析、调研、历史讨论、外部资料”，容易被误当主文档。

影响：

- 搜文档时极易命中旧结论。
- 新成员不知道哪些文档是“背景资料”，哪些是“当前系统说明”。

建议：

- 继续保留双树可以，但需要更显式的入口标识。
- 在父仓库 `docs/README.md` 顶部增加“非主实现文档树”的醒目标识。

### 7. 面试讲解稿和新成员导读混在同一层级

当前 `go-ut-bench/docs/` 根目录已有：

- `interview-project-walkthrough.md`

这份文档定位明确是“面试讲解稿”，但它和 overview 文档处于近似入口位置，容易被当主导读。

影响：

- 读者可能先读到“面试表述”，而不是“代码事实导读”。

建议：

- 保留这份文档，但为它补充导航标签。
- 新成员优先读 overview 导读，而不是面试稿。

## P2：帮助文案和实际命令面存在信息缺口

### 8. 顶层 usage 没有把 `assets` 的价值讲清楚

`main.go` 的帮助中列出了 `assets`，但旧文档和大多数概览文档并没有把它当第一层能力来讲。

定位：

- `cmd/utbench/main.go:37`
- `cmd/utbench/main.go:89`

影响：

- 用户不知道现在已经能查 subject 资产、generation 资产、evaluation 资产和 reuse 命中解释。

建议：

- 在 overview 和 user guide 里补一段 `assets` 子命令用途。

## P2：代码内部仍有历史兼容痕迹

### 9. `docker-image` 仍保留为废弃别名

Web 子命令里保留了：

```go
imageName := fs.String("docker-image", "", "Deprecated alias for --docker-eval-image")
```

定位：

- `cmd/utbench/main.go:429`

这不是 bug，但说明接口层仍有历史兼容层。

建议：

- 如果确认外部没有依赖，可以在未来版本删除。
- 如果继续保留，应在用户文档里明确只推荐 `--docker-eval-image`。

### 10. `utbench ingest` 被废弃，但仍作为错误入口保留

定位：

- `cmd/utbench/main.go:57`

这说明项目命令面经历过演化，当前建议路径是：

- `utbench db ingest-evaluation`
- `utbench db ingest-manifest`
- `utbench db ingest-report`
- `utbench db ingest-run`

建议：

- 文档里应统一只讲新命令面，不再出现旧命令。

## P2：测试验证当前被环境权限干扰

### 11. 关键包测试在当前环境里卡在 Go build cache 权限

本次尝试运行：

```text
go test ./internal/runner ./internal/reporter ./internal/evaluator
```

结果不是断言失败，而是：

- `C:\Users\wzd\AppData\Local\go-build\... Access is denied`

影响：

- 当前环境下很难快速判断“测试红是代码问题还是本机权限问题”。

建议：

- 在开发文档里补充 Windows / sandbox 下 Go build cache 的建议配置。
- 或统一在 `go-ut-bench` 下提供项目内 `GOCACHE` 的推荐启动方式。

## P3：命名和分类还存在一些“旧世界词汇”

### 12. 顶层 README 仍同时提到 `repo_level / module_level`

定位：

- `README.md:362`

虽然这里是在说明历史兼容，但对第一次接触项目的人来说，仍然会增加歧义。

建议：

- 改成“历史上曾使用 `module_level`，当前代码统一使用 `repo_level`”。

## 推荐处理顺序

建议先按下面顺序处理：

1. 修 CLI 帮助文案中的 `module_level`
2. 修 Quick Start 的 `.\tbench`
3. 统一 `run/generate` 默认 `--config`
4. 修 `project-file-map.md` 的 prompt 模式名
5. 重写 `project-introduction.md` 开头定位
6. 给双文档树加更强入口标识

## 本次新增关联文档

- 新成员导读：`docs/00-overview/new-member-code-walkthrough.md`

## 审计说明

本文只记录本次确认过的问题，不等同于完整缺陷列表。它主要覆盖：

- 代码与文档偏差
- CLI 与实际契约偏差
- 主入口认知偏差
- 环境验证噪音
