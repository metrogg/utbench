# go-ut-bench 测试方案设计文档

**文档版本：** v1.0
**创建日期：** 2026-04-26
**状态：** 待审核

---

## 1. 概述

### 1.1 目标

为 go-ut-bench 项目建立完善的测试体系，确保：
- 单元测试覆盖所有核心模块，基础覆盖率目标 30-50%
- Docker 环境下运行多语言（Python/Go/Java/C++）集成测试
- 验证完整流水线：dataset → runner → evaluator → reporter → store

### 1.2 现状分析

**已有测试模块：**
- `internal/dataset` — 5 个测试用例
- `internal/evaluator` — 21 个测试用例
- `internal/reporter` — 6 个测试用例
- `internal/runner` — 11 个测试用例
- `internal/store` — 1 个测试用例
- `cmd/utbench` — 2 个测试用例

**缺失测试模块：**
- `internal/orchestrator` — 流水线编排（P1 优先级）
- `internal/contracts` — 数据契约（P1 优先级）
- `internal/ctrl` — 控制层（P2 优先级）
- `internal/config` — 配置加载（P2 优先级）
- `internal/obs` — 日志工具（P3 优先级）
- `internal/web` — Web 服务（P3 优先级）

### 1.3 约束条件

- 单元测试覆盖率目标：30-50%（基础覆盖）
- 集成测试运行环境：Docker 容器（主要）
- 多语言组合：Python + Go + Java + C++
- 测试风格：Go table-driven 模式，与现有测试保持一致

---

## 2. 测试架构设计

### 2.1 目录结构

```
go-ut-bench/
│
├── internal/                      # 单元测试（各模块目录）
│   ├── orchestrator/
│   │   └── orchestrator_test.go   # P1 - 流水线编排测试
│   ├── contracts/
│   │   └── contracts_test.go      # P1 - 数据契约验证
│   ├── ctrl/
│   │   └── ctrl_test.go           # P2 - 控制层测试
│   ├── config/
│   │   └── config_test.go         # P2 - 配置加载测试
│   ├── obs/
│   │   └── obs_test.go            # P3 - 日志工具测试
│   └── web/
│   │   └── web_test.go            # P3 - Web 服务测试
│
├── tests/                         # 集成测试目录（新建）
│   ├── integration/
│   │   ├── pipeline_test.go       # 多语言集成测试
│   │   ├── checkpoint_test.go     # 断点续传测试
│   │   ├── mutation_test.go       # 变异测试集成
│   │   └── fixtures/              # 测试数据
│   │       ├── python_samples/    # Python 测试样本
│   │       ├── go_samples/        # Go 测试样本
│   │       ├── java_samples/      # Java 测试样本
│   │       ├── cpp_samples/       # C++ 测试样本
│   │       └── manifests/         # 测试 manifest 文件
│   │
│   └── e2e/
│       └── full_pipeline_test.go  # 完整端到端测试
│
├── scripts/
│   ├── run_integration.sh         # Docker 集成测试运行脚本
│   └── run_coverage.sh            # 覆盖率统计脚本
│
├── docker-compose.test.yml        # 测试专用 compose 文件
│
└── Makefile                       # 新增测试命令
```

### 2.2 测试分层

| 层级 | 类型 | 运行环境 | 命令 |
|------|------|----------|------|
| L1 | P1 单元测试 | 本地 | `make test-unit-p1` |
| L2 | P2 单元测试 | 本地 | `make test-unit-p2` |
| L3 | P3 单元测试 | 本地 | `make test-unit-p3` |
| L4 | 集成测试 | Docker | `make test-integration` |
| L5 | 端到端测试 | Docker | `make test-e2e` |

---

## 3. 单元测试设计

### 3.1 P1 单元测试（orchestrator + contracts）

#### 3.1.1 orchestrator 测试

**文件位置：** `internal/orchestrator/orchestrator_test.go`

**核心测试目标：**
- 流水线阶段调用顺序正确性
- 阶段失败时的错误处理
- checkpoint 加载和保存逻辑
- 并发控制（runner worker pool）

**测试用例：**

| 测试名称 | 测试内容 | 覆盖场景 |
|---------|---------|---------|
| `TestPipelineOrder` | 验证 dataset → runner → evaluator → reporter → store 调用顺序 | 正常流程 |
| `TestPipelineDatasetError` | dataset 阶段返回错误时，流水线终止 | 错误处理 |
| `TestPipelineRunnerPartialFail` | runner 部分任务失败时，继续执行后续阶段 | 部分失败 |
| `TestCheckpointResume` | 加载已有 checkpoint 时，跳过已完成任务 | 断点续传 |
| `TestCheckpointInvalidation` | 改变参数时，checkpoint 被正确失效 | 参数变更 |
| `TestWorkerPoolConcurrency` | 并发调用 LLM API 时，worker pool 正确调度 | 并发控制 |

**测试风格示例：**
```go
func TestPipelineOrder(t *testing.T) {
    // 使用 mock 依赖注入，验证各阶段调用顺序
    mockDataset := &MockDatasetService{}
    mockRunner := &MockRunnerService{}
    mockEvaluator := &MockEvaluatorService{}
    mockReporter := &MockReporterService{}
    mockStore := &MockStoreService{}

    orchestrator := NewOrchestrator(mockDataset, mockRunner, mockEvaluator, mockReporter, mockStore)
    err := orchestrator.Run(ctx, spec)

    require.NoError(t, err)
    // 验证调用顺序
    assert.Equal(t, []string{"dataset", "runner", "evaluator", "reporter", "store"}, mockDataset.callOrder)
}
```

#### 3.1.2 contracts 测试

**文件位置：** `internal/contracts/contracts_test.go`

**核心测试目标：**
- 数据结构 JSON 序列化/反序列化正确性
- 必填字段缺失时的错误处理
- SchemaVersion 常量一致性
- nil 字段处理（不替代为零值）

**测试用例：**

| 测试名称 | 测试内容 | 覆盖场景 |
|---------|---------|---------|
| `TestRunSpecJSONRoundTrip` | RunSpec 序列化后反序列化，数据完整 | JSON 处理 |
| `TestSampleRefRequiredFields` | 必填字段缺失时，解析返回错误 | 字段验证 |
| `TestEvaluationResultNilFields` | nil 字段（如 Coverage）序列化为 null | nil 处理 |
| `TestGeneratedManifestValidation` | manifest 文件结构验证 | 文件格式 |
| `TestDatasetClassConstants` | 验证 DatasetClass 常量值正确 | 常量一致性 |
| `TestSchemaVersionFormat` | 验证 SchemaVersion 格式符合规范 | 版本格式 |

**测试风格示例：**
```go
func TestEvaluationResultNilFields(t *testing.T) {
    row := contracts.EvaluationResult{
        CompilePass: true,
        Coverage:    nil,  // 应序列化为 null，而非 0
    }

    data, err := json.Marshal(row)
    require.NoError(t, err)

    // 验证 nil 字段序列化为 null
    assert.Contains(t, string(data), `"coverage":null`)
}
```

---

### 3.2 P2 单元测试（ctrl + config）

#### 3.2.1 ctrl 测试

**文件位置：** `internal/ctrl/ctrl_test.go`

**核心测试目标：**
- CLI 参数解析正确性
- 参数默认值处理
- 参数组合冲突检测
- 配置文件路径验证

**测试用例：**

| 测试名称 | 测试内容 | 覆盖场景 |
|---------|---------|---------|
| `TestParseArgsDefaults` | 验证默认参数值正确应用 | 默认值 |
| `TestParseArgsRequiredFlags` | 必填参数缺失时返回错误 | 参数验证 |
| `TestParseArgsModelList` | 多模型参数解析为数组 | 列表参数 |
| `TestParseArgsLanguageList` | 多语言参数解析为数组 | 列表参数 |
| `TestParseArgsInvalidCombination` | 无效参数组合时返回错误 | 参数冲突 |
| `TestConfigPathResolution` | 配置文件路径相对/绝对路径解析 | 路径处理 |

#### 3.2.2 config 测试

**文件位置：** `internal/config/config_test.go`

**核心测试目标：**
- models.yaml 配置文件解析
- API Key 环境变量读取
- 缺失配置项时的错误处理
- 配置项覆盖逻辑（命令行 > 环境变量 > 配置文件）

**测试用例：**

| 测试名称 | 测试内容 | 覆盖场景 |
|---------|---------|---------|
| `TestLoadModelsYAML` | 正确解析 models.yaml 结构 | 配置解析 |
| `TestLoadModelsYAMLInvalid` | 无效 YAML 格式时返回错误 | 错误处理 |
| `TestGetAPIKeyFromEnv` | 从环境变量读取 API Key | 环境变量 |
| `TestGetAPIKeyMissing` | API Key 缺失时返回错误 | 缺失处理 |
| `TestConfigOverridePriority` | 验证配置优先级正确 | 优先级逻辑 |
| `TestModelEndpointResolution` | 验证模型 endpoint URL 正确拼接 | URL 处理 |

---

### 3.3 P3 单元测试（obs + web）

#### 3.3.1 obs 测试

**文件位置：** `internal/obs/obs_test.go`

**核心测试目标：**
- Logger 初始化和配置
- 日志级别过滤
- slog 格式化输出正确性

**测试用例：**

| 测试名称 | 测试内容 | 覆盖场景 |
|---------|---------|---------|
| `TestLoggerInit` | Logger 初始化返回有效实例 | 初始化 |
| `TestLogLevelFilter` | 日志级别低于配置时不输出 | 级别过滤 |
| `TestLoggerFormat` | 输出格式符合 slog 结构 | 格式验证 |

#### 3.3.2 web 测试

**文件位置：** `internal/web/web_test.go`

**核心测试目标：**
- HTTP 路由正确性
- 静态文件服务
- API 响应格式

**测试用例：**

| 测试名称 | 测试内容 | 覆盖场景 |
|---------|---------|---------|
| `TestStaticFileServe` | 验证静态文件正确返回 | 静态资源 |
| `TestAPIReportEndpoint` | `/api/report` 返回正确 JSON | API 路由 |
| `TestWebServerStartup` | Web 服务启动和关闭 | 服务生命周期 |

---

## 4. 集成测试设计

### 4.1 Docker 测试环境

**配置文件：** `docker-compose.test.yml`

```yaml
version: '3.8'
services:
  utbench-test:
    build:
      context: .
      dockerfile: Dockerfile
    environment:
      - DEEPSEEK_API_KEY=${DEEPSEEK_API_KEY}
      - TEST_MODE=integration
    volumes:
      - ./tests/integration/fixtures:/test_fixtures:ro
      - ./test_artifacts:/test_artifacts
    command: ["./utbench", "run", "--config", "/app/configs/models.yaml", "--langs", "python,go,java,cpp", "--max-samples", "3", "--dry-run"]
```

### 4.2 多语言集成测试

**文件位置：** `tests/integration/pipeline_test.go`

**测试用例：**

| 测试名称 | 测试内容 | 运行条件 |
|---------|---------|---------|
| `TestPythonPipeline` | Python 数据集完整流程（编译→测试→覆盖率→变异） | Docker |
| `TestGoPipeline` | Go 数据集完整流程（gremlins 变异测试） | Docker |
| `TestJavaPipeline` | Java 数据集完整流程（pitest 变异测试） | Docker |
| `TestCppPipeline` | C++ 数据集完整流程（Mull 变异测试） | Docker |
| `TestMultiLanguageParallel` | 多语言并行运行，验证无冲突 | Docker |
| `TestReportGeneration` | 验证 HTML 报告正确生成 | Docker |

### 4.3 测试数据样本

**目录结构：** `tests/integration/fixtures/`

```
tests/integration/fixtures/
├── python_samples/
│   ├── simple_function.py       # 简单函数：def add(a, b): return a + b
│   └── boundary_001.py          # 边界条件测试样本
├── go_samples/
│   ├── simple_function.go       # 简单函数：func Add(a, b int) int { return a + b }
│   └── boundary_001.go          # 边界条件测试样本
├── java_samples/
│   ├── SimpleFunction.java      # 简单类：public int add(int a, int b) { return a + b; }
│   └── Boundary.java            # 边界条件测试样本
├── cpp_samples/
│   ├── simple_function.cpp      # 简单函数：int add(int a, int b) { return a + b; }
│   ├── simple_function.h        # 头文件
│   ├── boundary_001.cpp         # 边界条件测试样本
│   └── boundary_001.h           # 头文件
└── manifests/
    └── test_manifest.json       # 测试 manifest 文件
```

### 4.4 C++ 工具链依赖

| 工具 | 用途 | Dockerfile 安装状态 |
|------|------|---------------------|
| `g++` / `clang++` | C++ 编译器 | 已安装 |
| `GoogleTest` | C++ 测试框架 | 已安装 |
| `Mull` | C++ 变异测试工具 | 已安装（参考 cpp_eval_test.go） |
| `gcov` / `lcov` | C++ 覆盖率工具 | 需验证 |

### 4.5 Checkpoint 集成测试

**文件位置：** `tests/integration/checkpoint_test.go`

| 测试名称 | 测试内容 |
|---------|---------|
| `TestCheckpointSave` | 流程中断后 checkpoint 正确保存 |
| `TestCheckpointResume` | 恢复运行时跳过已完成任务 |
| `TestCheckpointInvalidation` | 参数变更后 checkpoint 失效重新运行 |

---

## 5. 测试运行命令

### 5.1 Makefile 配置

```makefile
# Makefile - 测试相关命令

.PHONY: test test-unit test-unit-p1 test-unit-p2 test-unit-p3 test-integration test-e2e coverage test-clean

# 运行所有测试
test: test-unit test-integration

# 单元测试（全部）
test-unit:
	go test ./internal/... ./cmd/... -v -count=1

# P1 单元测试（orchestrator + contracts）
test-unit-p1:
	go test ./internal/orchestrator/... ./internal/contracts/... -v -count=1

# P2 单元测试（ctrl + config）
test-unit-p2:
	go test ./internal/ctrl/... ./internal/config/... -v -count=1

# P3 单元测试（obs + web）
test-unit-p3:
	go test ./internal/obs/... ./internal/web/... -v -count=1

# 集成测试（Docker 环境）
test-integration:
	./scripts/run_integration.sh

# 端到端测试
test-e2e:
	./scripts/run_e2e.sh

# 覆盖率报告
coverage:
	go test ./internal/... ./cmd/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 清理测试产物
test-clean:
	rm -rf ./test_artifacts
	rm -rf ./coverage.out ./coverage.html
```

### 5.2 运行脚本

**文件位置：** `scripts/run_integration.sh`

```bash
#!/bin/bash
# 运行 Docker 集成测试

set -e

echo "Building test image..."
docker build -t utbench:test .

echo "Running integration tests..."
docker-compose -f docker-compose.test.yml up --abort-on-container-exit

echo "Collecting results..."
docker cp utbench-test:/test_artifacts ./test_artifacts

echo "Integration tests completed."
```

---

## 6. 测试执行流程

```
┌─────────────────────────────────────────────────────────────┐
│                      测试执行流程                            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  make test                                                  │
│       │                                                     │
│       ▼                                                     │
│  ┌─────────┐     ┌──────────────────┐                      │
│  │ 单元测试 │────▶│ Docker 集成测试  │                      │
│  └─────────┘     └──────────────────┘                      │
│       │                     │                              │
│       ▼                     ▼                              │
│  覆盖率统计            多语言测试                           │
│  (30-50%)          (Python/Go/Java/C++)                    │
│                                                             │
│  make coverage          make test-integration              │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 7. 预期成果

### 7.1 单元测试补充

| 模块 | 优先级 | 预期新增测试用例数 | 预期覆盖率 |
|------|--------|-------------------|------------|
| orchestrator | P1 | 6 | 35-45% |
| contracts | P1 | 6 | 40-50% |
| ctrl | P2 | 6 | 30-40% |
| config | P2 | 6 | 30-40% |
| obs | P3 | 3 | 30-40% |
| web | P3 | 3 | 30-40% |

### 7.2 集成测试新增

| 测试类型 | 预期新增测试用例数 |
|----------|-------------------|
| 多语言 Pipeline | 6 |
| Checkpoint | 3 |
| 端到端 | 1 |

### 7.3 文件清单

| 文件路径 | 类型 |
|----------|------|
| `internal/orchestrator/orchestrator_test.go` | 新增 |
| `internal/contracts/contracts_test.go` | 新增 |
| `internal/ctrl/ctrl_test.go` | 新增 |
| `internal/config/config_test.go` | 新增 |
| `internal/obs/obs_test.go` | 新增 |
| `internal/web/web_test.go` | 新增 |
| `tests/integration/pipeline_test.go` | 新增 |
| `tests/integration/checkpoint_test.go` | 新增 |
| `tests/integration/fixtures/*` | 新增 |
| `docker-compose.test.yml` | 新增 |
| `scripts/run_integration.sh` | 新增 |
| `Makefile` | 修改 |

---

## 8. 实施建议

### 8.1 实施顺序

1. **第一阶段**：P1 单元测试（orchestrator + contracts）
2. **第二阶段**：集成测试环境搭建（Docker + fixtures）
3. **第三阶段**：多语言集成测试实现
4. **第四阶段**：P2 单元测试（ctrl + config）
5. **第五阶段**：P3 单元测试（obs + web）

### 8.2 注意事项

1. **Mock 设计**：orchestrator 测试需要为各服务层设计 mock 接口
2. **Docker 工具链**：确保 Dockerfile 中已安装所有语言工具链
3. **测试数据**：fixtures 样本应足够简单，避免测试运行时间过长
4. **覆盖率验证**：实施后需验证实际覆盖率是否达到目标

---

## 9. 审核确认

- [ ] 测试架构设计确认
- [ ] 单元测试设计确认
- [ ] 集成测试设计确认
- [ ] Makefile 命令确认
- [ ] 实施计划确认