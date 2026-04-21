package main

import (
    "strings"
)

// You will be given a string of words separated by commas or spaces. Your task is
// to split the string into words and return an array of the words.
// 
// For example:
// WordsString("Hi, my name is John") == ["Hi", "my", "name", "is", "John"]
// WordsString("One, two, three, four, five, six") == ["One", "two", "three", "four", "five", "six"]
func WordsString(s string) []string {
    s_list := make([]rune, 0)

    for _, c := range s {
        if c == ',' {
            s_list = append(s_list, ' ')
        } else {
            s_list = append(s_list, c)
        }
    }
    return strings.Fields(string(s_list))
}
