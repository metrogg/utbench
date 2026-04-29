// orchestrator 包提供评测流程的编排功能
// 负责协调数据集发现、测试生成、评测执行、报告生成等阶段的执行
// 支持分阶段执行（generate/evaluate/report）和完整流水线运行
package orchestrator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go-ut-bench/internal/contracts"
	"go-ut-bench/internal/dataset"
	"go-ut-bench/internal/evaluator"
	"go-ut-bench/internal/reporter"
	"go-ut-bench/internal/runner"
	"go-ut-bench/internal/store"
)

// Service 编排服务
// 持有各子服务的引用，负责协调流水线各阶段的执行
type Service struct {
	dataset   *dataset.Service   // 数据集发现服务
	runner    *runner.Service    // 测试生成服务
	evaluator *evaluator.Service // 评测执行服务
	reporter  *reporter.Service  // 报告生成服务
}

// Options 编排运行选项
// 控制流水线的执行阶段和数据入库行为
type Options struct {
	Ingest bool // 是否将结果入库 SQLite 数据库
	DBPath string // SQLite 数据库路径
	// Phase 控制执行的流水线阶段
	// 支持值: "full"（默认，完整流水线）、"generate"（仅生成）、"evaluate"（仅评测）、"report"（仅报告）
	Phase string
	// SourceRunID 指定作为数据源的运行ID，用于 evaluate/report 阶段
	// 如果为空，使用当前 RunID（必须已有 artifacts）
	SourceRunID string
	// ManifestPath 覆盖 evaluate 阶段的默认 manifest 路径
	ManifestPath string
	// EvaluationPath 覆盖 report 阶段的默认 evaluation 结果路径
	EvaluationPath string
}

// Result 编排运行结果
// 包含各阶段输出的文件路径和入库状态
type Result struct {
	RunID          string // 运行唯一标识符
	ManifestPath   string // 生成清单文件路径
	EvaluationPath string // 评测结果文件路径
	ReportJSONPath string // 报告 JSON 文件路径
	ReportHTMLPath string // 报告 HTML 文件路径
	Ingested       bool   // 是否已入库数据库
}

// New 创建编排服务实例
// 参数:
//   - datasetSvc: 数据集发现服务
//   - runnerSvc: 测试生成服务
//   - evaluatorSvc: 评测执行服务
//   - reporterSvc: 报告生成服务
//
// 返回值:
//   - *Service: 编排服务实例
func New(
	datasetSvc *dataset.Service,
	runnerSvc *runner.Service,
	evaluatorSvc *evaluator.Service,
	reporterSvc *reporter.Service,
) *Service {
	return &Service{
		dataset:   datasetSvc,
		runner:    runnerSvc,
		evaluator: evaluatorSvc,
		reporter:  reporterSvc,
	}
}

// Run 执行评测流水线
// 根据 opts.Phase 参数执行对应阶段，支持分阶段执行或完整流水线
//
// 参数:
//   - ctx: 上下文，用于取消操作
//   - spec: 运行规格说明，包含模型、语言、数据集等配置
//   - opts: 编排选项，控制执行阶段和入库行为
//
// 返回值:
//   - Result: 运行结果，包含各阶段输出文件路径
//   - error: 执行过程中的错误
//
// 执行阶段（由 opts.Phase 控制）:
//  1. "generate" 或 "full": 数据集发现 → 测试生成
//  2. "evaluate" 或 "full": 编译 → 测试执行 → 覆盖率 → 变异测试
//  3. "report" 或 "full": 多维度聚合 → HTML报告生成
//  4. 如果 opts.Ingest=true: 将结果入库 SQLite
func (s *Service) Run(ctx context.Context, spec contracts.RunSpec, opts Options) (Result, error) {
	phase := opts.Phase
	if phase == "" {
		phase = "full"
	}

	// Determine the source run ID for artifact paths
	sourceRunID := opts.SourceRunID
	if sourceRunID == "" {
		sourceRunID = spec.RunID
	}

	// Build default artifact paths based on source run ID
	sourceRunDir := filepath.Join(spec.OutputRoot, "runs", sourceRunID)
	defaultManifestPath := filepath.Join(sourceRunDir, "generated", "generated_manifest.json")
	defaultEvaluationPath := filepath.Join(sourceRunDir, "evaluation", "evaluation_result.json")

	// Use explicit paths if provided, otherwise use defaults
	manifestPath := opts.ManifestPath
	if manifestPath == "" && (phase == "evaluate" || phase == "report") {
		manifestPath = defaultManifestPath
	}
	evaluationPath := opts.EvaluationPath
	if evaluationPath == "" && phase == "report" {
		evaluationPath = defaultEvaluationPath
	}

	// Phase "generate": discover samples and generate tests only.
	if phase == "generate" || phase == "full" {
		if err := s.dataset.ValidateSpec(spec); err != nil {
			return Result{}, err
		}
		samples, err := s.dataset.DiscoverSamples(spec)
		if err != nil {
			return Result{}, err
		}
		genOut, err := s.runner.Generate(ctx, spec, samples)
		if err != nil {
			return Result{}, err
		}
		manifestPath = genOut.ManifestPath
	}

	// Phase "evaluate": requires existing manifest; evaluate only.
	if phase == "evaluate" {
		// Verify manifest exists
		if !fileExists(manifestPath) {
			return Result{}, fmt.Errorf("manifest not found for evaluate phase: %s", manifestPath)
		}
	}

	// Phase "evaluate" or continuing from generate in full mode
	if phase == "evaluate" || phase == "full" {
		out, err := s.evaluator.Evaluate(ctx, spec, manifestPath)
		if err != nil {
			return Result{}, err
		}
		evaluationPath = out.ResultPath
	}

	// Phase "report": requires existing evaluation JSON; generate report only.
	if phase == "report" {
		if !fileExists(evaluationPath) {
			return Result{}, fmt.Errorf("evaluation JSON not found for report phase: %s", evaluationPath)
		}
	}

	// Phase "report" or continuing from evaluate in full mode
	var reportOut *reporter.Output
	if phase == "report" || phase == "full" {
		out, err := s.reporter.Generate(ctx, spec, evaluationPath)
		if err != nil {
			return Result{}, err
		}
		reportOut = &out
	}

	// Ingest the run directory into the v2 SQLite store. The store indexes the
	// manifest, evaluation result, report and linked artifacts when present.
	ingested := false
	if opts.Ingest {
		sqliteStore, err := store.OpenSQLite(opts.DBPath)
		if err != nil {
			return Result{}, err
		}
		defer sqliteStore.Close()

		if err := sqliteStore.Init(ctx); err != nil {
			return Result{}, err
		}
		runDir := filepath.Join(spec.OutputRoot, "runs", spec.RunID)
		if _, err := sqliteStore.IngestRun(ctx, store.IngestRunOptions{RunDir: runDir}); err != nil {
			return Result{}, err
		}
		ingested = true
	}

	// Build result paths
	result := Result{
		RunID:    spec.RunID,
		Ingested: ingested,
	}
	if manifestPath != "" {
		result.ManifestPath = manifestPath
	}
	if evaluationPath != "" {
		result.EvaluationPath = evaluationPath
	}
	if reportOut != nil {
		result.ReportJSONPath = reportOut.ReportJSONPath
		result.ReportHTMLPath = reportOut.ReportHTMLPath
	}

	// Write run_summary.json (use current run's directory, not source)
	runDir := filepath.Join(spec.OutputRoot, "runs", spec.RunID)
	runSummaryPath := filepath.Join(runDir, "run_summary.json")
	summaryData := map[string]any{
		"schema_version": contracts.SchemaVersion,
		"run_id":         spec.RunID,
		"created_at_utc": time.Now().UTC(),
		"spec":           spec,
		"phase":          phase,
		"ingested":       ingested,
		"db_path":        opts.DBPath,
	}
	if result.ManifestPath != "" {
		summaryData["manifest_path"] = result.ManifestPath
	}
	if result.EvaluationPath != "" {
		summaryData["evaluation_path"] = result.EvaluationPath
	}
	if result.ReportJSONPath != "" {
		summaryData["report_json_path"] = result.ReportJSONPath
	}
	if result.ReportHTMLPath != "" {
		summaryData["report_html_path"] = result.ReportHTMLPath
	}
	_ = contracts.WriteJSON(runSummaryPath, summaryData)

	return result, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
