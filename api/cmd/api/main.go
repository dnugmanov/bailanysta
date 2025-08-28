package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	"bailanysta/api/internal/config"
	httpRouter "bailanysta/api/internal/http"
	"bailanysta/api/internal/http/handlers"
	"bailanysta/api/internal/pkg/ai"
	"bailanysta/api/internal/pkg/auth"
	"bailanysta/api/internal/pkg/logger"
	"bailanysta/api/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	cfg.PrintConfig()

	// Initialize logger
	appLogger := logger.New(cfg.LogLevel, os.Stdout)

	// Connect to database
	dbpool, err := connectDB(cfg.DatabaseURL)
	if err != nil {
		appLogger.Fatal("Failed to connect to database", map[string]interface{}{
			"error": err.Error(),
		})
	}
	defer dbpool.Close()

	appLogger.Info("Connected to database")

	// Run migrations if enabled
	if cfg.MigrateOnStart {
		if err := runMigrations(cfg.DatabaseURL); err != nil {
			appLogger.Fatal("Failed to run migrations", map[string]interface{}{
				"error": err.Error(),
			})
		}
		appLogger.Info("Migrations completed")
	}

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(cfg.JwtSecret, cfg.JwtExpiry, cfg.RefreshExpiry)

	// Initialize AI client
	aiClient := ai.NewClient(cfg.OpenAIBaseURL, cfg.OpenAIApiKey)

	// Initialize services
	notificationsService := services.NewNotificationService(dbpool)
	authService := services.NewAuthService(dbpool, jwtManager)
	postsService := services.NewPostsService(dbpool, notificationsService)
	socialService := services.NewSocialService(dbpool, notificationsService)
	aiService := services.NewAIService(aiClient)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, appLogger)
	postsHandler := handlers.NewPostsHandler(postsService, appLogger, jwtManager)
	socialHandler := handlers.NewSocialHandler(socialService, appLogger, jwtManager)
	usersHandler := handlers.NewUsersHandler(authService, socialService, appLogger, jwtManager)
	searchHandler := handlers.NewSearchHandler(dbpool, appLogger, jwtManager)
	notificationsHandler := handlers.NewNotificationsHandler(notificationsService, appLogger, jwtManager)
	aiHandler := handlers.NewAIHandler(aiService, appLogger)

	handlers := &httpRouter.Handlers{
		Auth:          authHandler,
		Posts:         postsHandler,
		Social:        socialHandler,
		Users:         usersHandler,
		Search:        searchHandler,
		Notifications: notificationsHandler,
		AI:            aiHandler,
		Health:        &handlers.HealthHandler{Logger: appLogger},
	}

	// Create router
	router := httpRouter.NewRouter(&httpRouter.Deps{
		Config:     cfg,
		Logger:     appLogger,
		Handlers:   handlers,
		JWTManager: jwtManager,
	})

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  120 * time.Second, // Increased for AI requests
		WriteTimeout: 120 * time.Second, // Increased for AI requests
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for interrupt signal
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	// Register interrupt signals
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		appLogger.Info("Starting server", map[string]interface{}{
			"port": cfg.Port,
		})
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Server failed to start", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	// Wait for interrupt signal
	<-quit
	appLogger.Info("Server is shutting down...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", map[string]interface{}{
			"error": err.Error(),
		})
	}

	close(done)
	appLogger.Info("Server exited")
}

func connectDB(databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

func runMigrations(databaseURL string) error {
	m, err := migrate.New(
		"file://api/internal/db/migrations",
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
