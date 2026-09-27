package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	value int
}

func (c *Counter) Increment() int {

	c.value++
	return c.value
}

func concurrentSafeCounter() int {
	counter := Counter{}
	mutex := sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		
		go func() {
			defer wg.Done()
			mutex.Lock()
			counter.Increment()
			mutex.Unlock()
		}()

	}
	wg.Wait()
	return counter.value
}

func main() {
	fmt.Println(concurrentSafeCounter())
}
