package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type EmergencyRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewEmergencyRepository(uri, database, collection string) *EmergencyRepository {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	collectionRef := client.Database(database).Collection(collection)
	return &EmergencyRepository{
		client: client, collection: collectionRef,
	}
}

func (m *EmergencyRepository) SaveEmergencyEvent(event interface{}) error {
	_, err := m.collection.InsertOne(context.Background(), event)
	log.Println("Event successfully saved in MongoDB", event)
	return err
}
