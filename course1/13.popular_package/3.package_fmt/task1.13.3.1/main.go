package main

import "fmt"

// Пример кода, если необходимо
func generateMathString(operands []int, operator string) string {
	result := 0
	for _, val := range operands {
		result += val
	}
    return fmt.Sprintf("%d %s %d %s %d = %d", operands[0], operator, operands[1], operator, operands[2], result)
}

// Пример результата выполнения программы:
func main() {
    fmt.Println(generateMathString([]int{2, 4, 6}, "+")) // "2 + 4 + 6 = 12"
}