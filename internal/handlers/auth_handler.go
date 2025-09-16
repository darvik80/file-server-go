package handlers

import (
	"encoding/json"
	"net/http"

	"file-server-go/internal/auth"
	"file-server-go/internal/config"
	"file-server-go/internal/database"
	"file-server-go/internal/models"
)

// AuthHandler handles authentication
type AuthHandler struct {
	config *config.Config
	db     *database.DB
}

// NewAuthHandler creates new authentication handler
func NewAuthHandler(cfg *config.Config, db *database.DB) *AuthHandler {
	// Set JWT key from configuration
	auth.SetJWTKey(cfg.JWTKey)

	return &AuthHandler{
		config: cfg,
		db:     db,
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

	// Check credentials in database
	user := &models.User{}
	err := user.FindUserByUsername(h.db.DB, req.Username)
	if err != nil {
		// If user not found in database, check admin credentials from config
		if req.Username == h.config.Auth.Username && req.Password == h.config.Auth.Password {
			// Generate token for admin
			token, err := auth.GenerateToken(req.Username)
			if err != nil {
				sendError(w, "Token generation error", http.StatusInternalServerError)
				return
			}

			// Send token
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(LoginResponse{Token: token})
			return
		}

		sendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check password (in real application you should use password hashing)
	if req.Password != user.Password {
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

// UserInfoResponse represents user info response
type UserInfoResponse struct {
	ID        int64           `json:"id"`
	Username  string          `json:"username"`
	Role      models.UserRole `json:"role"`
	CreatedAt string          `json:"created_at"`
}

// GetUserInfo handles getting current user info
func (h *AuthHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	// Get username from context
	username := auth.GetUserFromContext(r.Context())
	if username == "" {
		sendError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Find user in database
	user := &models.User{}
	err := user.FindUserByUsername(h.db.DB, username)
	if err != nil {
		sendError(w, "User not found", http.StatusNotFound)
		return
	}

	// Send user info
	w.Header().Set("Content-Type", "application/json")
	response := UserInfoResponse{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		sendError(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// sendError function is defined in file_handler.go
