package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbCtx context.Context
var urlsCollection *mongo.Collection

// connects to the database and defines the context and the collection
func initDB(dbURL string) *mongo.Client {

	client, err := mongo.NewClient(options.Client().ApplyURI(dbURL))
	if err != nil {
		log.Fatal("Could not locate database url, error: ", err)
	}
	dbCtx = context.Background()
	err = client.Connect(dbCtx)
	if err != nil {
		log.Fatal("Could not connect to the database, error:", err)
	}

	urlsCollection = client.Database("urlshortener").Collection("url_list")
	return client
}

// Crud
func createURL(u url) error {

	newURL := bson.M{
		"id":     u.ID,
		"origin": u.Origin,
		"short":  u.Short,
		"usage":  u.Usage,
	}

	if _, err := urlsCollection.InsertOne(dbCtx, newURL); err != nil {
		return fmt.Errorf("Could not insert new item into database, error: %v", err)
	}
	return nil
}

// cRud
func readOneURL(sh string) (url, error) {

	opts := options.FindOne().SetSort(bson.D{{Key: "usage", Value: 1}})

	var u bson.M
	if err := urlsCollection.FindOne(dbCtx, bson.D{{Key: "short", Value: sh}}, opts).Decode(&u); err != nil {
		return url{}, fmt.Errorf("Failed to find item with short %v, error: %v", sh, err)
	}

	return url{
		ID:     u["id"].(string),
		Origin: u["origin"].(string),
		Short:  u["short"].(string),
		Usage:  u["usage"].(int64),
	}, nil
}

// cRud
func readAllURLs() ([]url, error) {
	var urls []url

	opts := options.Find().SetSort(bson.D{{Key: "usage", Value: 1}})

	urlsCursor, err := urlsCollection.Find(dbCtx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("Could not read from collection, error: %v", err)
	}

	// stores values in urlsMap
	var urlsMap []bson.M
	if err := urlsCursor.All(dbCtx, &urlsMap); err != nil {
		return nil, fmt.Errorf("Invalid value in database, error: %v", err)
	}

	// stores results in urls
	for _, urlItem := range urlsMap {

		u := url{
			ID:     urlItem["id"].(string),
			Origin: urlItem["origin"].(string),
			Short:  urlItem["short"].(string),
			Usage:  urlItem["usage"].(int64),
		}

		urls = append(urls, u)
	}
	return urls, nil
}

// crUd
func updateURL(id string, newURL url) error {

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{Key: "id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"id":     newURL.ID,
		"origin": newURL.Origin,
		"short":  newURL.Short,
		"usage":  newURL.Usage,
	}}}

	if _, err := urlsCollection.UpdateOne(dbCtx, filter, update, opts); err != nil {
		return fmt.Errorf("Update failed, error: %v", err)
	}
	return nil
}

// cruD
func deleteURL(id string) error {

	_, err := urlsCollection.DeleteOne(dbCtx, bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("Failed to remove item with id %v, error: %v", id, err)
	}
	return nil
}
