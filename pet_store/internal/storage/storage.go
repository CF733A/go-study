package storage

import (
	"context"
	"pet-store/internal/entities"
)

type PetRepository interface {
	Create(ctx context.Context, pet *entities.Pet) error
	GetByID(ctx context.Context, id int64) (*entities.Pet, error)
	Update(ctx context.Context, pet *entities.Pet) error
	Delete(ctx context.Context, id int64) error
	UpdateWithFormData(ctx context.Context, id int64, name string, status string) error
	FindByStatus(ctx context.Context, status entities.PetStatus) ([]*entities.Pet, error)
	AddImage(ctx context.Context, id int64, imageURL string) error
}

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	Update(ctx context.Context, username string, user *entities.User) error
	Delete(ctx context.Context, username string) error
	CreateMultiple(ctx context.Context, users []*entities.User) error
}

type OrderRepository interface {
	Create(ctx context.Context, order *entities.Order) error
	GetByID(ctx context.Context, id int64) (*entities.Order, error)
	Delete(ctx context.Context, id int64) error
}

type Repository struct {
	Pet   PetRepository
	User  UserRepository
	Order OrderRepository
}
