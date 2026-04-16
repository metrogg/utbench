package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

func GetLinuxVersion() (string, error) {
	etc := HostEtc()
	if value, err := getStringFromFile(`DISTRIB_DESCRIPTION="(.*)"`, filepath.Join(etc, "lsb-release")); err == nil {
		return value, nil
	}
	if value, err := getStringFromFile(`PRETTY_NAME="(.*)"`, filepath.Join(etc, "os-release")); err == nil {
		return value, nil
	}
	if value, err := ioutil.ReadFile(filepath.Join(etc, "centos-release")); err == nil {
		return string(value), nil
	}
	if value, err := ioutil.ReadFile(filepath.Join(etc, "redhat-release")); err == nil {
		return string(value), nil
	}
	if value, err := ioutil.ReadFile(filepath.Join(etc, "system-release")); err == nil {
		return string(value), nil
	}
	return "", errors.New("unable to find linux version")
}
