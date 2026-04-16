package contracts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type GeneratedCase struct {
	Model             string     `json:"model"`
	Language          string     `json:"language"`
	SampleID          string     `json:"sample_id"`
	SamplePath        string     `json:"sample_path"`
	GeneratedTestPath string     `json:"generated_test_path"`
	ResponsePath      string     `json:"response_path"`
	MetadataPath      string     `json:"metadata_path"`
	LatencyMS         int        `json:"latency_ms"`
	PromptTokens      *int       `json:"prompt_tokens,omitempty"`
	CompletionTokens  *int       `json:"completion_tokens,omitempty"`
	TotalTokens       *int       `json:"total_tokens,omitempty"`
	GeneratedAtUTC    time.Time  `json:"generated_at_utc"`
	Success           bool       `json:"success"`
	Error             *ErrorInfo `json:"error,omitempty"`
}

type GeneratedManifest struct {
	SchemaVersion string          `json:"schema_version"`
	RunID         string          `json:"run_id"`
	CreatedAtUTC  time.Time       `json:"created_at_utc"`
	Spec          RunSpec         `json:"spec"`
	Cases         []GeneratedCase `json:"cases"`
}

type EvaluationResult struct {
	Model              string   `json:"model"`
	Language           string   `json:"language"`
	SampleID           string   `json:"sample_id"`
	GeneratedTestPath  string   `json:"generated_test_path"`
	SourcePath         string   `json:"source_path"`
	CompilePass        bool     `json:"compile_pass"`
	TestPass           *bool    `json:"test_pass"`
	LineCoverage       *float64 `json:"line_coverage"`
	BranchCoverage     *float64 `json:"branch_coverage"`
	MutationScore      *float64 `json:"mutation_score"`
	MutationTotal      *int     `json:"mutation_total,omitempty"`
	MutationKilled     *int     `json:"mutation_killed,omitempty"`
	MutationSurvived   *int     `json:"mutation_survived,omitempty"`
	MutationNoTests    *int     `json:"mutation_no_tests,omitempty"`
	MutationTimeouts   *int     `json:"mutation_timeouts,omitempty"`
	MutationSkipped    *int     `json:"mutation_skipped,omitempty"`
	MutationSuspicious *int     `json:"mutation_suspicious,omitempty"`
	AssertionCount     *int     `json:"assertion_count"`
	TestCaseCount      *int     `json:"test_case_count"`
	AssertionDensity   *float64 `json:"assertion_density"`
	TestPassCount      *int     `json:"test_pass_count,omitempty"`
	TestTotalCount     *int     `json:"test_total_count,omitempty"`
	TestPassRate       *float64 `json:"test_pass_rate,omitempty"`
	RuntimeMS          *int     `json:"runtime_ms"`
	CompileError       string   `json:"compile_error,omitempty"`
	TestError          string   `json:"test_error,omitempty"`
	CoverageError      string   `json:"coverage_error,omitempty"`
	MutationError      string   `json:"mutation_error,omitempty"`
}

type EvaluationResultSet struct {
	SchemaVersion  string             `json:"schema_version"`
	RunID          string             `json:"run_id"`
	EvaluatedAtUTC time.Time          `json:"evaluated_at_utc"`
	ManifestPath   string             `json:"manifest_path"`
	Results        []EvaluationResult `json:"results"`
}

type ReportSummary struct {
	TotalSamples        int     `json:"total_samples"`
	CompilePassCount    int     `json:"compile_pass_count"`
	CompilePassRate     float64 `json:"compile_pass_rate"`
	TestPassCount       int     `json:"test_pass_count"`
	TestPassRate        float64 `json:"test_pass_rate"`
	AvgLineCoverage     float64 `json:"avg_line_coverage"`
	AvgMutationScore    float64 `json:"avg_mutation_score"`
	AvgAssertionDensity float64 `json:"avg_assertion_density"`
}

type ReportPayload struct {
	SchemaVersion    string        `json:"schema_version"`
	RunID            string        `json:"run_id"`
	GeneratedAtUTC   time.Time     `json:"generated_at_utc"`
	SourceEvaluation string        `json:"source_evaluation"`
	Summary          ReportSummary `json:"summary"`
	Dimensions       Dimensions    `json:"dimensions"`
	TopModels        []ModelRank   `json:"top_models"`
	Failures         []FailureRow  `json:"failures"`
	Thresholds       Thresholds    `json:"thresholds"`
}

type Dimensions struct {
	ByModel    []ModelDim    `json:"by_model"`
	ByLanguage []LanguageDim `json:"by_language"`
}

type ModelDim struct {
	Model             string  `json:"model"`
	TotalSamples      int     `json:"total_samples"`
	CompilePassRate   float64 `json:"compile_pass_rate"`
	AvgTestPassRate   float64 `json:"avg_test_pass_rate"`
	AvgLineCoverage   float64 `json:"avg_line_coverage"`
	AvgBranchCoverage float64 `json:"avg_branch_coverage"`
	AvgMutationScore  float64 `json:"avg_mutation_score"`
	AvgLatencyMS      float64 `json:"avg_latency_ms,omitempty"`
	AvgTokens         float64 `json:"avg_tokens,omitempty"`
}

type LanguageDim struct {
	Language          string  `json:"language"`
	TotalSamples      int     `json:"total_samples"`
	CompilePassRate   float64 `json:"compile_pass_rate"`
	AvgTestPassRate   float64 `json:"avg_test_pass_rate"`
	AvgLineCoverage   float64 `json:"avg_line_coverage"`
	AvgBranchCoverage float64 `json:"avg_branch_coverage"`
	AvgMutationScore  float64 `json:"avg_mutation_score"`
}

type ModelRank struct {
	Rank             int     `json:"rank"`
	Model            string  `json:"model"`
	AvgTestPassRate  float64 `json:"avg_test_pass_rate"`
	AvgLineCoverage  float64 `json:"avg_line_coverage"`
	AvgMutationScore float64 `json:"avg_mutation_score"`
	AvgLatencyMS     float64 `json:"avg_latency_ms,omitempty"`
	AvgTokens        float64 `json:"avg_tokens,omitempty"`
}

type FailureRow struct {
	Stage          string `json:"stage"`
	ErrorType      string `json:"error_type"`
	Count          int    `json:"count"`
	ExampleModel   string `json:"example_model,omitempty"`
	ExampleSample  string `json:"example_sample,omitempty"`
	ExampleMessage string `json:"example_message,omitempty"`
}

type Thresholds struct {
	CompilePassRate float64 `json:"compile_pass_rate"`
	TestPassRate    float64 `json:"test_pass_rate"`
	LineCoverage    float64 `json:"line_coverage"`
	BranchCoverage  float64 `json:"branch_coverage"`
	MutationScore   float64 `json:"mutation_score"`
}

func NewRunID() string {
	return time.Now().UTC().Format("20060102T150405.000000000Z")
}

func ReadGeneratedManifest(path string) (GeneratedManifest, error) {
	var payload GeneratedManifest
	if err := readJSON(path, &payload); err != nil {
		return GeneratedManifest{}, err
	}
	return payload, nil
}

func ReadEvaluationResultSet(path string) (EvaluationResultSet, error) {
	var payload EvaluationResultSet
	if err := readJSON(path, &payload); err != nil {
		return EvaluationResultSet{}, err
	}
	return payload, nil
}

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
