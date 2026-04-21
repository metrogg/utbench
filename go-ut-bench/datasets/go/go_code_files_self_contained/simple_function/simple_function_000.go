package main

func CalculateFibonacci(n int) int {
	if n < 0 {
		panic("n must be non-negative")
	}
	if n == 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	prev, curr := 0, 1
	for i := 2; i <= n; i++ {
		next := prev + curr
		prev = curr
		curr = next
	}
	return curr
}