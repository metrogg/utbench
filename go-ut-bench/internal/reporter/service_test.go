package reporter

import (
	"testing"

	"go-ut-bench/internal/contracts"
)

func TestBuildMutationBreakdown(t *testing.T) {
	total1 := 10
	killed1 := 6
	survived1 := 2
	noTests1 := 1
	timeouts1 := 0
	skipped1 := 0
	suspicious1 := 1

	total2 := 5
	killed2 := 3

	rows := []contracts.EvaluationResult{
		{
			MutationTotal:      &total1,
			MutationKilled:     &killed1,
			MutationSurvived:   &survived1,
			MutationNoTests:    &noTests1,
			MutationTimeouts:   &timeouts1,
			MutationSkipped:    &skipped1,
			MutationSuspicious: &suspicious1,
		},
		{
			MutationTotal:  &total2,
			MutationKilled: &killed2,
		},
	}

	b := buildMutationBreakdown(rows)
	if b.Total != 15 || b.Killed != 9 || b.Survived != 2 || b.NoTests != 1 || b.Suspicious != 1 {
		t.Fatalf("unexpected breakdown: %+v", b)
	}
}

func TestBuildSummaryUsesSampleLevelTestPassRateWhenCountsMissing(t *testing.T) {
	pass := true
	fail := false
	lineCov := 0.5

	rows := []contracts.EvaluationResult{
		{
			CompilePass:  true,
			TestPass:     &pass,
			LineCoverage: &lineCov,
		},
		{
			CompilePass: true,
			TestPass:    &fail,
		},
	}

	s := buildSummary(rows)
	if s.TotalSamples != 2 {
		t.Fatalf("expected total samples 2, got %d", s.TotalSamples)
	}
	if s.CompilePassRate != 1 {
		t.Fatalf("expected compile pass rate 1.0, got %v", s.CompilePassRate)
	}
	if s.TestPassCount != 1 {
		t.Fatalf("expected test pass count 1, got %d", s.TestPassCount)
	}
	if s.TestPassRate != 0.5 {
		t.Fatalf("expected test pass rate 0.5, got %v", s.TestPassRate)
	}
}
