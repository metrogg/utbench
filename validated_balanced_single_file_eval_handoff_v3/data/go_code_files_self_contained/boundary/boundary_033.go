package main

// Given a list of numbers, return whether or not they are sorted
// in ascending order. If list has more than 1 duplicate of the same
// number, return false. Assume no negative numbers and only integers.
// 
// Examples
// IsSorted([5]) ➞ true
// IsSorted([1, 2, 3, 4, 5]) ➞ true
// IsSorted([1, 3, 2, 4, 5]) ➞ false
// IsSorted([1, 2, 3, 4, 5, 6]) ➞ true
// IsSorted([1, 2, 3, 4, 5, 6, 7]) ➞ true
// IsSorted([1, 3, 2, 4, 5, 6, 7]) ➞ false
// IsSorted([1, 2, 2, 3, 3, 4]) ➞ true
// IsSorted([1, 2, 2, 2, 3, 4]) ➞ false
func IsSorted(lst []int) bool {
    count_digit := make(map[int]int)
    for _, i := range lst {
        count_digit[i] = 0
    }
    for _, i := range lst {
        count_digit[i]++
    }
    for _, i := range lst {
        if count_digit[i] > 2 {
            return false
        }
    }
    for i := 1;i < len(lst);i++ {
        if lst[i-1] > lst[i] {
            return false
        }
    }
    return true
}
