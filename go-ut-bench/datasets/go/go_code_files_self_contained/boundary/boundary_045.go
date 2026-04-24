package main

// Given list of numbers (of at least two elements), apply a linear transform to that list,
// such that the smallest number will become 0 and the largest will become 1
// >>> RescaleToUnit([1.0, 2.0, 3.0, 4.0, 5.0])
// [0.0, 0.25, 0.5, 0.75, 1.0]
func RescaleToUnit(numbers []float64) []float64 {
    smallest := numbers[0]
	largest := smallest
	for _, n := range numbers {
		if smallest > n {
			smallest = n
		}
		if largest < n {
			largest = n
		}
	}
	if smallest == largest {
		return numbers
	}
	for i, n := range numbers {
		numbers[i] = (n - smallest) / (largest - smallest)
	}
	return numbers
}
