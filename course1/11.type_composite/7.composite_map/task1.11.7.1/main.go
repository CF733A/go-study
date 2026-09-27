package main

import (
	"fmt"
	"strings"
)

func countWordOccurrences(text string) map[string]int {
	parts := strings.Split(text, " ")
	countMap := make(map[string]int)
	for _, val := range parts {
		countMap[val]++
	}
	return countMap
}

func main() {
    text := "Lorem ipsum dolor sit amet consectetur adipiscing elit ipsum"
    occurrences := countWordOccurrences(text)

    for word, count := range occurrences {
        fmt.Printf("%s: %d\n", word, count)
    }
}