// contracts 包定义了跨模块使用的数据结构和常量
// 这些类型和常量在整个评测工具的不同组件之间共享使用
package contracts

import "fmt"

// DatasetClass 定义数据集的类别类型
// 用于区分不同类型的数据集样本
type DatasetClass string

// 数据集类别的常量定义
const (
	// DatasetClassSelfContained 表示自包含类型的数据集
	// 这种类型的样本代码不依赖外部模块，可以独立编译和运行
	DatasetClassSelfContained DatasetClass = "self_contained"
	// DatasetClassRepoLevel 表示仓库级别类型的数据集
	// 这种类型的样本可能依赖项目内的其他模块，需要完整的项目结构
	DatasetClassRepoLevel DatasetClass = "repo_level"
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

// SupportedScenarios 定义支持的评测场景列表
// 用于报告聚合和前端筛选，新增场景只需修改此处
var SupportedScenarios = []string{"simple_function", "boundary", "complex_dependency", "interface_mock"}

// ScenarioLabels 定义场景的中英文标签映射
// 用于报告展示，与 SupportedScenarios 保持同步
var ScenarioLabels = map[string]string{
	"simple_function":    "简单函数 / Simple",
	"boundary":           "边界值 / Boundary",
	"complex_dependency": "复杂依赖 / Complex",
	"interface_mock":     "接口 Mock / Interface",
	"unknown":            "未知 / Unknown",
}

// ScoreWeights 定义综合评分的权重结构
// 综合得分 = Compile×W.Compile + Test×W.Test + Coverage×W.Coverage + Mutation×W.Mutation
type ScoreWeights struct {
	Compile   float64 // 编译通过率权重
	Test      float64 // 样本测试通过率权重
	Coverage  float64 // 行覆盖率权重
	Mutation  float64 // 变异分数权重
}

// DefaultWeights 是默认评分权重
var DefaultWeights = ScoreWeights{
	Compile:  0.3,
	Test:     0.3,
	Coverage: 0.2,
	Mutation: 0.2,
}

// String 返回评分公式的可读描述
func (w ScoreWeights) String() string {
	return fmt.Sprintf("编译×%.1f + 样本测试×%.1f + 行覆盖率×%.1f + 变异分数×%.1f",
		w.Compile, w.Test, w.Coverage, w.Mutation)
}
