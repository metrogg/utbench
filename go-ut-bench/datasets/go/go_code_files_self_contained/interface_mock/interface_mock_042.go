package main

import (
    "math"
)

// Given two lists operator, and operand. The first list has basic algebra operations, and
// the second list is a list of integers. Use the two given lists to build the algebric
// expression and return the evaluation of this expression.
// 
// The basic algebra operations:
// Addition ( + )
// Subtraction ( - )
// Multiplication ( * )
// Floor division ( // )
// Exponentiation ( ** )
// 
// Example:
// operator['+', '*', '-']
// array = [2, 3, 4, 5]
// result = 2 + 3 * 4 - 5
// => result = 9
// 
// Note:
// The length of operator list is equal to the length of operand list minus one.
// Operand is a list of of non-negative integers.
// Operator list has at least one operator, and operand list has at least two operands.
func DoAlgebra(operator []string, operand []int) int {
    higher := func(a, b string) bool {
        if b == "*" || b == "//" || b == "**" {
            return false
        }
        if a == "*" || a == "//" || a == "**" {
            return true
        }
        return false
    }
    for len(operand) > 1 {
        pos := 0
        sign := operator[0]
        for i, str := range operator {
            if higher(str, sign) {
                sign = str
                pos = i
            }
        }
        switch sign {
        case "+":
            operand[pos] += operand[pos+1]
        case "-":
            operand[pos] -= operand[pos+1]
        case "*":
            operand[pos] *= operand[pos+1]
        case "//":
            operand[pos] /= operand[pos+1]
        case "**":
            operand[pos] = int(math.Pow(float64(operand[pos]), float64(operand[pos+1])))
        }
        operator = append(operator[:pos], operator[pos+1:]...)
        operand = append(operand[:pos+1], operand[pos+2:]...)
    }
    return operand [0]
}
