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

CLI 默认会从 `./datasets` 读取样本。
