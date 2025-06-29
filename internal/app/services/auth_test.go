package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/i474232898/chatserver/api/types"
	"github.com/i474232898/chatserver/internal/app/repositories/models"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type mockUserRepository struct {
	users map[string]*models.User
}

func (r *mockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, ok := r.users[email]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}
func (r *mockUserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	r.users[user.Email] = user
	return user, nil
}
func (r *mockUserRepository) GetById(ctx context.Context, id int64) (*models.User, error) {
	user, ok := r.users["abc@email.com"]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func TestAuthService_Signin(t *testing.T) {
	r := &mockUserRepository{}
	password := "password123"
	email := "abc@email.com"
	setupMock(r, password, email)

	service := NewAuthService(r)

	request := &types.SigninRequest{
		Email:    openapi_types.Email(email),
		Password: password,
	}

	response, err := service.Signin(context.Background(), request)
	if err != nil {
		t.Fatalf("Signin failed: %v", err)
	}
	fmt.Println(response, "<<<<")
	// Check results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response == nil {
		t.Error("Expected response, got nil")
	}

	if response.Token == "" {
		t.Error("Expected token to be generated")
	}
}

func setupMock(r *mockUserRepository, pass, email string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	r.users = map[string]*models.User{
		email: {
			Model:    gorm.Model{ID: 1},
			Email:    email,
			Password: string(hashedPassword),
		},
	}
}
