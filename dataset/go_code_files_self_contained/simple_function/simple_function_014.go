package main

import (
	"strconv"
)

func EnforceMode() int {
	var enforce int

	enforceS, err := readCon(selinuxEnforcePath())
	if err != nil {
		return -1
	}

	enforce, err = strconv.Atoi(string(enforceS))
	if err != nil {
		return -1
	}
	return enforce
}
