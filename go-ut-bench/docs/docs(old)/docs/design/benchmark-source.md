# 评测集来源说明

## 一、说明目的

本文档用于系统说明当前多语言单元测试评测集的来源构成，重点回答以下问题：

- 基础评测集 `self_contained` 分别来自哪些原始数据集或源码源池
- 项目级评测集 `module_level` 分别来自哪些原始数据集或源码源池
- 每种语言的数据是如何从原始源池加工、筛选、分类并导出的
- 为什么基础评测集与项目级评测集不能被理解为“同一种难度的数据”

本文档适合用于：

- 对组员说明数据来源
- 对老师或评审解释评测集构成
- 在汇报中回答“数据从哪里来”的问题

## 二、总览：不是一套单一数据集，而是多源加工后的双轨评测集

当前评测体系并不是由单一原始数据集直接导出，而是将多个不同来源的数据源，按照统一规则加工成两类评测轨道：

- `self_contained`：基础评测集
- `module_level`：项目级评测集

更准确地说：

- `self_contained` 是一种“导出后的评测形态”
- `module_level` 也是一种“导出后的评测形态”
- 它们并不是原始数据集名称

真正的原始来源包括以下几类：

- 公开 benchmark 数据集
- 本地真实项目源码树
- 已安装 Python 包源码
- 单个大型开源源码文件
- 外部同步得到的 benchmark 子集

当前实际参与构建的原始源池主要包括：

- BigCodeBench
- Python 本地源码根目录
- go_all.jsonl
- hm-dianping
- CPP-UT-Bench
- lodash.js
- HumanEval-X
- AutoCodeBenchmark

因此，当前评测集本质上是一个“多源输入、统一加工、双轨导出”的评测体系。

## 三、基础评测集与项目级评测集的本质区别

### 1. 基础评测集 `self_contained`

基础评测集的目标是：

- 尽量保证单文件可运行
- 尽量减少外部依赖
- 适合统一 benchmark
- 适合多模型横向比较

因此，它的原始来源通常更偏向：

- 完整题解
- 可独立运行的函数实现
- 或者经过依赖闭包处理后变成单文件样本的源码

### 2. 项目级评测集 `module_level`

项目级评测集的目标是：

- 保留真实项目模块语义
- 保留项目上下文依赖
- 更接近真实研发环境
- 用于恢复项目环境后再执行测试

因此，它的原始来源通常更偏向：

- 真实项目源码文件
- 带项目来源信息的函数级样本
- 安装包源码
- 带 Ground Truth 测试线索的项目模块数据

换句话说：

- 基础评测集更偏 benchmark 逻辑
- 项目级评测集更偏真实工程逻辑

## 四、按语言详细说明数据来源

---

## 四点一、Python

### 1. Python 基础评测集 `self_contained` 的原始来源

Python 基础评测集主要来自：

- BigCodeBench

底层原始文件是：

- `data/bcb/bcb_all.jsonl`

构建脚本是：

- `build_self_contained_python_dataset.py`

### 2. Python 基础评测集的加工过程

构建流程大致如下：

1. 从 `bcb_all.jsonl` 中逐条读取原始样本
2. 将题目 prompt 与 `canonical_solution` 拼接成完整 Python 代码
3. 去除注释和 docstring
4. 做 AST 解析
5. 计算复杂度特征，包括：
   - 有效代码行数
   - 圈复杂度
   - 最大嵌套深度
   - 分支数
   - try/except 数
   - 是否递归
   - 是否有类
   - 是否存在状态型 OOP
6. 调用自包含检查逻辑，过滤掉不适合单文件评测的样本
7. 对剩余样本按四类任务打分
8. 按类别平衡采样，最终导出 200 条样本

### 3. Python 基础评测集过滤掉了什么

基础评测集在构建时重点过滤掉以下不适合单文件评测的情况：

- 相对导入
- 包内引用
- 第三方依赖

因此，Python 基础评测集可以准确描述为：

“来自 BigCodeBench 完整解答，经过去注释、自包含检查、复杂度筛选、类别打分和平衡采样后导出的 200 条 Python 自包含样本。”

### 4. Python 项目级评测集 `module_level` 的原始来源

Python 项目级评测集不是来自单一 benchmark，而是来自多个源码根目录扫描。

源码根目录包括：

- `bigcodebench`
- `LiveCodeBench`
- `data/bcb/code/ut-bench`
- `venv/Lib/site-packages`

构建脚本是：

- `build_python_module_level_dataset.py`

辅助规则来源于：

- `build_realworld_module_dataset.py`

### 5. Python 项目级评测集的加工过程

构建流程大致如下：

1. 遍历多个源码根目录中的 `.py` 文件
2. 排除不适合成为样本的文件，例如：
   - 测试目录
   - 缓存目录
   - `__init__.py`
3. 去除注释和 docstring
4. 做 AST 解析
5. 计算复杂度和结构特征
6. 要求样本满足一定的复杂度与长度阈值
7. 不强制要求其单文件自包含
8. 使用路径文本和代码文本做任务类别打分
9. 按四类任务平衡采样，最终导出 200 条样本

### 6. Python 项目级评测集的实际最终来源分布

虽然扫描源池包括多个目录，但当前最终入选样本几乎全部来自：

- `venv/Lib/site-packages`

实际统计结果为：

- `199` 条来自 `venv/Lib/site-packages`
- `1` 条来自 `bigcodebench`

说明当前 Python 项目级评测集本质上更接近：

“真实 Python 安装包生态中的模块源码样本集”

其中入选较多的包包括：

- pandas
- huggingface_hub
- numpy
- datasets
- fsspec
- httpcore
- _pytest
- aiohttp

因此，Python 项目级评测集可以准确描述为：

“来自多个 Python 源码根目录扫描、但最终几乎全部选自当前环境 site-packages 的真实模块源码，经复杂度筛选和类别平衡后导出的项目级评测集。”

---

## 四点二、Go

### 1. Go 基础评测集 `self_contained` 的原始来源

Go 基础评测集来自：

- `data/go/go_all.jsonl`

构建脚本是：

- `build_go_two_track_datasets.py`

### 2. Go 基础评测集的加工过程

构建流程大致如下：

1. 从 `go_all.jsonl` 中读取函数记录
2. 使用 `whole_func_string` 作为核心源码
3. 保留仓库名、函数名、函数路径等元信息
4. 先对样本做四类任务打分
5. 再做 Go 版本的自包含检查，包括：
   - 是否存在方法接收者
   - 是否依赖 testing
   - 是否只依赖标准库
   - 参数和返回值中是否出现不允许的自定义类型
6. 对通过自包含检查的样本自动包装成单文件形式
   - 自动补 `package main`
   - 自动生成标准库 import
7. 最终按四类平衡采样，导出 200 条样本

因此，Go 基础评测集可以描述为：

“来自 go_all.jsonl 的真实 Go 函数数据，经标准库约束、自包含检查、单文件包装和类别平衡采样后导出的 200 条基础评测样本。”

### 3. Go 项目级评测集 `module_level` 的原始来源

Go 项目级评测集也来自：

- `data/go/go_all.jsonl`

也就是说，Go 双轨数据是“同源分流”。

区别在于：

- 基础评测集要求自包含、可单文件运行
- 项目级评测集保留项目属性，不要求单文件自包含

因此，Go 项目级评测集可以描述为：

“来自 go_all.jsonl 的真实项目函数样本，按项目级模块评测标准筛选和平衡采样后导出的 200 条数据。”

---

## 四点三、Java

### 1. Java 基础评测集 `self_contained` 的原始来源

Java 基础评测集主要来自两个外部 benchmark：

- HumanEval-X
- AutoCodeBenchmark

这两个外部源会先同步到本地缓存目录，再参与构建。

同步脚本是：

- `sync_external_self_contained_sources.py`

同步后的本地中间文件包括：

- `humanevalx_java.jsonl`
- `autocodebenchmark_java.jsonl`

构建脚本是：

- `build_java_module_level_dataset.py`

### 2. Java 基础评测集的加工过程

构建流程大致如下：

1. 从 `HumanEval-X` 读取 Java 任务
2. 将 `prompt + canonical_solution` 拼成完整 Java 代码
3. 从 `AutoCodeBenchmark` 读取 Java `canonical_solution`
4. 过滤不满足自包含要求的代码，例如：
   - 非 `java.*` / `javax.*` import
   - `package` 声明
   - `static import`
   - 明显第三方包引用
5. 结合代码、问题描述、测试参考和难度信息做任务类别打分
6. 按四类平衡采样，最终导出 200 条样本

### 3. Java 基础评测集最终来源占比

当前最终统计为：

- `HumanEval-X`：`55`
- `AutoCodeBenchmark`：`145`

因此，Java 基础评测集可以描述为：

“来自 HumanEval-X 与 AutoCodeBenchmark 的 Java 完整实现，经单文件自包含约束过滤、类别打分和平衡采样后导出的 200 条基础评测样本。”

### 4. Java 项目级评测集 `module_level` 的原始来源

Java 项目级评测集来自真实 Java 项目源码：

- `hm-dianping`

源码根目录为：

- `data/bcb/code/hm-dianping/src/main/java`

构建脚本仍然是：

- `build_java_module_level_dataset.py`

### 5. Java 项目级评测集的加工过程

构建流程大致如下：

1. 遍历 `hm-dianping` 项目中的 `.java` 文件
2. 读取完整源码
3. 统计有效代码行数
4. 根据路径和代码语义做任务类别打分
5. 将每个源码文件归入得分最高的类别
6. 导出为项目级样本

### 6. Java 项目级评测集当前状态

Java 项目级评测集当前并没有平衡到 200 条，而是按当前可用项目源码直接导出。

当前总量为：

- `62`

全部来源于：

- `hm-dianping`

因此，Java 项目级评测集可以描述为：

“来自真实 Java 工程 hm-dianping 的源码文件级样本，经路径与代码语义分类后导出的项目级评测集。”

---

## 四点四、C++

### 1. C++ 基础评测集 `self_contained` 的原始来源

C++ 基础评测集主要来自两个外部 benchmark：

- HumanEval-X
- AutoCodeBenchmark

同步脚本同样是：

- `sync_external_self_contained_sources.py`

同步后的本地中间文件包括：

- `humanevalx_cpp.jsonl`
- `autocodebenchmark_cpp.jsonl`

构建脚本是：

- `build_cpp_module_level_dataset.py`

### 2. C++ 基础评测集的加工过程

构建流程大致如下：

1. 从两个外部源读取 C++ 完整实现
2. 检查是否满足单文件自包含要求
3. 过滤掉不适合独立运行的样本，例如：
   - 本地头文件引用
   - 明显第三方库依赖
   - gtest/gmock、absl、boost、grpc、protobuf 等工程依赖
4. 对剩余样本按四类任务打分
5. 按四类平衡采样，最终导出 200 条样本

### 3. C++ 基础评测集最终来源占比

当前最终统计为：

- `HumanEval-X`：`55`
- `AutoCodeBenchmark`：`145`

因此，C++ 基础评测集可以描述为：

“来自 HumanEval-X 与 AutoCodeBenchmark 的 C++ 完整实现，经头文件与第三方依赖过滤、自包含约束和类别平衡采样后导出的 200 条基础评测样本。”

### 4. C++ 项目级评测集 `module_level` 的原始来源

C++ 项目级评测集来自：

- CPP-UT-Bench

底层原始文件是：

- `data/cpp_ut/cpp_ut_all.jsonl`

构建脚本是：

- `build_cpp_module_level_dataset.py`

### 5. C++ 项目级评测集的加工过程

构建流程大致如下：

1. 从 `cpp_ut_all.jsonl` 中读取源码和真实测试参考
2. 保留项目来源信息
3. 根据源码和 Ground Truth 测试中的关键词做四类任务打分
4. 按四类平衡采样，导出 200 条项目级样本

因此，C++ 项目级评测集可以描述为：

“来自 CPP-UT-Bench 的真实项目模块及真实测试线索，经任务分类和平衡采样后导出的 200 条项目级评测样本。”

---

## 四点五、JavaScript

### 1. JavaScript 基础评测集 `self_contained` 的原始来源

JavaScript 基础评测集来自：

- `lodash.js`

构建脚本是：

- `build_javascript_two_track_datasets.py`

### 2. JavaScript 基础评测集的加工过程

构建流程大致如下：

1. 从 `lodash.js` 中抽取函数
2. 分析每个函数调用了哪些本地 helper
3. 递归展开依赖闭包
4. 如果闭包内所有依赖都能在本地源码中解析，就将其打包成一个 bundle
5. 如果存在无法消解的外部调用，则剔除
6. 将成功闭包展开后的 bundle 视为一个单文件样本
7. 再根据代码语义做任务分类

因此，JavaScript 基础评测集并不是原始就独立存在的 benchmark，而是：

“从真实项目源码 lodash.js 中，通过依赖闭包内联和本地 helper 打包加工得到的自包含 JavaScript 样本集。”

### 3. JavaScript 基础评测集当前状态

由于并不是所有函数都能成功形成可独立运行的闭包，因此当前基础评测集总量不是 200，而是：

- `88`

### 4. JavaScript 项目级评测集 `module_level` 的原始来源

JavaScript 项目级评测集也来自：

- `lodash.js`

它和基础评测集的区别在于：

- 基础评测集会做依赖闭包内联
- 项目级评测集保留为项目函数切片，不要求自包含

### 5. JavaScript 项目级评测集的加工过程

构建流程大致如下：

1. 从 `lodash.js` 中抽取原始函数
2. 保留其项目来源属性
3. 根据路径和代码语义做四类任务打分
4. 按四类平衡采样，导出 200 条项目级样本

因此，JavaScript 项目级评测集可以描述为：

“来自 lodash.js 的真实项目函数切片，经类别平衡采样后导出的项目级评测样本。”

## 五、五种语言来源总结表

| 语言 | 基础评测集来源 | 项目级评测集来源 |
|------|----------------|------------------|
| Python | BigCodeBench | 多源码根目录扫描，最终几乎全部来自 site-packages |
| Go | go_all.jsonl 中的可自包含函数 | go_all.jsonl 中的项目函数 |
| Java | HumanEval-X + AutoCodeBenchmark | hm-dianping 项目源码 |
| C++ | HumanEval-X + AutoCodeBenchmark | CPP-UT-Bench |
| JavaScript | lodash.js 经依赖闭包内联后的 bundle | lodash.js 原始函数切片 |

## 六、最适合对外讲解的总结表述

可以将整套来源概括为：

“当前评测集并非直接来自单一数据集，而是将多个原始源池按照统一规则加工成两类评测轨道：基础评测集主要来自可独立运行的完整实现或经过自包含处理后的代码样本，包括 BigCodeBench、HumanEval-X、AutoCodeBenchmark 以及从真实源码中闭包展开得到的 JavaScript bundle；项目级评测集主要来自真实项目模块源码，包括 Python 的 site-packages 模块、Go 项目函数、Java 的 hm-dianping、C++ 的 CPP-UT-Bench 和 JavaScript 的 lodash 函数切片。” 

## 七、必须强调的一点

由于两条轨道的来源逻辑不同，基础评测集与项目级评测集不应被当成同一难度的数据集。

更准确地说：

- 基础评测集强调可复现、可比较、可自动化
- 项目级评测集强调真实项目上下文中的可落地性

这也是为什么它们需要分开统计、分开解释、分开汇报。

