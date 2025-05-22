package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/i474232898/chatserver/internal/app/services"
)

type ContextKey string
const (
	JWTClaimsKey ContextKey = "jwt_claims"
)
var bearerPrefix = "Bearer "

func JWTAuthMiddleware(secretKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")

			if auth == "" || !strings.HasPrefix(auth, bearerPrefix) {
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(auth, bearerPrefix)

			data, err := jwt.ParseWithClaims(token, &services.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				return secretKey, nil
			})

			if err != nil || !data.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := data.Claims.(*services.CustomClaims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), JWTClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}


