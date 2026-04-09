package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func getMemory() (string, error) {
	bytes, err := ioutil.ReadFile(meminfo)
	if err != nil {
		return "", err
	}

	lines := string(bytes)

	for _, line := range strings.Split(lines, "\n") {
		if !strings.HasPrefix(line, "MemTotal") {
			continue
		}

		expectedFields := 2

		fields := strings.Split(line, ":")
		count := len(fields)

		if count != expectedFields {
			return "", fmt.Errorf("expected %d fields, got %d in line %q", expectedFields, count, line)
		}

		if fields[1] == "" {
			return "", fmt.Errorf("cannot determine total memory from line %q", line)
		}

		memTotal := strings.TrimSpace(fields[1])

		return memTotal, nil
	}

	return "", fmt.Errorf("no lines in file %q", meminfo)
}
