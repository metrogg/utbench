package main

import (
	"fmt"
)

func StartKernel() {
	client, err := getDockerClient()
	if err != nil {
		panic(fmt.Sprintf("can't create Docker client: %v", err))
	}
	go kernel(client, done)
}
