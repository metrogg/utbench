package main

import (
    "strconv"
)

// Given an integer. return a tuple that has the number of even and odd digits respectively.
// 
// Example:
// EvenOddCount(-12) ==> (1, 1)
// EvenOddCount(123) ==> (1, 2)
func EvenOddCount(num int) [2]int {
    even_count := 0
    odd_count := 0
    if num < 0 {
        num = -num
    }
    for _, r := range strconv.Itoa(num) {
        if r&1==0 {
            even_count++
        } else {
            odd_count++
        }
    }
    return [2]int{even_count, odd_count}
}
