package main

import (
	"fmt"
	"time"

	"github.com/alitto/pond"
)

func main() {
	// Create a buffered (non-blocking) pool that can scale up to 100 workers
	// and has a buffer capacity of 1000 tasks
	pool := pond.New(5, 5)

	// Submit 1000 tasks
	for i := 0; i < 10; i++ {
		n := i
		pool.Submit(func() {
			for j := 0; j < 5; j++ {
				fmt.Println(fmt.Sprintf("%v->\t%v", n, j))
				time.Sleep(time.Second)
			}
		})
	}

	// Stop the pool and wait for all submitted tasks to complete
	pool.StopAndWait()
}
