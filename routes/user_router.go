package routes

import (
	"go-mysql-backend/internal/handlers"
	"net/http"
)

func SetupRoutes(userHandler *handlers.UserHandler) {
	http.HandleFunc("/users", userHandler.GetUsers)
}
