package main

import (
	"fmt"
	"sync"
)

type UnsafeMap struct {
	mp    map[string]int
	mutex sync.Mutex
}

func NewUnsafeMap() *UnsafeMap {
	return &UnsafeMap{
		mp: make(map[string]int),
	}
}

func (sm *UnsafeMap) Set(key string, value int) {
	sm.mp[key] = value // (data race)
}

func (sm *UnsafeMap) Get(key string) (int, bool) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	res, ok := sm.mp[key]
	return res, ok
}

func main() {
	sm := NewUnsafeMap()
	wg := sync.WaitGroup{}

	count := 100
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			sm.Set(key, i)
		}(i)
	}

	wg.Wait()

	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key%d", i)
		if val, ok := sm.Get(key); ok {
			fmt.Printf("%s: %d\n", key, val)
		}
	}
}
