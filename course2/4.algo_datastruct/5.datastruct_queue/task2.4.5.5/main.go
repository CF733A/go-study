package main

import (
	"fmt"
)

type CircuitRinger interface {
	Add(val int)
	Get() (int, bool)
}

type RingBuffer struct {
	data     []int
	capacity int
}

func NewRingBuffer(cap int) *RingBuffer {
	return &RingBuffer{
		data:     make([]int, 0, cap),
		capacity: cap,
	}
}

func (rb *RingBuffer) Add(val int) {
	if len(rb.data) < rb.capacity {
		rb.data = append(rb.data, val)
	} else {
		rb.data = append(rb.data[1:], val)
	}

}

func (rb *RingBuffer) Get() (int, bool) {
	if len(rb.data) == 0 {
		return 0, false
	}
	val := rb.data[0]
	rb.data = rb.data[1:]
	return val, true
}

func main() {
	rb := NewRingBuffer(3)
	rb.Add(1)
	rb.Add(2)
	rb.Add(3)
	rb.Add(4) // Перезаписывает значение 1

	for val, ok := rb.Get(); ok; val, ok = rb.Get() {
		fmt.Println(val) // Выводит: 2, 3, 4
	}
	if _, ok := rb.Get(); !ok {
		fmt.Println("Buffer is empty") // Выводит: Buffer is empty
	}
}
