package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"bailanysta/api/internal/pkg/auth"
)

type AuthService struct {
	db         *pgxpool.Pool
	jwtManager *auth.JWTManager
}

type User struct {
	ID        uuid.UUID      `json:"id"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	Bio       sql.NullString `json:"bio"`
	AvatarURL sql.NullString `json:"avatar_url"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User   UserResponse   `json:"user"`
	Tokens auth.TokenPair `json:"tokens"`
}

type UserResponse struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Bio            string    `json:"bio"`
	AvatarURL      *string   `json:"avatar_url,omitempty"`
	FollowersCount int       `json:"followers_count,omitempty"`
	FollowingCount int       `json:"following_count,omitempty"`
	IsFollowing    bool      `json:"is_following,omitempty"`
}

func NewAuthService(db *pgxpool.Pool, jwtManager *auth.JWTManager) *AuthService {
	return &AuthService{
		db:         db,
		jwtManager: jwtManager,
	}
}

func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	var existingUser User
	err := s.db.QueryRow(ctx, "SELECT id, username, email FROM users WHERE email = $1", req.Email).Scan(&existingUser.ID, &existingUser.Username, &existingUser.Email)
	if err == nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	err = s.db.QueryRow(ctx, "SELECT id, username, email FROM users WHERE username = $1", req.Username).Scan(&existingUser.ID, &existingUser.Username, &existingUser.Email)
	if err == nil {
		return nil, fmt.Errorf("user with this username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	var user User
	err = s.db.QueryRow(ctx, `
		INSERT INTO users (username, email, password_hash, bio, avatar_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, username, email, bio, avatar_url`,
		req.Username, req.Email, string(hashedPassword), nil, nil).Scan(
		&user.ID, &user.Username, &user.Email, &user.Bio, &user.AvatarURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	tokens, err := s.jwtManager.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &AuthResponse{
		User: UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Bio:       getNullStringValue(user.Bio),
			AvatarURL: getNullStringPtr(user.AvatarURL),
		},
		Tokens: *tokens,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	// Get user by email
	var user User
	var passwordHash string
	err := s.db.QueryRow(ctx, `
		SELECT id, username, email, password_hash, bio, avatar_url
		FROM users WHERE email = $1`, req.Email).Scan(
		&user.ID, &user.Username, &user.Email, &passwordHash, &user.Bio, &user.AvatarURL)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate tokens
	tokens, err := s.jwtManager.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &AuthResponse{
		User: UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Bio:       getNullStringValue(user.Bio),
			AvatarURL: getNullStringPtr(user.AvatarURL),
		},
		Tokens: *tokens,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*auth.TokenPair, error) {
	// Validate refresh token
	if err := s.jwtManager.ValidateRefreshToken(refreshToken); err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// For now, we need to extract user ID from refresh token
	// In production, you'd store refresh tokens in database with user association
	return nil, fmt.Errorf("refresh token functionality not implemented yet")
}

func (s *AuthService) ValidateToken(tokenString string) (*auth.Claims, error) {
	return s.jwtManager.ValidateAccessToken(tokenString)
}

func (s *AuthService) GetDB() *pgxpool.Pool {
	return s.db
}

func (s *AuthService) GetCurrentUser(ctx context.Context, userID uuid.UUID) (*UserResponse, error) {
	var user User
	err := s.db.QueryRow(ctx, `
		SELECT id, username, email, bio, avatar_url
		FROM users WHERE id = $1`, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.Bio, &user.AvatarURL)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Bio:       getNullStringValue(user.Bio),
		AvatarURL: getNullStringPtr(user.AvatarURL),
	}, nil
}

// Helper functions
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// Helper functions for sql.NullString
func getNullStringValue(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func getNullStringPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}
