// Package db provides CRUD operations
package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoURL string = "mongodb+srv://ferda:ferdajekamarad@cluster0-9rknb.mongodb.net/test?retryWrites=true&w=majority"
var client *mongo.Client

var ctx context.Context
var motionsCollection *mongo.Collection

// Connect connects to the remote DB
func Connect() {
	ctx = context.Background()

	clientOpts := options.Client().ApplyURI(mongoURL)

	var err error
	client, err = mongo.Connect(ctx, clientOpts) // not redeclaring client
	if err != nil {
		log.Fatal(fmt.Errorf("Could not connect to the database: %v", err))
	}

	motionsCollection = client.Database("speed_it_up").Collection("motions")
}

// Disconnect disconnect from the remote DB
func Disconnect() {
	client.Disconnect(ctx)
}
