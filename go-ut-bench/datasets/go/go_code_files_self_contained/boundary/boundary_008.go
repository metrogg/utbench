package main

import (
	"io/ioutil"
)

func getMemoryCapacity() (uint64, error) {
	out, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	return parseCapacity(out, memoryCapacityRegexp)
}
