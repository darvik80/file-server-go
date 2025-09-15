package models

import (
	"database/sql"
	"time"
)

// Application represents application model
type Application struct {
	ID              int64     `json:"id"`
	AccessKeyID     string    `json:"access_key_id"`
	AccessKeySecret string    `json:"access_key_secret"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CreateApplication creates new application
func (a *Application) CreateApplication(db *sql.DB) error {
	query := `INSERT INTO applications (access_key_id, access_key_secret, name, description, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, a.AccessKeyID, a.AccessKeySecret, a.Name, a.Description, time.Now(), time.Now())
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
	query := `SELECT id, access_key_id, access_key_secret, name, description, created_at, updated_at 
	          FROM applications WHERE access_key_id = ?`
	row := db.QueryRow(query, accessKeyID)

	return row.Scan(&a.ID, &a.AccessKeyID, &a.AccessKeySecret, &a.Name, &a.Description, &a.CreatedAt, &a.UpdatedAt)
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
