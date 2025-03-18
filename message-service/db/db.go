package db

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/hoyci/ms-chat/message-service/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoRepository struct {
	client *mongo.Client
	config config.Config
}

func NewMongoRepository(cfg config.Config) *MongoRepository {
	client, err := mongo.Connect(options.Client().ApplyURI(cfg.DatabaseURL))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	return &MongoRepository{
		client: client,
		config: cfg,
	}
}

func Add[T any](repo *MongoRepository, ctx context.Context, collectionName string, document T) (bson.ObjectID, error) {
	collection := repo.client.Database(repo.config.DatabaseName).Collection(collectionName)
	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		return bson.NilObjectID, err
	}
	oid, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return bson.NilObjectID, errors.New("failed to convert the entered ID")
	}
	return oid, nil
}

func List[T any](repo *MongoRepository, ctx context.Context, collectionName string, filter bson.M) ([]T, error) {
	collection := repo.client.Database(repo.config.DatabaseName).Collection(collectionName)
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	for cursor.Next(ctx) {
		var elem T
		if err := cursor.Decode(&elem); err != nil {
			return nil, err
		}
		results = append(results, elem)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func GetByFilter[T any](repo *MongoRepository, ctx context.Context, collectionName string, filter bson.M) (*T, error) {
	collection := repo.client.Database(repo.config.DatabaseName).Collection(collectionName)

	var result T
	if err := collection.FindOne(ctx, filter).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func GetOrCreate[T any](repo *MongoRepository, ctx context.Context, collectionName string, filter bson.M, document T) (*T, error) {
	existing, err := GetByFilter[T](repo, ctx, collectionName, filter)
	if err == nil {
		return existing, nil
	}
	if err != mongo.ErrNoDocuments && err != bson.ErrInvalidHex {
		return nil, err
	}

	_, err = Add(repo, ctx, collectionName, document)
	if err == nil {
		return &document, nil
	}
	if mongo.IsDuplicateKeyError(err) {
		existing, err := GetByFilter[T](repo, ctx, collectionName, filter)
		if err == nil {
			return existing, nil
		}
		return nil, err
	}
	return nil, err
}
