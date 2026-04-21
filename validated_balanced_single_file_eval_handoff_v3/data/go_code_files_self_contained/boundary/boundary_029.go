package main

// Everyone knows Fibonacci sequence, it was studied deeply by mathematicians in
// the last couple centuries. However, what people don't know is Tribonacci sequence.
// Tribonacci sequence is defined by the recurrence:
// Tri(1) = 3
// Tri(n) = 1 + n / 2, if n is even.
// Tri(n) =  Tri(n - 1) + Tri(n - 2) + Tri(n + 1), if n is odd.
// For example:
// Tri(2) = 1 + (2 / 2) = 2
// Tri(4) = 3
// Tri(3) = Tri(2) + Tri(1) + Tri(4)
// = 2 + 3 + 3 = 8
// You are given a non-negative integer number n, you have to a return a list of the
// first n + 1 numbers of the Tribonacci sequence.
// Examples:
// Tri(3) = [1, 3, 2, 8]
func Tri(n int) []float64 {
    if n == 0 {
        return []float64{1}
    }
    my_tri := []float64{1, 3}
    for i := 2; i < n + 1; i++ {
        if i &1 == 0 {
            my_tri = append(my_tri, float64(i) / 2 + 1)
        } else {
            my_tri = append(my_tri, my_tri[i - 1] + my_tri[i - 2] + (float64(i) + 3) / 2)
        }
    }
    return my_tri
}
