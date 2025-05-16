package main

import (
	"log"
	"net/http"

	"go-mysql-backend/config"
	"go-mysql-backend/internal/db"
	"go-mysql-backend/internal/handler"
	"go-mysql-backend/internal/repository"
	"go-mysql-backend/internal/service"
	"go-mysql-backend/routes"

	"github.com/gorilla/mux"
)

func main() {
	// Load config (optional if you're using a config package)
	config.Load()

	// Initialize DB
	db.Init()

	// Setup repository, service, and handler layers
	userRepo := repository.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Setup router and routes
	r := mux.NewRouter()
	routes.RegisterUserRoutes(r, userHandler)

	// Start server
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
