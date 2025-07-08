package main

import (
	"fmt"
	"math/rand"
)

func producer(id int, dataCh chan<- int, stopCh <-chan struct{}) {
	for {
		select {
		case <-stopCh:
			fmt.Printf("生产者%d 退出\n", id)
			return
		default:
			data := rand.Intn(10)
			dataCh <- data
			fmt.Printf("生产者%d 发送: %d\n", id, data)
		}
	}
}

func consumer(id int, dataCh <-chan int, sumCh chan<- int) {
	for num := range dataCh {
		sumCh <- num // 将数据发送给sum协程
		fmt.Printf("消费者%d 接收: %d\n", id, num)
	}
}

func sumCounter(sumCh <-chan int, triggerCloseCh chan<- struct{}) {
	var globalSum int // 全局累加器
	for num := range sumCh {
		globalSum += num
		fmt.Printf("当前全局总和: %d\n", globalSum)

		if globalSum >= 300 {
			close(triggerCloseCh) // 直接关闭触发通道
			return
		}
	}
}

func coordinator(stopCh chan<- struct{}, triggerCloseCh <-chan struct{}) {
	<-triggerCloseCh // 等待触发关闭信号
	close(stopCh)
	fmt.Println("信号通道已关闭")
}

func main() {
	dataCh := make(chan int, 10)
	stopCh := make(chan struct{})
	triggerCloseCh := make(chan struct{}) // 改为无缓冲通道
	sumCh := make(chan int, 100)          // 新增sum通道

	go coordinator(stopCh, triggerCloseCh)

	// 启动sum协程
	go sumCounter(sumCh, triggerCloseCh)

	// 启动3个生产者
	for i := 0; i < 20; i++ {
		go producer(i, dataCh, stopCh)
	}

	// 启动3个消费者
	for i := 0; i < 20; i++ {
		go consumer(i, dataCh, sumCh)
	}

	select {} // 阻塞主协程
}
