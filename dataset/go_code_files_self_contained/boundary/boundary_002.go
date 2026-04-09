package main

import (
	"errors"
	"io/ioutil"
	"strings"
)

func getOSDistro() (string, error) {
	f := "/etc/lsb-release"
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read %q", f)
	}
	content := string(b)
	switch {
	case strings.Contains(content, "Ubuntu"):
		return "ubuntu", nil
	case strings.Contains(content, "Chrome OS"):
		return "cos", nil
	case strings.Contains(content, "CoreOS"):
		return "coreos", nil
	default:
		return "", errors.Errorf("failed to get OS distro: %s", content)
	}
}
