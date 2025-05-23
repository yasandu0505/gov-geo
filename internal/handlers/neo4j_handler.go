package handlers

import (
	apierrors "go-mysql-backend/internal/errors"
	"go-mysql-backend/internal/service"
	"net/http"
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
