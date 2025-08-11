package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/i474232898/chatserver/api/types"
	"github.com/i474232898/chatserver/internal/app"
	"github.com/i474232898/chatserver/internal/app/handlers/common"
	"github.com/i474232898/chatserver/internal/app/services"
	"github.com/i474232898/chatserver/internal/app/validations"
)

type AuthHandler interface {
	Signup(w http.ResponseWriter, r *http.Request)
	Signin(w http.ResponseWriter, r *http.Request)
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
		common.HandleValidationErrors(w, err)
		return
	}

	ctx := r.Context()
	user, err := handler.authService.Signup(ctx, &req)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Unable to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	common.EncodeResponse(w, user)
}

func (handler authHandler) Signin(w http.ResponseWriter, r *http.Request) {
	var req types.SigninRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Incorrect body", http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	if err := validations.ValidateSignin(&req); err != nil {
		common.HandleValidationErrors(w, err)
		return
	}

	ctx := r.Context()
	user, err := handler.authService.Signin(ctx, &req)
	if err != nil {
		if errors.Is(err, app.ErrUserNotFound) || errors.Is(err, app.ErrInvalidCredentials) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		slog.Error(err.Error())
		http.Error(w, "Unable to signin", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	common.EncodeResponse(w, user)
}
