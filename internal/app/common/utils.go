package common

import (
	"context"

	"github.com/i474232898/chatserver/internal/app"
	"github.com/i474232898/chatserver/internal/app/services"
)

func GetClaimsFromContext(ctx context.Context) (*services.CustomClaims, bool) {
	claims, ok := ctx.Value(app.JWTClaimsKey).(*services.CustomClaims)
	return claims, ok
}
