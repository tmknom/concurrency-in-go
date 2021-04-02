package main

import (
	"fmt"
	"sync"
)

func main() {
	foreach2()
}

func foreach2() {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "world", "bye"} {
		wg.Add(1)
		go func(salutation string) {
			defer wg.Done()
			fmt.Println(salutation)
		}(salutation)
	}
	wg.Wait()
}

func foreach1() {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "world", "bye"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(salutation)
		}()
	}
	wg.Wait()
}

func clojure() {
	var wg sync.WaitGroup
	salutation := "hello"

	wg.Add(1)
	go func() {
		defer wg.Done()
		salutation = "welcome"
	}()
	wg.Wait()
	fmt.Println(salutation)
}
