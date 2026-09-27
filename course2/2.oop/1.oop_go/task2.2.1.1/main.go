package main

import (
	"errors"
	"fmt"
)

type Accounter interface {
	Deposit(amount float64) error
	Withdraw(amount float64) error
	Balance() float64
}

type CurrentAccount struct {
	balance float64
}

type SavingsAccount struct {
	balance float64
}

func (c *CurrentAccount) Deposit(amount float64) error {
	if amount <= 0 {
		return errors.New("некорректная сумма пополнения")
	}
	c.balance += amount
	return nil
}

func (c *CurrentAccount) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("некорректная сумма снятия")
	}
	if amount > c.balance {
		return errors.New("недостаточно средств")
	}
	c.balance -= amount
	return nil
}

func (c *CurrentAccount) Balance() float64 {
	return c.balance
}

func (s *SavingsAccount) Deposit(amount float64) error {
	if amount <= 0 {
		return errors.New("некорректная сумма пополнения")
	}
	s.balance += amount
	return nil
}

func (s *SavingsAccount) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("некорректная сумма снятия")
	}
	if amount > s.balance || s.balance <= 500 {
		return errors.New("недостаточно средств")
	}
	s.balance -= amount
	return nil
}

func (s *SavingsAccount) Balance() float64 {
	return s.balance
}

func ProcessAccount(a Accounter) {
	err := a.Deposit(500)
	if err != nil {
		fmt.Println("Ошибка пополнения", err)
	}

	err = a.Withdraw(200)
	if err != nil {
		fmt.Println("Ошибка снятия", err)
	}

	fmt.Printf("Balance: %.2f\n", a.Balance())
}

func main() {
	c := &CurrentAccount{}
	s := &SavingsAccount{}
	ProcessAccount(c)
	ProcessAccount(s)
}
