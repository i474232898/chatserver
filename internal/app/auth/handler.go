package auth

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"password123"`
}

type SignupResponse struct {
	Message string `json:"message" example:"Signup successful"`
}
type ErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var validate = validator.New(validator.WithRequiredStructEnabled())

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
	var req SignupRequest

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Incorrect body", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		var errors []ErrorResponse

		for _, v := range err.(validator.ValidationErrors) {
			errors = append(errors, ErrorResponse{Message: v.Error(), Field: v.Field()})
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(SignupResponse{
		Message: "Signup successful",
	})
}
