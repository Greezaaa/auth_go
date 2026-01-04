package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, u *User) error {
	query := `
		INSERT INTO users (email, password, name, allowed_apps, is_admin, is_confirmed, is_email_confirmed, email_code, role) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`

	return r.pool.QueryRow(ctx, query,
		u.Email, u.Password, u.Name, u.AllowedApps, u.IsAdmin, u.IsConfirmed, u.IsEmailConfirmed, u.EmailConfirmationCode, u.Role,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *Repository) GetAll(ctx context.Context) ([]User, error) {
	query := `SELECT id, email, name, role, created_at FROM users`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `
		SELECT id, email, name, allowed_apps, is_admin, is_confirmed, 
		       is_email_confirmed, role, created_at, updated_at 
		FROM users WHERE id = $1`

	var u User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.Name, &u.AllowedApps, &u.IsAdmin,
		&u.IsConfirmed, &u.IsEmailConfirmed, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) VerifyEmail(ctx context.Context, email string, code string) error {
	query := `
		UPDATE users 
		SET is_email_confirmed = true, email_code = NULL 
		WHERE email = $1 AND email_code = $2`

	result, err := r.pool.Exec(ctx, query, email, code)
	if err != nil {
		return err
	}

	// Check if any row was actually updated
	if result.RowsAffected() == 0 {
		return fmt.Errorf("invalid email or verification code")
	}
	return nil
}
