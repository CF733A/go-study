package main

import (
	"encoding/json"
	"fmt"
)

type Comment struct {
	Text string `json:"text"`
}

type User struct {
	Name     string    `json:"name"`
	Age      int       `json:"age"`
	Comments []Comment `json:"comments"`
}

func main() {
	users := []User{
		{
			Name: "Alice",
			Age:  25,
			Comments: []Comment{
				{Text: "Привет!"},
				{Text: "Как дела?"},
			},
		},
		{
			Name: "Bob",
			Age:  30,
			Comments: []Comment{
				{Text: "Пока!"},
			},
		},
	}
	fmt.Println(getJSON(users))
}

func getJSON(data []User) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "err", err
	}
	return string(jsonData), nil
}
