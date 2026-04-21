package main

func SafeDivide(numerator, denominator float64) float64 {
	if denominator == 0 {
		panic("cannot divide by zero")
	}
	if numerator == 0 {
		return 0.0
	}
	result := numerator / denominator
	return result
}

func SafeSqrt(value int) int {
	if value < 0 {
		panic("cannot compute square root of negative number")
	}
	if value == 0 {
		return 0
	}
	i := 0
	for i * i <= value {
		i++
	}
	return i - 1
}