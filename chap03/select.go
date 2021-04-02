package main

import (
	"fmt"
	"time"
)

func main() {
	single()
	multi()
	timeout()
	loop()
}

func loop() {
	done := make(chan interface{})
	go func() {
		time.Sleep(500 * time.Millisecond)
		close(done)
	}()

	workCounter := 0

forloop:
	for {
		select {
		case <-done:
			break forloop
		default:
			fmt.Printf("work: %v\n", workCounter)
			workCounter++
			time.Sleep(100 * time.Millisecond)
		}
	}
	fmt.Printf("Achieved %v cycles of work before signalled to stop\n", workCounter)
}

func timeout() {
	var c <-chan int
	select {
	case <-c:
		fmt.Println(" Not executed")
	case <-time.After(1 * time.Second):
		fmt.Println("Timed out")
	}
}

func multi() {
	c1 := make(chan interface{})
	close(c1)
	c2 := make(chan interface{})
	close(c2)

	var c1Count, c2Count int
	for i := 0; i < 1000; i++ {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}
	fmt.Printf("c1Count: %d, c2Count: %d\n", c1Count, c2Count)

}

func single() {
	start := time.Now()
	c := make(chan interface{})
	go func() {
		time.Sleep(50 * time.Millisecond)
		close(c)
	}()

	fmt.Println("Blocking on read...")
	select {
	case <-c:
		fmt.Printf("Unblocked %v later\n", time.Since(start))
	}
}
