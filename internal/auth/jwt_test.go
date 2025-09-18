package auth

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestSetJWTKey(t *testing.T) {
	t.Run("valid key", func(t *testing.T) {
		testKey := "test-secret-key"
		SetJWTKey(testKey)

		if len(jwtKey) != len(testKey) {
			t.Errorf("Expected key length %d, got %d", len(testKey), len(jwtKey))
		}

		if string(jwtKey) != testKey {
			t.Errorf("Expected key %s, got %s", testKey, string(jwtKey))
		}
	})

	t.Run("empty key should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for empty key, but didn't panic")
			}
		}()

		SetJWTKey("")
	})
}

func TestGenerateToken(t *testing.T) {
	// Устанавливаем тестовый ключ
	SetJWTKey("test-secret-key-for-token-generation")

	t.Run("valid username", func(t *testing.T) {
		username := "testuser"
		token, err := GenerateToken(username)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if token == "" {
			t.Errorf("Expected non-empty token")
		}

		// Проверяем, что токен содержит правильные части (header.payload.signature)
		parts := strings.Split(token, ".")
		if len(parts) != 3 {
			t.Errorf("Expected token to have 3 parts, got %d", len(parts))
		}
	})

	t.Run("empty username", func(t *testing.T) {
		token, err := GenerateToken("")

		if err != nil {
			t.Fatalf("Expected no error for empty username, got %v", err)
		}

		if token == "" {
			t.Errorf("Expected non-empty token even for empty username")
		}
	})

	t.Run("no key set", func(t *testing.T) {
		// Сохраняем текущий ключ
		originalKey := make([]byte, len(jwtKey))
		copy(originalKey, jwtKey)

		// Очищаем ключ
		jwtKey = nil

		token, err := GenerateToken("testuser")

		// Восстанавливаем ключ
		jwtKey = originalKey

		if err != ErrKeyNotSet {
			t.Errorf("Expected ErrKeyNotSet, got %v", err)
		}

		if token != "" {
			t.Errorf("Expected empty token when key not set, got %s", token)
		}
	})
}

func TestValidateToken(t *testing.T) {
	// Устанавливаем тестовый ключ
	testKey := "test-secret-key-for-token-validation"
	SetJWTKey(testKey)

	t.Run("valid token", func(t *testing.T) {
		username := "testuser"

		// Генерируем токен
		token, err := GenerateToken(username)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		// Валидируем токен
		claims, err := ValidateToken(token)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if claims.Username != username {
			t.Errorf("Expected username %s, got %s", username, claims.Username)
		}

		// Проверяем время истечения
		if claims.ExpiresAt.Time.Before(time.Now()) {
			t.Errorf("Token should not be expired")
		}

		// Проверяем время выдачи
		if claims.IssuedAt.Time.After(time.Now()) {
			t.Errorf("Token issued time should be in the past")
		}
	})

	t.Run("invalid token format", func(t *testing.T) {
		invalidToken := "invalid.token.format"

		claims, err := ValidateToken(invalidToken)

		if err == nil {
			t.Errorf("Expected error for invalid token format")
		}

		if claims != nil {
			t.Errorf("Expected nil claims for invalid token")
		}
	})

	t.Run("malformed token", func(t *testing.T) {
		malformedToken := "not.a.jwt"

		claims, err := ValidateToken(malformedToken)

		if err == nil {
			t.Errorf("Expected error for malformed token")
		}

		if claims != nil {
			t.Errorf("Expected nil claims for malformed token")
		}
	})

	t.Run("token with wrong signature", func(t *testing.T) {
		// Генерируем токен с одним ключом
		SetJWTKey("first-key")
		token, err := GenerateToken("testuser")
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		// Пытаемся валидировать с другим ключом
		SetJWTKey("second-key")
		claims, err := ValidateToken(token)

		if err == nil {
			t.Errorf("Expected error for token with wrong signature")
		}

		if claims != nil {
			t.Errorf("Expected nil claims for token with wrong signature")
		}

		// Восстанавливаем исходный ключ
		SetJWTKey(testKey)
	})

	t.Run("expired token", func(t *testing.T) {
		// Создаем токен с истекшим сроком действия
		claims := &Claims{
			Username: "testuser",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Истек час назад
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)), // Выдан 2 часа назад
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			t.Fatalf("Failed to create expired token: %v", err)
		}

		validatedClaims, err := ValidateToken(tokenString)

		if err == nil {
			t.Errorf("Expected error for expired token")
		}

		if validatedClaims != nil {
			t.Errorf("Expected nil claims for expired token")
		}
	})

	t.Run("no key set", func(t *testing.T) {
		// Сохраняем текущий ключ
		originalKey := make([]byte, len(jwtKey))
		copy(originalKey, jwtKey)

		// Очищаем ключ
		jwtKey = nil

		claims, err := ValidateToken("some.token.here")

		// Восстанавливаем ключ
		jwtKey = originalKey

		if err != ErrKeyNotSet {
			t.Errorf("Expected ErrKeyNotSet, got %v", err)
		}

		if claims != nil {
			t.Errorf("Expected nil claims when key not set")
		}
	})

	t.Run("token with wrong signing method", func(t *testing.T) {
		// Создаем токен вручную для тестирования проверки метода подписи
		// (используем токен с RS256 алгоритмом вместо HS256)
		tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNjk5OTk5OTk5LCJpYXQiOjE2OTk5OTk5OTl9.invalid_signature"

		validatedClaims, err := ValidateToken(tokenString)

		if err == nil {
			t.Errorf("Expected error for token with wrong signing method")
		}

		if validatedClaims != nil {
			t.Errorf("Expected nil claims for token with wrong signing method")
		}
	})
}

func TestSetUserContext(t *testing.T) {
	t.Run("set user in context", func(t *testing.T) {
		ctx := context.Background()
		username := "testuser"

		newCtx := SetUserContext(ctx, username)

		if newCtx == ctx {
			t.Errorf("Expected new context, got same context")
		}

		// Проверяем, что значение установлено
		if value := newCtx.Value(userContextKey); value != username {
			t.Errorf("Expected username %s in context, got %v", username, value)
		}
	})

	t.Run("set empty username", func(t *testing.T) {
		ctx := context.Background()
		username := ""

		newCtx := SetUserContext(ctx, username)

		if value := newCtx.Value(userContextKey); value != username {
			t.Errorf("Expected empty username in context, got %v", value)
		}
	})
}

func TestGetUserFromContext(t *testing.T) {
	t.Run("get existing user from context", func(t *testing.T) {
		ctx := context.Background()
		username := "testuser"

		// Устанавливаем пользователя в контекст
		newCtx := SetUserContext(ctx, username)

		// Получаем пользователя из контекста
		retrievedUsername := GetUserFromContext(newCtx)

		if retrievedUsername != username {
			t.Errorf("Expected username %s, got %s", username, retrievedUsername)
		}
	})

	t.Run("get user from empty context", func(t *testing.T) {
		ctx := context.Background()

		username := GetUserFromContext(ctx)

		if username != "" {
			t.Errorf("Expected empty string for empty context, got %s", username)
		}
	})

	t.Run("get user from context with wrong type", func(t *testing.T) {
		ctx := context.Background()

		// Устанавливаем значение неправильного типа
		ctx = context.WithValue(ctx, userContextKey, 123)

		username := GetUserFromContext(ctx)

		if username != "" {
			t.Errorf("Expected empty string for wrong type in context, got %s", username)
		}
	})
}

func TestClaims(t *testing.T) {
	t.Run("claims structure", func(t *testing.T) {
		username := "testuser"
		now := time.Now()
		expiry := now.Add(24 * time.Hour)

		claims := &Claims{
			Username: username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiry),
				IssuedAt:  jwt.NewNumericDate(now),
			},
		}

		if claims.Username != username {
			t.Errorf("Expected username %s, got %s", username, claims.Username)
		}

		if claims.ExpiresAt.Time.Unix() != expiry.Unix() {
			t.Errorf("Expected expiry time %v, got %v", expiry.Unix(), claims.ExpiresAt.Time.Unix())
		}

		if claims.IssuedAt.Time.Unix() != now.Unix() {
			t.Errorf("Expected issued time %v, got %v", now.Unix(), claims.IssuedAt.Time.Unix())
		}
	})
}

func TestTokenLifecycle(t *testing.T) {
	t.Run("complete token lifecycle", func(t *testing.T) {
		// Устанавливаем ключ
		testKey := "test-lifecycle-key"
		SetJWTKey(testKey)

		username := "lifecycleuser"

		// 1. Генерируем токен
		token, err := GenerateToken(username)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		// 2. Валидируем токен
		claims, err := ValidateToken(token)
		if err != nil {
			t.Fatalf("Failed to validate token: %v", err)
		}

		if claims.Username != username {
			t.Errorf("Expected username %s, got %s", username, claims.Username)
		}

		// 3. Используем в контексте
		ctx := context.Background()
		ctx = SetUserContext(ctx, claims.Username)

		// 4. Получаем из контекста
		retrievedUsername := GetUserFromContext(ctx)

		if retrievedUsername != username {
			t.Errorf("Expected username %s from context, got %s", username, retrievedUsername)
		}
	})
}

// Benchmark тесты для измерения производительности
func BenchmarkGenerateToken(b *testing.B) {
	SetJWTKey("benchmark-secret-key")
	username := "benchmarkuser"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GenerateToken(username)
		if err != nil {
			b.Fatalf("Failed to generate token: %v", err)
		}
	}
}

func BenchmarkValidateToken(b *testing.B) {
	SetJWTKey("benchmark-secret-key")
	username := "benchmarkuser"

	// Генерируем токен для бенчмарка
	token, err := GenerateToken(username)
	if err != nil {
		b.Fatalf("Failed to generate token for benchmark: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ValidateToken(token)
		if err != nil {
			b.Fatalf("Failed to validate token: %v", err)
		}
	}
}

func BenchmarkSetUserContext(b *testing.B) {
	ctx := context.Background()
	username := "benchmarkuser"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SetUserContext(ctx, username)
	}
}

func BenchmarkGetUserFromContext(b *testing.B) {
	ctx := context.Background()
	username := "benchmarkuser"
	ctx = SetUserContext(ctx, username)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetUserFromContext(ctx)
	}
}
