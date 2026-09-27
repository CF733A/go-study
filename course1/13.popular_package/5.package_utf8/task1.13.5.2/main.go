package main

import (
	"fmt"
	"strings"
)

func main() {
	result := countRussianLetters("Привет, мир!")
	for key, value := range result {
        fmt.Printf("%c: %d ", key, value) // в: 1 е: 1 т: 1 м: 1 п: 1 р: 2 и: 2 
    }
}

func countRussianLetters(s string) map[rune]int {
    counts := make(map[rune]int)
    for _, char := range strings.ToLower(s) {
        if isRussianLetter(char) {
            counts[char]++
        }
    }
	
    return counts
}

func isRussianLetter(char rune) bool {
	alphabet := "абвгдеёжзийклмнопрстуфхцчшщъыьэюя"
	for _, r := range alphabet{
		if char == r {
			return true
		}
	}
	return false
}