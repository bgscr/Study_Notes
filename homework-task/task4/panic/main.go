package main

import (
	"fmt"
	"strconv"
)

func main() {
	for i := 0; i < 10; i++ {
		go SafeCall(i)
	}
}

func SafeCall(i int) {
	defer func() { // 每个 SafeCall 有自己的 recover
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()
	A(i) // 执行可能 panic 的操作
}

func A(i int) {
	panic("custom error" + strconv.Itoa(i))
}
