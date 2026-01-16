package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NikitaTumanov/Chat-Messaging-API/internal/pkg/model"
	"github.com/NikitaTumanov/Chat-Messaging-API/internal/pkg/server"
	"github.com/NikitaTumanov/Chat-Messaging-API/internal/pkg/server/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateChatHandler_OK(t *testing.T) {
	chatRepo := new(mocks.ChatRepositoryMock)

	handler := server.NewHandler(nil, chatRepo, nil)

	reqBody := model.CreateChatRequest{
		Title: "Test chat",
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	expectedChat := &model.Chat{
		ID:    1,
		Title: "Test chat",
	}

	chatRepo.
		On("Create", mock.Anything, "Test chat").
		Return(expectedChat, nil)

	req := httptest.NewRequest(http.MethodPost, "/chats/", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.Route(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var response model.Chat
	err = json.NewDecoder(res.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, expectedChat.ID, response.ID)
	assert.Equal(t, expectedChat.Title, response.Title)

	chatRepo.AssertExpectations(t)
}
