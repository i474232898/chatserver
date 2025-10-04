package handlers

import (
	"net/http"

	handlercommon "github.com/i474232898/chatserver/internal/app/handlers/common"
)

type HealthcheckHandler struct {
}

func NewHealthcheckHandler() *HealthcheckHandler {
	return &HealthcheckHandler{}
}

func (h *HealthcheckHandler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	handlercommon.EncodeResponse(w, "OK")
}
