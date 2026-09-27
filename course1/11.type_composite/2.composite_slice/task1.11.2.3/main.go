package main

import "fmt"

func bitwiseXOR(n, res int) int {
	return res ^ n
}

func findSingleNumber(numbers []int) int {
	result := 0
	for i := 0; i < len(numbers); i++ {
		result = bitwiseXOR(numbers[i], result)
	}
	return result
}

func main() {
	numbers := []int{1, 2, 3, 4, 5, 4, 3, 2, 1}
	singleNumber := findSingleNumber(numbers)
	fmt.Println(singleNumber) // 5
}
