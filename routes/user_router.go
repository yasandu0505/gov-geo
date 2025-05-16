package routes

import (
	"go-mysql-backend/internal/handler"

	"github.com/gorilla/mux"
)

func RegisterUserRoutes(r *mux.Router, h *handler.UserHandler) {
	r.HandleFunc("/users", h.GetUsers).Methods("GET")
}
