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
			MutationTool:       "mutmut",
		},
		{
			MutationTotal:  &total2,
			MutationKilled: &killed2,
			MutationTool:   "mull",
		},
	}

	b := buildMutationBreakdown(rows)
	if b.Total != 15 || b.Killed != 9 || b.Survived != 2 || b.NoTests != 1 || b.Suspicious != 1 {
		t.Fatalf("unexpected breakdown: %+v", b)
	}
	if len(b.ByTool) != 2 {
		t.Fatalf("expected two tool breakdowns, got %+v", b.ByTool)
	}
	if b.ByTool[0].Tool != "mull" || b.ByTool[0].Total != 5 || b.ByTool[0].Killed != 3 {
		t.Fatalf("unexpected first tool breakdown: %+v", b.ByTool[0])
	}
	if b.ByTool[1].Tool != "mutmut" || b.ByTool[1].Total != 10 || b.ByTool[1].Killed != 6 {
		t.Fatalf("unexpected second tool breakdown: %+v", b.ByTool[1])
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

func TestScoreEligibilityExcludesNonModelFailuresFromRanking(t *testing.T) {
	pass := true
	fail := false
	eligible := true
	excluded := false
	lineCov := 1.0

	rows := []contracts.EvaluationResult{
		{
			Model:         "m1",
			CompilePass:   true,
			TestPass:      &pass,
			LineCoverage:  &lineCov,
			ScoreEligible: &eligible,
		},
		{
			Model:                "m1",
			CompilePass:          false,
			TestPass:             &fail,
			ScoreEligible:        &excluded,
			FailureOrigin:        "environment",
			ScoreExclusionReason: "mvn not installed",
			CompileError:         "mvn not installed",
		},
	}

	s := buildSummary(rows)
	if s.TotalSamples != 2 || s.EligibleSamples != 1 || s.ExcludedSamples != 1 {
		t.Fatalf("unexpected summary eligibility counts: %+v", s)
	}
	if s.CompilePassRate != 1 {
		t.Fatalf("expected excluded row to be omitted from compile rate, got %v", s.CompilePassRate)
	}

	dims := buildDimensions(rows, nil)
	if len(dims.ByModel) != 1 || dims.ByModel[0].TotalSamples != 1 || dims.ByModel[0].CompilePassRate != 1 {
		t.Fatalf("unexpected eligible model dimensions: %+v", dims.ByModel)
	}

	exclusions := buildScoreExclusions(rows)
	if len(exclusions) != 1 || exclusions[0].Origin != "environment" || exclusions[0].Count != 1 {
		t.Fatalf("unexpected exclusions: %+v", exclusions)
	}
}
