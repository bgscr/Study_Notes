package main

import "fmt"

func main() {
	fmt.Println(isPalindrome(23222232))
	fmt.Println(isPalindromev2(123321))
}

func isPalindrome(x int) bool {

	revert := 0
	for x > revert {
		revert = revert*10 + x%10
		x /= 10
	}

	return x == revert || x == revert/10
}

func isPalindromev2(x int) bool {
	strarr := []rune{}
	for x > 0 {
		strarr = append(strarr, rune(x%10))
		x /= 10
	}
	for i := 0; i < len(strarr); i++ {
		if strarr[i] != strarr[len(strarr)-1-i] {
			return false
		}
	}
	return true

}
