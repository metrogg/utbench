package main

import (
    "sort"
    "strconv"
)

// Given a list of positive integers x. return a sorted list of all
// elements that hasn't any even digit.
// 
// Note: Returned list should be sorted in increasing order.
// 
// For example:
// >>> UniqueDigits([15, 33, 1422, 1])
// [1, 15, 33]
// >>> UniqueDigits([152, 323, 1422, 10])
// []
func UniqueDigits(x []int) []int {
    odd_digit_elements := make([]int, 0)
    OUTER:
    for _, i := range x {
        for _, c := range strconv.Itoa(i) {
            if (c - '0') % 2 == 0 {
                continue OUTER
            }
        }
            odd_digit_elements = append(odd_digit_elements, i)
    }
    sort.Slice(odd_digit_elements, func(i, j int) bool {
        return odd_digit_elements[i] < odd_digit_elements[j]
    })
    return odd_digit_elements
}
