package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/i474232898/chatserver/api/types"
	"github.com/i474232898/chatserver/internal/app/common"
	"github.com/i474232898/chatserver/internal/app/dto"
	handlercommon "github.com/i474232898/chatserver/internal/app/handlers/common"
	"github.com/i474232898/chatserver/internal/app/services"
	"github.com/i474232898/chatserver/internal/app/validations"
)

type ChatRoomHandler struct {
	chatRoomService services.ChatRoomService
}

func NewChatRoomHandler(chatRoomService services.ChatRoomService) *ChatRoomHandler {
	return &ChatRoomHandler{chatRoomService: chatRoomService}
}

func (handler *ChatRoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var room types.CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validations.ValidateCreateRoom(&room); err != nil {
		handlercommon.HandleValidationErrors(w, err)
		return
	}
	claims, ok := common.GetClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := claims.ID

	newRoom, err := handler.chatRoomService.Create(r.Context(), &dto.CreateRoomDTO{
		AdminID:   uint(userID),
		MemberIDs: room.MemberIDs,
		Name:      room.Name,
	})
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Unable to create chat", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	handlercommon.EncodeResponse(w, newRoom)
}

func (handler *ChatRoomHandler) DirectMessage(w http.ResponseWriter, r *http.Request) {
	var room types.CreateDirectRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validations.ValidateCreateDirectRoom(&room); err != nil {
		handlercommon.HandleValidationErrors(w, err)
		return
	}
	claims, ok := common.GetClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	adminID := claims.ID
	members := []int64{room.UserID}

	newRoom, err := handler.chatRoomService.Create(r.Context(), &dto.CreateRoomDTO{
		AdminID:   uint(adminID),
		MemberIDs: members,
	})
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Unable to create chat", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	handlercommon.EncodeResponse(w, newRoom)
}

func (handler *ChatRoomHandler) ListRooms(w http.ResponseWriter, r *http.Request) {
	claims, ok := common.GetClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := claims.ID

	rooms, err := handler.chatRoomService.ListRooms(r.Context(), uint64(userID))
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Unable to list rooms", http.StatusInternalServerError)
		return
	}
	handlercommon.EncodeResponse(w, rooms)
}
