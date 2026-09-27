package main

import "fmt"
//123
func sum(xs [8]int) int {
	summ := 0
	for i := 0; i < len(xs); i++ {
		summ += xs[i]
	}
	return summ
}

func average(xs [8]int) float64 {
	summ := 0
	for i := 0; i < len(xs); i++ {
		summ += xs[i]
	}
	return float64(summ) / float64(len(xs))
}

func averageFloat(xs [8]float64) float64 {
	var summ float64
	for i := 0; i < len(xs); i++ {
		summ += xs[i]
	}
	return summ / float64(len(xs))
}

func reverse(xs [8]int) [8]int {
	newArr := [8]int{}
	for i, j := 0, len(xs)-1; j > -1; i, j = i+1, j-1 {
		newArr[i] = xs[j]
	}
	return newArr
}

// Примеры использования функций:

func main() {
	xs := [8]int{1, 2, 3, 4, 5, 6, 7, 8}

	fmt.Println(sum(xs)) // Вывод: 36

	fmt.Println(average(xs)) // Вывод: 4.5

	ys := [8]float64{1, 2, 3, 4, 5, 6, 7, 8}

	fmt.Println(averageFloat(ys)) // 4.5

	fmt.Println(reverse(xs)) // Вывод: [8 7 6 5 4 3 2 1]
}
