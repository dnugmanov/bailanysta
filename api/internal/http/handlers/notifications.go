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

type NotificationsHandler struct {
	notificationsService *services.NotificationService
	logger               *logger.Logger
	jwtManager           *auth.JWTManager
}

func NewNotificationsHandler(notificationsService *services.NotificationService, logger *logger.Logger, jwtManager *auth.JWTManager) *NotificationsHandler {
	return &NotificationsHandler{
		notificationsService: notificationsService,
		logger:               logger,
		jwtManager:           jwtManager,
	}
}

func (h *NotificationsHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		h.respondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	limit := 20
	offset := 0
	unreadOnly := false

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

	if unreadParam := r.URL.Query().Get("unread_only"); unreadParam == "true" {
		unreadOnly = true
	}

	var notifications []*services.Notification
	if unreadOnly {
		// Get only unread notifications (limit to recent ones)
		allNotifications, err := h.notificationsService.GetUserNotifications(r.Context(), userID, 50, 0)
		if err != nil {
			h.logger.Error("Failed to get notifications", map[string]interface{}{
				"error":   err.Error(),
				"user_id": userID,
			})
			h.respondWithError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Filter unread notifications
		for _, notification := range allNotifications {
			if notification.ReadAt == nil {
				notifications = append(notifications, notification)
			}
		}

		// Apply limit/offset to unread notifications
		if offset >= len(notifications) {
			notifications = []*services.Notification{}
		} else {
			end := offset + limit
			if end > len(notifications) {
				end = len(notifications)
			}
			notifications = notifications[offset:end]
		}
	} else {
		notifications, err = h.notificationsService.GetUserNotifications(r.Context(), userID, limit, offset)
		if err != nil {
			h.logger.Error("Failed to get notifications", map[string]interface{}{
				"error":   err.Error(),
				"user_id": userID,
			})
			h.respondWithError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	h.respondWithJSON(w, map[string]interface{}{
		"notifications": notifications,
		"limit":         limit,
		"offset":        offset,
		"unread_only":   unreadOnly,
	}, http.StatusOK)
}

func (h *NotificationsHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		h.respondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	notificationIDParam := chi.URLParam(r, "id")
	notificationID, err := uuid.Parse(notificationIDParam)
	if err != nil {
		h.respondWithError(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	err = h.notificationsService.MarkAsRead(r.Context(), notificationID, userID)
	if err != nil {
		h.logger.Error("Failed to mark notification as read", map[string]interface{}{
			"error":           err.Error(),
			"user_id":         userID,
			"notification_id": notificationID,
		})
		h.respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Notification marked as read", map[string]interface{}{
		"notification_id": notificationID,
		"user_id":         userID,
	})

	h.respondWithJSON(w, map[string]interface{}{
		"message": "Notification marked as read",
	}, http.StatusOK)
}

func (h *NotificationsHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		h.respondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.notificationsService.MarkAllAsRead(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to mark all notifications as read", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		h.respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("All notifications marked as read", map[string]interface{}{
		"user_id": userID,
	})

	h.respondWithJSON(w, map[string]interface{}{
		"message": "All notifications marked as read",
	}, http.StatusOK)
}

func (h *NotificationsHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		h.respondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	count, err := h.notificationsService.GetUnreadCount(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get unread count", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		h.respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, map[string]interface{}{
		"unread_count": count,
	}, http.StatusOK)
}

func (h *NotificationsHandler) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		h.respondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	notificationIDParam := chi.URLParam(r, "id")
	notificationID, err := uuid.Parse(notificationIDParam)
	if err != nil {
		h.respondWithError(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	err = h.notificationsService.DeleteNotification(r.Context(), notificationID, userID)
	if err != nil {
		h.logger.Error("Failed to delete notification", map[string]interface{}{
			"error":           err.Error(),
			"user_id":         userID,
			"notification_id": notificationID,
		})
		h.respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Notification deleted", map[string]interface{}{
		"notification_id": notificationID,
		"user_id":         userID,
	})

	h.respondWithJSON(w, map[string]interface{}{
		"message": "Notification deleted successfully",
	}, http.StatusOK)
}

func (h *NotificationsHandler) respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *NotificationsHandler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	h.respondWithJSON(w, map[string]interface{}{
		"error": map[string]interface{}{
			"code":    getErrorCode(statusCode),
			"message": message,
		},
	}, statusCode)
}

func (h *NotificationsHandler) getUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	return h.jwtManager.GetUserIDFromContext(ctx)
}
