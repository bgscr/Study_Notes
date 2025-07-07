package main

import (
	"fmt"
)

func main() {
	intSlice := []int{0, 1, 1, 3, 4, 5, 6, 7, 7, 7, 8, 12, 12, 15, 26, 99}
	removeDuplicates(intSlice)
	fmt.Print(intSlice)
}

func removeDuplicates(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	slow := 1
	for fast := 1; fast < n; fast++ {
		if nums[fast] != nums[fast-1] {
			nums[slow] = nums[fast]
			slow++
		}
	}

	maxInt := int(^uint(0) >> 1)

	for i := range nums[slow:] {
		nums[slow+i] = maxInt
	}
	return slow
}
