package model

import "time"

type Message struct {
	ID        int       `gorm:"column:id;primaryKey"`
	ChatID    int       `gorm:"column:chat_id"`
	Text      string    `gorm:"column:text"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (Message) TableName() string {
	return "messages"
}
