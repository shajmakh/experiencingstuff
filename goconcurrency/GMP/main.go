package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
	"time"
)

func main() {
	// 1. Setup the Trace
	// We create a file to store the trace data
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Start tracing to that file
	if err := trace.Start(f); err != nil {
		panic(err)
	}
	defer trace.Stop()

	// 2. Control Parallelism
	// We limit Go to use 2 logical Processors (Ps) to make the visualization easier to read.
	runtime.GOMAXPROCS(2)

	var wg sync.WaitGroup

	// 3. Simulate "Worker" Goroutines
	// We launch 5 goroutines. Since we only have 2 Ps, they will have to compete.
	fmt.Println("Starting 5 workers...")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(id)
		}(i)
	}

	wg.Wait()
	fmt.Println("All done.")
}

func worker(id int) {
	fmt.Printf("Worker %d starting\n", id)

	// Phase A: Blocking Work (Simulating I/O)
	// This takes the G off the P, allowing other Gs to run.
	time.Sleep(100 * time.Millisecond)

	// Phase B: CPU Work
	// This hogs the P.
	count := 0
	for i := 0; i < 100000000; i++ {
		count += i
	}
	
	fmt.Printf("Worker %d finished\n", id)
}
