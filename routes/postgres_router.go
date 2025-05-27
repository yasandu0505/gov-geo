package routes

import (
	"go-mysql-backend/internal/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupPostgresOrgRoutes(router *mux.Router, handler *handlers.OrganizationHandler) {
	// Create a subrouter for API v1
	v1 := router.PathPrefix("/api/v1").Subrouter()

	// Ministries routes
	ministries := v1.PathPrefix("/ministries").Subrouter()
	ministries.HandleFunc("", handler.GetMinistriesWithDepartments).Methods(http.MethodGet, http.MethodOptions)
	ministries.HandleFunc("/paginated", handler.GetMinistriesWithDepartmentsPaginated).Methods(http.MethodGet, http.MethodOptions)
	ministries.HandleFunc("", handler.CreateMinistry).Methods(http.MethodPost, http.MethodOptions)
	ministries.HandleFunc("/{id}", handler.GetMinistryByIDWithDepartments).Methods(http.MethodGet, http.MethodOptions)

	// Departments routes
	departments := v1.PathPrefix("/departments").Subrouter()
	departments.HandleFunc("", handler.GetAllDepartments).Methods(http.MethodGet, http.MethodOptions)
	departments.HandleFunc("", handler.CreateDepartment).Methods(http.MethodPost, http.MethodOptions)
	departments.HandleFunc("/{id}", handler.GetDepartmentByID).Methods(http.MethodGet, http.MethodOptions)
}
