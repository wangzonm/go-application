package channel_test

import (
	"sync"
	"testing"
)

func BenchmarkBatchConsumer(b *testing.B) {
	eventQueue = make(chan interface{}, 4)
	for i := 0; i < b.N; i++ {
		for i := 0; i < workerSize; i++ {
			go func() {
				var batch []interface{}
				for data := range eventQueue {
					batch = append(batch, data)
					if len(batch) == batchSize {
						batchProcessor(batch)
						batch = make([]interface{}, 0)
					}
				}
			}()
		}
		for i := 0; i < eventSize; i++ {
			eventQueue <- i
		}
	}
}

func BenchmarkBatchConsumerByWg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		eventQueue = make(chan interface{}, 4)
		for i := 0; i < workerSize; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var batch []interface{}
				for data := range eventQueue {
					batch = append(batch, data)
					if len(batch) == batchSize {
						batchProcessor(batch)
						batch = make([]interface{}, 0)
					}
				}
			}()
		}

		for i := 0; i < eventSize; i++ {
			eventQueue <- i
		}
		close(eventQueue)
		wg.Wait()
	}
}
