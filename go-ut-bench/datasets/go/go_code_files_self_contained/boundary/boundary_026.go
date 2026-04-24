package main

// PairsSumToZero takes a list of integers as an input.
// it returns true if there are two distinct elements in the list that
// sum to zero, and false otherwise.
// >>> PairsSumToZero([1, 3, 5, 0])
// false
// >>> PairsSumToZero([1, 3, -2, 1])
// false
// >>> PairsSumToZero([1, 2, 3, 7])
// false
// >>> PairsSumToZero([2, 4, -5, 3, 5, 7])
// true
// >>> PairsSumToZero([1])
// false
func PairsSumToZero(l []int) bool {
    seen := map[int]bool{}
	for i := 0; i < len(l); i++ {
		for j := i + 1; j < len(l); j++ {
			if l[i] + l[j] == 0 {
				if _, ok := seen[l[i]]; !ok {
					seen[l[i]] = true
					return true
				}
				if _, ok := seen[l[j]]; !ok {
					seen[l[j]] = true
					return true
				}
			}
		}
	}
	return false
}
