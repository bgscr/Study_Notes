package main

import (
	"fmt"
)

func producer(ch chan<- int) {
	for i := 0; i < 55; i++ {
		ch <- i
	}
	close(ch) // 关闭通道触发消费者退出
}

func consumer(id int, ch <-chan int, done chan<- struct{}) {
	for num := range ch {
		fmt.Printf("消费者%d 接收: %d\n", id, num)
	}
	//零值 struct{}{}，占用内存0
	done <- struct{}{}
}

func main() {
	ch := make(chan int)
	done := make(chan struct{})

	go producer(ch)

	consumerCount := 3
	for i := 0; i < consumerCount; i++ {
		go consumer(i, ch, done)
	}

	for i := 0; i < consumerCount; i++ {
		<-done
	}
	fmt.Println("所有消费者结束")
}
