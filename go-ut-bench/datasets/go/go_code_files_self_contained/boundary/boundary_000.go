package main

import (
    "strings"
)

// Create a function which takes a string representing a file's name, and returns
// 'Yes' if the the file's name is valid, and returns 'No' otherwise.
// A file's name is considered to be valid if and only if all the following conditions
// are met:
// - There should not be more than three digits ('0'-'9') in the file's name.
// - The file's name contains exactly one dot '.'
// - The substring before the dot should not be empty, and it starts with a letter from
// the latin alphapet ('a'-'z' and 'A'-'Z').
// - The substring after the dot should be one of these: ['txt', 'exe', 'dll']
// Examples:
// FileNameCheck("example.txt") # => 'Yes'
// FileNameCheck("1example.dll") # => 'No' (the name should start with a latin alphapet letter)
func FileNameCheck(file_name string) string {
    suf := []string{"txt", "exe", "dll"}
    lst := strings.Split(file_name, ".")
    isInArray := func (arr []string, x string) bool {
        for _, y := range arr {
            if x == y {
                return true
            }
        }
        return false
    }
    switch {
    case len(lst) != 2:
        return "No"
    case !isInArray(suf, lst[1]):
        return "No"
    case len(lst[0]) == 0:
        return "No"
    case 'a' > strings.ToLower(lst[0])[0] || strings.ToLower(lst[0])[0] > 'z':
        return "No"
    }
    t := 0
    for _, c := range lst[0] {
        if '0' <= c && c <= '9' {
            t++
        }
    }
    if t > 3 {
        return "No"
    }
    return "Yes"
}
