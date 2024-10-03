package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type TelemetryRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewTelemetryRepository(uri, database, collection string) *TelemetryRepository {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	collectionRef := client.Database(database).Collection(collection)
	return &TelemetryRepository{client: client, collection: collectionRef}
}

func (m *TelemetryRepository) InsertEvent(event interface{}) error {
	_, err := m.collection.InsertOne(context.Background(), event)
	return err
}

func (m *TelemetryRepository) Close() {
	if err := m.client.Disconnect(context.Background()); err != nil {
		log.Fatalf("Failed to disconnect from MongoDB: %v", err)
	}
}
