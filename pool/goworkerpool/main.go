package main

import (
	"fmt"
	"github.com/Joker666/goworkerpool/workerpool"
	"time"
)

func main() {
	var allTask []*workerpool.Task
	for i := 1; i <= 10; i++ {
		ii := i
		task := workerpool.NewTask(func(data interface{}) error {
			//taskID := data.(int)
			for j := 0; j < 5; j++ {
				fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
				time.Sleep(time.Second)
			}
			return nil
		}, i)
		allTask = append(allTask, task)
	}

	pool := workerpool.NewPool(allTask, 5)
	pool.Run()

}
