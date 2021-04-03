package main

import "fmt"

func main() {
	orDone := func(done, c <-chan interface{}) <-chan interface{} {
		valueCh := make(chan interface{})
		go func() {
			defer close(valueCh)
			for {
				select {
				case <-done:
					return
				case v, ok := <-c:
					if !ok {
						return
					}
					select {
					case valueCh <- v:
					case <-done:
					}
				}
			}
		}()
		return valueCh
	}

	bridge := func(done <-chan interface{}, chanCh <-chan <-chan interface{}) <-chan interface{} {
		valueCh := make(chan interface{})
		go func() {
			defer close(valueCh)
			for {
				var ch <-chan interface{}
				select {
				case <-done:
					return
				case maybeCh, ok := <-chanCh:
					if !ok {
						return
					}
					ch = maybeCh
				}

				for val := range orDone(done, ch) {
					select {
					case <-done:
					case valueCh <- val:
					}
				}
			}
		}()
		return valueCh
	}

	genValues := func() <-chan <-chan interface{} {
		chanCh := make(chan (<-chan interface{}))
		go func() {
			defer close(chanCh)
			for i := 0; i < 10; i++ {
				ch := make(chan interface{}, 1)
				ch <- i
				close(ch)
				chanCh <- ch
			}
		}()
		return chanCh
	}

	done := make(chan interface{})
	defer close(done)

	for v := range bridge(done, genValues()) {
		fmt.Printf("%v ", v)
	}
	fmt.Println("")
}
