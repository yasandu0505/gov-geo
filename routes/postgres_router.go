package routes

import (
	"go-mysql-backend/internal/handlers"

	"github.com/gorilla/mux"
)

func SetupOrgRoutes(router *mux.Router, OrganizationHandler *handlers.OrganizationHandler) {
	router.HandleFunc("/ministries", OrganizationHandler.GetMinistriesWithDepartments).Methods("GET")
	router.HandleFunc("/ministries/paginated", OrganizationHandler.GetMinistriesWithDepartmentsPaginated).Methods("GET")
	router.HandleFunc("/ministries", OrganizationHandler.CreateMinistry).Methods("POST")
	router.HandleFunc("/departments", OrganizationHandler.CreateDepartment).Methods("POST")
	router.HandleFunc("/departments", OrganizationHandler.GetAllDepartments).Methods("GET")
	router.HandleFunc("/ministries/{id}", OrganizationHandler.GetMinistryByIDWithDepartments).Methods("GET")
	router.HandleFunc("/departments/{id}", OrganizationHandler.GetDepartmentByID).Methods("GET")

}
