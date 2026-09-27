package main

import "fmt"

func main(){
	var i interface{} = 42
fmt.Println(getType(i)) // Вывод: "int"

var j interface{} = "Hello, World!"
fmt.Println(getType(j)) // Вывод: "string"

var k interface{} = []int{1, 2, 3}
fmt.Println(getType(k)) // Вывод: "[]int"

var l interface{} = interface{}(nil)
fmt.Println(getType(l)) // Вывод: "Пустой интерфейс"
}

func getType(val interface{})string{
	result := ""
	switch val.(type){
	case int: result = "int"
	case string: result = "string"
	case []int: result = "[]int"
	case nil: result = "Пустой интерфейс"
	}
	return result
}