package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/i474232898/chatserver/internal/app"
	"github.com/i474232898/chatserver/internal/app/common"
)

func JWTAuthMiddleware(secretKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")

			if auth == "" || !strings.HasPrefix(auth, app.BearerPrefix) {
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(auth, app.BearerPrefix)
			claims, err := common.ParseJWT(token, secretKey)

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), app.JWTClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
