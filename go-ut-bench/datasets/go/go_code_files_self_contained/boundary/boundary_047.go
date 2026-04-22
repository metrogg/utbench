package main

import (
    "strconv"
)

// Circular shift the digits of the integer x, shift the digits right by shift
// and return the result as a string.
// If shift > number of digits, return digits reversed.
// >>> CircularShift(12, 1)
// "21"
// >>> CircularShift(12, 2)
// "12"
func CircularShift(x int, shift int) string {
    s := strconv.Itoa(x)
	if shift > len(s) {
		runes := make([]rune, 0)
		for i := len(s)-1; i >= 0; i-- {
			runes = append(runes, rune(s[i]))
		}
		return string(runes)
	}else{
		return s[len(s)-shift:]+s[:len(s)-shift]
	}
}
