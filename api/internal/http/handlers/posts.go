package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"bailanysta/api/internal/pkg/auth"
	"bailanysta/api/internal/pkg/logger"
	"bailanysta/api/internal/services"
)

type PostsHandler struct {
	postsService *services.PostsService
	logger       *logger.Logger
	validator    *validator.Validate
	jwtManager   *auth.JWTManager
}

func NewPostsHandler(postsService *services.PostsService, logger *logger.Logger, jwtManager *auth.JWTManager) *PostsHandler {
	return &PostsHandler{
		postsService: postsService,
		logger:       logger,
		validator:    validator.New(),
		jwtManager:   jwtManager,
	}
}

func (h *PostsHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		h.respondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req services.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode create post request", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Create post validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	post, err := h.postsService.CreatePost(r.Context(), userID, req)
	if err != nil {
		h.logger.Error("Failed to create post", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		h.respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Post created successfully", map[string]interface{}{
		"post_id": post.ID,
		"user_id": userID,
	})

	h.respondWithJSON(w, post, http.StatusCreated)
}

func (h *PostsHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	postIDParam := chi.URLParam(r, "id")
	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		h.respondWithError(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.postsService.GetPostByID(r.Context(), postID)
	if err != nil {
		h.logger.Warn("Post not found", map[string]interface{}{
			"post_id": postID,
			"error":   err.Error(),
		})
		h.respondWithError(w, "Post not found", http.StatusNotFound)
		return
	}

	h.respondWithJSON(w, post, http.StatusOK)
}

func (h *PostsHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		h.respondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postIDParam := chi.URLParam(r, "id")
	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		h.respondWithError(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var req services.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode update post request", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Update post validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	post, err := h.postsService.UpdatePost(r.Context(), userID, postID, req)
	if err != nil {
		h.logger.Error("Failed to update post", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
			"post_id": postID,
		})
		if err.Error() == "access denied" {
			h.respondWithError(w, "Access denied", http.StatusForbidden)
		} else {
			h.respondWithError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info("Post updated successfully", map[string]interface{}{
		"post_id": postID,
		"user_id": userID,
	})

	h.respondWithJSON(w, post, http.StatusOK)
}

func (h *PostsHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		h.respondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postIDParam := chi.URLParam(r, "id")
	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		h.respondWithError(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	err = h.postsService.DeletePost(r.Context(), userID, postID)
	if err != nil {
		h.logger.Error("Failed to delete post", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
			"post_id": postID,
		})
		if err.Error() == "access denied" {
			h.respondWithError(w, "Access denied", http.StatusForbidden)
		} else {
			h.respondWithError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info("Post deleted successfully", map[string]interface{}{
		"post_id": postID,
		"user_id": userID,
	})

	h.respondWithJSON(w, map[string]interface{}{"message": "Post deleted successfully"}, http.StatusOK)
}

func (h *PostsHandler) LikePost(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		h.respondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postIDParam := chi.URLParam(r, "id")
	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		h.respondWithError(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	err = h.postsService.LikePost(r.Context(), userID, postID)
	if err != nil {
		h.logger.Error("Failed to like post", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
			"post_id": postID,
		})
		h.respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Post liked successfully", map[string]interface{}{
		"post_id": postID,
		"user_id": userID,
	})

	h.respondWithJSON(w, map[string]interface{}{"message": "Post liked successfully"}, http.StatusOK)
}

func (h *PostsHandler) UnlikePost(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		h.respondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postIDParam := chi.URLParam(r, "id")
	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		h.respondWithError(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	err = h.postsService.UnlikePost(r.Context(), userID, postID)
	if err != nil {
		h.logger.Error("Failed to unlike post", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
			"post_id": postID,
		})
		h.respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Post unliked successfully", map[string]interface{}{
		"post_id": postID,
		"user_id": userID,
	})

	h.respondWithJSON(w, map[string]interface{}{"message": "Post unliked successfully"}, http.StatusOK)
}

func (h *PostsHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	postIDParam := chi.URLParam(r, "id")
	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		h.respondWithError(w, "Invalid post ID", http.StatusBadRequest)
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

	comments, err := h.postsService.GetComments(r.Context(), postID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get comments", map[string]interface{}{
			"error":   err.Error(),
			"post_id": postID,
		})
		h.respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, map[string]interface{}{
		"comments": comments,
		"limit":    limit,
		"offset":   offset,
	}, http.StatusOK)
}

func (h *PostsHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		h.respondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postIDParam := chi.URLParam(r, "id")
	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		h.respondWithError(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var req services.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode create comment request", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Create comment validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	comment, err := h.postsService.CreateComment(r.Context(), userID, postID, req)
	if err != nil {
		h.logger.Error("Failed to create comment", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
			"post_id": postID,
		})
		h.respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Comment created successfully", map[string]interface{}{
		"comment_id": comment.ID,
		"post_id":    postID,
		"user_id":    userID,
	})

	h.respondWithJSON(w, comment, http.StatusCreated)
}

func (h *PostsHandler) respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *PostsHandler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	h.respondWithJSON(w, map[string]interface{}{
		"error": map[string]interface{}{
			"code":    getErrorCode(statusCode),
			"message": message,
		},
	}, statusCode)
}

func (h *PostsHandler) getUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	return h.jwtManager.GetUserIDFromContext(ctx)
}
