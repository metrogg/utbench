package contracts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// GeneratedCase 表示单个测试生成案例的结果
// 记录了模型生成的单元测试的完整信息，包括文件路径、耗时、token使用情况等
type GeneratedCase struct {
	Model             string     `json:"model"`                       // 生成测试的模型名称，如"deepseek"、"qwen"等
	Language          string     `json:"language"`                    // 编程语言，如"python"、"go"、"java"等
	SampleID          string     `json:"sample_id"`                   // 数据集样本的唯一标识符
	SamplePath        string     `json:"sample_path"`                 // 原始源代码文件的路径
	GeneratedTestPath string     `json:"generated_test_path"`         // 生成的测试文件保存路径
	ResponsePath      string     `json:"response_path"`               // 模型API响应的JSON文件路径（用于调试）
	MetadataPath      string     `json:"metadata_path"`               // 元数据JSON文件路径
	LatencyMS         int        `json:"latency_ms"`                  // API调用耗时，单位毫秒
	PromptTokens      *int       `json:"prompt_tokens,omitempty"`     // 提示词token数量
	CompletionTokens  *int       `json:"completion_tokens,omitempty"` // 生成内容token数量
	TotalTokens       *int       `json:"total_tokens,omitempty"`      // 总token数量
	GeneratedAtUTC    time.Time  `json:"generated_at_utc"`            // 测试生成时间（UTC时间）
	Success           bool       `json:"success"`                     // 生成是否成功
	Truncated         bool       `json:"truncated,omitempty"`         // API响应是否因max_tokens截断（finish_reason="length")
	Error             *ErrorInfo `json:"error,omitempty"`             // 如果失败，记录错误详情
}

// GeneratedManifest 包含一组测试生成案例的清单
// 作为测试生成阶段的输出格式，包含完整的运行信息和所有案例列表
type GeneratedManifest struct {
	SchemaVersion string          `json:"schema_version"` // 数据结构版本号
	RunID         string          `json:"run_id"`         // 唯一运行ID，用于关联同一运行的所有数据
	CreatedAtUTC  time.Time       `json:"created_at_utc"` // 清单创建时间
	Spec          RunSpec         `json:"spec"`           // 运行时规格说明
	Cases         []GeneratedCase `json:"cases"`          // 所有生成的测试案例列表
}

// EvaluationResult 表示单个测试的评测结果
// 包含编译、运行、覆盖率、变异测试等全面的评测指标
type EvaluationResult struct {
	Model              string   `json:"model"`                         // 评测的模型名称
	Language           string   `json:"language"`                      // 编程语言
	SampleID           string   `json:"sample_id"`                     // 数据集样本ID
	GeneratedTestPath  string   `json:"generated_test_path"`           // 生成的测试文件路径
	SourcePath         string   `json:"source_path"`                   // 源代码文件路径
	CompilePass        bool     `json:"compile_pass"`                  // 编译是否通过
	TestPass           *bool    `json:"test_pass"`                     // 测试是否通过（nil表示未运行）
	Truncated          bool     `json:"truncated,omitempty"`           // API响应是否因max_tokens截断
	LineCoverage       *float64 `json:"line_coverage"`                 // 行覆盖率，范围0-100
	BranchCoverage     *float64 `json:"branch_coverage"`               // 分支覆盖率，范围0-100
	MutationScore      *float64 `json:"mutation_score"`                // 变异测试得分，范围0-100
	MutationTotal      *int     `json:"mutation_total,omitempty"`      // 变异体总数
	MutationKilled     *int     `json:"mutation_killed,omitempty"`     // 被杀死的变异体数量
	MutationSurvived   *int     `json:"mutation_survived,omitempty"`   // 存活的变异体数量
	MutationNoTests    *int     `json:"mutation_no_tests,omitempty"`   // 无法被测试检测的变异体数
	MutationTimeouts   *int     `json:"mutation_timeouts,omitempty"`   // 超时的变异体数量
	MutationSkipped    *int     `json:"mutation_skipped,omitempty"`    // 跳过的变异体数量
	MutationSuspicious *int     `json:"mutation_suspicious,omitempty"` // 可疑的变异体数量
	AssertionCount     *int     `json:"assertion_count"`               // 断言数量
	TestCaseCount      *int     `json:"test_case_count"`               // 测试用例数量
	AssertionDensity   *float64 `json:"assertion_density"`             // 断言密度（断言数/测试用例数）
	TestPassCount      *int     `json:"test_pass_count,omitempty"`     // 通过的测试用例数
	TestTotalCount     *int     `json:"test_total_count,omitempty"`    // 总测试用例数
	TestPassRate       *float64 `json:"test_pass_rate,omitempty"`      // 测试通过率
	RuntimeMS          *int     `json:"runtime_ms"`                    // 测试运行耗时（毫秒）
	PromptTokens       *int     `json:"prompt_tokens,omitempty"`       // 提示词token数量
	CompletionTokens   *int     `json:"completion_tokens,omitempty"`   // 生成token数量
	TotalTokens        *int     `json:"total_tokens,omitempty"`        // 总token数量
	CompileError       string   `json:"compile_error,omitempty"`       // 编译错误信息
	TestError          string   `json:"test_error,omitempty"`          // 测试运行错误信息
	CoverageError      string   `json:"coverage_error,omitempty"`      // 覆盖率收集错误信息
	MutationError      string   `json:"mutation_error,omitempty"`      // 变异测试错误信息
	MutationTool       string   `json:"mutation_tool,omitempty"`       // 使用的变异测试工具名称
}

type EvaluationResultSet struct {
	SchemaVersion  string             `json:"schema_version"`   // 数据结构版本号
	RunID          string             `json:"run_id"`           // 关联的运行ID
	EvaluatedAtUTC time.Time          `json:"evaluated_at_utc"` // 评测完成时间
	ManifestPath   string             `json:"manifest_path"`    // 关联的GeneratedManifest文件路径
	Results        []EvaluationResult `json:"results"`          // 所有评测结果列表
}

// ReportSummary 评测结果的汇总统计信息
// 用于快速了解整体评测效果，包含通过率、覆盖率等关键指标的平均值
type ReportSummary struct {
	TotalSamples        int     `json:"total_samples"`         // 总样本数量
	CompilePassCount    int     `json:"compile_pass_count"`    // 编译通过的样本数
	CompilePassRate     float64 `json:"compile_pass_rate"`     // 编译通过率（百分比）
	TestPassCount       int     `json:"test_pass_count"`       // 测试通过的样本数
	TestPassRate        float64 `json:"test_pass_rate"`        // 测试通过率（百分比）
	AvgLineCoverage     float64 `json:"avg_line_coverage"`     // 平均行覆盖率（百分比）
	AvgMutationScore    float64 `json:"avg_mutation_score"`    // 平均变异测试得分（百分比）
	AvgAssertionDensity float64 `json:"avg_assertion_density"` // 平均断言密度
}

// ReportPayload 报告的完整数据结构
// 包含汇总信息、多维度分析和失败案例详情，用于生成可视化报告
type ReportPayload struct {
	SchemaVersion    string             `json:"schema_version"`    // 数据结构版本号
	RunID            string             `json:"run_id"`            // 运行ID
	GeneratedAtUTC   time.Time          `json:"generated_at_utc"`  // 报告生成时间
	SourceEvaluation string             `json:"source_evaluation"` // 评测结果文件路径
	Summary          ReportSummary      `json:"summary"`           // 汇总统计
	Dimensions       Dimensions         `json:"dimensions"`        // 多维度分析数据
	TopModels        []ModelRank        `json:"top_models"`        // 模型排名列表
	ModelInfos       []ModelInfo        `json:"model_infos"`       // 模型详细信息（含型号）
	ByScenario       []ScenarioDim      `json:"by_scenario"`       // 按场景维度分析
	ByModelScenario  []ModelScenarioDim `json:"by_model_scenario"` // 模型+场景交叉分析
	TokenStats       TokenStats         `json:"token_stats"`       // Token 使用统计
	Failures         []FailureRow       `json:"failures"`          // 失败案例详情
	Thresholds       Thresholds         `json:"thresholds"`        // 评估阈值配置
	Prompts          map[string]string  `json:"prompts"`           // 按语言的提示词模板（key为语言，如"python"）
}

// Dimensions 多维度分析数据
// 从不同角度（如按模型、按语言、按场景）分析评测结果
type Dimensions struct {
	ByModel         []ModelDim         `json:"by_model"`          // 按模型维度的分析
	ByLanguage      []LanguageDim      `json:"by_language"`       // 按语言维度的分析
	ByScenario      []ScenarioDim      `json:"by_scenario"`       // 按场景维度的分析
	ByModelScenario []ModelScenarioDim `json:"by_model_scenario"` // 模型+场景交叉分析
}

// ModelDim 按模型维度的分析结果
// 统计特定模型在各指标上的表现
type ModelDim struct {
	Model               string  `json:"model"`                           // 模型名称
	ModelID             string  `json:"model_id,omitempty"`              // 具体型号（如 deepseek-chat）
	Provider            string  `json:"provider,omitempty"`              // 提供商
	TotalSamples        int     `json:"total_samples"`                   // 该模型的样本总数
	CompilePassRate     float64 `json:"compile_pass_rate"`               // 编译通过率
	AvgTestPassRate     float64 `json:"avg_test_pass_rate"`              // 平均测试通过率
	AvgLineCoverage     float64 `json:"avg_line_coverage"`               // 平均行覆盖率
	AvgBranchCoverage   float64 `json:"avg_branch_coverage"`             // 平均分支覆盖率
	AvgMutationScore    float64 `json:"avg_mutation_score"`              // 平均变异测试得分
	CompositeScore      float64 `json:"composite_score"`                 // 综合得分（加权：编译30%+测试30%+覆盖20%+变异20%）
	AvgLatencyMS        float64 `json:"avg_latency_ms,omitempty"`        // 平均API调用延迟
	AvgPromptTokens     float64 `json:"avg_prompt_tokens,omitempty"`     // 平均提示词Token
	AvgCompletionTokens float64 `json:"avg_completion_tokens,omitempty"` // 平均生成Token
	AvgTotalTokens      float64 `json:"avg_total_tokens,omitempty"`      // 平均总Token
	// 效率指标：每"通过样本"摊销的成本。通过样本 = 编译通过且测试通过率 > 0。
	// 数值越小越高效；无通过样本时为 0。
	TokensPerPass float64 `json:"tokens_per_pass,omitempty"` // 每通过样本平均 completion tokens
	MsPerPass     float64 `json:"ms_per_pass,omitempty"`     // 每通过样本平均耗时
	PassCount     int     `json:"pass_count,omitempty"`      // 通过样本数（用于效率分母）
}

// LanguageDim 按语言维度的分析结果
// 统计特定语言在各指标上的表现
type LanguageDim struct {
	Language          string  `json:"language"`            // 编程语言名称
	TotalSamples      int     `json:"total_samples"`       // 该语言的样本总数
	CompilePassRate   float64 `json:"compile_pass_rate"`   // 编译通过率
	AvgTestPassRate   float64 `json:"avg_test_pass_rate"`  // 平均测试通过率
	AvgLineCoverage   float64 `json:"avg_line_coverage"`   // 平均行覆盖率
	AvgBranchCoverage float64 `json:"avg_branch_coverage"` // 平均分支覆盖率
	AvgMutationScore  float64 `json:"avg_mutation_score"`  // 平均变异测试得分
}

// ModelRank 模型排名信息
// 用于展示模型的综合表现排名
type ModelRank struct {
	Rank                int     `json:"rank"`                            // 排名（1为最好）
	Model               string  `json:"model"`                           // 模型名称
	ModelID             string  `json:"model_id,omitempty"`              // 具体型号
	Provider            string  `json:"provider,omitempty"`              // 提供商
	CompilePassRate     float64 `json:"compile_pass_rate"`               // 编译通过率
	AvgTestPassRate     float64 `json:"avg_test_pass_rate"`              // 平均测试通过率
	AvgLineCoverage     float64 `json:"avg_line_coverage"`               // 平均行覆盖率
	AvgMutationScore    float64 `json:"avg_mutation_score"`              // 平均变异测试得分
	CompositeScore      float64 `json:"composite_score"`                 // 综合得分（加权）
	AvgLatencyMS        float64 `json:"avg_latency_ms,omitempty"`        // 平均延迟
	AvgPromptTokens     float64 `json:"avg_prompt_tokens,omitempty"`     // 平均提示词Token
	AvgCompletionTokens float64 `json:"avg_completion_tokens,omitempty"` // 平均生成Token
	AvgTotalTokens      float64 `json:"avg_total_tokens,omitempty"`      // 平均总Token
	TokensPerPass       float64 `json:"tokens_per_pass,omitempty"`       // 每通过样本平均 completion tokens
	MsPerPass           float64 `json:"ms_per_pass,omitempty"`           // 每通过样本平均耗时
	PassCount           int     `json:"pass_count,omitempty"`            // 通过样本数
}

// FailureRow 失败案例详情
// 记录特定类型错误的示例案例，用于问题诊断
type FailureRow struct {
	Stage          string         `json:"stage"`                     // 失败阶段："generate"、"evaluate"等
	ErrorType      string         `json:"error_type"`                // 错误类型："compile_error"、"test_error"等
	Count          int            `json:"count"`                     // 该类型错误的出现次数
	ByModel        map[string]int `json:"by_model,omitempty"`        // 按模型拆分的出现次数
	ExampleModel   string         `json:"example_model,omitempty"`   // 示例模型
	ExampleSample  string         `json:"example_sample,omitempty"`  // 示例样本ID
	ExampleMessage string         `json:"example_message,omitempty"` // 示例错误消息
}

// Thresholds 评估阈值配置
// 定义各项指标的及格线，用于判断模型是否达标
type Thresholds struct {
	CompilePassRate float64 `json:"compile_pass_rate"` // 编译通过率阈值
	TestPassRate    float64 `json:"test_pass_rate"`    // 测试通过率阈值
	LineCoverage    float64 `json:"line_coverage"`     // 行覆盖率阈值
	BranchCoverage  float64 `json:"branch_coverage"`   // 分支覆盖率阈值
	MutationScore   float64 `json:"mutation_score"`    // 变异测试得分阈值
}

// ModelInfo 模型详细信息
// 包含模型的名称、型号、Provider等完整信息
type ModelInfo struct {
	Name       string `json:"name"`        // 模型标识名，如 "deepseek"
	ModelID    string `json:"model_id"`    // 具体型号，如 "deepseek-chat"
	Provider   string `json:"provider"`    // 提供商，如 "deepseek", "dashscope"
	TotalCases int    `json:"total_cases"` // 该模型的评测样本数
}

// ScenarioDim 按场景维度的统计结果
// 统计特定场景（如 boundary、simple_function）在各指标上的表现
type ScenarioDim struct {
	Scenario          string  `json:"scenario"`            // 场景名称
	Language          string  `json:"language"`            // 编程语言
	TotalSamples      int     `json:"total_samples"`       // 该场景的样本总数
	CompilePassRate   float64 `json:"compile_pass_rate"`   // 编译通过率
	AvgTestPassRate   float64 `json:"avg_test_pass_rate"`  // 平均测试通过率
	AvgLineCoverage   float64 `json:"avg_line_coverage"`   // 平均行覆盖率
	AvgBranchCoverage float64 `json:"avg_branch_coverage"` // 平均分支覆盖率
	AvgMutationScore  float64 `json:"avg_mutation_score"`  // 平均变异测试得分
	AvgLatencyMS      float64 `json:"avg_latency_ms"`      // 平均耗时（毫秒）
	AvgTokens         float64 `json:"avg_tokens"`          // 平均 Token 使用量
}

// ModelScenarioDim 模型+场景交叉统计
// 统计特定模型在特定场景下的表现
type ModelScenarioDim struct {
	Model               string  `json:"model"`                 // 模型标识
	Scenario            string  `json:"scenario"`              // 场景名称
	Language            string  `json:"language"`              // 编程语言
	TotalSamples        int     `json:"total_samples"`         // 样本数
	CompilePassRate     float64 `json:"compile_pass_rate"`     // 编译通过率
	AvgTestPassRate     float64 `json:"avg_test_pass_rate"`    // 测试通过率
	AvgLineCoverage     float64 `json:"avg_line_coverage"`     // 行覆盖率
	AvgBranchCoverage   float64 `json:"avg_branch_coverage"`   // 分支覆盖率
	AvgMutationScore    float64 `json:"avg_mutation_score"`    // 变异得分
	AvgLatencyMS        float64 `json:"avg_latency_ms"`        // 平均耗时
	AvgPromptTokens     float64 `json:"avg_prompt_tokens"`     // 平均提示词Token
	AvgCompletionTokens float64 `json:"avg_completion_tokens"` // 平均生成Token
	AvgTotalTokens      float64 `json:"avg_total_tokens"`      // 平均总Token
}

// TokenStats Token 使用统计
type TokenStats struct {
	TotalPromptTokens     int     `json:"total_prompt_tokens"`     // 总提示词Token
	TotalCompletionTokens int     `json:"total_completion_tokens"` // 总生成Token
	TotalTokens           int     `json:"total_tokens"`            // 总Token
	AvgPromptTokens       float64 `json:"avg_prompt_tokens"`       // 平均提示词Token
	AvgCompletionTokens   float64 `json:"avg_completion_tokens"`   // 平均生成Token
	AvgTotalTokens        float64 `json:"avg_total_tokens"`        // 平均总Token
	SampleCount           int     `json:"sample_count"`            // 有Token记录的样本数
}

// NewRunID 生成一个新的唯一运行ID
// 返回值:
//   - string: 格式为"YYYYMMDDTHHMMSS.NANOSECONDSZ"的唯一ID
//
// 使用当前UTC时间生成，确保全球唯一性
func NewRunID() string {
	return time.Now().UTC().Format("20060102T150405.000000000Z")
}

// ReadGeneratedManifest 从指定路径读取GeneratedManifest文件
// 参数:
//   - path: JSON文件路径
//
// 返回值:
//   - GeneratedManifest: 解析后的清单结构
//   - error: 读取或解析失败时的错误
//
// 使用json.Unmarshal将文件内容反序列化为GeneratedManifest结构
func ReadGeneratedManifest(path string) (GeneratedManifest, error) {
	var payload GeneratedManifest
	if err := readJSON(path, &payload); err != nil {
		return GeneratedManifest{}, err
	}
	return payload, nil
}

// ReadEvaluationResultSet 从指定路径读取EvaluationResultSet文件
// 参数:
//   - path: JSON文件路径
//
// 返回值:
//   - EvaluationResultSet: 解析后的评测结果集
//   - error: 读取或解析失败时的错误
func ReadEvaluationResultSet(path string) (EvaluationResultSet, error) {
	var payload EvaluationResultSet
	if err := readJSON(path, &payload); err != nil {
		return EvaluationResultSet{}, err
	}
	return payload, nil
}

// WriteJSON 将任意数据结构写入JSON文件
// 参数:
//   - path: 目标文件路径（目录不存在时会自动创建）
//   - payload: 要序列化的数据结构
//
// 返回值:
//   - error: 写入失败时的错误
//
// 功能说明:
//   - 自动创建必要的目录结构（权限0o755）
//   - 使用json.MarshalIndent进行格式化输出（缩进2空格）
//   - 文件权限设置为0o644
func WriteJSON(path string, payload any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, raw, 0o644)
}

// readJSON 从JSON文件读取数据并反序列化
// 参数:
//   - path: JSON文件路径
//   - out: 目标结构指针
//
// 返回值:
//   - error: 读取或解析失败时的错误
//
// 内部使用的辅助函数，对外不可见
func readJSON(path string, out any) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("invalid json %s: %w", path, err)
	}
	return nil
}
