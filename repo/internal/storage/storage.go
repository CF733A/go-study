package storage

import (
	"context"
	"repo/internal/entities/user"
)

type Conditions struct {
	Limit  int
	Offset int
}

type UserRepository interface {
	Create(ctx context.Context, user *user.User) error
	GetByID(ctx context.Context, id string) (*user.User, error)
	Update(ctx context.Context, user *user.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, c Conditions) ([]*user.User, int, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
}

type Repository struct {
	User UserRepository
}

func NewRepository(userRepo UserRepository) *Repository {
	return &Repository{
		User: userRepo,
	}
}
