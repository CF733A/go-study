package main

import "fmt"

type Dish struct {
	Name  string
	Price float64
}

type Order struct {
	Dishes []Dish
	Total  float64
}

func (order *Order) AddDish(dish Dish) {
	order.Dishes = append(order.Dishes, Dish{
		Name:  dish.Name,
		Price: dish.Price,
	})
}

func (order *Order) RemoveDish(dish Dish) {
	for i, d := range order.Dishes {
		if dish.Name == d.Name {
			order.Dishes = append(order.Dishes[:i], order.Dishes[i+1:]...)
		}
	}
}

func (order *Order) CalculateTotal() {
	order.Total = 0
	for _, price := range order.Dishes {
		order.Total += price.Price
	}
}

func main() {
	order := Order{}
	dish1 := Dish{Name: "Pizza", Price: 10.99}
	dish2 := Dish{Name: "Burger", Price: 5.99}

	order.AddDish(dish1)
	order.AddDish(dish2)

	order.CalculateTotal()
	fmt.Println("Total:", order.Total)

	order.RemoveDish(dish1)

	order.CalculateTotal()
	fmt.Println("Total:", order.Total)
}
