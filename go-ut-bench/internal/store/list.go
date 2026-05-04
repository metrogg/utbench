// store 包提供 SQLite 数据库存储功能
// 本文件包含工具函数和数据库管理 List 方法
package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ────────────────────────────────────────────────────────────────
// 工具函数
// ────────────────────────────────────────────────────────────────

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

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
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

// ────────────────────────────────────────────────────────────────
// 数据库管理 List 方法
// ────────────────────────────────────────────────────────────────

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

func (s *SQLiteStore) ListAssetSubjects(ctx context.Context, limit int) ([]DBSubjectAssetItem, error) {
	if limit <= 0 || limit > 2000 {
		limit = 200
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT COALESCE(s.subject_id, gc.subject_id, gc.model) AS subject_id,
		       COALESCE(s.subject_kind, gc.subject_kind, ''),
		       COALESCE(s.framework, gc.agent_framework, ''),
		       COALESCE(s.model, gc.agent_model, ''),
		       COALESCE(s.skill, gc.skill_name, ''),
		       COUNT(DISTINCT gc.generated_case_id) AS generated_cases,
		       COUNT(DISTINCT er.evaluation_result_id) AS evaluation_results,
		       COALESCE(MAX(gc.generated_at_utc), '') AS latest_generated_at
		FROM generated_cases gc
		LEFT JOIN subjects s ON s.subject_id = gc.subject_id
		LEFT JOIN evaluation_results er ON er.generated_case_id = gc.generated_case_id
		GROUP BY subject_id
		ORDER BY latest_generated_at DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DBSubjectAssetItem
	for rows.Next() {
		var row DBSubjectAssetItem
		if err := rows.Scan(&row.SubjectID, &row.SubjectKind, &row.Framework, &row.Model, &row.Skill, &row.GeneratedCases, &row.EvaluationResults, &row.LatestGeneratedAt); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ListSubjectVersions(ctx context.Context, subjectID string, limit int) ([]DBSubjectVersionItem, error) {
	if limit <= 0 || limit > 2000 {
		limit = 200
	}
	where := []string{"1=1"}
	args := []any{}
	if subjectID != "" {
		where = append(where, "sv.subject_id = ?")
		args = append(args, subjectID)
	}
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx, `
		SELECT sv.subject_version_id, sv.subject_id, COALESCE(sv.model_config_id, ''),
		       COALESCE(sv.framework_config_sha256, ''), COALESCE(sv.skill_sha256, ''),
		       COALESCE(sv.agent_command_sha256, ''), COALESCE(sv.docker_image, ''),
		       COALESCE(sv.docker_image_digest, ''), COALESCE(sv.sandbox_fingerprint, ''),
		       COALESCE(sv.env_contract_sha256, ''), COALESCE(sv.created_at_utc, '')
		FROM subject_versions sv
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY sv.created_at_utc DESC
		LIMIT ?`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DBSubjectVersionItem
	for rows.Next() {
		var row DBSubjectVersionItem
		if err := rows.Scan(
			&row.SubjectVersionID, &row.SubjectID, &row.ModelConfigID,
			&row.FrameworkConfigSHA256, &row.SkillSHA256, &row.AgentCommandSHA256,
			&row.DockerImage, &row.DockerImageDigest, &row.SandboxFingerprint,
			&row.EnvContractSHA256, &row.CreatedAtUTC,
		); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ListAssetGenerations(ctx context.Context, subjectID, language, sampleID string, limit int) ([]DBGenerationAssetItem, error) {
	if limit <= 0 || limit > 5000 {
		limit = 200
	}
	where := []string{"1=1"}
	args := []any{}
	if subjectID != "" {
		where = append(where, "COALESCE(subject_id, model) = ?")
		args = append(args, subjectID)
	}
	if language != "" {
		where = append(where, "language = ?")
		args = append(args, language)
	}
	if sampleID != "" {
		where = append(where, "sample_id = ?")
		args = append(args, sampleID)
	}
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx, `
		SELECT gc.generated_case_id, gc.run_id, COALESCE(gc.subject_id, gc.model), gc.language, gc.sample_id,
		       COALESCE(gc.subject_version_id, ''), COALESCE(gc.sample_uid, ''),
		       COALESCE(gc.generation_key, ''), COALESCE(gc.dependency_fingerprint, ''),
		       COALESCE(gc.generation_env_fingerprint, ''), COALESCE(gc.sandbox_fingerprint, ''),
		       gc.success, gc.reused, COALESCE(gc.generated_at_utc, ''),
		       COALESCE(test_art.path, '')
		FROM generated_cases gc
		LEFT JOIN artifacts test_art ON test_art.artifact_id = gc.generated_test_artifact_id
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY gc.generated_at_utc DESC
		LIMIT ?`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DBGenerationAssetItem
	for rows.Next() {
		var row DBGenerationAssetItem
		var success, reused int
		if err := rows.Scan(
			&row.GeneratedCaseID, &row.RunID, &row.SubjectID, &row.Language, &row.SampleID,
			&row.SubjectVersionID, &row.SampleUID, &row.GenerationKey, &row.DependencyFingerprint,
			&row.GenerationEnvFingerprint, &row.SandboxFingerprint,
			&success, &reused, &row.GeneratedAtUTC, &row.GeneratedTestPath,
		); err != nil {
			return nil, err
		}
		row.Success = success != 0
		row.Reused = reused != 0
		row.GeneratedTestPath = resolveStoredPath(row.GeneratedTestPath)
		out = append(out, row)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ListAssetEvaluations(ctx context.Context, subjectID, language, sampleID string, limit int) ([]DBEvaluationAssetItem, error) {
	if limit <= 0 || limit > 5000 {
		limit = 200
	}
	where := []string{"1=1"}
	args := []any{}
	if subjectID != "" {
		where = append(where, "COALESCE(subject_id, model) = ?")
		args = append(args, subjectID)
	}
	if language != "" {
		where = append(where, "language = ?")
		args = append(args, language)
	}
	if sampleID != "" {
		where = append(where, "sample_id = ?")
		args = append(args, sampleID)
	}
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx, `
		SELECT evaluation_result_id, run_id, COALESCE(subject_id, model), language, sample_id,
		       COALESCE(evaluation_key, ''), COALESCE(evaluation_env_fingerprint, ''), compile_pass, reused, COALESCE(updated_db_at_utc, '')
		FROM evaluation_results
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY updated_db_at_utc DESC
		LIMIT ?`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DBEvaluationAssetItem
	for rows.Next() {
		var row DBEvaluationAssetItem
		var compilePass, reused int
		if err := rows.Scan(&row.EvaluationResultID, &row.RunID, &row.SubjectID, &row.Language, &row.SampleID, &row.EvaluationKey, &row.EvaluationEnvFingerprint, &compilePass, &reused, &row.CreatedAtUTC); err != nil {
			return nil, err
		}
		row.CompilePass = compilePass != 0
		row.Reused = reused != 0
		out = append(out, row)
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
