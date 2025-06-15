package common

import (
	"encoding/json"
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
	json.NewEncoder(w).Encode(errors)
}
