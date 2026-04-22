package main

// brackets is a string of "<" and ">".
// return true if every opening bracket has a corresponding closing bracket.
// 
// >>> CorrectBracketing("<")
// false
// >>> CorrectBracketing("<>")
// true
// >>> CorrectBracketing("<<><>>")
// true
// >>> CorrectBracketing("><<>")
// false
func CorrectBracketing(brackets string) bool {
    l := len(brackets)
	count := 0
	for index := 0; index < l; index++ {
		if brackets[index] == '<' {
			count++
		} else if brackets[index] == '>' {
			count--
		}
		if count < 0 {
			return false
		}
	}
    if count == 0 {
        return true
    } else {
        return false
    }
}
