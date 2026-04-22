package main

import (
    "strings"
)

// Given a string s, count the number of uppercase vowels in even indices.
// 
// For example:
// CountUpper('aBCdEf') returns 1
// CountUpper('abcdefg') returns 0
// CountUpper('dBBE') returns 0
func CountUpper(s string) int {
    count := 0
    runes := []rune(s)
    for i := 0; i < len(runes); i += 2 {
        if strings.ContainsRune("AEIOU", runes[i]) {
            count += 1
        }
    }
    return count
}
