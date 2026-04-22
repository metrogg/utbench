package main

import (
    "sort"
)

// Given list of integers, return list in strange order.
// Strange sorting, is when you start with the minimum value,
// then maximum of the remaining integers, then minimum and so on.
// 
// Examples:
// StrangeSortList([1, 2, 3, 4]) == [1, 4, 2, 3]
// StrangeSortList([5, 5, 5, 5]) == [5, 5, 5, 5]
// StrangeSortList([]) == []
func StrangeSortList(lst []int) []int {
    sort.Ints(lst)
	result := make([]int, 0)
	for i := 0; i < len(lst)/2; i++ {
		result = append(result, lst[i])
		result = append(result, lst[len(lst)-i-1])
	}
	if len(lst)%2 != 0 {
		result = append(result, lst[len(lst)/2])
	}
	return result
}
