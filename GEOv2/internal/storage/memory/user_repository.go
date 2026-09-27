package memory

import (
	"errors"
	"geo-service/internal/entities"
	"sync"
)

type UserRepository struct {
	users map[string]*entities.User
	mu    sync.RWMutex
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*entities.User),
	}
}

func (u *UserRepository) Create(user *entities.User) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if _, exists := u.users[user.Username]; exists {
		return errors.New("user already exists")
	}

	u.users[user.Username] = user
	return nil
}

func (u *UserRepository) FindByUsername(username string) (*entities.User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	user, exists := u.users[username]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}
