package main

import "fmt"

var stack []int

func push(value int) {
	stack = append(stack, value)
}

func pop() int {
	value := stack[len(stack)-1]
	stack = stack[:len(stack)-1]
	return value
}

func main() {
	push(5)
	push(3)
	result := pop() + pop()
	push(result)
	fmt.Println(stack[0]) // 8
}