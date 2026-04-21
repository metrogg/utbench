package main

import (
    "math"
)

// You are given a list of numbers.
// You need to return the sum of squared numbers in the given list,
// round each element in the list to the upper int(Ceiling) first.
// Examples:
// For lst = [1,2,3] the output should be 14
// For lst = [1,4,9] the output should be 98
// For lst = [1,3,5,7] the output should be 84
// For lst = [1.4,4.2,0] the output should be 29
// For lst = [-2.4,1,1] the output should be 6
func SumSquares(lst []float64) int {
    squared := 0
    for _, i := range lst {
        squared += int(math.Pow(math.Ceil(i), 2))
    }
    return squared
}
