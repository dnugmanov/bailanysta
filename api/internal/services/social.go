package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SocialService struct {
	db                   *pgxpool.Pool
	notificationsService *NotificationService
}

type FollowStats struct {
	FollowersCount int  `json:"followers_count"`
	FollowingCount int  `json:"following_count"`
	IsFollowing    bool `json:"is_following,omitempty"`
}

type FeedPost struct {
	ID           uuid.UUID    `json:"id"`
	AuthorID     uuid.UUID    `json:"author_id"`
	Text         string       `json:"text"`
	CourseID     *uuid.UUID   `json:"course_id,omitempty"`
	ModuleID     *uuid.UUID   `json:"module_id,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	LikeCount    int          `json:"like_count"`
	CommentCount int          `json:"comment_count"`
	Author       UserResponse `json:"author"`
	IsLiked      bool         `json:"is_liked"`
}

type FollowRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

func NewSocialService(db *pgxpool.Pool, notificationsService *NotificationService) *SocialService {
	return &SocialService{
		db:                   db,
		notificationsService: notificationsService,
	}
}

func (s *SocialService) FollowUser(ctx context.Context, followerID, followeeID uuid.UUID) error {
	if followerID == followeeID {
		return fmt.Errorf("cannot follow yourself")
	}

	// Check if already following
	var count int
	err := s.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM follows
		WHERE follower_id = $1 AND followee_id = $2`,
		followerID, followeeID).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check follow status: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("already following this user")
	}

	// Create follow relationship
	_, err = s.db.Exec(ctx, `
		INSERT INTO follows (follower_id, followee_id)
		VALUES ($1, $2)`, followerID, followeeID)
	if err != nil {
		return fmt.Errorf("failed to follow user: %w", err)
	}

	// Create notification
	if s.notificationsService != nil {
		err = s.notificationsService.NotifyFollow(ctx, followerID, followeeID)
		if err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Failed to create follow notification: %v\n", err)
		}
	}

	return nil
}

func (s *SocialService) UnfollowUser(ctx context.Context, followerID, followeeID uuid.UUID) error {
	result, err := s.db.Exec(ctx, `
		DELETE FROM follows
		WHERE follower_id = $1 AND followee_id = $2`,
		followerID, followeeID)
	if err != nil {
		return fmt.Errorf("failed to unfollow user: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("not following this user")
	}

	return nil
}

func (s *SocialService) GetFollowStats(ctx context.Context, userID, currentUserID uuid.UUID) (*FollowStats, error) {
	var stats FollowStats

	// Get followers count
	err := s.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM follows WHERE followee_id = $1`, userID).Scan(&stats.FollowersCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers count: %w", err)
	}

	// Get following count
	err = s.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM follows WHERE follower_id = $1`, userID).Scan(&stats.FollowingCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get following count: %w", err)
	}

	// Check if current user is following this user
	if currentUserID != uuid.Nil {
		var followCount int
		err = s.db.QueryRow(ctx, `
			SELECT COUNT(*) FROM follows
			WHERE follower_id = $1 AND followee_id = $2`,
			currentUserID, userID).Scan(&followCount)
		if err != nil {
			return nil, fmt.Errorf("failed to check follow status: %w", err)
		}
		stats.IsFollowing = followCount > 0
	}

	return &stats, nil
}

func (s *SocialService) GetFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*FeedPost, error) {
	rows, err := s.db.Query(ctx, `
		SELECT p.id, p.author_id, p.text, p.course_id, p.module_id, p.created_at, p.updated_at,
		       COUNT(DISTINCT l.user_id) as like_count,
		       COUNT(DISTINCT c.id) as comment_count,
		       u.username, u.email, u.bio, u.avatar_url,
		       CASE WHEN ul.user_id IS NOT NULL THEN true ELSE false END as is_liked
		FROM posts p
		JOIN users u ON p.author_id = u.id
		LEFT JOIN likes l ON p.id = l.post_id
		LEFT JOIN comments c ON p.id = c.post_id
		LEFT JOIN likes ul ON p.id = ul.post_id AND ul.user_id = $1
		WHERE p.author_id IN (
		    SELECT followee_id FROM follows WHERE follower_id = $1
		    UNION
		    SELECT $1
		)
		GROUP BY p.id, u.username, u.email, u.bio, u.avatar_url, ul.user_id
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get feed: %w", err)
	}
	defer rows.Close()

	var posts []*FeedPost
	for rows.Next() {
		var post FeedPost
		var courseID, moduleID pgtype.UUID
		var bio, avatarURL pgtype.Text

		err := rows.Scan(
			&post.ID, &post.AuthorID, &post.Text, &courseID, &moduleID,
			&post.CreatedAt, &post.UpdatedAt, &post.LikeCount, &post.CommentCount,
			&post.Author.Username, &post.Author.Email, &bio, &avatarURL, &post.IsLiked)
		if err != nil {
			return nil, fmt.Errorf("failed to scan feed post: %w", err)
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

		posts = append(posts, &post)
	}

	return posts, nil
}

func (s *SocialService) GetFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*UserResponse, error) {
	rows, err := s.db.Query(ctx, `
		SELECT u.id, u.username, u.email, u.bio, u.avatar_url
		FROM follows f
		JOIN users u ON f.follower_id = u.id
		WHERE f.followee_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}
	defer rows.Close()

	var users []*UserResponse
	for rows.Next() {
		var user UserResponse
		var bio, avatarURL pgtype.Text

		err := rows.Scan(&user.ID, &user.Username, &user.Email, &bio, &avatarURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan follower: %w", err)
		}

		user.Bio = getPgtypeTextValue(bio)
		user.AvatarURL = getPgtypeTextPtr(avatarURL)
		users = append(users, &user)
	}

	return users, nil
}

func (s *SocialService) GetFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*UserResponse, error) {
	rows, err := s.db.Query(ctx, `
		SELECT u.id, u.username, u.email, u.bio, u.avatar_url
		FROM follows f
		JOIN users u ON f.followee_id = u.id
		WHERE f.follower_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}
	defer rows.Close()

	var users []*UserResponse
	for rows.Next() {
		var user UserResponse
		var bio, avatarURL pgtype.Text

		err := rows.Scan(&user.ID, &user.Username, &user.Email, &bio, &avatarURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan following user: %w", err)
		}

		user.Bio = getPgtypeTextValue(bio)
		user.AvatarURL = getPgtypeTextPtr(avatarURL)
		users = append(users, &user)
	}

	return users, nil
}

func (s *SocialService) IsFollowing(ctx context.Context, followerID, followeeID uuid.UUID) (bool, error) {
	var count int
	err := s.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM follows
		WHERE follower_id = $1 AND followee_id = $2`,
		followerID, followeeID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check follow status: %w", err)
	}
	return count > 0, nil
}

func (s *SocialService) GetCourses(ctx context.Context) ([]*Course, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, title, description
		FROM courses
		ORDER BY title`)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}
	defer rows.Close()

	var courses []*Course
	for rows.Next() {
		var course Course
		var description pgtype.Text

		err := rows.Scan(&course.ID, &course.Title, &description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan course: %w", err)
		}

		course.Description = getPgtypeTextValue(description)
		courses = append(courses, &course)
	}

	return courses, nil
}

func (s *SocialService) GetModulesByCourse(ctx context.Context, courseID uuid.UUID) ([]*Module, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, course_id, title, "order"
		FROM modules
		WHERE course_id = $1
		ORDER BY "order"`, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get modules: %w", err)
	}
	defer rows.Close()

	var modules []*Module
	for rows.Next() {
		var module Module
		err := rows.Scan(&module.ID, &module.CourseID, &module.Title, &module.Order)
		if err != nil {
			return nil, fmt.Errorf("failed to scan module: %w", err)
		}
		modules = append(modules, &module)
	}

	return modules, nil
}

// Additional types
type Course struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

type Module struct {
	ID       uuid.UUID `json:"id"`
	CourseID uuid.UUID `json:"course_id"`
	Title    string    `json:"title"`
	Order    int       `json:"order"`
}
