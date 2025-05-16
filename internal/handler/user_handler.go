package handler

import (
	"encoding/json"
	"net/http"

	"go-mysql-backend/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{s}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}
