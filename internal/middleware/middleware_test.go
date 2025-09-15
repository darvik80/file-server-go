package middleware

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"file-server-go/internal/auth"
)

func TestCORS(t *testing.T) {
	// Создаем тестовый обработчик
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Оборачиваем в CORS middleware
	corsHandler := CORS(handler)

	t.Run("OPTIONS request", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/", nil)
		w := httptest.NewRecorder()

		corsHandler.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 for OPTIONS, got %d", resp.StatusCode)
		}

		// Проверяем заголовки CORS
		checkCORSHeaders(t, resp)
	})

	t.Run("Regular request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		corsHandler.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 for GET, got %d", resp.StatusCode)
		}

		// Проверяем заголовки CORS
		checkCORSHeaders(t, resp)
	})
}

func checkCORSHeaders(t *testing.T, resp *http.Response) {
	headers := map[string]string{
		"Access-Control-Allow-Origin":   "*",
		"Access-Control-Allow-Methods":  "GET, POST, PUT, DELETE, OPTIONS",
		"Access-Control-Allow-Headers":  "Content-Type, Authorization",
		"Access-Control-Expose-Headers": "Authorization",
	}

	for header, expected := range headers {
		if actual := resp.Header.Get(header); actual != expected {
			t.Errorf("Expected %s: %s, got %s", header, expected, actual)
		}
	}
}

func TestAuth(t *testing.T) {
	// Устанавливаем JWT ключ для тестов
	auth.SetJWTKey("test-secret-key-for-testing")

	// Создаем тестовый обработчик
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Оборачиваем в Auth middleware
	authHandler := Auth(handler)

	t.Run("Valid token", func(t *testing.T) {
		// Генерируем валидный токен
		token, err := auth.GenerateToken("testuser")
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		authHandler.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 with valid token, got %d", resp.StatusCode)
		}
	})

	t.Run("Missing token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		w := httptest.NewRecorder()

		authHandler.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401 with missing token, got %d", resp.StatusCode)
		}

		checkErrorResponse(t, resp, "Отсутствует токен авторизации")
	})

	t.Run("Invalid token format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "InvalidToken")
		w := httptest.NewRecorder()

		authHandler.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401 with invalid token format, got %d", resp.StatusCode)
		}

		checkErrorResponse(t, resp, "Неверный формат токена")
	})

	t.Run("Invalid token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		w := httptest.NewRecorder()

		authHandler.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401 with invalid token, got %d", resp.StatusCode)
		}

		checkErrorResponse(t, resp, "Недействительный токен")
	})

	t.Run("Exempt paths", func(t *testing.T) {
		exemptPaths := []string{"/login", "/", "/static/test.css"}

		for _, path := range exemptPaths {
			t.Run(path, func(t *testing.T) {
				req := httptest.NewRequest("GET", path, nil)
				w := httptest.NewRecorder()

				authHandler.ServeHTTP(w, req)

				resp := w.Result()
				if resp.StatusCode != http.StatusOK {
					t.Errorf("Expected status 200 for exempt path %s, got %d", path, resp.StatusCode)
				}
			})
		}
	})
}

func checkErrorResponse(t *testing.T, resp *http.Response, expectedError string) {
	var errorResp map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	if errorResp["error"] != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, errorResp["error"])
	}
}

func TestLogging(t *testing.T) {
	// Создаем тестовый обработчик
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Оборачиваем в Logging middleware
	loggingHandler := Logging(handler)

	// Захватываем вывод логов
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("User-Agent", "test-agent")
	w := httptest.NewRecorder()

	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		log.Printf("Request: %s %s, Status: %d, User-Agent: %s, Duration: %s", req.Method, req.URL.Path, w.Code, req.UserAgent(), duration)
	}()

	loggingHandler.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Проверяем что лог содержит ожидаемую информацию
	logOutput := logBuffer.String()
	if !strings.Contains(logOutput, "GET") {
		t.Errorf("Expected log to contain 'GET', got: %s", logOutput)
	}
	if !strings.Contains(logOutput, "/test") {
		t.Errorf("Expected log to contain '/test', got: %s", logOutput)
	}
	if !strings.Contains(logOutput, "200") {
		t.Errorf("Expected log to contain '200', got: %s", logOutput)
	}
	if !strings.Contains(logOutput, "test-agent") {
		t.Errorf("Expected log to contain 'test-agent', got: %s", logOutput)
	}
	// Проверяем что в логе есть информация о времени выполнения (любая единица времени)
	hasTimeUnit := strings.Contains(logOutput, "ns") ||
		strings.Contains(logOutput, "µs") ||
		strings.Contains(logOutput, "ms") ||
		strings.Contains(logOutput, "s")
	if !hasTimeUnit {
		t.Errorf("Expected log to contain time duration, got: %s", logOutput)
	}
}
