package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/NikitaTumanov/Chat-Messaging-API/internal/pkg/database/repository"
	"github.com/NikitaTumanov/Chat-Messaging-API/internal/pkg/model"
)

const (
	titleLimit   = 200
	textLimit    = 5000
	defaultLimit = 20
	maxLimit     = 100
)

var (
	errMethodNotAllowed = errors.New("method not allowed")
	errInvalidChatID    = errors.New("invalid chat id")
	errInvalidJSON      = errors.New("invalid json")
)

type HTTPHandler interface {
	Route(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	log      *slog.Logger
	chats    repository.ChatRepo
	messages repository.MessageRepo
}

func NewHandler(log *slog.Logger, chatRepo repository.ChatRepo, msgRepo repository.MessageRepo) *Handler {
	return &Handler{
		log:      log,
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
			http.Error(w, errMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		}
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil || id < 0 {
		http.Error(w, errInvalidChatID.Error(), http.StatusBadRequest)
		return
	}

	if len(parts) == 1 {
		switch r.Method {
		case http.MethodGet:
			s.GetChatHandler(w, r, id)
		case http.MethodDelete:
			s.DeleteChatHandler(w, r, id)
		default:
			http.Error(w, errMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		}
		return
	}

	if runes := []rune(r.URL.Path); len(parts) == 2 && parts[1] == "messages" && runes[len(runes)-1] == '/' {
		switch r.Method {
		case http.MethodPost:
			s.SendMessageHandler(w, r, id)
		default:
			http.Error(w, errMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		}
		return
	}

	http.NotFound(w, r)
}

func writeDBError(w http.ResponseWriter, err error) {
	switch err {
	case repository.ErrNotFound:
		http.Error(w, "not found", http.StatusNotFound)

	case repository.ErrAlreadyExists:
		http.Error(w, "already exists", http.StatusConflict)

	case repository.ErrConflict:
		http.Error(w, "conflict", http.StatusConflict)

	case repository.ErrInvalidInput:
		http.Error(w, "invalid input", http.StatusBadRequest)

	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (s *Handler) CreateChatHandler(w http.ResponseWriter, r *http.Request) {
	var request model.CreateChatRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, errInvalidJSON.Error(), http.StatusBadRequest)
		return
	}

	if decoder.More() {
		http.Error(w, errInvalidJSON.Error(), http.StatusBadRequest)
		return
	}

	request.Title = strings.TrimSpace(request.Title)
	if request.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(request.Title) > titleLimit {
		http.Error(w, fmt.Sprintf("title is longer than %d characters", titleLimit), http.StatusBadRequest)
		return
	}

	chat, err := s.chats.Create(r.Context(), request.Title)
	if err != nil {
		s.log.Error("failed to create chat",
			"error", err,
			"title", request.Title,
		)
		writeDBError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chat)
}

func (s *Handler) SendMessageHandler(w http.ResponseWriter, r *http.Request, id int) {
	var request model.SendMessageRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, errInvalidJSON.Error(), http.StatusBadRequest)
		return
	}

	if decoder.More() {
		http.Error(w, errInvalidJSON.Error(), http.StatusBadRequest)
		return
	}

	request.Text = strings.TrimSpace(request.Text)
	if request.Text == "" {
		http.Error(w, "text is required", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(request.Text) > textLimit {
		http.Error(w, fmt.Sprintf("text is longer than %d characters", textLimit), http.StatusBadRequest)
		return
	}

	chat, err := s.chats.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			s.log.Error("failed to get chat by ID",
				"error", err,
				"text", request.Text,
			)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	message, err := s.messages.Create(r.Context(), chat.ID, request.Text)
	if err != nil {
		s.log.Error("failed to add message to chat",
			"error", err,
			"text", request.Text,
		)
		writeDBError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func (s *Handler) GetChatHandler(w http.ResponseWriter, r *http.Request, id int) {
	var request model.GetChatRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, errInvalidJSON.Error(), http.StatusBadRequest)
		return
	}

	if decoder.More() {
		http.Error(w, errInvalidJSON.Error(), http.StatusBadRequest)
		return
	}

	if request.Limit == 0 {
		request.Limit = defaultLimit
	}
	if request.Limit > maxLimit {
		http.Error(w, fmt.Sprintf("limit is greater than %d", maxLimit), http.StatusBadRequest)
		return
	}
	if request.Limit < 0 {
		http.Error(w, fmt.Sprintf("limit should be between %d and %d", 1, maxLimit), http.StatusBadRequest)
		return
	}

	chat, err := s.chats.GetByIDWithMessages(r.Context(), id, request.Limit)
	if err != nil {
		s.log.Error("failed to get chat by ID with messages",
			"error", err,
			"limit", request.Limit,
		)
		writeDBError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chat)
}

func (s *Handler) DeleteChatHandler(w http.ResponseWriter, r *http.Request, id int) {
	err := s.chats.Delete(r.Context(), id)
	if err != nil {
		s.log.Error("failed to delete chat by ID with messages",
			"error", err,
		)
		writeDBError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
