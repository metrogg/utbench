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

func TestBuildSummarySeparatesSampleAndTestCasePassRate(t *testing.T) {
	pass := true
	fail := false
	lineCov := 0.8
	casePass1, caseTotal1 := 1, 2
	casePass2, caseTotal2 := 3, 3

	rows := []contracts.EvaluationResult{
		{
			Model:          "m1",
			Language:       "python",
			SampleID:       "boundary_000",
			CompilePass:    true,
			TestPass:       &fail,
			TestPassCount:  &casePass1,
			TestTotalCount: &caseTotal1,
			LineCoverage:   &lineCov,
		},
		{
			Model:          "m1",
			Language:       "python",
			SampleID:       "boundary_001",
			CompilePass:    true,
			TestPass:       &pass,
			TestPassCount:  &casePass2,
			TestTotalCount: &caseTotal2,
			LineCoverage:   &lineCov,
		},
	}

	s := buildSummary(rows)
	if s.SampleTestPassCount != 1 || s.SampleTestPassRate != 0.5 {
		t.Fatalf("expected sample-level pass rate 0.5, got %+v", s)
	}
	if s.TestPassCount != 1 || s.TestPassRate != 0.5 {
		t.Fatalf("expected compatibility sample-level test metrics, got %+v", s)
	}
	if s.TestCasePassCount != 4 || s.TestCasePassRate != (4.0/5.0) {
		t.Fatalf("expected case-level pass rate 0.8, got %+v", s)
	}

	dims := buildDimensions(rows, nil)
	if len(dims.ByModel) != 1 {
		t.Fatalf("expected single model dimension, got %+v", dims.ByModel)
	}
	if dims.ByModel[0].AvgTestPassRate != 0.5 {
		t.Fatalf("expected model sample-level rate 0.5, got %+v", dims.ByModel[0])
	}
	if dims.ByModel[0].AvgTestCasePassRate != 0.8 {
		t.Fatalf("expected model case-level rate 0.8, got %+v", dims.ByModel[0])
	}

	top := buildTopModels(dims.ByModel)
	if len(top) != 1 || top[0].AvgTestPassRate != 0.5 || top[0].AvgTestCasePassRate != 0.8 {
		t.Fatalf("unexpected top model metrics: %+v", top)
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

func TestClassifyMutationError(t *testing.T) {
	cases := map[string]string{
		"Mull: baseline tests failed, skipping mutation":                           "mutation_skipped_baseline_failed",
		"mutmut: generated tests do not import mutation target, skipping mutation": "mutation_target_not_exercised",
		"go-mutesting: go-mutesting no results to report":                          "mutation_no_results",
		"pitest: pitest no killed/survived results":                                "mutation_no_coverage",
		"mutmut produced zero mutants":                                             "mutation_no_effective_mutants",
		"mull timed out after 120s":                                                "mutation_timeout",
		"pitest parse error: stats file not found":                                 "mutation_tool_error",
		"unexpected mutation failure from generated tests":                         "mutation_error",
	}

	for input, want := range cases {
		if got := classifyMutationError(input); got != want {
			t.Fatalf("classifyMutationError(%q)=%q, want %q", input, got, want)
		}
	}
}

func TestBuildFailureRowsUsesMutationErrorCategories(t *testing.T) {
	rows := []contracts.EvaluationResult{
		{Model: "m1", SampleID: "s1", MutationError: "Mull: baseline tests failed, skipping mutation"},
		{Model: "m1", SampleID: "s2", MutationError: "go-mutesting no results to report"},
	}

	failures := buildFailureRows(rows)
	got := map[string]int{}
	for _, f := range failures {
		got[f.ErrorType] = f.Count
	}
	if got["mutation_skipped_baseline_failed"] != 1 || got["mutation_no_results"] != 1 {
		t.Fatalf("expected categorized mutation failures, got %+v", failures)
	}
	if got["mutation_error"] != 0 {
		t.Fatalf("did not expect generic mutation_error bucket, got %+v", failures)
	}
}

func TestBuildAgentAndSkillComparisons(t *testing.T) {
	pass := true
	baseLine := 0.5
	agentLine := 0.8
	skillLine := 0.9
	baseMut := 0.4
	agentMut := 0.7
	skillMut := 0.85
	baseLatency := 100
	agentLatency := 160
	skillLatency := 170

	rows := []contracts.EvaluationResult{
		{
			Model:          "model_api__deepseek__no_skill",
			SubjectID:      "model_api__deepseek__no_skill",
			AgentFramework: "model_api",
			AgentModel:     "deepseek",
			SkillName:      "no_skill",
			Language:       "python",
			SampleID:       "boundary_000",
			CompilePass:    true,
			TestPass:       &pass,
			LineCoverage:   &baseLine,
			MutationScore:  &baseMut,
			LatencyMS:      &baseLatency,
		},
		{
			Model:          "aider__deepseek__no_skill",
			SubjectID:      "aider__deepseek__no_skill",
			SubjectKind:    "cli_agent",
			AgentFramework: "aider",
			AgentModel:     "deepseek",
			SkillName:      "no_skill",
			Language:       "python",
			SampleID:       "boundary_000",
			CompilePass:    true,
			TestPass:       &pass,
			LineCoverage:   &agentLine,
			MutationScore:  &agentMut,
			LatencyMS:      &agentLatency,
		},
		{
			Model:          "aider__deepseek__unit_test_skill",
			SubjectID:      "aider__deepseek__unit_test_skill",
			SubjectKind:    "cli_agent",
			AgentFramework: "aider",
			AgentModel:     "deepseek",
			SkillName:      "unit_test_skill",
			SkillVersion:   "1",
			Language:       "python",
			SampleID:       "boundary_000",
			CompilePass:    true,
			TestPass:       &pass,
			LineCoverage:   &skillLine,
			MutationScore:  &skillMut,
			LatencyMS:      &skillLatency,
		},
	}

	agentRows := buildAgentComparisons(rows)
	if len(agentRows) != 1 {
		t.Fatalf("expected one agent comparison, got %+v", agentRows)
	}
	if agentRows[0].LineCoverageDelta != 0.3 || agentRows[0].MutationScoreDelta != 0.3 || agentRows[0].LatencyMSDelta != 60 {
		t.Fatalf("unexpected agent comparison: %+v", agentRows[0])
	}

	skillRows := buildSkillUplifts(rows)
	if len(skillRows) != 1 {
		t.Fatalf("expected one skill uplift, got %+v", skillRows)
	}
	if skillRows[0].LineCoverageDelta != 0.1 || skillRows[0].MutationScoreDelta != 0.15 || skillRows[0].LatencyMSDelta != 10 {
		t.Fatalf("unexpected skill uplift: %+v", skillRows[0])
	}
}
