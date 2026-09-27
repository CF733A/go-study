package main

import "fmt"

func main() {
	fmt.Println(MaxDifference([]int{5,-6}))
}

func MaxDifference(numbers []int) int {
	if len(numbers) <= 1 {
		return 0
	} 
	max := numbers[0]
	min := numbers[0]
	for _, nums := range numbers {
		if max < nums {
			max = nums
		}
		if min > nums {
			min = nums
		}
	}
	return max - min
}
