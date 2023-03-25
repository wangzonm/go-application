package concurrentmap

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

var (
	num = 100000
	bN  = 100
)

/*
goos: darwin
goarch: amd64
pkg: map/concurrentmap
cpu: Intel(R) Core(TM) i5-5287U CPU @ 2.90GHz
BenchmarkConcurrentMap/0-4 	     100	  67790847 ns/op	 3999277 B/op	  499276 allocs/op
BenchmarkConcurrentMap/1-4 	     100	  51639550 ns/op	 3998421 B/op	  499269 allocs/op
BenchmarkConcurrentMap/2-4 	     100	  58651565 ns/op	 3998027 B/op	  499273 allocs/op
BenchmarkConcurrentMap/3-4 	     100	  52596867 ns/op	 3997472 B/op	  499259 allocs/op
BenchmarkConcurrentMap/4-4 	     100	  87013303 ns/op	 3998057 B/op	  499270 allocs/op
BenchmarkSyncMap/0-4       	     100	  69154131 ns/op	 7194667 B/op	  699240 allocs/op
BenchmarkSyncMap/1-4       	     100	  84881340 ns/op	 7195059 B/op	  699240 allocs/op
BenchmarkSyncMap/2-4       	     100	  67451353 ns/op	 7195029 B/op	  699246 allocs/op
BenchmarkSyncMap/3-4       	     100	  80630496 ns/op	 7194917 B/op	  699248 allocs/op
BenchmarkSyncMap/4-4       	     100	 115386835 ns/op	 7195906 B/op	  699254 allocs/op
PASS
ok  	map/concurrentmap	74.638s
*/
func BenchmarkConcurrentMap(b *testing.B) {
	cm := New()
	for i := 1; i <= num; i++ {
		cm.Set(fmt.Sprintf("k-%d", i), i)
	}
	b.ResetTimer()
	for i := 0; i < 5; i++ {
		b.Run(strconv.Itoa(i), func(b *testing.B) {
			b.N = bN
			wg := sync.WaitGroup{}
			wg.Add(2 * b.N)
			for i := 0; i < b.N; i++ {
				go func() {
					defer wg.Done()
					for i := 0; i < num; i++ {
						cm.Set(fmt.Sprintf("k-%d", rand.Intn(num)), i)
					}
				}()
				go func() {
					defer wg.Done()
					for i := 0; i < num; i++ {
						cm.Get(fmt.Sprintf("k-%d", rand.Intn(num)))
					}
				}()
			}
			wg.Wait()
		})
	}
}

func BenchmarkSyncMap(b *testing.B) {
	sm := sync.Map{}
	for i := 1; i <= num; i++ {
		sm.Store(fmt.Sprintf("k-%d", i), i)
	}
	b.ResetTimer()
	for i := 0; i < 5; i++ {
		b.Run(strconv.Itoa(i), func(b *testing.B) {
			b.N = bN
			wg := sync.WaitGroup{}
			wg.Add(2 * b.N)
			for i := 0; i < b.N; i++ {
				go func() {
					defer wg.Done()
					for i := 0; i < num; i++ {
						sm.Store(fmt.Sprintf("k-%d", rand.Intn(num)), i)
					}
				}()
				go func() {
					defer wg.Done()
					for i := 0; i < num; i++ {
						sm.Load(fmt.Sprintf("k-%d", rand.Intn(num)))
					}
				}()
			}
			wg.Wait()
		})
	}
}
