package contracts

type DatasetClass string

const (
	DatasetClassSelfContained DatasetClass = "self_contained"
	DatasetClassModuleLevel   DatasetClass = "module_level"
)

type RunMode string

const (
	RunModeFull        RunMode = "full"
	RunModeIncremental RunMode = "incremental"
)

const SchemaVersion = "v0.1.0"

var SupportedLanguages = []string{"python", "java", "go", "cpp"}
