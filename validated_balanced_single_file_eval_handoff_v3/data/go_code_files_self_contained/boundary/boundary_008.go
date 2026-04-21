package main

import (
    "math"
    "strconv"
    "strings"
)

// Create a function that takes a value (string) representing a number
// and returns the closest integer to it. If the number is equidistant
// from two integers, round it away from zero.
// 
// Examples
// >>> ClosestInteger("10")
// 10
// >>> ClosestInteger("15.3")
// 15
// 
// Note:
// Rounding away from zero means that if the given number is equidistant
// from two integers, the one you should return is the one that is the
// farthest from zero. For example ClosestInteger("14.5") should
// return 15 and ClosestInteger("-14.5") should return -15.
func ClosestInteger(value string) int {
    if strings.Count(value, ".") == 1 {
        // remove trailing zeros
        for value[len(value)-1] == '0' {
            value = value[:len(value)-1]
        }
    }
    var res float64
    num, _ := strconv.ParseFloat(value, 64)
    if len(value) >= 2 && value[len(value)-2:] == ".5" {
        if num > 0 {
            res = math.Ceil(num)
        } else {
            res = math.Floor(num)
        }
    } else if len(value) > 0 {
        res = math.Round(num)
    } else {
        res = 0
    }

    return int(res)
}
