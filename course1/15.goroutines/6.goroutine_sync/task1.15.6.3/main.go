package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Пример структуры счетчика
type Counter struct {
	count int64
}

// Функция для увеличения значения счетчика на 1
func (c *Counter) Increment() {
	atomic.AddInt64(&c.count, 1)
}

// Функция для получения текущего значения счетчика
func (c *Counter) GetCount() int64 {
	return atomic.LoadInt64(&c.count)
}

func main() {
	counter := Counter{}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			counter.Increment()

		}(i)
		for j := 0; j < 3; j++ {
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				fmt.Printf("number inc: %d , number read: %d, counter: %d\n", i, j, counter.GetCount())
			}(i, j)
		}
	}
	wg.Wait()
}
