package main

import (
	"fmt"
	"strings"
)

func main(){
	result := ReplaceSymbols("Hello, world!", 'o', '0')
	fmt.Println(result)
}

func ReplaceSymbols(text string, old rune, new rune)string{
	return strings.ReplaceAll(text, string(old), string(new))
}