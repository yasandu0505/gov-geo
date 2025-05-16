package routes

import (
	"go-mysql-backend/internal/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRouter() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/users", handler.GetUsers).Methods("GET")
	return r
}
