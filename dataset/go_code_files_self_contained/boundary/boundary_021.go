package main

import (
	"io/ioutil"
	"strconv"
	"strings"
)

func GetMaxThreads() (int, error) {
	sysMaxThreadsStr, err := ioutil.ReadFile("/proc/sys/kernel/threads-max")
	if err != nil {
		return 0, err
	}
	sysMaxThreads, err := strconv.Atoi(strings.TrimSpace(string(sysMaxThreadsStr)))
	if err != nil {
		return 0, err
	}
	return sysMaxThreads, nil
}
