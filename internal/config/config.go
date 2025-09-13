package config

import (
	"crypto/rand"
	"encoding/hex"
	"os"
)

// Config содержит конфигурацию приложения
type Config struct {
	Port        string
	UploadDir   string
	JWTKey      string
	Auth        AuthConfig
	MaxFileSize int64 // в байтах
}

// AuthConfig содержит настройки авторизации
type AuthConfig struct {
	Username string
	Password string
}

// generateSecretKey генерирует случайный ключ
func generateSecretKey() string {
	bytes := make([]byte, 32) // 256 бит
	if _, err := rand.Read(bytes); err != nil {
		panic("Не удалось сгенерировать секретный ключ: " + err.Error())
	}
	return hex.EncodeToString(bytes)
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	// Генерируем случайный ключ, если не указан в переменных окружения
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
		MaxFileSize: 10 << 20, // 10MB по умолчанию
	}

	// Создаем папку для загрузок если её нет
	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		panic("Не удалось создать папку для загрузок: " + err.Error())
	}

	return cfg
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
