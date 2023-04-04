package main

import (
	"sync"
	"testing"
	"time"

	"github.com/Joker666/goworkerpool/workerpool"
)

const (
	RunTimes           = 1e3
	PoolCap            = 5e2
	BenchParam         = 10
	DefaultExpiredTime = 10 * time.Second
)

func BenchmarkPool(b *testing.B) {
	var wg sync.WaitGroup

	b.ResetTimer()
	var allTask []*workerpool.Task
	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for i := 0; i <= RunTimes; i++ {
			task := workerpool.NewTask(func(i interface{}) error {
				demoFunc()
				wg.Done()
				return nil
			}, i)
			allTask = append(allTask, task)
		}
		wg.Wait()
	}
	pool := workerpool.NewPool(allTask, 5)
	pool.Run()
}

func demoFunc() error {
	time.Sleep(time.Duration(BenchParam) * time.Millisecond)
	return nil
}
