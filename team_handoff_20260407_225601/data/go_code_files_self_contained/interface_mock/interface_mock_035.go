package main

import (
	"os"
)

func HasSupport() bool {
	t := os.HostOS()
	for _, v := range osSupport {
		if v == t {
			return true
		}
	}
	return false
}
