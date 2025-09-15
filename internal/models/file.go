package models

import (
	"database/sql"
	"time"
)

// File represents uploaded file model
type File struct {
	ID             int64     `json:"id"`
	ApplicationID  int64     `json:"application_id"`
	Filename       string    `json:"filename"`
	Filepath       string    `json:"filepath"`
	URL            string    `json:"url"`
	ExpirationTime time.Time `json:"expiration_time,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// CreateFile creates new file record in database
func (f *File) CreateFile(db *sql.DB) error {
	query := `INSERT INTO files (application_id, filename, filepath, url, expiration_time, created_at) 
	          VALUES (?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, f.ApplicationID, f.Filename, f.Filepath, f.URL, f.ExpirationTime, time.Now())
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	f.ID = id
	return nil
}

// FindFileByURL finds file by URL in database
func (f *File) FindFileByURL(db *sql.DB, url string) error {
	query := `SELECT id, application_id, filename, filepath, url, expiration_time, created_at 
	          FROM files WHERE url = ?`
	row := db.QueryRow(query, url)

	return row.Scan(&f.ID, &f.ApplicationID, &f.Filename, &f.Filepath, &f.URL, &f.ExpirationTime, &f.CreatedAt)
}

// DeleteExpiredFiles deletes expired files from database
func (f *File) DeleteExpiredFiles(db *sql.DB) (int64, error) {
	query := `DELETE FROM files WHERE expiration_time IS NOT NULL AND expiration_time < ?`
	result, err := db.Exec(query, time.Now())
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
