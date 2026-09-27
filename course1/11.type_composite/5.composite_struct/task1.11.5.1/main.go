package main

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
)

type User struct {
	Name string
	Age  int
}

func main() {
	users := getUsers()
	result := preparePrint(users)
	fmt.Println(result)
}

func getUsers() []User {
	var fakeUsers []User
	for i := 0; i < 10; i++ {
		fakeUsers = append(fakeUsers, User{
			Name: gofakeit.Name(),
			Age:  gofakeit.IntRange(18, 60),
		})
	}
	return fakeUsers
}

func preparePrint(fakeUsers []User) string {
	var result string
	for _, user := range fakeUsers {
		result += fmt.Sprintf("Имя: %s, Возраст: %d\n", user.Name, user.Age)
	}
	return result
}
