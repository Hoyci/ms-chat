package types

import (
	"time"

	coreTypes "github.com/hoyci/ms-chat/core/types"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusDelivered Status = "delivered"
)

type Message struct {
	ID        string     `json:"_id"  validate:"required"`
	RoomID    string     `json:"room_id"  validate:"required"`
	UserID    int        `json:"user_id" validate:"required"`
	Content   string     `json:"content" validate:"required"`
	Status    Status     `json:"status"  validate:"required"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type BroadcastMessage struct {
	UserID    int                 `json:"user_id"`
	Messages  []coreTypes.Message `json:"messages"`
	Timestamp time.Time           `json:"timestamp"`
}
