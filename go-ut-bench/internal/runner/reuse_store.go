package runner

import (
	"context"

	"go-ut-bench/internal/contracts"
)

// GenerationReuseStore 生成结果复用存储接口
// 用于解耦 runner 对 store 包的直接依赖
// 实现者负责根据 generationKey 查找可复用的已生成测试
type GenerationReuseStore interface {
	// FindReusableGeneratedAsset 根据 generationKey 查找可复用的生成结果
	// 返回值:
	//   - ReusableGeneratedCase: 可复用的结果（仅当 ok=true 时有效）
	//   - ok: 是否找到匹配的可复用结果
	//   - err: 查询过程中的错误
	FindReusableGeneratedAsset(ctx context.Context, generationKey string) (contracts.ReusableGeneratedCase, bool, error)
}
