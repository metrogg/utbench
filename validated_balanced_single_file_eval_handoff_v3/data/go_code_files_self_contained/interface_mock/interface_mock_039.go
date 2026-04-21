package main

// You are given a list of two strings, both strings consist of open
// parentheses '(' or close parentheses ')' only.
// Your job is to check if it is possible to concatenate the two strings in
// some order, that the resulting string will be good.
// A string S is considered to be good if and only if all parentheses in S
// are balanced. For example: the string '(())()' is good, while the string
// '())' is not.
// Return 'Yes' if there's a way to make a good string, and return 'No' otherwise.
// 
// Examples:
// MatchParens(['()(', ')']) == 'Yes'
// MatchParens([')', ')']) == 'No'
func MatchParens(lst []string) string {
    check := func(s string) bool {
        val := 0
        for _, i := range s {
            if i == '(' {
                val++
            } else {
                val--
            }
            if val < 0 {
                return false
            }
        }
        return val == 0
    }

    S1 := lst[0] + lst[1]
    S2 := lst[1] + lst[0]
    if check(S1) || check(S2) {
        return "Yes"
    }
    return "No"
}
