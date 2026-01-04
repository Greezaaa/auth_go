package user

import (
	"context"
	"crypto/rand"
	"io"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetAll(ctx context.Context) ([]User, error)
	VerifyEmail(ctx context.Context, email string, code string) error
}

type Service struct {
	repo   UserRepository
	logger *log.Logger
}

func NewService(repo UserRepository, logger *log.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func generateVerificationCode() string {
	table := [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, 6)
	io.ReadAtLeast(rand.Reader, b, 6)
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func (s *Service) Register(ctx context.Context, req *CreateUserRequest) (*User, error) {
	s.logger.Printf("Service: Register attempt for %s", req.Email)

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:                 req.Email,
		Password:              string(hashed),
		Name:                  req.Name,
		AllowedApps:           req.AllowedApps,
		IsAdmin:               false,
		IsConfirmed:           false,
		IsEmailConfirmed:      false,
		EmailConfirmationCode: generateVerificationCode(),
		Role:                  RoleUser,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		s.logger.Printf("Service: DB error for %s: %v", req.Email, err)
		return nil, err
	}

	s.logger.Printf("Service: User %s created", req.Email)
	return user, nil
}

func (s *Service) FindByID(ctx context.Context, idStr string) (*User, error) {
	s.logger.Printf("Service: Searching for user ID: %s", idStr)

	uid, err := uuid.Parse(idStr)
	if err != nil {
		s.logger.Printf("Service: Invalid UUID format: %s", idStr)
		return nil, err
	}

	user, err := s.repo.GetByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	return user, nil
}
