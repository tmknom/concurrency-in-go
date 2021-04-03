package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	fanout()
}

func fanout() {
	repeatFn := func(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
		valueCh := make(chan interface{})
		go func() {
			defer close(valueCh)
			for {
				select {
				case <-done:
					return
				case valueCh <- fn():
				}
			}
		}()
		return valueCh
	}

	take := func(done <-chan interface{}, valueCh <-chan interface{}, num int) <-chan interface{} {
		takeCh := make(chan interface{})
		go func() {
			defer close(takeCh)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeCh <- <-valueCh:
				}
			}
		}()
		return takeCh
	}

	toInt := func(done <-chan interface{}, valueCh <-chan interface{}) <-chan int {
		intCh := make(chan int)
		go func() {
			defer close(intCh)
			for v := range valueCh {
				select {
				case <-done:
					return
				case intCh <- v.(int):
				}
			}
		}()
		return intCh
	}

	primeFinder := func(done <-chan interface{}, intCh <-chan int) <-chan interface{} {
		primeCh := make(chan interface{})
		go func() {
			defer close(primeCh)
			for integer := range intCh {
				integer -= 1
				prime := true
				for divisor := integer - 1; divisor > 1; divisor-- {
					if integer%divisor == 0 {
						prime = false
						break
					}
				}

				if prime {
					select {
					case <-done:
						return
					case primeCh <- integer:
					}
				}
			}
		}()
		return primeCh
	}

	fanIn := func(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
		var wg sync.WaitGroup
		multiplexedCh := make(chan interface{})

		multiplex := func(c <-chan interface{}) {
			defer wg.Done()
			for i := range c {
				select {
				case <-done:
					return
				case multiplexedCh <- i:
				}
			}
		}

		wg.Add(len(channels))
		for _, c := range channels {
			go multiplex(c)
		}

		go func() {
			wg.Wait()
			close(multiplexedCh)
		}()

		return multiplexedCh
	}

	done := make(chan interface{})
	defer close(done)

	start := time.Now()
	rand := func() interface{} { return rand.Intn(50000000) }
	randIntCh := toInt(done, repeatFn(done, rand))

	numFinders := runtime.NumCPU()
	fmt.Printf("Spinning up %d prime finders\n", numFinders)
	finders := make([]<-chan interface{}, numFinders)
	fmt.Println("Primes:")
	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntCh)
	}
	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime)
	}
	fmt.Printf("Search took: %v\n", time.Since(start))
}

func noFanout() {
	repeatFn := func(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
		valueCh := make(chan interface{})
		go func() {
			defer close(valueCh)
			for {
				select {
				case <-done:
					return
				case valueCh <- fn():
				}
			}
		}()
		return valueCh
	}

	take := func(done <-chan interface{}, valueCh <-chan interface{}, num int) <-chan interface{} {
		takeCh := make(chan interface{})
		go func() {
			defer close(takeCh)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeCh <- <-valueCh:
				}
			}
		}()
		return takeCh
	}

	toInt := func(done <-chan interface{}, valueCh <-chan interface{}) <-chan int {
		intCh := make(chan int)
		go func() {
			defer close(intCh)
			for v := range valueCh {
				select {
				case <-done:
					return
				case intCh <- v.(int):
				}
			}
		}()
		return intCh
	}

	primeFinder := func(done <-chan interface{}, intCh <-chan int) <-chan interface{} {
		primeCh := make(chan interface{})
		go func() {
			defer close(primeCh)
			for integer := range intCh {
				integer -= 1
				prime := true
				for divisor := integer - 1; divisor > 1; divisor-- {
					if integer%divisor == 0 {
						prime = false
						break
					}
				}

				if prime {
					select {
					case <-done:
						return
					case primeCh <- integer:
					}
				}
			}
		}()
		return primeCh
	}

	done := make(chan interface{})
	defer close(done)

	rand := func() interface{} { return rand.Intn(50000000) }
	start := time.Now()

	randIntCh := toInt(done, repeatFn(done, rand))
	fmt.Println("Primes:")
	for prime := range take(done, primeFinder(done, randIntCh), 10) {
		fmt.Printf("\t%d\n", prime)
	}
	fmt.Printf("Search took: %v\n", time.Since(start))
}
