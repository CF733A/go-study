package domain

import (
	"pet-store/config"
	auth "pet-store/internal/auth"
	service "pet-store/internal/services"
	auserv "pet-store/internal/services/domain/auth"
	"pet-store/internal/services/domain/order"
	"pet-store/internal/services/domain/pet"
	"pet-store/internal/services/domain/user"

	"pet-store/internal/storage"
)

func NewServices(repo *storage.Repository, cfg *config.Config) *service.Service {
	jwtManager := auth.NewJWTManager(cfg.JWTSecret)

	petService := pet.NewPetService(repo.Pet)
	userService := user.NewUserService(repo.User)
	orderService := order.NewOrderService(repo.Order, repo.Pet)

	authService := auserv.NewAuthService(userService, jwtManager)

	return &service.Service{
		Pet:   petService,
		User:  userService,
		Order: orderService,
		Auth:  authService,
		JWT:   jwtManager,
	}
}
