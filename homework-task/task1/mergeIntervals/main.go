package main

import "fmt"

func main() {
	//request := [][]int{{1, 3}, {2, 6}, {15, 18}, {8, 10}}
	request := [][]int{{2, 3}, {4, 5}, {6, 7}, {8, 9}, {1, 10}}
	result := merge(request)
	fmt.Println(request)
	fmt.Println(result)
}

func merge(intervals [][]int) [][]int {
	merge := make([][]int, 0, len(intervals))

	for i := 0; i < len(intervals); i++ {
		for j := i + 1; j < len(intervals); j++ {
			if intervals[j][0] < intervals[i][0] {
				intervals[i], intervals[j] = intervals[j], intervals[i]
			}
		}
	}

	merge = append(merge, intervals[0])

	for _, v := range intervals[1:] {
		pre := merge[len(merge)-1]
		if pre[1] >= v[0] {
			if pre[1] < v[1] {
				pre[1] = v[1]
			}
		} else {
			merge = append(merge, v)
		}
	}

	return merge
}
