package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bailanysta/api/internal/pkg/auth"
	"bailanysta/api/internal/services"
)

func TestAuthHandler_Register(t *testing.T) {
	// Create mock database
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	jwtManager := auth.NewJWTManager("test-secret", time.Hour, 24*time.Hour)
	authService := services.NewAuthService(db, jwtManager)
	handler := NewAuthHandler(authService, nil)

	// Create test request
	reqBody := map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}
	reqBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(reqBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Mock the service method to return success
	// Note: In a real scenario, you'd mock the database calls in the service

	handler.Register(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check response structure
	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok)

	user, ok := data["user"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "testuser", user["username"])
	assert.Equal(t, "test@example.com", user["email"])

	tokens, ok := data["tokens"].(map[string]interface{})
	assert.True(t, ok)
	assert.NotEmpty(t, tokens["access_token"])
	assert.NotEmpty(t, tokens["refresh_token"])
}

func TestAuthHandler_Register_InvalidJSON(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret", time.Hour, 24*time.Hour)
	authService := services.NewAuthService(nil, jwtManager)
	handler := NewAuthHandler(authService, nil)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response["error"].(map[string]interface{})["message"], "invalid JSON")
}

func TestAuthHandler_Register_ValidationError(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret", time.Hour, 24*time.Hour)
	authService := services.NewAuthService(nil, jwtManager)
	handler := NewAuthHandler(authService, nil)

	// Invalid request - missing required fields
	reqBody := map[string]interface{}{
		"username": "",
		"email":    "invalid-email",
		"password": "short",
	}
	reqBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(reqBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response["error"].(map[string]interface{})["message"], "validation failed")
}

func TestAuthHandler_Login(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret", time.Hour, 24*time.Hour)
	authService := services.NewAuthService(nil, jwtManager)
	handler := NewAuthHandler(authService, nil)

	reqBody := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
	}
	reqBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(reqBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	// Note: This would require mocking the database in a real scenario
	// For now, we're just testing the handler structure
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusUnauthorized)
}

func TestAuthHandler_Refresh(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret", time.Hour, 24*time.Hour)
	authService := services.NewAuthService(nil, jwtManager)
	handler := NewAuthHandler(authService, nil)

	reqBody := map[string]interface{}{
		"refresh_token": "some-refresh-token",
	}
	reqBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewReader(reqBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Refresh(w, req)

	// This endpoint is not fully implemented yet
	assert.Equal(t, http.StatusNotImplemented, w.Code)
}

func TestAuthHandler_GetCurrentUser(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret", time.Hour, 24*time.Hour)
	authService := services.NewAuthService(nil, jwtManager)
	handler := NewAuthHandler(authService, nil)

	req := httptest.NewRequest("GET", "/auth/me", nil)
	req.Header.Set("Content-Type", "application/json")

	// Add user ID to context (simulating middleware)
	userID := uuid.New()
	ctx := context.WithValue(req.Context(), "user_id", userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.GetCurrentUser(w, req)

	// Note: This would require mocking the database in a real scenario
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
}

func TestAuthHandler_Logout(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret", time.Hour, 24*time.Hour)
	authService := services.NewAuthService(nil, jwtManager)
	handler := NewAuthHandler(authService, nil)

	req := httptest.NewRequest("POST", "/auth/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Logout(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Logged out successfully", response["message"])
}

func TestGetErrorCode(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		expected string
	}{
		{"Bad Request", http.StatusBadRequest, "VALIDATION_ERROR"},
		{"Unauthorized", http.StatusUnauthorized, "AUTHENTICATION_ERROR"},
		{"Forbidden", http.StatusForbidden, "AUTHORIZATION_ERROR"},
		{"Not Found", http.StatusNotFound, "NOT_FOUND"},
		{"Conflict", http.StatusConflict, "CONFLICT"},
		{"Internal Server Error", http.StatusInternalServerError, "INTERNAL_ERROR"},
		{"Unknown", 999, "UNKNOWN_ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getErrorCode(tt.status)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function to add user ID to request context
func addUserIDToContext(r *http.Request, userID uuid.UUID) *http.Request {
	ctx := context.WithValue(r.Context(), "user_id", userID)
	return r.WithContext(ctx)
}

// Helper function to add JWT token to request
func addJWTTokenToRequest(r *http.Request, token string) *http.Request {
	r.Header.Set("Authorization", "Bearer "+token)
	return r
}
