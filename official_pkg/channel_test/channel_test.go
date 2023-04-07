package channel_test

import (
	"sync"
	"testing"
)

var (
	eventSize = 100000
	eventQueue chan interface{}
	batchSize = 8
	workerSize = 3
	batchProcessor = func(message interface{}) {
		//fmt.Printf("%v\n", message)
	}
)

func TestBatchConsumer(t *testing.T) {
	eventQueue = make(chan interface{}, 4)
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
	close(eventQueue)
}

func TestBatchConsumerByWg(t *testing.T) {
	eventQueue = make(chan interface{}, 4)
	var wg sync.WaitGroup
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