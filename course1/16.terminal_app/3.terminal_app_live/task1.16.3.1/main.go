package main

import (
	"fmt"
	"time"

	"github.com/gosuri/uilive"
)

func main() {
	writer := uilive.New()
	writer.Start()
	defer writer.Stop()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		currentTime := time.Now()
		fmt.Fprintf(writer, "%s\n", currentTime.Format("15:04:05"))
		fmt.Fprintf(writer, "%s\n", currentTime.Format("2006-01-02"))
	}
}
