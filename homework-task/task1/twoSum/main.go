package main

import "fmt"

func main() {
	fmt.Println(twoSum([]int{1, 2, 3, 4, 5, 6}, 6))

	fmt.Println(twoSum([]int{3, 3}, 6))

	fmt.Println(twoSum([]int{3, 2, 4}, 6))
}

func twoSum(nums []int, target int) []int {
	indexMap := make(map[int]int)
	for i, v := range nums {
		if index, ok := indexMap[target-v]; ok {
			return []int{index, i}
		}
		indexMap[v] = i
	}

	return nil
}
