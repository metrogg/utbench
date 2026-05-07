// store 包提供 SQLite 数据库存储功能
// 负责持久化评测运行数据、索引 artifacts、提供查询接口
// 支持生成结果复用、历史查询、跨运行对比报告
package store

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	"go-ut-bench/internal/contracts"

	_ "modernc.org/sqlite"
)

// schemaVersion 数据库 Schema 版本
const schemaVersion = "store.v2"

// SQLiteStore SQLite 存储实例
// 持有数据库连接和路径信息
type SQLiteStore struct {
	db   *sql.DB // 数据库连接
	path string  // 数据库文件路径
}

// DBOverview 数据库概览统计
// 包含各表记录数和最新运行列表
type DBOverview struct {
	DBPath            string      `json:"db_path"`            // 数据库文件路径
	SchemaVersion     string      `json:"schema_version"`     // Schema 版本
	GenerationRuns    int         `json:"generation_runs"`    // 生成运行数
	EvaluationRuns    int         `json:"evaluation_runs"`    // 评测运行数
	GeneratedCases    int         `json:"generated_cases"`    // 生成案例数
	EvaluationResults int         `json:"evaluation_results"` // 评测结果数
	Artifacts         int         `json:"artifacts"`          // artifacts 数
	Reports           int         `json:"reports"`            // 报告数
	LatestRuns        []DBRunItem `json:"latest_runs"`        // 最新运行列表
}

// DBRunItem 运行记录摘要
// 用于概览中展示运行基本信息
type DBRunItem struct {
	RunID             string `json:"run_id"`                     // 运行ID
	ExperimentID      string `json:"experiment_id,omitempty"`    // 实验ID
	ExperimentName    string `json:"experiment_name,omitempty"`  // 实验名称
	CreatedAtUTC      string `json:"created_at_utc,omitempty"`   // 创建时间
	EvaluatedAtUTC    string `json:"evaluated_at_utc,omitempty"` // 评测时间
	Models            int    `json:"models"`                     // 模型数
	Languages         int    `json:"languages"`                  // 语言数
	GeneratedCases    int    `json:"generated_cases"`            // 生成案例数
	EvaluationResults int    `json:"evaluation_results"`         // 评测结果数
	Reports           int    `json:"reports"`                    // 报告数
}

// DBResultItem 评测结果记录
// 用于查询结果列表展示
type DBResultItem struct {
	EvaluationResultID   string   `json:"evaluation_result_id"`             // 评测结果ID
	EvaluationRunID      string   `json:"evaluation_run_id"`                // 评测运行ID
	RunID                string   `json:"run_id"`                           // 运行ID
	Model                string   `json:"model"`                            // 模型名
	Language             string   `json:"language"`                         // 语言
	SampleID             string   `json:"sample_id"`                        // 样本ID
	CompilePass          bool     `json:"compile_pass"`                     // 编译是否通过
	TestPass             *bool    `json:"test_pass,omitempty"`              // 测试是否通过
	LineCoverage         *float64 `json:"line_coverage,omitempty"`          // 行覆盖率
	MutationScore        *float64 `json:"mutation_score,omitempty"`         // 变异得分
	MutationTotal        *int     `json:"mutation_total,omitempty"`         // 变异体总数
	FailureOrigin        string   `json:"failure_origin,omitempty"`         // 失败归因
	ScoreEligible        bool     `json:"score_eligible"`                   // 是否计入排名
	ScoreExclusionReason string   `json:"score_exclusion_reason,omitempty"` // 排名剔除原因
	RuntimeMS            *int     `json:"runtime_ms,omitempty"`             // 运行耗时（毫秒）
}

// DBArtifactItem artifact 记录
// 用于索引文件产物信息
type DBArtifactItem struct {
	ArtifactID   string `json:"artifact_id"`    // artifact ID
	Kind         string `json:"kind"`           // 类型（manifest/evaluation/report等）
	Path         string `json:"path"`           // 文件路径
	SizeBytes    int64  `json:"size_bytes"`     // 文件大小
	SHA256       string `json:"sha256"`         // SHA256 哈希
	Redacted     bool   `json:"redacted"`       // 是否已脱敏
	CreatedAtUTC string `json:"created_at_utc"` // 创建时间
	DeletedAtUTC string `json:"deleted_at_utc,omitempty"`
}

type DBReportFilter struct {
	RunIDs            []string `json:"run_ids,omitempty"`
	EvaluationRunIDs  []string `json:"evaluation_run_ids,omitempty"`
	Models            []string `json:"models,omitempty"`
	Languages         []string `json:"languages,omitempty"`
	ScoreEligibleOnly bool     `json:"score_eligible_only,omitempty"`
}

type DBReportFacets struct {
	Runs      []DBReportRunOption `json:"runs"`
	Models    []string            `json:"models"`
	Languages []string            `json:"languages"`
	EnvGroups []DBReportEnvOption `json:"env_groups"`
}

type DBReportRunOption struct {
	RunID             string `json:"run_id"`
	EvaluationRunID   string `json:"evaluation_run_id"`
	EvaluatedAtUTC    string `json:"evaluated_at_utc"`
	EnvID             string `json:"env_id"`
	PromptStrategy    string `json:"prompt_strategy,omitempty"`
	PromptVersionID   string `json:"prompt_version_id,omitempty"`
	EvaluationResults int    `json:"evaluation_results"`
}

type DBReportEnvOption struct {
	EnvID             string `json:"env_id"`
	Fingerprint       string `json:"fingerprint"`
	EvaluationRuns    int    `json:"evaluation_runs"`
	EvaluationResults int    `json:"evaluation_results"`
}

// Phase 3: Database management list item types

type DBGenerationRunItem struct {
	RunID              string `json:"run_id"`
	ExperimentID       string `json:"experiment_id,omitempty"`
	SchemaVersion      string `json:"schema_version"`
	CreatedAtUTC       string `json:"created_at_utc"`
	PromptStrategy     string `json:"prompt_strategy,omitempty"`
	PromptVersionID    string `json:"prompt_version_id,omitempty"`
	DatasetFingerprint string `json:"dataset_fingerprint,omitempty"`
	CreatedDBAtUTC     string `json:"created_db_at_utc"`
}

type DBGeneratedCaseItem struct {
	GeneratedCaseID  string `json:"generated_case_id"`
	RunID            string `json:"run_id"`
	Model            string `json:"model"`
	Language         string `json:"language"`
	SampleID         string `json:"sample_id"`
	Success          bool   `json:"success"`
	Truncated        bool   `json:"truncated"`
	LatencyMS        *int   `json:"latency_ms,omitempty"`
	PromptTokens     *int   `json:"prompt_tokens,omitempty"`
	CompletionTokens *int   `json:"completion_tokens,omitempty"`
	GeneratedAtUTC   string `json:"generated_at_utc,omitempty"`
}

type DBPromptRenderingItem struct {
	PromptRenderingID string `json:"prompt_rendering_id"`
	RunID             string `json:"run_id"`
	Model             string `json:"model"`
	Language          string `json:"language"`
	SampleID          string `json:"sample_id"`
	PromptVersionID   string `json:"prompt_version_id,omitempty"`
	PromptMode        string `json:"prompt_mode,omitempty"`
	CreatedAtUTC      string `json:"created_at_utc"`
}

type DBEvaluationRunItem struct {
	EvaluationRunID string `json:"evaluation_run_id"`
	RunID           string `json:"run_id"`
	SchemaVersion   string `json:"schema_version"`
	EvaluatedAtUTC  string `json:"evaluated_at_utc"`
	EnvID           string `json:"env_id,omitempty"`
	ScorePolicyID   string `json:"score_policy_id,omitempty"`
	CreatedDBAtUTC  string `json:"created_db_at_utc"`
}

type DBEvaluationStageItem struct {
	StageResultID      string `json:"stage_result_id"`
	EvaluationResultID string `json:"evaluation_result_id"`
	Stage              string `json:"stage"`
	Status             string `json:"status"`
	ExitCode           *int   `json:"exit_code,omitempty"`
	DurationMS         *int   `json:"duration_ms,omitempty"`
	CreatedAtUTC       string `json:"created_at_utc"`
}

type DBDatasetSampleItem struct {
	SampleUID    string `json:"sample_uid"`
	SampleID     string `json:"sample_id"`
	Language     string `json:"language"`
	Class        string `json:"class,omitempty"`
	Scenario     string `json:"scenario,omitempty"`
	Path         string `json:"path"`
	CreatedAtUTC string `json:"created_at_utc"`
}

type DBDatasetSnapshotItem struct {
	SnapshotID   string `json:"snapshot_id"`
	Fingerprint  string `json:"fingerprint"`
	SampleCount  int    `json:"sample_count"`
	CreatedAtUTC string `json:"created_at_utc"`
}

type DBModelConfigItem struct {
	ModelConfigID string `json:"model_config_id"`
	ModelName     string `json:"model_name"`
	Provider      string `json:"provider,omitempty"`
	ModelID       string `json:"model_id,omitempty"`
	CreatedAtUTC  string `json:"created_at_utc"`
}

type DBPromptProfileItem struct {
	ProfileID    string `json:"profile_id"`
	Strategy     string `json:"strategy,omitempty"`
	VersionID    string `json:"version_id,omitempty"`
	CreatedAtUTC string `json:"created_at_utc"`
}

type DBEvaluationEnvItem struct {
	EnvID        string `json:"env_id"`
	Fingerprint  string `json:"fingerprint"`
	CreatedAtUTC string `json:"created_at_utc"`
}

type DBSubjectVersionItem struct {
	SubjectVersionID      string `json:"subject_version_id"`
	SubjectID             string `json:"subject_id"`
	ModelConfigID         string `json:"model_config_id,omitempty"`
	FrameworkConfigSHA256 string `json:"framework_config_sha256,omitempty"`
	SkillSHA256           string `json:"skill_sha256,omitempty"`
	AgentCommandSHA256    string `json:"agent_command_sha256,omitempty"`
	DockerImage           string `json:"docker_image,omitempty"`
	DockerImageDigest     string `json:"docker_image_digest,omitempty"`
	SandboxFingerprint    string `json:"sandbox_fingerprint,omitempty"`
	EnvContractSHA256     string `json:"env_contract_sha256,omitempty"`
	CreatedAtUTC          string `json:"created_at_utc"`
}

type DBScorePolicyItem struct {
	ScorePolicyID string `json:"score_policy_id"`
	Name          string `json:"name"`
	CreatedAtUTC  string `json:"created_at_utc"`
}

type DBReportItem struct {
	ReportID       string `json:"report_id"`
	RunID          string `json:"run_id"`
	GeneratedAtUTC string `json:"generated_at_utc"`
	CreatedDBAtUTC string `json:"created_db_at_utc"`
}

type DBRunArtifactItem struct {
	RunID        string `json:"run_id"`
	ArtifactID   string `json:"artifact_id"`
	Role         string `json:"role"`
	CreatedAtUTC string `json:"created_at_utc"`
}

type DBExperimentItem struct {
	ExperimentID string `json:"experiment_id"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	CreatedAtUTC string `json:"created_at_utc"`
	UpdatedAtUTC string `json:"updated_at_utc"`
}

type DBSubjectAssetItem struct {
	SubjectID         string `json:"subject_id"`
	SubjectKind       string `json:"subject_kind,omitempty"`
	Framework         string `json:"framework,omitempty"`
	Model             string `json:"model,omitempty"`
	Skill             string `json:"skill,omitempty"`
	GeneratedCases    int    `json:"generated_cases"`
	EvaluationResults int    `json:"evaluation_results"`
	LatestGeneratedAt string `json:"latest_generated_at,omitempty"`
}

type DBGenerationAssetItem struct {
	GeneratedCaseID          string `json:"generated_case_id"`
	RunID                    string `json:"run_id"`
	SubjectID                string `json:"subject_id,omitempty"`
	SubjectVersionID         string `json:"subject_version_id,omitempty"`
	Language                 string `json:"language"`
	SampleID                 string `json:"sample_id"`
	SampleUID                string `json:"sample_uid,omitempty"`
	GenerationKey            string `json:"generation_key,omitempty"`
	DependencyFingerprint    string `json:"dependency_fingerprint,omitempty"`
	GenerationEnvFingerprint string `json:"generation_env_fingerprint,omitempty"`
	SandboxFingerprint       string `json:"sandbox_fingerprint,omitempty"`
	Success                  bool   `json:"success"`
	Reused                   bool   `json:"reused"`
	GeneratedAtUTC           string `json:"generated_at_utc,omitempty"`
	GeneratedTestPath        string `json:"generated_test_path,omitempty"`
}

type DBEvaluationAssetItem struct {
	EvaluationResultID       string `json:"evaluation_result_id"`
	RunID                    string `json:"run_id"`
	SubjectID                string `json:"subject_id,omitempty"`
	Language                 string `json:"language"`
	SampleID                 string `json:"sample_id"`
	EvaluationKey            string `json:"evaluation_key,omitempty"`
	EvaluationEnvFingerprint string `json:"evaluation_env_fingerprint,omitempty"`
	CompilePass              bool   `json:"compile_pass"`
	Reused                   bool   `json:"reused"`
	CreatedAtUTC             string `json:"created_at_utc,omitempty"`
}

// ReusableGeneratedCase 可复用的生成结果（别名，实际定义在 contracts 包）
type ReusableGeneratedCase = contracts.ReusableGeneratedCase

type ReusableEvaluationResult struct {
	EvaluationResultID string                     `json:"evaluation_result_id"`
	RunID              string                     `json:"run_id"`
	Result             contracts.EvaluationResult `json:"result"`
	UpdatedAtUTC       string                     `json:"updated_at_utc,omitempty"`
}

type IngestRunOptions struct {
	RunDir string
}

type IngestSummary struct {
	RunID                string `json:"run_id"`
	ManifestIngested     bool   `json:"manifest_ingested"`
	EvaluationIngested   bool   `json:"evaluation_ingested"`
	ReportIngested       bool   `json:"report_ingested"`
	GenerationCases      int    `json:"generation_cases"`
	EvaluationResults    int    `json:"evaluation_results"`
	ArtifactsIndexed     int    `json:"artifacts_indexed"`
	UnavailableArtifacts int    `json:"unavailable_artifacts"`
}

type ingestContext struct {
	runID   string
	runDir  string
	spec    contracts.RunSpec
	seen    map[string]struct{}
	summary *IngestSummary
}

// OpenSQLite 打开 SQLite 数据库
// 如果目录不存在会自动创建，如果数据库文件不存在会在 Init 时创建
//
// 参数:
//   - path: 数据库文件路径
//
// 返回值:
//   - *SQLiteStore: 存储实例
//   - error: 打开过程中的错误
func OpenSQLite(path string) (*SQLiteStore, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	// 启用 WAL 与 busy_timeout，避免并发 generate/evaluate 写入时
	// 报 "database is locked (5) (SQLITE_BUSY)"。
	// - journal_mode=WAL：读写并发不再互斥
	// - busy_timeout=5000：写写仍串行，但驱动会重试 5s
	// - synchronous=NORMAL：WAL 下安全，且写入更快
	dsn := path + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=synchronous(NORMAL)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	// 写入串行化：SQLite 同一时刻只允许一个写事务，把池上限设为 1 可彻底
	// 消除多 goroutine 抢锁导致的 BUSY；读路径走 WAL 不受影响。
	db.SetMaxOpenConns(1)
	return &SQLiteStore{db: db, path: path}, nil
}

// Close 关闭数据库连接
func (s *SQLiteStore) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *SQLiteStore) ensureColumn(ctx context.Context, table, columnDef string) error {
	columnName := strings.Fields(columnDef)[0]
	rows, err := s.db.QueryContext(ctx, `PRAGMA table_info(`+table+`)`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var cid int
		var name, typ string
		var notNull int
		var defaultValue any
		var pk int
		if err := rows.Scan(&cid, &name, &typ, &notNull, &defaultValue, &pk); err != nil {
			return err
		}
		if name == columnName {
			return nil
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, `ALTER TABLE `+table+` ADD COLUMN `+columnDef)
	return err
}

// Init 初始化数据库表结构
// 创建所有必要的表（experiments、generation_runs、evaluation_runs、generated_cases、evaluation_results 等）
//
// 参数:
//   - ctx: 上下文
//
// 返回值:
//   - error: 初始化过程中的错误
func (s *SQLiteStore) Init(ctx context.Context) error {
	ddl := []string{
		`PRAGMA foreign_keys = ON;`,
		`CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS experiments (
			experiment_id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			created_at_utc TEXT NOT NULL,
			updated_at_utc TEXT NOT NULL,
			deleted_at_utc TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS artifacts (
			artifact_id TEXT PRIMARY KEY,
			kind TEXT NOT NULL,
			path TEXT NOT NULL,
			size_bytes INTEGER NOT NULL,
			sha256 TEXT NOT NULL,
			redacted INTEGER NOT NULL DEFAULT 0,
			created_at_utc TEXT NOT NULL,
			deleted_at_utc TEXT,
			UNIQUE(kind, sha256)
		);`,
		`CREATE TABLE IF NOT EXISTS run_artifacts (
			run_id TEXT NOT NULL,
			artifact_id TEXT NOT NULL,
			role TEXT NOT NULL,
			created_at_utc TEXT NOT NULL,
			PRIMARY KEY(run_id, artifact_id, role)
		);`,
		`CREATE TABLE IF NOT EXISTS dataset_samples (
			sample_uid TEXT PRIMARY KEY,
			sample_id TEXT NOT NULL,
			language TEXT NOT NULL,
			class TEXT,
			scenario TEXT,
			path TEXT NOT NULL,
			source_md5 TEXT,
			source_sha256 TEXT,
			source_artifact_id TEXT,
			risk_flags_json TEXT,
			created_at_utc TEXT NOT NULL,
			updated_at_utc TEXT NOT NULL,
			UNIQUE(language, sample_id, path)
		);`,
		`CREATE TABLE IF NOT EXISTS dataset_snapshots (
			snapshot_id TEXT PRIMARY KEY,
			fingerprint TEXT NOT NULL UNIQUE,
			created_at_utc TEXT NOT NULL,
			sample_count INTEGER NOT NULL,
			spec_json TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS dataset_snapshot_members (
			snapshot_id TEXT NOT NULL,
			sample_uid TEXT NOT NULL,
			PRIMARY KEY(snapshot_id, sample_uid)
		);`,
		`CREATE TABLE IF NOT EXISTS model_configs (
			model_config_id TEXT PRIMARY KEY,
			model_name TEXT NOT NULL,
			provider TEXT,
			model_id TEXT,
			config_json TEXT,
			created_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS subjects (
			subject_id TEXT PRIMARY KEY,
			subject_kind TEXT,
			framework TEXT,
			model TEXT,
			skill TEXT,
			display_name TEXT,
			enabled INTEGER NOT NULL DEFAULT 1,
			tags_json TEXT,
			created_at_utc TEXT NOT NULL,
			updated_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS subject_versions (
			subject_version_id TEXT PRIMARY KEY,
			subject_id TEXT NOT NULL,
			model_config_id TEXT,
			framework_config_sha256 TEXT,
			skill_sha256 TEXT,
			agent_command_sha256 TEXT,
			docker_image TEXT,
			docker_image_digest TEXT,
			sandbox_fingerprint TEXT,
			env_contract_sha256 TEXT,
			created_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS prompt_profiles (
			profile_id TEXT PRIMARY KEY,
			strategy TEXT,
			version_id TEXT,
			template_hash TEXT,
			rules_json TEXT,
			created_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS prompt_renderings (
			prompt_rendering_id TEXT PRIMARY KEY,
			run_id TEXT NOT NULL,
			model TEXT NOT NULL,
			language TEXT NOT NULL,
			sample_id TEXT NOT NULL,
			prompt_artifact_id TEXT,
			prompt_version_id TEXT,
			prompt_mode TEXT,
			created_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS generation_runs (
			run_id TEXT PRIMARY KEY,
			experiment_id TEXT,
			schema_version TEXT NOT NULL,
			created_at_utc TEXT NOT NULL,
			spec_json TEXT NOT NULL,
			dataset_snapshot_id TEXT,
			dataset_fingerprint TEXT,
			prompt_strategy TEXT,
			prompt_version_id TEXT,
			prompt_snapshot_dir TEXT,
			manifest_artifact_id TEXT,
			created_db_at_utc TEXT NOT NULL,
			updated_db_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS generated_cases (
			generated_case_id TEXT PRIMARY KEY,
			run_id TEXT NOT NULL,
			model TEXT NOT NULL,
			subject_id TEXT,
			subject_kind TEXT,
			agent_framework TEXT,
			agent_model TEXT,
			skill_name TEXT,
			skill_version TEXT,
			language TEXT NOT NULL,
			sample_id TEXT NOT NULL,
			sample_uid TEXT,
			sample_path TEXT,
			prompt_rendering_id TEXT,
			generated_test_artifact_id TEXT,
			response_artifact_id TEXT,
			metadata_artifact_id TEXT,
			trace_artifact_id TEXT,
			workspace_diff_artifact_id TEXT,
			sandbox_fingerprint TEXT,
			subject_version_id TEXT,
			framework_config_sha256 TEXT,
			skill_sha256 TEXT,
			agent_command_sha256 TEXT,
			docker_image TEXT,
			docker_image_digest TEXT,
			env_contract_sha256 TEXT,
			generation_key TEXT,
			dependency_fingerprint TEXT,
			generation_env_fingerprint TEXT,
			reused INTEGER NOT NULL DEFAULT 0,
			reuse_stage TEXT,
			reuse_key TEXT,
			reuse_reason TEXT,
			reused_from_run_id TEXT,
			reused_from_case_id TEXT,
			latency_ms INTEGER,
			prompt_tokens INTEGER,
			completion_tokens INTEGER,
			total_tokens INTEGER,
			token_source TEXT,
			estimated_cost_usd REAL,
			cost_source TEXT,
			generated_at_utc TEXT,
			success INTEGER NOT NULL,
			truncated INTEGER NOT NULL DEFAULT 0,
			error_json TEXT,
			created_db_at_utc TEXT NOT NULL,
			updated_db_at_utc TEXT NOT NULL,
			UNIQUE(run_id, model, language, sample_id)
		);`,
		`CREATE TABLE IF NOT EXISTS evaluation_envs (
			env_id TEXT PRIMARY KEY,
			fingerprint TEXT NOT NULL UNIQUE,
			env_json TEXT NOT NULL,
			created_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS score_policies (
			score_policy_id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			policy_json TEXT NOT NULL,
			created_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS evaluation_runs (
			evaluation_run_id TEXT PRIMARY KEY,
			run_id TEXT NOT NULL,
			experiment_id TEXT,
			generation_run_id TEXT,
			schema_version TEXT NOT NULL,
			evaluated_at_utc TEXT NOT NULL,
			manifest_path TEXT,
			result_artifact_id TEXT,
			env_id TEXT,
			score_policy_id TEXT,
			created_db_at_utc TEXT NOT NULL,
			updated_db_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS evaluation_results (
			evaluation_result_id TEXT PRIMARY KEY,
			evaluation_run_id TEXT NOT NULL,
			run_id TEXT NOT NULL,
			model TEXT NOT NULL,
			subject_id TEXT,
			subject_kind TEXT,
			agent_framework TEXT,
			agent_model TEXT,
			skill_name TEXT,
			skill_version TEXT,
			language TEXT NOT NULL,
			sample_id TEXT NOT NULL,
			sample_uid TEXT,
			generated_case_id TEXT,
			generated_test_artifact_id TEXT,
			source_artifact_id TEXT,
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
			prompt_tokens INTEGER,
			completion_tokens INTEGER,
			total_tokens INTEGER,
			token_source TEXT,
			estimated_cost_usd REAL,
			cost_source TEXT,
			truncated INTEGER NOT NULL DEFAULT 0,
			mutation_tool TEXT,
			failure_origin TEXT,
			score_eligible INTEGER NOT NULL DEFAULT 1,
			score_exclusion_reason TEXT,
			compile_error TEXT,
			test_error TEXT,
			coverage_error TEXT,
			mutation_error TEXT,
			trace_artifact_id TEXT,
			workspace_diff_artifact_id TEXT,
			sandbox_fingerprint TEXT,
			evaluation_env_fingerprint TEXT,
			evaluation_key TEXT,
			evaluator_version TEXT,
			mutation_config_sha256 TEXT,
			reused INTEGER NOT NULL DEFAULT 0,
			reuse_stage TEXT,
			reuse_key TEXT,
			reuse_reason TEXT,
			reused_from_run_id TEXT,
			reused_from_result_id TEXT,
			created_db_at_utc TEXT NOT NULL,
			updated_db_at_utc TEXT NOT NULL,
			UNIQUE(evaluation_run_id, model, language, sample_id)
		);`,
		`CREATE TABLE IF NOT EXISTS evaluation_stage_results (
			stage_result_id TEXT PRIMARY KEY,
			evaluation_result_id TEXT NOT NULL,
			stage TEXT NOT NULL,
			command TEXT,
			workdir TEXT,
			exit_code INTEGER,
			duration_ms INTEGER,
			status TEXT NOT NULL,
			stdout_artifact_id TEXT,
			stderr_artifact_id TEXT,
			summary_json TEXT,
			created_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS report_snapshots (
			report_id TEXT PRIMARY KEY,
			run_id TEXT NOT NULL,
			evaluation_run_id TEXT,
			generated_at_utc TEXT NOT NULL,
			source_evaluation TEXT,
			report_json_artifact_id TEXT,
			report_html_artifact_id TEXT,
			summary_json TEXT,
			filters_json TEXT,
			created_db_at_utc TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS report_result_members (
			report_id TEXT NOT NULL,
			evaluation_result_id TEXT NOT NULL,
			PRIMARY KEY(report_id, evaluation_result_id)
		);`,
		`CREATE TABLE IF NOT EXISTS automation_schedules (
			schedule_id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			enabled INTEGER NOT NULL DEFAULT 1,
			trigger_type TEXT NOT NULL,
			cron_expr TEXT,
			interval_seconds INTEGER NOT NULL DEFAULT 0,
			timezone TEXT NOT NULL DEFAULT 'Asia/Shanghai',
			concurrency_policy TEXT NOT NULL DEFAULT 'skip',
			use_docker INTEGER NOT NULL DEFAULT 1,
			run_spec_json TEXT NOT NULL,
			orchestrator_options_json TEXT NOT NULL DEFAULT '{}',
			notify_policy_json TEXT NOT NULL DEFAULT '{}',
			last_fire_at_utc TEXT,
			next_fire_at_utc TEXT,
			created_at_utc TEXT NOT NULL,
			updated_at_utc TEXT NOT NULL,
			deleted_at_utc TEXT
		);`,
		`CREATE INDEX IF NOT EXISTS idx_automation_schedules_due
			ON automation_schedules(enabled, next_fire_at_utc)
			WHERE deleted_at_utc IS NULL;`,
		`CREATE TABLE IF NOT EXISTS automation_runs (
			automation_run_id TEXT PRIMARY KEY,
			schedule_id TEXT NOT NULL,
			run_id TEXT,
			trigger_source TEXT NOT NULL,
			status TEXT NOT NULL,
			started_at_utc TEXT,
			ended_at_utc TEXT,
			error TEXT,
			summary_json TEXT,
			created_at_utc TEXT NOT NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_automation_runs_schedule_created
			ON automation_runs(schedule_id, created_at_utc DESC);`,
		`CREATE TABLE IF NOT EXISTS notification_channels (
			channel_id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			enabled INTEGER NOT NULL DEFAULT 1,
			config_json TEXT NOT NULL,
			created_at_utc TEXT NOT NULL,
			updated_at_utc TEXT NOT NULL,
			deleted_at_utc TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS notification_deliveries (
			delivery_id TEXT PRIMARY KEY,
			automation_run_id TEXT,
			channel_id TEXT NOT NULL,
			status TEXT NOT NULL,
			attempt INTEGER NOT NULL DEFAULT 1,
			request_summary TEXT,
			response_summary TEXT,
			error TEXT,
			created_at_utc TEXT NOT NULL,
			delivered_at_utc TEXT
		);`,
		`CREATE INDEX IF NOT EXISTS idx_notification_deliveries_run
			ON notification_deliveries(automation_run_id, created_at_utc DESC);`,
		`CREATE TABLE IF NOT EXISTS log_events (
			log_event_id TEXT PRIMARY KEY,
			run_id TEXT NOT NULL,
			ts_utc TEXT,
			level TEXT,
			message TEXT NOT NULL,
			fields_json TEXT,
			created_at_utc TEXT NOT NULL
		);`,
	}
	for _, stmt := range ddl {
		if _, err := s.db.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}
	for _, col := range []string{
		"framework_config_sha256 TEXT",
		"skill_sha256 TEXT",
		"agent_command_sha256 TEXT",
		"docker_image TEXT",
		"docker_image_digest TEXT",
		"sandbox_fingerprint TEXT",
		"env_contract_sha256 TEXT",
	} {
		if err := s.ensureColumn(ctx, "subject_versions", col); err != nil {
			return err
		}
	}
	for _, col := range []string{
		"subject_id TEXT",
		"subject_kind TEXT",
		"agent_framework TEXT",
		"agent_model TEXT",
		"skill_name TEXT",
		"skill_version TEXT",
	} {
		if err := s.ensureColumn(ctx, "generated_cases", col); err != nil {
			return err
		}
		if err := s.ensureColumn(ctx, "evaluation_results", col); err != nil {
			return err
		}
	}
	for _, col := range []string{
		"sample_uid TEXT",
		"trace_artifact_id TEXT",
		"workspace_diff_artifact_id TEXT",
		"sandbox_fingerprint TEXT",
		"subject_version_id TEXT",
		"framework_config_sha256 TEXT",
		"skill_sha256 TEXT",
		"agent_command_sha256 TEXT",
		"docker_image TEXT",
		"docker_image_digest TEXT",
		"env_contract_sha256 TEXT",
		"generation_key TEXT",
		"dependency_fingerprint TEXT",
		"generation_env_fingerprint TEXT",
		"reused INTEGER NOT NULL DEFAULT 0",
		"reuse_stage TEXT",
		"reuse_key TEXT",
		"reuse_reason TEXT",
		"reused_from_run_id TEXT",
		"reused_from_case_id TEXT",
		"token_source TEXT",
		"estimated_cost_usd REAL",
		"cost_source TEXT",
	} {
		if err := s.ensureColumn(ctx, "generated_cases", col); err != nil {
			return err
		}
	}
	for _, col := range []string{
		"framework_config_sha256 TEXT",
		"skill_sha256 TEXT",
		"agent_command_sha256 TEXT",
		"docker_image TEXT",
		"docker_image_digest TEXT",
		"env_contract_sha256 TEXT",
		"evaluation_env_fingerprint TEXT",
		"sample_uid TEXT",
		"trace_artifact_id TEXT",
		"workspace_diff_artifact_id TEXT",
		"sandbox_fingerprint TEXT",
		"evaluation_key TEXT",
		"evaluator_version TEXT",
		"mutation_config_sha256 TEXT",
		"reused INTEGER NOT NULL DEFAULT 0",
		"reuse_stage TEXT",
		"reuse_key TEXT",
		"reuse_reason TEXT",
		"reused_from_run_id TEXT",
		"reused_from_result_id TEXT",
		"token_source TEXT",
		"estimated_cost_usd REAL",
		"cost_source TEXT",
	} {
		if err := s.ensureColumn(ctx, "evaluation_results", col); err != nil {
			return err
		}
	}
	for _, stmt := range []string{
		`CREATE INDEX IF NOT EXISTS idx_generated_cases_generation_key_success_time
			ON generated_cases(generation_key, success, generated_at_utc)`,
		`CREATE INDEX IF NOT EXISTS idx_generated_cases_subject_language_sample
			ON generated_cases(subject_id, language, sample_id)`,
		`CREATE INDEX IF NOT EXISTS idx_evaluation_results_evaluation_key_time
			ON evaluation_results(evaluation_key, updated_db_at_utc)`,
		`CREATE INDEX IF NOT EXISTS idx_evaluation_results_subject_language_sample
			ON evaluation_results(subject_id, language, sample_id)`,
	} {
		if _, err := s.db.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}
	_, err := s.db.ExecContext(ctx, `INSERT OR IGNORE INTO schema_migrations(version, applied_at_utc) VALUES(?, ?)`, schemaVersion, nowUTC())
	return err
}
