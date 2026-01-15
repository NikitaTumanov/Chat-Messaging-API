package repository

import (
	"context"
	"errors"

	"github.com/NikitaTumanov/Chat-Messaging-API/internal/pkg/model"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	errNotFound      = errors.New("entity not found")
	errAlreadyExists = errors.New("entity already exists")
	errConflict      = errors.New("conflict")
	errInvalidInput  = errors.New("invalid input")
	errInternal      = errors.New("internal error")
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

	err := c.db.WithContext(ctx).First(&chat, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errNotFound
		}
		return nil, errInternal
	}

	return &chat, nil
}

func (c *ChatRepository) GetByIDWithMessages(ctx context.Context, id int, limit int) (*model.Chat, error) {
	var chat model.Chat

	err := c.db.WithContext(ctx).Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC").Limit(limit)
	}).First(&chat, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errNotFound
		}
		return nil, errInternal
	}

	return &chat, nil
}

func (c *ChatRepository) Create(ctx context.Context, title string) (*model.Chat, error) {
	chat := model.Chat{
		Title: title,
	}

	err := c.db.WithContext(ctx).Create(&chat).Error
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

	return &chat, nil
}

func (c *ChatRepository) Delete(ctx context.Context, id int) error {
	result := c.db.WithContext(ctx).Delete(&model.Chat{}, id)
	if result.Error != nil {
		return errInternal
	}
	if result.RowsAffected == 0 {
		return errNotFound
	}

	return nil
}
