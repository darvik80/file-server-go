package services

import (
	"fmt"
	"time"

	"file-server-go/internal/database"
	"file-server-go/internal/models"
)

// FileService provides file-related operations
type FileService struct {
	db *database.DB
}

// NewFileService creates new file service
func NewFileService(db *database.DB) *FileService {
	return &FileService{db: db}
}

// CreateFileRecord creates new file record in database
func (s *FileService) CreateFileRecord(appID int64, filename, filepath, url string, ttl time.Duration) (*models.File, error) {
	file := &models.File{
		ApplicationID: appID,
		Filename:      filename,
		Filepath:      filepath,
		URL:           url,
	}
	
	// Set expiration time if TTL is specified
	if ttl > 0 {
		file.ExpirationTime = time.Now().Add(ttl)
	}
	
	// Save to database
	if err := file.CreateFile(s.db.DB); err != nil {
		return nil, fmt.Errorf("failed to create file record: %w", err)
	}
	
	return file, nil
}

// FindFileByURL finds file by URL
func (s *FileService) FindFileByURL(url string) (*models.File, error) {
	file := &models.File{}
	
	if err := file.FindFileByURL(s.db.DB, url); err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}
	
	// Check if file has expired
	if !file.ExpirationTime.IsZero() && time.Now().After(file.ExpirationTime) {
		return nil, fmt.Errorf("file has expired")
	}
	
	return file, nil
}

// DeleteExpiredFiles deletes expired files
func (s *FileService) DeleteExpiredFiles() (int64, error) {
	file := &models.File{}
	return file.DeleteExpiredFiles(s.db.DB)
}