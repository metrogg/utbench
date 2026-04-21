package main

import (
	"time"
)

func otherSource() chan int {
	c := make(chan int, 0)
	go func() {
		defer close(c)

		i := 1
		for {
			c <- i
			i++
			time.Sleep(time.Second)
		}
	}()
	return c
}
