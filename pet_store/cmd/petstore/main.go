// @title PetStore API
// @version 0.1
// @description REST API for PetStore application with user authentication and pet management

// @contact.name API Support
// @contact.url http://localhost:8080

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Token
package main

import (
	"pet-store/internal/server/http"
	"pet-store/run"
)


func main() {
	app, err := run.NewApp()
	if err != nil {
		panic(err)
	}

	server := http.NewServer(app.Config.ServerPort, app.Router)

	lifecycle := http.NewLifecycleManager(server)
	lifecycle.Run()
}
