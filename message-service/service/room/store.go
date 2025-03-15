package room

import (
	"context"

	"github.com/hoyci/ms-chat/message-service/db"
	"github.com/hoyci/ms-chat/message-service/types"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RoomStore struct {
	dbRepo *db.MongoRepository
}

func NewRoomStore(dbRepo *db.MongoRepository) *RoomStore {
	return &RoomStore{
		dbRepo: dbRepo,
	}
}

func (s *RoomStore) Create(ctx context.Context, newRoom types.CreateRoomPayload) (bson.ObjectID, error) {
	result, err := db.Add(s.dbRepo, ctx, "rooms", newRoom)

	if err != nil {
		return bson.NilObjectID, err
	}

	return result, nil
}

func (s *RoomStore) GetByID(ctx context.Context, roomID string) (*types.Room, error) {
	result, err := db.GetByID[types.Room](s.dbRepo, ctx, "rooms", roomID)

	if err != nil {
		return nil, err
	}

	return result, nil
}
