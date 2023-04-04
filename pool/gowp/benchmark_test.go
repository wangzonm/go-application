package main

import (
	"sync"
	"testing"
	"time"

	"github.com/xxjwxc/gowp/workpool"
)

const (
	RunTimes           = 1e6
	PoolCap            = 5e4
	BenchParam         = 10
	DefaultExpiredTime = 10 * time.Second
)

func BenchmarkPool(b *testing.B) {
	var wg sync.WaitGroup

	wp := workpool.New(PoolCap)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			wp.Do(func() error {
				demoFunc()
				wg.Done()
				return nil
			})
		}
		wg.Wait()
	}
	wp.Wait()
}

func demoFunc() error {
	time.Sleep(time.Duration(BenchParam) * time.Millisecond)
	return nil
}
