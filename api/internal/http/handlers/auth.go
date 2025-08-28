package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"

	"bailanysta/api/internal/pkg/logger"
	"bailanysta/api/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
	logger      *logger.Logger
	validator   *validator.Validate
}

func NewAuthHandler(authService *services.AuthService, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
		validator:   validator.New(),
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req services.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode register request", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Register validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Register user
	response, err := h.authService.Register(r.Context(), req)
	if err != nil {
		h.logger.Error("Registration failed", map[string]interface{}{
			"error": err.Error(),
			"email": req.Email,
		})
		h.respondWithError(w, err.Error(), http.StatusConflict)
		return
	}

	h.logger.Info("User registered successfully", map[string]interface{}{
		"user_id":  response.User.ID,
		"username": response.User.Username,
	})

	h.respondWithJSON(w, response, http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req services.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode login request", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Login validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Login user
	response, err := h.authService.Login(r.Context(), req)
	if err != nil {
		h.logger.Warn("Login failed", map[string]interface{}{
			"error": err.Error(),
			"email": req.Email,
		})
		h.respondWithError(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	h.logger.Info("User logged in successfully", map[string]interface{}{
		"user_id":  response.User.ID,
		"username": response.User.Username,
	})

	h.respondWithJSON(w, response, http.StatusOK)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement refresh token functionality
	h.logger.Info("Refresh token requested")
	h.respondWithError(w, "Not implemented yet", http.StatusNotImplemented)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// For stateless JWT, logout is handled client-side
	// In production, you might want to blacklist tokens
	h.logger.Info("User logged out")
	h.respondWithJSON(w, map[string]interface{}{"message": "Logged out successfully"}, http.StatusOK)
}

func (h *AuthHandler) respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *AuthHandler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	h.respondWithJSON(w, map[string]interface{}{
		"error": map[string]interface{}{
			"code":    getErrorCode(statusCode),
			"message": message,
		},
	}, statusCode)
}

func getErrorCode(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusConflict:
		return "CONFLICT"
	case http.StatusInternalServerError:
		return "INTERNAL_SERVER_ERROR"
	default:
		return "UNKNOWN_ERROR"
	}
}
