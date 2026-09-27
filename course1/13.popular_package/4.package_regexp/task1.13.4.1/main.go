package main

import (
	"fmt"
	"regexp"
)

func main() {
	email := "test@example.com"
	valid := isValidEmail(email)
	if valid {
		fmt.Printf("%s является валидным email-адресом\n", email)
	} else {
		fmt.Printf("%s не является валидным email-адресом\n", email)
	}
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`[a-zA-Z0-9]+@[a-zA-Z0-9]+\.[a-zA-Z]{2,}`)
	return re.MatchString(email)
}
