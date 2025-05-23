package routes

import (
	"go-mysql-backend/internal/handlers"

	"github.com/gorilla/mux"
)

func SetupNeo4JRoutes(router *mux.Router, Neo4JHandler *handlers.Neo4JHandler) {
	router.HandleFunc("/ministries", Neo4JHandler.GetMinistriesWithDepartments).Methods("GET")
	router.HandleFunc("/ministries/{id}", Neo4JHandler.GetMinistryByIDWithDepartments).Methods("GET")
	router.HandleFunc("/seed", Neo4JHandler.SeedDummyData).Methods("POST")

}
