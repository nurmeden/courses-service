package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupDatabase(ctx context.Context, mongoURI string) (*mongo.Client, error) {
	co := options.Client().ApplyURI("mongodb://coursesdb:27017")
	client, err := mongo.Connect(ctx, co)
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Printf("Failed to ping MongoDB: %v", err)
		return nil, err
	}
	log.Printf("Connected to MongoDB: %v", client)
	return client, nil
}
