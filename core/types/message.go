package types

import "time"

type Status string

const (
	StatusPending   Status = "pending"
	StatusDelivered Status = "delivered"
)

type Message struct {
	ID         string     `json:"_id" bson:"_id"`
	RoomID     string     `json:"room_id" bson:"room_id" validate:"required"`
	SenderID   int        `json:"sender_id" bson:"sender_id" validate:"required"`
	ReceiverID int        `json:"receiver_id" bson:"receiver_id" validate:"required"`
	Content    string     `json:"content" bson:"content" validate:"required"`
	Status     Status     `json:"status" bson:"status"`
	CreatedAt  time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at" bson:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at" bson:"deleted_at"`
	ClientID   string     `json:"-" bson:"-"`
}
