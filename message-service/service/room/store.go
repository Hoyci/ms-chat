package room

import (
	"context"
	"sync"
	"time"

	"github.com/hoyci/ms-chat/message-service/db"
	"github.com/hoyci/ms-chat/message-service/types"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var (
	instance *RoomStore
	once     sync.Once
)

type RoomStore struct {
	dbRepo *db.MongoRepository
}

func GetRoomStore(dbRepo *db.MongoRepository) *RoomStore {
	once.Do(func() {
		instance = &RoomStore{dbRepo: dbRepo}
	})
	return instance
}

func (s *RoomStore) Create(ctx context.Context, newRoom types.Room) (bson.ObjectID, error) {
	result, err := db.Add(
		s.dbRepo,
		ctx,
		"rooms",
		newRoom,
	)

	if err != nil {
		return bson.NilObjectID, err
	}

	return result, nil
}

func (s *RoomStore) GetByID(ctx context.Context, roomID string) (*types.Room, error) {
	objectID, err := bson.ObjectIDFromHex(roomID)
	if err != nil {
		return nil, err
	}

	result, err := db.GetByFilter[types.Room](s.dbRepo, ctx, "rooms", bson.M{"_id": objectID})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *RoomStore) GetOrCreate(ctx context.Context, roomID string, users []int) (*types.Room, error) {
	result, err := db.GetOrCreate(s.dbRepo, ctx, "rooms", roomID, types.Room{
		ID:        bson.NewObjectID(),
		Users:     users,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
		DeletedAt: nil,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
