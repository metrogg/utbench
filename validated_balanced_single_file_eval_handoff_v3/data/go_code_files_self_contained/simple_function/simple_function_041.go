package main

// Return list of all prefixes from shortest to longest of the input string
// >>> AllPrefixes('abc')
// ['a', 'ab', 'abc']
func AllPrefixes(str string) []string{
    prefixes := make([]string, 0, len(str))
	for i := 0; i < len(str); i++ {
		prefixes = append(prefixes, str[:i+1])
	}
	return prefixes
}
