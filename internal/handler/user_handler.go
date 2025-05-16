package handler

import (
	"encoding/json"
	"go-mysql-backend/internal/service"
	"net/http"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := service.FetchUsers()
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
