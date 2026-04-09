package main

import (
	"os"
)

func supportsEditing() bool {
	term := os.Getenv("TERM")

	for _, t := range unsupported {
		if t == term {
			return false
		}
	}

	return true
}
