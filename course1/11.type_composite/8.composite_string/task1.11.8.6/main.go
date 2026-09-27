package main

import (
	"fmt"
	"strings"
)

func main() {
	count := CountVowels("Привет, мир!")
	fmt.Println(count)

	count = CountVowels("Hello, world!")
	fmt.Println(count)
}

func CountVowels(text string) int {
	count := 0
	vowels := "aeiouаеёиоуыэюя"
	for _, runes := range strings.ToLower(text) {
		if strings.Contains(vowels, string(runes)) {
			count++
		}
	}
	return count
}
