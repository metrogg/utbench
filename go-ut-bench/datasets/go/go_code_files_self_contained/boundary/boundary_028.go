package main

// You are given a string s.
// Your task is to check if the string is happy or not.
// A string is happy if its length is at least 3 and every 3 consecutive letters are distinct
// For example:
// IsHappy(a) => false
// IsHappy(aa) => false
// IsHappy(abcd) => true
// IsHappy(aabb) => false
// IsHappy(adb) => true
// IsHappy(xyy) => false
func IsHappy(s string) bool {
    if len(s) < 3 {
        return false
    }
    for i := 0; i < len(s)-2; i++ {
        if s[i] == s[i+1] || s[i+1] == s[i+2] || s[i] == s[i+2] {
            return false
        }
    }
    return true
}
