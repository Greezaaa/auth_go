package user

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleGuest      Role = "guest"
	RoleUser       Role = "user"
	RoleEditor     Role = "editor"
	RoleAdmin      Role = "admin"
	RoleSuperAdmin Role = "superAdmin"
)

type User struct {
	ID                    uuid.UUID `json:"id"`
	Email                 string    `json:"email"`
	Password              string    `json:"-"`
	Name                  string    `json:"name"`
	AllowedApps           []string  `json:"allowed_apps"`
	IsAdmin               bool      `json:"is_admin"`
	IsConfirmed           bool      `json:"is_confirmed"`
	IsEmailConfirmed      bool      `json:"is_email_confirmed"`
	EmailConfirmationCode string    `json:"-"`
	Role                  Role      `json:"role"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Name        string   `json:"name"`
	AllowedApps []string `json:"allowed_apps"`
}
