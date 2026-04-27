package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"go-ut-bench/internal/contracts"

	"gopkg.in/yaml.v3"
	_ "modernc.org/sqlite"
)

const schemaVersion = "store.v2"

type SQLiteStore struct {
	db   *sql.DB
	path string
}

type DBOverview struct {
	DBPath            string      `json:"db_path"`
	SchemaVersion     string      `json:"schema_version"`
	GenerationRuns    int         `json:"generation_runs"`
	EvaluationRuns    int         `json:"evaluation_runs"`
	GeneratedCases    int         `json:"generated_cases"`
	EvaluationResults int         `json:"evaluation_results"`
	Artifacts         int         `json:"artifacts"`
	Reports           int         `json:"reports"`
	LatestRuns        []DBRunItem `json:"latest_runs"`
}

type DBRunItem struct {
	RunID             string `json:"run_id"`
	ExperimentID      string `json:"experiment_id,omitempty"`
	ExperimentName    string `json:"experiment_name,omitempty"`
	CreatedAtUTC      string `json:"created_at_utc,omitempty"`
	EvaluatedAtUTC    string `json:"evaluated_at_utc,omitempty"`
	Models            int    `json:"models"`
	Languages         int    `json:"languages"`
	GeneratedCases    int    `json:"generated_cases"`
	EvaluationResults int    `json:"evaluation_results"`
	Reports           int    `json:"reports"`
}

type DBResultItem struct {
	EvaluationResultID   string   `json:"evaluation_result_id"`
	EvaluationRunID      string   `json:"evaluation_run_id"`
	RunID                string   `json:"run_id"`
	Model                string   `json:"model"`
	Language             string   `json:"language"`
	SampleID             string   `json:"sample_id"`
	CompilePass          bool     `json:"compile_pass"`
	TestPass             *bool    `json:"test_pass,omitempty"`
	LineCoverage         *float64 `json:"line_coverage,omitempty"`
	MutationScore        *float64 `json:"mutation_score,omitempty"`
	MutationTotal        *int     `json:"mutation_total,omitempty"`
	FailureOrigin        string   `json:"failure_origin,omitempty"`
	ScoreEligible        bool     `json:"score_eligible"`
	ScoreExclusionReason string   `json:"score_exclusion_reason,omitempty"`
	RuntimeMS            *int     `json:"runtime_ms,omitempty"`
}

type DBArtifactItem struct {
	ArtifactID   string `json:"artifact_id"`
	Kind         string `json:"kind"`
	Path         string `json:"path"`
	SizeBytes    int64  `json:"size_bytes"`
	SHA256       string `json:"sha256"`
	Redacted     bool   `json:"redacted"`
	CreatedAtUTC string `json:"created_at_utc"`
	DeletedAtUTC string `json:"deleted_at_utc,omitempty"`
}

type DBReportFilter struct {
	RunIDs            []string `json:"run_ids,omitempty"`
	EvaluationRunIDs  []string `json:"evaluation_run_ids,omitempty"`
	Models            []string `json:"models,omitempty"`
	Languages         []string `json:"languages,omitempty"`
	ScoreEligibleOnly bool     `json:"score_eligible_only,omitempty"`
}

type DBReportFacets struct {
	Runs      []DBReportRunOption `json:"runs"`
	Models    []string            `json:"models"`
	Languages []string            `json:"languages"`
	EnvGroups []DBReportEnvOption `json:"env_groups"`
}

type DBReportRunOption struct {
	RunID             string `json:"run_id"`
	EvaluationRunID   string `json:"evaluation_run_id"`
	EvaluatedAtUTC    string `json:"evaluated_at_utc"`
	EnvID             string `json:"env_id"`
	PromptStrategy    string `json:"prompt_strategy,omitempty"`
	PromptVersionID   string `json:"prompt_version_id,omitempty"`
	EvaluationResults int    `json:"evaluation_results"`
}

type DBReportEnvOption struct {
	EnvID             string `json:"env_id"`
	Fingerprint       string `json:"fingerprint"`
	EvaluationRuns    int    `json:"evaluation_runs"`
	EvaluationResults int    `json:"evaluation_results"`
}

// Phase 3: Database management list item types

type DBGenerationRunItem struct {
	RunID              string `json:"run_id"`
	ExperimentID       string `json:"experiment_id,omitempty"`
	SchemaVersion      string `json:"schema_version"`
	CreatedAtUTC       string `json:"created_at_utc"`
	PromptStrategy     string `json:"prompt_strategy,omitempty"`
	PromptVersionID    string `json:"prompt_version_id,omitempty"`
	DatasetFingerprint string `json:"dataset_fingerprint,omitempty"`
	CreatedDBAtUTC     string `json:"created_db_at_utc"`
}

type DBGeneratedCaseItem struct {
	GeneratedCaseID  string `json:"generated_case_id"`
	RunID            string `json:"run_id"`
	Model            string `json:"model"`
	Language         string `json:"language"`
	SampleID         string `json:"sample_id"`
	Success          bool   `json:"success"`
	Truncated        bool   `json:"truncated"`
	LatencyMS        *int   `json:"latency_ms,omitempty"`
	PromptTokens     *int   `json:"prompt_tokens,omitempty"`
	CompletionTokens *int   `json:"completion_tokens,omitempty"`
	GeneratedAtUTC   string `json:"generated_at_utc,omitempty"`
}

type DBPromptRenderingItem struct {
	PromptRenderingID string `json:"prompt_rendering_id"`
	RunID             string `json:"run_id"`
	Model             string `json:"model"`
	Language          string `json:"language"`
	SampleID          string `json:"sample_id"`
	PromptVersionID   string `json:"prompt_version_id,omitempty"`
	PromptMode        string `json:"prompt_mode,omitempty"`
	CreatedAtUTC      string `json:"created_at_utc"`
}

type DBEvaluationRunItem struct {
	EvaluationRunID string `json:"evaluation_run_id"`
	RunID           string `json:"run_id"`
	SchemaVersion   string `json:"schema_version"`
	EvaluatedAtUTC  string `json:"evaluated_at_utc"`
	EnvID           string `json:"env_id,omitempty"`
	ScorePolicyID   string `json:"score_policy_id,omitempty"`
	CreatedDBAtUTC  string `json:"created_db_at_utc"`
}

type DBEvaluationStageItem struct {
	StageResultID      string `json:"stage_result_id"`
	EvaluationResultID string `json:"evaluation_result_id"`
	Stage              string `json:"stage"`
	Status             string `json:"status"`
	ExitCode           *int   `json:"exit_code,omitempty"`
	DurationMS         *int   `json:"duration_ms,omitempty"`
	CreatedAtUTC       string `json:"created_at_utc"`
}

type DBDatasetSampleItem struct {
	SampleUID    string `json:"sample_uid"`
	SampleID     string `json:"sample_id"`
	Language     string `json:"language"`
	Class        string `json:"class,omitempty"`
	Scenario     string `json:"scenario,omitempty"`
	Path         string `json:"path"`
	CreatedAtUTC string `json:"created_at_utc"`
}

type DBDatasetSnapshotItem struct {
	SnapshotID   string `json:"snapshot_id"`
	Fingerprint  string `json:"fingerprint"`
	SampleCount  int    `json:"sample_count"`
	CreatedAtUTC string `json:"created_at_utc"`
}

type DBModelConfigItem struct {
	ModelConfigID string `json:"model_config_id"`
	ModelName     string `json:"model_name"`
	Provider      string `json:"provider,omitempty"`
	ModelID       string `json:"model_id,omitempty"`
	CreatedAtUTC  string `json:"created_at_utc"`
}

type DBPromptProfileItem struct {
	ProfileID    string `json:"profile_id"`
	Strategy     string `json:"strategy,omitempty"`
	VersionID    string `json:"version_id,omitempty"`
	CreatedAtUTC string `json:"created_at_utc"`
}

type DBEvaluationEnvItem struct {
	EnvID        string `json:"env_id"`
	Fingerprint  string `json:"fingerprint"`
	CreatedAtUTC string `json:"created_at_utc"`
}

type DBScorePolicyItem struct {
	ScorePolicyID string `json:"score_policy_id"`
	Name          string `json:"name"`
	CreatedAtUTC  string `json:"created_at_utc"`
}

type DBReportItem struct {
	ReportID       string `json:"report_id"`
	RunID          string `json:"run_id"`
	GeneratedAtUTC string `json:"generated_at_utc"`
	CreatedDBAtUTC string `json:"created_db_at_utc"`
}

type DBRunArtifactItem struct {
	RunID        string `json:"run_id"`
	ArtifactID   string `json:"artifact_id"`
	Role         string `json:"role"`
	CreatedAtUTC string `json:"created_at_utc"`
}

type DBExperimentItem struct {
	ExperimentID string `json:"experiment_id"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	CreatedAtUTC string `json:"created_at_utc"`
	UpdatedAtUTC string `json:"updated_at_utc"`
}

type ReusableGeneratedCase struct {
	GeneratedCaseID       string `json:"generated_case_id"`
	RunID                 string `json:"run_id"`
	Model                 string `json:"model"`
	Language              string `json:"language"`
	SampleID              string `json:"sample_id"`
	GeneratedTestPath     string `json:"generated_test_path"`
	ResponsePath          string `json:"response_path,omitempty"`
	MetadataPath          string `json:"metadata_path,omitempty"`
	PromptVersionID       string `json:"prompt_version_id,omitempty"`
	PromptMode            string `json:"prompt_mode,omitempty"`
	LatencyMS             int    `json:"latency_ms,omitempty"`
	PromptTokens          *int   `json:"prompt_tokens,omitempty"`
	CompletionTokens      *int   `json:"completion_tokens,omitempty"`
	TotalTokens           *int   `json:"total_tokens,omitempty"`
	GeneratedAtUTC        string `json:"generated_at_utc,omitempty"`
	GeneratedTestArtifact string `json:"generated_test_artifact_id,omitempty"`
}

type IngestRunOptions struct {
	RunDir string
}

type IngestSummary struct {
	RunID                string `json:"run_id"`
	ManifestIngested     bool   `json:"manifest_ingested"`
	EvaluationIngested   bool   `json:"evaluation_ingested"`
	ReportIngested       bool   `json:"report_ingested"`
	GenerationCases      int    `json:"generation_cases"`
	EvaluationResults    int    `json:"evaluation_results"`
	ArtifactsIndexed     int    `json:"artifacts_indexed"`
	UnavailableArtifacts int    `json:"unavailable_artifacts"`
}

type ingestContext struct {
	runID   string
	runDir  string
	spec    contracts.RunSpec
	seen    map[string]struct{}
	summary *IngestSummary
}

func OpenSQLite(path string) (*SQLiteStore, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	return &SQLiteStore{db: db, path: path}, nil
}

func (s *SQLiteStore) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *SQLiteStore) Init(ctx context.Context) error {
	ddl := []string{
		`PRAGMA foreign_keys = ON;`,
		`CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS experiments (
			experiment_id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			created_at_utc TEXT NOT NULL,
			updated_at_utc TEXT NOT NULL,
			deleted_at_utc TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS artifacts (
			artifact_id TEXT PRIMARY KEY,
			kind TEXT NOT NULL,
			path TEXT NOT NULL,
			size_bytes INTEGER NOT NULL,
			sha256 TEXT NOT NULL,
			redacted INTEGER NOT NULL DEFAULT 0,
			created_at_utc TEXT NOT NULL,
			deleted_at_utc TEXT,
			UNIQUE(kind, sha256)
		);`,
		`CREATE TABLE IF NOT EXISTS run_artifacts (
			run_id TEXT NOT NULL,
			artifact_id TEXT NOT NULL,
			role TEXT NOT NULL,
			created_at_utc TEXT NOT NULL,
			PRIMARY KEY(run_id, artifact_id, role)
		);`,
		`CREATE TABLE IF NOT EXISTS dataset_samples (
			sample_uid TEXT PRIMARY KEY,
			sample_id TEXT NOT NULL,
			language TEXT NOT NULL,
			class TEXT,
			scenario TEXT,
			path TEXT NOT NULL,
			source_md5 TEXT,
			source_sha256 TEXT,
			source_artifact_id TEXT,
			risk_flags_json TEXT,
			created_at_utc TEXT NOT NULL,
			updated_at_utc TEXT NOT NULL,
			UNIQUE(language, sample_id, path)
		);`,
		`CREATE TABLE IF NOT EXISTS dataset_snapshots (
			snapshot_id TEXT PRIMARY KEY,
			fingerprint TEXT NOT NULL UNIQUE,
			created_at_utc TEXT NOT NULL,
			sample_count INTEGER NOT NULL,
			spec_json TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS dataset_snapshot_members (
			snapshot_id TEXT NOT NULL,
			sample_uid TEXT NOT NULL,
			PRIMARY KEY(snapshot_id, sample_uid)
		);`,
		`CREATE TABLE IF NOT EXISTS model_configs (
			model_config_id TEXT PRIMARY KEY,
			model_name TEXT NOT NULL,
			provider TEXT,
			model_id TEXT,
			config_json TEXT,
			created_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS prompt_profiles (
			profile_id TEXT PRIMARY KEY,
			strategy TEXT,
			version_id TEXT,
			template_hash TEXT,
			rules_json TEXT,
			created_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS prompt_renderings (
			prompt_rendering_id TEXT PRIMARY KEY,
			run_id TEXT NOT NULL,
			model TEXT NOT NULL,
			language TEXT NOT NULL,
			sample_id TEXT NOT NULL,
			prompt_artifact_id TEXT,
			prompt_version_id TEXT,
			prompt_mode TEXT,
			created_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS generation_runs (
			run_id TEXT PRIMARY KEY,
			experiment_id TEXT,
			schema_version TEXT NOT NULL,
			created_at_utc TEXT NOT NULL,
			spec_json TEXT NOT NULL,
			dataset_snapshot_id TEXT,
			dataset_fingerprint TEXT,
			prompt_strategy TEXT,
			prompt_version_id TEXT,
			prompt_snapshot_dir TEXT,
			manifest_artifact_id TEXT,
			created_db_at_utc TEXT NOT NULL,
			updated_db_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS generated_cases (
			generated_case_id TEXT PRIMARY KEY,
			run_id TEXT NOT NULL,
			model TEXT NOT NULL,
			language TEXT NOT NULL,
			sample_id TEXT NOT NULL,
			sample_uid TEXT,
			sample_path TEXT,
			prompt_rendering_id TEXT,
			generated_test_artifact_id TEXT,
			response_artifact_id TEXT,
			metadata_artifact_id TEXT,
			latency_ms INTEGER,
			prompt_tokens INTEGER,
			completion_tokens INTEGER,
			total_tokens INTEGER,
			generated_at_utc TEXT,
			success INTEGER NOT NULL,
			truncated INTEGER NOT NULL DEFAULT 0,
			error_json TEXT,
			created_db_at_utc TEXT NOT NULL,
			updated_db_at_utc TEXT NOT NULL,
			UNIQUE(run_id, model, language, sample_id)
		);`,
		`CREATE TABLE IF NOT EXISTS evaluation_envs (
			env_id TEXT PRIMARY KEY,
			fingerprint TEXT NOT NULL UNIQUE,
			env_json TEXT NOT NULL,
			created_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS score_policies (
			score_policy_id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			policy_json TEXT NOT NULL,
			created_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS evaluation_runs (
			evaluation_run_id TEXT PRIMARY KEY,
			run_id TEXT NOT NULL,
			experiment_id TEXT,
			generation_run_id TEXT,
			schema_version TEXT NOT NULL,
			evaluated_at_utc TEXT NOT NULL,
			manifest_path TEXT,
			result_artifact_id TEXT,
			env_id TEXT,
			score_policy_id TEXT,
			created_db_at_utc TEXT NOT NULL,
			updated_db_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS evaluation_results (
			evaluation_result_id TEXT PRIMARY KEY,
			evaluation_run_id TEXT NOT NULL,
			run_id TEXT NOT NULL,
			model TEXT NOT NULL,
			language TEXT NOT NULL,
			sample_id TEXT NOT NULL,
			generated_case_id TEXT,
			generated_test_artifact_id TEXT,
			source_artifact_id TEXT,
			compile_pass INTEGER NOT NULL,
			test_pass INTEGER,
			test_pass_count INTEGER,
			test_total_count INTEGER,
			test_pass_rate REAL,
			line_coverage REAL,
			branch_coverage REAL,
			mutation_score REAL,
			mutation_total INTEGER,
			mutation_killed INTEGER,
			mutation_survived INTEGER,
			mutation_no_tests INTEGER,
			mutation_timeouts INTEGER,
			mutation_skipped INTEGER,
			mutation_suspicious INTEGER,
			assertion_count INTEGER,
			test_case_count INTEGER,
			assertion_density REAL,
			runtime_ms INTEGER,
			prompt_tokens INTEGER,
			completion_tokens INTEGER,
			total_tokens INTEGER,
			truncated INTEGER NOT NULL DEFAULT 0,
			mutation_tool TEXT,
			failure_origin TEXT,
			score_eligible INTEGER NOT NULL DEFAULT 1,
			score_exclusion_reason TEXT,
			compile_error TEXT,
			test_error TEXT,
			coverage_error TEXT,
			mutation_error TEXT,
			created_db_at_utc TEXT NOT NULL,
			updated_db_at_utc TEXT NOT NULL,
			UNIQUE(evaluation_run_id, model, language, sample_id)
		);`,
		`CREATE TABLE IF NOT EXISTS evaluation_stage_results (
			stage_result_id TEXT PRIMARY KEY,
			evaluation_result_id TEXT NOT NULL,
			stage TEXT NOT NULL,
			command TEXT,
			workdir TEXT,
			exit_code INTEGER,
			duration_ms INTEGER,
			status TEXT NOT NULL,
			stdout_artifact_id TEXT,
			stderr_artifact_id TEXT,
			summary_json TEXT,
			created_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS report_snapshots (
			report_id TEXT PRIMARY KEY,
			run_id TEXT NOT NULL,
			evaluation_run_id TEXT,
			generated_at_utc TEXT NOT NULL,
			source_evaluation TEXT,
			report_json_artifact_id TEXT,
			report_html_artifact_id TEXT,
			summary_json TEXT,
			filters_json TEXT,
			created_db_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS report_result_members (
			report_id TEXT NOT NULL,
			evaluation_result_id TEXT NOT NULL,
			PRIMARY KEY(report_id, evaluation_result_id)
		);`,
		`CREATE TABLE IF NOT EXISTS log_events (
			log_event_id TEXT PRIMARY KEY,
			run_id TEXT NOT NULL,
			ts_utc TEXT,
			level TEXT,
			message TEXT NOT NULL,
			fields_json TEXT,
			created_at_utc TEXT NOT NULL
		);`,
	}
	for _, stmt := range ddl {
		if _, err := s.db.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}
	_, err := s.db.ExecContext(ctx, `INSERT OR IGNORE INTO schema_migrations(version, applied_at_utc) VALUES(?, ?)`, schemaVersion, nowUTC())
	return err
}

func (s *SQLiteStore) IngestManifestFile(ctx context.Context, path string) (IngestSummary, error) {
	manifest, err := contracts.ReadGeneratedManifest(path)
	if err != nil {
		return IngestSummary{}, err
	}
	sum := IngestSummary{RunID: manifest.RunID, ManifestIngested: true, GenerationCases: len(manifest.Cases)}
	ictx := newIngestContext(manifest.RunID, path, manifest.Spec, &sum)
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return IngestSummary{}, err
	}
	defer func() { _ = tx.Rollback() }()
	if err := s.ingestManifestTx(ctx, tx, path, manifest, ictx); err != nil {
		return IngestSummary{}, err
	}
	if err := tx.Commit(); err != nil {
		return IngestSummary{}, err
	}
	return sum, nil
}

func (s *SQLiteStore) IngestEvaluationFile(ctx context.Context, path string) (IngestSummary, error) {
	set, err := contracts.ReadEvaluationResultSet(path)
	if err != nil {
		return IngestSummary{}, err
	}
	sum := IngestSummary{RunID: set.RunID, EvaluationIngested: true, EvaluationResults: len(set.Results)}
	ictx := newIngestContext(set.RunID, path, contracts.RunSpec{RunID: set.RunID}, &sum)
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return IngestSummary{}, err
	}
	defer func() { _ = tx.Rollback() }()
	if set.ManifestPath != "" {
		manifestPath := ictx.resolvePath(set.ManifestPath)
		if manifest, readErr := contracts.ReadGeneratedManifest(manifestPath); readErr == nil {
			ictx.spec = manifest.Spec
			if err := s.ingestManifestTx(ctx, tx, manifestPath, manifest, ictx); err != nil {
				return IngestSummary{}, err
			}
			sum.ManifestIngested = true
			sum.GenerationCases = len(manifest.Cases)
		}
	}
	if err := s.ingestEvaluationTx(ctx, tx, path, set, ictx); err != nil {
		return IngestSummary{}, err
	}
	if err := tx.Commit(); err != nil {
		return IngestSummary{}, err
	}
	return sum, nil
}

func (s *SQLiteStore) IngestReportFile(ctx context.Context, path string) (IngestSummary, error) {
	var payload contracts.ReportPayload
	if err := readJSON(path, &payload); err != nil {
		return IngestSummary{}, err
	}
	sum := IngestSummary{RunID: payload.RunID, ReportIngested: true}
	ictx := newIngestContext(payload.RunID, path, contracts.RunSpec{RunID: payload.RunID}, &sum)
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return IngestSummary{}, err
	}
	defer func() { _ = tx.Rollback() }()
	if err := s.ingestReportTx(ctx, tx, path, payload, ictx); err != nil {
		return IngestSummary{}, err
	}
	if err := tx.Commit(); err != nil {
		return IngestSummary{}, err
	}
	return sum, nil
}

func (s *SQLiteStore) IngestRun(ctx context.Context, opts IngestRunOptions) (IngestSummary, error) {
	runDir := opts.RunDir
	if runDir == "" {
		return IngestSummary{}, fmt.Errorf("run dir is required")
	}
	manifestPath := filepath.Join(runDir, "generated", "generated_manifest.json")
	evaluationPath := filepath.Join(runDir, "evaluation", "evaluation_result.json")
	reportPath := filepath.Join(runDir, "report", "report_summary.json")
	var total IngestSummary
	if fileExists(manifestPath) {
		sum, err := s.IngestManifestFile(ctx, manifestPath)
		if err != nil {
			return IngestSummary{}, err
		}
		total = mergeIngestSummaries(total, sum)
	}
	if fileExists(evaluationPath) {
		sum, err := s.IngestEvaluationFile(ctx, evaluationPath)
		if err != nil {
			return IngestSummary{}, err
		}
		total = mergeIngestSummaries(total, sum)
	}
	if fileExists(reportPath) {
		sum, err := s.IngestReportFile(ctx, reportPath)
		if err != nil {
			return IngestSummary{}, err
		}
		total = mergeIngestSummaries(total, sum)
	}
	if total.RunID == "" {
		total.RunID = filepath.Base(runDir)
	}
	if !total.ManifestIngested && !total.EvaluationIngested && !total.ReportIngested {
		return IngestSummary{}, fmt.Errorf("no known artifacts found in %s", runDir)
	}
	return total, nil
}

// IngestEvaluation keeps the in-process orchestrator path simple. File-based
// ingest is preferred because it can also index linked manifest/artifacts.
func (s *SQLiteStore) IngestEvaluation(ctx context.Context, set contracts.EvaluationResultSet) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	sum := IngestSummary{RunID: set.RunID, EvaluationIngested: true, EvaluationResults: len(set.Results)}
	ictx := newIngestContext(set.RunID, "", contracts.RunSpec{RunID: set.RunID}, &sum)
	if err := s.ingestEvaluationTx(ctx, tx, "", set, ictx); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *SQLiteStore) Overview(ctx context.Context, limit int) (DBOverview, error) {
	if limit <= 0 {
		limit = 10
	}
	out := DBOverview{DBPath: s.path, SchemaVersion: schemaVersion}
	counts := []struct {
		dst   *int
		table string
	}{
		{&out.GenerationRuns, "generation_runs"},
		{&out.EvaluationRuns, "evaluation_runs"},
		{&out.GeneratedCases, "generated_cases"},
		{&out.EvaluationResults, "evaluation_results"},
		{&out.Artifacts, "artifacts"},
		{&out.Reports, "report_snapshots"},
	}
	for _, c := range counts {
		if err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+c.table).Scan(c.dst); err != nil {
			return DBOverview{}, err
		}
	}
	runs, err := s.ListRuns(ctx, limit)
	if err != nil {
		return DBOverview{}, err
	}
	out.LatestRuns = runs
	return out, nil
}

func (s *SQLiteStore) ListRuns(ctx context.Context, limit int) ([]DBRunItem, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			gr.run_id,
			COALESCE(gr.experiment_id, ''),
			COALESCE(e.name, ''),
			COALESCE(gr.created_at_utc, ''),
			COALESCE(MAX(erun.evaluated_at_utc), ''),
			COUNT(DISTINCT gc.model),
			COUNT(DISTINCT gc.language),
			COUNT(DISTINCT gc.generated_case_id),
			COUNT(DISTINCT er.evaluation_result_id),
			COUNT(DISTINCT rs.report_id)
		FROM generation_runs gr
		LEFT JOIN experiments e ON e.experiment_id = gr.experiment_id
		LEFT JOIN generated_cases gc ON gc.run_id = gr.run_id
		LEFT JOIN evaluation_runs erun ON erun.run_id = gr.run_id
		LEFT JOIN evaluation_results er ON er.evaluation_run_id = erun.evaluation_run_id
		LEFT JOIN report_snapshots rs ON rs.run_id = gr.run_id
		GROUP BY gr.run_id
		ORDER BY COALESCE(MAX(erun.evaluated_at_utc), gr.created_at_utc) DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DBRunItem
	for rows.Next() {
		var r DBRunItem
		if err := rows.Scan(&r.RunID, &r.ExperimentID, &r.ExperimentName, &r.CreatedAtUTC, &r.EvaluatedAtUTC, &r.Models, &r.Languages, &r.GeneratedCases, &r.EvaluationResults, &r.Reports); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ListResults(ctx context.Context, runID, model, language string, limit int) ([]DBResultItem, error) {
	if limit <= 0 || limit > 1000 {
		limit = 200
	}
	where := []string{"1=1"}
	args := []any{}
	if runID != "" {
		where = append(where, "run_id = ?")
		args = append(args, runID)
	}
	if model != "" {
		where = append(where, "model = ?")
		args = append(args, model)
	}
	if language != "" {
		where = append(where, "language = ?")
		args = append(args, language)
	}
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx, `
		SELECT evaluation_result_id, evaluation_run_id, run_id, model, language, sample_id,
		       compile_pass, test_pass, line_coverage, mutation_score, mutation_total,
		       COALESCE(failure_origin, ''), score_eligible, COALESCE(score_exclusion_reason, ''), runtime_ms
		FROM evaluation_results
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY run_id DESC, model, language, sample_id
		LIMIT ?`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DBResultItem
	for rows.Next() {
		var item DBResultItem
		var testPass sql.NullInt64
		var lineCov, mutScore sql.NullFloat64
		var mutTotal, runtime sql.NullInt64
		var compilePass, eligible int
		if err := rows.Scan(
			&item.EvaluationResultID, &item.EvaluationRunID, &item.RunID, &item.Model, &item.Language, &item.SampleID,
			&compilePass, &testPass, &lineCov, &mutScore, &mutTotal,
			&item.FailureOrigin, &eligible, &item.ScoreExclusionReason, &runtime,
		); err != nil {
			return nil, err
		}
		item.CompilePass = compilePass != 0
		item.ScoreEligible = eligible != 0
		if testPass.Valid {
			v := testPass.Int64 != 0
			item.TestPass = &v
		}
		if lineCov.Valid {
			v := lineCov.Float64
			item.LineCoverage = &v
		}
		if mutScore.Valid {
			v := mutScore.Float64
			item.MutationScore = &v
		}
		if mutTotal.Valid {
			v := int(mutTotal.Int64)
			item.MutationTotal = &v
		}
		if runtime.Valid {
			v := int(runtime.Int64)
			item.RuntimeMS = &v
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ListArtifacts(ctx context.Context, runID, kind string, limit int) ([]DBArtifactItem, error) {
	if limit <= 0 || limit > 1000 {
		limit = 200
	}
	args := []any{}
	query := `SELECT DISTINCT a.artifact_id, a.kind, a.path, a.size_bytes, a.sha256, a.redacted, a.created_at_utc, COALESCE(a.deleted_at_utc, '')
		FROM artifacts a`
	where := []string{"1=1"}
	if runID != "" {
		query += ` JOIN run_artifacts ra ON ra.artifact_id = a.artifact_id`
		where = append(where, "ra.run_id = ?")
		args = append(args, runID)
	}
	if kind != "" {
		where = append(where, "a.kind = ?")
		args = append(args, kind)
	}
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx, query+` WHERE `+strings.Join(where, " AND ")+` ORDER BY a.created_at_utc DESC LIMIT ?`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DBArtifactItem
	for rows.Next() {
		var item DBArtifactItem
		var redacted int
		if err := rows.Scan(&item.ArtifactID, &item.Kind, &item.Path, &item.SizeBytes, &item.SHA256, &redacted, &item.CreatedAtUTC, &item.DeletedAtUTC); err != nil {
			return nil, err
		}
		item.Redacted = redacted != 0
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ReportFacets(ctx context.Context, limit int) (DBReportFacets, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	var facets DBReportFacets
	runRows, err := s.db.QueryContext(ctx, `
		SELECT ev.run_id, ev.evaluation_run_id, COALESCE(ev.evaluated_at_utc, ''), COALESCE(ev.env_id, ''),
		       COALESCE(gr.prompt_strategy, ''), COALESCE(gr.prompt_version_id, ''),
		       COUNT(er.evaluation_result_id) AS result_count
		FROM evaluation_runs ev
		JOIN evaluation_results er ON er.evaluation_run_id = ev.evaluation_run_id
		LEFT JOIN generation_runs gr ON gr.run_id = ev.generation_run_id
		GROUP BY ev.run_id, ev.evaluation_run_id, ev.evaluated_at_utc, ev.env_id, gr.prompt_strategy, gr.prompt_version_id
		ORDER BY ev.evaluated_at_utc DESC
		LIMIT ?`, limit)
	if err != nil {
		return DBReportFacets{}, err
	}
	defer runRows.Close()
	for runRows.Next() {
		var row DBReportRunOption
		if err := runRows.Scan(&row.RunID, &row.EvaluationRunID, &row.EvaluatedAtUTC, &row.EnvID, &row.PromptStrategy, &row.PromptVersionID, &row.EvaluationResults); err != nil {
			return DBReportFacets{}, err
		}
		facets.Runs = append(facets.Runs, row)
	}
	if err := runRows.Err(); err != nil {
		return DBReportFacets{}, err
	}
	models, err := s.distinctStrings(ctx, `SELECT DISTINCT model FROM evaluation_results ORDER BY model`)
	if err != nil {
		return DBReportFacets{}, err
	}
	langs, err := s.distinctStrings(ctx, `SELECT DISTINCT language FROM evaluation_results ORDER BY language`)
	if err != nil {
		return DBReportFacets{}, err
	}
	facets.Models = models
	facets.Languages = langs
	envRows, err := s.db.QueryContext(ctx, `
		SELECT COALESCE(ev.env_id, ''), COALESCE(env.fingerprint, ''),
		       COUNT(DISTINCT ev.evaluation_run_id), COUNT(er.evaluation_result_id)
		FROM evaluation_runs ev
		JOIN evaluation_results er ON er.evaluation_run_id = ev.evaluation_run_id
		LEFT JOIN evaluation_envs env ON env.env_id = ev.env_id
		GROUP BY ev.env_id, env.fingerprint
		ORDER BY COUNT(er.evaluation_result_id) DESC`)
	if err != nil {
		return DBReportFacets{}, err
	}
	defer envRows.Close()
	for envRows.Next() {
		var row DBReportEnvOption
		if err := envRows.Scan(&row.EnvID, &row.Fingerprint, &row.EvaluationRuns, &row.EvaluationResults); err != nil {
			return DBReportFacets{}, err
		}
		facets.EnvGroups = append(facets.EnvGroups, row)
	}
	return facets, envRows.Err()
}

func (s *SQLiteStore) distinctStrings(ctx context.Context, query string) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		if strings.TrimSpace(v) != "" {
			out = append(out, v)
		}
	}
	return out, rows.Err()
}

func (s *SQLiteStore) SelectEvaluationResultSet(ctx context.Context, runID string, filter DBReportFilter) (contracts.EvaluationResultSet, error) {
	if runID == "" {
		runID = contracts.NewRunID()
	}
	where := []string{"1=1"}
	args := []any{}
	addInFilter := func(column string, values []string) {
		cleaned := nonEmptyStrings(values)
		if len(cleaned) == 0 {
			return
		}
		placeholders := make([]string, len(cleaned))
		for i, v := range cleaned {
			placeholders[i] = "?"
			args = append(args, v)
		}
		where = append(where, column+" IN ("+strings.Join(placeholders, ",")+")")
	}
	addInFilter("er.run_id", filter.RunIDs)
	addInFilter("er.evaluation_run_id", filter.EvaluationRunIDs)
	addInFilter("er.model", filter.Models)
	addInFilter("er.language", filter.Languages)
	if filter.ScoreEligibleOnly {
		where = append(where, "er.score_eligible = 1")
	}
	query := `
		SELECT
			er.model, er.language, er.sample_id,
			COALESCE(g.path, ''), COALESCE(src.path, ''),
			er.compile_pass, er.test_pass, er.truncated,
			er.line_coverage, er.branch_coverage, er.mutation_score,
			er.mutation_total, er.mutation_killed, er.mutation_survived,
			er.mutation_no_tests, er.mutation_timeouts, er.mutation_skipped, er.mutation_suspicious,
			er.assertion_count, er.test_case_count, er.assertion_density,
			er.test_pass_count, er.test_total_count, er.test_pass_rate,
			er.runtime_ms, er.prompt_tokens, er.completion_tokens, er.total_tokens,
			COALESCE(er.compile_error, ''), COALESCE(er.test_error, ''), COALESCE(er.coverage_error, ''), COALESCE(er.mutation_error, ''),
			COALESCE(er.mutation_tool, ''), COALESCE(er.failure_origin, ''), er.score_eligible, COALESCE(er.score_exclusion_reason, ''),
				gc.latency_ms
		FROM evaluation_results er
		LEFT JOIN artifacts g ON g.artifact_id = er.generated_test_artifact_id
		LEFT JOIN artifacts src ON src.artifact_id = er.source_artifact_id
		LEFT JOIN generated_cases gc ON gc.generated_case_id = er.generated_case_id
		WHERE ` + strings.Join(where, " AND ") + `
		ORDER BY er.run_id, er.model, er.language, er.sample_id`
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return contracts.EvaluationResultSet{}, err
	}
	defer rows.Close()
	var results []contracts.EvaluationResult
	for rows.Next() {
		var r contracts.EvaluationResult
		var compilePass, truncated, eligible int
		var testPass sql.NullInt64
		var lineCov, branchCov, mutationScore, assertionDensity, testPassRate sql.NullFloat64
		var mutationTotal, mutationKilled, mutationSurvived, mutationNoTests, mutationTimeouts, mutationSkipped, mutationSuspicious sql.NullInt64
		var assertionCount, testCaseCount, testPassCount, testTotalCount, runtimeMS, latencyMS, promptTokens, completionTokens, totalTokens sql.NullInt64
		if err := rows.Scan(
			&r.Model, &r.Language, &r.SampleID,
			&r.GeneratedTestPath, &r.SourcePath,
			&compilePass, &testPass, &truncated,
			&lineCov, &branchCov, &mutationScore,
			&mutationTotal, &mutationKilled, &mutationSurvived,
			&mutationNoTests, &mutationTimeouts, &mutationSkipped, &mutationSuspicious,
			&assertionCount, &testCaseCount, &assertionDensity,
			&testPassCount, &testTotalCount, &testPassRate,
			&runtimeMS, &promptTokens, &completionTokens, &totalTokens,
			&r.CompileError, &r.TestError, &r.CoverageError, &r.MutationError,
			&r.MutationTool, &r.FailureOrigin, &eligible, &r.ScoreExclusionReason,
				&latencyMS,
		); err != nil {
			return contracts.EvaluationResultSet{}, err
		}
		r.CompilePass = compilePass != 0
		r.Truncated = truncated != 0
		r.TestPass = nullableIntBool(testPass)
		r.LineCoverage = nullableSQLFloat(lineCov)
		r.BranchCoverage = nullableSQLFloat(branchCov)
		r.MutationScore = nullableSQLFloat(mutationScore)
		r.MutationTotal = nullableSQLInt(mutationTotal)
		r.MutationKilled = nullableSQLInt(mutationKilled)
		r.MutationSurvived = nullableSQLInt(mutationSurvived)
		r.MutationNoTests = nullableSQLInt(mutationNoTests)
		r.MutationTimeouts = nullableSQLInt(mutationTimeouts)
		r.MutationSkipped = nullableSQLInt(mutationSkipped)
		r.MutationSuspicious = nullableSQLInt(mutationSuspicious)
		r.AssertionCount = nullableSQLInt(assertionCount)
		r.TestCaseCount = nullableSQLInt(testCaseCount)
		r.AssertionDensity = nullableSQLFloat(assertionDensity)
		r.TestPassCount = nullableSQLInt(testPassCount)
		r.TestTotalCount = nullableSQLInt(testTotalCount)
		r.TestPassRate = nullableSQLFloat(testPassRate)
		r.RuntimeMS = nullableSQLInt(runtimeMS)
			r.LatencyMS = nullableSQLInt(latencyMS)
		r.PromptTokens = nullableSQLInt(promptTokens)
		r.CompletionTokens = nullableSQLInt(completionTokens)
		r.TotalTokens = nullableSQLInt(totalTokens)
		scoreEligible := eligible != 0
		r.ScoreEligible = &scoreEligible
		results = append(results, r)
	}
	if err := rows.Err(); err != nil {
		return contracts.EvaluationResultSet{}, err
	}
	if len(results) == 0 {
		return contracts.EvaluationResultSet{}, fmt.Errorf("no evaluation results match database report filter")
	}
	evaluatedAt := time.Now().UTC()
	var manifestPath string
	_ = s.db.QueryRowContext(ctx, `
		SELECT COALESCE(ev.manifest_path, '')
		FROM evaluation_runs ev
		WHERE ev.evaluation_run_id IN (
			SELECT DISTINCT er.evaluation_run_id
			FROM evaluation_results er
			WHERE `+strings.Join(where, " AND ")+`
		)
		ORDER BY evaluated_at_utc DESC
		LIMIT 1`, args...).Scan(&manifestPath)
	return contracts.EvaluationResultSet{
		SchemaVersion:  contracts.SchemaVersion,
		RunID:          runID,
		EvaluatedAtUTC: evaluatedAt,
		ManifestPath:   manifestPath,
		Results:        results,
	}, nil
}

func (s *SQLiteStore) FindReusableGeneratedCase(ctx context.Context, model, language, sampleID, sourceSHA256, promptVersionID string) (ReusableGeneratedCase, bool, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT gc.generated_case_id, gc.run_id, gc.model, gc.language, gc.sample_id,
		       COALESCE(test_art.path, ''), COALESCE(resp_art.path, ''), COALESCE(meta_art.path, ''),
		       COALESCE(pr.prompt_version_id, ''), COALESCE(pr.prompt_mode, ''),
		       COALESCE(gc.latency_ms, 0), gc.prompt_tokens, gc.completion_tokens, gc.total_tokens,
		       COALESCE(gc.generated_at_utc, ''), COALESCE(gc.generated_test_artifact_id, '')
		FROM generated_cases gc
		JOIN dataset_samples ds ON ds.sample_uid = gc.sample_uid
		LEFT JOIN prompt_renderings pr ON pr.prompt_rendering_id = gc.prompt_rendering_id
		LEFT JOIN artifacts test_art ON test_art.artifact_id = gc.generated_test_artifact_id
		LEFT JOIN artifacts resp_art ON resp_art.artifact_id = gc.response_artifact_id
		LEFT JOIN artifacts meta_art ON meta_art.artifact_id = gc.metadata_artifact_id
		WHERE gc.success = 1
		  AND gc.model = ?
		  AND gc.language = ?
		  AND gc.sample_id = ?
		  AND COALESCE(ds.source_sha256, '') = ?
		  AND COALESCE(pr.prompt_version_id, '') = ?
		  AND COALESCE(test_art.deleted_at_utc, '') = ''
		ORDER BY gc.generated_at_utc DESC
		LIMIT 1`, model, language, sampleID, sourceSHA256, promptVersionID)
	var out ReusableGeneratedCase
	var promptTokens, completionTokens, totalTokens sql.NullInt64
	if err := row.Scan(
		&out.GeneratedCaseID, &out.RunID, &out.Model, &out.Language, &out.SampleID,
		&out.GeneratedTestPath, &out.ResponsePath, &out.MetadataPath,
		&out.PromptVersionID, &out.PromptMode,
		&out.LatencyMS, &promptTokens, &completionTokens, &totalTokens,
		&out.GeneratedAtUTC, &out.GeneratedTestArtifact,
	); err != nil {
		if err == sql.ErrNoRows {
			return ReusableGeneratedCase{}, false, nil
		}
		return ReusableGeneratedCase{}, false, err
	}
	out.GeneratedTestPath = resolveStoredPath(out.GeneratedTestPath)
	out.ResponsePath = resolveStoredPath(out.ResponsePath)
	out.MetadataPath = resolveStoredPath(out.MetadataPath)
	out.PromptTokens = nullableSQLInt(promptTokens)
	out.CompletionTokens = nullableSQLInt(completionTokens)
	out.TotalTokens = nullableSQLInt(totalTokens)
	if !fileExists(out.GeneratedTestPath) {
		return ReusableGeneratedCase{}, false, nil
	}
	return out, true, nil
}

func (s *SQLiteStore) ingestManifestTx(ctx context.Context, tx *sql.Tx, path string, manifest contracts.GeneratedManifest, ictx *ingestContext) error {
	now := nowUTC()
	manifestArtifactID, err := s.putArtifact(ctx, tx, ictx, "generated_manifest", path, "manifest")
	if err != nil {
		return err
	}
	specJSON := mustJSON(manifest.Spec)
	sampleUIDs, fingerprint, snapshotID, err := s.upsertDatasetSnapshot(ctx, tx, manifest, ictx)
	if err != nil {
		return err
	}
	profileID := stableID("prompt_profile", manifest.PromptStrategy, manifest.PromptVersionID)
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO prompt_profiles(profile_id, strategy, version_id, template_hash, rules_json, created_at_utc)
		VALUES(?, ?, ?, ?, ?, ?)
		ON CONFLICT(profile_id) DO UPDATE SET strategy=excluded.strategy, version_id=excluded.version_id`,
		profileID, manifest.PromptStrategy, manifest.PromptVersionID, "", "{}", now); err != nil {
		return fmt.Errorf("upsert prompt profile: %w", err)
	}
	// 入库模型配置：从models.yaml读取完整配置
	modelConfigs := s.readModelConfigsFromYAML(manifest.Spec.ConfigPath)
	seenModels := map[string]struct{}{}
	for _, c := range manifest.Cases {
		seenModels[c.Model] = struct{}{}
	}
	for model := range seenModels {
		cfg := modelConfigs[model]
		modelConfigID := stableID("model_config", model, cfg.Provider, cfg.ModelID)
		configJSON := mustJSON(cfg)
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO model_configs(model_config_id, model_name, provider, model_id, config_json, created_at_utc)
			VALUES(?, ?, ?, ?, ?, ?)
			ON CONFLICT(model_config_id) DO UPDATE SET provider=excluded.provider, model_id=excluded.model_id, config_json=excluded.config_json`,
			modelConfigID, model, cfg.Provider, cfg.ModelID, configJSON, now); err != nil {
			return fmt.Errorf("upsert model config: %w", err)
		}
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO generation_runs(run_id, experiment_id, schema_version, created_at_utc, spec_json,
			dataset_snapshot_id, dataset_fingerprint, prompt_strategy, prompt_version_id, prompt_snapshot_dir,
			manifest_artifact_id, created_db_at_utc, updated_db_at_utc)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(run_id) DO UPDATE SET
			schema_version=excluded.schema_version,
			created_at_utc=excluded.created_at_utc,
			spec_json=excluded.spec_json,
			dataset_snapshot_id=excluded.dataset_snapshot_id,
			dataset_fingerprint=excluded.dataset_fingerprint,
			prompt_strategy=excluded.prompt_strategy,
			prompt_version_id=excluded.prompt_version_id,
			prompt_snapshot_dir=excluded.prompt_snapshot_dir,
			manifest_artifact_id=excluded.manifest_artifact_id,
			updated_db_at_utc=excluded.updated_db_at_utc`,
		manifest.RunID, defaultExperimentID(manifest.RunID), manifest.SchemaVersion, formatTime(manifest.CreatedAtUTC), specJSON,
		snapshotID, fingerprint, manifest.PromptStrategy, manifest.PromptVersionID, manifest.PromptSnapshotDir,
		manifestArtifactID, now, now); err != nil {
		return fmt.Errorf("upsert generation run: %w", err)
	}
	if err := s.ensureExperiment(ctx, tx, defaultExperimentID(manifest.RunID), manifest.RunID); err != nil {
		return err
	}
	for _, c := range manifest.Cases {
		sampleUID := sampleUIDs[caseKey(c.Model, c.Language, c.SampleID)]
		promptArtifactID, _ := s.putOptionalArtifact(ctx, tx, ictx, "prompt_rendering", c.PromptPath, "prompt")
		promptRenderingID := stableID("prompt_rendering", manifest.RunID, c.Model, c.Language, c.SampleID)
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO prompt_renderings(prompt_rendering_id, run_id, model, language, sample_id,
				prompt_artifact_id, prompt_version_id, prompt_mode, created_at_utc)
			VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(prompt_rendering_id) DO UPDATE SET
				prompt_artifact_id=excluded.prompt_artifact_id,
				prompt_version_id=excluded.prompt_version_id,
				prompt_mode=excluded.prompt_mode`,
			promptRenderingID, manifest.RunID, c.Model, c.Language, c.SampleID, nullString(promptArtifactID), c.PromptVersionID, c.PromptMode, now); err != nil {
			return fmt.Errorf("upsert prompt rendering: %w", err)
		}
		generatedTestArtifactID, _ := s.putOptionalArtifact(ctx, tx, ictx, "generated_test", c.GeneratedTestPath, "generated_test")
		responseArtifactID, _ := s.putOptionalArtifact(ctx, tx, ictx, "model_response", c.ResponsePath, "response")
		metadataArtifactID, _ := s.putOptionalArtifact(ctx, tx, ictx, "generation_metadata", c.MetadataPath, "metadata")
		caseID := stableID("generated_case", manifest.RunID, c.Model, c.Language, c.SampleID)
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO generated_cases(generated_case_id, run_id, model, language, sample_id, sample_uid, sample_path,
				prompt_rendering_id, generated_test_artifact_id, response_artifact_id, metadata_artifact_id,
				latency_ms, prompt_tokens, completion_tokens, total_tokens, generated_at_utc, success, truncated,
				error_json, created_db_at_utc, updated_db_at_utc)
			VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(run_id, model, language, sample_id) DO UPDATE SET
				sample_uid=excluded.sample_uid,
				sample_path=excluded.sample_path,
				prompt_rendering_id=excluded.prompt_rendering_id,
				generated_test_artifact_id=excluded.generated_test_artifact_id,
				response_artifact_id=excluded.response_artifact_id,
				metadata_artifact_id=excluded.metadata_artifact_id,
				latency_ms=excluded.latency_ms,
				prompt_tokens=excluded.prompt_tokens,
				completion_tokens=excluded.completion_tokens,
				total_tokens=excluded.total_tokens,
				generated_at_utc=excluded.generated_at_utc,
				success=excluded.success,
				truncated=excluded.truncated,
				error_json=excluded.error_json,
				updated_db_at_utc=excluded.updated_db_at_utc`,
			caseID, manifest.RunID, c.Model, c.Language, c.SampleID, nullString(sampleUID), c.SamplePath,
			promptRenderingID, nullString(generatedTestArtifactID), nullString(responseArtifactID), nullString(metadataArtifactID),
			c.LatencyMS, nullableInt(c.PromptTokens), nullableInt(c.CompletionTokens), nullableInt(c.TotalTokens), formatTime(c.GeneratedAtUTC),
			boolToInt(c.Success), boolToInt(c.Truncated), nullString(string(mustJSON(c.Error))), now, now); err != nil {
			return fmt.Errorf("upsert generated case: %w", err)
		}
	}
	return nil
}

func (s *SQLiteStore) ingestEvaluationTx(ctx context.Context, tx *sql.Tx, path string, set contracts.EvaluationResultSet, ictx *ingestContext) error {
	now := nowUTC()
	resultArtifactID := ""
	var err error
	if path != "" {
		resultArtifactID, err = s.putArtifact(ctx, tx, ictx, "evaluation_result", path, "evaluation")
		if err != nil {
			return err
		}
	}
	// 使用真实的环境指纹（如果存在）
	envFingerprint := set.EnvironmentFingerprint
	envJSON := set.EnvironmentJSON
	if envFingerprint == "" {
		envFingerprint = "unknown"
		envJSON = `{"fingerprint":"unknown","note":"environment capture not available"}`
	}
	envID := stableID("evaluation_env", envFingerprint)
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO evaluation_envs(env_id, fingerprint, env_json, created_at_utc)
		VALUES(?, ?, ?, ?)
		ON CONFLICT(fingerprint) DO NOTHING`,
		envID, envFingerprint, envJSON, now); err != nil {
		return fmt.Errorf("upsert evaluation env: %w", err)
	}
	policyID := stableID("score_policy", "default-v2")
	policyJSON := `{"compile":0.30,"test":0.30,"coverage":0.20,"mutation":0.20,"excludes_score_eligible_false":true}`
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO score_policies(score_policy_id, name, policy_json, created_at_utc)
		VALUES(?, ?, ?, ?)
		ON CONFLICT(score_policy_id) DO UPDATE SET policy_json=excluded.policy_json`,
		policyID, "default-v2", policyJSON, now); err != nil {
		return fmt.Errorf("upsert score policy: %w", err)
	}
	evalRunID := stableID("evaluation_run", set.RunID, formatTime(set.EvaluatedAtUTC), path)
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO evaluation_runs(evaluation_run_id, run_id, experiment_id, generation_run_id, schema_version,
			evaluated_at_utc, manifest_path, result_artifact_id, env_id, score_policy_id, created_db_at_utc, updated_db_at_utc)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(evaluation_run_id) DO UPDATE SET
			schema_version=excluded.schema_version,
			evaluated_at_utc=excluded.evaluated_at_utc,
			manifest_path=excluded.manifest_path,
			result_artifact_id=excluded.result_artifact_id,
			env_id=excluded.env_id,
			score_policy_id=excluded.score_policy_id,
			updated_db_at_utc=excluded.updated_db_at_utc`,
		evalRunID, set.RunID, defaultExperimentID(set.RunID), set.RunID, set.SchemaVersion,
		formatTime(set.EvaluatedAtUTC), set.ManifestPath, nullString(resultArtifactID), envID, policyID, now, now); err != nil {
		return fmt.Errorf("upsert evaluation run: %w", err)
	}
	if err := s.ensureExperiment(ctx, tx, defaultExperimentID(set.RunID), set.RunID); err != nil {
		return err
	}
	for _, row := range set.Results {
		sourceArtifactID, _ := s.putOptionalArtifact(ctx, tx, ictx, "dataset_source", row.SourcePath, "source")
		generatedTestArtifactID, _ := s.putOptionalArtifact(ctx, tx, ictx, "generated_test", row.GeneratedTestPath, "generated_test")
		generatedCaseID := stableID("generated_case", set.RunID, row.Model, row.Language, row.SampleID)
		resultID := stableID("evaluation_result", evalRunID, row.Model, row.Language, row.SampleID)
		eligible := true
		if row.ScoreEligible != nil {
			eligible = *row.ScoreEligible
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO evaluation_results(evaluation_result_id, evaluation_run_id, run_id, model, language, sample_id,
				generated_case_id, generated_test_artifact_id, source_artifact_id,
				compile_pass, test_pass, test_pass_count, test_total_count, test_pass_rate,
				line_coverage, branch_coverage, mutation_score, mutation_total, mutation_killed, mutation_survived,
				mutation_no_tests, mutation_timeouts, mutation_skipped, mutation_suspicious,
				assertion_count, test_case_count, assertion_density, runtime_ms,
				prompt_tokens, completion_tokens, total_tokens, truncated, mutation_tool,
				failure_origin, score_eligible, score_exclusion_reason,
				compile_error, test_error, coverage_error, mutation_error,
				created_db_at_utc, updated_db_at_utc)
			VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(evaluation_run_id, model, language, sample_id) DO UPDATE SET
				generated_case_id=excluded.generated_case_id,
				generated_test_artifact_id=excluded.generated_test_artifact_id,
				source_artifact_id=excluded.source_artifact_id,
				compile_pass=excluded.compile_pass,
				test_pass=excluded.test_pass,
				test_pass_count=excluded.test_pass_count,
				test_total_count=excluded.test_total_count,
				test_pass_rate=excluded.test_pass_rate,
				line_coverage=excluded.line_coverage,
				branch_coverage=excluded.branch_coverage,
				mutation_score=excluded.mutation_score,
				mutation_total=excluded.mutation_total,
				mutation_killed=excluded.mutation_killed,
				mutation_survived=excluded.mutation_survived,
				mutation_no_tests=excluded.mutation_no_tests,
				mutation_timeouts=excluded.mutation_timeouts,
				mutation_skipped=excluded.mutation_skipped,
				mutation_suspicious=excluded.mutation_suspicious,
				assertion_count=excluded.assertion_count,
				test_case_count=excluded.test_case_count,
				assertion_density=excluded.assertion_density,
				runtime_ms=excluded.runtime_ms,
				prompt_tokens=excluded.prompt_tokens,
				completion_tokens=excluded.completion_tokens,
				total_tokens=excluded.total_tokens,
				truncated=excluded.truncated,
				mutation_tool=excluded.mutation_tool,
				failure_origin=excluded.failure_origin,
				score_eligible=excluded.score_eligible,
				score_exclusion_reason=excluded.score_exclusion_reason,
				compile_error=excluded.compile_error,
				test_error=excluded.test_error,
				coverage_error=excluded.coverage_error,
				mutation_error=excluded.mutation_error,
				updated_db_at_utc=excluded.updated_db_at_utc`,
			resultID, evalRunID, set.RunID, row.Model, row.Language, row.SampleID,
			generatedCaseID, nullString(generatedTestArtifactID), nullString(sourceArtifactID),
			boolToInt(row.CompilePass), ptrBoolToNullableInt(row.TestPass), nullableInt(row.TestPassCount), nullableInt(row.TestTotalCount), nullableFloat(row.TestPassRate),
			nullableFloat(row.LineCoverage), nullableFloat(row.BranchCoverage), nullableFloat(row.MutationScore), nullableInt(row.MutationTotal), nullableInt(row.MutationKilled), nullableInt(row.MutationSurvived),
			nullableInt(row.MutationNoTests), nullableInt(row.MutationTimeouts), nullableInt(row.MutationSkipped), nullableInt(row.MutationSuspicious),
			nullableInt(row.AssertionCount), nullableInt(row.TestCaseCount), nullableFloat(row.AssertionDensity), nullableInt(row.RuntimeMS),
			nullableInt(row.PromptTokens), nullableInt(row.CompletionTokens), nullableInt(row.TotalTokens), boolToInt(row.Truncated), nullString(row.MutationTool),
			nullString(row.FailureOrigin), boolToInt(eligible), nullString(row.ScoreExclusionReason),
			nullString(row.CompileError), nullString(row.TestError), nullString(row.CoverageError), nullString(row.MutationError),
			now, now); err != nil {
			return fmt.Errorf("upsert evaluation result: %w", err)
		}
	}
	return nil
}

func (s *SQLiteStore) ingestReportTx(ctx context.Context, tx *sql.Tx, path string, payload contracts.ReportPayload, ictx *ingestContext) error {
	now := nowUTC()
	reportJSONArtifactID, err := s.putArtifact(ctx, tx, ictx, "report_summary", path, "report_json")
	if err != nil {
		return err
	}
	htmlPath := filepath.Join(filepath.Dir(path), "report.html")
	reportHTMLArtifactID, _ := s.putOptionalArtifact(ctx, tx, ictx, "report_html", htmlPath, "report_html")
	evalRunID := ""
	if payload.SourceEvaluation != "" {
		evalRunID = stableID("evaluation_run", payload.RunID, "", ictx.resolvePath(payload.SourceEvaluation))
	}
	reportID := stableID("report", payload.RunID, formatTime(payload.GeneratedAtUTC), path)
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO report_snapshots(report_id, run_id, evaluation_run_id, generated_at_utc, source_evaluation,
			report_json_artifact_id, report_html_artifact_id, summary_json, filters_json, created_db_at_utc)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(report_id) DO UPDATE SET
			report_json_artifact_id=excluded.report_json_artifact_id,
			report_html_artifact_id=excluded.report_html_artifact_id,
			summary_json=excluded.summary_json,
			filters_json=excluded.filters_json`,
		reportID, payload.RunID, nullString(evalRunID), formatTime(payload.GeneratedAtUTC), payload.SourceEvaluation,
		reportJSONArtifactID, nullString(reportHTMLArtifactID), string(mustJSON(payload.Summary)), `{}`, now); err != nil {
		return fmt.Errorf("upsert report: %w", err)
	}
	return nil
}

func (s *SQLiteStore) upsertDatasetSnapshot(ctx context.Context, tx *sql.Tx, manifest contracts.GeneratedManifest, ictx *ingestContext) (map[string]string, string, string, error) {
	type sampleInfo struct {
		key      string
		uid      string
		id       string
		language string
		class    string
		scenario string
		path     string
		sha      string
		md5      string
		artifact string
	}
	byKey := map[string]sampleInfo{}
	for _, c := range manifest.Cases {
		key := caseKey(c.Model, c.Language, c.SampleID)
		if _, ok := byKey[key]; ok {
			continue
		}
		class, scenario := inferClassScenario(c.SamplePath, c.SampleID)
		sourceArtifactID, sha, _ := s.putSourceArtifact(ctx, tx, ictx, c.SamplePath)
		uid := stableID("dataset_sample", c.Language, c.SampleID, c.SamplePath, sha)
		byKey[key] = sampleInfo{
			key: key, uid: uid, id: c.SampleID, language: c.Language, class: class, scenario: scenario,
			path: c.SamplePath, sha: sha, artifact: sourceArtifactID,
		}
	}
	infos := make([]sampleInfo, 0, len(byKey))
	for _, v := range byKey {
		infos = append(infos, v)
	}
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].language+"|"+infos[i].id+"|"+infos[i].path < infos[j].language+"|"+infos[j].id+"|"+infos[j].path
	})
	fingerprintParts := make([]string, 0, len(infos))
	now := nowUTC()
	sampleUIDs := make(map[string]string, len(infos))
	for _, info := range infos {
		fingerprintParts = append(fingerprintParts, info.language+"|"+info.id+"|"+info.path+"|"+info.sha)
		sampleUIDs[info.key] = info.uid
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO dataset_samples(sample_uid, sample_id, language, class, scenario, path, source_md5, source_sha256,
				source_artifact_id, risk_flags_json, created_at_utc, updated_at_utc)
			VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(language, sample_id, path) DO UPDATE SET
				class=excluded.class,
				scenario=excluded.scenario,
				source_sha256=excluded.source_sha256,
				source_artifact_id=excluded.source_artifact_id,
				updated_at_utc=excluded.updated_at_utc`,
			info.uid, info.id, info.language, nullString(info.class), nullString(info.scenario), info.path, nullString(info.md5), nullString(info.sha),
			nullString(info.artifact), "[]", now, now); err != nil {
			return nil, "", "", fmt.Errorf("upsert dataset sample: %w", err)
		}
	}
	fingerprint := hashString(strings.Join(fingerprintParts, "\n"))
	snapshotID := stableID("dataset_snapshot", fingerprint)
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO dataset_snapshots(snapshot_id, fingerprint, created_at_utc, sample_count, spec_json)
		VALUES(?, ?, ?, ?, ?)
		ON CONFLICT(fingerprint) DO UPDATE SET sample_count=excluded.sample_count, spec_json=excluded.spec_json`,
		snapshotID, fingerprint, now, len(infos), string(mustJSON(manifest.Spec))); err != nil {
		return nil, "", "", fmt.Errorf("upsert dataset snapshot: %w", err)
	}
	for _, info := range infos {
		if _, err := tx.ExecContext(ctx, `
			INSERT OR IGNORE INTO dataset_snapshot_members(snapshot_id, sample_uid) VALUES(?, ?)`,
			snapshotID, info.uid); err != nil {
			return nil, "", "", fmt.Errorf("upsert snapshot member: %w", err)
		}
	}
	return sampleUIDs, fingerprint, snapshotID, nil
}

func (s *SQLiteStore) putSourceArtifact(ctx context.Context, tx *sql.Tx, ictx *ingestContext, path string) (artifactID string, sha string, err error) {
	resolved := ictx.resolvePath(path)
	if !fileExists(resolved) {
		ictx.summary.UnavailableArtifacts++
		return "", "", nil
	}
	sha, _, err = fileSHA256(resolved)
	if err != nil {
		ictx.summary.UnavailableArtifacts++
		return "", "", nil
	}
	artifactID, err = s.putArtifact(ctx, tx, ictx, "dataset_source", resolved, "source")
	return artifactID, sha, err
}

func (s *SQLiteStore) putOptionalArtifact(ctx context.Context, tx *sql.Tx, ictx *ingestContext, kind, path, role string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", nil
	}
	resolved := ictx.resolvePath(path)
	if !fileExists(resolved) {
		ictx.summary.UnavailableArtifacts++
		return "", nil
	}
	return s.putArtifact(ctx, tx, ictx, kind, resolved, role)
}

func (s *SQLiteStore) putArtifact(ctx context.Context, tx *sql.Tx, ictx *ingestContext, kind, path, role string) (string, error) {
	resolved := ictx.resolvePath(path)
	sha, size, err := fileSHA256(resolved)
	if err != nil {
		ictx.summary.UnavailableArtifacts++
		return "", fmt.Errorf("hash artifact %s: %w", path, err)
	}
	artifactID := stableID("artifact", kind, sha)
	now := nowUTC()
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO artifacts(artifact_id, kind, path, size_bytes, sha256, redacted, created_at_utc)
		VALUES(?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(kind, sha256) DO UPDATE SET
			path=excluded.path,
			size_bytes=excluded.size_bytes,
			deleted_at_utc=NULL`,
		artifactID, kind, portablePath(resolved), size, sha, 0, now); err != nil {
		return "", fmt.Errorf("upsert artifact: %w", err)
	}
	if ictx != nil && ictx.runID != "" {
		if _, err := tx.ExecContext(ctx, `
			INSERT OR IGNORE INTO run_artifacts(run_id, artifact_id, role, created_at_utc)
			VALUES(?, ?, ?, ?)`,
			ictx.runID, artifactID, role, now); err != nil {
			return "", fmt.Errorf("link artifact: %w", err)
		}
		if _, seen := ictx.seen[artifactID]; !seen {
			ictx.seen[artifactID] = struct{}{}
			ictx.summary.ArtifactsIndexed++
		}
	}
	return artifactID, nil
}

func (s *SQLiteStore) ensureExperiment(ctx context.Context, tx *sql.Tx, id, runID string) error {
	now := nowUTC()
	_, err := tx.ExecContext(ctx, `
		INSERT INTO experiments(experiment_id, name, description, created_at_utc, updated_at_utc)
		VALUES(?, ?, ?, ?, ?)
		ON CONFLICT(experiment_id) DO UPDATE SET updated_at_utc=excluded.updated_at_utc`,
		id, runID, "auto-created from run_id", now, now)
	return err
}

func newIngestContext(runID, anchorPath string, spec contracts.RunSpec, summary *IngestSummary) *ingestContext {
	if runID == "" {
		runID = spec.RunID
	}
	runDir := ""
	if anchorPath != "" {
		runDir = inferRunDir(anchorPath, runID)
	}
	if summary == nil {
		summary = &IngestSummary{RunID: runID}
	}
	return &ingestContext{runID: runID, runDir: runDir, spec: spec, seen: map[string]struct{}{}, summary: summary}
}

func (c *ingestContext) resolvePath(path string) string {
	if path == "" {
		return path
	}
	if fileExists(path) {
		return path
	}
	p := filepath.FromSlash(path)
	if fileExists(p) {
		return p
	}
	if c.runDir != "" && c.runID != "" {
		marker := filepath.Join("runs", c.runID) + string(filepath.Separator)
		if idx := strings.Index(p, marker); idx >= 0 {
			candidate := filepath.Join(c.runDir, p[idx+len(marker):])
			if fileExists(candidate) {
				return candidate
			}
		}
	}
	if c.runDir != "" {
		runsDir := filepath.Dir(c.runDir)
		marker := string(filepath.Separator) + "runs" + string(filepath.Separator)
		if idx := strings.Index(p, marker); idx >= 0 {
			candidate := filepath.Join(runsDir, p[idx+len(marker):])
			if fileExists(candidate) {
				return candidate
			}
		}
		if strings.HasPrefix(p, "runs"+string(filepath.Separator)) {
			candidate := filepath.Join(filepath.Dir(runsDir), p)
			if fileExists(candidate) {
				return candidate
			}
		}
	}
	if !filepath.IsAbs(p) && c.runDir != "" {
		candidate := filepath.Join(c.runDir, p)
		if fileExists(candidate) {
			return candidate
		}
	}
	return p
}

func inferRunDir(anchorPath, runID string) string {
	abs, err := filepath.Abs(anchorPath)
	if err != nil {
		abs = anchorPath
	}
	if runID != "" {
		parts := strings.Split(filepath.Clean(abs), string(filepath.Separator))
		for i := len(parts) - 1; i >= 0; i-- {
			if parts[i] == runID {
				return strings.Join(parts[:i+1], string(filepath.Separator))
			}
		}
	}
	return filepath.Dir(filepath.Dir(abs))
}

func inferClassScenario(path, sampleID string) (string, string) {
	p := strings.ToLower(filepath.ToSlash(path))
	class := ""
	switch {
	case strings.Contains(p, "self_contained"):
		class = "self_contained"
	case strings.Contains(p, "module_level"):
		class = "module_level"
	}
	scenario := filepath.Base(filepath.Dir(path))
	if scenario == "." || scenario == "" {
		scenario = sampleID
		if idx := strings.LastIndex(scenario, "_"); idx > 0 {
			scenario = scenario[:idx]
		}
	}
	return class, scenario
}

func mergeIngestSummaries(a, b IngestSummary) IngestSummary {
	if a.RunID == "" {
		a.RunID = b.RunID
	}
	a.ManifestIngested = a.ManifestIngested || b.ManifestIngested
	a.EvaluationIngested = a.EvaluationIngested || b.EvaluationIngested
	a.ReportIngested = a.ReportIngested || b.ReportIngested
	if b.GenerationCases > a.GenerationCases {
		a.GenerationCases = b.GenerationCases
	}
	if b.EvaluationResults > a.EvaluationResults {
		a.EvaluationResults = b.EvaluationResults
	}
	a.ArtifactsIndexed += b.ArtifactsIndexed
	a.UnavailableArtifacts += b.UnavailableArtifacts
	return a
}

func stableID(parts ...string) string {
	return parts[0] + "_" + hashString(strings.Join(parts[1:], "\x00"))[:16]
}

func hashString(v string) string {
	sum := sha256.Sum256([]byte(v))
	return hex.EncodeToString(sum[:])
}

func fileSHA256(path string) (string, int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()
	h := sha256.New()
	n, err := io.Copy(h, f)
	if err != nil {
		return "", 0, err
	}
	return hex.EncodeToString(h.Sum(nil)), n, nil
}

func portablePath(path string) string {
	if path == "" {
		return path
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	cwd, err := os.Getwd()
	if err == nil {
		if rel, relErr := filepath.Rel(cwd, abs); relErr == nil && rel != "." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && rel != ".." {
			return filepath.ToSlash(rel)
		}
	}
	return filepath.ToSlash(abs)
}

func resolveStoredPath(path string) string {
	if strings.TrimSpace(path) == "" {
		return ""
	}
	p := filepath.FromSlash(path)
	if fileExists(p) {
		return p
	}
	if !filepath.IsAbs(p) {
		if cwd, err := os.Getwd(); err == nil {
			candidate := filepath.Join(cwd, p)
			if fileExists(candidate) {
				return candidate
			}
		}
	}
	return p
}

func readJSON(path string, dst any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

func mustJSON(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		return []byte("{}")
	}
	return data
}

func defaultExperimentID(runID string) string {
	return stableID("experiment", runID)
}

func caseKey(model, language, sampleID string) string {
	return model + "|" + language + "|" + sampleID
}

func nullableInt(v *int) any {
	if v == nil {
		return nil
	}
	return *v
}

func nullableFloat(v *float64) any {
	if v == nil {
		return nil
	}
	return *v
}

func nullString(v string) any {
	if strings.TrimSpace(v) == "" || v == "null" {
		return nil
	}
	return v
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func ptrBoolToNullableInt(v *bool) any {
	if v == nil {
		return nil
	}
	if *v {
		return 1
	}
	return 0
}

func nullableIntBool(v sql.NullInt64) *bool {
	if !v.Valid {
		return nil
	}
	out := v.Int64 != 0
	return &out
}

func nullableSQLInt(v sql.NullInt64) *int {
	if !v.Valid {
		return nil
	}
	out := int(v.Int64)
	return &out
}

func nullableSQLFloat(v sql.NullFloat64) *float64 {
	if !v.Valid {
		return nil
	}
	out := v.Float64
	return &out
}

func nonEmptyStrings(values []string) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

func nowUTC() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339Nano)
}

func fileExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

// Phase 3: Database management list methods

func (s *SQLiteStore) ListGenerationRuns(ctx context.Context, limit int) ([]DBGenerationRunItem, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT run_id, COALESCE(experiment_id, ''), schema_version, COALESCE(created_at_utc, ''),
		       COALESCE(prompt_strategy, ''), COALESCE(prompt_version_id, ''),
		       COALESCE(dataset_fingerprint, ''), created_db_at_utc
		FROM generation_runs
		ORDER BY created_db_at_utc DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DBGenerationRunItem
	for rows.Next() {
		var r DBGenerationRunItem
		if err := rows.Scan(&r.RunID, &r.ExperimentID, &r.SchemaVersion, &r.CreatedAtUTC,
			&r.PromptStrategy, &r.PromptVersionID, &r.DatasetFingerprint, &r.CreatedDBAtUTC); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ListGeneratedCases(ctx context.Context, runID, model, language string, limit int) ([]DBGeneratedCaseItem, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	where := []string{"1=1"}
	args := []any{}
	if runID != "" {
		where = append(where, "run_id = ?")
		args = append(args, runID)
	}
	if model != "" {
		where = append(where, "model = ?")
		args = append(args, model)
	}
	if language != "" {
		where = append(where, "language = ?")
		args = append(args, language)
	}
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx, `
		SELECT generated_case_id, run_id, model, language, sample_id, success, truncated,
		       latency_ms, prompt_tokens, completion_tokens, COALESCE(generated_at_utc, '')
		FROM generated_cases
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY created_db_at_utc DESC
		LIMIT ?`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DBGeneratedCaseItem
	for rows.Next() {
		var r DBGeneratedCaseItem
		var success, truncated int
		var latency, promptTok, compTok sql.NullInt64
		if err := rows.Scan(&r.GeneratedCaseID, &r.RunID, &r.Model, &r.Language, &r.SampleID,
			&success, &truncated, &latency, &promptTok, &compTok, &r.GeneratedAtUTC); err != nil {
			return nil, err
		}
		r.Success = success != 0
		r.Truncated = truncated != 0
		r.LatencyMS = nullableSQLInt(latency)
		r.PromptTokens = nullableSQLInt(promptTok)
		r.CompletionTokens = nullableSQLInt(compTok)
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ListPromptRenderings(ctx context.Context, runID string, limit int) ([]DBPromptRenderingItem, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	where := "1=1"
	args := []any{}
	if runID != "" {
		where = "run_id = ?"
		args = append(args, runID)
	}
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx, `
		SELECT prompt_rendering_id, run_id, model, language, sample_id,
		       COALESCE(prompt_version_id, ''), COALESCE(prompt_mode, ''), created_at_utc
		FROM prompt_renderings
		WHERE `+where+`
		ORDER BY created_at_utc DESC
		LIMIT ?`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DBPromptRenderingItem
	for rows.Next() {
		var r DBPromptRenderingItem
		if err := rows.Scan(&r.PromptRenderingID, &r.RunID, &r.Model, &r.Language, &r.SampleID,
			&r.PromptVersionID, &r.PromptMode, &r.CreatedAtUTC); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ListEvaluationRuns(ctx context.Context, runID string, limit int) ([]DBEvaluationRunItem, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	where := "1=1"
	args := []any{}
	if runID != "" {
		where = "run_id = ?"
		args = append(args, runID)
	}
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx, `
		SELECT evaluation_run_id, run_id, schema_version, COALESCE(evaluated_at_utc, ''),
		       COALESCE(env_id, ''), COALESCE(score_policy_id, ''), created_db_at_utc
		FROM evaluation_runs
		WHERE `+where+`
		ORDER BY evaluated_at_utc DESC
		LIMIT ?`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DBEvaluationRunItem
	for rows.Next() {
		var r DBEvaluationRunItem
		if err := rows.Scan(&r.EvaluationRunID, &r.RunID, &r.SchemaVersion, &r.EvaluatedAtUTC,
			&r.EnvID, &r.ScorePolicyID, &r.CreatedDBAtUTC); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ListEvaluationStages(ctx context.Context, evaluationRunID string, limit int) ([]DBEvaluationStageItem, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	where := "1=1"
	args := []any{}
	if evaluationRunID != "" {
		where = "er.evaluation_run_id = ?"
		args = append(args, evaluationRunID)
	}
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx, `
		SELECT sr.stage_result_id, sr.evaluation_result_id, sr.stage, sr.status,
		       sr.exit_code, sr.duration_ms, sr.created_at_utc
		FROM evaluation_stage_results sr
		JOIN evaluation_results er ON er.evaluation_result_id = sr.evaluation_result_id
		WHERE `+where+`
		ORDER BY sr.created_at_utc DESC
		LIMIT ?`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DBEvaluationStageItem
	for rows.Next() {
		var r DBEvaluationStageItem
		var exitCode, duration sql.NullInt64
		if err := rows.Scan(&r.StageResultID, &r.EvaluationResultID, &r.Stage, &r.Status,
			&exitCode, &duration, &r.CreatedAtUTC); err != nil {
			return nil, err
		}
		r.ExitCode = nullableSQLInt(exitCode)
		r.DurationMS = nullableSQLInt(duration)
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ListDatasetSamples(ctx context.Context, language, class string, limit int) ([]DBDatasetSampleItem, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	where := []string{"1=1"}
	args := []any{}
	if language != "" {
		where = append(where, "language = ?")
		args = append(args, language)
	}
	if class != "" {
		where = append(where, "class = ?")
		args = append(args, class)
	}
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx, `
		SELECT sample_uid, sample_id, language, COALESCE(class, ''), COALESCE(scenario, ''),
		       path, created_at_utc
		FROM dataset_samples
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY created_at_utc DESC
		LIMIT ?`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DBDatasetSampleItem
	for rows.Next() {
		var r DBDatasetSampleItem
		if err := rows.Scan(&r.SampleUID, &r.SampleID, &r.Language, &r.Class, &r.Scenario,
			&r.Path, &r.CreatedAtUTC); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ListDatasetSnapshots(ctx context.Context, limit int) ([]DBDatasetSnapshotItem, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT snapshot_id, fingerprint, sample_count, created_at_utc
		FROM dataset_snapshots
		ORDER BY created_at_utc DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DBDatasetSnapshotItem
	for rows.Next() {
		var r DBDatasetSnapshotItem
		if err := rows.Scan(&r.SnapshotID, &r.Fingerprint, &r.SampleCount, &r.CreatedAtUTC); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ListModelConfigs(ctx context.Context, limit int) ([]DBModelConfigItem, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT model_config_id, model_name, COALESCE(provider, ''), COALESCE(model_id, ''), created_at_utc
		FROM model_configs
		ORDER BY created_at_utc DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DBModelConfigItem
	for rows.Next() {
		var r DBModelConfigItem
		if err := rows.Scan(&r.ModelConfigID, &r.ModelName, &r.Provider, &r.ModelID, &r.CreatedAtUTC); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ListPromptProfiles(ctx context.Context, limit int) ([]DBPromptProfileItem, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT profile_id, COALESCE(strategy, ''), COALESCE(version_id, ''), created_at_utc
		FROM prompt_profiles
		ORDER BY created_at_utc DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DBPromptProfileItem
	for rows.Next() {
		var r DBPromptProfileItem
		if err := rows.Scan(&r.ProfileID, &r.Strategy, &r.VersionID, &r.CreatedAtUTC); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ListEvaluationEnvs(ctx context.Context, limit int) ([]DBEvaluationEnvItem, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT env_id, fingerprint, created_at_utc
		FROM evaluation_envs
		ORDER BY created_at_utc DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DBEvaluationEnvItem
	for rows.Next() {
		var r DBEvaluationEnvItem
		if err := rows.Scan(&r.EnvID, &r.Fingerprint, &r.CreatedAtUTC); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ListScorePolicies(ctx context.Context, limit int) ([]DBScorePolicyItem, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT score_policy_id, name, created_at_utc
		FROM score_policies
		ORDER BY created_at_utc DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DBScorePolicyItem
	for rows.Next() {
		var r DBScorePolicyItem
		if err := rows.Scan(&r.ScorePolicyID, &r.Name, &r.CreatedAtUTC); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ListReports(ctx context.Context, runID string, limit int) ([]DBReportItem, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	where := "1=1"
	args := []any{}
	if runID != "" {
		where = "run_id = ?"
		args = append(args, runID)
	}
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx, `
		SELECT report_id, run_id, COALESCE(generated_at_utc, ''), created_db_at_utc
		FROM report_snapshots
		WHERE `+where+`
		ORDER BY generated_at_utc DESC
		LIMIT ?`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DBReportItem
	for rows.Next() {
		var r DBReportItem
		if err := rows.Scan(&r.ReportID, &r.RunID, &r.GeneratedAtUTC, &r.CreatedDBAtUTC); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ListRunArtifacts(ctx context.Context, runID string, limit int) ([]DBRunArtifactItem, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	where := "1=1"
	args := []any{}
	if runID != "" {
		where = "run_id = ?"
		args = append(args, runID)
	}
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx, `
		SELECT run_id, artifact_id, role, created_at_utc
		FROM run_artifacts
		WHERE `+where+`
		ORDER BY created_at_utc DESC
		LIMIT ?`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DBRunArtifactItem
	for rows.Next() {
		var r DBRunArtifactItem
		if err := rows.Scan(&r.RunID, &r.ArtifactID, &r.Role, &r.CreatedAtUTC); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ListExperiments(ctx context.Context, limit int) ([]DBExperimentItem, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT experiment_id, name, COALESCE(description, ''), created_at_utc, updated_at_utc
		FROM experiments
		ORDER BY created_at_utc DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DBExperimentItem
	for rows.Next() {
		var r DBExperimentItem
		if err := rows.Scan(&r.ExperimentID, &r.Name, &r.Description, &r.CreatedAtUTC, &r.UpdatedAtUTC); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// ModelConfigInfo 模型配置信息（用于入库）
type ModelConfigInfo struct {
	Provider    string         `json:"provider"`
	ModelID     string         `json:"model_id"`
	APIEndpoint string         `json:"api_endpoint"`
	Parameters  map[string]any `json:"parameters,omitempty"`
}

// readModelConfigsFromYAML 从models.yaml读取模型配置
func (s *SQLiteStore) readModelConfigsFromYAML(configPath string) map[string]ModelConfigInfo {
	result := make(map[string]ModelConfigInfo)
	if configPath == "" {
		return result
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return result
	}
	var root map[string]any
	if err := yaml.Unmarshal(data, &root); err != nil {
		return result
	}
	models, ok := root["models"].(map[string]any)
	if !ok {
		return result
	}
	for name, v := range models {
		node, ok := v.(map[string]any)
		if !ok {
			continue
		}
		cfg := ModelConfigInfo{}
		cfg.Provider, _ = node["provider"].(string)
		if config, ok := node["config"].(map[string]any); ok {
			cfg.ModelID, _ = config["model"].(string)
			cfg.APIEndpoint, _ = config["api_endpoint"].(string)
			if params, ok := config["parameters"].(map[string]any); ok {
				cfg.Parameters = params
			}
		}
		result[name] = cfg
	}
	return result
}
