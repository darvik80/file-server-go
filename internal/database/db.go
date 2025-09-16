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

	// Create default admin user if users table is empty
	if err := db.createDefaultAdminUser(); err != nil {
		log.Printf("Warning: failed to create default admin user: %v", err)
	}

	log.Printf("Connected to database at %s", dbPath)
	return db, nil
}

// runMigrations creates database tables if they don't exist
func (db *DB) runMigrations() error {
	// Create users table with role column
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'reader',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(usersTable)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Add role column to existing users table if it doesn't exist
	_, err = db.Exec("ALTER TABLE users ADD COLUMN role TEXT NOT NULL DEFAULT 'reader'")
	if err != nil {
		// Column might already exist, ignore error
		log.Printf("Note: role column might already exist in users table")
	}

	// Create applications table with permissions column
	appsTable := `
	CREATE TABLE IF NOT EXISTS applications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		access_key_id TEXT UNIQUE NOT NULL,
		access_key_secret TEXT NOT NULL,
		name TEXT NOT NULL,
		description TEXT,
		permissions TEXT NOT NULL DEFAULT 'read',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(appsTable)
	if err != nil {
		return fmt.Errorf("failed to create applications table: %w", err)
	}

	// Add permissions column to existing applications table if it doesn't exist
	_, err = db.Exec("ALTER TABLE applications ADD COLUMN permissions TEXT NOT NULL DEFAULT 'read'")
	if err != nil {
		// Column might already exist, ignore error
		log.Printf("Note: permissions column might already exist in applications table")
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

// createDefaultAdminUser creates default admin user if users table is empty
func (db *DB) createDefaultAdminUser() error {
	// Check if users table is empty
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	// If users table is not empty, do nothing
	if count > 0 {
		return nil
	}

	// Create default admin user
	_, err = db.Exec(
		"INSERT INTO users (username, password, role, created_at) VALUES (?, ?, ?, ?)",
		"admin",
		"admin",
		"admin",
		"CURRENT_TIMESTAMP",
	)
	if err != nil {
		return fmt.Errorf("failed to create default admin user: %w", err)
	}

	log.Println("Created default admin user (admin:admin)")
	return nil
}

// Close closes database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
