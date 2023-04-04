package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

var sum int32

func myFunc(i interface{}) {
	n := i.(int32)
	for j := 0; j < 5; j++ {
		fmt.Println(fmt.Sprintf("%v->\t%v", n, j))
		time.Sleep(time.Second)
	}
}

func main() {
	defer ants.Release()

	runTimes := 10

	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(5, func(i interface{}) {
		myFunc(i)
		wg.Done()
	})
	defer p.Release()
	// Submit tasks one by one.
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = p.Invoke(int32(i))
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", p.Running())
	fmt.Printf("finish all tasks, result is %d\n", sum)
}
