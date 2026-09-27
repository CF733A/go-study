package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
    filePath := "course1/13.popular_package/7.package_os/task1.13.7.2/file.txt"

    file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        return
    }
    defer file.Close()
    
    err = WriteFile(strings.NewReader("Hello, World!"), file)
    if err != nil {
        fmt.Println("Ошибка при записи файла:", err)
    }
}

func WriteFile(data io.Reader, fd io.Writer) error {
	_, err := io.Copy(fd, data)
	return err
}
