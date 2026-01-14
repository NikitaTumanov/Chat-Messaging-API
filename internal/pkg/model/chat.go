package model

import "time"

type Chat struct {
	ID        int       `gorm:"column:id;primaryKey"`
	Title     string    `gorm:"column:title"`
	CreatedAt time.Time `gorm:"column:created_at"`

	Messages []Message `gorm:"foreignKey:ChatID"`
}

func (Chat) TableName() string {
	return "chats"
}
