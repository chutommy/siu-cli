// Package db provides CRUD operations
package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoURL          = "mongodb+srv://ferda:ferdajekamarad@cluster0-9rknb.mongodb.net/test?retryWrites=true&w=majority"
	client            *mongo.Client
	ctx               context.Context
	motionsCollection *mongo.Collection
)

// Connect connects to the remote DB.
func Connect() {
	ctx = context.Background()

	clientOpts := options.Client().ApplyURI(mongoURL)

	client, err := mongo.Connect(ctx, clientOpts) // not re-declaring client
	if err != nil {
		log.Fatal(fmt.Errorf("could not connect to the database: %w", err))
	}

	motionsCollection = client.Database("speed_it_up").Collection("motions")
}

// Disconnect disconnects from the remote DB.
func Disconnect() error {
	return client.Disconnect(ctx)
}
