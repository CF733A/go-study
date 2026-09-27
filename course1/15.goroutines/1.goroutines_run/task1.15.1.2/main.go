package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	data := NotifyEvery(ticker, 5*time.Second, "Таймер сработал")

	for v := range data {
		fmt.Println(v)
	}

	fmt.Println("Программа завершена")
}

func NotifyEvery(ticker *time.Ticker, d time.Duration, message string) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		timeout := time.After(d)
		for {
			
			select {
			case <-ticker.C:
				ch <- message
			case <-timeout:
				return
			}
		}
	}()

	return ch
}
