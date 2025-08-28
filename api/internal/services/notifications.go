package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationType string

const (
	NotificationTypeLike    NotificationType = "like"
	NotificationTypeComment NotificationType = "comment"
	NotificationTypeFollow  NotificationType = "follow"
	NotificationTypeMention NotificationType = "mention"
	NotificationTypeNewPost NotificationType = "new_post"
)

type NotificationService struct {
	db *pgxpool.Pool
}

type Notification struct {
	ID        uuid.UUID              `json:"id"`
	UserID    uuid.UUID              `json:"user_id"`
	Type      NotificationType       `json:"type"`
	EntityID  *uuid.UUID             `json:"entity_id"`
	Payload   map[string]interface{} `json:"payload"`
	ReadAt    *time.Time             `json:"read_at"`
	CreatedAt time.Time              `json:"created_at"`

	// Additional data for display
	Actor *UserResponse `json:"actor,omitempty"`
	Post  *Post         `json:"post,omitempty"`
}

type CreateNotificationRequest struct {
	UserID   uuid.UUID              `json:"user_id"`
	Type     NotificationType       `json:"type"`
	EntityID *uuid.UUID             `json:"entity_id"`
	Payload  map[string]interface{} `json:"payload"`
}

func NewNotificationService(db *pgxpool.Pool) *NotificationService {
	return &NotificationService{db: db}
}

func (s *NotificationService) CreateNotification(ctx context.Context, req CreateNotificationRequest) (*Notification, error) {
	payloadJSON, err := json.Marshal(req.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	var notification Notification
	var entityID pgtype.UUID
	if req.EntityID != nil {
		var bytes [16]byte
		copy(bytes[:], req.EntityID[:])
		entityID = pgtype.UUID{Bytes: bytes, Valid: true}
	}

	err = s.db.QueryRow(ctx, `
		INSERT INTO notifications (user_id, type, entity_id, payload_json)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, type, entity_id, payload_json, read_at, created_at`,
		req.UserID, req.Type, entityID, payloadJSON).Scan(
		&notification.ID, &notification.UserID, &notification.Type,
		&entityID, &payloadJSON, &notification.ReadAt, &notification.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	if entityID.Valid {
		entityUUID := uuid.UUID(entityID.Bytes)
		notification.EntityID = &entityUUID
	}

	err = json.Unmarshal(payloadJSON, &notification.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	return &notification, nil
}

func (s *NotificationService) GetUserNotifications(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Notification, error) {
	rows, err := s.db.Query(ctx, `
		SELECT n.id, n.user_id, n.type, n.entity_id, n.payload_json, n.read_at, n.created_at
		FROM notifications n
		WHERE n.user_id = $1
		ORDER BY n.created_at DESC
		LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*Notification
	for rows.Next() {
		var notification Notification
		var entityID pgtype.UUID
		var payloadJSON []byte

		err := rows.Scan(
			&notification.ID, &notification.UserID, &notification.Type,
			&entityID, &payloadJSON, &notification.ReadAt, &notification.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}

		if entityID.Valid {
			entityUUID := uuid.UUID(entityID.Bytes)
			notification.EntityID = &entityUUID
		}

		err = json.Unmarshal(payloadJSON, &notification.Payload)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		// Populate additional data based on notification type
		err = s.populateNotificationData(ctx, &notification)
		if err != nil {
			return nil, fmt.Errorf("failed to populate notification data: %w", err)
		}

		notifications = append(notifications, &notification)
	}

	return notifications, nil
}

func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID, userID uuid.UUID) error {
	result, err := s.db.Exec(ctx, `
		UPDATE notifications
		SET read_at = now()
		WHERE id = $1 AND user_id = $2 AND read_at IS NULL`, notificationID, userID)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("notification not found or already read")
	}

	return nil
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	_, err := s.db.Exec(ctx, `
		UPDATE notifications
		SET read_at = now()
		WHERE user_id = $1 AND read_at IS NULL`, userID)
	if err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}

	return nil
}

func (s *NotificationService) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := s.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM notifications
		WHERE user_id = $1 AND read_at IS NULL`, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return count, nil
}

func (s *NotificationService) DeleteNotification(ctx context.Context, notificationID, userID uuid.UUID) error {
	result, err := s.db.Exec(ctx, `
		DELETE FROM notifications
		WHERE id = $1 AND user_id = $2`, notificationID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

// Notification triggers - called when certain actions happen

func (s *NotificationService) NotifyLike(ctx context.Context, likerID, postID uuid.UUID) error {
	// Get post author
	var postAuthorID uuid.UUID
	var postText string
	err := s.db.QueryRow(ctx, `
		SELECT author_id, text FROM posts WHERE id = $1`, postID).Scan(&postAuthorID, &postText)
	if err != nil {
		return fmt.Errorf("failed to get post info: %w", err)
	}

	// Don't notify if user likes their own post
	if likerID == postAuthorID {
		return nil
	}

	payload := map[string]interface{}{
		"liker_id":  likerID,
		"post_id":   postID,
		"post_text": truncateText(postText, 100),
	}

	_, err = s.CreateNotification(ctx, CreateNotificationRequest{
		UserID:   postAuthorID,
		Type:     NotificationTypeLike,
		EntityID: &postID,
		Payload:  payload,
	})

	return err
}

func (s *NotificationService) NotifyComment(ctx context.Context, commenterID, postID uuid.UUID, commentText string) error {
	// Get post author
	var postAuthorID uuid.UUID
	var postText string
	err := s.db.QueryRow(ctx, `
		SELECT author_id, text FROM posts WHERE id = $1`, postID).Scan(&postAuthorID, &postText)
	if err != nil {
		return fmt.Errorf("failed to get post info: %w", err)
	}

	// Don't notify if user comments on their own post
	if commenterID == postAuthorID {
		return nil
	}

	payload := map[string]interface{}{
		"commenter_id": commenterID,
		"post_id":      postID,
		"comment_text": truncateText(commentText, 100),
		"post_text":    truncateText(postText, 100),
	}

	_, err = s.CreateNotification(ctx, CreateNotificationRequest{
		UserID:   postAuthorID,
		Type:     NotificationTypeComment,
		EntityID: &postID,
		Payload:  payload,
	})

	return err
}

func (s *NotificationService) NotifyFollow(ctx context.Context, followerID, followeeID uuid.UUID) error {
	payload := map[string]interface{}{
		"follower_id": followerID,
	}

	_, err := s.CreateNotification(ctx, CreateNotificationRequest{
		UserID:   followeeID,
		Type:     NotificationTypeFollow,
		EntityID: &followerID,
		Payload:  payload,
	})

	return err
}

func (s *NotificationService) NotifyNewPost(ctx context.Context, authorID, postID uuid.UUID, postText string) error {
	// Get all followers of the author
	rows, err := s.db.Query(ctx, `
		SELECT follower_id FROM follows WHERE followee_id = $1`, authorID)
	if err != nil {
		return fmt.Errorf("failed to get followers: %w", err)
	}
	defer rows.Close()

	var followerIDs []uuid.UUID
	for rows.Next() {
		var followerID uuid.UUID
		if err := rows.Scan(&followerID); err != nil {
			return fmt.Errorf("failed to scan follower ID: %w", err)
		}
		followerIDs = append(followerIDs, followerID)
	}

	// Create notifications for all followers
	for _, followerID := range followerIDs {
		payload := map[string]interface{}{
			"author_id": authorID,
			"post_id":   postID,
			"post_text": truncateText(postText, 100),
		}

		_, err = s.CreateNotification(ctx, CreateNotificationRequest{
			UserID:   followerID,
			Type:     NotificationTypeNewPost,
			EntityID: &postID,
			Payload:  payload,
		})
		if err != nil {
			// Log error but continue with other notifications
			fmt.Printf("Failed to create new post notification for user %s: %v\n", followerID, err)
		}
	}

	return nil
}

// Helper methods

func (s *NotificationService) populateNotificationData(ctx context.Context, notification *Notification) error {
	switch notification.Type {
	case NotificationTypeLike:
		return s.populateLikeData(ctx, notification)
	case NotificationTypeComment:
		return s.populateCommentData(ctx, notification)
	case NotificationTypeFollow:
		return s.populateFollowData(ctx, notification)
	case NotificationTypeNewPost:
		return s.populateNewPostData(ctx, notification)
	}
	return nil
}

func (s *NotificationService) populateLikeData(ctx context.Context, notification *Notification) error {
	if notification.EntityID == nil {
		return nil
	}

	likerID, ok := notification.Payload["liker_id"].(string)
	if !ok {
		return nil
	}

	likerUUID, err := uuid.Parse(likerID)
	if err != nil {
		return err
	}

	// Get liker info
	var liker UserResponse
	var bio, avatarURL pgtype.Text
	err = s.db.QueryRow(ctx, `
		SELECT username, email, bio, avatar_url
		FROM users WHERE id = $1`, likerUUID).Scan(
		&liker.Username, &liker.Email, &bio, &avatarURL)
	if err != nil {
		return err
	}

	liker.ID = likerUUID
	liker.Bio = getPgtypeTextValue(bio)
	liker.AvatarURL = getPgtypeTextPtr(avatarURL)

	notification.Actor = &liker

	// Get post info
	var post Post
	var courseID, moduleID pgtype.UUID
	var postBio, postAvatarURL pgtype.Text
	err = s.db.QueryRow(ctx, `
		SELECT p.id, p.author_id, p.text, p.course_id, p.module_id, p.created_at, p.updated_at,
		       u.username, u.email, u.bio, u.avatar_url
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE p.id = $1`, *notification.EntityID).Scan(
		&post.ID, &post.AuthorID, &post.Text, &courseID, &moduleID, &post.CreatedAt, &post.UpdatedAt,
		&post.Author.Username, &post.Author.Email, &postBio, &postAvatarURL)
	if err != nil {
		return err
	}

	if courseID.Valid {
		courseUUID := uuid.UUID(courseID.Bytes)
		post.CourseID = &courseUUID
	}
	if moduleID.Valid {
		moduleUUID := uuid.UUID(moduleID.Bytes)
		post.ModuleID = &moduleUUID
	}
	post.Author.Bio = getPgtypeTextValue(postBio)
	post.Author.AvatarURL = getPgtypeTextPtr(postAvatarURL)

	notification.Post = &post

	return nil
}

func (s *NotificationService) populateCommentData(ctx context.Context, notification *Notification) error {
	if notification.EntityID == nil {
		return nil
	}

	commenterID, ok := notification.Payload["commenter_id"].(string)
	if !ok {
		return nil
	}

	commenterUUID, err := uuid.Parse(commenterID)
	if err != nil {
		return err
	}

	// Get commenter info
	var commenter UserResponse
	var bio, avatarURL pgtype.Text
	err = s.db.QueryRow(ctx, `
		SELECT username, email, bio, avatar_url
		FROM users WHERE id = $1`, commenterUUID).Scan(
		&commenter.Username, &commenter.Email, &bio, &avatarURL)
	if err != nil {
		return err
	}

	commenter.ID = commenterUUID
	commenter.Bio = getPgtypeTextValue(bio)
	commenter.AvatarURL = getPgtypeTextPtr(avatarURL)

	notification.Actor = &commenter

	// Get post info (same as in populateLikeData)
	return s.populateLikeData(ctx, notification)
}

func (s *NotificationService) populateFollowData(ctx context.Context, notification *Notification) error {
	if notification.EntityID == nil {
		return nil
	}

	// Get follower info
	var follower UserResponse
	var bio, avatarURL pgtype.Text
	err := s.db.QueryRow(ctx, `
		SELECT username, email, bio, avatar_url
		FROM users WHERE id = $1`, *notification.EntityID).Scan(
		&follower.Username, &follower.Email, &bio, &avatarURL)
	if err != nil {
		return err
	}

	follower.ID = *notification.EntityID
	follower.Bio = getPgtypeTextValue(bio)
	follower.AvatarURL = getPgtypeTextPtr(avatarURL)

	notification.Actor = &follower

	return nil
}

func (s *NotificationService) populateNewPostData(ctx context.Context, notification *Notification) error {
	if notification.EntityID == nil {
		return nil
	}

	// Get post info
	var post Post
	var bio, avatarURL pgtype.Text
	err := s.db.QueryRow(ctx, `
		SELECT p.id, p.author_id, p.text, p.created_at, 
			   u.username, u.email, u.bio, u.avatar_url
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE p.id = $1`, *notification.EntityID).Scan(
		&post.ID, &post.AuthorID, &post.Text, &post.CreatedAt,
		&post.Author.Username, &post.Author.Email, &bio, &avatarURL)
	if err != nil {
		return fmt.Errorf("failed to get post info: %w", err)
	}

	// Convert pgtype to regular types
	post.Author.Bio = getPgtypeTextValue(bio)
	post.Author.AvatarURL = getPgtypeTextPtr(avatarURL)

	notification.Post = &post
	notification.Actor = &post.Author

	return nil
}

// Utility functions
func truncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength] + "..."
}
