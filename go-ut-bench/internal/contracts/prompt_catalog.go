package contracts

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// PromptMode 提示词模式类型
type PromptMode string

// 三种提示词模式常量
const (
	PromptModeFullFile   PromptMode = "full_file"  // 完整文件模式
	PromptModeCompletion PromptMode = "completion"  // 续写模式
	PromptModeRepoLevel  PromptMode = "repo_level"  // 仓库级模式
)

// PromptCatalog 提示词目录
// 包含版本信息、系统消息和各语言各模式的提示词模板
type PromptCatalog struct {
	Strategy      string                           `json:"strategy"`
	VersionID     string                           `json:"version_id"`
	SystemMessage string                           `json:"system_message"`
	Modes         []PromptMode                     `json:"modes"`
	Templates     map[string]map[PromptMode]string `json:"templates"`
}

// LoadPromptCatalog 从目录加载提示词目录
// 读取 prompt_catalog.json 文件
func LoadPromptCatalog(dir string) (PromptCatalog, error) {
	var catalog PromptCatalog
	raw, err := os.ReadFile(filepath.Join(dir, "prompt_catalog.json"))
	if err != nil {
		return PromptCatalog{}, err
	}
	if err := json.Unmarshal(raw, &catalog); err != nil {
		return PromptCatalog{}, err
	}
	return catalog, nil
}
