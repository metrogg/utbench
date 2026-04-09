package main

import (
	"strconv"
	"strings"
)

func GetVersionNumber() (int64, int) {
	ver := strings.Replace(GetVersionString(), "-", ".", -1)
	parts := strings.Split(ver, ".")
	multiplier := int64(100 * 100) // major, minor, patch
	version := int64(0)
	buildNum := int(0)
	for i, subVerString := range parts {
		if i == 3 {
			buildNum, _ = strconv.Atoi(subVerString)
			break // done
		}
		subVer, err := strconv.Atoi(subVerString)
		if err != nil {
			break // cannot parse
		}
		version += int64(subVer) * multiplier
		multiplier /= 100

	}
	return version, buildNum
}
