package main

import (
	"bytes"
	"io/ioutil"
	"os"
)

func DetectDistro() string {
	if b, e := ioutil.ReadFile("/etc/lsb-release"); e == nil && bytes.Contains(b, []byte("Ubuntu")) {
		return "ubuntu"
	} else if isRedhat() {
		return "redhat"
	} else if _, err := os.Stat("/etc/alpine-release"); err == nil {
		return "alpine"
	} else if _, err := os.Stat("/etc/arch-release"); err == nil {
		return "arch"
	} else if _, err := os.Stat("/etc/debian_version"); err == nil {
		return "debian"
	}
	return ""
}
