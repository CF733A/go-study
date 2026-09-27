package user

import (
	"context"
	"errors"
	"fmt"
	"pet-store/internal/entities"
	"pet-store/internal/storage"
)

type UserService struct {
	UserRepo storage.UserRepository
}

func NewUserService(userRepo storage.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (u *UserService) CreateUser(ctx context.Context, user *entities.User) error {
	userExist, _ := u.UserRepo.GetByUsername(ctx, user.Username)
	if userExist != nil {
		return errors.New("user already exists")
	}

	if user.Username == "" {
		return errors.New("username is required")
	}

	if user.Email == "" {
		return errors.New("email is required")
	}

	return u.UserRepo.Create(ctx, user)
}

func (u *UserService) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	user, err := u.UserRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

func (u *UserService) UpdateUser(ctx context.Context, username string, user *entities.User) error {
	currentUser, err := u.UserRepo.GetByUsername(ctx, username)
	if err != nil {
		return err
	}

	if err := currentUser.Update(user); err != nil {
		return err
	}

	return u.UserRepo.Update(ctx, username, currentUser)
}

func (u *UserService) DeleteUser(ctx context.Context, username string) error {
	_, err := u.UserRepo.GetByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	return u.UserRepo.Delete(ctx, username)
}

func (u *UserService) CreateUsers(ctx context.Context, users []*entities.User) error {
	usernames := make(map[string]bool)
	for _, user := range users {
		if usernames[user.Username] {
			return fmt.Errorf("duplicate username in request: %s", user.Username)
		}
		usernames[user.Username] = true
	}

	for _, user := range users {
		existingUser, _ := u.UserRepo.GetByUsername(ctx, user.Username)
		if existingUser != nil {
			return fmt.Errorf("username already exists in database: %s", user.Username)
		}
	}

	return u.UserRepo.CreateMultiple(ctx, users)
}

