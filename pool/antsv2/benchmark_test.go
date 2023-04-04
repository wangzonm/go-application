package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/panjf2000/ants/v2"
)

const (
	RunTimes           = 1e3
	PoolCap            = 5e2
	BenchParam         = 10
	DefaultExpiredTime = 10 * time.Second
)

func TestDemo(t *testing.T) {
	fmt.Println("this is TestDemo")
}

func BenchmarkPool(b *testing.B) {
	var wg sync.WaitGroup
	p, _ := ants.NewPool(PoolCap, ants.WithExpiryDuration(DefaultExpiredTime))
	defer p.Release()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			_ = p.Submit(func() {
				demoFunc()
				wg.Done()
			})
		}
		wg.Wait()
	}
}

func demoFunc() {
	time.Sleep(time.Duration(BenchParam) * time.Millisecond)
}
