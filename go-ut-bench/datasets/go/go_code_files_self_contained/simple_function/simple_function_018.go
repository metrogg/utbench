package main

import (
    "fmt"
)

// Given a list of strings, where each string consists of only digits, return a list.
// Each element i of the output should be "the number of odd elements in the
// string i of the input." where all the i's should be replaced by the number
// of odd digits in the i'th string of the input.
// 
// >>> OddCount(['1234567'])
// ["the number of odd elements 4n the str4ng 4 of the 4nput."]
// >>> OddCount(['3',"11111111"])
// ["the number of odd elements 1n the str1ng 1 of the 1nput.",
// "the number of odd elements 8n the str8ng 8 of the 8nput."]
func OddCount(lst []string) []string {
    res := make([]string, 0, len(lst))
    for _, arr := range lst {
        n := 0
        for _, d := range arr {
            if (d - '0') % 2 == 1 {
                n++
            }
        }
        res = append(res, fmt.Sprintf("the number of odd elements %dn the str%dng %d of the %dnput.", n,n,n,n))
    }
    return res
}
