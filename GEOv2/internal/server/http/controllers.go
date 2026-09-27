package http

import (
	"geo-service/internal/server/http/handler/auth"
	"geo-service/internal/server/http/handler/geo"
	"geo-service/internal/server/http/handler/home"
	"geo-service/internal/service"
	"geo-service/pkg/responder"
)

type Controllers struct{
	Auth *auth.AuthHandler
	Geo *geo.GeoHandler
	Home *home.HomeHandler
}

func NewControllers(services *service.Service, resp responder.Responder) *Controllers {
	return &Controllers{
		Auth: auth.NewAuthHandler(services.Auth, resp),
		Geo: geo.NewGeoHandler(services.Geo, resp),
		Home: home.NewHomeHandler(resp),
	}
}