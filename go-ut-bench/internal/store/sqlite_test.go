package store

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go-ut-bench/internal/contracts"

	_ "modernc.org/sqlite"
)

func TestIngestEvaluationUpsertV2(t *testing.T) {
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
	if err := db.QueryRow(`SELECT COUNT(*) FROM evaluation_results WHERE run_id='run_1' AND model='deepseek' AND language='python' AND sample_id='s1'`).Scan(&cnt); err != nil {
		t.Fatal(err)
	}
	if cnt != 1 {
		t.Fatalf("expected one upserted row, got %d", cnt)
	}

	var gotRate float64
	if err := db.QueryRow(`SELECT test_pass_rate FROM evaluation_results WHERE run_id='run_1' AND model='deepseek' AND language='python' AND sample_id='s1'`).Scan(&gotRate); err != nil {
		t.Fatal(err)
	}
	if gotRate != 1.0 {
		t.Fatalf("expected updated test_pass_rate=1.0, got %f", gotRate)
	}
}

func TestIngestManifestAndOverview(t *testing.T) {
	ctx := context.Background()
	tmp := t.TempDir()
	dbPath := filepath.Join(tmp, "utbench.db")
	runDir := filepath.Join(tmp, "artifacts", "runs", "run_manifest")
	srcPath := filepath.Join(tmp, "datasets", "python", "python_code_files_self_contained", "simple_function", "simple_function_000.py")
	promptPath := filepath.Join(runDir, "generated", "prompts", "rendered", "deepseek", "python", "simple_function_000.md")
	testPath := filepath.Join(runDir, "generated", "tests", "deepseek", "python", "simple_function_000_test.py")
	responsePath := filepath.Join(runDir, "generated", "responses", "deepseek", "python", "simple_function_000.json")
	metadataPath := filepath.Join(runDir, "generated", "metadata", "deepseek", "python", "simple_function_000.json")
	writeTestFile(t, srcPath, "def add(a, b):\n    return a + b\n")
	writeTestFile(t, promptPath, "Write tests")
	writeTestFile(t, testPath, "from module_under_test import add\n")
	writeTestFile(t, responsePath, `{"content":"ok"}`)
	writeTestFile(t, metadataPath, `{"latency_ms":12}`)

	manifest := contracts.GeneratedManifest{
		SchemaVersion:   contracts.SchemaVersion,
		RunID:           "run_manifest",
		CreatedAtUTC:    time.Now().UTC(),
		PromptStrategy:  "fairness-v1",
		PromptVersionID: "prompt-v1",
		Spec: contracts.RunSpec{
			RunID:          "run_manifest",
			Models:         []string{"deepseek"},
			Languages:      []string{"python"},
			DatasetClasses: []string{"self_contained"},
			OutputRoot:     filepath.Join(tmp, "artifacts"),
		},
		Cases: []contracts.GeneratedCase{
			{
				Model:             "deepseek",
				Language:          "python",
				SampleID:          "simple_function_000",
				SamplePath:        srcPath,
				PromptVersionID:   "prompt-v1",
				PromptMode:        "self_contained",
				PromptPath:        promptPath,
				GeneratedTestPath: testPath,
				ResponsePath:      responsePath,
				MetadataPath:      metadataPath,
				GeneratedAtUTC:    time.Now().UTC(),
				Success:           true,
			},
		},
	}
	manifestPath := filepath.Join(runDir, "generated", "generated_manifest.json")
	if err := contracts.WriteJSON(manifestPath, manifest); err != nil {
		t.Fatal(err)
	}
	pass := true
	evalPath := filepath.Join(runDir, "evaluation", "evaluation_result.json")
	evalSet := contracts.EvaluationResultSet{
		SchemaVersion:  contracts.SchemaVersion,
		RunID:          "run_manifest",
		EvaluatedAtUTC: time.Now().UTC(),
		ManifestPath:   manifestPath,
		Results: []contracts.EvaluationResult{
			{
				Model:             "deepseek",
				Language:          "python",
				SampleID:          "simple_function_000",
				GeneratedTestPath: testPath,
				SourcePath:        srcPath,
				CompilePass:       true,
				TestPass:          &pass,
			},
		},
	}
	if err := contracts.WriteJSON(evalPath, evalSet); err != nil {
		t.Fatal(err)
	}
	reportPath := filepath.Join(runDir, "report", "report_summary.json")
	report := contracts.ReportPayload{
		SchemaVersion:    contracts.SchemaVersion,
		RunID:            "run_manifest",
		GeneratedAtUTC:   time.Now().UTC(),
		SourceEvaluation: evalPath,
		Summary:          contracts.ReportSummary{TotalSamples: 1, EligibleSamples: 1},
	}
	if err := contracts.WriteJSON(reportPath, report); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(runDir, "report", "report.html"), "<html></html>")

	s, err := OpenSQLite(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if err := s.Init(ctx); err != nil {
		t.Fatal(err)
	}
	sum, err := s.IngestManifestFile(ctx, manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if !sum.ManifestIngested || sum.GenerationCases != 1 {
		t.Fatalf("unexpected summary: %+v", sum)
	}
	runSum, err := s.IngestRun(ctx, IngestRunOptions{RunDir: runDir})
	if err != nil {
		t.Fatal(err)
	}
	if !runSum.EvaluationIngested || !runSum.ReportIngested || runSum.EvaluationResults != 1 {
		t.Fatalf("unexpected run summary: %+v", runSum)
	}
	overview, err := s.Overview(ctx, 10)
	if err != nil {
		t.Fatal(err)
	}
	if overview.GenerationRuns != 1 || overview.GeneratedCases != 1 || overview.EvaluationResults != 1 || overview.Reports != 1 || overview.Artifacts < 7 {
		t.Fatalf("unexpected overview: %+v", overview)
	}

	selected, err := s.SelectEvaluationResultSet(ctx, "db_report", DBReportFilter{RunIDs: []string{"run_manifest"}})
	if err != nil {
		t.Fatal(err)
	}
	if selected.RunID != "db_report" || len(selected.Results) != 1 {
		t.Fatalf("unexpected selected result set: %+v", selected)
	}
	if selected.Results[0].GeneratedTestPath == "" || selected.Results[0].SourcePath == "" {
		t.Fatalf("expected selected rows to keep artifact paths: %+v", selected.Results[0])
	}
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
