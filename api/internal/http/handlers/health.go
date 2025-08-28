package handlers

import (
	"encoding/json"
	"net/http"

	"bailanysta/api/internal/pkg/logger"
)

type HealthHandler struct {
	Logger *logger.Logger
}

type HealthResponse struct {
	OK bool `json:"ok"`
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{OK: true}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.Logger.Error("Failed to encode health response", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Health check requested", map[string]interface{}{
		"path":   r.URL.Path,
		"method": r.Method,
	})
}
