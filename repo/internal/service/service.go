package service

import (
	"context"
	"repo/internal/entities/user"
	"repo/internal/storage"
)

type UserService interface {
	CreateUser(ctx context.Context, name, email string) (*user.User, error)
	GetUser(ctx context.Context, id string) (*user.User, error)
	UpdateUser(ctx context.Context, id, name, email string) (*user.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, c storage.Conditions) ([]*user.User, int, error)
}

type Service struct {
	User UserService
}

func NewService(userServ UserService) *Service {
	return &Service{
		User: userServ,
	}
}
