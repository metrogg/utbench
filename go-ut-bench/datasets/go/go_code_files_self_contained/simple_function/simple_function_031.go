package main

import (
	"fmt"
)

func MyIP() (string, error) {

	maxIface := 5
	var err error
	for i := 0; i < maxIface; i++ {
		var ip string
		ip, err = IPByInterface(fmt.Sprintf("eth%d", i))
		if err == nil {
			return ip, nil
		}
	}

	return "0.0.0.0", err
}
