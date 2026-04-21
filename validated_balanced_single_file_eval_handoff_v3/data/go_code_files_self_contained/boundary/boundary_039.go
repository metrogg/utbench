package main

// This function takes two positive numbers x and y and returns the
// biggest even integer number that is in the range [x, y] inclusive. If
// there's no such number, then the function should return -1.
// 
// For example:
// ChooseNum(12, 15) = 14
// ChooseNum(13, 12) = -1
func ChooseNum(x, y int) int {
    if x > y {
        return -1
    }
    if y % 2 == 0 {
        return y
    }
    if x == y {
        return -1
    }
    return y - 1
}
