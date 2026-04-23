package main

import (
    "regexp"
)

// You'll be given a string of words, and your task is to count the number
// of boredoms. A boredom is a sentence that starts with the word "I".
// Sentences are delimited by '.', '?' or '!'.
// 
// For example:
// >>> IsBored("Hello world")
// 0
// >>> IsBored("The sky is blue. The sun is shining. I love this weather")
// 1
func IsBored(S string) int {
    r, _ := regexp.Compile(`[.?!]\s*`)
    sentences := r.Split(S, -1)
    sum := 0
    for _, s := range sentences {
        if len(s) >= 2 && s[:2] == "I " {
            sum++
        }
    }
    return sum
}
