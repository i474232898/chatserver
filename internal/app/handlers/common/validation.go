package common

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/i474232898/chatserver/api/types"
)

func HandleValidationErrors(w http.ResponseWriter, err error) {
	var errors types.BadRequestError
	for _, v := range err.(validator.ValidationErrors) {
		errors = append(errors, types.ValidationError{Message: v.Error(), Field: v.Field()})
	}
	w.WriteHeader(http.StatusBadRequest)
	err = json.NewEncoder(w).Encode(errors)
	if err != nil {
		slog.Error("Failed to encode errors: " + err.Error())
		http.Error(w, "Failed to encode errors", http.StatusInternalServerError)
		return
	}
}
