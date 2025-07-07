package main

import "fmt"

func main() {

	fmt.Println(isValid(""))
}

func isValid(s string) bool {

	mapParenthese := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}

	stack := []byte{}

	for i := 0; i < len(s); i++ {
		if mapParenthese[s[i]] == 0 {
			stack = append(stack, s[i])
			continue
		}
		if len(stack) == 0 || stack[len(stack)-1] != mapParenthese[s[i]] {
			return false
		}
		stack = stack[:len(stack)-1]
	}

	return len(stack) == 0
}
