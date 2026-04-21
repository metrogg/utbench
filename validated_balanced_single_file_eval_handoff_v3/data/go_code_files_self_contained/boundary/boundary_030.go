package main

import (
    "strings"
)

// brackets is a string of "(" and ")".
// return true if every opening bracket has a corresponding closing bracket.
// 
// >>> CorrectBracketing("(")
// false
// >>> CorrectBracketing("()")
// true
// >>> CorrectBracketing("(()())")
// true
// >>> CorrectBracketing(")(()")
// false
func CorrectBracketing(brackets string) bool {
    brackets = strings.Replace(brackets, "(", " ( ", -1)
	brackets = strings.Replace(brackets, ")", ") ", -1)
	open := 0
	for _, b := range brackets {
		if b == '(' {
			open++
		} else if b == ')' {
			open--
		}
		if open < 0 {
			return false
		}
	}
	return open == 0
}
