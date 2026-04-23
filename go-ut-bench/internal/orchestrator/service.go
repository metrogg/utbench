package orchestrator

import (
	"context"
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

	evalOut, err := s.evaluator.Evaluate(ctx, spec, genOut.ManifestPath)
	if err != nil {
		return Result{}, err
	}

	reportOut, err := s.reporter.Generate(ctx, spec, evalOut.ResultPath)
	if err != nil {
		return Result{}, err
	}

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
		if err := sqliteStore.IngestEvaluation(ctx, evalOut.Result); err != nil {
			return Result{}, err
		}
		ingested = true
	}

	runSummaryPath := filepath.Join(spec.OutputRoot, "runs", spec.RunID, "run_summary.json")
	_ = contracts.WriteJSON(runSummaryPath, map[string]any{
		"schema_version": contracts.SchemaVersion,
		"run_id":         spec.RunID,
		"created_at_utc": time.Now().UTC(),
		"spec":           spec,
		"artifacts": map[string]string{
			"manifest":    genOut.ManifestPath,
			"evaluation":  evalOut.ResultPath,
			"report_json": reportOut.ReportJSONPath,
			"report_html": reportOut.ReportHTMLPath,
		},
		"ingested": ingested,
		"db_path":  opts.DBPath,
	})

	return Result{
		RunID:          spec.RunID,
		ManifestPath:   genOut.ManifestPath,
		EvaluationPath: evalOut.ResultPath,
		ReportJSONPath: reportOut.ReportJSONPath,
		ReportHTMLPath: reportOut.ReportHTMLPath,
		Ingested:       ingested,
	}, nil
}
