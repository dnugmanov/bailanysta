package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"bailanysta/api/internal/pkg/auth"
	"bailanysta/api/internal/pkg/logger"
	"bailanysta/api/internal/services"
)

type SearchHandler struct {
	db         *pgxpool.Pool
	logger     *logger.Logger
	jwtManager *auth.JWTManager
}

type SearchResult struct {
	Posts      []*services.Post         `json:"posts"`
	Users      []*services.UserResponse `json:"users"`
	Query      string                   `json:"query"`
	TotalPosts int                      `json:"total_posts"`
	TotalUsers int                      `json:"total_users"`
}

func NewSearchHandler(db *pgxpool.Pool, logger *logger.Logger, jwtManager *auth.JWTManager) *SearchHandler {
	return &SearchHandler{
		db:         db,
		logger:     logger,
		jwtManager: jwtManager,
	}
}

func (h *SearchHandler) SearchPosts(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("query"))
	if query == "" {
		h.respondWithError(w, "Query parameter is required", http.StatusBadRequest)
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

	// Get current user if authenticated
	currentUserID := uuid.Nil
	if userID, err := h.getUserIDFromContext(r.Context()); err == nil {
		currentUserID = userID
	}

	result := &SearchResult{
		Query:      query,
		Posts:      []*services.Post{},
		Users:      []*services.UserResponse{},
		TotalPosts: 0,
		TotalUsers: 0,
	}

	// Search posts - always use text search for better results
	posts, total, err := h.searchPostsByText(r.Context(), query, currentUserID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to search posts by text", map[string]interface{}{
			"error": err.Error(),
			"query": query,
		})
		h.respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if posts != nil {
		result.Posts = posts
	}
	result.TotalPosts = total

	// Search users
	users, userTotal, err := h.searchUsers(r.Context(), query, currentUserID, 10, 0)
	if err != nil {
		h.logger.Error("Failed to search users", map[string]interface{}{
			"error": err.Error(),
			"query": query,
		})
		h.respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if users != nil {
		result.Users = users
	}
	result.TotalUsers = userTotal

	h.logger.Info("Search completed", map[string]interface{}{
		"query":       query,
		"posts_found": len(result.Posts),
		"users_found": len(result.Users),
	})

	h.respondWithJSON(w, result, http.StatusOK)
}

func (h *SearchHandler) searchPostsByText(ctx context.Context, query string, currentUserID uuid.UUID, limit, offset int) ([]*services.Post, int, error) {
	var total int
	err := h.db.QueryRow(ctx, "SELECT COUNT(*) FROM posts WHERE text ILIKE '%' || $1 || '%'", query).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := h.db.Query(ctx, `
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
		WHERE p.text ILIKE '%' || $2 || '%'
		GROUP BY p.id, u.username, u.email, u.bio, u.avatar_url, ul.user_id
		ORDER BY p.created_at DESC
		LIMIT $3 OFFSET $4`, currentUserID, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []*services.Post
	for rows.Next() {
		var post services.Post
		var courseID, moduleID pgtype.UUID
		var bio, avatarURL pgtype.Text

		err := rows.Scan(
			&post.ID, &post.AuthorID, &post.Text, &courseID, &moduleID,
			&post.CreatedAt, &post.UpdatedAt, &post.LikeCount, &post.CommentCount,
			&post.Author.Username, &post.Author.Email, &bio, &avatarURL, &post.IsLiked)
		if err != nil {
			return nil, 0, err
		}

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

	return posts, total, nil
}

func (h *SearchHandler) searchPostsByHashtag(ctx context.Context, hashtag string, currentUserID uuid.UUID, limit, offset int) ([]*services.Post, int, error) {
	var total int
	err := h.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM posts p
		JOIN post_hashtags ph ON p.id = ph.post_id
		JOIN hashtags h ON ph.hashtag_id = h.id
		WHERE h.tag = $1`, hashtag).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := h.db.Query(ctx, `
		SELECT p.id, p.author_id, p.text, p.course_id, p.module_id, p.created_at, p.updated_at,
		       COUNT(DISTINCT l.user_id) as like_count,
		       COUNT(DISTINCT c.id) as comment_count,
		       u.username, u.email, u.bio, u.avatar_url,
		       CASE WHEN ul.user_id IS NOT NULL THEN true ELSE false END as is_liked
		FROM posts p
		JOIN users u ON p.author_id = u.id
		JOIN post_hashtags ph ON p.id = ph.post_id
		JOIN hashtags h ON ph.hashtag_id = h.id
		LEFT JOIN likes l ON p.id = l.post_id
		LEFT JOIN comments c ON p.id = c.post_id
		LEFT JOIN likes ul ON p.id = ul.post_id AND ul.user_id = $1
		WHERE h.tag = $2
		GROUP BY p.id, u.username, u.email, u.bio, u.avatar_url, ul.user_id
		ORDER BY p.created_at DESC
		LIMIT $3 OFFSET $4`, currentUserID, hashtag, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []*services.Post
	for rows.Next() {
		var post services.Post
		var courseID, moduleID pgtype.UUID
		var bio, avatarURL pgtype.Text

		err := rows.Scan(
			&post.ID, &post.AuthorID, &post.Text, &courseID, &moduleID,
			&post.CreatedAt, &post.UpdatedAt, &post.LikeCount, &post.CommentCount,
			&post.Author.Username, &post.Author.Email, &bio, &avatarURL, &post.IsLiked)
		if err != nil {
			return nil, 0, err
		}

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

	return posts, total, nil
}

func (h *SearchHandler) searchUsers(ctx context.Context, query string, currentUserID uuid.UUID, limit, offset int) ([]*services.UserResponse, int, error) {
	var total int
	err := h.db.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE username ILIKE '%' || $1 || '%' OR bio ILIKE '%' || $1 || '%'", query).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := h.db.Query(ctx, `
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
		WHERE u.username ILIKE '%' || $2 || '%' OR u.bio ILIKE '%' || $2 || '%'
		ORDER BY u.username
		LIMIT $3 OFFSET $4`, currentUserID, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*services.UserResponse
	for rows.Next() {
		var user services.UserResponse
		var bio, avatarURL pgtype.Text
		var followersCount, followingCount int

		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &bio, &avatarURL,
			&followersCount, &followingCount, &user.IsFollowing)
		if err != nil {
			return nil, 0, err
		}

		user.Bio = getPgtypeTextValue(bio)
		user.AvatarURL = getPgtypeTextPtr(avatarURL)
		user.FollowersCount = followersCount
		user.FollowingCount = followingCount

		users = append(users, &user)
	}

	return users, total, nil
}

func (h *SearchHandler) respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *SearchHandler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	h.respondWithJSON(w, map[string]interface{}{
		"error": map[string]interface{}{
			"code":    getErrorCode(statusCode),
			"message": message,
		},
	}, statusCode)
}

func (h *SearchHandler) getUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	return h.jwtManager.GetUserIDFromContext(ctx)
}

// Helper functions (duplicated from posts service, should be moved to shared package)
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
