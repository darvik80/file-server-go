package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"file-server-go/internal/auth"
	"file-server-go/internal/database"
	"file-server-go/internal/models"
)

// UserHandler handles user-related requests
type UserHandler struct {
	db *database.DB
}

// NewUserHandler creates new user handler
func NewUserHandler(db *database.DB) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

// GetUsersResponse represents users list response
type GetUsersResponse struct {
	Users []models.User `json:"users"`
}

// GetUsers handles getting all users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	// Get current user from context
	currentUsername := auth.GetUserFromContext(r.Context())
	if currentUsername == "" {
		sendError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Find current user
	currentUser := &models.User{}
	if err := currentUser.FindUserByUsername(h.db.DB, currentUsername); err != nil {
		sendError(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if user has admin role
	if !currentUser.HasAdminRole() {
		sendError(w, "Access denied. Admin rights required", http.StatusForbidden)
		return
	}

	// Get all users
	user := &models.User{}
	users, err := user.GetAllUsers(h.db.DB)
	if err != nil {
		sendError(w, "Failed to get users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	response := GetUsersResponse{
		Users: users,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		sendError(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// DeleteUser handles deleting user
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		sendError(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	// Get current user from context
	currentUsername := auth.GetUserFromContext(r.Context())
	if currentUsername == "" {
		sendError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Find current user
	currentUser := &models.User{}
	if err := currentUser.FindUserByUsername(h.db.DB, currentUsername); err != nil {
		sendError(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if user has admin role
	if !currentUser.HasAdminRole() {
		sendError(w, "Access denied. Admin rights required", http.StatusForbidden)
		return
	}

	// Get user ID from query parameter
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		sendError(w, "User ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Find user to be deleted
	userToDelete := &models.User{}
	if err := userToDelete.FindUserByID(h.db.DB, id); err != nil {
		sendError(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if user tries to delete themselves
	if userToDelete.Username == currentUsername {
		sendError(w, "You cannot delete yourself", http.StatusForbidden)
		return
	}

	// Delete user
	if err := userToDelete.DeleteUser(h.db.DB, id); err != nil {
		sendError(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send success response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "User deleted successfully",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		sendError(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// CreateUserRequest represents create user request
type CreateUserRequest struct {
	Username string          `json:"username"`
	Password string          `json:"password"`
	Role     models.UserRole `json:"role"`
}

// CreateUser handles creating new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	// Get current user from context
	currentUsername := auth.GetUserFromContext(r.Context())
	if currentUsername == "" {
		sendError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Find current user
	currentUser := &models.User{}
	if err := currentUser.FindUserByUsername(h.db.DB, currentUsername); err != nil {
		sendError(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if user has admin role
	if !currentUser.HasAdminRole() {
		sendError(w, "Access denied. Admin rights required", http.StatusForbidden)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" {
		sendError(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = models.ReaderRole
	}

	// Validate role
	validRoles := map[models.UserRole]bool{
		models.AdminRole:  true,
		models.ReaderRole: true,
		models.WriterRole: true,
	}

	if !validRoles[req.Role] {
		sendError(w, "Invalid role. Valid roles are: admin, reader, writer", http.StatusBadRequest)
		return
	}

	// Create user
	user := &models.User{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
	}

	if err := user.CreateUser(h.db.DB); err != nil {
		sendError(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "User created successfully",
		"user": map[string]interface{}{
			"id":         user.ID,
			"username":   user.Username,
			"role":       user.Role,
			"created_at": user.CreatedAt,
		},
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		sendError(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// ChangePasswordRequest represents change password request
type ChangePasswordRequest struct {
	Username    string `json:"username"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// ChangePassword handles changing user password
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	// Get current user from context
	currentUsername := auth.GetUserFromContext(r.Context())
	if currentUsername == "" {
		sendError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Find current user
	currentUser := &models.User{}
	if err := currentUser.FindUserByUsername(h.db.DB, currentUsername); err != nil {
		sendError(w, "User not found", http.StatusNotFound)
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Username == "" || req.NewPassword == "" {
		sendError(w, "Username and new password are required", http.StatusBadRequest)
		return
	}

	// Find user in database
	user := &models.User{}
	err := user.FindUserByUsername(h.db.DB, req.Username)
	if err != nil {
		sendError(w, "User not found", http.StatusNotFound)
		return
	}

	// Проверка прав доступа: администратор может менять пароль любому пользователю,
	// обычный пользователь может менять только свой пароль и должен предоставить старый пароль
	if currentUser.HasAdminRole() {
		// Администратор может менять пароль любому пользователю без проверки старого пароля
	} else if currentUser.Username == user.Username {
		// Обычный пользователь может менять только свой пароль и должен предоставить старый пароль
		if req.OldPassword == "" {
			sendError(w, "Old password is required", http.StatusBadRequest)
			return
		}
		// Check old password
		if req.OldPassword != user.Password {
			sendError(w, "Invalid old password", http.StatusUnauthorized)
			return
		}
	} else {
		// Обычный пользователь пытается изменить пароль другого пользователя
		sendError(w, "Access denied. You can only change your own password", http.StatusForbidden)
		return
	}

	// Update password
	user.Password = req.NewPassword
	query := `UPDATE users SET password = ? WHERE id = ?`
	_, err = h.db.DB.Exec(query, user.Password, user.ID)
	if err != nil {
		sendError(w, "Failed to update password: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send success response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "Password changed successfully",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		sendError(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
