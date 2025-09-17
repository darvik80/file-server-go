package utils

import (
	"fmt"
	"net/http"

	"file-server-go/internal/auth"
	"file-server-go/internal/database"
	"file-server-go/internal/models"
)

// GetCurrentUser retrieves current authenticated user
func GetCurrentUser(r *http.Request, db *database.DB) (*models.User, error) {
	currentUsername := auth.GetUserFromContext(r.Context())
	if currentUsername == "" {
		return nil, fmt.Errorf("user not authenticated")
	}

	currentUser := &models.User{}
	if err := currentUser.FindUserByUsername(db.DB, currentUsername); err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return currentUser, nil
}

// CheckAdminRole checks if current user has admin role
func CheckAdminRole(r *http.Request, db *database.DB) (*models.User, error) {
	user, err := GetCurrentUser(r, db)
	if err != nil {
		return nil, err
	}

	if !user.HasAdminRole() {
		return nil, fmt.Errorf("access denied. Admin rights required")
	}

	return user, nil
}

// CheckWriterRole checks if current user has writer role or higher
func CheckWriterRole(r *http.Request, db *database.DB) (*models.User, error) {
	user, err := GetCurrentUser(r, db)
	if err != nil {
		return nil, err
	}

	if !user.HasWriterRole() {
		return nil, fmt.Errorf("access denied. Writer or admin rights required")
	}

	return user, nil
}
