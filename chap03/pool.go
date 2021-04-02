package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	pool()
	calc()
}

func pool() {
	myPool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating new instance")
			return struct{}{}
		},
	}

	myPool.Get()
	instance := myPool.Get()
	myPool.Put(instance)
	myPool.Get()
}

func calc() {
	var numCalcsCreated int

	calcPool := &sync.Pool{
		New: func() interface{} {
			numCalcsCreated += 1
			mem := make([]byte, 1024)
			return &mem
		},
	}
	calcPool.Put(calcPool.Get())
	calcPool.Put(calcPool.Get())
	calcPool.Put(calcPool.Get())
	calcPool.Put(calcPool.Get())

	const numWorkers = 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			mem := calcPool.Get().(*[]byte)
			defer calcPool.Put(mem)
			time.Sleep(1)
		}()
	}

	wg.Wait()
	fmt.Printf("created: %d, run: %d\n", numCalcsCreated, numWorkers)
}
