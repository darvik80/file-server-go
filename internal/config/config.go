package config

import (
	"os"
)

// Config содержит конфигурацию приложения
type Config struct {
	Port      string
	UploadDir string
	MaxFileSize int64 // в байтах
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		UploadDir:   getEnv("UPLOAD_DIR", "./uploads"),
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