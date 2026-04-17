package main

import (
	"strconv"
)

func MaxBrokerClients() int {
	c, err := strconv.Atoi(maxBrokerClients)
	if err != nil {
		return 50000
	}

	return c
}
