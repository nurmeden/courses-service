package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupDatabase() *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: %v", err)
		return nil
	}
	err = client.Ping(context.Background(), nil)
	log.Println(err)
	if err != nil {
		return nil
	}
	fmt.Printf("client: %v\n", client)
	return client
}
