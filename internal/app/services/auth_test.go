package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/i474232898/chatserver/api/types"
	"github.com/i474232898/chatserver/internal/app"
	"github.com/i474232898/chatserver/internal/app/dto"
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
		return nil, gorm.ErrRecordNotFound
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
	tests := []struct {
		name          string
		request       *types.SigninRequest
		setupMock     func(*mockUserRepository)
		expectedError error
		checkResponse func(*testing.T, *dto.SigninResponse)
	}{
		{
			name: "successful signin",
			request: &types.SigninRequest{
				Email:    openapi_types.Email("test@example.com"),
				Password: "password123",
			},
			setupMock: func(mock *mockUserRepository) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				mock.users = map[string]*models.User{
					"test@example.com": {
						Model:    gorm.Model{ID: 1},
						Email:    "test@example.com",
						Password: string(hashedPassword),
					},
				}
			},
			expectedError: nil,
			checkResponse: func(t *testing.T, response *dto.SigninResponse) {
				if response.Token == "" {
					t.Error("Expected token to be generated")
				}

				// Verify token is valid JWT
				token, err := jwt.ParseWithClaims(response.Token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte("secret"), nil
				})
				if err != nil {
					t.Errorf("Expected valid JWT token, got error: %v", err)
				}

				claims, ok := token.Claims.(*CustomClaims)
				if !ok {
					t.Error("Expected CustomClaims")
				}

				if claims.ID != 1 {
					t.Errorf("Expected user ID 1, got %d", claims.ID)
				}

				if claims.ExpiresAt.Time.Before(time.Now()) {
					t.Error("Token should not be expired")
				}
			},
		},
		{
			name: "user not found",
			request: &types.SigninRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			setupMock: func(mock *mockUserRepository) {
				mock.users = make(map[string]*models.User)
			},
			expectedError: app.ErrUserNotFound,
			checkResponse: func(t *testing.T, response *dto.SigninResponse) {
				if response != nil {
					t.Error("Expected nil response when user not found")
				}
			},
		},
		{
			name: "invalid password",
			request: &types.SigninRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setupMock: func(mock *mockUserRepository) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
				mock.users = map[string]*models.User{
					"test@example.com": {
						Model: gorm.Model{ID: 1},
						Email: "test@example.com",
						Password: string(hashedPassword),
					},
				}
			},
			expectedError: app.ErrInvalidCredentials,
			checkResponse: func(t *testing.T, response *dto.SigninResponse) {
				if response != nil {
					t.Error("Expected nil response when credentials are invalid")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUserRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			service := NewAuthService(mockRepo)
			response, err := service.Signin(context.Background(), tt.request)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.expectedError)
					return
				}
				if !errors.Is(err, tt.expectedError) {
					t.Errorf("Expected error %v, got %v", tt.expectedError, err)
					return
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
					return
				}
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}
