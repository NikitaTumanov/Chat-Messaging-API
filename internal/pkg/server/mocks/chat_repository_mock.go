package mocks

import (
	"context"

	"github.com/NikitaTumanov/Chat-Messaging-API/internal/pkg/model"
	"github.com/stretchr/testify/mock"
)

type ChatRepositoryMock struct {
	mock.Mock
}

func (c *ChatRepositoryMock) GetByID(ctx context.Context, id int) (*model.Chat, error) {
	args := c.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Chat), args.Error(1)
}

func (c *ChatRepositoryMock) GetByIDWithMessages(ctx context.Context, id int, limit int) (*model.Chat, error) {
	args := c.Called(ctx, id, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Chat), args.Error(1)
}

func (c *ChatRepositoryMock) Create(ctx context.Context, title string) (*model.Chat, error) {
	args := c.Called(ctx, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Chat), args.Error(1)
}

func (c *ChatRepositoryMock) Delete(ctx context.Context, id int) error {
	args := c.Called(ctx, id)
	return args.Error(0)
}
