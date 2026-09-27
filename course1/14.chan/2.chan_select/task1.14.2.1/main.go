package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	value := 15
	go func() {
		result := trySend(c, value)
		fmt.Println("result:", result)
	}()
	fmt.Println(<-c)
	time.Sleep(1 * time.Second)
}

func trySend(ch chan int, v int) bool {
	select {
	case ch <- v:
		return true
	default:
		return false
	}

}
