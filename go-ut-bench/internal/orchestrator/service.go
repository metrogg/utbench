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

type Service struct {
	dataset   *dataset.Service
	runner    *runner.Service
	evaluator *evaluator.Service
	reporter  *reporter.Service
}

type Options struct {
	Ingest bool
	DBPath string
	// Phase controls which pipeline stage(s) to execute.
	// Supported: "full" (default), "generate", "evaluate", "report".
	Phase string
	// SourceRunID specifies the run ID to use as data source for evaluate/report phases.
	// If empty, uses the current RunID (which must have existing artifacts).
	SourceRunID string
	// ManifestPath overrides the default manifest path for evaluate phase.
	ManifestPath string
	// EvaluationPath overrides the default evaluation path for report phase.
	EvaluationPath string
}

type Result struct {
	RunID          string
	ManifestPath   string
	EvaluationPath string
	ReportJSONPath string
	ReportHTMLPath string
	Ingested       bool
}

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

	var evalOut *evaluator.Output

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
		evalOut = &out
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

	// Ingest only makes sense for full or evaluate phase
	ingested := false
	if opts.Ingest && (phase == "full" || phase == "evaluate") && evalOut != nil {
		sqliteStore, err := store.OpenSQLite(opts.DBPath)
		if err != nil {
			return Result{}, err
		}
		defer sqliteStore.Close()

		if err := sqliteStore.Init(ctx); err != nil {
			return Result{}, err
		}
		if err := sqliteStore.IngestEvaluation(ctx, evalOut.Result); err != nil {
			return Result{}, err
		}
		ingested = true
	}

	// Build result paths
	result := Result{
		RunID:     spec.RunID,
		Ingested:  ingested,
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
