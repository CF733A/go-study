package main

import (
	"fmt"
)

func main() {
	bytesSl := getBytes("Привет, мир!")
	fmt.Println(bytesSl)

	runesSl := getRunes("Привет, мир!")
	fmt.Println(runesSl)
}

func getBytes(s string) []byte {
	return []byte(s)
}

func getRunes(s string) []rune {
	return []rune(s)
}
