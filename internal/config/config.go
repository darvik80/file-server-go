package config

import (
	"crypto/rand"
	"encoding/hex"
	"os"
)

// Config contains application configuration
type Config struct {
	Port        string
	UploadDir   string
	JWTKey      string
	Auth        AuthConfig
	MaxFileSize int64 // in bytes
}

// AuthConfig contains authentication settings
type AuthConfig struct {
	Username string
	Password string
}

// generateSecretKey generates a random key
func generateSecretKey() string {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		panic("Failed to generate secret key: " + err.Error())
	}
	return hex.EncodeToString(bytes)
}

// Load loads configuration from environment variables
func Load() *Config {
	// Generate random key if not specified in environment variables
	jwtKey := getEnv("JWT_KEY", "")
	if jwtKey == "" {
		jwtKey = generateSecretKey()
	}

	cfg := &Config{
		Port:      getEnv("PORT", "8080"),
		UploadDir: getEnv("UPLOAD_DIR", "./uploads"),
		JWTKey:    jwtKey,
		Auth: AuthConfig{
			Username: getEnv("ADMIN_USERNAME", "admin"),
			Password: getEnv("ADMIN_PASSWORD", "admin"),
		},
		MaxFileSize: 10 << 20, // 10MB by default
	}

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		panic("Failed to create upload directory: " + err.Error())
	}

	return cfg
}

// getEnv returns environment variable value or default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
