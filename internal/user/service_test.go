package user

import (
	"context"
	"log"
	"os"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestService_Register(t *testing.T) {
	mockRepo := &MockRepository{
		OnCreate: func(u *User) error {
			return nil
		},
	}

	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)
	service := NewService(mockRepo, logger)

	t.Run("should hash password and generate code", func(t *testing.T) {
		req := &CreateUserRequest{
			Email:    "service_test@example.com",
			Password: "plain_password",
			Name:     "Service Tester",
		}

		user, err := service.Register(context.Background(), req)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			t.Errorf("password was not correctly hashed: %v", err)
		}

		if len(user.EmailConfirmationCode) != 6 {
			t.Errorf("expected 6-digit verification code, got %s", user.EmailConfirmationCode)
		}

		if user.IsConfirmed || user.IsEmailConfirmed {
			t.Error("newly registered user should not be confirmed yet")
		}
	})
}
