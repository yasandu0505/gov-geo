package handlers

import (
	"encoding/json"
	apierrors "go-mysql-backend/internal/errors"
	"go-mysql-backend/internal/models"
	"go-mysql-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
		respondWithError(w, apierrors.ErrInternal)
		return
	}
	respondWithJSON(w, http.StatusOK, ministries)
}

func (h *OrganizationHandler) GetMinistriesWithDepartmentsPaginated(w http.ResponseWriter, r *http.Request) {
	// Default values
	limit := 10
	offset := 0

	// Parse from query params
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil {
			limit = parsedLimit
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil {
			offset = parsedOffset
		}
	}

	ministries, err := h.Service.GetMinistriesWithDepartmentsPaginated(limit, offset)
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
		respondWithError(w, apierrors.ErrInvalidInput)
		return
	}

	if ministry.Name == "" {
		respondWithError(w, apierrors.ErrMissingField)
		return
	}

	id, err := h.Service.CreateMinistry(ministry)
	if err != nil {
		respondWithError(w, apierrors.ErrInternal)
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Ministry created successfully",
		"id":      id,
	})
}

func (h *OrganizationHandler) GetAllDepartments(w http.ResponseWriter, r *http.Request) {
	departments, err := h.Service.GetAllDepartments()
	if err != nil {
		respondWithError(w, apierrors.ErrInternal)
		return
	}
	respondWithJSON(w, http.StatusOK, departments)
}

func (h *OrganizationHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var dept models.Department
	if err := json.NewDecoder(r.Body).Decode(&dept); err != nil {
		respondWithError(w, apierrors.ErrDepartmentNotFound)
		return
	}

	if dept.Name == "" || dept.MinistryID == 0 {
		respondWithError(w, apierrors.ErrMissingField)
		return
	}

	id, err := h.Service.CreateDepartment(dept)
	if err != nil {
		respondWithError(w, apierrors.ErrInternal)
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Department created successfully",
		"id":      id,
	})
}

func (h *OrganizationHandler) GetMinistryByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	ministryID, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, apierrors.ErrInvalidInput)
		return
	}

	ministry, err := h.Service.GetMinistryByID(ministryID)
	if err != nil {
		respondWithError(w, apierrors.ErrMinistryNotFound)
		return
	}
	if ministry.ID == 0 {
		respondWithError(w, apierrors.ErrMinistryNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, ministry)
}

func (h *OrganizationHandler) GetMinistryByIDWithDepartments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	ministryID, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, apierrors.ErrInvalidInput)
		return
	}

	ministry, err := h.Service.GetMinistryByIDWithDepartments(ministryID)
	if err != nil {
		respondWithError(w, apierrors.ErrMinistryNotFound)
		return
	}
	if ministry.ID == 0 {
		respondWithError(w, apierrors.ErrMinistryNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, ministry)
}

func (h *OrganizationHandler) GetDepartmentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, apierrors.ErrInvalidInput)
		return
	}

	dept, err := h.Service.GetDepartmentByID(id)
	if err != nil {
		respondWithError(w, apierrors.ErrInternal)
		return
	}
	if dept == nil {
		respondWithError(w, apierrors.ErrDepartmentNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, dept)
}
