package main

import (
	"os"
)

func unixLocale() string {
	for _, env := range []string{"LC_ALL", "LC_CTYPE", "LANG"} {
		if locale := os.Getenv(env); locale != "" {
			return locale
		}
	}
	return "US-ASCII"
}
