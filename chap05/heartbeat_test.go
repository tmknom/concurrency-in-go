package main

import (
	"testing"
	"time"
)

func TestDoWork(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	heartbeat, results := DoWork(done, intSlice...)

	<-heartbeat

	i := 0
	for result := range results {
		if expected := intSlice[i]; result != expected {
			t.Errorf("index %v: expected %v, but received %v,", i, expected, result)
		}
		i++
	}
}

func TestDoWork2(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	const timeout = 2 * time.Second
	heartbeat, results := DoWork2(done, timeout/2, intSlice...)

	<-heartbeat

	i := 0
	for {
		select {
		case result, ok := <-results:
			if !ok {
				return
			} else if expected := intSlice[i]; result != expected {
				t.Errorf("index %v: expected %v, but received %v,", i, expected, result)
			}
			i++
		case <-heartbeat:
		case <-time.After(timeout):
			t.Fatal("test time out")
		}
	}
}

func DoWork(done <-chan interface{}, nums ...int) (<-chan interface{}, <-chan int) {
	heartbeatCh := make(chan interface{}, 1)
	intCh := make(chan int)

	go func() {
		defer close(heartbeatCh)
		defer close(intCh)

		time.Sleep(2 * time.Second)

		for _, num := range nums {
			select {
			case heartbeatCh <- struct{}{}:
			default:
			}

			select {
			case <-done:
				return
			case intCh <- num:
			}
		}
	}()
	return heartbeatCh, intCh
}

func DoWork2(done <-chan interface{}, pulseInterval time.Duration, nums ...int) (<-chan interface{}, <-chan int) {
	heartbeatCh := make(chan interface{}, 1)
	intCh := make(chan int)

	go func() {
		defer close(heartbeatCh)
		defer close(intCh)

		time.Sleep(2 * time.Second)

		pulse := time.Tick(pulseInterval)

	numLoop:
		for _, num := range nums {
			for {
				select {
				case <-done:
					return
				case <-pulse:
					select {
					case heartbeatCh <- struct{}{}:
					default:
					}
				case intCh <- num:
					continue numLoop
				}
			}
		}
	}()
	return heartbeatCh, intCh
}
