package main

import (
    "fmt"
    "math"
)

// You are given two positive integers n and m, and your task is to compute the
// average of the integers from n through m (including n and m).
// Round the answer to the nearest integer and convert that to binary.
// If n is greater than m, return -1.
// Example:
// RoundedAvg(1, 5) => "0b11"
// RoundedAvg(7, 5) => -1
// RoundedAvg(10, 20) => "0b1111"
// RoundedAvg(20, 33) => "0b11010"
func RoundedAvg(n, m int) interface{} {
    if m < n {
        return -1
    }
    summation := 0
    for i := n;i < m+1;i++{
        summation += i
    }
    return fmt.Sprintf("0b%b", int(math.Round(float64(summation)/float64(m - n + 1))))
}
