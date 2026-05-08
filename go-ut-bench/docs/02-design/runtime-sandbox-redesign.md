# UT-Bench 运行时 / Docker / 沙箱重构设计

## 背景

UT-Bench 在加入 Agent 评测后，运行时已经不再是“一个 Docker 镜像 + 一条命令”这么简单。当前系统实际上同时承担了三类职责：

1. **控制面（control plane）**
   - CLI / Web
   - 任务编排
   - 资产复用
   - 报告生成
   - SQLite 入库
2. **评测执行面（evaluation plane）**
   - compile / test / coverage / mutation
   - 语言工具链与评测依赖
   - 超时与进程清理
3. **Agent 沙箱执行面（agent sandbox plane）**
   - 运行 OpenCode / 未来 Claude Code / 其他 CLI Agent
   - 隔离 agent 对 workspace 的读写
   - 采集 trace / usage / diff

现在的主要问题不是“有没有分 Docker”，而是**这三层虽然客观上已经部分分开了，但设计上还没有被系统化表达出来**。

例如：

- 外层 `utbench:latest` 同时承担控制面和评测执行面。
- 内层 `utbench-agent-opencode-*` 已经是 Agent 沙箱镜像，但配置协议仍然主要围绕 `docker_image(s)` 这类旧字段展开。
- Web、本地后端、DOOD、`UTBENCH_SANDBOX_HOST_OUTPUT_ROOT`、`docker.sock` 这些约束散落在多处逻辑里。

结果是：

- nested Docker 容易出错；
- 用户难以理解“哪个镜像是干什么的”；
- 未来接入 repo-level 任务、远端 runtime、E2B、gVisor、Firecracker 时改造成本会很高。

## 业界做法

### 1. SWE-bench：把评测 harness 当成一等公民

SWE-bench 的核心不是网页，而是 **task contract + evaluation harness**。官方 README 直接强调其 Docker 化评测 harness，用固定环境复现实验结果。[SWE-bench README](https://github.com/SWE-bench/SWE-bench/blob/main/README.md)

值得学习的点：

- 评测环境不是顺手写的脚本，而是正式的 execution contract。
- benchmark 的可信度来自“可复现环境 + 固定任务契约”，不是来自界面。

### 2. OpenHands：明确区分 runtime 类型

OpenHands 官方文档把 runtime 分成 Docker、Local、Remote，并且把 runtime/sandbox 视为独立架构层，而不是 Agent 逻辑的一部分。[Runtime Overview](https://docs.all-hands.dev/openhands/usage/runtimes/overview) / [Runtime Architecture](https://docs.all-hands.dev/openhands/usage/architecture/runtime)

值得学习的点：

- runtime 是独立抽象；
- sandbox image 可定制；
- 运行时依赖、挂载、网络策略、镜像选择不应该埋在业务逻辑里。

### 3. E2B：把沙箱模板和版本化做成产品能力

E2B 不是拿一个通用 Docker 镜像硬跑，而是把 sandbox template、base image、template tags/versioning 做成正式能力。[Base Image](https://e2b.dev/docs/template/base-image) / [Tags & Versioning](https://e2b.dev/docs/template/tags)

值得学习的点：

- “沙箱环境”应该可版本化；
- 环境变更应该产生新 fingerprint，而不是悄悄污染旧结果；
- sandbox 不等于 evaluator，也不等于 control plane。

### 4. MLE-bench / Terminal-Bench：执行环境和评分环境是明确分层的

MLE-bench 有自己的 grader、rule violation detector、submission packaging；Terminal-Bench 则围绕真实终端任务设计 execution harness。[MLE-bench](https://github.com/openai/mle-bench) / [Terminal-Bench](https://github.com/laude-institute/terminal-bench)

值得学习的点：

- 执行环境、评分环境、违规检测最好明确分层；
- 失败类型必须区分：能力失败、环境失败、规则失败。

## 对当前 UT-Bench 的判断

### 当前实际拓扑

```text
宿主机
├── ./utbench web / ./utbench run    ← 控制面（Web / CLI）直接运行在宿主机
│   └── generation in-process        ← Agent 沙箱容器由宿主机直接启动
│
├── utbench:latest 容器              ← 评测执行面
│   └── compile / test / coverage / mutation（不需要 docker.sock）
│
└── utbench-agent-opencode 容器      ← Agent 沙箱（统一镜像）
    └── Python + Go + Java + C++ + OpenCode CLI
```

Web 模式下流水线自动拆分：generation 在宿主机 in-process 执行（Agent 沙箱可用宿主机 Docker），
evaluation 在 eval 容器内执行（语言工具链齐全）。eval 容器不需要 docker.sock。

这个拓扑 **不是错的**，但有两个问题：

1. **控制面和评测执行面耦合过重**
   - 外层镜像过胖；
   - Web 也继承了评测工具链依赖；
   - rebuild 成本高，职责边界不清。
2. **Agent 沙箱已经独立，但协议仍偏“Docker 参数拼装”**
   - 缺少正式的 sandbox runtime 配置模型；
   - 未来切换到 remote sandbox / stronger isolation 时，改动面会很大。

## 重构目标

目标不是“一次性上最复杂方案”，而是把运行时边界先钉死。

### 总体目标

将运行时明确拆为三层：

```text
Control Plane
  负责任务编排、Web、资产库、报告

Evaluation Plane
  负责语言工具链、compile/test/coverage/mutation

Agent Sandbox Plane
  负责 agent 执行、workspace 隔离、trace/usage 采集
```

### 核心原则

1. **控制面不再默认背所有语言评测依赖**
2. **评测环境和 Agent 沙箱环境分别版本化**
3. **Agent framework 配置表达“sandbox runtime”，而不是只表达“docker image”**
4. **DOOD 只作为过渡方案，不作为长期架构中心**
5. **能力失败 / 环境失败 / sandbox 失败 / policy 失败分层记录**

## 目标架构

### Phase A：短期可落地架构

```text
宿主机
├── utbench-control:latest
│   ├── CLI / Web
│   ├── SQLite / assets / report
│   └── 通过本地进程或 docker 启动 evaluator / agent sandbox
├── utbench-eval:latest 或 utbench-eval-{lang}:latest
│   └── compile / test / coverage / mutation
└── utbench-agent-{framework}-{lang}:latest
    └── 只包含 Agent 运行所需环境
```

说明：

- **控制面镜像**不再强依赖 Mull、mutmut、JDK、Go、GCC 全家桶。
- **评测镜像**专门承载评测工具链。
- **Agent 沙箱镜像**专门承载 Agent CLI 及任务语言依赖。

### Phase B：中期目标

将 Agent sandbox provider 从“默认 Docker”抽象成多后端：

```text
provider = local | docker | remote | e2b
```

初期仍先实现：

```text
provider = docker
mode = local | docker
```

但配置模型上不再把“docker”写死为唯一世界观。

### Phase C：长期目标

- evaluator 支持独立 worker image / queue
- sandbox 支持 stronger isolation（gVisor / Firecracker / Cube / E2B）
- 控制面成为真正的 orchestration + assets + reporting layer

## 设计决策

### 决策 1：保留“评测镜像”和“Agent 沙箱镜像”分离

结论：**应该分开，而且继续强化分开。**

原因：

- 评测依赖与 Agent 依赖不是一回事。
- 评测镜像强调 compile/test/coverage/mutation 的稳定性。
- Agent 沙箱镜像强调 CLI agent 可运行、受限、可追踪。
- 两者生命周期不同，缓存策略不同，安全边界也不同。

### 决策 2：控制面镜像应从胖镜像中拆出

结论：**应该拆。**

现有 `utbench:latest` 作为“一切都装进去”的镜像，适合 MVP，不适合长期平台化。

建议：

- `utbench-control:latest`
  - Go binary
  - Web 静态资源
  - SQLite
  - 少量系统工具
  - 可选 docker CLI
- `utbench-eval:latest`
  - 语言工具链
  - coverage / mutation
- `utbench-agent-...`
  - agent CLI + 语言运行时

### 决策 3：sandbox 配置要成为独立块

现有写法：

```yaml
frameworks:
  opencode:
    sandbox_mode: docker
    docker_images:
      python: ...
```

目标写法：

```yaml
frameworks:
  opencode:
    sandbox:
      provider: docker
      mode: docker
      images:
        python: ...
        go: ...
      timeout_seconds: 600
      network_disabled: false
      cpu: "2"
      memory: 2g
```

这样做的价值：

- “framework 配置”与“sandbox runtime 配置”语义分离；
- 后续接入 `provider: e2b` 或 `provider: remote` 时，不需要推翻配置协议；
- 资产指纹也能更准确表达 sandbox 版本。

### 决策 4：flat 架构 —— 流水线拆分，generation 在宿主机，evaluation 在容器

已完成：

- 控制面（Web/CLI）直接运行在宿主机，不依赖 Docker 镜像
- Web 模式下流水线自动拆分：generation 在宿主机 in-process 执行，evaluation 在 eval 容器内执行
- Agent 沙箱容器由宿主机直接启动（宿主机有 Docker），eval 容器不需要 docker.sock
- 评测面和 Agent 沙箱容器是同级 peer 容器，不存在嵌套 DOOD
- Agent 沙箱容器（utbench-agent-opencode:latest）包含 OpenCode CLI 和所有语言运行时

## 存在的其他问题

### 1. 失败分类还不够硬

当前不少失败仍然是“最终失败了”，但对用户而言应至少区分：

- agent capability failure
- sandbox preflight failure
- sample environment setup failure
- evaluation tool failure
- policy violation
- infra / docker topology failure

### 2. evaluator 环境和 sandbox 环境 fingerprint 还没形成完整双轨

现在已经有环境指纹，但应明确分成：

- `evaluator_env_fingerprint`
- `sandbox_env_fingerprint`

并分别进入复用判断。

### 3. Web / CLI 对运行时拓扑的可见性不足

用户应能一眼看到本次运行的：

- control backend
- evaluator backend
- sandbox provider
- sandbox image / evaluator image
- 是否走 DOOD

### 4. Docker 设计与资产设计还没完全打通

镜像 digest、sandbox provider、topology mode，应该进入：

- `subject_version`
- `generation_key`
- `evaluation_key`

否则“看起来同一次实验”，实际环境可能并不相同。

## 改造方案

### 第一阶段：协议收敛，不改核心拓扑

目标：

- 引入独立 `sandbox` 配置块；
- 保持兼容旧字段；
- 在 fingerprint / trace / Web 展示里明确 sandbox provider / image；
- 文档层明确三层架构。

这阶段不追求拆镜像，只先把“表达方式”校正。

### 第二阶段：拆控制面镜像

目标：

- 新增 `utbench-control`；
- 保留现有 `utbench-eval` 作为完整评测镜像；
- Web/CLI 默认基于 control image 启动；
- Docker Guide 区分：
  - control only
  - control + evaluator
  - control + nested agent sandbox

### 第三阶段：evaluator 运行时独立

目标：

- evaluator 支持独立 image / backend
- 可以按语言选择 `utbench-eval-go`、`utbench-eval-java`、`utbench-eval-cpp`、`utbench-eval-python`
- control plane 不再默认自带全部评测依赖

### 第四阶段：引入更强 sandbox provider

优先级建议：

1. Docker + gVisor
2. Remote sandbox provider
3. E2B / Cube / Firecracker 类方案

## 建议的近期实施顺序

### P0

1. 配置协议改造：`framework.sandbox.*`
2. trace / fingerprint 增加 sandbox provider / image 语义
3. Web/CLI 显示当前 runtime topology

### P1

4. 新增 `utbench-control` Dockerfile
5. 保留现有胖镜像为 `utbench-eval`
6. 更新启动脚本和文档

### P2

7. evaluator backend 抽象
8. 语言粒度 evaluator image
9. 更强隔离 sandbox provider

## 本轮要先做什么

本轮先做 **第一阶段**：

1. 写下这份设计文档
2. 将 agent framework 配置中的沙箱字段抽象成独立 `sandbox` 块
3. 保持兼容旧配置
4. 为后续拆 evaluator runtime 和 control image 留接口

## 当前进展

### 已完成

- **第一阶段（协议收敛）** ✅
  - sandbox 配置块已独立（`frameworks.{fw}.sandbox.*`），兼容旧字段
  - sandbox provider / image 已进入 trace / generated case / evaluation result / report runtime summary
  - Web 顶部状态栏区分：Docker、拓扑模式、控制镜像、评测镜像
  - 构建镜像弹窗支持 eval / control 目标选择

- **第二阶段（拆控制面镜像）** ✅ 初步完成
  - `Dockerfile.control` 已建立，使用 `docker:28-cli` 多阶段构建复制 CLI 二进制（替代 `docker.io` 重依赖）
  - 主 `Dockerfile`（eval image）也改用同样方式获取 docker CLI
  - `DockerConfig` 已拆为 `ControlImage` / `EvalImage`
  - Web 启动参数支持 `--docker-control-image` / `--docker-eval-image`
  - DOCKER_GUIDE.md 已重写，明确三种镜像的角色、构建命令和使用场景

- **Web runtime topology 展示** ✅
  - Run 详情页头部新增：沙箱 provider、Agent 框架、评测环境指纹
  - 环境检查页新增"运行时拓扑"卡片：拓扑模式、Docker 状态、镜像状态、嵌套沙箱就绪状态
  - 构建弹窗添加 eval/control 用途说明

- **Evaluator backend 抽象** ✅ 完成
  - 新增 `EvalBackend` 接口（`internal/evaluator/backend.go`）
  - `LocalBackend` 实现委托给平台级 `runCommandLocal`（Unix: 进程组 kill；Windows: CommandContext）
  - `DockerBackend` 实现：在评测容器内执行编译/测试/覆盖率/变异命令
  - 兼容包装函数 `runCommandWithProcessGroupKill` 委托给 `evalBackend.RunCommand`
  - 所有 25+ 个现有调用点（go_eval / java_eval / cpp_eval / mutation / service）自动通过 backend 执行
  - 全局 `SetEvalBackend()` / `GetEvalBackend()` 供后续切换
  - CLI 支持 `--eval-backend local|docker` 和 `--eval-docker-image` 参数

- **Flat 架构落地** ✅ 完成
  - 控制面（Web/CLI）直接运行在宿主机，不依赖 Docker 镜像
  - Web 模式下流水线自动拆分：generation 在宿主机 in-process，evaluation 在 eval 容器
  - Agent 沙箱容器由宿主机直接启动，eval 容器不需要 docker.sock
  - Agent 沙箱在独立容器内以 docker 模式运行，包含完整 OpenCode CLI 和语言运行时

- **统一沙箱镜像** ✅ 完成
  - 4 个语言沙箱镜像合并为 1 个 `utbench-agent-opencode:latest`
  - 包含 Python + Go + Java + C++ + OpenCode CLI 全部运行时
  - `frameworkSandboxImage()` 优先级：`Sandbox.Image` → `DockerImage` → `Sandbox.Images[lang]`
  - `build.sh` / `build.ps1` 简化为 3 个镜像构建目标

- **源码只读保护** ✅ 完成
  - Docker 沙箱通过 `ReadOnlyMounts` 将源文件以 `:ro` 挂载到容器
  - 保留 `chmod 0444` 作为纵深防御（local 模式也需要）
  - sandbox fingerprint 包含 readonly_mounts 信息

### 待完成

- **P2-9**: 更强隔离 sandbox provider（gVisor / E2B / Firecracker / Cube）
- **长期**: 远程 sandbox provider、evaluator 独立 worker queue

