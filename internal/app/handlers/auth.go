package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/i474232898/chatserver/api/types"
	"github.com/i474232898/chatserver/internal/app/services"
	"github.com/i474232898/chatserver/internal/app/validations"
)

type ErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type AuthHandler interface {
	Signup(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) AuthHandler {
	return authHandler{authService: authService}
}

func (handler authHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req types.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Incorrect body", http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	if err := validations.ValidateUser(&req); err != nil {
		var errors []ErrorResponse

		for _, v := range err.(validator.ValidationErrors) {
			errors = append(errors, ErrorResponse{Message: v.Error(), Field: v.Field()})
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)
		return
	}

	ctx := context.Background()
	user, err := handler.authService.Signup(ctx, &req)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
