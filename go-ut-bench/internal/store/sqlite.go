package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go-ut-bench/internal/contracts"

	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db   *sql.DB
	path string
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
		`CREATE TABLE IF NOT EXISTS runs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			run_id TEXT NOT NULL UNIQUE,
			schema_version TEXT NOT NULL,
			evaluated_at_utc TEXT NOT NULL,
			created_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS sample_results (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			run_id TEXT NOT NULL,
			model TEXT NOT NULL,
			language TEXT NOT NULL,
			sample_id TEXT NOT NULL,
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
			compile_error TEXT,
			test_error TEXT,
			coverage_error TEXT,
			mutation_error TEXT,
			generated_test_path TEXT,
			source_path TEXT,
			UNIQUE(run_id, model, language, sample_id)
		);`,
	}

	for _, stmt := range ddl {
		if _, err := s.db.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}

	alterDDLs := []string{
		`ALTER TABLE sample_results ADD COLUMN test_pass_count INTEGER;`,
		`ALTER TABLE sample_results ADD COLUMN test_total_count INTEGER;`,
		`ALTER TABLE sample_results ADD COLUMN test_pass_rate REAL;`,
		`ALTER TABLE sample_results ADD COLUMN mutation_total INTEGER;`,
		`ALTER TABLE sample_results ADD COLUMN mutation_killed INTEGER;`,
		`ALTER TABLE sample_results ADD COLUMN mutation_survived INTEGER;`,
		`ALTER TABLE sample_results ADD COLUMN mutation_no_tests INTEGER;`,
		`ALTER TABLE sample_results ADD COLUMN mutation_timeouts INTEGER;`,
		`ALTER TABLE sample_results ADD COLUMN mutation_skipped INTEGER;`,
		`ALTER TABLE sample_results ADD COLUMN mutation_suspicious INTEGER;`,
	}
	for _, stmt := range alterDDLs {
		_, _ = s.db.ExecContext(ctx, stmt)
	}
	return nil
}

func (s *SQLiteStore) IngestEvaluation(ctx context.Context, set contracts.EvaluationResultSet) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO runs (run_id, schema_version, evaluated_at_utc, created_at_utc)
		 VALUES (?, ?, ?, ?)
		 ON CONFLICT(run_id) DO UPDATE SET
		   schema_version=excluded.schema_version,
		   evaluated_at_utc=excluded.evaluated_at_utc`,
		set.RunID,
		set.SchemaVersion,
		set.EvaluatedAtUTC.Format(time.RFC3339Nano),
		time.Now().UTC().Format(time.RFC3339Nano),
	); err != nil {
		return fmt.Errorf("insert run: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO sample_results (
			run_id, model, language, sample_id,
			compile_pass, test_pass, test_pass_count, test_total_count, test_pass_rate,
			line_coverage, branch_coverage, mutation_score,
			mutation_total, mutation_killed, mutation_survived, mutation_no_tests, mutation_timeouts, mutation_skipped, mutation_suspicious,
			assertion_count, test_case_count, assertion_density, runtime_ms,
			compile_error, test_error, coverage_error, mutation_error,
			generated_test_path, source_path
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(run_id, model, language, sample_id) DO UPDATE SET
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
			compile_error=excluded.compile_error,
			test_error=excluded.test_error,
			coverage_error=excluded.coverage_error,
			mutation_error=excluded.mutation_error,
			generated_test_path=excluded.generated_test_path,
			source_path=excluded.source_path
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, row := range set.Results {
		if _, err := stmt.ExecContext(
			ctx,
			set.RunID,
			row.Model,
			row.Language,
			row.SampleID,
			boolToInt(row.CompilePass),
			ptrBoolToNullableInt(row.TestPass),
			row.TestPassCount,
			row.TestTotalCount,
			row.TestPassRate,
			row.LineCoverage,
			row.BranchCoverage,
			row.MutationScore,
			row.MutationTotal,
			row.MutationKilled,
			row.MutationSurvived,
			row.MutationNoTests,
			row.MutationTimeouts,
			row.MutationSkipped,
			row.MutationSuspicious,
			row.AssertionCount,
			row.TestCaseCount,
			row.AssertionDensity,
			row.RuntimeMS,
			nullIfEmpty(row.CompileError),
			nullIfEmpty(row.TestError),
			nullIfEmpty(row.CoverageError),
			nullIfEmpty(row.MutationError),
			row.GeneratedTestPath,
			row.SourcePath,
		); err != nil {
			return fmt.Errorf("insert sample result: %w", err)
		}
	}

	return tx.Commit()
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

func nullIfEmpty(v string) any {
	if v == "" {
		return nil
	}
	return v
}
