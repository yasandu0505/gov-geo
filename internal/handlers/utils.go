package handlers

import (
	"encoding/json"
	"net/http"

	apierrors "go-mysql-backend/internal/errors"
)

// respondWithError sends a structured error response
func respondWithError(w http.ResponseWriter, err error) {
	if apiErr, ok := err.(*apierrors.APIError); ok {
		respondWithJSON(w, apiErr.Code, map[string]string{"error": apiErr.Message})
	} else {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
}

// respondWithJSON sends a generic JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func respondJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
