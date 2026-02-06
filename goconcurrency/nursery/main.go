package main

import (
	"fmt"
	"sync"
	"time"
)

// it has to have a something to control the main (parent) flow

type mutexNursery struct {
	mu   sync.Mutex
	body func(args ...interface{})
}

func NewNursery(f func(args ...interface{})) *mutexNursery {
	return &mutexNursery{
		mu:   sync.Mutex{},
		body: f,
	}
}

func (n *mutexNursery) Run() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.body()
}

/* the problem with the abpve approach is that nested irresponsible (no cancellation) goroutines can still be found.
we want to avoid this.
we want to wrap the standard goroutine into a "responsible" goroutine.
we need a cancellation signal to be set such that the caller (parent) will wait until the child is finished.
*/
///////////////////////////////////////////////////////////////

//channel-controlled-nursery

// why channels? because they control teh flow and the order.
// "the nursery block doesn't exit until all the tasks inside it have exited –
// if the parent task reaches the end of the block before all the children are finished,
// then it pauses there and waits for them. The nursery automatically expands to hold the children."
//
// at the same time of implementing this, you want to allow concurrent execution of the tasks. tasks here would be goroutines.
// it is up to the user to define the concurrency level. so this would be getting the "body" of the nursery from the user.
// this works under the assumption that the tasks functions do not use nested goroutines.
type nurseryBundle struct {
	tasks []func(args ...interface{})
	// will have a limit that is equal to the number of tasks, when that is full the nursery will be done and thus exit.
	// the signals should come from the task exec it self and since we want that to  be concurrent we'll call the tasks
	// through go routines
	done chan int
}

func NewNurseryBundle(tasks []func(args ...interface{})) *nurseryBundle {
	return &nurseryBundle{
		tasks: tasks,
		done:  make(chan int, len(tasks)),
	}
}

func (b *nurseryBundle) Run() {
	for _, task := range b.tasks {
		go func() {
			task()
			b.done <- 1
		}()
	}

	//while the channel is not full, keep waiting
	for {
		if len(b.done) == len(b.tasks) {
			close(b.done)
			return
		}
	}
}

// now we use it
func main() {
	fmt.Println("Starting nursery")
	nursery := NewNurseryBundle([]func(args ...interface{}){
		func(args ...interface{}) {
			time.Sleep(1 * time.Second)
			fmt.Println("task 1")
		},
		func(args ...interface{}) {
			time.Sleep(5 * time.Second)
			fmt.Println("task 2")
		},
		func(args ...interface{}) {
			fmt.Println("task 3")
		},
	})
	nursery.Run() // expected 3 1 2
	fmt.Println("Nursery finished")

	// while with bare go routines the main will exit before the tasks are finished
	go func() {
		fmt.Println("task 1")
	}()
	go func() {
		time.Sleep(10 * time.Second)
		fmt.Println("task 2")
	}()
}
