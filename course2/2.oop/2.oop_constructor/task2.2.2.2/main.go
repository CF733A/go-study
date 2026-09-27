package main

import (
	"fmt"
)

type User struct {
	ID       int
	Username string
	Email    string
	Role     string
}

type UserOption func(*User)

func NewUser(id int, opts ...UserOption) *User {
	user := &User{ID: id}
	for _, ops := range opts {
		ops(user)
	}
	return user
}

func WithUsername(username string) UserOption {
	return func(u *User) {
		u.Username = username
	}
}

func WithEmail(mail string) UserOption {
	return func(u *User) {
		u.Email = mail
	}
}

func WithRole(role string) UserOption {
	return func(u *User) {
		u.Role = role
	}
}

func main() {
	user := NewUser(1, WithUsername("testuser"), WithEmail("testuser@example.com"), WithRole("admin"))
	fmt.Printf("User: %+v\n", user)
}
