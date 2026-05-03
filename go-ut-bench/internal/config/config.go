// config 包提供应用程序配置相关的功能
// 负责定义默认配置值、配置验证以及配置项的默认值管理
package config

import (
	"errors"
	"strings"

	"go-ut-bench/internal/contracts"
)

// AppConfig 定义应用程序的默认配置结构
// 包含数据集路径、输出路径、数据库路径、默认模型、默认语言等配置项
// 这些值可以在CLI启动时被命令行参数覆盖
type AppConfig struct {
	DefaultDatasetRoot string                 // 数据集根目录路径，默认为"./datasets"
	DefaultOutputRoot  string                 // 评测结果输出根目录，默认为"./artifacts"
	DefaultDBPath      string                 // SQLite数据库文件路径，默认为"./storage/utbench.db"
	DefaultModels      []string               // 默认启用的模型列表，默认为["deepseek"]
	DefaultLanguages   []string               // 默认评测的编程语言列表，默认为["python"]
	DefaultClass       contracts.DatasetClass // 数据集类别，默认为self_contained（自包含模式）
	DefaultMode        contracts.RunMode      // 运行模式，默认为full（完整流程）
}

// Default 返回应用程序的默认配置
// 返回值:
//   - AppConfig: 包含所有默认值的配置对象
//
// 默认配置包含：
//   - 数据集根目录: ./datasets
//   - 输出根目录: ./artifacts
//   - 数据库路径: ./storage/utbench.db
//   - 默认模型: deepseek
//   - 默认语言: python
//   - 数据集类别: self_contained
//   - 运行模式: full
func Default() AppConfig {
	return AppConfig{
		DefaultDatasetRoot: "./datasets",
		DefaultOutputRoot:  "./artifacts",
		DefaultDBPath:      "./storage/utbench.db",
		DefaultModels:      []string{"deepseek"},
		DefaultLanguages:   []string{"python"},
		DefaultClass:       contracts.DatasetClassSelfContained,
		DefaultMode:        contracts.RunModeFull,
	}
}

// ValidateClass 验证数据集类别配置的有效性
// 参数:
//   - v: 数据集类别字符串，支持逗号分隔的多个值
//
// 返回值:
//   - error: 如果验证失败返回错误信息，否则返回nil
//
// 支持的数据集类别：
//   - self_contained: 自包含模式，代码文件不依赖外部模块
//   - repo_level: 仓库级别模式，代码文件可能依赖项目内的其他模块
//
// 示例：
//   - "self_contained" - 有效
//   - "repo_level" - 有效
//   - "self_contained,repo_level" - 有效
//   - "invalid" - 无效，返回错误
func ValidateClass(v string) error {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	validClasses := map[string]bool{
		"self_contained": true,
		"repo_level":     true,
	}
	for _, c := range strings.Split(v, ",") {
		c = strings.TrimSpace(c)
		if !validClasses[c] {
			return errors.New("dataset class must be self_contained or repo_level (comma-separated allowed)")
		}
	}
	return nil
}
