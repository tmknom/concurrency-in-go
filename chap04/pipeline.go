package main

import (
	"fmt"
	"math/rand"
)

func main() {
	typeAssertion()
}

func typeAssertion() {
	repeat := func(done <-chan interface{}, values ...interface{}) <-chan interface{} {
		valueCh := make(chan interface{})
		go func() {
			defer close(valueCh)
			for {
				for _, value := range values {
					select {
					case <-done:
						return
					case valueCh <- value:
					}
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

	toString := func(done <-chan interface{}, valueCh <-chan interface{}) <-chan string {
		stringCh := make(chan string)
		go func() {
			defer close(stringCh)
			for v := range valueCh {
				select {
				case <-done:
					return
				case stringCh <- v.(string):
				}
			}
		}()
		return stringCh
	}

	done := make(chan interface{})
	defer close(done)

	var message string
	for token := range toString(done, take(done, repeat(done, "I", "am."), 5)) {
		message += token
	}
	fmt.Printf("message: %s...\n", message)
}

func repeatFunction() {
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

	done := make(chan interface{})
	defer close(done)

	rand := func() interface{} { return rand.Int() }

	for v := range take(done, repeatFn(done, rand), 10) {
		fmt.Println(v)
	}
}

func repeat() {
	repeat := func(done <-chan interface{}, values ...interface{}) <-chan interface{} {
		valueCh := make(chan interface{})
		go func() {
			defer close(valueCh)
			for {
				for _, value := range values {
					select {
					case <-done:
						return
					case valueCh <- value:
					}
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

	done := make(chan interface{})
	defer close(done)

	for v := range take(done, repeat(done, 1), 10) {
		fmt.Println(v)
	}
}

func generator() {
	generator := func(done <-chan interface{}, integers ...int) <-chan int {
		intCh := make(chan int, len(integers))
		go func() {
			defer close(intCh)
			for _, integer := range integers {
				select {
				case <-done:
					return
				case intCh <- integer:
				}
			}
		}()
		return intCh
	}

	multiply := func(done <-chan interface{}, intCh <-chan int, multiplier int) <-chan int {
		multipliedCh := make(chan int)
		go func() {
			defer close(multipliedCh)
			for value := range intCh {
				select {
				case <-done:
					return
				case multipliedCh <- value * multiplier:
				}
			}
		}()
		return multipliedCh
	}

	add := func(done <-chan interface{}, intCh <-chan int, additive int) <-chan int {
		addedCh := make(chan int)
		go func() {
			defer close(addedCh)
			for value := range intCh {
				select {
				case <-done:
					return
				case addedCh <- value + additive:
				}
			}
		}()
		return addedCh
	}

	done := make(chan interface{})
	defer close(done)

	intCh := generator(done, 1, 2, 3, 4)
	pipeline := multiply(done, add(done, multiply(done, intCh, 2), 1), 2)

	for v := range pipeline {
		fmt.Println(v)
	}
}

func pipeline() {
	multiply := func(values []int, multiplier int) []int {
		multipliedValues := make([]int, len(values))
		for i, value := range values {
			multipliedValues[i] = value * multiplier
		}
		return multipliedValues
	}

	add := func(values []int, additive int) []int {
		addedValues := make([]int, len(values))
		for i, value := range values {
			addedValues[i] = value + additive
		}
		return addedValues
	}

	ints := []int{1, 2, 3, 4}
	for _, v := range add(multiply(ints, 2), 1) {
		fmt.Println(v)
	}

	fmt.Println("")

	for _, v := range multiply(add(multiply(ints, 2), 1), 2) {
		fmt.Println(v)
	}
}
