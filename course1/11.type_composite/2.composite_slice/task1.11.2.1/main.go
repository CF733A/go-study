package main

import "fmt"

func getSubSlice(xs []int, start, end int) []int {
	subslice := xs[start:end]
	return subslice
}

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	subSlice := getSubSlice(numbers, 2, 6)
	fmt.Println(subSlice) // Вывод: [3 4 5 6]
}
