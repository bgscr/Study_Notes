package main

import (
	"sync"
	"sync/atomic"
)

func main() {
	//mutexTest()
	atomicTest()
}

func mutexTest() {

	m := sync.Mutex{}
	count := 0
	var w sync.WaitGroup
	w.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer w.Done()
			for j := 0; j < 1000; j++ {
				m.Lock()
				count++
				m.Unlock()

			}
		}()
	}

	w.Wait()
	println(count)
}

func atomicTest() {
	var count uint32 = 0
	var countPoint *uint32 = &count
	var wg sync.WaitGroup

	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddUint32(countPoint, 1)
			}
		}()
	}

	wg.Wait()
	println(atomic.LoadUint32(countPoint))

}
