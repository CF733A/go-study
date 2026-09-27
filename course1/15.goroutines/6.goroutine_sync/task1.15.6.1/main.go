package main

import (
	"fmt"
	"sync"
)

func waitGroupExample(goroutines ...func() string) string {
	var result string
	var wg sync.WaitGroup
	var mutex sync.Mutex

	wg.Add(len(goroutines))

	for _, val := range goroutines {
		go func(val func() string) {
			defer wg.Done()
			mutex.Lock()
			result += val()+"\n"
			mutex.Unlock()
		}(val)

	}

	wg.Wait()

	return result
}

func main() {
	count := 1000
	goroutines := make([]func() string, count)

	for i := 0; i < count; i++ {
		j := i
		goroutines[i] = func() string {
			return fmt.Sprintf("goroutine %d", j)
		}
	}

	fmt.Println(waitGroupExample(goroutines...))
}
