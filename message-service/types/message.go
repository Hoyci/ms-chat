package types

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusDelivered Status = "delivered"
)

type Message struct {
	ID        bson.ObjectID `json:"_id" bson:"_id"`
	RoomID    bson.ObjectID `json:"room_id"  bson:"room_id" validate:"required"`
	UserID    int           `json:"user_id" bson:"user_id" validate:"required"`
	Content   string        `json:"content" bson:"content" validate:"required"`
	Status    Status        `json:"status" bson:"status" validate:"required"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time    `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time    `json:"deleted_at" bson:"deleted_at"`
}

func (r Message) MarshalJSON() ([]byte, error) {
	type Alias Message
	return json.Marshal(&struct {
		ID     string `json:"_id"`
		RoomID string `json:"room_id"`
		Alias
	}{
		ID:     r.ID.Hex(),
		RoomID: r.RoomID.Hex(),
		Alias:  (Alias)(r),
	})
}
