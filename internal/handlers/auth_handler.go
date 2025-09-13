package handlers

import (
	"encoding/json"
	"net/http"

	"file-server-go/internal/auth"
	"file-server-go/internal/config"
)

// AuthHandler обрабатывает аутентификацию
type AuthHandler struct {
	config *config.Config
}

// NewAuthHandler создает новый обработчик аутентификации
func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		config: cfg,
	}
}

// LoginRequest представляет данные запроса на вход
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse представляет ответ на успешный вход
type LoginResponse struct {
	Token string `json:"token"`
}

// Login обрабатывает запрос на аутентификацию
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Проверяем учетные данные
	if req.Username != h.config.Auth.Username || req.Password != h.config.Auth.Password {
		sendError(w, "Неверные учетные данные", http.StatusUnauthorized)
		return
	}

	// Генерируем токен
	token, err := auth.GenerateToken(req.Username)
	if err != nil {
		sendError(w, "Ошибка создания токена", http.StatusInternalServerError)
		return
	}

	// Отправляем токен
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}
