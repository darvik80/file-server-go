package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// DB represents database connection
type DB struct {
	*sql.DB
}

// New creates new database connection
func New() (*DB, error) {
	// Create data directory if it doesn't exist
	dataDir := "./data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Open database connection
	dbPath := filepath.Join(dataDir, "fileserver.db")
	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * 60) // 5 minutes

	db := &DB{sqlDB}
	
	// Run migrations
	if err := db.runMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Printf("Connected to database at %s", dbPath)
	return db, nil
}

// runMigrations creates database tables if they don't exist
func (db *DB) runMigrations() error {
	// Create users table
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(usersTable)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create applications table
	appsTable := `
	CREATE TABLE IF NOT EXISTS applications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		access_key_id TEXT UNIQUE NOT NULL,
		access_key_secret TEXT NOT NULL,
		name TEXT NOT NULL,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(appsTable)
	if err != nil {
		return fmt.Errorf("failed to create applications table: %w", err)
	}

	// Create files table for tracking uploaded files
	filesTable := `
	CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		application_id INTEGER NOT NULL,
		filename TEXT NOT NULL,
		filepath TEXT NOT NULL,
		url TEXT NOT NULL,
		expiration_time TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (application_id) REFERENCES applications (id)
	);`

	_, err = db.Exec(filesTable)
	if err != nil {
		return fmt.Errorf("failed to create files table: %w", err)
	}

	// Create indexes
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_files_application_id ON files (application_id);",
		"CREATE INDEX IF NOT EXISTS idx_files_expiration ON files (expiration_time);",
		"CREATE INDEX IF NOT EXISTS idx_apps_access_key_id ON applications (access_key_id);",
	}

	for _, index := range indexes {
		_, err = db.Exec(index)
		if err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}

// Close closes database connection
func (db *DB) Close() error {
	return db.DB.Close()
}