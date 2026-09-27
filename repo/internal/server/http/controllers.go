package http

import (
	"repo/internal/server/http/handler/user"
	"repo/internal/service"
)

type Controllers struct {
	User *user.UserHandler
}

func NewControllers(service *service.Service) *Controllers {
	return &Controllers{
		User: user.NewUserHandler(service.User),
	}
}
