package main

import (
    "sort"
)

// Return sorted Unique elements in a list
// >>> Unique([5, 3, 5, 2, 3, 3, 9, 0, 123])
// [0, 2, 3, 5, 9, 123]
func Unique(l []int) []int {
    set := make(map[int]interface{})
	for _, i := range l {
		set[i]=nil
	}
	l = make([]int,0)
	for i, _ := range set {
		l = append(l, i)
	}
	sort.Ints(l)
	return l
}
