package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostsService struct {
	db                   *pgxpool.Pool
	notificationsService *NotificationService
}

type Post struct {
	ID           uuid.UUID    `json:"id"`
	AuthorID     uuid.UUID    `json:"author_id"`
	Text         string       `json:"text"`
	CourseID     *uuid.UUID   `json:"course_id,omitempty"`
	ModuleID     *uuid.UUID   `json:"module_id,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	LikeCount    int          `json:"like_count"`
	CommentCount int          `json:"comment_count"`
	Author       UserResponse `json:"author,omitempty"`
	IsLiked      bool         `json:"is_liked"`
}

type Comment struct {
	ID        uuid.UUID    `json:"id"`
	PostID    uuid.UUID    `json:"post_id"`
	AuthorID  uuid.UUID    `json:"author_id"`
	Text      string       `json:"text"`
	CreatedAt time.Time    `json:"created_at"`
	Author    UserResponse `json:"author,omitempty"`
}

type Like struct {
	UserID    uuid.UUID `json:"user_id"`
	PostID    uuid.UUID `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CreatePostRequest struct {
	Text     string     `json:"text" validate:"required,min=1,max=5000"`
	CourseID *uuid.UUID `json:"course_id,omitempty"`
	ModuleID *uuid.UUID `json:"module_id,omitempty"`
}

type UpdatePostRequest struct {
	Text     string     `json:"text" validate:"required,min=1,max=5000"`
	CourseID *uuid.UUID `json:"course_id,omitempty"`
	ModuleID *uuid.UUID `json:"module_id,omitempty"`
}

type CreateCommentRequest struct {
	Text string `json:"text" validate:"required,min=1,max=1000"`
}

func NewPostsService(db *pgxpool.Pool, notificationsService *NotificationService) *PostsService {
	return &PostsService{
		db:                   db,
		notificationsService: notificationsService,
	}
}

func (s *PostsService) CreatePost(ctx context.Context, userID uuid.UUID, req CreatePostRequest) (*Post, error) {
	var post Post

	// Extract hashtags from text
	hashtags := extractHashtags(req.Text)

	// Begin transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Create post
	err = tx.QueryRow(ctx, `
		INSERT INTO posts (author_id, text, course_id, module_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, author_id, text, course_id, module_id, created_at, updated_at`,
		userID, req.Text, req.CourseID, req.ModuleID).Scan(
		&post.ID, &post.AuthorID, &post.Text, &post.CourseID, &post.ModuleID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// Add hashtags
	for _, hashtag := range hashtags {
		// Insert or get hashtag
		var hashtagID uuid.UUID
		err = tx.QueryRow(ctx, `
			INSERT INTO hashtags (tag)
			VALUES ($1)
			ON CONFLICT (tag) DO UPDATE SET tag = EXCLUDED.tag
			RETURNING id`, hashtag).Scan(&hashtagID)
		if err != nil {
			return nil, fmt.Errorf("failed to create hashtag: %w", err)
		}

		// Link post to hashtag
		_, err = tx.Exec(ctx, `
			INSERT INTO post_hashtags (post_id, hashtag_id)
			VALUES ($1, $2)`, post.ID, hashtagID)
		if err != nil {
			return nil, fmt.Errorf("failed to link post to hashtag: %w", err)
		}
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	post.LikeCount = 0
	post.CommentCount = 0

	// Create notifications for followers
	if s.notificationsService != nil {
		err = s.notificationsService.NotifyNewPost(ctx, userID, post.ID, post.Text)
		if err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Failed to create new post notifications: %v\n", err)
		}
	}

	return &post, nil
}

func (s *PostsService) GetPostByID(ctx context.Context, postID uuid.UUID) (*Post, error) {
	var post Post
	var courseID, moduleID pgtype.UUID
	var bio, avatarURL pgtype.Text

	err := s.db.QueryRow(ctx, `
		SELECT p.id, p.author_id, p.text, p.course_id, p.module_id, p.created_at, p.updated_at,
		       COUNT(DISTINCT l.user_id) as like_count,
		       COUNT(DISTINCT c.id) as comment_count,
		       u.username, u.email, u.bio, u.avatar_url
		FROM posts p
		JOIN users u ON p.author_id = u.id
		LEFT JOIN likes l ON p.id = l.post_id
		LEFT JOIN comments c ON p.id = c.post_id
		WHERE p.id = $1
		GROUP BY p.id, u.username, u.email, u.bio, u.avatar_url`, postID).Scan(
		&post.ID, &post.AuthorID, &post.Text, &courseID, &moduleID, &post.CreatedAt, &post.UpdatedAt,
		&post.LikeCount, &post.CommentCount,
		&post.Author.Username, &post.Author.Email, &bio, &avatarURL)
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}

	// Convert pgtype to regular types
	if courseID.Valid {
		courseUUID := uuid.UUID(courseID.Bytes)
		post.CourseID = &courseUUID
	}
	if moduleID.Valid {
		moduleUUID := uuid.UUID(moduleID.Bytes)
		post.ModuleID = &moduleUUID
	}
	post.Author.Bio = getPgtypeTextValue(bio)
	post.Author.AvatarURL = getPgtypeTextPtr(avatarURL)

	return &post, nil
}

func (s *PostsService) UpdatePost(ctx context.Context, userID, postID uuid.UUID, req UpdatePostRequest) (*Post, error) {
	// Check if user owns the post
	var authorID uuid.UUID
	err := s.db.QueryRow(ctx, "SELECT author_id FROM posts WHERE id = $1", postID).Scan(&authorID)
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}
	if authorID != userID {
		return nil, fmt.Errorf("access denied")
	}

	// Update post
	var courseID, moduleID uuid.NullUUID
	if req.CourseID != nil {
		courseID = uuid.NullUUID{UUID: *req.CourseID, Valid: true}
	}
	if req.ModuleID != nil {
		moduleID = uuid.NullUUID{UUID: *req.ModuleID, Valid: true}
	}

	var post Post
	err = s.db.QueryRow(ctx, `
		UPDATE posts
		SET text = $1, course_id = $2, module_id = $3, updated_at = now()
		WHERE id = $4 AND author_id = $5
		RETURNING id, author_id, text, course_id, module_id, created_at, updated_at`,
		req.Text, courseID, moduleID, postID, userID).Scan(
		&post.ID, &post.AuthorID, &post.Text, &post.CourseID, &post.ModuleID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	// Get counts
	err = s.db.QueryRow(ctx, `
		SELECT COUNT(DISTINCT l.user_id), COUNT(DISTINCT c.id)
		FROM posts p
		LEFT JOIN likes l ON p.id = l.post_id
		LEFT JOIN comments c ON p.id = c.post_id
		WHERE p.id = $1
		GROUP BY p.id`, postID).Scan(&post.LikeCount, &post.CommentCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get counts: %w", err)
	}

	return &post, nil
}

func (s *PostsService) DeletePost(ctx context.Context, userID, postID uuid.UUID) error {
	// Check if user owns the post
	var authorID uuid.UUID
	err := s.db.QueryRow(ctx, "SELECT author_id FROM posts WHERE id = $1", postID).Scan(&authorID)
	if err != nil {
		return fmt.Errorf("post not found: %w", err)
	}
	if authorID != userID {
		return fmt.Errorf("access denied")
	}

	// Delete post (cascade will handle related records)
	_, err = s.db.Exec(ctx, "DELETE FROM posts WHERE id = $1 AND author_id = $2", postID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

func (s *PostsService) GetUserPosts(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Post, error) {
	rows, err := s.db.Query(ctx, `
		SELECT p.id, p.author_id, p.text, p.course_id, p.module_id, p.created_at, p.updated_at,
		       COUNT(DISTINCT l.user_id) as like_count,
		       COUNT(DISTINCT c.id) as comment_count,
		       u.username, u.email, u.bio, u.avatar_url
		FROM posts p
		JOIN users u ON p.author_id = u.id
		LEFT JOIN likes l ON p.id = l.post_id
		LEFT JOIN comments c ON p.id = c.post_id
		WHERE p.author_id = $1
		GROUP BY p.id, u.username, u.email, u.bio, u.avatar_url
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user posts: %w", err)
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID, &post.AuthorID, &post.Text, &post.CourseID, &post.ModuleID, &post.CreatedAt, &post.UpdatedAt,
			&post.LikeCount, &post.CommentCount,
			&post.Author.Username, &post.Author.Email, &post.Author.Bio, &post.Author.AvatarURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (s *PostsService) CreateComment(ctx context.Context, userID, postID uuid.UUID, req CreateCommentRequest) (*Comment, error) {
	var comment Comment
	err := s.db.QueryRow(ctx, `
		INSERT INTO comments (post_id, author_id, text)
		VALUES ($1, $2, $3)
		RETURNING id, post_id, author_id, text, created_at`,
		postID, userID, req.Text).Scan(
		&comment.ID, &comment.PostID, &comment.AuthorID, &comment.Text, &comment.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Get author info
	var bio, avatarURL pgtype.Text
	err = s.db.QueryRow(ctx, `
		SELECT username, email, bio, avatar_url
		FROM users WHERE id = $1`, userID).Scan(
		&comment.Author.Username, &comment.Author.Email, &bio, &avatarURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get author info: %w", err)
	}

	// Convert pgtype to regular types
	comment.Author.Bio = getPgtypeTextValue(bio)
	comment.Author.AvatarURL = getPgtypeTextPtr(avatarURL)

	// Create notification
	if s.notificationsService != nil {
		err = s.notificationsService.NotifyComment(ctx, userID, postID, req.Text)
		if err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Failed to create comment notification: %v\n", err)
		}
	}

	return &comment, nil
}

func (s *PostsService) GetComments(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*Comment, error) {
	rows, err := s.db.Query(ctx, `
		SELECT c.id, c.post_id, c.author_id, c.text, c.created_at,
		       u.username, u.email, u.bio, u.avatar_url
		FROM comments c
		JOIN users u ON c.author_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC
		LIMIT $2 OFFSET $3`, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		var comment Comment
		var bio, avatarURL pgtype.Text
		err := rows.Scan(
			&comment.ID, &comment.PostID, &comment.AuthorID, &comment.Text, &comment.CreatedAt,
			&comment.Author.Username, &comment.Author.Email, &bio, &avatarURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}

		// Convert pgtype to regular types
		comment.Author.Bio = getPgtypeTextValue(bio)
		comment.Author.AvatarURL = getPgtypeTextPtr(avatarURL)

		comments = append(comments, &comment)
	}

	return comments, nil
}

func (s *PostsService) LikePost(ctx context.Context, userID, postID uuid.UUID) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO likes (user_id, post_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, post_id) DO NOTHING`, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to like post: %w", err)
	}

	// Create notification (will be handled by trigger in production)
	if s.notificationsService != nil {
		err = s.notificationsService.NotifyLike(ctx, userID, postID)
		if err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Failed to create like notification: %v\n", err)
		}
	}

	return nil
}

func (s *PostsService) UnlikePost(ctx context.Context, userID, postID uuid.UUID) error {
	_, err := s.db.Exec(ctx, `
		DELETE FROM likes WHERE user_id = $1 AND post_id = $2`, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to unlike post: %w", err)
	}
	return nil
}

func (s *PostsService) IsPostLiked(ctx context.Context, userID, postID uuid.UUID) (bool, error) {
	var count int
	err := s.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM likes WHERE user_id = $1 AND post_id = $2`, userID, postID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check like status: %w", err)
	}
	return count > 0, nil
}

// Helper functions
func extractHashtags(text string) []string {
	re := regexp.MustCompile(`#\w+`)
	matches := re.FindAllString(text, -1)
	var hashtags []string
	for _, match := range matches {
		hashtags = append(hashtags, strings.TrimPrefix(match, "#"))
	}
	return hashtags
}

func getPgtypeTextValue(pt pgtype.Text) string {
	if pt.Valid {
		return pt.String
	}
	return ""
}

func getPgtypeTextPtr(pt pgtype.Text) *string {
	if pt.Valid {
		return &pt.String
	}
	return nil
}
