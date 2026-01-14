package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/NikitaTumanov/Chat-Messaging-API/internal/pkg/database/repository"
)

type HTTPHandler interface {
	Route(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	chats    *repository.ChatRepository
	messages *repository.MessageRepository
}

func NewHandler(chatRepo *repository.ChatRepository, msgRepo *repository.MessageRepository) *Handler {
	return &Handler{
		chats:    chatRepo,
		messages: msgRepo,
	}
}

func (s *Handler) Route(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/chats/")
	path = strings.Trim(path, "/")

	parts := strings.Split(path, "/")

	if len(parts) == 1 && parts[0] == "" {
		switch r.Method {
		case http.MethodPost:
			s.CreateChatHandler(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "invalid chat id", http.StatusBadRequest)
		return
	}

	if len(parts) == 1 {
		switch r.Method {
		case http.MethodGet:
			s.GetChatHandler(w, r, id)
		case http.MethodDelete:
			s.DeleteChatHandler(w, r, id)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	if runes := []rune(r.URL.Path); len(parts) == 2 && parts[1] == "messages" && runes[len(runes)-1] == '/' {
		switch r.Method {
		case http.MethodPost:
			s.SendMessageHandler(w, r, id)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	http.NotFound(w, r)
}

func (s *Handler) CreateChatHandler(w http.ResponseWriter, r *http.Request) {
	chat, err := s.chats.Create(r.Context(), "qwerty")
	if err != nil {
		w.Write([]byte("ascasc"))
	}
	json.NewEncoder(w).Encode(chat)
}

func (s *Handler) SendMessageHandler(w http.ResponseWriter, r *http.Request, id int) {
	message, err := s.messages.Create(r.Context(), 1, "sdfsdfsdf")
	if err != nil {

	}
	json.NewEncoder(w).Encode(message)
}

func (s *Handler) GetChatHandler(w http.ResponseWriter, r *http.Request, id int) {
	chat, err := s.chats.GetByID(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(chat)
}

func (s *Handler) DeleteChatHandler(w http.ResponseWriter, r *http.Request, id int) {
	err := s.chats.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
}
