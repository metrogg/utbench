package main

// Write a function that returns true if the object q will fly, and false otherwise.
// The object q will fly if it's balanced (it is a palindromic list) and the sum of its elements is less than or equal the maximum possible weight w.
// 
// Example:
// WillItFly([1, 2], 5) ➞ false
// 1+2 is less than the maximum possible weight, but it's unbalanced.
// 
// WillItFly([3, 2, 3], 1) ➞ false
// it's balanced, but 3+2+3 is more than the maximum possible weight.
// 
// WillItFly([3, 2, 3], 9) ➞ true
// 3+2+3 is less than the maximum possible weight, and it's balanced.
// 
// WillItFly([3], 5) ➞ true
// 3 is less than the maximum possible weight, and it's balanced.
func WillItFly(q []int,w int) bool {
    sum := 0
	for i := 0; i < len(q); i++ {
		sum += q[i]
	}
	if sum <= w && isPalindrome(q) {
		return true
	}
	return false
}

func isPalindrome(arr []int) bool {
	for i := 0; i < (len(arr) / 2); i++ {
		if arr[i] != arr[len(arr) - i - 1] {
			return false
		}
	}
	return true
}
