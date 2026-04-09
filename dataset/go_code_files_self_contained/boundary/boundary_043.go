package main

import (
	"os"
)

func Path() string {
	p, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return p
}
