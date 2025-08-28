package http

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"golang.org/x/time/rate"

	"bailanysta/api/internal/config"
	"bailanysta/api/internal/http/handlers"
	"bailanysta/api/internal/pkg/auth"
	"bailanysta/api/internal/pkg/logger"
)

type Router struct {
	*chi.Mux
}

type Deps struct {
	Config     *config.Config
	Logger     *logger.Logger
	Handlers   *Handlers
	JWTManager *auth.JWTManager
}

type Handlers struct {
	Auth          *handlers.AuthHandler
	Users         *handlers.UsersHandler
	Posts         *handlers.PostsHandler
	Social        *handlers.SocialHandler
	Search        *handlers.SearchHandler
	Notifications *handlers.NotificationsHandler
	AI            *handlers.AIHandler
	Health        *handlers.HealthHandler
}

func NewRouter(deps *Deps) *Router {
	r := chi.NewRouter()

	// Basic middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(loggerMiddleware(deps.Logger))

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{deps.Config.CORSOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Rate limiting middleware
	r.Use(rateLimitMiddleware(deps.Config.RateLimitRPM))

	// Health endpoint (no auth required)
	r.Get("/health", deps.Handlers.Health.HealthCheck)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Auth routes (no auth required)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", deps.Handlers.Auth.Register)
			r.Post("/login", deps.Handlers.Auth.Login)
			r.Post("/refresh", deps.Handlers.Auth.Refresh)
			r.Post("/logout", deps.Handlers.Auth.Logout)
		})

		// Public routes (no auth required)
		r.Get("/courses", deps.Handlers.Social.GetCourses)
		r.Get("/courses/{id}/modules", deps.Handlers.Social.GetModulesByCourse)
		r.Get("/search", deps.Handlers.Search.SearchPosts)

		// Protected routes
		r.Route("/", func(r chi.Router) {
			r.Use(AuthMiddleware(deps.JWTManager, deps.Logger))

			// User routes
			r.Get("/me", deps.Handlers.Users.GetCurrentUser)
			r.Patch("/me", deps.Handlers.Users.UpdateCurrentUser)
			r.Get("/users", deps.Handlers.Users.GetAllUsers)
			r.Get("/users/{id}", deps.Handlers.Users.GetUserByID)
			r.Post("/users/{id}/follow", deps.Handlers.Social.FollowUser)
			r.Delete("/users/{id}/follow", deps.Handlers.Social.UnfollowUser)

			// Posts routes
			r.Post("/posts", deps.Handlers.Posts.CreatePost)
			r.Get("/posts/{id}", deps.Handlers.Posts.GetPostByID)
			r.Patch("/posts/{id}", deps.Handlers.Posts.UpdatePost)
			r.Delete("/posts/{id}", deps.Handlers.Posts.DeletePost)
			r.Post("/posts/{id}/like", deps.Handlers.Posts.LikePost)
			r.Delete("/posts/{id}/like", deps.Handlers.Posts.UnlikePost)
			r.Get("/posts/{id}/comments", deps.Handlers.Posts.GetComments)
			r.Post("/posts/{id}/comments", deps.Handlers.Posts.CreateComment)

			// Feed
			r.Get("/feed", deps.Handlers.Social.GetFeed)

			// Notifications
			r.Get("/notifications", deps.Handlers.Notifications.GetNotifications)
			r.Post("/notifications/mark-read", deps.Handlers.Notifications.MarkAllAsRead)
			r.Get("/notifications/unread-count", deps.Handlers.Notifications.GetUnreadCount)
			r.Post("/notifications/{id}/mark-read", deps.Handlers.Notifications.MarkAsRead)
			r.Delete("/notifications/{id}", deps.Handlers.Notifications.DeleteNotification)

			// AI
			r.Post("/ai/generate", deps.Handlers.AI.GenerateText)
			r.Post("/ai/generate-post", deps.Handlers.AI.GeneratePost)
			r.Post("/ai/generate-comment", deps.Handlers.AI.GenerateComment)
			r.Post("/ai/generate-study-notes", deps.Handlers.AI.GenerateStudyNotes)
			r.Post("/ai/generate-quiz", deps.Handlers.AI.GenerateQuiz)
			r.Post("/ai/explain-concept", deps.Handlers.AI.ExplainConcept)
		})
	})

	return &Router{Mux: r}
}

func rateLimitMiddleware(rpm int) func(http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(rpm)/60, rpm/4) // burst size = rpm/4

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func loggerMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a custom ResponseWriter to capture status code
			rw := &responseWriter{ResponseWriter: w, statusCode: 200}

			next.ServeHTTP(rw, r)

			log.Info("HTTP request", map[string]interface{}{
				"method":      r.Method,
				"path":        r.URL.Path,
				"status":      rw.statusCode,
				"duration_ms": time.Since(start).Milliseconds(),
				"user_agent":  r.Header.Get("User-Agent"),
			})
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func AuthMiddleware(jwtManager *auth.JWTManager, logger *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn("Missing authorization header", map[string]interface{}{
					"path": r.URL.Path,
				})
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
				logger.Warn("Invalid authorization header format", map[string]interface{}{
					"path": r.URL.Path,
				})
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := authHeader[7:]

			// Validate JWT token
			claims, err := jwtManager.ValidateAccessToken(tokenString)
			if err != nil {
				logger.Warn("Invalid JWT token", map[string]interface{}{
					"path":  r.URL.Path,
					"error": err.Error(),
				})
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Add user ID to context
			ctx := context.WithValue(r.Context(), "user_id", claims.UserID.String())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
