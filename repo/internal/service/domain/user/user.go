package user

import (
	"context"
	"fmt"
	"log"
	"repo/internal/entities/user"
	"repo/internal/storage"
)

type UserService struct {
	UserRepo storage.UserRepository
}

func NewUserService(repo storage.UserRepository) *UserService {
	return &UserService{
		UserRepo: repo,
	}
}

func (u *UserService) CreateUser(ctx context.Context, name, email string) (*user.User, error) {
	existingUser, err := u.UserRepo.GetByEmail(ctx, email)
    if err != nil {
        return nil, fmt.Errorf("error checking email existence: %w", err)
    }
    if existingUser != nil {
        return nil, fmt.Errorf("user with email %s already exists", email)
    }
	
	usr := user.NewUser(name, email)

	if err := u.UserRepo.Create(ctx, usr); err != nil {
		return nil, fmt.Errorf("error when creating a user ")
	}

	log.Println("user created")
	return usr, nil
}

func (u *UserService) GetUser(ctx context.Context, id string) (*user.User, error) {
	usr, err := u.UserRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (u *UserService) UpdateUser(ctx context.Context, id, name, email string) (*user.User, error) {
	usr, err := u.UserRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	usr.Update(name, email)

	if err := u.UserRepo.Update(ctx, usr); err != nil {
		return nil, fmt.Errorf("user update err")
	}

	return usr, nil
}

func (u *UserService) DeleteUser(ctx context.Context, id string) error {
	return u.UserRepo.Delete(ctx, id)
}

func (u *UserService) ListUsers(ctx context.Context, c storage.Conditions) ([]*user.User, int, error) {
	return u.UserRepo.List(ctx, c)
}
