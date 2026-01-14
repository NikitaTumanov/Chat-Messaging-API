package repository

import (
	"context"

	"github.com/NikitaTumanov/Chat-Messaging-API/internal/pkg/model"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

func (m *MessageRepository) Create(ctx context.Context, chatID int, text string) (*model.Message, error) {
	message := model.Message{
		ChatID: chatID,
		Text:   text,
	}
	err := m.db.WithContext(ctx).Create(&message).Error
	return &message, err
}
