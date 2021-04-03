package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreetingDeadline(ctx); err != nil {
			fmt.Printf("cannot print greeting: %v\n", err)
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewellDeadline(ctx); err != nil {
			fmt.Printf("cannot print farewell: %v\n", err)
		}
	}()

	wg.Wait()
}

func printGreetingDeadline(ctx context.Context) error {
	greeting, err := genGreetingDeadline(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world\n", greeting)
	return nil
}

func printFarewellDeadline(ctx context.Context) error {
	greeting, err := genFarewellDeadline(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world\n", greeting)
	return nil
}

func genGreetingDeadline(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	switch locale, err := localeDeadline(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func genFarewellDeadline(ctx context.Context) (string, error) {
	switch locale, err := localeDeadline(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func localeDeadline(ctx context.Context) (string, error) {
	if deadline, ok := ctx.Deadline(); ok {
		if deadline.Sub(time.Now().Add(10*time.Second)) <= 0 {
			return "", context.DeadlineExceeded
		}
	}
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(10 * time.Second):
	}
	return "EN/US", nil
}
