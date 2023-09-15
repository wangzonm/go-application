package singleflight_test

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"golang.org/x/sync/singleflight"
)

var (
	testKey   = "test_key"
	wg        sync.WaitGroup
	sf        singleflight.Group
	groupNums = 5
)

func TestDemo(t *testing.T) {
	wg.Add(groupNums)
	for i := 0; i < groupNums; i++ {
		go func(idx int) {
			defer wg.Done()
			time.Sleep(50 * time.Millisecond)
			log.Println("groutine.no=", idx, " start to set")
			v, _, shared := sf.Do(testKey, func() (interface{}, error) {
				log.Println("groutine.no=", idx, " setting value")
				log.Println("groutine.no=", idx, " set value success")
				return "testValue", nil
			})
			log.Println("groutine.no=", idx, " get value=", v.(string), " shared=", shared)
		}(i)
	}
	wg.Wait()
}

func TestDoCall(t *testing.T) {
	doCall()
}

//
//  doCall
//  @Description: singleflight.doCall
//=== RUN   TestDoCall
//defer 2 normalReturn false recovered false
//defer 1 normalReturn true recovered false
//--- PASS: TestDoCall (0.00s)
//PASS
//
func doCall() {
	normalReturn := false
	recovered := false
	defer func() {
		fmt.Println("defer 1 normalReturn", normalReturn, "recovered", recovered)
	}()
	func() {
		defer func() {
			fmt.Println("defer 2 normalReturn", normalReturn, "recovered", recovered)
			normalReturn = true
		}()
	}()
	if !normalReturn {
		recovered = true
	}
}

func TestChan(t *testing.T) {
	wg.Add(groupNums)
	ctx := context.Background()
	done := make(chan struct{})
	for i := 1; i <= groupNums; i++ {
		go func(ctx context.Context, idx int) {
			defer wg.Done()
			time.Sleep(time.Duration(idx) * time.Second)
			log.Println("groutine.no=", idx, " sleep", idx, " sec")
		}(ctx, i)
	}
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		log.Println("work group done")
	case <-time.After(3 * time.Second):
		log.Println("work group timeout")
	}
}

func TestChan2(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	done := make(chan struct{})
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			time.Sleep(time.Duration(i) * time.Second)
			log.Println("groutine.no=", i, " sleep", i, " sec")
			if i == 3 {
				cancel()
			}
		}(i)
	}
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		log.Println("work done")
	case <-ctx.Done():
		log.Println("ctx done=", ctx.Err())
	}
}
