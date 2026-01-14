package repository

import (
	"context"

	"github.com/NikitaTumanov/Chat-Messaging-API/internal/pkg/model"
	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{
		db: db,
	}
}

func (c *ChatRepository) GetByID(ctx context.Context, id int) (*model.Chat, error) {
	var chat model.Chat

	err := c.db.WithContext(ctx).Preload("Messages").First(&chat, id).Error
	return &chat, err
}

func (c *ChatRepository) Create(ctx context.Context, title string) (*model.Chat, error) {
	chat := model.Chat{
		Title: title,
	}

	err := c.db.WithContext(ctx).Create(&chat).Error
	return &chat, err
}

func (c *ChatRepository) Delete(ctx context.Context, id int) error {
	return c.db.WithContext(ctx).Delete(&model.Chat{}, id).Error
}
