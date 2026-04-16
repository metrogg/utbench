package store

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"go-ut-bench/internal/contracts"

	_ "modernc.org/sqlite"
)

func TestIngestEvaluationUpsert(t *testing.T) {
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "utbench.db")

	s, err := OpenSQLite(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	if err := s.Init(ctx); err != nil {
		t.Fatal(err)
	}

	pass := true
	pc := 1
	tc := 2
	pr := 0.5
	set := contracts.EvaluationResultSet{
		SchemaVersion:  contracts.SchemaVersion,
		RunID:          "run_1",
		EvaluatedAtUTC: time.Now().UTC(),
		Results: []contracts.EvaluationResult{
			{
				Model:          "deepseek",
				Language:       "python",
				SampleID:       "s1",
				CompilePass:    true,
				TestPass:       &pass,
				TestPassCount:  &pc,
				TestTotalCount: &tc,
				TestPassRate:   &pr,
			},
		},
	}
	if err := s.IngestEvaluation(ctx, set); err != nil {
		t.Fatal(err)
	}

	pc2 := 2
	tc2 := 2
	pr2 := 1.0
	set.Results[0].TestPassCount = &pc2
	set.Results[0].TestTotalCount = &tc2
	set.Results[0].TestPassRate = &pr2
	if err := s.IngestEvaluation(ctx, set); err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var cnt int
	if err := db.QueryRow(`SELECT COUNT(*) FROM sample_results WHERE run_id='run_1' AND model='deepseek' AND language='python' AND sample_id='s1'`).Scan(&cnt); err != nil {
		t.Fatal(err)
	}
	if cnt != 1 {
		t.Fatalf("expected one upserted row, got %d", cnt)
	}

	var gotRate float64
	if err := db.QueryRow(`SELECT test_pass_rate FROM sample_results WHERE run_id='run_1' AND model='deepseek' AND language='python' AND sample_id='s1'`).Scan(&gotRate); err != nil {
		t.Fatal(err)
	}
	if gotRate != 1.0 {
		t.Fatalf("expected updated test_pass_rate=1.0, got %f", gotRate)
	}
}
