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

type SocialHandler struct {
	socialService *services.SocialService
	logger        *logger.Logger
	jwtManager    *auth.JWTManager
}

func NewSocialHandler(socialService *services.SocialService, logger *logger.Logger, jwtManager *auth.JWTManager) *SocialHandler {
	return &SocialHandler{
		socialService: socialService,
		logger:        logger,
		jwtManager:    jwtManager,
	}
}

func (h *SocialHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		h.respondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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

	posts, err := h.socialService.GetFeed(r.Context(), userID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get feed", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		h.respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, map[string]interface{}{
		"posts":  posts,
		"limit":  limit,
		"offset": offset,
	}, http.StatusOK)
}

func (h *SocialHandler) GetCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := h.socialService.GetCourses(r.Context())
	if err != nil {
		h.logger.Error("Failed to get courses", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, map[string]interface{}{
		"courses": courses,
	}, http.StatusOK)
}

func (h *SocialHandler) GetModulesByCourse(w http.ResponseWriter, r *http.Request) {
	courseIDParam := chi.URLParam(r, "id")
	courseID, err := uuid.Parse(courseIDParam)
	if err != nil {
		h.respondWithError(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	modules, err := h.socialService.GetModulesByCourse(r.Context(), courseID)
	if err != nil {
		h.logger.Error("Failed to get modules", map[string]interface{}{
			"error":     err.Error(),
			"course_id": courseID,
		})
		h.respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, map[string]interface{}{
		"modules": modules,
	}, http.StatusOK)
}

func (h *SocialHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	followerID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		h.respondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userIDParam := chi.URLParam(r, "id")
	followeeID, err := uuid.Parse(userIDParam)
	if err != nil {
		h.respondWithError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.socialService.FollowUser(r.Context(), followerID, followeeID)
	if err != nil {
		h.logger.Error("Failed to follow user", map[string]interface{}{
			"error":       err.Error(),
			"follower_id": followerID,
			"followee_id": followeeID,
		})
		h.respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info("User followed successfully", map[string]interface{}{
		"follower_id": followerID,
		"followee_id": followeeID,
	})

	h.respondWithJSON(w, map[string]interface{}{
		"message": "User followed successfully",
	}, http.StatusOK)
}

func (h *SocialHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	followerID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		h.respondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userIDParam := chi.URLParam(r, "id")
	followeeID, err := uuid.Parse(userIDParam)
	if err != nil {
		h.respondWithError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.socialService.UnfollowUser(r.Context(), followerID, followeeID)
	if err != nil {
		h.logger.Error("Failed to unfollow user", map[string]interface{}{
			"error":       err.Error(),
			"follower_id": followerID,
			"followee_id": followeeID,
		})
		h.respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info("User unfollowed successfully", map[string]interface{}{
		"follower_id": followerID,
		"followee_id": followeeID,
	})

	h.respondWithJSON(w, map[string]interface{}{
		"message": "User unfollowed successfully",
	}, http.StatusOK)
}

func (h *SocialHandler) respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *SocialHandler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	h.respondWithJSON(w, map[string]interface{}{
		"error": map[string]interface{}{
			"code":    getErrorCode(statusCode),
			"message": message,
		},
	}, statusCode)
}

func (h *SocialHandler) getUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	return h.jwtManager.GetUserIDFromContext(ctx)
}
