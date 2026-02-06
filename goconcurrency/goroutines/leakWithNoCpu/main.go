package main

import "time"

func main() {
	c := make(chan int)
	for i := 0; i < 100000; i++ {
		go func() {
			ints := make([]int, 1000000)
			for i := range ints {
				ints[i] = i
			}
			<-c
		}()
	}
	time.Sleep(10 * time.Minute)
}
