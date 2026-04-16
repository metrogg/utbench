package main

import (
	"errors"
	"os"
	"path/filepath"
)

func makeGvisorDirs() error {
	// Make /run/containerd/runsc to hold logs
	fp := filepath.Join(nodeDir, "run/containerd/runsc")
	if err := os.MkdirAll(fp, 0755); err != nil {
		return errors.Wrap(err, "creating runsc dir")
	}

	// Make /usr/local/bin to store the runsc binary
	fp = filepath.Join(nodeDir, "usr/local/bin")
	if err := os.MkdirAll(fp, 0755); err != nil {
		return errors.Wrap(err, "creating usr/local/bin dir")
	}

	// Make /tmp/runsc to also hold logs
	fp = filepath.Join(nodeDir, "tmp/runsc")
	if err := os.MkdirAll(fp, 0755); err != nil {
		return errors.Wrap(err, "creating runsc logs dir")
	}

	return nil
}
