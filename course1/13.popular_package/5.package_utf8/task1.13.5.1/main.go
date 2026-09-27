package main

import (
	"fmt"
)

func main() {
	myStr := "Hello, 世界!"
	fmt.Println(countUniqueUTF8Chars(myStr))
}

func countUniqueUTF8Chars(s string) int {
	mymap := make(map[rune]int)
	for _, r := range s {
		mymap[r]++
	}
	return len(mymap)
}
