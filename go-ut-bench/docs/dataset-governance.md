# 数据集治理说明（MVP）

## 1. 背景

现有仓库中数据集来源较多、目录组织不完全统一。
为了快速推进工具开发，MVP 采用“逻辑治理优先”的方案。

## 2. 当前分类口径

MVP 当前按两层组织：

- 大类：`self_contained` / `module_level`
- 子类：`boundary` / `simple_function` / `complex_dependency` / `interface_mock`

## 3. 当前判定规则（临时）

由样本 ID + 相对路径联合判定：

- 包含 `self_contained` -> `self_contained`
- 包含 `module_level` -> `module_level`
- 包含子类关键词（`boundary` / `simple_function` / `complex_dependency` / `interface_mock`）但未命中大类 -> `self_contained`
- 其他 -> `self_contained`

说明：

- 这是过渡规则，后续会被 manifest/csv 索引规则替换。
- 判定逻辑集中在 `internal/dataset/service.go`，便于统一维护。

## 4. 后续演进

建议下一步新增：

1. `dataset_index.json`：统一样本索引
2. `manifest_l1.json`：明确 L1 样本清单
3. `dataset validate --strict`：校验命名、重复 ID、分类一致性

## 5. 当前实现

- 已提供示例清单：`configs/dataset_l1.json`
- Go 工具本地数据集目录：`datasets/`
- 当前目录可兼容你现在这类结构：`<lang>/<lang>_code_files_{self_contained|module_level}/<scenario>/*.ext`
- 执行优先级：`--dataset-manifest` > `--level` 对应默认清单 > 目录扫描
- 默认 level 清单命名约定：`configs/dataset_<level>.json`

推荐流程：

1. 原始文件只放 `datasets/`
2. 执行 `utbench dataset index` 生成索引（`configs/dataset_index.json`）
3. 执行 `utbench dataset manifest` 产出 L1/L2 清单
4. `run/generate/evaluate` 固定使用 manifest，避免目录扫描口径漂移

备注：

- `dataset_index.json` 为中间产物，通常不作为长期维护文件。
