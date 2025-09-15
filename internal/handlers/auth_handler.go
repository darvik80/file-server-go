package handlers

import (
	"encoding/json"
	"net/http"

	"file-server-go/internal/auth"
	"file-server-go/internal/config"
)

// AuthHandler handles authentication
type AuthHandler struct {
	config *config.Config
}

// NewAuthHandler creates new authentication handler
func NewAuthHandler(cfg *config.Config) *AuthHandler {
	// Set JWT key from configuration
	auth.SetJWTKey(cfg.JWTKey)

	return &AuthHandler{
		config: cfg,
	}
}

// LoginRequest represents login request data
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents successful login response
type LoginResponse struct {
	Token string `json:"token"`
}

// Login handles authentication request
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Check credentials
	if req.Username != h.config.Auth.Username || req.Password != h.config.Auth.Password {
		sendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := auth.GenerateToken(req.Username)
	if err != nil {
		sendError(w, "Token generation error", http.StatusInternalServerError)
		return
	}

	// Send token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}

// sendError function is defined in file_handler.go
