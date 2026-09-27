package main

import (
	"errors"
	"fmt"
)

type PaymentMethod interface {
	Pay(amount float64) error
}

type CreditCard struct {
	balance float64
}

func (c *CreditCard) Pay(amount float64) error {
	err := MakeAPayment(&c.balance, amount)
	if err != nil {
		return err
	}
	fmt.Printf("Оплачено %0.2f с помощью кредитной карты\n", amount)
	return nil
}

type Bitcoin struct {
	balance float64
}

func (b *Bitcoin) Pay(amount float64) error {
	err := MakeAPayment(&b.balance, amount)
	if err != nil {
		return err
	}
	fmt.Printf("Оплачено %0.2f с помощью биткоина\n", amount)
	return nil
}

func MakeAPayment(balance *float64, amount float64) error {
	if amount <= 0 {
		return errors.New("недопустимая сумма платежа")
	}

	if *balance < amount {
		return errors.New("недостаточный баланс")
	}

	*balance -= amount
	return nil
}

func ProcessPayment(p PaymentMethod, amount float64) {
	err := p.Pay(amount)
	if err != nil {
		fmt.Println("Не удалось обработать платеж:", err)
	}
}

func main() {
	cc := &CreditCard{balance: 500.00}
	btc := &Bitcoin{balance: 2.00}

	ProcessPayment(cc, 200.00)
	ProcessPayment(btc, 1.00)
}
