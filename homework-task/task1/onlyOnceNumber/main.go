package main

import "fmt"

func main() {
	nums := []int{2, 2, 4, 12, 4, 7, 7, 9, 6, 9, 6}
	fmt.Println(singleNumber(nums))

}

func singleNumber(nums []int) int {

	// res := 0
	// for i := 0; i < len(nums); i++ {
	// 	res ^= nums[i]
	// }
	// return res

	numMap := make(map[int]int)
	for i := 0; i < len(nums); i++ {
		numMap[nums[i]]++
	}

	for key, value := range numMap {
		if value == 1 {
			return key
		}
	}

	return -1

}
