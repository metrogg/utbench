package contracts

import "time"

// RunSpec 定义评测运行的完整规格说明
// 包含所有影响运行行为的配置参数，用于追踪每次运行的配置
type RunSpec struct {
	RunID           string       `json:"run_id"`            // 唯一运行标识符，由NewRunID生成
	Models          []string     `json:"models"`            // 要评测的模型列表，如["deepseek", "qwen"]
	Languages       []string     `json:"languages"`        // 要评测的编程语言，如["python", "go"]
	DatasetClasses  []string     `json:"dataset_classes"`   // 数据集类别，如["self_contained", "module_level"]
	DatasetScenario string       `json:"dataset_scenario,omitempty"`    // 数据集场景过滤（可选），如"boundary"
	DatasetLevel    string       `json:"dataset_level,omitempty"`      // 数据集级别（可选），如"l1"
	DatasetRoot     string       `json:"dataset_root"`      // 数据集根目录路径
	DatasetManifest string       `json:"dataset_manifest,omitempty"`   // 数据集清单文件路径（可选）
	ConfigPath      string       `json:"config_path"`       // 模型配置文件路径
	Mode            RunMode      `json:"mode"`              // 运行模式：full或incremental
	DryRun          bool         `json:"dry_run"`           // 是否为试运行模式（不调用真实API）
	ResetCheckpoint bool         `json:"reset_checkpoint"`  // 是否重置checkpoint，强制重新运行
	MutationEnabled bool         `json:"mutation_enabled"`  // 是否启用变异测试
	MutationTimeout int          `json:"mutation_timeout_seconds"` // 变异测试超时时间（秒）
	MutationPolicy  string       `json:"mutation_error_policy"`     // 变异测试错误策略："skip"或"fail"
	TestTimeout     int          `json:"test_timeout_seconds"`       // 测试执行超时时间（秒）
	MaxSamples      int          `json:"max_samples"`       // 最大样本数量限制（0表示不限制）
	Workers         int          `json:"workers"`           // 并发worker数量（0表示使用默认值）
	OutputRoot      string       `json:"output_root"`       // 输出根目录
	CreatedAtUTC    time.Time    `json:"created_at_utc"`    // 运行创建时间（UTC）
}

// SampleRef 数据集样本的引用信息
// 用于唯一标识和定位数据集中的样本
type SampleRef struct {
	ID        string       `json:"id"`        // 样本唯一标识符，格式为"<场景>_<编号>"
	Language  string       `json:"language"`  // 编程语言，如"python"、"go"
	Category  DatasetClass `json:"category"`  // 数据集类别：self_contained或module_level
	Scenario  string       `json:"scenario"`  // 场景名称，如"boundary"、"simple_function"
	Path      string       `json:"path"`      // 样本文件的绝对路径
	SourceMD5 string       `json:"source_md5"` // 源代码的MD5哈希值，用于校验
}

// ModuleLevelMeta 模块级别样本的元数据
// 包含module_level类型样本的额外信息，用于正确执行和导入
type ModuleLevelMeta struct {
	SampleID       string   `json:"sample_id"`        // 样本ID
	ModuleImport   string   `json:"module_import"`    // Python模块导入路径，如"dateutil.parser"
	PackageName    string   `json:"package_name"`     // 包名，如"dateutil"
	TargetFile     string   `json:"target_file"`     // 待测试的目标文件，相对于workspace_root
	WorkspaceRoot  string   `json:"workspace_root"`  // 工作区根目录
	Requirements   []string `json:"requirements,omitempty"` // 依赖的Python包列表
}
