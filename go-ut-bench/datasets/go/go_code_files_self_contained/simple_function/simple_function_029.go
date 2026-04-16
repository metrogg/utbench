package main

import (
	"fmt"
	"os"
)

func sanityCheck() bool {
	//sanity check
	if token := os.Getenv(envBinCheck); token != "" {
		fmt.Fprint(os.Stdout, token)
		return true
	}
	//legacy sanity check using old env var
	if token := os.Getenv(envBinCheckLegacy); token != "" {
		fmt.Fprint(os.Stdout, token)
		return true
	}
	return false
}
