package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {

	bytes := countBytes("Привет, мир!")
	fmt.Println(bytes)

	
	symbols := countSymbols("Привет, мир!")
	fmt.Println(symbols)
}

func countBytes(text string) int {
	return len(text)
}

func countSymbols(text string) int {
	return utf8.RuneCountInString(text)
}
