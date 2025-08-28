package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"bailanysta/api/internal/pkg/auth"
	"bailanysta/api/internal/pkg/logger"
	"bailanysta/api/internal/services"
)

type UsersHandler struct {
	authService   *services.AuthService
	socialService *services.SocialService
	logger        *logger.Logger
	jwtManager    *auth.JWTManager
}

func NewUsersHandler(authService *services.AuthService, socialService *services.SocialService, logger *logger.Logger, jwtManager *auth.JWTManager) *UsersHandler {
	return &UsersHandler{
		authService:   authService,
		socialService: socialService,
		logger:        logger,
		jwtManager:    jwtManager,
	}
}

func (h *UsersHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		h.respondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.GetCurrentUser(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get current user", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		h.respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, user, http.StatusOK)
}

func (h *UsersHandler) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		h.respondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Bio       *string `json:"bio"`
		AvatarURL *string `json:"avatar_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Implement user update
	h.logger.Info("User update requested", map[string]interface{}{
		"user_id": userID,
	})

	h.respondWithJSON(w, map[string]interface{}{
		"message": "User update functionality not implemented yet",
	}, http.StatusOK)
}

func (h *UsersHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	currentUserID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		h.respondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse pagination parameters
	limit := 20
	offset := 0

	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
		if parsedOffset, err := strconv.Atoi(offsetParam); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Get total count
	var total int
	err = h.authService.GetDB().QueryRow(r.Context(), "SELECT COUNT(*) FROM users WHERE id != $1", currentUserID).Scan(&total)
	if err != nil {
		h.respondWithError(w, "Failed to get users count", http.StatusInternalServerError)
		return
	}

	// Get users with follow stats
	rows, err := h.authService.GetDB().Query(r.Context(), `
		SELECT u.id, u.username, u.email, u.bio, u.avatar_url,
		       COALESCE(f.followers_count, 0), COALESCE(ff.following_count, 0),
		       CASE WHEN fl.follower_id IS NOT NULL THEN true ELSE false END as is_following
		FROM users u
		LEFT JOIN (
		    SELECT followee_id, COUNT(*) as followers_count
		    FROM follows GROUP BY followee_id
		) f ON u.id = f.followee_id
		LEFT JOIN (
		    SELECT follower_id, COUNT(*) as following_count
		    FROM follows GROUP BY follower_id
		) ff ON u.id = ff.follower_id
		LEFT JOIN follows fl ON fl.followee_id = u.id AND fl.follower_id = $1
		WHERE u.id != $1
		ORDER BY u.username
		LIMIT $2 OFFSET $3`, currentUserID, limit, offset)
	if err != nil {
		h.respondWithError(w, "Failed to get users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []*services.UserResponse
	for rows.Next() {
		var user services.UserResponse
		var bio, avatarURL *string
		var followersCount, followingCount int

		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &bio, &avatarURL,
			&followersCount, &followingCount, &user.IsFollowing)
		if err != nil {
			h.respondWithError(w, "Failed to scan user", http.StatusInternalServerError)
			return
		}

		if bio != nil {
			user.Bio = *bio
		}
		user.AvatarURL = avatarURL
		user.FollowersCount = followersCount
		user.FollowingCount = followingCount

		users = append(users, &user)
	}

	response := map[string]interface{}{
		"users":  users,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}

	h.respondWithJSON(w, response, http.StatusOK)
}

func (h *UsersHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userIDParam := chi.URLParam(r, "id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		h.respondWithError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	currentUserID, _ := h.getUserIDFromContext(r.Context())

	// Get basic user info
	var user services.UserResponse
	var bio, avatarURL string
	err = h.authService.GetDB().QueryRow(r.Context(), `
		SELECT username, email, bio, avatar_url
		FROM users WHERE id = $1`, userID).Scan(
		&user.Username, &user.Email, &bio, &avatarURL)
	if err != nil {
		h.respondWithError(w, "User not found", http.StatusNotFound)
		return
	}

	user.ID = userID
	user.Bio = bio
	user.AvatarURL = &avatarURL

	// Get follow stats
	stats, err := h.socialService.GetFollowStats(r.Context(), userID, currentUserID)
	if err != nil {
		h.logger.Error("Failed to get follow stats", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		user.FollowersCount = stats.FollowersCount
		user.FollowingCount = stats.FollowingCount
		user.IsFollowing = stats.IsFollowing
	}

	h.respondWithJSON(w, user, http.StatusOK)
}

func (h *UsersHandler) respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *UsersHandler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	h.respondWithJSON(w, map[string]interface{}{
		"error": map[string]interface{}{
			"code":    getErrorCode(statusCode),
			"message": message,
		},
	}, statusCode)
}

func (h *UsersHandler) getUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	return h.jwtManager.GetUserIDFromContext(ctx)
}
