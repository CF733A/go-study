package main

import "fmt"

type Bank struct {
	queue []string
}

func (b *Bank) AddClient(client string) {
	b.queue = append(b.queue, client)
}

func (b *Bank) ServeNextClient() string {
	if len(b.queue) == 0 {
		return "No clients in the queue"
	}
	res := b.queue[0]
	b.queue = b.queue[1:]
	return res
}

func main() {
	bank := Bank{}

	bank.AddClient("Client 1")
	bank.AddClient("Client 2")
	bank.AddClient("Client 3")

	fmt.Println(bank.ServeNextClient()) // Output: Client 1
	fmt.Println(bank.ServeNextClient()) // Output: Client 2
	fmt.Println(bank.ServeNextClient()) // Output: Client 3
	fmt.Println(bank.ServeNextClient()) // Output: No clients in the queue
}
