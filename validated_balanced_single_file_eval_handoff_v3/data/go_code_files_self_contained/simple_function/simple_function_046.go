package main

// Return only positive numbers in the list.
// >>> GetPositive([-1, 2, -4, 5, 6])
// [2, 5, 6]
// >>> GetPositive([5, 3, -5, 2, -3, 3, 9, 0, 123, 1, -10])
// [5, 3, 2, 3, 9, 123, 1]
func GetPositive(l []int) []int {
    res := make([]int, 0)
    for _, x := range l {
        if x > 0 {
            res = append(res, x)
        }
    }
    return res
}
