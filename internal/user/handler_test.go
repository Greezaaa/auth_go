package user

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
)

type MockRepository struct {
	OnCreate func(u *User) error
}

func (m *MockRepository) Create(ctx context.Context, u *User) error {
	return m.OnCreate(u)
}
func (m *MockRepository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) { return nil, nil }
func (m *MockRepository) GetAll(ctx context.Context) ([]User, error)               { return nil, nil }
func (m *MockRepository) VerifyEmail(ctx context.Context, e, c string) error       { return nil }

func TestRegisterHandler(t *testing.T) {
	mockRepo := &MockRepository{
		OnCreate: func(u *User) error {
			u.ID = uuid.New()
			return nil
		},
	}

	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)
	service := NewService(mockRepo, logger)
	h := NewHandler(service, logger)

	t.Run("successful registration", func(t *testing.T) {
		userReq := CreateUserRequest{
			Email:    "test@example.com",
			Password: "password123",
			Name:     "Test User",
		}
		body, _ := json.Marshal(userReq)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		h.Register(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", rr.Code)
		}

		var savedUser User
		json.NewDecoder(rr.Body).Decode(&savedUser)
		if savedUser.Email != userReq.Email {
			t.Errorf("expected email %s, got %s", userReq.Email, savedUser.Email)
		}
	})
}
