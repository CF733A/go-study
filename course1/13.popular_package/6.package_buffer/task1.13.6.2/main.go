package main

import (
	"bufio"
	"bytes"
)


func main() {
	data := []byte("Hello\n,\n World!")
	buffer := bytes.NewBuffer(data)

	scanner := getScanner(buffer)

	if scanner == nil {
		panic("Expected non-nil reader, got nil")
	}
	for scanner.Scan() {
		println(scanner.Text())
	}
}

func getScanner(b *bytes.Buffer) *bufio.Scanner {
    return bufio.NewScanner(b)
}