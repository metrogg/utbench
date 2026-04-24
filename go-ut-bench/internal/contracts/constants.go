// contracts 包定义了跨模块使用的数据结构和常量
// 这些类型和常量在整个评测工具的不同组件之间共享使用
package contracts

// DatasetClass 定义数据集的类别类型
// 用于区分不同类型的数据集样本
type DatasetClass string

// 数据集类别的常量定义
const (
	// DatasetClassSelfContained 表示自包含类型的数据集
	// 这种类型的样本代码不依赖外部模块，可以独立编译和运行
	DatasetClassSelfContained DatasetClass = "self_contained"
	// DatasetClassModuleLevel 表示模块级别类型的数据集
	// 这种类型的样本可能依赖项目内的其他模块，需要完整的项目结构
	DatasetClassModuleLevel DatasetClass = "module_level"
)

// RunMode 定义评测工具的运行模式
type RunMode string

// 运行模式的常量定义
const (
	// RunModeFull 表示完整运行模式
	// 会执行完整的评测流程：生成测试 -> 评测 -> 报告
	RunModeFull RunMode = "full"
	// RunModeIncremental 表示增量运行模式
	// 支持从上次中断处继续执行，使用checkpoint机制保存进度
	RunModeIncremental RunMode = "incremental"
)

// SchemaVersion 定义当前数据结构的版本号
// 用于版本控制和兼容性检查
const SchemaVersion = "v0.1.0"

// SupportedLanguages 定义支持的编程语言列表
// 当前支持：Python、Java、Go、C++ 四种语言
var SupportedLanguages = []string{"python", "java", "go", "cpp"}
