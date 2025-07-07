package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan int, 3)

	var wg sync.WaitGroup

	wg.Add(2)

	go func(ch chan<- int) {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			ch <- i
		}
		close(ch)
	}(ch)

	go func(ch <-chan int) {
		defer wg.Done()
		for v := range ch {
			fmt.Println(v)
		}
	}(ch)

	wg.Wait()
}
