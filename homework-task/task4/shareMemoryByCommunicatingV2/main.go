package main

import (
	"fmt"
)

func producer(id int, dataCh chan<- int, stopCh <-chan struct{}) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获异常生产者%d 发送失败:\n", id)
		}
	}()
	for {
		select {
		case <-stopCh:
			fmt.Printf("生产者%d 退出\n", id)
			return
		default:
			//data := rand.Intn(10)
			data := 1
			dataCh <- data
			fmt.Printf("生产者%d 发送: %d\n", id, data)
		}
	}
}

func consumer(id int, dataCh <-chan int, sumCh chan<- int, doneCh chan<- struct{}) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获异常：消费者%d 无法发送数据到sumCh\n", id)
		}
		fmt.Printf("消费者%d退出\n", id)
		doneCh <- struct{}{}
	}()

	for num := range dataCh {
		sumCh <- num
		fmt.Printf("消费者%d 接收: %d\n", id, num)
	}

}

func sumCounter(sumCh <-chan int, triggerCloseCh chan<- struct{}) {
	var globalSum int
	for num := range sumCh {
		globalSum += num
		fmt.Printf("当前全局总和: %d\n", globalSum)

		if globalSum >= 300 {
			close(triggerCloseCh)
			return
		}
	}
}

func coordinator(stopCh chan<- struct{}, dataCh chan<- int, triggerCloseCh <-chan struct{}, sumCh chan<- int) {
	<-triggerCloseCh // 等待触发关闭信号
	//time.Sleep(500 * time.Millisecond)
	close(stopCh)
	close(dataCh)
	close(sumCh)
	fmt.Println("信号通道已关闭")
}

func main() {
	dataCh := make(chan int, 10)
	stopCh := make(chan struct{})
	triggerCloseCh := make(chan struct{}) // 改为无缓冲通道
	sumCh := make(chan int, 3)

	go coordinator(stopCh, dataCh, triggerCloseCh, sumCh)

	// 启动sum协程
	go sumCounter(sumCh, triggerCloseCh)

	// 启动个生产者
	producerCount := 10
	for i := 0; i < producerCount; i++ {
		go producer(i, dataCh, stopCh)
	}

	// 启动consumerCount个消费者
	consumerCount := 20
	doneCh := make(chan struct{})
	for i := 0; i < consumerCount; i++ {
		go consumer(i, dataCh, sumCh, doneCh)
	}

	for i := 0; i < consumerCount; i++ {
		<-doneCh
	}
	close(doneCh)
}
