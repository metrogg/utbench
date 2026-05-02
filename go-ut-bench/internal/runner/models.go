// runner/models.go 提供模型配置加载功能
// 从 YAML 文件解析模型配置，支持模型选择和过滤
package runner

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// modelConfig 单个模型的配置信息
// 包含 API 调用所需的所有参数
type modelConfig struct {
	Name      string         // 模型名称（配置文件中的键）
	Provider  string         // 提供商类型（openai、dashscope、volcengine 等）
	Endpoint  string         // API 端点 URL
	Model     string         // 实际调用的模型 ID
	APIKeyEnv string         // API 密钥环境变量名
	Enabled   bool           // 是否启用
	Params    map[string]any // 额外参数（temperature、max_tokens 等）
	Pricing   modelPricing   // 可选的成本配置
}

type modelPricing struct {
	PromptPer1KUSD     float64
	CompletionPer1KUSD float64
}

// modelsFile YAML 配置文件结构
// 对应 models.yaml 的顶层结构
type modelsFile struct {
	Models map[string]struct { // 模型配置映射（键为模型名称）
		Enabled  bool   `yaml:"enabled"`   // 是否启用
		Provider string `yaml:"provider"`  // 提供商类型
		Config   struct {
			APIEndpoint string         `yaml:"api_endpoint"` // API 端点
			Model       string         `yaml:"model"`        // 模型 ID
			APIKeyEnv   string         `yaml:"api_key_env"`  // 密钥环境变量
			Parameters  map[string]any `yaml:"parameters"`   // 额外参数
			Pricing     struct {
				PromptPer1KUSD     float64 `yaml:"prompt_per_1k_usd"`
				CompletionPer1KUSD float64 `yaml:"completion_per_1k_usd"`
			} `yaml:"pricing"`
		} `yaml:"config"`
	} `yaml:"models"`
}

// loadModelConfigs 从 YAML 文件加载模型配置
// 支持按名称筛选和启用状态过滤
//
// 参数:
//   - configPath: 配置文件路径
//   - selected: 选中的模型名称列表（空表示加载所有启用的模型）
//
// 返回值:
//   - []modelConfig: 模型配置列表
//   - error: 加载或验证错误
func loadModelConfigs(configPath string, selected []string) ([]modelConfig, error) {
	raw, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg modelsFile
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}
	if len(cfg.Models) == 0 {
		return nil, errors.New("models config is empty")
	}

	selectedSet := map[string]struct{}{}
	for _, item := range selected {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			selectedSet[trimmed] = struct{}{}
		}
	}

	out := make([]modelConfig, 0)
	for name, item := range cfg.Models {
		if len(selectedSet) > 0 {
			if _, ok := selectedSet[name]; !ok {
				continue
			}
		}
		if !item.Enabled {
			continue
		}
		if item.Config.APIKeyEnv == "" {
			return nil, fmt.Errorf("model %s missing api_key_env", name)
		}
		out = append(out, modelConfig{
			Name:      name,
			Provider:  strings.ToLower(strings.TrimSpace(item.Provider)),
			Endpoint:  strings.TrimSuffix(strings.TrimSpace(item.Config.APIEndpoint), "/"),
			Model:     strings.TrimSpace(item.Config.Model),
			APIKeyEnv: strings.TrimSpace(item.Config.APIKeyEnv),
			Enabled:   item.Enabled,
			Params:    item.Config.Parameters,
			Pricing: modelPricing{
				PromptPer1KUSD:     item.Config.Pricing.PromptPer1KUSD,
				CompletionPer1KUSD: item.Config.Pricing.CompletionPer1KUSD,
			},
		})
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("no enabled models selected: %v", selected)
	}
	return out, nil
}
