package main

func ReverseString(s string) string {
	if s == "" {
		return ""
	}
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		result[i] = s[len(s)-1-i]
	}
	return string(result)
}

func IsPalindrome(s string) bool {
	if s == "" {
		return false
	}
	reversed := ReverseString(s)
	return s == reversed
}