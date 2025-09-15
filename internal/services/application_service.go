package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"file-server-go/internal/database"
	"file-server-go/internal/models"
)

// ApplicationService provides application-related operations
type ApplicationService struct {
	db *database.DB
}

// NewApplicationService creates new application service
func NewApplicationService(db *database.DB) *ApplicationService {
	return &ApplicationService{db: db}
}

// RegisterApplication registers new application
func (s *ApplicationService) RegisterApplication(name, description string) (*models.Application, error) {
	// Generate unique access key ID and secret
	accessKeyID, err := generateRandomString(16)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access key ID: %w", err)
	}
	
	accessKeySecret, err := generateRandomString(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access key secret: %w", err)
	}
	
	// Create application model
	app := &models.Application{
		AccessKeyID:     accessKeyID,
		AccessKeySecret: accessKeySecret,
		Name:            name,
		Description:     description,
	}
	
	// Save to database
	if err := app.CreateApplication(s.db.DB); err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}
	
	return app, nil
}

// AuthenticateApplication authenticates application by credentials
func (s *ApplicationService) AuthenticateApplication(accessKeyID, accessKeySecret string) (*models.Application, error) {
	app := &models.Application{}
	
	// Find application by access key ID
	if err := app.FindApplicationByAccessKeyID(s.db.DB, accessKeyID); err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}
	
	// Validate credentials
	valid, err := app.ValidateCredentials(s.db.DB, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to validate credentials: %w", err)
	}
	
	if !valid {
		return nil, fmt.Errorf("invalid credentials")
	}
	
	return app, nil
}

// generateRandomString generates random string of specified length
func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}