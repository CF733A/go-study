package main

import (
	"fmt"
)

type BrowserHistory struct {
	stack []string
}

func (h *BrowserHistory) Visit(url string) {
	h.stack = append(h.stack, url)
	fmt.Printf("Посещение %s\n", url)
}

func (h *BrowserHistory) Back() {
	if len(h.stack) == 0 {
		fmt.Println("Нет больше истории для возврата")
	}
	h.stack = h.stack[:len(h.stack)-1]
	fmt.Printf("Возврат к %s\n", h.stack[len(h.stack)-1])
}

func (h *BrowserHistory) PrintHistory() {
	for i := len(h.stack) - 1; i >= 0; i-- {
		fmt.Println(h.stack[i])
	}
}

func main() {
	history := &BrowserHistory{}
	history.Visit("www.google.com")
	history.Visit("www.github.com")
	history.Visit("www.openai.com")
	history.Back()
	history.PrintHistory()
}
