package storage

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type MongoStorage struct {
	taskCollection *mongo.Collection
	userCollection *mongo.Collection
	ctx            context.Context

	Client *mongo.Client
}

func NewMongoStorage(uri string, ctx context.Context) (ms MongoStorage) {
	log.Println("Mongo database connecting...")

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	taskCollection := client.Database("s1_tts").Collection("users")
	userCollection := client.Database("s1_tts").Collection("tasks")
	ms = MongoStorage{
		taskCollection: taskCollection,
		userCollection: userCollection,
		ctx:            ctx,
		Client:         client,
	}
	log.Println("Mongo database Connected.")
	return
}

func (ms MongoStorage) Disconnect() error {
	if err := ms.Client.Disconnect(context.TODO()); err != nil {
		log.Fatal(fmt.Sprintf(`Mongo storage disconnect error %s`, err))
		return err
	}
	return nil
}

func (ms MongoStorage) Ping() bool {
	if err := ms.Client.Ping(context.TODO(), nil); err != nil {
		log.Fatal(fmt.Sprintf(`Mongo storage ping error %s`, err))
		return false
	}
	return true
}

func (ms MongoStorage) TaskCollection() *mongo.Collection {
	return ms.taskCollection
}

func (ms MongoStorage) UserCollection() *mongo.Collection {
	return ms.userCollection
}
