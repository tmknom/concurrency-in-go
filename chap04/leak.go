package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	terminate()
	random()
}

func random() {
	newRandCh := func(done <-chan interface{}) <-chan int {
		randCh := make(chan int)
		go func() {
			defer fmt.Println("newRandCh exited")
			defer close(randCh)
			fmt.Println("\nnewRandCh started")
			for {
				select {
				case randCh <- rand.Int():
				case <-done:
					return
				}
			}
		}()
		return randCh
	}

	done := make(chan interface{})
	randCh := newRandCh(done)
	for i := 0; i < 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randCh)
	}

	close(done)
	time.Sleep(1 * time.Second)
}

func terminate() {
	doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited")
			defer close(terminated)
			fmt.Println("doWork started")
			for {
				select {
				case s := <-strings:
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan interface{})
	terminated := doWork(done, nil)

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()

	<-terminated
	fmt.Println("Done")
}

func leak() {
	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited")
			defer close(completed)
			fmt.Println("doWork started")
			for s := range strings {
				fmt.Println(s)
			}
		}()
		return completed
	}

	doWork(nil)
	time.Sleep(1 * time.Second)
	fmt.Println("Done")
}
