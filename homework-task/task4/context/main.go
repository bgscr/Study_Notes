package main

import (
	"context"
	"fmt"
	"myContext/myCtx"
	"time"
)

func worker(ctx context.Context, id int) {
	fmt.Printf("worker %d start\n", id)
	defer fmt.Printf("worker %d exit\n", id)

	select {
	case <-ctx.Done():
		fmt.Printf("worker %d canceled, err=%v,value=%v\n", id, ctx.Err(), ctx.Value(&trace_id))
	case <-time.After(time.Duration(id) * time.Second):
		fmt.Printf("worker %d done normally\n", id)
	}
}

var trace_id string = "trace_id"
var trace_id2 string = "trace_id2"

func main() {
	// 1. 创建根 Context
	root := context.Background()

	// 2. 派生一个可取消的子 Context
	cancelCtx, cancel := myCtx.WithMyCancel(root)

	// 3. 再派生带值的子 Context
	valueCtx1 := myCtx.WithMyValue(cancelCtx, &trace_id, "parent value")

	// 4. 启动 多个 goroutine 监听同一个 Context
	i := 0
	for ; i < 4; i++ {

		go func(j int) {
			worker(valueCtx1, j)
		}(i)
	}

	valueCtx2 := myCtx.WithMyValue(valueCtx1, &trace_id2, "child value")
	for ; i < 6; i++ {

		go func(j int) {
			worker(valueCtx2, j)
		}(i)
	}

	// 5. 2 秒后统一取消
	time.Sleep(2 * time.Second)
	fmt.Println("main cancel")
	cancel()

	// 6. 等 goroutine 打印
	time.Sleep(10 * time.Second)
}
