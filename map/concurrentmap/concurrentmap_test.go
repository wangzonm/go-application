package concurrentmap

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/orcaman/concurrent-map/v2"
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
BenchmarkConcurrentMap/0-4   	     100	  95876518 ns/op	 3999620 B/op	  499280 allocs/op
BenchmarkConcurrentMap/1-4   	     100	  64931553 ns/op	 3999169 B/op	  499275 allocs/op
BenchmarkConcurrentMap/2-4   	     100	  66856476 ns/op	 3998267 B/op	  499273 allocs/op
BenchmarkConcurrentMap/3-4   	     100	 117061303 ns/op	 3997555 B/op	  499259 allocs/op
BenchmarkConcurrentMap/4-4   	     100	  81363032 ns/op	 3998085 B/op	  499270 allocs/op
BenchmarkSyncMap/0-4         	     100	 100979086 ns/op	 7195090 B/op	  699244 allocs/op
BenchmarkSyncMap/1-4         	     100	  95963044 ns/op	 7195327 B/op	  699242 allocs/op
BenchmarkSyncMap/2-4         	     100	  84743453 ns/op	 7195498 B/op	  699250 allocs/op
BenchmarkSyncMap/3-4         	     100	  81156601 ns/op	 7194961 B/op	  699248 allocs/op
BenchmarkSyncMap/4-4         	     100	  76454971 ns/op	 7194985 B/op	  699246 allocs/op
BenchmarkConcurrentMapV2/0-4 	     100	  53158915 ns/op	 3200622 B/op	  399528 allocs/op
BenchmarkConcurrentMapV2/1-4 	     100	  53396233 ns/op	 3200128 B/op	  399528 allocs/op
BenchmarkConcurrentMapV2/2-4 	     100	  53511853 ns/op	 3200216 B/op	  399527 allocs/op
BenchmarkConcurrentMapV2/3-4 	     100	  55735606 ns/op	 3200328 B/op	  399529 allocs/op
BenchmarkConcurrentMapV2/4-4 	     100	  53184668 ns/op	 3200624 B/op	  399529 allocs/op
PASS
ok  	map/concurrentmap	114.788s
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

func BenchmarkConcurrentMapV2(b *testing.B) {
	m := cmap.New[int]()
	for i := 1; i <= num; i++ {
		m.Set(fmt.Sprintf("k-%d", i), i)
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
						m.Set(fmt.Sprintf("k-%d", rand.Intn(num)), i)
					}
				}()
				go func() {
					defer wg.Done()
					for i := 0; i < num; i++ {
						_, _ = m.Get(fmt.Sprintf("k-%d", rand.Intn(num)))
					}
				}()
			}
			wg.Wait()
		})
	}
}
