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
