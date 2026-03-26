package main

import (
	"fmt"
	"sync"
)

func runSyncMap() {
	var safeMap sync.Map
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()

			safeMap.Store("key", val)
		}(i)
	}
	wg.Wait()

	if val, ok := safeMap.Load("key"); ok {
		fmt.Printf("sync.Map Result: %d\n", val)
	}
}

type SafeMutexMap struct {
	mu sync.RWMutex
	m  map[string]int
}

func runMutexMap() {
	sm := SafeMutexMap{
		m: make(map[string]int),
	}
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			sm.mu.Lock()
			sm.m["key"] = val
			sm.mu.Unlock()
		}(i)
	}
	wg.Wait()

	sm.mu.RLock()
	val := sm.m["key"]
	sm.mu.RUnlock()

	fmt.Printf("sync.RWMutex Result: %d\n", val)
}

func main() {
	fmt.Println("Способ 1: sync.Map")
	runSyncMap()

	fmt.Println("\nСпособ 2: sync.RWMutex")
	runMutexMap()
}
