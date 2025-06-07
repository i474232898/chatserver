package app

type ContextKey string

const (
	JWTClaimsKey ContextKey = "jwt_claims"
	BearerPrefix string     = "Bearer "
)
