package main

import (
	"sort"
)

func BuiltinFingerprints() []string {
	fingerprints := make([]string, 0, len(hostFingerprinters))
	for k := range hostFingerprinters {
		fingerprints = append(fingerprints, k)
	}
	sort.Strings(fingerprints)
	for k := range envFingerprinters {
		fingerprints = append(fingerprints, k)
	}
	return fingerprints
}
