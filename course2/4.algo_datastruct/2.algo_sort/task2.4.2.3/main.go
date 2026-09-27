package main

import (
	"fmt"
	"sort"
)

// Структура пользователя
type User struct {
	ID   int
	Name string
	Age  int
}

type SortByID []User

func (s SortByID) Len() int {
	return len(s)
}

func (s SortByID) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortByID) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

// Функция слияния двух отсортированных массивов пользователей
func Merge(arr1 []User, arr2 []User) []User {
	arr1 = append(arr1, arr2...)
	sort.Sort(SortByID(arr1))
	return arr1
}

func main() {
	users1 := []User{
		{ID: 2, Name: "Алексей", Age: 30},
		{ID: 1, Name: "Иван", Age: 25},
		{ID: 4, Name: "Сергей", Age: 35},
		{ID: 3, Name: "Дмитрий", Age: 28},
		{ID: 5, Name: "Анна", Age: 22},
	}
	users2 := []User{
		{ID: 9, Name: "Елена", Age: 40},
		{ID: 7, Name: "Олег", Age: 32},
		{ID: 10, Name: "Мария", Age: 27},
		{ID: 6, Name: "Николай", Age: 29},
		{ID: 8, Name: "Татьяна", Age: 31},
	}
	mergedUsers := Merge(users1, users2)
	fmt.Println(mergedUsers)
}
