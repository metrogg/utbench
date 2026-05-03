// dataset/manifest.go 提供数据集清单加载功能
// 从 JSON 文件解析样本列表，处理路径规范化
package dataset

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"go-ut-bench/internal/contracts"
)

// datasetManifest 数据集清单结构
// 定义清单文件格式，包含级别和样本列表
type datasetManifest struct {
	Level   string `json:"level"`   // 数据集级别（如 l1、l2）
	Samples []struct {              // 样本列表
		ID       string                 `json:"id"`       // 样本唯一标识
		Language string                 `json:"language"` // 编程语言
		Category contracts.DatasetClass `json:"category"` // 数据集类别
		Scenario string                 `json:"scenario"` // 场景名称
		Path     string                 `json:"path"`     // 样本文件路径
	} `json:"samples"`
}

// loadManifest 加载并解析数据集清单文件
// 处理路径规范化，支持相对路径转换为绝对路径
//
// 参数:
//   - path: 清单文件路径
//   - datasetRoot: 数据集根目录
//
// 返回值:
//   - datasetManifest: 解析后的清单
//   - error: 加载错误
func loadManifest(path string, datasetRoot string) (datasetManifest, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return datasetManifest{}, err
	}
	var mf datasetManifest
	if err := json.Unmarshal(raw, &mf); err != nil {
		return datasetManifest{}, err
	}
	for i := range mf.Samples {
		sample := &mf.Samples[i]
		if sample.Path == "" {
			sample.Path = filepath.Join(sample.Language, sample.ID+languageExtByName(sample.Language))
		}
		if !filepath.IsAbs(sample.Path) {
			candidate := filepath.Clean(sample.Path)
			if _, err := os.Stat(candidate); err == nil {
				sample.Path = candidate
			} else {
				sample.Path = filepath.Join(datasetRoot, sample.Path)
			}
		}
		sample.Path = filepath.Clean(sample.Path)
		sample.Language = strings.ToLower(strings.TrimSpace(sample.Language))
		sample.Scenario = strings.TrimSpace(sample.Scenario)
	}
	return mf, nil
}
