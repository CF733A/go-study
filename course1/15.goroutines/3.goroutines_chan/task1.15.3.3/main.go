package main

import (
	"fmt"
	"sync"
)

type sema chan struct{}

func New(n int) sema {
	ch := make(chan struct{}, n)
	for i := 0; i < n; i++ {
		ch <- struct{}{}
	}
	return ch
}

func (s sema) Inc(k int) {
	for i := 0; i < k; i++ {
		if len(s) < cap(s) {
			s <- struct{}{}
		}
	}
}

func (s sema) Dec(k int) {
	for i := 0; i < k; i++ {
		<-s
	}
}

func main() {
	numbers := []int{1, 2, 3, 4, 5}
	n := len(numbers)

	sem := New(n)
	var wg sync.WaitGroup
	for _, num := range numbers {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			sem.Dec(1)
			fmt.Println(n)
			sem.Inc(1)
		}(num)
	}
	wg.Wait()
	sem.Dec(n)
}
