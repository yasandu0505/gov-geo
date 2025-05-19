package routes

import (
	"go-mysql-backend/internal/handlers"
	"net/http"
)

func SetupOrgRoutes(OrganizationHandler *handlers.OrganizationHandler) {
	http.HandleFunc("/ministries", OrganizationHandler.GetMinistriesWithDepartments)
	http.HandleFunc("/ministries/create", OrganizationHandler.CreateMinistry)
}
