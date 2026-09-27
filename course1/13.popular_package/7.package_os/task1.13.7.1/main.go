package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	err := WriteFile("path/to/file.txt", []byte("Hello, World!"), os.FileMode(0644))
	if err != nil {
    	fmt.Println(err)
		return
	}


}

func WriteFile(filePath string, data []byte, perm os.FileMode) error {
    dirpath := filepath.Dir(filePath)
	err := os.MkdirAll(dirpath, 0755)
	if err != nil {
		fmt.Println(err)
		return err
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	file.Write(data)

	return nil
}
