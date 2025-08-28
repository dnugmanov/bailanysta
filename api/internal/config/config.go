package config

import (
	"fmt"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port           string        `envconfig:"PORT" default:"8080"`
	DatabaseURL    string        `envconfig:"DATABASE_URL" required:"true"`
	JwtSecret      string        `envconfig:"JWT_SECRET" required:"true"`
	JwtExpiry      time.Duration `envconfig:"JWT_EXPIRY" default:"15m"`
	RefreshExpiry  time.Duration `envconfig:"REFRESH_EXPIRY" default:"168h"`
	CORSOrigin     string        `envconfig:"CORS_ORIGIN" default:"http://localhost:3000"`
	MigrateOnStart bool          `envconfig:"MIGRATE_ON_START" default:"false"`
	LogLevel       string        `envconfig:"LOG_LEVEL" default:"info"`

	// AI Configuration
	OpenAIBaseURL string `envconfig:"OPENAI_BASE_URL" default:"https://api.openai.com/v1"`
	OpenAIApiKey  string `envconfig:"OPENAI_API_KEY"`

	// Rate limiting
	RateLimitRPM int `envconfig:"RATE_LIMIT_RPM" default:"100"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.JwtSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.Port == "" {
		return fmt.Errorf("PORT is required")
	}
	return nil
}

func (c *Config) PrintConfig() {
	log.Printf("Configuration loaded:")
	log.Printf("  Port: %s", c.Port)
	log.Printf("  Database URL: %s", maskPassword(c.DatabaseURL))
	log.Printf("  JWT Secret: %s", maskSecret(c.JwtSecret))
	log.Printf("  JWT Expiry: %v", c.JwtExpiry)
	log.Printf("  Refresh Expiry: %v", c.RefreshExpiry)
	log.Printf("  CORS Origin: %s", c.CORSOrigin)
	log.Printf("  Migrate on Start: %v", c.MigrateOnStart)
	log.Printf("  Log Level: %s", c.LogLevel)
	log.Printf("  OpenAI Base URL: %s", c.OpenAIBaseURL)
	log.Printf("  OpenAI API Key: %s", maskSecret(c.OpenAIApiKey))
	log.Printf("  Rate Limit RPM: %d", c.RateLimitRPM)
}

func maskPassword(url string) string {
	return "***"
}

func maskSecret(secret string) string {
	if len(secret) <= 4 {
		return "***"
	}
	return secret[:4] + "***"
}
