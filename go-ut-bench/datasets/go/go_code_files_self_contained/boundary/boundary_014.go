package main

// The Fib4 number sequence is a sequence similar to the Fibbonacci sequnece that's defined as follows:
// Fib4(0) -> 0
// Fib4(1) -> 0
// Fib4(2) -> 2
// Fib4(3) -> 0
// Fib4(n) -> Fib4(n-1) + Fib4(n-2) + Fib4(n-3) + Fib4(n-4).
// Please write a function to efficiently compute the n-th element of the Fib4 number sequence.  Do not use recursion.
// >>> Fib4(5)
// 4
// >>> Fib4(6)
// 8
// >>> Fib4(7)
// 14
func Fib4(n int) int {
    switch n {
	case 0:
		return 0
	case 1:
		return 0
	case 2:
		return 2
	case 3:
		return 0
	default:
		return Fib4(n-1) + Fib4(n-2) + Fib4(n-3) + Fib4(n-4)
	}
}
