package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	//printNumber()

	tasks := []func(){
		func() {
			time.Sleep(1 * time.Second)
			fmt.Println("sleep 1 second")
		},
		func() {
			time.Sleep(2 * time.Second)
			fmt.Println("sleep 2 second")
		},
		func() {
			time.Sleep(500 * time.Millisecond)
			fmt.Println("sleep 5 millisecond")
		},
	}

	scheduler := &TaskScheduler{
		results: make(map[int]time.Duration),
	}

	for _, task := range tasks {
		scheduler.AddTask(task)
	}

	result := scheduler.Wait()

	for id, duration := range result {
		fmt.Printf("Task %d: %v\n", id, duration)
	}
}

func printNumber() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 2; i < 11; i += 2 {
			fmt.Println("even number:", i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 1; i < 11; i += 2 {
			fmt.Println("odd number:", i)
		}
	}()

	wg.Wait()
}

// 任务调度器
type TaskScheduler struct {
	wg      sync.WaitGroup        // 协程同步
	mu      sync.Mutex            // 并发安全锁
	results map[int]time.Duration // 任务耗时统计
	taskID  int                   // 任务标识
}

// 添加任务
func (ts *TaskScheduler) AddTask(task func()) {
	ts.wg.Add(1)
	id := ts.taskID
	ts.taskID++

	go func() {
		defer ts.wg.Done()
		start := time.Now()

		task() // 执行任务

		elapsed := time.Since(start)

		ts.mu.Lock() // 写结果时加锁
		defer ts.mu.Unlock()
		ts.results[id] = elapsed
	}()
}

// 等待所有任务完成并返回结果
func (ts *TaskScheduler) Wait() map[int]time.Duration {
	ts.wg.Wait()
	return ts.results
}
