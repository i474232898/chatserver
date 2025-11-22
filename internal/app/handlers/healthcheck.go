package handlers

import (
	"net/http"

	"github.com/i474232898/chatserver/internal/app"
	handlercommon "github.com/i474232898/chatserver/internal/app/handlers/common"
)

type HealthcheckHandler struct {
}

func NewHealthcheckHandler() *HealthcheckHandler {
	return &HealthcheckHandler{}
}

func (h *HealthcheckHandler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	if app.IsShuttingDown.Load() {
		w.WriteHeader(http.StatusServiceUnavailable)
		handlercommon.EncodeResponse(w, "Shutting down")
		return
	}
	w.WriteHeader(http.StatusOK)
	handlercommon.EncodeResponse(w, "OK")
}
