package main

import (
	"fmt"
	"time"
)

func timeout(timeout time.Duration) func() bool {
	ch := make(chan struct{})
	timer := time.NewTimer(timeout)

	return func() bool {
		
		go func() {
			<-timer.C
			close(ch)
		}()

		select {
		case <-ch:
			return true
		default:
			return false
		}
	}
}

// Пример использования функции timeout
func main() {
	timeoutFunc := timeout(1 * time.Second)
	since := time.NewTimer(3 * time.Second)
	for {
		select {
		case <-since.C:
			fmt.Println("Функция не выполнена вовремя")
			return
		default:
			if timeoutFunc() {
				fmt.Println("Функция выполнена вовремя")
				return
			}
		}
	}
}

// Output: Функция выполнена вовремя
