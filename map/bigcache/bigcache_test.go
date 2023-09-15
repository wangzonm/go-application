package bigcache

import (
	"fmt"
	"github.com/allegro/bigcache/v2"
	"runtime"
	"runtime/debug"
	"testing"
	"time"
)

func TestBigCacheGcPause(t *testing.T) {
	entries := 20000000
	valueSize := 100
	repeat := 50
	debug.SetGCPercent(10)
	fmt.Println("Number of entries: ", entries)
	fmt.Println("Number of repeats: ", repeat)
	fmt.Println("Value size:        ", valueSize)
	bigCache(entries, valueSize)
	fmt.Println("GC pause for startup: ", gcPause())
	for i := 0; i < repeat; i++ {
		bigCache(entries, valueSize)
	}

	fmt.Printf("GC pause for %s: %s\n", "bigcache", gcPause())
}

var previousPause time.Duration

func gcPause() time.Duration {
	runtime.GC()
	var stats debug.GCStats
	debug.ReadGCStats(&stats)
	pause := stats.PauseTotal - previousPause
	previousPause = stats.PauseTotal
	return pause
}

func bigCache(entries, valueSize int) {
	config := bigcache.Config{
		Shards:             256,
		LifeWindow:         100 * time.Minute,
		MaxEntriesInWindow: entries,
		MaxEntrySize:       200,
		Verbose:            true,
	}

	bigcache, _ := bigcache.NewBigCache(config)
	for i := 0; i < entries; i++ {
		key, val := generateKeyValue(i, valueSize)
		bigcache.Set(key, val)
	}

	firstKey, _ := generateKeyValue(1, valueSize)
	v, err := bigcache.Get(firstKey)
	checkFirstElement(valueSize, v, err)
}

func stdMap(entries, valSize int) {
	mapCache := make(map[string][]byte)
	for i := 0; i < entries; i++ {
		key, val := generateKeyValue(i, valSize)
		mapCache[key] = val
	}
}

func generateKeyValue(index int, valSize int) (string, []byte) {
	key := fmt.Sprintf("key-%010d", index)
	fixedNumber := []byte(fmt.Sprintf("%010d", index))
	val := append(make([]byte, valSize-10), fixedNumber...)

	return key, val
}

func checkFirstElement(valueSize int, val []byte, err error) {
	_, expectedVal := generateKeyValue(1, valueSize)
	if err != nil {
		fmt.Println("Error in get: ", err.Error())
	} else if string(val) != string(expectedVal) {
		fmt.Println("Wrong first element: ", string(val))
	}
}
