package main

import "fmt"

func main() {
	var num int = 10
	var numPoint *int = &num
	receivePoint(numPoint)
	fmt.Println(numPoint)
	fmt.Println(*numPoint)

	arraysPoint := &[]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	receivePointV2(arraysPoint)
	fmt.Println(*arraysPoint)
}

func receivePoint(num *int) {
	*num += 10

}

func receivePointV2(arrarys *[]int) {
	for i := range *arrarys {
		(*arrarys)[i] *= 2
	}
}
