package main

// You are given two intervals,
// where each interval is a pair of integers. For example, interval = (start, end) = (1, 2).
// The given intervals are closed which means that the interval (start, end)
// includes both start and end.
// For each given interval, it is assumed that its start is less or equal its end.
// Your task is to determine whether the length of Intersection of these two
// intervals is a prime number.
// Example, the Intersection of the intervals (1, 3), (2, 4) is (2, 3)
// which its length is 1, which not a prime number.
// If the length of the Intersection is a prime number, return "YES",
// otherwise, return "NO".
// If the two intervals don't intersect, return "NO".
// 
// 
// [input/output] samples:
// Intersection((1, 2), (2, 3)) ==> "NO"
// Intersection((-1, 1), (0, 4)) ==> "NO"
// Intersection((-3, -1), (-5, 5)) ==> "YES"
func Intersection(interval1 [2]int, interval2 [2]int) string {
    is_prime := func(num int) bool {
        if num == 1 || num == 0 {
            return false
        }
        if num == 2 {
            return true
        }
        for i := 2;i < num;i++ {
            if num%i == 0 {
                return false
            }
        }
        return true
    }
    l := interval1[0]
    if interval2[0] > l {
        l = interval2[0]
    }
    r := interval1[1]
    if interval2[1] < r {
        r = interval2[1]
    }
    length := r - l
    if length > 0 && is_prime(length) {
        return "YES"
    }
    return "NO"
}
