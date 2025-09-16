package models

import (
	"database/sql"
	"time"
)

// UserRole represents user role
type UserRole string

const (
	AdminRole  UserRole = "admin"
	ReaderRole UserRole = "reader"
	WriterRole UserRole = "writer"
)

// User represents user model
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateUser creates new user in database
func (u *User) CreateUser(db *sql.DB) error {
	query := `INSERT INTO users (username, password, role, created_at) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, u.Username, u.Password, u.Role, time.Now())
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
	query := `SELECT id, username, password, role, created_at FROM users WHERE username = ?`
	row := db.QueryRow(query, username)

	return row.Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.CreatedAt)
}

// FindUserByID finds user by ID in database
func (u *User) FindUserByID(db *sql.DB, id int64) error {
	query := `SELECT id, username, password, role, created_at FROM users WHERE id = ?`
	row := db.QueryRow(query, id)

	return row.Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.CreatedAt)
}

// DeleteUser deletes user by ID from database
func (u *User) DeleteUser(db *sql.DB, id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// GetAllUsers retrieves all users from database
func (u *User) GetAllUsers(db *sql.DB) ([]User, error) {
	query := `SELECT id, username, role, created_at FROM users`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Role, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// HasAdminRole checks if user has admin role
func (u *User) HasAdminRole() bool {
	return u.Role == AdminRole
}

// HasWriterRole checks if user has writer role or higher
func (u *User) HasWriterRole() bool {
	return u.Role == AdminRole || u.Role == WriterRole
}

// HasReaderRole checks if user has reader role or higher
func (u *User) HasReaderRole() bool {
	return u.Role == AdminRole || u.Role == WriterRole || u.Role == ReaderRole
}
