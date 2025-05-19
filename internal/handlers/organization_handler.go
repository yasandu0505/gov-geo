package handlers

import (
	"encoding/json"
	"go-mysql-backend/internal/models"
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

func (h *OrganizationHandler) CreateMinistry(w http.ResponseWriter, r *http.Request) {
	var ministry models.Ministry
	if err := json.NewDecoder(r.Body).Decode(&ministry); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	id, err := h.Service.CreateMinistry(ministry)
	if err != nil {
		http.Error(w, "Error creating ministry", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Ministry created successfully",
		"id":      id,
	})
}

func (h *OrganizationHandler) GetAllDepartments(w http.ResponseWriter, r *http.Request) {
	departments, err := h.Service.GetAllDepartments()
	if err != nil {
		http.Error(w, "Error fetching departments", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(departments)
}

func (h *OrganizationHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var dept models.Department
	if err := json.NewDecoder(r.Body).Decode(&dept); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check for required fields
	if dept.Name == "" || dept.MinistryID == 0 {
		http.Error(w, "Missing department name or ministry_id", http.StatusBadRequest)
		return
	}

	id, err := h.Service.CreateDepartment(dept)
	if err != nil {
		http.Error(w, "Error creating department", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Department created successfully",
		"id":      id,
	})
}
