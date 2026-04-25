package main

import (
	"path/filepath"
	"testing"
)

func TestDoctorCanarySourceCoversSupportedLanguages(t *testing.T) {
	for _, lang := range []string{"python", "go", "java", "cpp"} {
		source, test, sourceName, testName := doctorCanarySource(lang)
		if source == "" || test == "" || sourceName == "" || testName == "" {
			t.Fatalf("missing canary for %s", lang)
		}
	}
}

func TestWriteDoctorCanaryFilesBuildsCases(t *testing.T) {
	root := t.TempDir()
	cases, err := writeDoctorCanaryFiles(root, []string{"python", "go"})
	if err != nil {
		t.Fatal(err)
	}
	if len(cases) != 2 {
		t.Fatalf("expected 2 canary cases, got %d", len(cases))
	}
	for _, c := range cases {
		if c.Model != "doctor" || !c.Success {
			t.Fatalf("unexpected canary case: %+v", c)
		}
		if !filepath.IsAbs(c.SamplePath) || !filepath.IsAbs(c.GeneratedTestPath) {
			t.Fatalf("expected absolute canary paths: %+v", c)
		}
	}
}
