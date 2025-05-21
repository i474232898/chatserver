package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/i474232898/chatserver/api/types"
	"github.com/i474232898/chatserver/internal/app"
	"github.com/i474232898/chatserver/internal/app/dto"
	"github.com/i474232898/chatserver/internal/app/repositories"
	"github.com/i474232898/chatserver/internal/app/repositories/models"
)

type AuthService interface {
	Signup(ctx context.Context, user *types.SignupRequest) (*dto.SignupResponse, error)
	Signin(ctx context.Context, user *types.SigninRequest) (*dto.SigninResponse, error)
}

type authService struct {
	userRepository repositories.UserRepository
}

func NewAuthService(userRepository repositories.UserRepository) AuthService {
	return &authService{userRepository: userRepository}
}

func (serv *authService) Signup(ctx context.Context, user *types.SignupRequest) (*dto.SignupResponse, error) {
	hashedPassword, err := serv.hashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	newUser, err := serv.userRepository.Create(ctx, &models.User{Email: string(user.Email), Password: hashedPassword})

	if err != nil {
		return nil, fmt.Errorf("failed to create user. %w", err)
	}

	return &dto.SignupResponse{Email: newUser.Email, ID: newUser.ID}, nil
}

func (serv *authService) Signin(ctx context.Context, user *types.SigninRequest) (*dto.SigninResponse, error) {
	dbUser, err := serv.userRepository.GetByEmail(ctx, string(user.Email))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app.ErrUserNotFound
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return nil, app.ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  dbUser.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}

	return &dto.SigninResponse{Token: tokenString}, nil
}

func (serv *authService) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("hashing failed: %w", err)
	}

	return string(hashedPassword), nil
}
