package main

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {
	channel()
	closeRead()
	wait()
	buffer()
	owner()
}

func owner() {
	fmt.Println("\nCreate chanOwner")
	chanOwner := func() <-chan int {
		resultCh := make(chan int, 5)
		go func() {
			defer close(resultCh)
			defer fmt.Println("Closed resultCh")
			for i := 0; i < 5; i++ {
				fmt.Printf("Sending %d\n", i)
				resultCh <- i
			}
		}()
		return resultCh
	}

	readCh := chanOwner()
	for integer := range readCh {
		fmt.Printf("Received %v\n", integer)
	}
	fmt.Println("Done receiving")
}

func buffer() {
	var stdoutBuff bytes.Buffer
	defer stdoutBuff.WriteTo(os.Stdout)

	intCh := make(chan int, 4)
	go func() {
		defer close(intCh)
		defer fmt.Fprintln(&stdoutBuff, "Producer Done")
		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdoutBuff, "Sending %d\n", i)
			intCh <- i
		}
	}()

	for integer := range intCh {
		fmt.Fprintf(&stdoutBuff, "Received %v\n", integer)
	}
	fmt.Println("")
}

func wait() {
	begin := make(chan interface{})
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Printf("%v has created\n", i)
			<-begin
			fmt.Printf("%v has begun\n", i)
		}(i)
	}
	time.Sleep(100)
	fmt.Println("Unblocking goroutines...")
	close(begin)
	wg.Wait()
}

func closeRead() {
	intCh := make(chan int)
	go func() {
		defer close(intCh)
		for i := 0; i < 5; i++ {
			intCh <- i
		}
	}()
	for integer := range intCh {
		fmt.Printf("%v ", integer)
	}
	fmt.Println("")
}

func channel() {
	stringCh := make(chan string)
	go func() {
		stringCh <- "Hello channels"
	}()
	result, ok := <-stringCh
	fmt.Printf("(%v) %v\n", ok, result)
}
