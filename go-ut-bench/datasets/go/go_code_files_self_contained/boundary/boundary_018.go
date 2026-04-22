package main

import (
    "math"
)

// Given an array representing a branch of a tree that has non-negative integer nodes
// your task is to Pluck one of the nodes and return it.
// The Plucked node should be the node with the smallest even value.
// If multiple nodes with the same smallest even value are found return the node that has smallest index.
// 
// The Plucked node should be returned in a list, [ smalest_value, its index ],
// If there are no even values or the given array is empty, return [].
// 
// Example 1:
// Input: [4,2,3]
// Output: [2, 1]
// Explanation: 2 has the smallest even value, and 2 has the smallest index.
// 
// Example 2:
// Input: [1,2,3]
// Output: [2, 1]
// Explanation: 2 has the smallest even value, and 2 has the smallest index.
// 
// Example 3:
// Input: []
// Output: []
// 
// Example 4:
// Input: [5, 0, 3, 0, 4, 2]
// Output: [0, 1]
// Explanation: 0 is the smallest value, but  there are two zeros,
// so we will choose the first zero, which has the smallest index.
// 
// Constraints:
// * 1 <= nodes.length <= 10000
// * 0 <= node.value
func Pluck(arr []int) []int {
    result := make([]int, 0)
	if len(arr) == 0 {
		return result
	}
	evens := make([]int, 0)
	min := math.MaxInt64
	minIndex := 0
	for i, x := range arr {
		if x%2 == 0 {
			evens = append(evens, x)
			if x < min {
				min = x
				minIndex = i
			}
		}
	}
	if len(evens) == 0 {
		return result
	}
	result = []int{min, minIndex}
	return result
}
