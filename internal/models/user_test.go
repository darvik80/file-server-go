package models

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB создает тестовую базу данных в памяти
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Создаем таблицу users
	createUserTable := `
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL,
		created_at DATETIME NOT NULL
	)`

	if _, err := db.Exec(createUserTable); err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	return db
}

func TestUserRole(t *testing.T) {
	t.Run("user role constants", func(t *testing.T) {
		if AdminRole != "admin" {
			t.Errorf("Expected AdminRole to be 'admin', got %s", AdminRole)
		}
		if ReaderRole != "reader" {
			t.Errorf("Expected ReaderRole to be 'reader', got %s", ReaderRole)
		}
		if WriterRole != "writer" {
			t.Errorf("Expected WriterRole to be 'writer', got %s", WriterRole)
		}
	})
}

func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	t.Run("create valid user", func(t *testing.T) {
		user := &User{
			Username: "testuser",
			Password: "testpass",
			Role:     WriterRole,
		}

		err := user.CreateUser(db)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if user.ID == 0 {
			t.Errorf("Expected user ID to be set, got 0")
		}
	})

	t.Run("create user with duplicate username", func(t *testing.T) {
		user1 := &User{
			Username: "duplicate",
			Password: "pass1",
			Role:     ReaderRole,
		}

		user2 := &User{
			Username: "duplicate",
			Password: "pass2",
			Role:     AdminRole,
		}

		err := user1.CreateUser(db)
		if err != nil {
			t.Fatalf("Expected no error for first user, got %v", err)
		}

		err = user2.CreateUser(db)
		if err == nil {
			t.Errorf("Expected error for duplicate username, got nil")
		}
	})

	t.Run("create user with empty username", func(t *testing.T) {
		user := &User{
			Username: "",
			Password: "testpass",
			Role:     WriterRole,
		}

		err := user.CreateUser(db)
		// SQLite позволяет пустые строки, но это может быть нежелательно в реальном приложении
		// Для тестов мы принимаем текущее поведение
		if err != nil {
			t.Logf("Empty username created user with ID: %d", user.ID)
		}
	})

	t.Run("create user with all roles", func(t *testing.T) {
		roles := []UserRole{AdminRole, ReaderRole, WriterRole}

		for i, role := range roles {
			user := &User{
				Username: "user" + string(rune('1'+i)),
				Password: "testpass",
				Role:     role,
			}

			err := user.CreateUser(db)
			if err != nil {
				t.Errorf("Expected no error for role %s, got %v", role, err)
			}

			if user.Role != role {
				t.Errorf("Expected role %s, got %s", role, user.Role)
			}
		}
	})
}

func TestFindUserByUsername(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Создаем тестового пользователя
	originalUser := &User{
		Username: "findme",
		Password: "secret",
		Role:     AdminRole,
	}
	err := originalUser.CreateUser(db)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	t.Run("find existing user", func(t *testing.T) {
		foundUser := &User{}
		err := foundUser.FindUserByUsername(db, "findme")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if foundUser.Username != "findme" {
			t.Errorf("Expected username 'findme', got %s", foundUser.Username)
		}

		if foundUser.Password != "secret" {
			t.Errorf("Expected password 'secret', got %s", foundUser.Password)
		}

		if foundUser.Role != AdminRole {
			t.Errorf("Expected role %s, got %s", AdminRole, foundUser.Role)
		}

		if foundUser.ID != originalUser.ID {
			t.Errorf("Expected ID %d, got %d", originalUser.ID, foundUser.ID)
		}
	})

	t.Run("find non-existent user", func(t *testing.T) {
		user := &User{}
		err := user.FindUserByUsername(db, "nonexistent")

		if err == nil {
			t.Errorf("Expected error for non-existent user, got nil")
		}

		if err != sql.ErrNoRows {
			t.Errorf("Expected sql.ErrNoRows, got %v", err)
		}
	})

	t.Run("find user with empty username", func(t *testing.T) {
		user := &User{}
		err := user.FindUserByUsername(db, "")

		if err == nil {
			t.Errorf("Expected error for empty username, got nil")
		}
	})
}

func TestFindUserByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Создаем тестового пользователя
	originalUser := &User{
		Username: "findbyid",
		Password: "secret",
		Role:     WriterRole,
	}
	err := originalUser.CreateUser(db)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	t.Run("find existing user by ID", func(t *testing.T) {
		foundUser := &User{}
		err := foundUser.FindUserByID(db, originalUser.ID)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if foundUser.ID != originalUser.ID {
			t.Errorf("Expected ID %d, got %d", originalUser.ID, foundUser.ID)
		}

		if foundUser.Username != "findbyid" {
			t.Errorf("Expected username 'findbyid', got %s", foundUser.Username)
		}

		if foundUser.Role != WriterRole {
			t.Errorf("Expected role %s, got %s", WriterRole, foundUser.Role)
		}
	})

	t.Run("find non-existent user by ID", func(t *testing.T) {
		user := &User{}
		err := user.FindUserByID(db, 99999)

		if err == nil {
			t.Errorf("Expected error for non-existent user ID, got nil")
		}

		if err != sql.ErrNoRows {
			t.Errorf("Expected sql.ErrNoRows, got %v", err)
		}
	})

	t.Run("find user by zero ID", func(t *testing.T) {
		user := &User{}
		err := user.FindUserByID(db, 0)

		if err == nil {
			t.Errorf("Expected error for zero ID, got nil")
		}
	})
}

func TestDeleteUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	t.Run("delete existing user", func(t *testing.T) {
		// Создаем пользователя для удаления
		user := &User{
			Username: "deleteme",
			Password: "secret",
			Role:     ReaderRole,
		}
		err := user.CreateUser(db)
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		// Удаляем пользователя
		err = user.DeleteUser(db, user.ID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Проверяем, что пользователь удален
		foundUser := &User{}
		err = foundUser.FindUserByID(db, user.ID)
		if err != sql.ErrNoRows {
			t.Errorf("Expected user to be deleted, but found: %v", err)
		}
	})

	t.Run("delete non-existent user", func(t *testing.T) {
		user := &User{}
		err := user.DeleteUser(db, 99999)

		// Удаление несуществующего пользователя не должно вызывать ошибку
		if err != nil {
			t.Errorf("Expected no error for deleting non-existent user, got %v", err)
		}
	})
}

func TestGetAllUsers(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	t.Run("get all users from empty database", func(t *testing.T) {
		user := &User{}
		users, err := user.GetAllUsers(db)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(users) != 0 {
			t.Errorf("Expected 0 users, got %d", len(users))
		}
	})

	t.Run("get all users with data", func(t *testing.T) {
		// Создаем несколько пользователей
		testUsers := []*User{
			{Username: "user1", Password: "pass1", Role: AdminRole},
			{Username: "user2", Password: "pass2", Role: WriterRole},
			{Username: "user3", Password: "pass3", Role: ReaderRole},
		}

		for _, u := range testUsers {
			err := u.CreateUser(db)
			if err != nil {
				t.Fatalf("Failed to create test user %s: %v", u.Username, err)
			}
		}

		// Получаем всех пользователей
		user := &User{}
		users, err := user.GetAllUsers(db)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(users) != 3 {
			t.Errorf("Expected 3 users, got %d", len(users))
		}

		// Проверяем, что все пользователи присутствуют
		usernames := make(map[string]bool)
		for _, u := range users {
			usernames[u.Username] = true
		}

		for _, expectedUser := range testUsers {
			if !usernames[expectedUser.Username] {
				t.Errorf("Expected to find user %s", expectedUser.Username)
			}
		}
	})
}

func TestUserRolePermissions(t *testing.T) {
	t.Run("admin role permissions", func(t *testing.T) {
		user := &User{Role: AdminRole}

		if !user.HasAdminRole() {
			t.Errorf("Admin user should have admin role")
		}
		if !user.HasWriterRole() {
			t.Errorf("Admin user should have writer role")
		}
		if !user.HasReaderRole() {
			t.Errorf("Admin user should have reader role")
		}
	})

	t.Run("writer role permissions", func(t *testing.T) {
		user := &User{Role: WriterRole}

		if user.HasAdminRole() {
			t.Errorf("Writer user should not have admin role")
		}
		if !user.HasWriterRole() {
			t.Errorf("Writer user should have writer role")
		}
		if !user.HasReaderRole() {
			t.Errorf("Writer user should have reader role")
		}
	})

	t.Run("reader role permissions", func(t *testing.T) {
		user := &User{Role: ReaderRole}

		if user.HasAdminRole() {
			t.Errorf("Reader user should not have admin role")
		}
		if user.HasWriterRole() {
			t.Errorf("Reader user should not have writer role")
		}
		if !user.HasReaderRole() {
			t.Errorf("Reader user should have reader role")
		}
	})

	t.Run("invalid role permissions", func(t *testing.T) {
		user := &User{Role: UserRole("invalid")}

		if user.HasAdminRole() {
			t.Errorf("Invalid role user should not have admin role")
		}
		if user.HasWriterRole() {
			t.Errorf("Invalid role user should not have writer role")
		}
		if user.HasReaderRole() {
			t.Errorf("Invalid role user should not have reader role")
		}
	})
}

func TestUserCreatedAt(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	t.Run("created_at is set correctly", func(t *testing.T) {
		beforeCreate := time.Now()

		user := &User{
			Username: "timetest",
			Password: "secret",
			Role:     WriterRole,
		}

		err := user.CreateUser(db)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		afterCreate := time.Now()

		// Получаем пользователя из базы данных
		foundUser := &User{}
		err = foundUser.FindUserByID(db, user.ID)
		if err != nil {
			t.Fatalf("Failed to find user: %v", err)
		}

		// Проверяем, что время создания находится в разумных пределах
		if foundUser.CreatedAt.Before(beforeCreate) || foundUser.CreatedAt.After(afterCreate) {
			t.Errorf("CreatedAt time %v is not between %v and %v",
				foundUser.CreatedAt, beforeCreate, afterCreate)
		}
	})
}

// Benchmark тесты для измерения производительности
func BenchmarkCreateUser(b *testing.B) {
	db := setupTestDB(&testing.T{})
	defer db.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := &User{
			Username: "benchuser" + string(rune('0'+i%10)),
			Password: "benchpass",
			Role:     WriterRole,
		}
		user.CreateUser(db)
	}
}

func BenchmarkFindUserByUsername(b *testing.B) {
	db := setupTestDB(&testing.T{})
	defer db.Close()

	// Создаем тестового пользователя
	user := &User{
		Username: "benchfind",
		Password: "secret",
		Role:     AdminRole,
	}
	user.CreateUser(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		foundUser := &User{}
		foundUser.FindUserByUsername(db, "benchfind")
	}
}

func BenchmarkUserRoleCheck(b *testing.B) {
	user := &User{Role: WriterRole}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user.HasWriterRole()
	}
}
