package main

import (
	"errors"
	"fmt"
)

type Order interface {
	AddItem(item string, quantity int) error
	RemoveItem(item string) error
	GetOrderDetails() map[string]int
}

type DineInOrder struct {
	orderDetails map[string]int
}

type TakeAwayOrder struct {
	orderDetails map[string]int
}

func ManageOrder(o Order) {
	o.AddItem("Pizza", 2)
	o.AddItem("Burger", 1)
	o.AddItem("Pizza", 2)
	o.RemoveItem("Pizza")
	fmt.Println(o.GetOrderDetails())
}

func (d *DineInOrder) AddItem(item string, quantity int) error {
	if quantity > 0 {
		d.orderDetails[item] += quantity
		fmt.Println(d.orderDetails)
		return nil
	} else {
		return errors.New("the quantity must be greater than 0")
	}
}

func (d *DineInOrder) RemoveItem(item string) error {
	if _, exists := d.orderDetails[item]; exists {
		delete(d.orderDetails, item)
		fmt.Println(d.orderDetails)
		return nil
	} else {
		return errors.New("this item is not in the order")
	}
}
func (d *DineInOrder) GetOrderDetails() map[string]int {
	return d.orderDetails
}

func (t *TakeAwayOrder) AddItem(item string, quantity int) error {
	if quantity > 0 {
		t.orderDetails[item] += quantity
		fmt.Println(t.orderDetails)
		return nil
	} else {
		return errors.New("the quantity must be greater than 0")
	}
}

func (t *TakeAwayOrder) RemoveItem(item string) error {
	if _, exists := t.orderDetails[item]; exists {
		delete(t.orderDetails, item)
		fmt.Println(t.orderDetails)
		return nil
	} else {
		return errors.New("this item is not in the order")
	}
}
func (t *TakeAwayOrder) GetOrderDetails() map[string]int {
	return t.orderDetails
}

func main() {
	dineIn := &DineInOrder{orderDetails: make(map[string]int)}
	takeAway := &TakeAwayOrder{orderDetails: make(map[string]int)}

	ManageOrder(dineIn)
	ManageOrder(takeAway)
}
