package main

type Stack struct {
	items []int
}

func NewStack() *Stack {
	return &Stack{items: []int{}}
}

func (s *Stack) Push(item int) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() int {
	if s.IsEmpty() {
		panic("stack is empty")
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

func (s *Stack) Peek() int {
	if s.IsEmpty() {
		panic("stack is empty")
	}
	return s.items[len(s.items)-1]
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack) Size() int {
	return len(s.items)
}

func EvaluatePostfix(expression string) int {
	stack := NewStack()
	tokens := splitTokens(expression)
	
	for _, token := range tokens {
		if isNumeric(token) {
			stack.Push(parseInt(token))
		} else {
			b := stack.Pop()
			a := stack.Pop()
			result := applyOperator(token, a, b)
			stack.Push(result)
		}
	}
	
	return stack.Pop()
}

func splitTokens(s string) []string {
	result := []string{}
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ' ' {
			if start < i {
				result = append(result, s[start:i])
			}
			start = i + 1
		}
	}
	return result
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

func parseInt(s string) int {
	result := 0
	for _, c := range s {
		result = result * 10 + int(c - '0')
	}
	return result
}

func applyOperator(op string, a, b int) int {
	switch op {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		return a / b
	default:
		panic("unknown operator: " + op)
	}
}