package main

import (
	"fmt"
	"strings"
)

func filterSentence(sentence string, filter map[string]bool) string {
	parts := strings.Split(sentence, " ")
	filteredParts := []string{}
	for _, parts := range parts {
		if !filter[parts] {
			filteredParts = append(filteredParts, parts)
		}
	}

	return strings.Join(filteredParts, " ")
}

func main() {
    sentence := "Lorem ipsum dolor sit amet consectetur adipiscing elit ipsum"
    filter := map[string]bool{"ipsum": true, "elit": true}
    
    filteredSentence := filterSentence(sentence, filter)
    fmt.Println(filteredSentence)
}