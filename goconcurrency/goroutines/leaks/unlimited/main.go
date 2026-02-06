package main

import "time"

func main() {
	// run FOREVER an infinite goroutines without cancellation
	for {
		go func() {
			for {
				time.Sleep(time.Second)
			}
		}()
	}
}
