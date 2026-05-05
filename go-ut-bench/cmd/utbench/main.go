package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"go-ut-bench/internal/contracts"
	"go-ut-bench/internal/dataset"
	"go-ut-bench/internal/evaluator"
	"go-ut-bench/internal/obs"
	"go-ut-bench/internal/orchestrator"
	"go-ut-bench/internal/reporter"
	"go-ut-bench/internal/runner"
	"go-ut-bench/internal/store"
	"go-ut-bench/internal/web"

	"gopkg.in/yaml.v3"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	var args []string
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	var err error
	switch cmd {
	case "run":
		err = runRun(args)
	case "generate":
		err = runGenerate(args)
	case "evaluate":
		err = runEvaluate(args)
	case "report":
		err = runReport(args)
	case "db":
		err = runDB(args)
	case "assets":
		err = runAssets(args)
	case "ingest":
		err = fmt.Errorf("utbench ingest has been replaced by `utbench db ingest-evaluation --evaluation <path>`")
	case "dataset":
		err = runDataset(args)
	case "doctor":
		err = runDoctor(args)
	case "web":
		err = runWeb(args)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`utbench - unified test-bench CLI

Usage:
  utbench run          Run full pipeline (generate -> evaluate -> report)
  utbench generate     Generate unit tests only
  utbench evaluate     Evaluate existing generated tests
  utbench report       Generate reports from evaluation results
  utbench db           Manage SQLite benchmark database
  utbench assets       Query reusable subject/sample assets
  utbench dataset      Dataset management (index, manifest, stats)
  utbench doctor       Check evaluator toolchains with canary tests
  utbench web          Launch Web management UI
  utbench help         Show this help

Run "utbench <command> --help" for more details on a command.

Quick Start:
  1. cp .env.example .env        # 填入你的 API Key
  2. utbench doctor --langs python  # 检查环境是否就绪
  3. utbench web --addr :8080       # 启动 Web UI 开始使用`)
}

func runAssets(args []string) error {
	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		printAssetsUsage()
		return nil
	}
	switch args[0] {
	case "subjects":
		return runAssetsSubjects(args[1:])
	case "generations":
		return runAssetsGenerations(args[1:])
	case "evaluations":
		return runAssetsEvaluations(args[1:])
	case "explain-reuse":
		return runAssetsExplainReuse(args[1:])
	case "help", "-h", "--help":
		printAssetsUsage()
		return nil
	default:
		return fmt.Errorf("unknown assets subcommand: %s", args[0])
	}
}

func printAssetsUsage() {
	fmt.Println(`Usage: utbench assets <subcommand> [flags]

Subcommands:
  subjects           List subject-level assets
  generations        List generated test assets
  evaluations        List evaluation result assets
  explain-reuse      Show the latest reusable generation candidate for a subject/sample

Common flags:
  --db-path ./storage/utbench.db
  --json
  --limit 50`)
}

func runAssetsSubjects(args []string) error {
	fs := flag.NewFlagSet("utbench assets subjects", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println("Usage: utbench assets subjects [flags]")
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}
	dbPath := fs.String("db-path", "./storage/utbench.db", "SQLite database path")
	limit := fs.Int("limit", 100, "Limit")
	jsonOut := fs.Bool("json", false, "Print JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	ctx := context.Background()
	sqliteStore, err := openInitializedStore(ctx, *dbPath)
	if err != nil {
		return err
	}
	defer sqliteStore.Close()
	rows, err := sqliteStore.ListAssetSubjects(ctx, *limit)
	if err != nil {
		return err
	}
	if *jsonOut {
		return printJSON(rows)
	}
	for _, r := range rows {
		fmt.Printf("%s\tkind=%s\tframework=%s\tmodel=%s\tskill=%s\tgenerated=%d\tevaluated=%d\tlatest=%s\n",
			r.SubjectID, r.SubjectKind, r.Framework, r.Model, r.Skill, r.GeneratedCases, r.EvaluationResults, r.LatestGeneratedAt)
	}
	return nil
}

func runAssetsGenerations(args []string) error {
	fs := flag.NewFlagSet("utbench assets generations", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println("Usage: utbench assets generations [flags]")
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}
	dbPath := fs.String("db-path", "./storage/utbench.db", "SQLite database path")
	subjectID := fs.String("subject", "", "Subject ID filter")
	lang := fs.String("lang", "", "Language filter")
	sample := fs.String("sample", "", "Sample ID filter")
	limit := fs.Int("limit", 50, "Limit")
	jsonOut := fs.Bool("json", false, "Print JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	ctx := context.Background()
	sqliteStore, err := openInitializedStore(ctx, *dbPath)
	if err != nil {
		return err
	}
	defer sqliteStore.Close()
	rows, err := sqliteStore.ListAssetGenerations(ctx, *subjectID, *lang, *sample, *limit)
	if err != nil {
		return err
	}
	if *jsonOut {
		return printJSON(rows)
	}
	for _, r := range rows {
		fmt.Printf("%s\t%s\t%s\t%s\tsuccess=%v\treused=%v\tgenerated=%s\tkey=%s\tpath=%s\n",
			r.RunID, r.SubjectID, r.Language, r.SampleID, r.Success, r.Reused, r.GeneratedAtUTC, r.GenerationKey, r.GeneratedTestPath)
	}
	return nil
}

func runAssetsEvaluations(args []string) error {
	fs := flag.NewFlagSet("utbench assets evaluations", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println("Usage: utbench assets evaluations [flags]")
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}
	dbPath := fs.String("db-path", "./storage/utbench.db", "SQLite database path")
	subjectID := fs.String("subject", "", "Subject ID filter")
	lang := fs.String("lang", "", "Language filter")
	sample := fs.String("sample", "", "Sample ID filter")
	limit := fs.Int("limit", 50, "Limit")
	jsonOut := fs.Bool("json", false, "Print JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	ctx := context.Background()
	sqliteStore, err := openInitializedStore(ctx, *dbPath)
	if err != nil {
		return err
	}
	defer sqliteStore.Close()
	rows, err := sqliteStore.ListAssetEvaluations(ctx, *subjectID, *lang, *sample, *limit)
	if err != nil {
		return err
	}
	if *jsonOut {
		return printJSON(rows)
	}
	for _, r := range rows {
		fmt.Printf("%s\t%s\t%s\t%s\tcompile=%v\treused=%v\tupdated=%s\tkey=%s\n",
			r.RunID, r.SubjectID, r.Language, r.SampleID, r.CompilePass, r.Reused, r.CreatedAtUTC, r.EvaluationKey)
	}
	return nil
}

func runAssetsExplainReuse(args []string) error {
	fs := flag.NewFlagSet("utbench assets explain-reuse", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println("Usage: utbench assets explain-reuse --subject <id> --lang <lang> --sample <sample-id> [flags]")
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}
	dbPath := fs.String("db-path", "./storage/utbench.db", "SQLite database path")
	subjectID := fs.String("subject", "", "Subject ID")
	lang := fs.String("lang", "", "Language")
	sample := fs.String("sample", "", "Sample ID")
	generationKey := fs.String("generation-key", "", "Optional current generation key to compare")
	dependencyFingerprint := fs.String("dependency-fingerprint", "", "Optional current dependency fingerprint to compare")
	generationEnvFingerprint := fs.String("generation-env-fingerprint", "", "Optional current generation environment fingerprint to compare")
	limit := fs.Int("limit", 20, "Candidate scan limit")
	jsonOut := fs.Bool("json", false, "Print JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *subjectID == "" {
		return fmt.Errorf("--subject is required")
	}
	if *lang == "" {
		return fmt.Errorf("--lang is required")
	}
	if *sample == "" {
		return fmt.Errorf("--sample is required")
	}

	ctx := context.Background()
	sqliteStore, err := openInitializedStore(ctx, *dbPath)
	if err != nil {
		return err
	}
	defer sqliteStore.Close()
	rows, err := sqliteStore.ListAssetGenerations(ctx, *subjectID, *lang, *sample, *limit)
	if err != nil {
		return err
	}
	type reuseExplanation struct {
		SubjectID      string                        `json:"subject_id"`
		Language       string                        `json:"language"`
		SampleID       string                        `json:"sample_id"`
		Matched        bool                          `json:"matched"`
		MissReason     string                        `json:"miss_reason,omitempty"`
		GenerationKey  string                        `json:"generation_key,omitempty"`
		Comparisons    map[string]string             `json:"environment_fingerprint_comparison,omitempty"`
		LatestReusable *store.DBGenerationAssetItem  `json:"latest_reusable,omitempty"`
		Candidates     []store.DBGenerationAssetItem `json:"candidates,omitempty"`
	}
	explained := reuseExplanation{
		SubjectID:  *subjectID,
		Language:   *lang,
		SampleID:   *sample,
		Candidates: rows,
		MissReason: "no_successful_generation_asset",
	}
	for _, row := range rows {
		if row.Success && row.GeneratedTestPath != "" && row.GenerationKey != "" {
			reusable := row
			explained.Matched = true
			explained.MissReason = ""
			explained.GenerationKey = row.GenerationKey
			explained.LatestReusable = &reusable
			explained.Comparisons = map[string]string{
				"generation_key":             compareOptionalFingerprint(*generationKey, row.GenerationKey),
				"dependency_fingerprint":     compareOptionalFingerprint(*dependencyFingerprint, row.DependencyFingerprint),
				"generation_env_fingerprint": compareOptionalFingerprint(*generationEnvFingerprint, row.GenerationEnvFingerprint),
				"stored_subject_version_id":  row.SubjectVersionID,
				"stored_sandbox_fingerprint": row.SandboxFingerprint,
				"stored_sample_uid":          row.SampleUID,
			}
			break
		}
	}
	if *jsonOut {
		return printJSON(explained)
	}
	fmt.Printf("subject: %s\n", explained.SubjectID)
	fmt.Printf("language: %s\n", explained.Language)
	fmt.Printf("sample: %s\n", explained.SampleID)
	if explained.Matched && explained.LatestReusable != nil {
		fmt.Printf("matched: true\n")
		fmt.Printf("generation_key: %s\n", explained.GenerationKey)
		fmt.Printf("latest_reusable_run: %s\n", explained.LatestReusable.RunID)
		fmt.Printf("generated_case_id: %s\n", explained.LatestReusable.GeneratedCaseID)
		fmt.Printf("generated_test_path: %s\n", explained.LatestReusable.GeneratedTestPath)
		fmt.Println("environment_fingerprint_comparison:")
		for _, key := range []string{"generation_key", "dependency_fingerprint", "generation_env_fingerprint", "stored_subject_version_id", "stored_sandbox_fingerprint", "stored_sample_uid"} {
			fmt.Printf("  %s: %s\n", key, explained.Comparisons[key])
		}
		return nil
	}
	fmt.Printf("matched: false\n")
	fmt.Printf("miss_reason: %s\n", explained.MissReason)
	if len(rows) > 0 {
		fmt.Println("candidates:")
		for _, row := range rows {
			fmt.Printf("  %s\tcase=%s\tsuccess=%v\treused=%v\tkey=%s\tpath=%s\n",
				row.RunID, row.GeneratedCaseID, row.Success, row.Reused, row.GenerationKey, row.GeneratedTestPath)
		}
	}
	return nil
}

func compareOptionalFingerprint(current, stored string) string {
	current = strings.TrimSpace(current)
	stored = strings.TrimSpace(stored)
	if current == "" {
		if stored == "" {
			return "not_recorded"
		}
		return "current_not_provided; stored=" + stored
	}
	if stored == "" {
		return "stored_not_recorded; current=" + current
	}
	if current == stored {
		return "match"
	}
	return "mismatch; current=" + current + "; stored=" + stored
}

func parseCommaList(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func normalizeMutationPolicy(v string) (string, error) {
	policy := strings.ToLower(strings.TrimSpace(v))
	if policy == "" {
		return "", nil
	}
	switch policy {
	case "warn", "fail", "skip":
		return policy, nil
	default:
		return "", fmt.Errorf("unsupported mutation policy: %s", v)
	}
}

func ensureRunID(spec *contracts.RunSpec) {
	if spec.RunID == "" {
		spec.RunID = fmt.Sprintf("run_%d", time.Now().UnixMilli())
	}
}

func newCtx() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}

func withSignal(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()
	return ctx, cancel
}

func runWeb(args []string) error {
	fs := flag.NewFlagSet("utbench web", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println("Usage: utbench web [flags]")
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}

	addr := fs.String("addr", ":8080", "HTTP listen address")
	configPath := fs.String("config", "./configs/models.yaml", "Model config path")
	datasetRoot := fs.String("dataset-root", "./datasets", "Dataset root directory")
	outputRoot := fs.String("output-root", "./artifacts", "Output root directory")
	dbPath := fs.String("db-path", "./storage/utbench.db", "SQLite database path")
	agentsConfigPath := fs.String("agents-config", "", "Agent/skill config path (auto-detected from --config dir if omitted)")
	imageName := fs.String("docker-image", "", "Deprecated alias for --docker-eval-image")
	evalImageName := fs.String("docker-eval-image", "utbench:latest", "Docker image for containerized evaluation runs")
	projectRoot := fs.String("project-root", ".", "Project root mounted into Docker")
	envFile := fs.String("env-file", "./.env", "Environment file passed to Docker runs")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Auto-detect agents config from the same directory as models.yaml
	if *agentsConfigPath == "" {
		cfgDir := filepath.Dir(*configPath)
		for _, candidate := range []string{"agents.yaml", "agents.example.yaml"} {
			p := filepath.Join(cfgDir, candidate)
			if _, err := os.Stat(p); err == nil {
				*agentsConfigPath = p
				break
			}
		}
	}

	absProjectRoot, err := filepath.Abs(*projectRoot)
	if err != nil {
		return fmt.Errorf("resolve project root: %w", err)
	}

	// When running inside a container, --project-root defaults to the
	// container-internal working directory (e.g. "/app").  Docker daemon
	// running on the host cannot resolve that path for bind mounts.
	// If UTBENCH_SANDBOX_HOST_OUTPUT_ROOT is set, derive the host project
	// root from it (the parent of the "artifacts" directory).
	if *projectRoot == "." {
		if hostOutputRoot := strings.TrimSpace(os.Getenv("UTBENCH_SANDBOX_HOST_OUTPUT_ROOT")); hostOutputRoot != "" {
			if _, err := os.Stat("/.dockerenv"); err == nil || strings.TrimSpace(os.Getenv("container")) != "" {
				absProjectRoot = filepath.Dir(hostOutputRoot)
			}
		}
	}
	absEnvFile := *envFile
	if absEnvFile != "" {
		absEnvFile, err = filepath.Abs(absEnvFile)
		if err != nil {
			return fmt.Errorf("resolve env file: %w", err)
		}
		// When running inside a container with the default env file path,
		// map it to the host path so Docker daemon can access it.
		if *envFile == "./.env" && absProjectRoot != filepath.Clean(".") {
			absEnvFile = filepath.Join(absProjectRoot, ".env")
		}
	}

	// Load .env file into process environment so os.Getenv() can read API keys
	if err := web.LoadEnvFile(absEnvFile); err != nil {
		return fmt.Errorf("load env file: %w", err)
	}

	// Check API keys and warn about missing ones
	checkModelAPIKeys(*configPath)

	if strings.TrimSpace(*imageName) != "" {
		*evalImageName = *imageName
	}
	dockerCfg := web.DockerConfig{
		EvalImageName: *evalImageName,
		ProjectRoot:   absProjectRoot,
		EnvFile:       absEnvFile,
	}
	if *agentsConfigPath != "" {
		fmt.Printf("[web] agents config: %s\n", *agentsConfigPath)
	} else {
		fmt.Printf("[web] agents config: not found (searched in %s)\n", filepath.Dir(*configPath))
	}
	mgr := web.NewRunManager(*configPath, *agentsConfigPath, *datasetRoot, *outputRoot, *dbPath, dockerCfg)
	bld := web.NewBuildManager(absProjectRoot)
	server, err := web.NewServer(mgr, bld, *configPath, *outputRoot, *dbPath, dockerCfg)
	if err != nil {
		return fmt.Errorf("create web server: %w", err)
	}
	defer server.Close()

	// 优雅关闭：监听 SIGINT/SIGTERM，收到信号后先清理再退出
	return server.StartGraceful(*addr)
}

// checkModelAPIKeys 读取 models.yaml 并检查各模型的 API Key 环境变量是否已设置。
// 缺失的 Key 会在启动时打印警告，帮助用户尽早发现问题。
func checkModelAPIKeys(configPath string) {
	raw, err := os.ReadFile(configPath)
	if err != nil {
		return // config 读取失败不阻塞启动
	}
	var cfg struct {
		Models map[string]struct {
			Enabled bool `yaml:"enabled"`
			Config  struct {
				APIKeyEnv string `yaml:"api_key_env"`
			} `yaml:"config"`
		} `yaml:"models"`
	}
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return
	}
	var missing []string
	for name, m := range cfg.Models {
		if !m.Enabled {
			continue
		}
		key := strings.TrimSpace(m.Config.APIKeyEnv)
		if key == "" {
			continue
		}
		val := strings.TrimSpace(os.Getenv(key))
		lower := strings.ToLower(val)
		if val == "" || strings.Contains(lower, "your_") || strings.Contains(lower, "placeholder") {
			missing = append(missing, fmt.Sprintf("  %-22s → %s", name, key))
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		fmt.Fprintf(os.Stderr, "\n⚠  以下模型的 API Key 未配置，运行时会失败:\n%s\n", strings.Join(missing, "\n"))
		fmt.Fprintf(os.Stderr, "   请在 .env 文件中设置对应变量，参考 .env.example\n\n")
	}
}

func runRun(args []string) error {
	fs := flag.NewFlagSet("utbench run", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println("Usage: utbench run [flags]")
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}

	ingest := fs.Bool("ingest", false, "Ingest results into SQLite after run")
	dbPath := fs.String("db-path", "./storage/utbench.db", "SQLite database path")
	verbose := fs.Bool("v", false, "Verbose output")
	config := fs.String("config", "../benchmark/config/models.yaml", "Model config path")
	agentsConfig := fs.String("agents-config", "", "Agent/skill config path")
	outputRoot := fs.String("output-root", "./artifacts", "Output root directory")
	datasetRoot := fs.String("dataset-root", "./datasets", "Dataset root directory")
	datasetManifest := fs.String("dataset-manifest", "", "Dataset manifest path")
	datasetLevel := fs.String("level", "", "Dataset level")
	datasetClass := fs.String("class", "self_contained", "Dataset class(es), comma-separated (self_contained, module_level)")
	datasetScenario := fs.String("scenario", "", "Dataset scenario (boundary, simple_function, complex_dependency, interface_mock)")
	maxSamples := fs.Int("max-samples", 0, "Max samples")
	mode := fs.String("mode", "full", "Run mode (full, incremental)")
	resetCheckpoint := fs.Bool("reset-checkpoint", false, "Reset checkpoint")
	dryRun := fs.Bool("dry-run", false, "Dry run (skip API calls)")
	reuseGenerated := fs.Bool("reuse-generated", true, "Reuse matching generated tests from SQLite before calling models")
	reuseEvaluation := fs.Bool("reuse-evaluation", false, "Reuse matching evaluation results from SQLite when environment keys match")
	mutationEnabled := fs.Bool("mutation-enabled", true, "Enable mutation testing")
	mutationTimeout := fs.Int("mutation-timeout", 600, "Mutation timeout (seconds)")
	mutationPolicy := fs.String("mutation-policy", "warn", "Mutation policy (warn, fail)")
	testTimeout := fs.Int("test-timeout", 180, "Test execution timeout (seconds)")
	workers := fs.Int("workers", 16, "Number of concurrent workers (default 16)")
	models := fs.String("models", "", "Comma-separated models")
	subjects := fs.String("subjects", "", "Comma-separated subjects (framework__model__skill)")
	langs := fs.String("langs", "", "Comma-separated languages")
	runID := fs.String("run-id", "", "Run ID")
	evalBackendFlag := fs.String("eval-backend", "local", "Evaluation backend (local, docker)")
	evalDockerImage := fs.String("eval-docker-image", "utbench:latest", "Docker image for docker eval backend")

	if err := fs.Parse(args); err != nil {
		return err
	}
	policy, err := normalizeMutationPolicy(*mutationPolicy)
	if err != nil {
		return err
	}

	spec := contracts.RunSpec{
		ConfigPath:       *config,
		AgentsConfigPath: *agentsConfig,
		OutputRoot:       *outputRoot,
		DatasetRoot:      *datasetRoot,
		DatasetManifest:  *datasetManifest,
		DatasetLevel:     *datasetLevel,
		DatasetClasses:   parseCommaList(*datasetClass),
		DatasetScenario:  *datasetScenario,
		MaxSamples:       *maxSamples,
		Workers:          *workers,
		Mode:             contracts.RunMode(*mode),
		ResetCheckpoint:  *resetCheckpoint,
		DryRun:           *dryRun,
		ReuseGenerated:   *reuseGenerated,
		ReuseEvaluation:  *reuseEvaluation,
		DBPath:           *dbPath,
		MutationEnabled:  *mutationEnabled,
		MutationTimeout:  *mutationTimeout,
		MutationPolicy:   policy,
		TestTimeout:      *testTimeout,
		Models:           parseCommaList(*models),
		Subjects:         parseCommaList(*subjects),
		Languages:        parseCommaList(*langs),
		RunID:            *runID,
		CreatedAtUTC:     time.Now().UTC(),
	}
	ensureRunID(&spec)

	logDir := filepath.Join(*outputRoot, "runs", spec.RunID, "logs")
	logger := obs.NewLogger(*verbose, logDir)
	ctx, cancel := withSignal(context.Background())
	defer cancel()

	datasetSvc := dataset.NewService()
	runnerSvc := runner.NewService(logger)
	evaluatorSvc := evaluator.NewService(logger)
	reporterSvc := reporter.NewService(logger, &runner.DefaultPromptMetaProvider{})

	// 设置评测后端
	switch strings.ToLower(strings.TrimSpace(*evalBackendFlag)) {
	case "docker":
		evaluator.SetEvalBackend(evaluator.NewDockerBackend(*evalDockerImage))
		logger.Info("eval backend set to docker", "image", *evalDockerImage)
	case "local", "":
		// 默认 LocalBackend，无需设置
	default:
		return fmt.Errorf("unknown eval-backend: %s (supported: local, docker)", *evalBackendFlag)
	}

	svc := orchestrator.New(datasetSvc, runnerSvc, evaluatorSvc, reporterSvc)
	opts := orchestrator.Options{Ingest: *ingest, DBPath: *dbPath}

	result, err := svc.Run(ctx, spec, opts)
	if err != nil {
		return fmt.Errorf("run failed: %w", err)
	}

	fmt.Printf("Run completed: %s\n", result.RunID)
	fmt.Printf("  Manifest:   %s\n", result.ManifestPath)
	fmt.Printf("  Evaluation: %s\n", result.EvaluationPath)
	fmt.Printf("  Report JSON: %s\n", result.ReportJSONPath)
	fmt.Printf("  Report HTML: %s\n", result.ReportHTMLPath)
	if result.Ingested {
		fmt.Printf("  Ingested:   yes\n")
	}
	return nil
}

func runGenerate(args []string) error {
	fs := flag.NewFlagSet("utbench generate", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println("Usage: utbench generate [flags]")
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}

	verbose := fs.Bool("v", false, "Verbose output")
	config := fs.String("config", "../benchmark/config/models.yaml", "Model config path")
	agentsConfig := fs.String("agents-config", "", "Agent/skill config path")
	outputRoot := fs.String("output-root", "./artifacts", "Output root directory")
	datasetRoot := fs.String("dataset-root", "./datasets", "Dataset root directory")
	datasetManifest := fs.String("dataset-manifest", "", "Dataset manifest path")
	datasetLevel := fs.String("level", "", "Dataset level")
	datasetClass := fs.String("class", "self_contained", "Dataset class")
	datasetScenario := fs.String("scenario", "", "Dataset scenario")
	maxSamples := fs.Int("max-samples", 0, "Max samples")
	mode := fs.String("mode", "full", "Run mode")
	resetCheckpoint := fs.Bool("reset-checkpoint", false, "Reset checkpoint")
	dryRun := fs.Bool("dry-run", false, "Dry run")
	reuseGenerated := fs.Bool("reuse-generated", true, "Reuse matching generated tests from SQLite before calling models")
	dbPath := fs.String("db-path", "./storage/utbench.db", "SQLite database path")
	models := fs.String("models", "", "Comma-separated models")
	subjects := fs.String("subjects", "", "Comma-separated subjects (framework__model__skill)")
	langs := fs.String("langs", "", "Comma-separated languages")
	runID := fs.String("run-id", "", "Run ID")

	if err := fs.Parse(args); err != nil {
		return err
	}

	spec := contracts.RunSpec{
		ConfigPath:       *config,
		AgentsConfigPath: *agentsConfig,
		OutputRoot:       *outputRoot,
		DatasetRoot:      *datasetRoot,
		DatasetManifest:  *datasetManifest,
		DatasetLevel:     *datasetLevel,
		DatasetClasses:   parseCommaList(*datasetClass),
		DatasetScenario:  *datasetScenario,
		MaxSamples:       *maxSamples,
		Mode:             contracts.RunMode(*mode),
		ResetCheckpoint:  *resetCheckpoint,
		DryRun:           *dryRun,
		ReuseGenerated:   *reuseGenerated,
		DBPath:           *dbPath,
		Models:           parseCommaList(*models),
		Subjects:         parseCommaList(*subjects),
		Languages:        parseCommaList(*langs),
		RunID:            *runID,
		CreatedAtUTC:     time.Now().UTC(),
	}
	ensureRunID(&spec)

	logger := obs.NewLogger(*verbose, "")
	ctx, cancel := withSignal(context.Background())
	defer cancel()

	datasetSvc := dataset.NewService()
	runnerSvc := runner.NewService(logger)

	if err := datasetSvc.ValidateSpec(spec); err != nil {
		return fmt.Errorf("invalid spec: %w", err)
	}

	samples, err := datasetSvc.DiscoverSamples(spec)
	if err != nil {
		return fmt.Errorf("discover samples: %w", err)
	}

	var reuseStore runner.GenerationReuseStore
	if spec.ReuseGenerated && strings.TrimSpace(spec.DBPath) != "" && !spec.DryRun {
		if db, openErr := store.OpenSQLite(spec.DBPath); openErr == nil {
			if initErr := db.Init(ctx); initErr == nil {
				reuseStore = db
			} else {
				_ = db.Close()
			}
		}
	}
	output, err := runnerSvc.Generate(ctx, spec, samples, reuseStore)
	if reuseStore != nil {
		if closer, ok := reuseStore.(interface{ Close() error }); ok {
			_ = closer.Close()
		}
	}
	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	fmt.Printf("Generated %d cases\n", len(output.Manifest.Cases))
	fmt.Printf("Manifest: %s\n", output.ManifestPath)
	return nil
}

func runEvaluate(args []string) error {
	fs := flag.NewFlagSet("utbench evaluate", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println("Usage: utbench evaluate [flags]")
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}

	verbose := fs.Bool("v", false, "Verbose output")
	outputRoot := fs.String("output-root", "./artifacts", "Output root directory")
	manifestPath := fs.String("manifest", "", "Path to generated_manifest.json (required)")
	dbPath := fs.String("db-path", "./storage/utbench.db", "SQLite database path")
	reuseEvaluation := fs.Bool("reuse-evaluation", false, "Reuse matching evaluation results from SQLite when environment keys match")
	mutationEnabled := fs.Bool("mutation-enabled", true, "Enable mutation testing")
	mutationTimeout := fs.Int("mutation-timeout", 600, "Mutation timeout (seconds)")
	mutationPolicy := fs.String("mutation-policy", "warn", "Mutation policy")
	testTimeout := fs.Int("test-timeout", 180, "Test execution timeout (seconds)")
	runID := fs.String("run-id", "", "Run ID")
	evalBackendFlag := fs.String("eval-backend", "local", "Evaluation backend (local, docker)")
	evalDockerImage := fs.String("eval-docker-image", "utbench:latest", "Docker image for docker eval backend")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *manifestPath == "" {
		return fmt.Errorf("--manifest is required")
	}
	policy, err := normalizeMutationPolicy(*mutationPolicy)
	if err != nil {
		return err
	}

	spec := contracts.RunSpec{
		OutputRoot:      *outputRoot,
		DBPath:          *dbPath,
		ReuseEvaluation: *reuseEvaluation,
		MutationEnabled: *mutationEnabled,
		MutationTimeout: *mutationTimeout,
		MutationPolicy:  policy,
		TestTimeout:     *testTimeout,
		RunID:           *runID,
		CreatedAtUTC:    time.Now().UTC(),
	}
	ensureRunID(&spec)

	logDir := filepath.Join(spec.OutputRoot, "runs", spec.RunID, "logs")
	logger := obs.NewLogger(*verbose, logDir)
	ctx, cancel := withSignal(context.Background())
	defer cancel()

	// 设置评测后端
	switch strings.ToLower(strings.TrimSpace(*evalBackendFlag)) {
	case "docker":
		evaluator.SetEvalBackend(evaluator.NewDockerBackend(*evalDockerImage))
	case "local", "":
		// 默认 LocalBackend
	default:
		return fmt.Errorf("unknown eval-backend: %s (supported: local, docker)", *evalBackendFlag)
	}

	evaluatorSvc := evaluator.NewService(logger)
	output, err := evaluatorSvc.Evaluate(ctx, spec, *manifestPath)
	if err != nil {
		return fmt.Errorf("evaluate: %w", err)
	}

	fmt.Printf("Evaluated %d results\n", len(output.Result.Results))
	fmt.Printf("Result: %s\n", output.ResultPath)
	return nil
}

func runReport(args []string) error {
	fs := flag.NewFlagSet("utbench report", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println("Usage: utbench report [flags]")
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}

	verbose := fs.Bool("v", false, "Verbose output")
	outputRoot := fs.String("output-root", "./artifacts", "Output root directory")
	evaluationPath := fs.String("evaluation", "", "Path to evaluation_result.json (required)")
	runID := fs.String("run-id", "", "Run ID")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *evaluationPath == "" {
		return fmt.Errorf("--evaluation is required")
	}

	spec := contracts.RunSpec{
		OutputRoot: *outputRoot,
		RunID:      *runID,
	}
	ensureRunID(&spec)

	logDir := filepath.Join(spec.OutputRoot, "runs", spec.RunID, "logs")
	logger := obs.NewLogger(*verbose, logDir)
	reporterSvc := reporter.NewService(logger, &runner.DefaultPromptMetaProvider{})
	output, err := reporterSvc.Generate(context.Background(), spec, *evaluationPath)
	if err != nil {
		return fmt.Errorf("report: %w", err)
	}

	fmt.Printf("Report JSON: %s\n", output.ReportJSONPath)
	fmt.Printf("Report HTML: %s\n", output.ReportHTMLPath)
	return nil
}

func runDB(args []string) error {
	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		printDBUsage()
		return nil
	}
	switch args[0] {
	case "init":
		return runDBInit(args[1:])
	case "ingest-manifest":
		return runDBIngestManifest(args[1:])
	case "ingest-evaluation":
		return runDBIngestEvaluation(args[1:])
	case "ingest-report":
		return runDBIngestReport(args[1:])
	case "ingest-run":
		return runDBIngestRun(args[1:])
	case "overview":
		return runDBOverview(args[1:])
	case "list-runs":
		return runDBListRuns(args[1:])
	case "list-results":
		return runDBListResults(args[1:])
	case "report":
		return runDBReport(args[1:])
	case "help", "-h", "--help":
		printDBUsage()
		return nil
	default:
		return fmt.Errorf("unknown db subcommand: %s", args[0])
	}
}

func printDBUsage() {
	fmt.Println(`Usage: utbench db <subcommand> [flags]

Subcommands:
  init                Initialize SQLite schema
  ingest-manifest     Ingest generated_manifest.json and linked artifacts
  ingest-evaluation   Ingest evaluation_result.json and linked manifest/artifacts
  ingest-report       Ingest report_summary.json and linked report HTML
  ingest-run          Ingest all known artifacts from artifacts/runs/<run-id>
  overview            Print database counts and latest runs
  list-runs           List runs stored in database
  list-results        List evaluation results stored in database
  report              Generate report from database-selected results`)
}

func openInitializedStore(ctx context.Context, dbPath string) (*store.SQLiteStore, error) {
	sqliteStore, err := store.OpenSQLite(dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if err := sqliteStore.Init(ctx); err != nil {
		_ = sqliteStore.Close()
		return nil, fmt.Errorf("init sqlite: %w", err)
	}
	return sqliteStore, nil
}

func runDBInit(args []string) error {
	fs := flag.NewFlagSet("utbench db init", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println("Usage: utbench db init [flags]")
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}
	dbPath := fs.String("db-path", "./storage/utbench.db", "SQLite database path")
	if err := fs.Parse(args); err != nil {
		return err
	}
	ctx := context.Background()
	sqliteStore, err := openInitializedStore(ctx, *dbPath)
	if err != nil {
		return err
	}
	defer sqliteStore.Close()
	fmt.Printf("Initialized database: %s\n", *dbPath)
	return nil
}

func runDBIngestManifest(args []string) error {
	fs := flag.NewFlagSet("utbench db ingest-manifest", flag.ContinueOnError)
	dbPath := fs.String("db-path", "./storage/utbench.db", "SQLite database path")
	manifestPath := fs.String("manifest", "", "Path to generated_manifest.json (required)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *manifestPath == "" {
		return fmt.Errorf("--manifest is required")
	}
	ctx := context.Background()
	sqliteStore, err := openInitializedStore(ctx, *dbPath)
	if err != nil {
		return err
	}
	defer sqliteStore.Close()
	sum, err := sqliteStore.IngestManifestFile(ctx, *manifestPath)
	if err != nil {
		return fmt.Errorf("ingest manifest: %w", err)
	}
	printIngestSummary(sum, *dbPath)
	return nil
}

func runDBIngestEvaluation(args []string) error {
	fs := flag.NewFlagSet("utbench db ingest-evaluation", flag.ContinueOnError)
	dbPath := fs.String("db-path", "./storage/utbench.db", "SQLite database path")
	evaluationPath := fs.String("evaluation", "", "Path to evaluation_result.json (required)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *evaluationPath == "" {
		return fmt.Errorf("--evaluation is required")
	}
	ctx := context.Background()
	sqliteStore, err := openInitializedStore(ctx, *dbPath)
	if err != nil {
		return err
	}
	defer sqliteStore.Close()
	sum, err := sqliteStore.IngestEvaluationFile(ctx, *evaluationPath)
	if err != nil {
		return fmt.Errorf("ingest evaluation: %w", err)
	}
	printIngestSummary(sum, *dbPath)
	return nil
}

func runDBIngestReport(args []string) error {
	fs := flag.NewFlagSet("utbench db ingest-report", flag.ContinueOnError)
	dbPath := fs.String("db-path", "./storage/utbench.db", "SQLite database path")
	reportPath := fs.String("report", "", "Path to report_summary.json (required)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *reportPath == "" {
		return fmt.Errorf("--report is required")
	}
	ctx := context.Background()
	sqliteStore, err := openInitializedStore(ctx, *dbPath)
	if err != nil {
		return err
	}
	defer sqliteStore.Close()
	sum, err := sqliteStore.IngestReportFile(ctx, *reportPath)
	if err != nil {
		return fmt.Errorf("ingest report: %w", err)
	}
	printIngestSummary(sum, *dbPath)
	return nil
}

func runDBIngestRun(args []string) error {
	fs := flag.NewFlagSet("utbench db ingest-run", flag.ContinueOnError)
	dbPath := fs.String("db-path", "./storage/utbench.db", "SQLite database path")
	outputRoot := fs.String("output-root", "./artifacts", "Output root directory")
	runID := fs.String("run-id", "", "Run ID")
	runDir := fs.String("run-dir", "", "Run directory")
	if err := fs.Parse(args); err != nil {
		return err
	}
	dir := *runDir
	if dir == "" {
		if *runID == "" {
			return fmt.Errorf("--run-id or --run-dir is required")
		}
		dir = filepath.Join(*outputRoot, "runs", *runID)
	}
	ctx := context.Background()
	sqliteStore, err := openInitializedStore(ctx, *dbPath)
	if err != nil {
		return err
	}
	defer sqliteStore.Close()
	sum, err := sqliteStore.IngestRun(ctx, store.IngestRunOptions{RunDir: dir})
	if err != nil {
		return fmt.Errorf("ingest run: %w", err)
	}
	printIngestSummary(sum, *dbPath)
	return nil
}

func runDBOverview(args []string) error {
	fs := flag.NewFlagSet("utbench db overview", flag.ContinueOnError)
	dbPath := fs.String("db-path", "./storage/utbench.db", "SQLite database path")
	limit := fs.Int("limit", 10, "Latest run limit")
	jsonOut := fs.Bool("json", false, "Print JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	ctx := context.Background()
	sqliteStore, err := openInitializedStore(ctx, *dbPath)
	if err != nil {
		return err
	}
	defer sqliteStore.Close()
	overview, err := sqliteStore.Overview(ctx, *limit)
	if err != nil {
		return err
	}
	if *jsonOut {
		return printJSON(overview)
	}
	fmt.Printf("Database: %s\n", overview.DBPath)
	fmt.Printf("  schema: %s\n", overview.SchemaVersion)
	fmt.Printf("  generation_runs: %d\n", overview.GenerationRuns)
	fmt.Printf("  evaluation_runs: %d\n", overview.EvaluationRuns)
	fmt.Printf("  generated_cases: %d\n", overview.GeneratedCases)
	fmt.Printf("  evaluation_results: %d\n", overview.EvaluationResults)
	fmt.Printf("  artifacts: %d\n", overview.Artifacts)
	fmt.Printf("  reports: %d\n", overview.Reports)
	for _, r := range overview.LatestRuns {
		fmt.Printf("  [run] %s generated=%d evaluated=%d models=%d langs=%d\n", r.RunID, r.GeneratedCases, r.EvaluationResults, r.Models, r.Languages)
	}
	return nil
}

func runDBListRuns(args []string) error {
	fs := flag.NewFlagSet("utbench db list-runs", flag.ContinueOnError)
	dbPath := fs.String("db-path", "./storage/utbench.db", "SQLite database path")
	limit := fs.Int("limit", 50, "Limit")
	jsonOut := fs.Bool("json", false, "Print JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	ctx := context.Background()
	sqliteStore, err := openInitializedStore(ctx, *dbPath)
	if err != nil {
		return err
	}
	defer sqliteStore.Close()
	rows, err := sqliteStore.ListRuns(ctx, *limit)
	if err != nil {
		return err
	}
	if *jsonOut {
		return printJSON(rows)
	}
	for _, r := range rows {
		fmt.Printf("%s\tgenerated=%d\tevaluated=%d\tmodels=%d\tlangs=%d\n", r.RunID, r.GeneratedCases, r.EvaluationResults, r.Models, r.Languages)
	}
	return nil
}

func runDBListResults(args []string) error {
	fs := flag.NewFlagSet("utbench db list-results", flag.ContinueOnError)
	dbPath := fs.String("db-path", "./storage/utbench.db", "SQLite database path")
	runID := fs.String("run-id", "", "Run ID filter")
	model := fs.String("model", "", "Model filter")
	lang := fs.String("lang", "", "Language filter")
	limit := fs.Int("limit", 50, "Limit")
	jsonOut := fs.Bool("json", false, "Print JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	ctx := context.Background()
	sqliteStore, err := openInitializedStore(ctx, *dbPath)
	if err != nil {
		return err
	}
	defer sqliteStore.Close()
	rows, err := sqliteStore.ListResults(ctx, *runID, *model, *lang, *limit)
	if err != nil {
		return err
	}
	if *jsonOut {
		return printJSON(rows)
	}
	for _, r := range rows {
		testPass := "nil"
		if r.TestPass != nil {
			testPass = fmt.Sprintf("%v", *r.TestPass)
		}
		fmt.Printf("%s\t%s\t%s\t%s\tcompile=%v\ttest=%s\tmutation=%s\n", r.RunID, r.Model, r.Language, r.SampleID, r.CompilePass, testPass, formatNullableFloat(r.MutationScore))
	}
	return nil
}

func runDBReport(args []string) error {
	fs := flag.NewFlagSet("utbench db report", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println("Usage: utbench db report [flags]")
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}
	dbPath := fs.String("db-path", "./storage/utbench.db", "SQLite database path")
	outputRoot := fs.String("output-root", "./artifacts", "Output root directory")
	configPath := fs.String("config", "./configs/models.yaml", "Model config path")
	runIDs := fs.String("run-ids", "", "Comma-separated source run IDs")
	evaluationRunIDs := fs.String("evaluation-run-ids", "", "Comma-separated evaluation run IDs")
	models := fs.String("models", "", "Comma-separated model filter")
	langs := fs.String("langs", "", "Comma-separated language filter")
	runID := fs.String("run-id", "", "Output report run ID")
	scoreEligibleOnly := fs.Bool("score-eligible-only", false, "Only include score-eligible rows")
	ingestReport := fs.Bool("ingest-report", true, "Ingest generated DB report back into SQLite")
	verbose := fs.Bool("v", false, "Verbose output")
	if err := fs.Parse(args); err != nil {
		return err
	}
	outRunID := strings.TrimSpace(*runID)
	if outRunID == "" {
		outRunID = contracts.NewRunID() + "_db_report"
	}
	ctx := context.Background()
	sqliteStore, err := openInitializedStore(ctx, *dbPath)
	if err != nil {
		return err
	}
	defer sqliteStore.Close()
	filter := store.DBReportFilter{
		RunIDs:            parseCommaList(*runIDs),
		EvaluationRunIDs:  parseCommaList(*evaluationRunIDs),
		Models:            parseCommaList(*models),
		Languages:         parseCommaList(*langs),
		ScoreEligibleOnly: *scoreEligibleOnly,
	}
	set, err := sqliteStore.SelectEvaluationResultSet(ctx, outRunID, filter)
	if err != nil {
		return fmt.Errorf("select db results: %w", err)
	}
	logger := obs.NewLogger(*verbose, filepath.Join(*outputRoot, "runs", outRunID, "logs"))
	reporterSvc := reporter.NewService(logger, &runner.DefaultPromptMetaProvider{})
	spec := contracts.RunSpec{
		RunID:      outRunID,
		OutputRoot: *outputRoot,
		ConfigPath: *configPath,
	}
	source := "db://" + *dbPath
	output, err := reporterSvc.GenerateFromResultSet(spec, set, source)
	if err != nil {
		return fmt.Errorf("generate db report: %w", err)
	}
	if *ingestReport {
		if _, err := sqliteStore.IngestReportFile(ctx, output.ReportJSONPath); err != nil {
			return fmt.Errorf("ingest db report: %w", err)
		}
	}
	fmt.Printf("DB report generated from %d results\n", len(set.Results))
	fmt.Printf("Report JSON: %s\n", output.ReportJSONPath)
	fmt.Printf("Report HTML: %s\n", output.ReportHTMLPath)
	return nil
}

func printIngestSummary(sum store.IngestSummary, dbPath string) {
	fmt.Printf("Ingested run %s into %s\n", sum.RunID, dbPath)
	fmt.Printf("  manifest=%v evaluation=%v report=%v\n", sum.ManifestIngested, sum.EvaluationIngested, sum.ReportIngested)
	fmt.Printf("  generated_cases=%d evaluation_results=%d artifacts=%d unavailable_artifacts=%d\n", sum.GenerationCases, sum.EvaluationResults, sum.ArtifactsIndexed, sum.UnavailableArtifacts)
}

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func formatNullableFloat(v *float64) string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprintf("%.2f", *v)
}

func runDataset(args []string) error {
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		svc := dataset.NewService()
		return runDatasetSubcommand(svc, args[0], args[1:])
	}

	fs := flag.NewFlagSet("utbench dataset", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println("Usage: utbench dataset <subcommand> [flags]")
		fmt.Println("Subcommands:")
		fmt.Println("  stats            Show dataset file counts")
		fmt.Println("  index            Build dataset index")
		fmt.Println("  manifest         Build dataset manifest (L1/L2)")
		fmt.Println("\nRun 'utbench dataset <subcommand> --help' for subcommand flags.")
	}

	subCmd := fs.String("cmd", "stats", "Sub-command: stats, index, manifest")

	if err := fs.Parse(args); err != nil {
		return err
	}

	remaining := fs.Args()
	svc := dataset.NewService()
	return runDatasetSubcommand(svc, *subCmd, remaining)
}

func runDatasetSubcommand(svc *dataset.Service, subCmd string, remaining []string) error {
	switch subCmd {
	case "stats":
		return datasetStats(svc, remaining)
	case "index":
		return datasetIndex(svc, remaining)
	case "manifest":
		return datasetManifest(svc, remaining)
	case "validate":
		return datasetValidate(svc, remaining)
	default:
		return fmt.Errorf("unknown sub-command: %s (use: stats, index, manifest, validate)", subCmd)
	}
}

func datasetStats(svc *dataset.Service, args []string) error {
	fs := flag.NewFlagSet("utbench dataset stats", flag.ContinueOnError)
	datasetRoot := fs.String("dataset-root", "./datasets", "Dataset root directory")
	fs.Parse(args)

	langs := []string{"python", "java", "go", "cpp"}
	total := 0
	for _, lang := range langs {
		langDir := filepath.Join(*datasetRoot, lang)
		count := 0
		filepath.WalkDir(langDir, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return nil
			}
			if !d.IsDir() {
				count++
			}
			return nil
		})
		fmt.Printf("  %s: %d files\n", lang, count)
		total += count
	}
	fmt.Printf("  total: %d files\n", total)
	return nil
}

func datasetIndex(svc *dataset.Service, args []string) error {
	fs := flag.NewFlagSet("utbench dataset index", flag.ContinueOnError)
	datasetRoot := fs.String("dataset-root", "./datasets", "Dataset root directory")
	output := fs.String("output", "./configs/dataset_index.json", "Index output path")
	fs.Parse(args)

	summary, err := svc.BuildIndex(*datasetRoot, *output)
	if err != nil {
		return fmt.Errorf("build index: %w", err)
	}
	fmt.Printf("Indexed %d samples -> %s\n", summary.Total, summary.Path)
	return nil
}

func datasetManifest(svc *dataset.Service, args []string) error {
	fs := flag.NewFlagSet("utbench dataset manifest", flag.ContinueOnError)
	indexPath := fs.String("index", "./configs/dataset_index.json", "Index path")
	level := fs.String("level", "l1", "Manifest level (l1, l2, ...)")
	output := fs.String("output", "", "Manifest output path")
	limit := fs.Int("limit-per-scenario", 20, "Limit per scenario (0 = all)")
	langs := fs.String("langs", "", "Comma-separated languages to include")
	classFilter := fs.String("class", "", "Dataset class filter")
	scenarioFilter := fs.String("scenario", "", "Dataset scenario filter")
	fs.Parse(args)

	if *output == "" {
		*output = fmt.Sprintf("./configs/dataset_%s.json", *level)
	}

	langList := parseCommaList(*langs)
	if len(langList) == 0 {
		langList = contracts.SupportedLanguages
	}

	opts := dataset.ManifestBuildOptions{
		IndexPath:        *indexPath,
		Level:            *level,
		OutputPath:       *output,
		Languages:        langList,
		ClassFilter:      *classFilter,
		ScenarioFilter:   *scenarioFilter,
		LimitPerScenario: *limit,
	}

	summary, err := svc.BuildManifest(opts)
	if err != nil {
		return fmt.Errorf("build manifest: %w", err)
	}
	fmt.Printf("Built manifest with %d samples -> %s\n", summary.Total, summary.Path)
	return nil
}

func datasetValidate(svc *dataset.Service, args []string) error {
	fs := flag.NewFlagSet("utbench dataset validate", flag.ContinueOnError)
	datasetRoot := fs.String("dataset-root", "./datasets", "Dataset root directory")
	langs := fs.String("langs", "", "Comma-separated languages to include")
	classFilter := fs.String("class", "", "Dataset class filter")
	scenarioFilter := fs.String("scenario", "", "Dataset scenario filter")
	strict := fs.Bool("strict", false, "Fail on validation errors")
	jsonOut := fs.String("json", "", "Write validation report JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	report := svc.ValidateReadiness(dataset.ValidateOptions{
		DatasetRoot: *datasetRoot,
		Languages:   parseCommaList(*langs),
		Classes:     parseCommaList(*classFilter),
		Scenario:    *scenarioFilter,
		Strict:      *strict,
	})
	if *jsonOut != "" {
		if err := contracts.WriteJSON(*jsonOut, report); err != nil {
			return err
		}
	}
	fmt.Printf("Dataset validation: %s\n", mapBool(report.OK, "OK", "FAILED"))
	fmt.Printf("  samples: %d | errors: %d | warnings: %d\n", report.Total, len(report.Errors), len(report.Warnings))
	for _, count := range report.Counts {
		fmt.Printf("  %s/%s/%s: %d\n", count.Language, count.Class, count.Scenario, count.Count)
	}
	for _, issue := range append(report.Errors, firstIssues(report.Warnings, 10)...) {
		fmt.Printf("  [%s] %s %s %s\n", issue.Severity, issue.Code, issue.Language, issue.Path)
	}
	if *strict && !report.OK {
		return fmt.Errorf("dataset validation failed with %d errors", len(report.Errors))
	}
	return nil
}

type doctorToolRow struct {
	Name    string `json:"name"`
	Command string `json:"command,omitempty"`
	Version string `json:"version,omitempty"`
	OK      bool   `json:"ok"`
	Error   string `json:"error,omitempty"`
}

type doctorCanaryRow struct {
	Language      string   `json:"language"`
	CompilePass   bool     `json:"compile_pass"`
	TestPass      bool     `json:"test_pass"`
	LineCoverage  *float64 `json:"line_coverage,omitempty"`
	MutationScore *float64 `json:"mutation_score,omitempty"`
	Error         string   `json:"error,omitempty"`
}

type doctorReport struct {
	OK       bool              `json:"ok"`
	Tools    []doctorToolRow   `json:"tools"`
	Canaries []doctorCanaryRow `json:"canaries"`
}

func runDoctor(args []string) error {
	fs := flag.NewFlagSet("utbench doctor", flag.ContinueOnError)
	langs := fs.String("langs", "python,go,java,cpp", "Comma-separated languages to check")
	mutationEnabled := fs.Bool("mutation-enabled", true, "Enable mutation canary checks")
	mutationTimeout := fs.Int("mutation-timeout", 120, "Mutation timeout (seconds)")
	testTimeout := fs.Int("test-timeout", 60, "Canary test timeout (seconds)")
	jsonOut := fs.String("json", "", "Write doctor report JSON")
	verbose := fs.Bool("v", false, "Verbose output")
	if err := fs.Parse(args); err != nil {
		return err
	}

	langList := parseCommaList(*langs)
	if len(langList) == 0 {
		langList = contracts.SupportedLanguages
	}

	report := doctorReport{}
	report.Tools = checkDoctorTools(langList, *mutationEnabled)
	canaries, canaryErr := runDoctorCanaries(langList, *mutationEnabled, *mutationTimeout, *testTimeout, *verbose)
	report.Canaries = canaries
	report.OK = canaryErr == nil
	for _, tool := range report.Tools {
		if !tool.OK {
			report.OK = false
		}
	}
	for _, row := range report.Canaries {
		if row.Error != "" || !row.CompilePass || !row.TestPass || row.LineCoverage == nil || (*mutationEnabled && row.MutationScore == nil) {
			report.OK = false
		}
	}

	if *jsonOut != "" {
		if err := contracts.WriteJSON(*jsonOut, report); err != nil {
			return err
		}
	}
	fmt.Printf("Doctor: %s\n", mapBool(report.OK, "OK", "FAILED"))
	for _, tool := range report.Tools {
		if tool.OK {
			fmt.Printf("  [tool ok] %s: %s\n", tool.Name, firstLine(tool.Version))
		} else {
			fmt.Printf("  [tool fail] %s: %s\n", tool.Name, tool.Error)
		}
	}
	for _, row := range report.Canaries {
		fmt.Printf("  [canary] %s compile=%v test=%v coverage=%v mutation=%v %s\n",
			row.Language, row.CompilePass, row.TestPass, row.LineCoverage != nil, row.MutationScore != nil, row.Error)
	}
	if canaryErr != nil {
		return canaryErr
	}
	if !report.OK {
		return fmt.Errorf("doctor checks failed")
	}
	return nil
}

func checkDoctorTools(langs []string, mutationEnabled bool) []doctorToolRow {
	need := map[string]bool{}
	for _, lang := range langs {
		need[strings.ToLower(strings.TrimSpace(lang))] = true
	}
	var rows []doctorToolRow
	if need["python"] {
		py := findPythonCommand()
		rows = append(rows, versionRow("python", py, "--version"))
		rows = append(rows, pythonModuleVersionRow("pytest", py, "pytest", "--version"))
		rows = append(rows, pythonModuleVersionRow("coverage", py, "coverage", "--version"))
		if mutationEnabled {
			rows = append(rows, pythonModuleVersionRow("mutmut", py, "mutmut", "--version"))
		}
	}
	if need["go"] {
		rows = append(rows, versionRow("go", "go", "version"))
		if mutationEnabled {
			rows = append(rows, versionRow("go-mutesting", findGoMutestingCommand(), "--help"))
		}
	}
	if need["java"] {
		rows = append(rows, versionRow("java", "java", "-version"))
		rows = append(rows, versionRow("mvn", "mvn", "-version"))
		if mutationEnabled {
			rows = append(rows, configuredToolRow("pitest", "PITest Maven plugin 1.19.6 configured by evaluator"))
		}
	}
	if need["cpp"] {
		rows = append(rows, versionRow("cmake", "cmake", "--version"))
		rows = append(rows, versionRow("gcov", "gcov", "--version"))
		if mutationEnabled {
			rows = append(rows, versionRow("mull", findMullCommand(), "--version"))
		}
	}
	return rows
}

func runDoctorCanaries(langs []string, mutationEnabled bool, mutationTimeout, testTimeout int, verbose bool) ([]doctorCanaryRow, error) {
	tmpRoot, err := os.MkdirTemp("", "utbench_doctor_")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpRoot)

	cases, err := writeDoctorCanaryFiles(tmpRoot, langs)
	if err != nil {
		return nil, err
	}
	manifest := contracts.GeneratedManifest{
		SchemaVersion: contracts.SchemaVersion,
		RunID:         "doctor",
		CreatedAtUTC:  time.Now().UTC(),
		Cases:         cases,
	}
	manifestPath := filepath.Join(tmpRoot, "generated_manifest.json")
	if err := contracts.WriteJSON(manifestPath, manifest); err != nil {
		return nil, err
	}

	logger := obs.NewLogger(verbose, "")
	evaluatorSvc := evaluator.NewService(logger)
	output, err := evaluatorSvc.Evaluate(context.Background(), contracts.RunSpec{
		RunID:           "doctor",
		OutputRoot:      filepath.Join(tmpRoot, "artifacts"),
		MutationEnabled: mutationEnabled,
		MutationTimeout: mutationTimeout,
		MutationPolicy:  "warn",
		TestTimeout:     testTimeout,
		Workers:         1,
	}, manifestPath)
	if err != nil {
		return nil, err
	}
	rows := make([]doctorCanaryRow, 0, len(output.Result.Results))
	for _, r := range output.Result.Results {
		testPass := r.TestPass != nil && *r.TestPass
		errMsg := firstNonEmptyMain(r.CompileError, r.TestError, r.CoverageError, r.MutationError)
		rows = append(rows, doctorCanaryRow{
			Language:      r.Language,
			CompilePass:   r.CompilePass,
			TestPass:      testPass,
			LineCoverage:  r.LineCoverage,
			MutationScore: r.MutationScore,
			Error:         errMsg,
		})
	}
	return rows, nil
}

func writeDoctorCanaryFiles(root string, langs []string) ([]contracts.GeneratedCase, error) {
	var cases []contracts.GeneratedCase
	for _, lang := range langs {
		lang = strings.ToLower(strings.TrimSpace(lang))
		if lang == "" {
			continue
		}
		source, test, sourceName, testName := doctorCanarySource(lang)
		if source == "" {
			continue
		}
		langDir := filepath.Join(root, lang)
		if err := os.MkdirAll(langDir, 0o755); err != nil {
			return nil, err
		}
		sourcePath := filepath.Join(langDir, sourceName)
		testPath := filepath.Join(langDir, testName)
		if err := os.WriteFile(sourcePath, []byte(source), 0o644); err != nil {
			return nil, err
		}
		if err := os.WriteFile(testPath, []byte(test), 0o644); err != nil {
			return nil, err
		}
		cases = append(cases, contracts.GeneratedCase{
			Model:             "doctor",
			Language:          lang,
			SampleID:          "doctor_" + lang,
			SamplePath:        sourcePath,
			GeneratedTestPath: testPath,
			GeneratedAtUTC:    time.Now().UTC(),
			Success:           true,
		})
	}
	return cases, nil
}

func doctorCanarySource(lang string) (source, test, sourceName, testName string) {
	switch lang {
	case "python":
		return "def add(a, b):\n    return a + b\n\ndef is_positive(x):\n    return x > 0\n",
			"from doctor_python import add, is_positive\n\n\ndef test_add():\n    assert add(2, 3) == 5\n\n\ndef test_is_positive():\n    assert is_positive(1) is True\n",
			"doctor_python.py", "doctor_python_test.py"
	case "go":
		return "package main\n\nfunc Add(a int, b int) int { return a + b }\nfunc IsPositive(x int) bool { return x > 0 }\n",
			"package main\n\nimport \"testing\"\n\nfunc TestAdd(t *testing.T) {\n\tif Add(2, 3) != 5 { t.Fatalf(\"unexpected add\") }\n}\n\nfunc TestIsPositive(t *testing.T) {\n\tif !IsPositive(1) { t.Fatalf(\"expected positive\") }\n}\n",
			"doctor_go.go", "doctor_go_test.go"
	case "java":
		return "public class DoctorJava {\n    public int add(int a, int b) { return a + b; }\n    public boolean isPositive(int x) { return x > 0; }\n}\n",
			"import org.junit.Test;\nimport static org.junit.Assert.*;\n\npublic class DoctorJavaTest {\n    @Test public void testAdd() { assertEquals(5, new DoctorJava().add(2, 3)); }\n    @Test public void testIsPositive() { assertTrue(new DoctorJava().isPositive(1)); }\n}\n",
			"DoctorJava.java", "DoctorJavaTest.java"
	case "cpp":
		return "int add(int a, int b) { return a + b; }\nbool is_positive(int x) { return x > 0; }\n",
			"#include <gtest/gtest.h>\n\nTEST(DoctorCpp, Add) { EXPECT_EQ(add(2, 3), 5); }\nTEST(DoctorCpp, IsPositive) { EXPECT_TRUE(is_positive(1)); }\n",
			"doctor_cpp.cpp", "doctor_cpp_test.cpp"
	default:
		return "", "", "", ""
	}
}

func versionRow(name, command string, args ...string) doctorToolRow {
	row := doctorToolRow{Name: name, Command: command}
	if strings.TrimSpace(command) == "" {
		row.Error = "command not found"
		return row
	}
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	row.Version = strings.TrimSpace(string(out))
	if err != nil {
		row.Error = err.Error()
		return row
	}
	row.OK = true
	return row
}

func pythonModuleVersionRow(name, py, module string, args ...string) doctorToolRow {
	fullArgs := append([]string{"-m", module}, args...)
	return versionRow(name, py, fullArgs...)
}

func configuredToolRow(name, version string) doctorToolRow {
	return doctorToolRow{Name: name, Version: version, OK: true}
}

func findPythonCommand() string {
	for _, candidate := range []string{"/opt/venv/bin/python", "/opt/venv/bin/python3", "python3", "python"} {
		if path, err := exec.LookPath(candidate); err == nil {
			if commandRuns(path, "--version") {
				return path
			}
		}
		if _, err := os.Stat(candidate); err == nil {
			if commandRuns(candidate, "--version") {
				return candidate
			}
		}
	}
	return ""
}

func findGoMutestingCommand() string {
	for _, candidate := range []string{"go-mutesting", filepath.Join(os.Getenv("HOME"), "go", "bin", "go-mutesting"), "/root/go/bin/go-mutesting"} {
		if path, err := exec.LookPath(candidate); err == nil {
			return path
		}
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return ""
}

func findMullCommand() string {
	for _, candidate := range []string{"mull-runner-19", "mull-runner-18", "mull-runner"} {
		if path, err := exec.LookPath(candidate); err == nil {
			return path
		}
	}
	return ""
}

func commandRuns(command string, args ...string) bool {
	if strings.TrimSpace(command) == "" {
		return false
	}
	return exec.Command(command, args...).Run() == nil
}

func mapBool(ok bool, trueVal, falseVal string) string {
	if ok {
		return trueVal
	}
	return falseVal
}

func firstIssues(items []dataset.ValidationIssue, max int) []dataset.ValidationIssue {
	if len(items) <= max {
		return items
	}
	return items[:max]
}

func firstLine(v string) string {
	v = strings.TrimSpace(v)
	if idx := strings.IndexAny(v, "\r\n"); idx >= 0 {
		return v[:idx]
	}
	return v
}

func firstNonEmptyMain(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
