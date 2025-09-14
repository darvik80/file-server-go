package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"file-server-go/internal/auth"
)

// CORS добавляет заголовки CORS
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Auth проверяет JWT токен
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Пропускаем авторизацию только для логина и корневой страницы
		if r.URL.Path == "/login" || r.URL.Path == "/" || strings.HasPrefix(r.URL.Path, "/js") || strings.HasPrefix(r.URL.Path, "/css") {
			next.ServeHTTP(w, r)
			return
		}

		// Получаем токен из заголовка
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("Missing Authorization header for path: %s", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Отсутствует токен авторизации"})
			return
		}

		// Убираем префикс "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // префикс не был найден
			log.Printf("Invalid token format - missing Bearer prefix")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Неверный формат токена"})
			return
		}

		// Проверяем токен
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			log.Printf("Token validation failed: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Недействительный токен"})
			return
		}

		// Если токен валиден, продолжаем выполнение запроса
		log.Printf("Successful authentication for user: %s, path: %s", claims.Username, r.URL.Path)
		ctx := r.Context()
		ctx = auth.SetUserContext(ctx, claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Logging логирует запросы
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lw, r)
		log.Printf(
			"%s %s %d %v %s",
			r.Method,
			r.RequestURI,
			lw.statusCode,
			time.Since(start),
			r.UserAgent(),
		)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	if !lrw.written {
		lrw.statusCode = code
		lrw.ResponseWriter.WriteHeader(code)
		lrw.written = true
	}
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	if !lrw.written {
		// If no status has been written yet, assume 200 OK
		lrw.statusCode = http.StatusOK
		lrw.written = true
	}
	n, err := lrw.ResponseWriter.Write(b)
	return n, err
}
