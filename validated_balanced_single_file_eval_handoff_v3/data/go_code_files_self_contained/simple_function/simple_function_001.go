package main

// Write a function that accepts two lists of strings and returns the list that has
// total number of chars in the all strings of the list less than the other list.
// 
// if the two lists have the same number of chars, return the first list.
// 
// Examples
// TotalMatch([], []) ➞ []
// TotalMatch(['hi', 'admin'], ['hI', 'Hi']) ➞ ['hI', 'Hi']
// TotalMatch(['hi', 'admin'], ['hi', 'hi', 'admin', 'project']) ➞ ['hi', 'admin']
// TotalMatch(['hi', 'admin'], ['hI', 'hi', 'hi']) ➞ ['hI', 'hi', 'hi']
// TotalMatch(['4'], ['1', '2', '3', '4', '5']) ➞ ['4']
func TotalMatch(lst1 []string,lst2 []string) []string {
    var numchar1 = 0
	var numchar2 = 0
	for _, item := range lst1 {
		numchar1 += len(item)
	}
	for _, item := range lst2 {
		numchar2 += len(item)
	}
	if numchar1 <= numchar2 {
		return lst1
	} else {
		return lst2
	}
}
