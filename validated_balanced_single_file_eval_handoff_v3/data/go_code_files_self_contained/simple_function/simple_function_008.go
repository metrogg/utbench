package main

import (
	"strconv"
	"strings"
)

// Return the number of times the digit 7 appears in integers less than n which are divisible by 11 or 13.
// >>> FizzBuzz(50)
// 0
// >>> FizzBuzz(78)
// 2
// >>> FizzBuzz(79)
// 3
func FizzBuzz(n int) int {
    ns := make([]int, 0)
	for i := 0; i < n; i++ {
		if i%11 == 0 || i%13 == 0 {
			ns = append(ns, i)
		}
	}
	temp := make([]string, 0)
	for _, i := range ns {
		temp = append(temp, strconv.Itoa(i))
	}
	join := strings.Join(temp, "")
	ans := 0
	for _, c := range join {
		if c == '7' {
			ans++
		}
	}
	return ans
}
