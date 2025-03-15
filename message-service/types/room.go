package types

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Room struct {
	ObjectID bson.ObjectID `json:"-" bson:"_id"`
	RoomName string        `json:"room_name" bson:"room_name"`
}

func (r Room) MarshalJSON() ([]byte, error) {
	type Alias Room
	return json.Marshal(&struct {
		RoomID string `json:"room_id"`
		Alias
	}{
		RoomID: r.ObjectID.Hex(),
		Alias:  (Alias)(r),
	})
}

type CreateRoomPayload struct {
	RoomName string `bson:"room_name" json:"room_name" validate:"required,min=5"`
}

type CreateRoomResponse struct {
	ObjectID bson.ObjectID `json:"-" bson:"_id"`
	RoomID   string        `json:"room_id" bson:"-"`
	Message  string        `json:"message"`
}
