package handlers

import (
	apierrors "go-mysql-backend/internal/errors"
	"go-mysql-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Neo4JHandler struct {
	Service *service.Neo4JService
}

func NewNeo4JHandler(service *service.Neo4JService) *Neo4JHandler {
	return &Neo4JHandler{Service: service}
}

func (h *Neo4JHandler) GetMinistriesWithDepartments(w http.ResponseWriter, r *http.Request) {
	ministries, err := h.Service.GetMinistriesWithDepartments()
	if err != nil {
		respondWithError(w, apierrors.ErrInternal)
		return
	}
	respondJSON(w, http.StatusOK, ministries)
}

func (h *Neo4JHandler) GetMinistryByIDWithDepartments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	ministryID, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, apierrors.ErrInvalidInput)
		return
	}
	ministries, err := h.Service.GetMinistryByIDWithDepartments(ministryID)
	if err != nil {
		respondWithError(w, apierrors.ErrInternal)
		return
	}
	respondJSON(w, http.StatusOK, ministries)
}
