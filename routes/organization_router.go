package routes

import (
	"go-mysql-backend/internal/handlers"

	"github.com/gorilla/mux"
)

func SetupOrgRoutes(router *mux.Router, OrganizationHandler *handlers.OrganizationHandler) {
	router.HandleFunc("/ministries", OrganizationHandler.GetMinistriesWithDepartments).Methods("GET")
	router.HandleFunc("/ministries", OrganizationHandler.CreateMinistry).Methods("POST")
	router.HandleFunc("/departments", OrganizationHandler.CreateDepartment).Methods("POST")
	router.HandleFunc("/ministries/{id}", OrganizationHandler.GetMinistryByID).Methods("GET")
}
