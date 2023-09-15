package list_test

import (
	"container/list"
	"testing"
)

var (
	ordersList *list.List
	ordersMap  map[int]int
	TotalCount = 1000000
	WriteCount = 10000
)

func init() {
	ordersList = list.New()
	for i := 0; i < TotalCount; i++ {
		ordersList.PushBack(i)
	}
	//for e := ordersList.Front(); e != nil; e = e.Next() {
	//	fmt.Println(e.Value)
	//}
	ordersMap = make(map[int]int)
	for i := 0; i < TotalCount; i++ {
		ordersMap[i] = i
	}
}

func BenchmarkList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for e := ordersList.Front(); e != nil; e = e.Next() {
			if e.Value.(int)-WriteCount < 0 {
				e.Value = -(e.Value).(int)
			}
		}
	}
}

func BenchmarkMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for k, v := range ordersMap {
			if v-WriteCount < 0 {
				ordersMap[k] = -v
			}
		}
	}
}
