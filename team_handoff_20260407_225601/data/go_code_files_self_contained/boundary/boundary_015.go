package main

import (
	"regexp"
)

func init() {
	var err error
	validTagRe, err = regexp.Compile(`^([a-z][a-z0-9]*(\-[a-z0-9]+)*)(:([a-z0-9]+(\-[a-z0-9]+)*))*$`)
	if err != nil {
		panic(err)
	}
}
