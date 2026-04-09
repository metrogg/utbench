package main

import (
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	go client_task()
	go client_task()
	go client_task()
	go server_task()

	//  Run for 5 seconds then quit
	time.Sleep(5 * time.Second)
}
