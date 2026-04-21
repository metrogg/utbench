package main

import (
	"time"
)

func Fib() func() time.Duration {
	a, b := 0, 1
	return func() time.Duration {
		a, b = b, a+b
		return time.Duration(a*10) * time.Millisecond
	}
}
