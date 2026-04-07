# MEMORY.md - ut-bench 项目长期记忆

## 项目概述
- 多模型单元测试生成效果横向评测工具
- 当前分支: `feat/evaluator`
- Python 3.10+，依赖 pyyaml、pytest、coverage、mutmut

## 架构
- Runner: 调用模型 API 生成测试代码
- Evaluator: 编译→执行→覆盖率→变异测试→聚合
- Reporter: 生成 JSON/CSV/HTML 报告

## 当前聚焦
- Python 评测链路（其他语言 toolchain 未实现）
- 4 个 baseline 样本: boundary, complex_dependency, interface_mock, simple_function

## 关键设计决策
- Self-contained 分层: allowlist 区分可评测样本
- AST import 重写: 修复 LLM 生成测试的导入问题
- mutmut multiprocessing workaround: sitecustomize.py monkey-patch

## 用户偏好
- 中文文档优先
- 需要 PEP 8 代码风格
- 技术指导型协作
