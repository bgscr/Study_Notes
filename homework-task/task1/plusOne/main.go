package main

import "fmt"

func main() {

	fmt.Println(plusOne([]int{1, 2, 3}))
	fmt.Println(plusOne([]int{1, 2, 9}))
	fmt.Println(plusOne([]int{1, 9, 9, 9, 9}))
	fmt.Println(plusOne([]int{9, 9, 9, 9, 9}))
	fmt.Println(plusOne([]int{9, 9, 9, 9, 8}))

}

func plusOne(digits []int) []int {
	for i := len(digits) - 1; i >= 0; i-- {
		digits[i]++
		if digits[i] < 10 {
			return digits
		}
		digits[i] = 0

	}

	if digits[0] == 0 {
		digits = make([]int, len(digits)+1)
		digits[0] = 1
	}
	return digits
}
