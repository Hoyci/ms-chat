package message

import (
	"context"
	"sync"

	"github.com/hoyci/ms-chat/message-service/db"
	"github.com/hoyci/ms-chat/message-service/types"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var (
	instance *MessageStore
	once     sync.Once
)

type MessageStore struct {
	dbRepo *db.MongoRepository
}

func GetMessageStore(dbRepo *db.MongoRepository) *MessageStore {
	once.Do(func() {
		instance = &MessageStore{dbRepo: dbRepo}
	})
	return instance
}

func (s *MessageStore) Create(ctx context.Context, newMessage types.Message) (bson.ObjectID, error) {
	result, err := db.Add(
		s.dbRepo,
		ctx,
		"messages",
		newMessage,
	)

	if err != nil {
		return bson.NilObjectID, err
	}

	return result, nil
}

func (s *MessageStore) List(ctx context.Context, roomId string) ([]types.Message, error) {
	result, err := db.List[types.Message](
		s.dbRepo,
		ctx,
		"messages",
		bson.M{"room_id": roomId},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
