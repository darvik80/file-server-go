package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"file-server-go/internal/database"
	"file-server-go/internal/models"
	"file-server-go/internal/services"
)

// AppAuthContextKey is the context key for application authentication
type AppAuthContextKey string

const (
	// AppContextKey is the context key for application
	AppContextKey AppAuthContextKey = "app"
)

// AppAuth проверяет учетные данные приложения
func AppAuth(db *database.DB) func(http.Handler) http.Handler {
	appService := services.NewApplicationService(db)
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем заголовок авторизации
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// Если нет заголовка авторизации, передаем запрос дальше
				next.ServeHTTP(w, r)
				return
			}

			// Проверяем, является ли это авторизацией приложения
			if strings.HasPrefix(authHeader, "App ") {
				// Извлекаем учетные данные
				credentials := strings.TrimPrefix(authHeader, "App ")
				parts := strings.Split(credentials, ":")
				if len(parts) != 2 {
					log.Printf("Invalid application credentials format")
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(map[string]string{"error": "Неверный формат учетных данных приложения"})
					return
				}

				accessKeyID := parts[0]
				accessKeySecret := parts[1]

				// Аутентифицируем приложение
				app, err := appService.AuthenticateApplication(accessKeyID, accessKeySecret)
				if err != nil {
					log.Printf("Application authentication failed: %v", err)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(map[string]string{"error": "Недействительные учетные данные приложения"})
					return
				}

				// Если аутентификация успешна, добавляем приложение в контекст
				log.Printf("Successful application authentication for app: %s", app.Name)
				ctx := context.WithValue(r.Context(), AppContextKey, app)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Если это не авторизация приложения, передаем запрос дальше
			next.ServeHTTP(w, r)
		})
	}
}

// GetAppFromContext получает приложение из контекста
func GetAppFromContext(ctx context.Context) *models.Application {
	if app, ok := ctx.Value(AppContextKey).(*models.Application); ok {
		return app
	}
	return nil
}