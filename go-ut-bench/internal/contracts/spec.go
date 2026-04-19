package contracts

import "time"

type RunSpec struct {
	RunID           string       `json:"run_id"`
	Models          []string     `json:"models"`
	Languages       []string     `json:"languages"`
	DatasetClasses  []string     `json:"dataset_classes"`
	DatasetScenario string       `json:"dataset_scenario,omitempty"`
	DatasetLevel    string       `json:"dataset_level,omitempty"`
	DatasetRoot     string       `json:"dataset_root"`
	DatasetManifest string       `json:"dataset_manifest,omitempty"`
	ConfigPath      string       `json:"config_path"`
	Mode            RunMode      `json:"mode"`
	DryRun          bool         `json:"dry_run"`
	ResetCheckpoint bool         `json:"reset_checkpoint"`
	MutationEnabled bool         `json:"mutation_enabled"`
	MutationTimeout int          `json:"mutation_timeout_seconds"`
	MutationPolicy  string       `json:"mutation_error_policy"`
	MaxSamples      int          `json:"max_samples"`
	OutputRoot      string       `json:"output_root"`
	CreatedAtUTC    time.Time    `json:"created_at_utc"`
}

type SampleRef struct {
	ID        string       `json:"id"`
	Language  string       `json:"language"`
	Category  DatasetClass `json:"category"`
	Scenario  string       `json:"scenario"`
	Path      string       `json:"path"`
	SourceMD5 string       `json:"source_md5"`
}

type ModuleLevelMeta struct {
	SampleID       string   `json:"sample_id"`
	ModuleImport   string   `json:"module_import"`
	PackageName    string   `json:"package_name"`
	TargetFile     string   `json:"target_file"`
	WorkspaceRoot  string   `json:"workspace_root"`
	Requirements   []string `json:"requirements,omitempty"`
}
