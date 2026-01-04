package user

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	connStr := "postgres://postgres:pass@localhost:5436/auth_test_db?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	_, err = pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            email TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL,
            name TEXT,
            allowed_apps TEXT[],
            is_admin BOOLEAN,
            is_confirmed BOOLEAN,
            is_email_confirmed BOOLEAN,
            email_code TEXT,
            role TEXT,
            created_at TIMESTAMPTZ DEFAULT NOW(),
            updated_at TIMESTAMPTZ DEFAULT NOW()
        )`)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	_, err = pool.Exec(context.Background(), "DELETE FROM users")
	if err != nil {
		t.Fatalf("Failed to clean test database: %v", err)
	}

	return pool
}

func TestRepository_CreateAndGet(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	repo := NewRepository(pool)

	ctx := context.Background()
	u := &User{
		Email:       "repo_test@example.com",
		Password:    "hashed_password",
		Name:        "Repo Tester",
		AllowedApps: []string{"app1"},
		Role:        RoleUser,
	}

	t.Run("Create User", func(t *testing.T) {
		err := repo.Create(ctx, u)
		if err != nil {
			t.Errorf("Failed to create user: %v", err)
		}
		if u.ID == uuid.Nil {
			t.Error("Expected ID to be populated by RETURNING clause")
		}
	})

	t.Run("Get User By ID", func(t *testing.T) {
		found, err := repo.GetByID(ctx, u.ID)
		if err != nil {
			t.Errorf("Failed to get user: %v", err)
		}
		if found.Email != u.Email {
			t.Errorf("Expected email %s, got %s", u.Email, found.Email)
		}
	})
}
