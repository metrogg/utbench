package main

func Clamp(value, minVal, maxVal int) int {
	if minVal > maxVal {
		panic("minVal must be less than or equal to maxVal")
	}
	if value < minVal {
		return minVal
	}
	if value > maxVal {
		return maxVal
	}
	return value
}