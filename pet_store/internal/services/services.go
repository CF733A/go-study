package service

import (
	"context"
	"pet-store/internal/entities"
)

type JWTManager interface {
	Generate(userID int64, username string) (string, error)
	Validate(token string) (int64, error)
	RevokeToken(token string) error
	RevokeALLToken(userID int64) error
}

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
	Logout(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, token string) (int64, error)
}

type PetService interface {
	CreatePet(ctx context.Context, pet *entities.Pet) error
	GetPetByID(ctx context.Context, id int64) (*entities.Pet, error)
	UpdatePet(ctx context.Context, pet *entities.Pet) error
	DeletePet(ctx context.Context, id int64) error
	UpdatePetWithForm(ctx context.Context, id int64, name, status string) error
	FindPetsByStatus(ctx context.Context, status entities.PetStatus) ([]*entities.Pet, error)
	UploadPetImage(ctx context.Context, id int64, imageURL string) error
}

type UserService interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetUserByUsername(ctx context.Context, username string) (*entities.User, error)
	UpdateUser(ctx context.Context, username string, user *entities.User) error
	DeleteUser(ctx context.Context, username string) error
	CreateUsers(ctx context.Context, users []*entities.User) error
}

type OrderService interface {
	CreateOrder(ctx context.Context, order *entities.Order) error
	GetOrderByID(ctx context.Context, id int64) (*entities.Order, error)
	DeleteOrder(ctx context.Context, id int64) error
	GetInventory(ctx context.Context) (map[string]int32, error)
}

type Service struct {
	Pet   PetService
	User  UserService
	Order OrderService
	Auth  AuthService
	JWT   JWTManager
}
