package main

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"sort"
	"time"
)

type Product struct {
	Name      string
	Price     float64
	CreatedAt time.Time
	Count     int
}

func (p Product) String() string {
	return fmt.Sprintf("Name: %s, Price: %f, Count: %v", p.Name, p.Price, p.Count)
}

func generateProducts(n int) []Product {
	gofakeit.Seed(time.Now().UnixNano())
	products := make([]Product, n)
	for i := range products {
		products[i] = Product{
			Name:      gofakeit.Word(),
			Price:     gofakeit.Price(1.0, 100.0),
			CreatedAt: gofakeit.Date(),
			Count:     gofakeit.Number(1, 100),
		}
	}
	return products
}

type ByPrice []Product

func (bp ByPrice) Len() int {
	return len(bp)
}

func (bp ByPrice) Swap(i, j int) {
	bp[i], bp[j] = bp[j], bp[i]
}

func (bp ByPrice) Less(i, j int) bool {
	return bp[i].Price < bp[j].Price
}

type ByCreatedAt []Product

func (ca ByCreatedAt) Len() int {
	return len(ca)
}

func (ca ByCreatedAt) Swap(i, j int) {
	ca[i], ca[j] = ca[j], ca[i]
}

func (ca ByCreatedAt) Less(i, j int) bool {
	return ca[i].CreatedAt.Before(ca[j].CreatedAt)
}

type ByCount []Product

func (c ByCount) Len() int {
	return len(c)
}

func (c ByCount) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c ByCount) Less(i, j int) bool {
	return c[i].Count < c[j].Count
}

func main() {
	products := generateProducts(10)

	fmt.Println("Исходный список:")
	fmt.Println(products)

	// Сортировка продуктов по цене
	sort.Sort(ByPrice(products))
	fmt.Println("\nОтсортировано по цене:")
	fmt.Println(products)

	// Сортировка продуктов по дате создания
	sort.Sort(ByCreatedAt(products))
	fmt.Println("\nОтсортировано по дате создания:")
	fmt.Println(products)

	// Сортировка продуктов по количеству
	sort.Sort(ByCount(products))
	fmt.Println("\nОтсортировано по количеству:")
	fmt.Println(products)
}
