package repository

import (
	"context"
	"errors"

	"github.com/NikitaTumanov/Chat-Messaging-API/internal/pkg/model"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrNotFound      = errors.New("entity not found")
	ErrAlreadyExists = errors.New("entity already exists")
	ErrConflict      = errors.New("conflict")
	ErrInvalidInput  = errors.New("invalid input")
	ErrInternal      = errors.New("internal error")
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
			return nil, ErrNotFound
		}
		return nil, ErrInternal
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
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}

	return &chat, nil
}

func validateSpecificDBError(err error) error {
	if errors.Is(err, gorm.ErrInvalidTransaction) {
		return ErrInternal
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return ErrAlreadyExists
		case "23502":
			return ErrInvalidInput
		case "23503":
			return ErrConflict
		default:
			return ErrInternal
		}
	}

	return nil
}

func (c *ChatRepository) Create(ctx context.Context, title string) (*model.Chat, error) {
	chat := model.Chat{
		Title: title,
	}

	err := c.db.WithContext(ctx).Create(&chat).Error
	if err != nil {
		return nil, validateSpecificDBError(err)
	}

	return &chat, nil
}

func (c *ChatRepository) Delete(ctx context.Context, id int) error {
	result := c.db.WithContext(ctx).Delete(&model.Chat{}, id)
	if result.Error != nil {
		return ErrInternal
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
