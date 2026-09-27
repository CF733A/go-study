package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(createUniqueText("bar bar bar foo foo baz"))
}

func createUniqueText(text string) string {
	parts := strings.Split(text, " ")
	myMap := make(map[string]int)
	newParts := []string{}
	for _, word := range parts {
		if myMap[word] == 0 {
			myMap[word] = 1
			newParts = append(newParts, word)
		}
	}
	return strings.Join(newParts, " ")
}