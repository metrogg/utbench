// store 包提供 SQLite 数据库存储功能
// 本文件包含查询相关的函数
package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"go-ut-bench/internal/contracts"
)

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
			er.model, COALESCE(er.subject_id, ''), COALESCE(er.subject_kind, ''), COALESCE(er.agent_framework, ''),
			COALESCE(er.agent_model, ''), COALESCE(er.skill_name, ''), COALESCE(er.skill_version, ''),
			er.language, er.sample_id,
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
			COALESCE(trace.path, ''), COALESCE(diff.path, ''), COALESCE(er.sandbox_fingerprint, ''),
				gc.latency_ms
		FROM evaluation_results er
		LEFT JOIN artifacts g ON g.artifact_id = er.generated_test_artifact_id
		LEFT JOIN artifacts src ON src.artifact_id = er.source_artifact_id
		LEFT JOIN artifacts trace ON trace.artifact_id = er.trace_artifact_id
		LEFT JOIN artifacts diff ON diff.artifact_id = er.workspace_diff_artifact_id
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
			&r.Model, &r.SubjectID, &r.SubjectKind, &r.AgentFramework, &r.AgentModel, &r.SkillName, &r.SkillVersion,
			&r.Language, &r.SampleID,
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
			&r.TracePath, &r.WorkspaceDiffPath, &r.SandboxFingerprint,
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

func (s *SQLiteStore) FindReusableGeneratedAsset(ctx context.Context, generationKey string) (ReusableGeneratedCase, bool, error) {
	generationKey = strings.TrimSpace(generationKey)
	if generationKey == "" {
		return ReusableGeneratedCase{}, false, nil
	}
	row := s.db.QueryRowContext(ctx, `
		SELECT gc.generated_case_id, gc.run_id, gc.model, gc.language, gc.sample_id,
		       COALESCE(test_art.path, ''), COALESCE(test_art.sha256, ''),
		       COALESCE(resp_art.path, ''), COALESCE(meta_art.path, ''),
		       COALESCE(trace_art.path, ''), COALESCE(diff_art.path, ''),
		       COALESCE(gc.sandbox_fingerprint, ''), COALESCE(gc.latency_ms, 0),
		       gc.prompt_tokens, gc.completion_tokens, gc.total_tokens,
		       COALESCE(gc.token_source, ''), gc.estimated_cost_usd, COALESCE(gc.cost_source, ''),
		       COALESCE(gc.generated_at_utc, ''), COALESCE(gc.generated_test_artifact_id, '')
		FROM generated_cases gc
		LEFT JOIN artifacts test_art ON test_art.artifact_id = gc.generated_test_artifact_id
		LEFT JOIN artifacts resp_art ON resp_art.artifact_id = gc.response_artifact_id
		LEFT JOIN artifacts meta_art ON meta_art.artifact_id = gc.metadata_artifact_id
		LEFT JOIN artifacts trace_art ON trace_art.artifact_id = gc.trace_artifact_id
		LEFT JOIN artifacts diff_art ON diff_art.artifact_id = gc.workspace_diff_artifact_id
		WHERE gc.success = 1
		  AND gc.generation_key = ?
		  AND COALESCE(test_art.deleted_at_utc, '') = ''
		ORDER BY gc.generated_at_utc DESC
		LIMIT 1`, generationKey)
	var out ReusableGeneratedCase
	var promptTokens, completionTokens, totalTokens sql.NullInt64
	var estimatedCost sql.NullFloat64
	if err := row.Scan(
		&out.GeneratedCaseID, &out.RunID, &out.Model, &out.Language, &out.SampleID,
		&out.GeneratedTestPath, &out.GeneratedTestSHA256,
		&out.ResponsePath, &out.MetadataPath,
		&out.TracePath, &out.WorkspaceDiffPath,
		&out.SandboxFingerprint, &out.LatencyMS,
		&promptTokens, &completionTokens, &totalTokens,
		&out.TokenSource, &estimatedCost, &out.CostSource,
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
	out.TracePath = resolveStoredPath(out.TracePath)
	out.WorkspaceDiffPath = resolveStoredPath(out.WorkspaceDiffPath)
	out.PromptTokens = nullableSQLInt(promptTokens)
	out.CompletionTokens = nullableSQLInt(completionTokens)
	out.TotalTokens = nullableSQLInt(totalTokens)
	out.EstimatedCostUSD = nullableSQLFloat(estimatedCost)
	if !fileExists(out.GeneratedTestPath) {
		return ReusableGeneratedCase{}, false, nil
	}
	if out.GeneratedTestSHA256 != "" {
		sha, _, err := fileSHA256(out.GeneratedTestPath)
		if err != nil || sha != out.GeneratedTestSHA256 {
			return ReusableGeneratedCase{}, false, nil
		}
	}
	return out, true, nil
}

func (s *SQLiteStore) FindReusableEvaluationAsset(ctx context.Context, evaluationKey string) (ReusableEvaluationResult, bool, error) {
	evaluationKey = strings.TrimSpace(evaluationKey)
	if evaluationKey == "" {
		return ReusableEvaluationResult{}, false, nil
	}
	row := s.db.QueryRowContext(ctx, `
		SELECT er.evaluation_result_id, er.run_id,
		       er.compile_pass, er.test_pass, er.test_pass_count, er.test_total_count, er.test_pass_rate,
		       er.line_coverage, er.branch_coverage,
		       er.mutation_score, er.mutation_total, er.mutation_killed, er.mutation_survived,
		       er.mutation_no_tests, er.mutation_timeouts, er.mutation_skipped, er.mutation_suspicious,
		       er.assertion_count, er.test_case_count, er.assertion_density, er.runtime_ms,
		       er.prompt_tokens, er.completion_tokens, er.total_tokens,
		       COALESCE(er.token_source, ''), er.estimated_cost_usd, COALESCE(er.cost_source, ''),
		       er.truncated, COALESCE(er.mutation_tool, ''), COALESCE(er.failure_origin, ''),
		       er.score_eligible, COALESCE(er.score_exclusion_reason, ''),
		       COALESCE(er.compile_error, ''), COALESCE(er.test_error, ''), COALESCE(er.coverage_error, ''), COALESCE(er.mutation_error, ''),
		       COALESCE(trace_art.path, ''), COALESCE(diff_art.path, ''), COALESCE(er.sandbox_fingerprint, ''), COALESCE(er.evaluation_env_fingerprint, ''),
		       COALESCE(er.evaluation_key, ''), COALESCE(er.evaluator_version, ''), COALESCE(er.mutation_config_sha256, ''),
		       COALESCE(er.updated_db_at_utc, '')
		FROM evaluation_results er
		LEFT JOIN artifacts trace_art ON trace_art.artifact_id = er.trace_artifact_id
		LEFT JOIN artifacts diff_art ON diff_art.artifact_id = er.workspace_diff_artifact_id
		WHERE er.evaluation_key = ?
		  AND er.compile_pass = 1
		  AND COALESCE(er.test_pass, 1) = 1
		  AND er.score_eligible = 1
		ORDER BY er.updated_db_at_utc DESC
		LIMIT 1`, evaluationKey)

	var out ReusableEvaluationResult
	var compilePass, truncated int
	var testPass sql.NullInt64
	var testPassCount, testTotalCount, mutationTotal, mutationKilled, mutationSurvived sql.NullInt64
	var mutationNoTests, mutationTimeouts, mutationSkipped, mutationSuspicious sql.NullInt64
	var assertionCount, testCaseCount, runtimeMS, promptTokens, completionTokens, totalTokens sql.NullInt64
	var testPassRate, lineCoverage, branchCoverage, mutationScore, assertionDensity, estimatedCost sql.NullFloat64
	var scoreEligible sql.NullInt64
	if err := row.Scan(
		&out.EvaluationResultID, &out.RunID,
		&compilePass, &testPass, &testPassCount, &testTotalCount, &testPassRate,
		&lineCoverage, &branchCoverage,
		&mutationScore, &mutationTotal, &mutationKilled, &mutationSurvived,
		&mutationNoTests, &mutationTimeouts, &mutationSkipped, &mutationSuspicious,
		&assertionCount, &testCaseCount, &assertionDensity, &runtimeMS,
		&promptTokens, &completionTokens, &totalTokens,
		&out.Result.TokenSource, &estimatedCost, &out.Result.CostSource,
		&truncated, &out.Result.MutationTool, &out.Result.FailureOrigin,
		&scoreEligible, &out.Result.ScoreExclusionReason,
		&out.Result.CompileError, &out.Result.TestError, &out.Result.CoverageError, &out.Result.MutationError,
		&out.Result.TracePath, &out.Result.WorkspaceDiffPath, &out.Result.SandboxFingerprint, &out.Result.EvaluationEnvFingerprint,
		&out.Result.EvaluationKey, &out.Result.EvaluatorVersion, &out.Result.MutationConfigSHA256,
		&out.UpdatedAtUTC,
	); err != nil {
		if err == sql.ErrNoRows {
			return ReusableEvaluationResult{}, false, nil
		}
		return ReusableEvaluationResult{}, false, err
	}
	out.Result.CompilePass = compilePass != 0
	out.Result.TestPass = nullableIntBool(testPass)
	out.Result.TestPassCount = nullableSQLInt(testPassCount)
	out.Result.TestTotalCount = nullableSQLInt(testTotalCount)
	out.Result.TestPassRate = nullableSQLFloat(testPassRate)
	out.Result.LineCoverage = nullableSQLFloat(lineCoverage)
	out.Result.BranchCoverage = nullableSQLFloat(branchCoverage)
	out.Result.MutationScore = nullableSQLFloat(mutationScore)
	out.Result.MutationTotal = nullableSQLInt(mutationTotal)
	out.Result.MutationKilled = nullableSQLInt(mutationKilled)
	out.Result.MutationSurvived = nullableSQLInt(mutationSurvived)
	out.Result.MutationNoTests = nullableSQLInt(mutationNoTests)
	out.Result.MutationTimeouts = nullableSQLInt(mutationTimeouts)
	out.Result.MutationSkipped = nullableSQLInt(mutationSkipped)
	out.Result.MutationSuspicious = nullableSQLInt(mutationSuspicious)
	out.Result.AssertionCount = nullableSQLInt(assertionCount)
	out.Result.TestCaseCount = nullableSQLInt(testCaseCount)
	out.Result.AssertionDensity = nullableSQLFloat(assertionDensity)
	out.Result.RuntimeMS = nullableSQLInt(runtimeMS)
	out.Result.PromptTokens = nullableSQLInt(promptTokens)
	out.Result.CompletionTokens = nullableSQLInt(completionTokens)
	out.Result.TotalTokens = nullableSQLInt(totalTokens)
	out.Result.EstimatedCostUSD = nullableSQLFloat(estimatedCost)
	out.Result.Truncated = truncated != 0
	out.Result.ScoreEligible = nullableIntBool(scoreEligible)
	out.Result.TracePath = resolveStoredPath(out.Result.TracePath)
	out.Result.WorkspaceDiffPath = resolveStoredPath(out.Result.WorkspaceDiffPath)
	return out, true, nil
}
