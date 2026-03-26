package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func fixWithMutex() {
	var counter int
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Printf("Result using sync.Mutex: %d\n", counter)
}

func fixWithAtomic() {
	var counter int64
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			atomic.AddInt64(&counter, 1)
		}()
	}
	wg.Wait()
	fmt.Printf("Result using sync/atomic: %d\n", counter)
}

func main() {
	fixWithMutex()
	fixWithAtomic()
}
