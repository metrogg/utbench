package main

import (
	"io/ioutil"
	"strings"
)

func AppArmorProfile() string {
	contents, err := ioutil.ReadFile("/proc/self/attr/current")
	if err == nil {
		return strings.TrimSpace(string(contents))
	}

	return ""
}
