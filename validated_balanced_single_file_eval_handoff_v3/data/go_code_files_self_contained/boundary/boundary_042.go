package main

// Write a function that takes a string and returns true if the string
// length is a prime number or false otherwise
// Examples
// PrimeLength('Hello') == true
// PrimeLength('abcdcba') == true
// PrimeLength('kittens') == true
// PrimeLength('orange') == false
func PrimeLength(s string) bool {
    l := len(s)
    if l == 0 || l == 1 {
        return false
    }
    for i := 2; i < l; i++ {
        if l%i == 0 {
            return false
        }
    }
    return true
}
