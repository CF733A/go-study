package main

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
)

type Animal struct {
    Type   string
    Name   string
    Age    int
}

func getAnimals() []Animal {
    var anilist []Animal
    for i := 0; i < 3; i++ {
        anilist = append(anilist, Animal{
            Type: gofakeit.Animal(),
            Name: gofakeit.PetName(),
            Age: gofakeit.IntRange(1, 20),
        })
    } 
    return anilist
}

func preparePrint(animals []Animal) string {
    var result string
    for _, anim := range animals {
        result += fmt.Sprintf("Тип: %s, Имя: %s, Возраст: %d\n", anim.Type, anim.Name, anim.Age)
    }
    return result
}

func main(){
    animals := getAnimals()
    result := preparePrint(animals)
    fmt.Println(result)
}