package handlers

import (
	"encoding/json"
	"go-mysql-backend/internal/service"
	"net/http"
)

type OrganizationHandler struct {
	Service *service.OrganizationService
}

func NewOrganizationHandler(service *service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{Service: service}
}

func (h *OrganizationHandler) GetMinistriesWithDepartments(w http.ResponseWriter, r *http.Request) {
	ministries, err := h.Service.GetMinistriesWithDepartments()
	if err != nil {
		http.Error(w, "Error fetching ministries", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ministries)
}
