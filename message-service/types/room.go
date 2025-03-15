package types

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Room struct {
	ID        bson.ObjectID `json:"_id" bson:"_id"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time    `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time    `json:"deleted_at" bson:"deleted_at"`
}

func (r Room) MarshalJSON() ([]byte, error) {
	type Alias Room
	return json.Marshal(&struct {
		ID string `json:"_id"`
		Alias
	}{
		ID:    r.ID.Hex(),
		Alias: (Alias)(r),
	})
}

// type CreateRoomPayload struct {
// 	RoomName string `bson:"room_name" json:"room_name" validate:"required,min=5"`
// }

// type CreateRoomResponse struct {
// 	ObjectID bson.ObjectID `json:"-" bson:"_id"`
// 	RoomID   string        `json:"room_id" bson:"-"`
// 	Message  string        `json:"message"`
// }
