# datasets

go-ut-bench 本地数据集目录。

当前目录约定：

- `datasets/python`
- `datasets/java`
- `datasets/go`
- `datasets/cpp`

命名建议：

- 自包含样本：`boundary_xxx`、`simple_function_xxx`
- 复杂依赖样本：`complex_dependency_xxx`

标准结构（建议）：

- `datasets/<lang>/<lang>_code_files_self_contained/<scenario>/<sample_id>.<ext>`
- `datasets/<lang>/<lang>_code_files_module_level/<scenario>/<sample_id>.<ext>`

其中：

- `<scenario>` 固定为：`boundary` / `simple_function` / `complex_dependency` / `interface_mock`
- `<sample_id>` 建议统一为：`<scenario>_<index3>`（如 `simple_function_000`）

module_level（预恢复样本）管理建议：

- 源码入口文件仍按上述位置落在 `datasets/` 下（用于 runner/evaluator 的样本引用）
- 同时在样本元数据里维护 `module_import`、`workspace_root`、`requirements` 等信息
- 不建议把样本 ID 直接命名成 `dateutil`，应保留统一编号和场景前缀

CLI 默认会从 `./datasets` 读取样本。
