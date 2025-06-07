package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/i474232898/chatserver/api/types"
	"github.com/i474232898/chatserver/internal/app/common"
	"github.com/i474232898/chatserver/internal/app/dto"
	"github.com/i474232898/chatserver/internal/app/services"
)

type ChatRoomHandler struct {
	chatRoomService services.ChatRoomService
}

func NewChatRoomHandler(chatRoomService services.ChatRoomService) *ChatRoomHandler {
	return &ChatRoomHandler{chatRoomService: chatRoomService}
}

func (handler *ChatRoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var room types.CreateRoomJSONBody
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	claims, ok := common.GetClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := claims.ID

	newRoom, err := handler.chatRoomService.Create(r.Context(), &dto.NewRoomDTO{
		AdminID:   uint(userID),
		MemberIDs: room.MemberIDs,
		CreateRoomRequest: dto.CreateRoomRequest{
			Name: room.Name,
		},
	})
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Failed to create chat", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newRoom)
}

func (handler *ChatRoomHandler) DirectMessage(w http.ResponseWriter, r *http.Request) {
	var room types.CreateDirectRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	claims, ok := common.GetClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	adminID := claims.ID
	members := []int64{room.UserID}

	newRoom, err := handler.chatRoomService.Create(r.Context(), &dto.NewRoomDTO{
		AdminID:   uint(adminID),
		MemberIDs: &members,
	})
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Failed to create chat", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newRoom)
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
		http.Error(w, "Failed to list rooms", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(rooms)
}
