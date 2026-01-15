package repository

import (
	"context"
	"errors"

	"github.com/NikitaTumanov/Chat-Messaging-API/internal/pkg/model"
	"github.com/jackc/pgx/v5/pgconn"
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
	if err != nil {
		if errors.Is(err, gorm.ErrInvalidTransaction) {
			return nil, errInternal
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return nil, errAlreadyExists
			case "23502":
				return nil, errInvalidInput
			case "23503":
				return nil, errConflict
			default:
				return nil, errInternal
			}
		}
	}

	return &message, nil
}
