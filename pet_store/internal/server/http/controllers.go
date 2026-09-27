package http

import (
	authHandler "pet-store/internal/server/http/handler/auth"
	orderHandler "pet-store/internal/server/http/handler/order"
	petHandler "pet-store/internal/server/http/handler/pet"
	userHandler "pet-store/internal/server/http/handler/user"
	service "pet-store/internal/services"
)

type Controllers struct {
	Auth  *authHandler.AuthHandler
	User  *userHandler.UserHandler
	Pet   *petHandler.PetHandler
	Order *orderHandler.OrderHandler
}

func NewControllers(services *service.Service) *Controllers {
	return &Controllers{
		Auth:  authHandler.NewAuthHandler(services.Auth),
		User:  userHandler.NewUserHandler(services.User),
		Pet:   petHandler.NewPetHandler(services.Pet),
		Order: orderHandler.NewOrderHandler(services.Order),
	}
}
