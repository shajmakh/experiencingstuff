package main

import "time"

func main() {
	// run FOREVER infinite goroutines with cancellation
	for i := 0; i < 10; i++ {
		go func() {
			for {
				time.Sleep(2 * time.Minute)
			}
		}()
	}
	time.Sleep(20 * time.Minute)
}
