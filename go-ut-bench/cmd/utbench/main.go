package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
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
	case "ingest":
		err = runIngest(args)
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
  utbench ingest       Ingest evaluation results into SQLite
  utbench dataset      Dataset management (index, manifest, stats)
  utbench doctor       Check evaluator toolchains with canary tests
  utbench web          Launch Web management UI
  utbench tui          Launch interactive TUI interface
  utbench help         Show this help

Run "utbench <command> --help" for more details on a command.`)
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
	signal.Notify(sigChan, os.Interrupt)
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
	imageName := fs.String("docker-image", "utbench:latest", "Docker image for containerized runs")
	projectRoot := fs.String("project-root", ".", "Project root mounted into Docker")
	envFile := fs.String("env-file", "./.env", "Environment file passed to Docker runs")

	if err := fs.Parse(args); err != nil {
		return err
	}

	absProjectRoot, err := filepath.Abs(*projectRoot)
	if err != nil {
		return fmt.Errorf("resolve project root: %w", err)
	}
	absEnvFile := *envFile
	if absEnvFile != "" {
		absEnvFile, err = filepath.Abs(absEnvFile)
		if err != nil {
			return fmt.Errorf("resolve env file: %w", err)
		}
	}

	// Load .env file into process environment so os.Getenv() can read API keys
	if err := web.LoadEnvFile(absEnvFile); err != nil {
		return fmt.Errorf("load env file: %w", err)
	}

	dockerCfg := web.DockerConfig{
		ImageName:   *imageName,
		ProjectRoot: absProjectRoot,
		EnvFile:     absEnvFile,
	}
	mgr := web.NewRunManager(*configPath, *datasetRoot, *outputRoot, *dbPath, dockerCfg)
	bld := web.NewBuildManager(absProjectRoot)
	server := web.NewServer(mgr, bld, *configPath, *outputRoot, dockerCfg)
	return server.Start(*addr)
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
	mutationEnabled := fs.Bool("mutation-enabled", true, "Enable mutation testing")
	mutationTimeout := fs.Int("mutation-timeout", 600, "Mutation timeout (seconds)")
	mutationPolicy := fs.String("mutation-policy", "warn", "Mutation policy (warn, fail)")
	testTimeout := fs.Int("test-timeout", 180, "Test execution timeout (seconds)")
	workers := fs.Int("workers", 16, "Number of concurrent workers (default 16)")
	models := fs.String("models", "", "Comma-separated models")
	langs := fs.String("langs", "", "Comma-separated languages")
	runID := fs.String("run-id", "", "Run ID")

	if err := fs.Parse(args); err != nil {
		return err
	}
	policy, err := normalizeMutationPolicy(*mutationPolicy)
	if err != nil {
		return err
	}

	spec := contracts.RunSpec{
		ConfigPath:      *config,
		OutputRoot:      *outputRoot,
		DatasetRoot:     *datasetRoot,
		DatasetManifest: *datasetManifest,
		DatasetLevel:    *datasetLevel,
		DatasetClasses:  parseCommaList(*datasetClass),
		DatasetScenario: *datasetScenario,
		MaxSamples:      *maxSamples,
		Workers:         *workers,
		Mode:            contracts.RunMode(*mode),
		ResetCheckpoint: *resetCheckpoint,
		DryRun:          *dryRun,
		MutationEnabled: *mutationEnabled,
		MutationTimeout: *mutationTimeout,
		MutationPolicy:  policy,
		TestTimeout:     *testTimeout,
		Models:          parseCommaList(*models),
		Languages:       parseCommaList(*langs),
		RunID:           *runID,
		CreatedAtUTC:    time.Now().UTC(),
	}
	ensureRunID(&spec)

	logDir := filepath.Join(*outputRoot, "runs", spec.RunID, "logs")
	logger := obs.NewLogger(*verbose, logDir)
	ctx, cancel := withSignal(context.Background())
	defer cancel()

	datasetSvc := dataset.NewService()
	runnerSvc := runner.NewService(logger)
	evaluatorSvc := evaluator.NewService(logger)
	reporterSvc := reporter.NewService(logger)

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
	models := fs.String("models", "", "Comma-separated models")
	langs := fs.String("langs", "", "Comma-separated languages")
	runID := fs.String("run-id", "", "Run ID")

	if err := fs.Parse(args); err != nil {
		return err
	}

	spec := contracts.RunSpec{
		ConfigPath:      *config,
		OutputRoot:      *outputRoot,
		DatasetRoot:     *datasetRoot,
		DatasetManifest: *datasetManifest,
		DatasetLevel:    *datasetLevel,
		DatasetClasses:  parseCommaList(*datasetClass),
		DatasetScenario: *datasetScenario,
		MaxSamples:      *maxSamples,
		Mode:            contracts.RunMode(*mode),
		ResetCheckpoint: *resetCheckpoint,
		DryRun:          *dryRun,
		Models:          parseCommaList(*models),
		Languages:       parseCommaList(*langs),
		RunID:           *runID,
		CreatedAtUTC:    time.Now().UTC(),
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

	output, err := runnerSvc.Generate(ctx, spec, samples)
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
	mutationEnabled := fs.Bool("mutation-enabled", true, "Enable mutation testing")
	mutationTimeout := fs.Int("mutation-timeout", 600, "Mutation timeout (seconds)")
	mutationPolicy := fs.String("mutation-policy", "warn", "Mutation policy")
	runID := fs.String("run-id", "", "Run ID")

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
		MutationEnabled: *mutationEnabled,
		MutationTimeout: *mutationTimeout,
		MutationPolicy:  policy,
		RunID:           *runID,
		CreatedAtUTC:    time.Now().UTC(),
	}
	ensureRunID(&spec)

	logDir := filepath.Join(spec.OutputRoot, "runs", spec.RunID, "logs")
	logger := obs.NewLogger(*verbose, logDir)
	ctx, cancel := withSignal(context.Background())
	defer cancel()

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
	reporterSvc := reporter.NewService(logger)
	output, err := reporterSvc.Generate(context.Background(), spec, *evaluationPath)
	if err != nil {
		return fmt.Errorf("report: %w", err)
	}

	fmt.Printf("Report JSON: %s\n", output.ReportJSONPath)
	fmt.Printf("Report HTML: %s\n", output.ReportHTMLPath)
	return nil
}

func runIngest(args []string) error {
	fs := flag.NewFlagSet("utbench ingest", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println("Usage: utbench ingest [flags]")
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}

	verbose := fs.Bool("v", false, "Verbose output")
	dbPath := fs.String("db-path", "./storage/utbench.db", "SQLite database path")
	evaluationPath := fs.String("evaluation", "", "Path to evaluation_result.json (required)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *evaluationPath == "" {
		return fmt.Errorf("--evaluation is required")
	}

	logger := obs.NewLogger(*verbose, "")
	_ = logger

	sqliteStore, err := store.OpenSQLite(*dbPath)
	if err != nil {
		return fmt.Errorf("open sqlite: %w", err)
	}
	defer sqliteStore.Close()

	ctx := context.Background()
	if err := sqliteStore.Init(ctx); err != nil {
		return fmt.Errorf("init sqlite: %w", err)
	}

	resultSet, err := contracts.ReadEvaluationResultSet(*evaluationPath)
	if err != nil {
		return fmt.Errorf("read evaluation: %w", err)
	}

	if err := sqliteStore.IngestEvaluation(ctx, resultSet); err != nil {
		return fmt.Errorf("ingest: %w", err)
	}

	fmt.Printf("Ingested %d results into %s\n", len(resultSet.Results), *dbPath)
	return nil
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
