package main

// From a list of integers, remove all elements that occur more than once.
// Keep order of elements left the same as in the input.
// >>> RemoveDuplicates([1, 2, 3, 2, 4])
// [1, 3, 4]
func RemoveDuplicates(numbers []int) []int {
    c := make(map[int] int)
	for _, number := range numbers {
		if i, ok := c[number]; ok {
			c[number] = i + 1
		} else {
			c[number] = 1
		}
	}
	result := make([]int, 0)
	for _, number := range numbers {
		if c[number] <= 1 {
			result = append(result, number)
		}
	}
	return result
}
