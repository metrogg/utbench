// store 包提供 SQLite 数据库存储功能
// 本文件包含数据入库(ingest)相关的所有函数
package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"go-ut-bench/internal/contracts"

	"gopkg.in/yaml.v3"
)

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

// IngestRun 入库整个运行目录
// 自动检测并入库 manifest、evaluation、report 文件
//
// 参数:
//   - ctx: 上下文
//   - opts: 入库选项，包含 RunDir 路径
//
// 返回值:
//   - IngestSummary: 入库统计摘要
//   - error: 入库过程中的错误
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
		seenModels[firstNonEmpty(c.AgentModel, c.Model)] = struct{}{}
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
		if c.SampleUID != "" {
			sampleUID = c.SampleUID
		}
		subjectID := firstNonEmpty(c.SubjectID, c.Model)
		subjectModel := firstNonEmpty(c.AgentModel, c.Model)
		cfg := modelConfigs[subjectModel]
		modelConfigID := stableID("model_config", subjectModel, cfg.Provider, cfg.ModelID)
		subjectVersionID := firstNonEmpty(c.SubjectVersionID, stableID("subject_version", subjectID, modelConfigID, c.SandboxFingerprint, c.SkillVersion))
		if err := s.upsertSubjectTx(ctx, tx, subjectID, c.SubjectKind, c.AgentFramework, subjectModel, c.SkillName, now); err != nil {
			return err
		}
		if err := s.upsertSubjectVersionTx(ctx, tx, DBSubjectVersionItem{
			SubjectVersionID:      subjectVersionID,
			SubjectID:             subjectID,
			ModelConfigID:         modelConfigID,
			FrameworkConfigSHA256: c.FrameworkConfigSHA256,
			SkillSHA256:           c.SkillSHA256,
			AgentCommandSHA256:    c.AgentCommandSHA256,
			DockerImage:           c.DockerImage,
			DockerImageDigest:     c.DockerImageDigest,
			SandboxFingerprint:    c.SandboxFingerprint,
			EnvContractSHA256:     c.EnvContractSHA256,
			CreatedAtUTC:          now,
		}); err != nil {
			return err
		}
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
		traceArtifactID, _ := s.putOptionalArtifact(ctx, tx, ictx, "agent_trace", c.TracePath, "trace")
		workspaceDiffArtifactID, _ := s.putOptionalArtifact(ctx, tx, ictx, "agent_workspace_diff", c.WorkspaceDiffPath, "workspace_diff")
		caseID := stableID("generated_case", manifest.RunID, c.Model, c.Language, c.SampleID)
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO generated_cases(generated_case_id, run_id, model, subject_id, subject_kind, agent_framework, agent_model, skill_name, skill_version,
				language, sample_id, sample_uid, sample_path,
				prompt_rendering_id, generated_test_artifact_id, response_artifact_id, metadata_artifact_id,
				trace_artifact_id, workspace_diff_artifact_id, sandbox_fingerprint,
				subject_version_id, framework_config_sha256, skill_sha256, agent_command_sha256, docker_image, docker_image_digest, env_contract_sha256,
				generation_key, dependency_fingerprint, generation_env_fingerprint,
				reused, reuse_stage, reuse_key, reuse_reason, reused_from_run_id, reused_from_case_id,
				latency_ms, prompt_tokens, completion_tokens, total_tokens, token_source, estimated_cost_usd, cost_source,
				generated_at_utc, success, truncated, error_json, created_db_at_utc, updated_db_at_utc)
			VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(run_id, model, language, sample_id) DO UPDATE SET
				subject_id=excluded.subject_id,
				subject_kind=excluded.subject_kind,
				agent_framework=excluded.agent_framework,
				agent_model=excluded.agent_model,
				skill_name=excluded.skill_name,
				skill_version=excluded.skill_version,
				sample_uid=excluded.sample_uid,
				sample_path=excluded.sample_path,
				prompt_rendering_id=excluded.prompt_rendering_id,
				generated_test_artifact_id=excluded.generated_test_artifact_id,
				response_artifact_id=excluded.response_artifact_id,
				metadata_artifact_id=excluded.metadata_artifact_id,
				trace_artifact_id=excluded.trace_artifact_id,
				workspace_diff_artifact_id=excluded.workspace_diff_artifact_id,
				sandbox_fingerprint=excluded.sandbox_fingerprint,
				subject_version_id=excluded.subject_version_id,
				framework_config_sha256=excluded.framework_config_sha256,
				skill_sha256=excluded.skill_sha256,
				agent_command_sha256=excluded.agent_command_sha256,
				docker_image=excluded.docker_image,
				docker_image_digest=excluded.docker_image_digest,
				env_contract_sha256=excluded.env_contract_sha256,
				generation_key=excluded.generation_key,
				dependency_fingerprint=excluded.dependency_fingerprint,
				generation_env_fingerprint=excluded.generation_env_fingerprint,
				reused=excluded.reused,
				reuse_stage=excluded.reuse_stage,
				reuse_key=excluded.reuse_key,
				reuse_reason=excluded.reuse_reason,
				reused_from_run_id=excluded.reused_from_run_id,
				reused_from_case_id=excluded.reused_from_case_id,
				latency_ms=excluded.latency_ms,
				prompt_tokens=excluded.prompt_tokens,
				completion_tokens=excluded.completion_tokens,
				total_tokens=excluded.total_tokens,
				token_source=excluded.token_source,
				estimated_cost_usd=excluded.estimated_cost_usd,
				cost_source=excluded.cost_source,
				generated_at_utc=excluded.generated_at_utc,
				success=excluded.success,
				truncated=excluded.truncated,
				error_json=excluded.error_json,
				updated_db_at_utc=excluded.updated_db_at_utc`,
			caseID, manifest.RunID, c.Model, nullString(subjectID), nullString(c.SubjectKind), nullString(c.AgentFramework), nullString(c.AgentModel), nullString(c.SkillName), nullString(c.SkillVersion),
			c.Language, c.SampleID, nullString(sampleUID), c.SamplePath,
			promptRenderingID, nullString(generatedTestArtifactID), nullString(responseArtifactID), nullString(metadataArtifactID),
			nullString(traceArtifactID), nullString(workspaceDiffArtifactID), nullString(c.SandboxFingerprint),
			nullString(subjectVersionID), nullString(c.FrameworkConfigSHA256), nullString(c.SkillSHA256), nullString(c.AgentCommandSHA256), nullString(c.DockerImage), nullString(c.DockerImageDigest), nullString(c.EnvContractSHA256),
			nullString(c.GenerationKey), nullString(c.DependencyFingerprint), nullString(c.GenerationEnvFingerprint),
			boolToInt(c.Reused), nullString(c.ReuseStage), nullString(c.ReuseKey), nullString(c.ReuseReason), nullString(c.ReusedFromRunID), nullString(c.ReusedFromCaseID),
			c.LatencyMS, nullableInt(c.PromptTokens), nullableInt(c.CompletionTokens), nullableInt(c.TotalTokens), nullString(c.TokenSource), nullableFloat(c.EstimatedCostUSD), nullString(c.CostSource), formatTime(c.GeneratedAtUTC),
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
	policyJSON := fmt.Sprintf(`{"compile":%.2f,"test":%.2f,"coverage":%.2f,"mutation":%.2f,"excludes_score_eligible_false":true}`,
		contracts.DefaultWeights.Compile, contracts.DefaultWeights.Test, contracts.DefaultWeights.Coverage, contracts.DefaultWeights.Mutation)
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
		traceArtifactID, _ := s.putOptionalArtifact(ctx, tx, ictx, "agent_trace", row.TracePath, "trace")
		workspaceDiffArtifactID, _ := s.putOptionalArtifact(ctx, tx, ictx, "agent_workspace_diff", row.WorkspaceDiffPath, "workspace_diff")
		generatedCaseID := stableID("generated_case", set.RunID, row.Model, row.Language, row.SampleID)
		resultID := stableID("evaluation_result", evalRunID, row.Model, row.Language, row.SampleID)
		eligible := true
		if row.ScoreEligible != nil {
			eligible = *row.ScoreEligible
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO evaluation_results(evaluation_result_id, evaluation_run_id, run_id, model, subject_id, subject_kind, agent_framework, agent_model, skill_name, skill_version,
				language, sample_id, sample_uid,
				generated_case_id, generated_test_artifact_id, source_artifact_id,
				compile_pass, test_pass, test_pass_count, test_total_count, test_pass_rate,
				line_coverage, branch_coverage, mutation_score, mutation_total, mutation_killed, mutation_survived,
				mutation_no_tests, mutation_timeouts, mutation_skipped, mutation_suspicious,
				assertion_count, test_case_count, assertion_density, runtime_ms,
				prompt_tokens, completion_tokens, total_tokens, token_source, estimated_cost_usd, cost_source, truncated, mutation_tool,
				failure_origin, score_eligible, score_exclusion_reason,
				compile_error, test_error, coverage_error, mutation_error,
			trace_artifact_id, workspace_diff_artifact_id, sandbox_fingerprint,
			evaluation_env_fingerprint,
			evaluation_key, evaluator_version, mutation_config_sha256,
			reused, reuse_stage, reuse_key, reuse_reason, reused_from_run_id, reused_from_result_id,
			created_db_at_utc, updated_db_at_utc)
			VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(evaluation_run_id, model, language, sample_id) DO UPDATE SET
				subject_id=excluded.subject_id,
				subject_kind=excluded.subject_kind,
				agent_framework=excluded.agent_framework,
				agent_model=excluded.agent_model,
				skill_name=excluded.skill_name,
				skill_version=excluded.skill_version,
				sample_uid=excluded.sample_uid,
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
				token_source=excluded.token_source,
				estimated_cost_usd=excluded.estimated_cost_usd,
				cost_source=excluded.cost_source,
				truncated=excluded.truncated,
				mutation_tool=excluded.mutation_tool,
				failure_origin=excluded.failure_origin,
				score_eligible=excluded.score_eligible,
				score_exclusion_reason=excluded.score_exclusion_reason,
				compile_error=excluded.compile_error,
				test_error=excluded.test_error,
				coverage_error=excluded.coverage_error,
				mutation_error=excluded.mutation_error,
				trace_artifact_id=excluded.trace_artifact_id,
				workspace_diff_artifact_id=excluded.workspace_diff_artifact_id,
				sandbox_fingerprint=excluded.sandbox_fingerprint,
				evaluation_env_fingerprint=excluded.evaluation_env_fingerprint,
				evaluation_key=excluded.evaluation_key,
				evaluator_version=excluded.evaluator_version,
				mutation_config_sha256=excluded.mutation_config_sha256,
				reused=excluded.reused,
				reuse_stage=excluded.reuse_stage,
				reuse_key=excluded.reuse_key,
				reuse_reason=excluded.reuse_reason,
				reused_from_run_id=excluded.reused_from_run_id,
				reused_from_result_id=excluded.reused_from_result_id,
				updated_db_at_utc=excluded.updated_db_at_utc`,
			resultID, evalRunID, set.RunID, row.Model, nullString(firstNonEmpty(row.SubjectID, row.Model)), nullString(row.SubjectKind), nullString(row.AgentFramework), nullString(row.AgentModel), nullString(row.SkillName), nullString(row.SkillVersion),
			row.Language, row.SampleID, nullString(row.SampleUID),
			generatedCaseID, nullString(generatedTestArtifactID), nullString(sourceArtifactID),
			boolToInt(row.CompilePass), ptrBoolToNullableInt(row.TestPass), nullableInt(row.TestPassCount), nullableInt(row.TestTotalCount), nullableFloat(row.TestPassRate),
			nullableFloat(row.LineCoverage), nullableFloat(row.BranchCoverage), nullableFloat(row.MutationScore), nullableInt(row.MutationTotal), nullableInt(row.MutationKilled), nullableInt(row.MutationSurvived),
			nullableInt(row.MutationNoTests), nullableInt(row.MutationTimeouts), nullableInt(row.MutationSkipped), nullableInt(row.MutationSuspicious),
			nullableInt(row.AssertionCount), nullableInt(row.TestCaseCount), nullableFloat(row.AssertionDensity), nullableInt(row.RuntimeMS),
			nullableInt(row.PromptTokens), nullableInt(row.CompletionTokens), nullableInt(row.TotalTokens), nullString(row.TokenSource), nullableFloat(row.EstimatedCostUSD), nullString(row.CostSource), boolToInt(row.Truncated), nullString(row.MutationTool),
			nullString(row.FailureOrigin), boolToInt(eligible), nullString(row.ScoreExclusionReason),
			nullString(row.CompileError), nullString(row.TestError), nullString(row.CoverageError), nullString(row.MutationError),
			nullString(traceArtifactID), nullString(workspaceDiffArtifactID), nullString(row.SandboxFingerprint),
			nullString(row.EvaluationEnvFingerprint),
			nullString(row.EvaluationKey), nullString(row.EvaluatorVersion), nullString(row.MutationConfigSHA256),
			boolToInt(row.Reused), nullString(row.ReuseStage), nullString(row.ReuseKey), nullString(row.ReuseReason), nullString(row.ReusedFromRunID), nullString(row.ReusedFromResultID),
			now, now); err != nil {
			return fmt.Errorf("upsert evaluation result: %w", err)
		}
	}
	return nil
}

func (s *SQLiteStore) upsertSubjectTx(ctx context.Context, tx *sql.Tx, subjectID, subjectKind, framework, model, skill, now string) error {
	if strings.TrimSpace(subjectID) == "" {
		return nil
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO subjects(subject_id, subject_kind, framework, model, skill, display_name, enabled, tags_json, created_at_utc, updated_at_utc)
		VALUES(?, ?, ?, ?, ?, ?, 1, ?, ?, ?)
		ON CONFLICT(subject_id) DO UPDATE SET
			subject_kind=excluded.subject_kind,
			framework=excluded.framework,
			model=excluded.model,
			skill=excluded.skill,
			display_name=excluded.display_name,
			updated_at_utc=excluded.updated_at_utc`,
		subjectID, nullString(subjectKind), nullString(framework), nullString(model), nullString(skill), subjectID, "[]", now, now); err != nil {
		return fmt.Errorf("upsert subject: %w", err)
	}
	return nil
}

func (s *SQLiteStore) upsertSubjectVersionTx(ctx context.Context, tx *sql.Tx, row DBSubjectVersionItem) error {
	if strings.TrimSpace(row.SubjectVersionID) == "" || strings.TrimSpace(row.SubjectID) == "" {
		return nil
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO subject_versions(subject_version_id, subject_id, model_config_id,
			framework_config_sha256, skill_sha256, agent_command_sha256, docker_image, docker_image_digest,
			sandbox_fingerprint, env_contract_sha256, created_at_utc)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(subject_version_id) DO UPDATE SET
			subject_id=excluded.subject_id,
			model_config_id=excluded.model_config_id,
			framework_config_sha256=excluded.framework_config_sha256,
			skill_sha256=excluded.skill_sha256,
			agent_command_sha256=excluded.agent_command_sha256,
			docker_image=excluded.docker_image,
			docker_image_digest=excluded.docker_image_digest,
			sandbox_fingerprint=excluded.sandbox_fingerprint,
			env_contract_sha256=excluded.env_contract_sha256`,
		row.SubjectVersionID, row.SubjectID, nullString(row.ModelConfigID),
		nullString(row.FrameworkConfigSHA256), nullString(row.SkillSHA256), nullString(row.AgentCommandSHA256), nullString(row.DockerImage), nullString(row.DockerImageDigest),
		nullString(row.SandboxFingerprint), nullString(row.EnvContractSHA256), firstNonEmpty(row.CreatedAtUTC, nowUTC())); err != nil {
		return fmt.Errorf("upsert subject version: %w", err)
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
		uid := strings.TrimSpace(c.SampleUID)
		if uid == "" {
			uid = stableID("dataset_sample", c.Language, c.SampleID, c.SamplePath, sha)
		}
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
	case strings.Contains(p, "repo_level"):
		class = "repo_level"
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
