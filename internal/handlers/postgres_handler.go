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
	query := r.URL.Query()

	limitStr := query.Get("limit")
	offsetStr := query.Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		respondWithError(w, apierrors.ErrInvalidInput)
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		respondWithError(w, apierrors.ErrInvalidInput)
		return
	}

	ministries, err := h.Service.GetMinistriesWithDepartmentsPaginated(limit, offset)
	if err != nil {
		respondWithError(w, apierrors.ErrInternal)
		return
	}

	respondWithJSON(w, http.StatusOK, ministries)
}

func (h *OrganizationHandler) CreateMinistry(w http.ResponseWriter, r *http.Request) {
	var ministry models.Ministry
	if err := json.NewDecoder(r.Body).Decode(&ministry); err != nil {
		respondWithError(w, apierrors.ErrInvalidInput)
		return
	}
	defer r.Body.Close()

	if err := validateMinistry(ministry); err != nil {
		respondWithError(w, err)
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
		respondWithError(w, apierrors.ErrInvalidInput)
		return
	}
	defer r.Body.Close()

	if err := validateDepartment(dept); err != nil {
		respondWithError(w, err)
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
	id, err := getIDFromRequest(r)
	if err != nil {
		respondWithError(w, apierrors.ErrInvalidInput)
		return
	}

	ministry, err := h.Service.GetMinistryByID(id)
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
	id, err := getIDFromRequest(r)
	if err != nil {
		respondWithError(w, apierrors.ErrInvalidInput)
		return
	}

	ministry, err := h.Service.GetMinistryByIDWithDepartments(id)
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
	id, err := getIDFromRequest(r)
	if err != nil {
		respondWithError(w, apierrors.ErrInvalidInput)
		return
	}

	dept, err := h.Service.GetDepartmentByID(id)
	if err != nil {
		respondWithError(w, apierrors.ErrDepartmentNotFound)
		return
	}
	if dept == nil {
		respondWithError(w, apierrors.ErrDepartmentNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, dept)
}

// Helper functions
func getIDFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	return strconv.Atoi(vars["id"])
}

func validateMinistry(ministry models.Ministry) error {
	if ministry.Name == "" {
		return apierrors.ErrMissingField
	}
	return nil
}

func validateDepartment(dept models.Department) error {
	if dept.Name == "" {
		return apierrors.ErrMissingField
	}
	if dept.MinistryID == 0 {
		return apierrors.ErrMissingField
	}
	return nil
}
