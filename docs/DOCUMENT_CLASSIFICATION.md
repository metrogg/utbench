# 外层 docs 分类说明

外层 `docs/` 已经降级，不再作为当前实现的主文档树。

## 当前规则

- 当前实现、使用方式、开发说明：看 `go-ut-bench/docs/`
- 研究资料：通过 `go-ut-bench/docs/06-research/README.md` 进入
- 外部工具资料：通过 `go-ut-bench/docs/05-integrations/README.md` 进入
- 旧分析和历史规格：通过 `go-ut-bench/docs/99-archive/legacy-root-docs.md` 进入

## 为什么不直接删掉

- 这些资料里有研究背景和外部资料，直接删除会损失上下文。
- 在实体文件完全迁移前，先通过统一导航去重，避免继续维护两套入口。

## 后续目标

- 最终只保留 `go-ut-bench/docs/` 作为唯一文档树。
- 外层 `docs/` 在实体迁移完成后，可以进一步缩成纯占位目录，或者完全删除。
