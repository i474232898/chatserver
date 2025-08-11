package handlers

import (
	"net/http"

	"github.com/i474232898/chatserver/internal/app/common"
	handlercommon "github.com/i474232898/chatserver/internal/app/handlers/common"
	"github.com/i474232898/chatserver/internal/app/services"
)

type UserHandler interface {
	Me(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) UserHandler {
	return &userHandler{userService: userService}
}

func (handler *userHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := common.GetClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := handler.userService.Me(r.Context(), int64(claims.ID))
	if err != nil {
		http.Error(w, "Unable to get user info", http.StatusInternalServerError)
		return
	}

	handlercommon.EncodeResponse(w, user)
}
