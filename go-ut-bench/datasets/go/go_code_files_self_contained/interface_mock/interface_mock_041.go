package main

import (
    "strings"
)

// Given a string representing a space separated lowercase letters, return a dictionary
// of the letter with the most repetition and containing the corresponding count.
// If several letters have the same occurrence, return all of them.
// 
// Example:
// Histogram('a b c') == {'a': 1, 'b': 1, 'c': 1}
// Histogram('a b b a') == {'a': 2, 'b': 2}
// Histogram('a b c a b') == {'a': 2, 'b': 2}
// Histogram('b b b b a') == {'b': 4}
// Histogram('') == {}
func Histogram(test string) map[rune]int {
    dict1 := make(map[rune]int)
    list1 := strings.Fields(test)
    t := 0
    count := func(lst []string, v string) int {
        cnt := 0
        for _, i := range lst {
            if i == v {
                cnt++
            }
        }
        return cnt
    }
    for _, i := range list1 {
        if c := count(list1, i); c>t && i!="" {
            t=c
        }
    }
    if t>0 {
        for _, i := range list1 {
            if count(list1, i)==t {
                dict1[[]rune(i)[0]]=t
            }
        }
    }
    return dict1
}
