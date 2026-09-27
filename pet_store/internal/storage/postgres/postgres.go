package postgres

import (
	"pet-store/internal/storage"
	"pet-store/internal/storage/postgres/order"
	"pet-store/internal/storage/postgres/pet"
	"pet-store/internal/storage/postgres/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresRepository(db *pgxpool.Pool) *storage.Repository {
	petR := pet.NewPetRepository(db)
	userR := user.NewUserRepository(db)
	orderR := order.NewOrderRepository(db)

	return &storage.Repository{
		Pet:   petR,
		User:  userR,
		Order: orderR,
	}
}
