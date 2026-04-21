package main

// Filter given list of any values only for integers
// >>> FilterIntegers(['a', 3.14, 5])
// [5]
// >>> FilterIntegers([1, 2, 3, 'abc', {}, []])
// [1, 2, 3]
func FilterIntegers(values []interface{}) []int {
    result := make([]int, 0)
    for _, val := range values {
        switch i := val.(type) {
        case int:
            result = append(result, i)
        }
    }
    return result
}
