package common

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/i474232898/chatserver/internal/app/services"
)

func ParseJWT(token string, secretKey []byte) (services.CustomClaims, error) {
	claims := services.CustomClaims{}
	data, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !data.Valid {
		return claims, errors.New("invalid token")
	}

	cl, ok := data.Claims.(*services.CustomClaims)
	if !ok {
		return claims, errors.New("invalid token claims")
	}

	return *cl, nil
}
