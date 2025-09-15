package models

import (
	"database/sql"
	"time"
)

// User represents user model
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateUser creates new user in database
func (u *User) CreateUser(db *sql.DB) error {
	query := `INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`
	result, err := db.Exec(query, u.Username, u.Password, time.Now())
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = id
	return nil
}

// FindUserByUsername finds user by username in database
func (u *User) FindUserByUsername(db *sql.DB, username string) error {
	query := `SELECT id, username, password, created_at FROM users WHERE username = ?`
	row := db.QueryRow(query, username)

	return row.Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt)
}
