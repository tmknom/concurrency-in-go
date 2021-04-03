package main

import "fmt"

func main() {
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

	tee := func(done <-chan interface{}, in <-chan interface{}) (_, _ <-chan interface{}) {
		out1Ch := make(chan interface{})
		out2Ch := make(chan interface{})
		go func() {
			defer close(out1Ch)
			defer close(out2Ch)
			for val := range orDone(done, in) {
				var out1Ch, out2Ch = out1Ch, out2Ch
				for i := 0; i < 2; i++ {
					select {
					case out1Ch <- val:
						out1Ch = nil
					case out2Ch <- val:
						out2Ch = nil
					}
				}
			}
		}()
		return out1Ch, out2Ch
	}

	done := make(chan interface{})
	defer close(done)

	out1, out2 := tee(done, take(done, repeat(done, 1, 2, 3), 5))
	for val := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val, 2*(<-out2).(int))
	}
}
