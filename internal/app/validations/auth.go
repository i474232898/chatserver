package validations

import (
	"github.com/i474232898/chatserver/api/types"
	"github.com/i474232898/chatserver/internal/app/dto"
)

func ValidateUser(user *types.SignupRequest) error {
	data := dto.SignupRequest{
		Email:    string(user.Email),
		Password: user.Password,
	}

	return Validator.Struct(data)
}

func ValidateSignin(user *types.SigninRequest) error {
	data := dto.SigninRequest{
		Email:    string(user.Email),
		Password: user.Password,
	}

	return Validator.Struct(data)
}
