package common

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func EncodeResponse(w http.ResponseWriter, data any) {
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		slog.Error("Unable to encode response: " + err.Error())
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
	}
}
