package main

import (
	"os"
)

func SupportANSI() bool {
	term := os.Getenv("TERM")
	if term == "" {
		return false
	}

	for _, v := range shellsWithoutANSI {
		if v == term {
			return false
		}
	}
	return true
}
