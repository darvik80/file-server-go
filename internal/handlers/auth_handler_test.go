package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"file-server-go/internal/config"
)

func TestAuthHandler(t *testing.T) {
	// Создаем конфигурацию для тестов
	cfg := &config.Config{
		Port:      "8080",
		UploadDir: "./uploads",
		JWTKey:    "testsecret",
		Auth: config.AuthConfig{
			Username: "testuser",
			Password: "testpass",
		},
		MaxFileSize: 10 << 20, // 10MB
	}

	handler := NewAuthHandler(cfg)

	t.Run("ValidLogin", func(t *testing.T) {
		testValidLogin(t, handler)
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		testInvalidCredentials(t, handler)
	})

	t.Run("InvalidMethod", func(t *testing.T) {
		testInvalidMethod(t, handler)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		testInvalidJSON(t, handler)
	})
}

func testValidLogin(t *testing.T, handler *AuthHandler) {
	loginData := LoginRequest{
		Username: "testuser",
		Password: "testpass",
	}
	body, _ := json.Marshal(loginData)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var response LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Token == "" {
		t.Errorf("Expected token in response, got empty string")
	}
}

func testInvalidCredentials(t *testing.T, handler *AuthHandler) {
	tests := []struct {
		name     string
		username string
		password string
	}{
		{"wrong username", "wronguser", "testpass"},
		{"wrong password", "testuser", "wrongpass"},
		{"both wrong", "wronguser", "wrongpass"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loginData := LoginRequest{
				Username: tt.username,
				Password: tt.password,
			}
			body, _ := json.Marshal(loginData)

			req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Login(w, req)

			resp := w.Result()
			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("Expected status 401, got %d", resp.StatusCode)
			}
		})
	}
}

func testInvalidMethod(t *testing.T, handler *AuthHandler) {
	req := httptest.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()

	handler.Login(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", resp.StatusCode)
	}
}

func testInvalidJSON(t *testing.T, handler *AuthHandler) {
	req := httptest.NewRequest("POST", "/login", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}
