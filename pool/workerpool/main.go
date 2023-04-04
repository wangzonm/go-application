package main

import (
	"fmt"
	"time"

	"github.com/gammazero/workerpool"
)

func main() {
	wp := workerpool.New(5)

	for i := 0; i < 10; i++ {
		ii := i
		wp.Submit(func() {
			for j := 0; j < 5; j++ {
				fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
				time.Sleep(time.Second)
			}
		})
	}

	wp.StopWait()
}
