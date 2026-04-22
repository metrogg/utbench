package main

import (
    "math"
)

// Given a list of numbers, return the sum of squares of the numbers
// in the list that are odd. Ignore numbers that are negative or not integers.
// 
// DoubleTheDifference([1, 3, 2, 0]) == 1 + 9 + 0 + 0 = 10
// DoubleTheDifference([-1, -2, 0]) == 0
// DoubleTheDifference([9, -2]) == 81
// DoubleTheDifference([0]) == 0
// 
// If the input list is empty, return 0.
func DoubleTheDifference(lst []float64) int {
    sum := 0
    for _, i := range lst {
        if i > 0 && math.Mod(i, 2) != 0 && i == float64(int(i)) {
            sum += int(math.Pow(i, 2))
        }
    }
    return sum
}
