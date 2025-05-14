package auth

import (
	"net/http"
)

type SignupRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

type SignupResponse struct {
	Message string `json:"message" example:"User created successfully"`
}

// @Summary      Sign up a new user
// @Description  Create a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body SignupRequest true "Signup credentials"
// @Success      201 {object} SignupResponse
// @Failure      400 {object} map[string]string
// @Router       /auth/signup [post]
func SignupH(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Signup endpoint"))
}
