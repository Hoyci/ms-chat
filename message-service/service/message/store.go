package message

import (
	"context"
	"sync"

	coreTypes "github.com/hoyci/ms-chat/core/types"
	"github.com/hoyci/ms-chat/message-service/db"
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

func (s *MessageStore) Create(ctx context.Context, newMessage coreTypes.Message) (bson.ObjectID, error) {
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

func (s *MessageStore) List(ctx context.Context, filter bson.M) ([]coreTypes.Message, error) {
	result, err := db.List[coreTypes.Message](
		s.dbRepo,
		ctx,
		"messages",
		filter,
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
