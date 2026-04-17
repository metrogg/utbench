# 组员交接说明

## 1. 交接包里有什么

交接包保留了可直接运行的原始相对路径，并额外补了说明文档和环境文件。核心结构如下：

- `data/`
  - 各语言的 `self_contained` / `module_level` 数据与代码文件
  - Java / C++ 新增外部源缓存 `data/external_sources/java_cpp_self_contained`
- 根目录脚本
  - 数据构建脚本
  - 数据校验脚本
  - Python 项目级恢复与执行脚本
  - Java / C++ 外部源同步脚本
- `eval/`
  - 现有评测辅助脚本
- `env/`
  - Python 项目级评测环境依赖文件
- `examples/`
  - 一个已经实际跑通的 Python 项目级评测示例
- `docs/`
  - 多语言双轨数据说明
  - Python 数据说明
  - 环境初始化脚本

## 2. 两类数据如何理解

- `self_contained`
  - 面向单文件评测
  - 输入是一段尽量自包含的代码
  - 适合统一 prompt、统一编译运行、统一 coverage 采集
- `module_level`
  - 面向项目级 / 模块级评测
  - 输入是“真实项目中的一个模块”
  - 允许存在包内引用、相对导入、第三方依赖
  - 不能按单文件沙箱方式直接跑

一句话区分：

- `self_contained` 是基础 benchmark
- `module_level` 是真实项目能力验证

## 3. 推荐如何使用

### 如果是继续写单文件评测脚本

直接使用：

- `data/python_ut_dataset/python_dataset_self_contained.json`
- `data/go_ut_dataset/go_dataset_self_contained.json`
- `data/java_ut_dataset/java_dataset_self_contained.json`
- `data/cpp_ut_dataset/cpp_dataset_self_contained.json`
- `data/javascript_ut_dataset/javascript_dataset_self_contained.json`

### 如果是继续写项目级评测脚本

直接使用：

- `data/python_ut_dataset/python_dataset_module_level.json`
- `data/go_ut_dataset/go_dataset_module_level.json`
- `data/java_ut_dataset/java_dataset_module_level.json`
- `data/cpp_ut_dataset/cpp_dataset_module_level.json`
- `data/javascript_ut_dataset/javascript_dataset_module_level.json`

说明：

- `module_level` 评测时，不要把代码文件当成独立单文件运行
- 应该按“恢复项目上下文 -> 生成测试 -> 在项目环境执行 -> 采集目标模块 coverage”的方式跑

## 3.1 项目级评测当前完成度

下面这张表是最重要的现状说明，组员可以直接按这个判断接下来该做什么：

| 语言 | `module_level` 数据集 | 原项目/环境恢复状态 | 是否已能直接做项目级评测 | 还需要补什么 |
|------|------------------------|----------------------|--------------------------|--------------|
| Python | 已准备好 | 部分已具备 | 可以做演示评测 | 批量化执行器、更多项目恢复规则 |
| Go | 已准备好 | 未完成 | 不能直接完整评测 | `go.mod` 恢复、包结构恢复、测试注入、coverage 执行器 |
| Java | 已准备好 | 未完成 | 不能直接完整评测 | Maven/Spring 项目恢复、依赖安装、测试落盘、coverage 执行器 |
| C++ | 已准备好 | 未完成 | 不能直接完整评测 | 头文件/构建系统恢复、依赖库安装、测试 target、coverage 执行器 |
| JavaScript | 已准备好 | 未完成 | 不能直接完整评测 | `package.json` 依赖恢复、模块结构恢复、测试注入、coverage 执行器 |

一句话理解：

- `module_level` 数据已经齐了
- 但除 Python 外，其他语言目前还没有“开箱即跑”的项目级恢复执行器
- 所以组员不能只拿 `module_level.json` 就直接完整评测，还需要继续补恢复脚本

## 3.2 组员是否需要另外找原项目

需要，但不同语言程度不同：

- Python
  - 大多数样本来自当前工作区的 `venv/Lib/site-packages`
  - 所以很多样本不需要额外去网上找原项目，先用当前环境就能做恢复演示
- Go / Java / C++ / JavaScript
  - 目前 `module_level` 更多是“目标模块定义”和“代码文件集合”
  - 还不是完整 repo 快照
  - 所以组员通常还需要补项目结构、依赖环境、测试入口，必要时需要回到原项目来源恢复

## 3.3 组员拿到包后怎么判断能不能直接跑

可以直接按下面规则判断：

- 如果任务是单文件评测：
  - 直接使用 `self_contained`，可以开始写执行脚本
- 如果任务是 Python 项目级评测：
  - 可以先用现有 `python_project_module_eval.py` 做演示和扩展
- 如果任务是 Go / Java / C++ / JavaScript 项目级评测：
  - 不要直接开跑
  - 先补“项目恢复 + 测试注入 + coverage 执行器”

## 4. Java / C++ 新增数据源说明

这次新增了两条外部数据源，用来补齐 Java 和 C++ 的 `self_contained`：

- `HumanEval-X`
- `AutoCodeBenchmark`

同步后的本地缓存目录是：

- `data/external_sources/java_cpp_self_contained/humanevalx_java.jsonl`
- `data/external_sources/java_cpp_self_contained/humanevalx_cpp.jsonl`
- `data/external_sources/java_cpp_self_contained/autocodebenchmark_java.jsonl`
- `data/external_sources/java_cpp_self_contained/autocodebenchmark_cpp.jsonl`

如果后续需要重新同步，运行：

```powershell
python sync_external_self_contained_sources.py
```

## 5. Python 项目级评测怎么启动

### 5.1 环境准备

建议在交接包根目录执行：

```powershell
.\docs\setup_python_handoff_env.ps1
```

这会做两件事：

- 创建 `venv/`
- 安装 `env/requirements_python_project_eval.txt`

注意：

- 当前 Python `module_level` 样本里，绝大多数 `source_path` 都指向 `venv/Lib/site-packages/...`
- 所以请保留虚拟环境目录名为 `venv`
- 如果把环境目录改成别的名字，现有恢复脚本会找不到源码路径

### 5.2 运行一个项目级样例

```powershell
venv\Scripts\python.exe python_project_module_eval.py `
  --sample-id python_simple_function_0 `
  --test-file examples\python_project_eval\generated_test_dateutil_relativedelta.py `
  --model-name demo_model
```

运行后会在：

- `eval/project_level_runs/<sample_id>/<timestamp>/`

生成：

- 恢复后的工作区
- 生成测试文件
- `pytest` 结果
- coverage JSON
- 汇总报告

## 6. 已知限制

- Python `module_level` 数据集中有 `199/200` 条来自 `venv/Lib/site-packages`
- 还有 `1/200` 条来自 `bigcodebench/analysis/get_results.py`
- 这意味着：
  - 对大多数 Python 模块级样本，只要创建好 `venv` 并装好依赖，就可以恢复
  - 极少数本地 repo 来源样本，仍然可能需要额外源码目录或额外依赖

- Java / C++ 的 `self_contained` 虽然已经补齐，但数据源来自外部 benchmark，不等同于真实项目源码
- Go / JavaScript 的项目级恢复脚本还需要继续开发

## 6.1 项目级评测特别提醒

请组员不要误解 `module_level` 的含义：

- `module_level` 不等于“完整项目已打包完成”
- `module_level` 的意思是“评测对象是项目中的模块”
- 真正执行时，通常还需要恢复包结构、依赖和测试环境

所以：

- 这份交接包已经把“评测样本定义”准备好了
- 但“完整项目恢复执行器”目前只对 Python 给出了现成样例

## 7. 推荐后续优先级

1. 先基于 `self_contained` 数据补齐统一单文件评测流水线
2. 再基于 Python `module_level` 把项目级评测流程正式化
3. 然后把 Python 的流程迁移到 Go / JavaScript
4. 最后再考虑 Java / C++ 的项目环境恢复

## 8. 交接包里最重要的文件

- `docs/multilang_two_track_handoff.md`
- `docs/python_eval_dataset_handoff.md`
- `python_project_module_eval.py`
- `sync_external_self_contained_sources.py`
- `verify_dataset.py`
- `env/requirements_python_project_eval.txt`
- `examples/python_project_eval/example_report.json`

## 9. 重新打包命令

如果后续还要重新生成一份新的交接包，在仓库根目录运行：

```powershell
python package_team_handoff.py --zip
```

输出目录在：

- `dist/team_handoff_<timestamp>/`
- `dist/team_handoff_<timestamp>.zip`
