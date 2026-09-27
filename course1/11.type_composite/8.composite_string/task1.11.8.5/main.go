package main

import (
	"fmt"
)

func main(){
	fmt.Println(ReverseString("Hello, world!"))
	fmt.Println(ReverseString("12345"))
	fmt.Println(ReverseString("Привет мир!")) 
}

func ReverseString(text string) string {
    runes := []rune(text)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}