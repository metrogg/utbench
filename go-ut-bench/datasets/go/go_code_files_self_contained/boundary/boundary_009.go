package main

import (
    "math"
)

// This function will take a list of integers. For all entries in the list, the function shall square the integer entry if its index is a
// multiple of 3 and will cube the integer entry if its index is a multiple of 4 and not a multiple of 3. The function will not
// change the entries in the list whose indexes are not a multiple of 3 or 4. The function shall then return the sum of all entries.
// 
// Examples:
// For lst = [1,2,3] the output should be 6
// For lst = []  the output should be 0
// For lst = [-1,-5,2,-1,-5]  the output should be -126
func SumSquares(lst []int) int {
    result := make([]int, 0)
    for i := 0;i < len(lst);i++ {
        switch {
        case i %3 == 0:
            result = append(result, int(math.Pow(float64(lst[i]), 2)))
        case i % 4 == 0 && i%3 != 0:
            result = append(result, int(math.Pow(float64(lst[i]), 3)))
        default:
            result = append(result, lst[i])
        }
    }
    sum := 0
    for _, x := range result {
        sum += x
    }
    return sum
}
