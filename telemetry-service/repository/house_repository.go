package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type HouseRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewHouseRepository(uri, database, collection string) *HouseRepository {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	collectionRef := client.Database(database).Collection(collection)
	return &HouseRepository{
		client: client, collection: collectionRef,
	}
}

func (m *HouseRepository) InsertHouse(event interface{}) error {
	_, err := m.collection.InsertOne(context.Background(), event)
	log.Println("Event successfully saved in MongoDB", event)
	return err
}
