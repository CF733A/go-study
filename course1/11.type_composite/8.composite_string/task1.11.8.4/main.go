package main

import (
	"fmt"
	"strings"
)

func main() {
	myString := concatStrings("Hello", " ", "world!")
	fmt.Println(myString)
}

func concatStrings(xs ...string) string {
	return strings.Join(xs, "")
}
