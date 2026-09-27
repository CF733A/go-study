package entities

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Password   string `json:"-"`
	Phone      string `json:"phone"`
	UserStatus int32  `json:"userStatus"`
}

func NewUser(username, firstName, lastName, email, password, phone string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		Username:   username,
		FirstName:  firstName,
		LastName:   lastName,
		Email:      email,
		Password:   string(hashedPassword),
		Phone:      phone,
		UserStatus: 1,
	}, nil
}

func (u *User) Update(updatedUser *User) error {
	u.Username = updatedUser.Username
	u.FirstName = updatedUser.FirstName
	u.LastName = updatedUser.LastName
	u.Email = updatedUser.Email
	u.Phone = updatedUser.Phone
	u.UserStatus = updatedUser.UserStatus

	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}

	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) Activate() {
	u.UserStatus = 1
}

func (u *User) Deactivate() {
	u.UserStatus = 0
}

func (u *User) IsActive() bool {
	return u.UserStatus == 1
}
