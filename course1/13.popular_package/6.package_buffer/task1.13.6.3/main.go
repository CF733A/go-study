package main

import (
	"bytes"
	"fmt"
)

func main() {
	buffer := bytes.NewBufferString("Hello, World!")

	result := getDataString(buffer)

	expected := "Hello, World!"
	if result != expected {
		panic(fmt.Sprintf("Expected %s, but got %s", expected, result))
	}
}

func getDataString(b *bytes.Buffer) string {
	return b.String()
}
