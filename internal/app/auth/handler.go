package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"log/slog"

	"github.com/i474232898/chatserver/pkg/db"
)

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"password123"`
}

type SignupResponse struct {
	UserId int `json:"userId" example:"1"`
}
type ErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var validate = validator.New(validator.WithRequiredStructEnabled())

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

	pool := db.GetPool()

	password, err := HashPassword(req.Password)
	if err != nil {
		slog.Error("Password hashing failed", "error", err.Error())
	}
	var newUserID int
	insertQuery := `insert into "users" ("email", "password") values ($1, $2) returning user_id`

	err = pool.QueryRow(context.Background(), insertQuery, req.Email, password).Scan(&newUserID)
	if err != nil {
		//todo: check if duplicate
		fmt.Println(">>11>", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(SignupResponse{
		UserId: newUserID,
	})
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		slog.Error("Hashing failed", "error", err.Error())
		return "", err
	}

	return string(hashedPassword), nil
}
