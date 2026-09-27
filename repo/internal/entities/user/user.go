package user

import (
	"time"

	"github.com/google/uuid"
)


type User struct {
	ID        string     `db:"id"`
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	IsDeleted bool       `db:"is_deleted"`
}

func NewUser(name, email string) *User {
	return &User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}
}

func (u *User) MarkAsDeleted() {
	now := time.Now()
	u.UpdatedAt = now
	u.DeletedAt = &now
	u.IsDeleted = true
}

func (u *User) Update(name, email string) {
	u.Name = name
	u.Email = email
	u.UpdatedAt = time.Now()
}
