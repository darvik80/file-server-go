package models

import (
	"database/sql"
	"strings"
	"time"
)

// AppPermission represents application permissions
type AppPermission string

const (
	ReadPermission  AppPermission = "read"
	WritePermission AppPermission = "write"
)

// Application represents application model
type Application struct {
	ID              int64           `json:"id"`
	AccessKeyID     string          `json:"access_key_id"`
	AccessKeySecret string          `json:"access_key_secret"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Permissions     []AppPermission `json:"permissions"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// CreateApplication creates new application
func (a *Application) CreateApplication(db *sql.DB) error {
	// Convert permissions to string for storage
	permStr := ""
	for i, perm := range a.Permissions {
		if i > 0 {
			permStr += ","
		}
		permStr += string(perm)
	}

	query := `INSERT INTO applications (access_key_id, access_key_secret, name, description, permissions, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, a.AccessKeyID, a.AccessKeySecret, a.Name, a.Description, permStr, time.Now(), time.Now())
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	a.ID = id
	return nil
}

// FindApplicationByAccessKeyID finds application by access key ID
func (a *Application) FindApplicationByAccessKeyID(db *sql.DB, accessKeyID string) error {
	query := `SELECT id, access_key_id, access_key_secret, name, description, permissions, created_at, updated_at 
	          FROM applications WHERE access_key_id = ?`
	row := db.QueryRow(query, accessKeyID)

	var permStr string
	err := row.Scan(&a.ID, &a.AccessKeyID, &a.AccessKeySecret, &a.Name, &a.Description, &permStr, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return err
	}

	// Parse permissions from string
	a.Permissions = parsePermissions(permStr)
	return nil
}

// FindApplicationByID finds application by ID
func (a *Application) FindApplicationByID(db *sql.DB, id int64) error {
	query := `SELECT id, access_key_id, access_key_secret, name, description, permissions, created_at, updated_at 
	          FROM applications WHERE id = ?`
	row := db.QueryRow(query, id)

	var permStr string
	err := row.Scan(&a.ID, &a.AccessKeyID, &a.AccessKeySecret, &a.Name, &a.Description, &permStr, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return err
	}

	// Parse permissions from string
	a.Permissions = parsePermissions(permStr)
	return nil
}

// DeleteApplication deletes application by ID
func (a *Application) DeleteApplication(db *sql.DB, id int64) error {
	query := `DELETE FROM applications WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// GetAllApplications retrieves all applications
func (a *Application) GetAllApplications(db *sql.DB) ([]Application, error) {
	query := `SELECT id, access_key_id, name, description, permissions, created_at, updated_at 
	          FROM applications`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []Application
	for rows.Next() {
		var app Application
		var permStr string
		err := rows.Scan(&app.ID, &app.AccessKeyID, &app.Name, &app.Description, &permStr, &app.CreatedAt, &app.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Parse permissions from string
		app.Permissions = parsePermissions(permStr)
		applications = append(applications, app)
	}

	return applications, nil
}

// ValidateCredentials validates application credentials
func (a *Application) ValidateCredentials(db *sql.DB, accessKeyID, accessKeySecret string) (bool, error) {
	query := `SELECT COUNT(*) FROM applications WHERE access_key_id = ? AND access_key_secret = ?`
	var count int
	err := db.QueryRow(query, accessKeyID, accessKeySecret).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// HasReadPermission checks if application has read permission
func (a *Application) HasReadPermission() bool {
	for _, perm := range a.Permissions {
		if perm == ReadPermission {
			return true
		}
	}
	return false
}

// HasWritePermission checks if application has write permission
func (a *Application) HasWritePermission() bool {
	for _, perm := range a.Permissions {
		if perm == WritePermission {
			return true
		}
	}
	return false
}

// parsePermissions parses permissions from comma-separated string
func parsePermissions(permStr string) []AppPermission {
	if permStr == "" {
		return []AppPermission{}
	}

	var permissions []AppPermission
	// Simple split by comma
	parts := strings.Split(permStr, ",")
	for _, part := range parts {
		perm := AppPermission(strings.TrimSpace(part))
		permissions = append(permissions, perm)
	}

	return permissions
}
