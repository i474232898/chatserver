
package common

import (
	"context"

	"github.com/i474232898/chatserver/internal/app/middlewares"
	"github.com/i474232898/chatserver/internal/app/services"
)

func GetClaimsFromContext(ctx context.Context) (*services.CustomClaims, bool) {
	claims, ok := ctx.Value(middlewares.JWTClaimsKey).(*services.CustomClaims)
	return claims, ok
}
