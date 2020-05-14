package handlers

import (
	"context"
	"fmt"
	"log"
	"speedit/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DbCtx context for crud operation
var DbCtx context.Context

// UrlsCollection holds all urls
var UrlsCollection *mongo.Collection

// InitDB connects to the database and defines the context and the collection
func InitDB(ctx context.Context, dbURL string) *mongo.Client {
	DbCtx = ctx

	client, err := mongo.NewClient(options.Client().ApplyURI(dbURL))
	if err != nil {
		log.Fatal("Could not locate database url, error: ", err)
	}
	DbCtx = context.Background()
	err = client.Connect(DbCtx)
	if err != nil {
		log.Fatal("Could not connect to the database, error:", err)
	}

	UrlsCollection = client.Database("urlshortener").Collection("url_list")
	return client
}

// CreateURL create a new document in the UrlsCollection
func CreateURL(u models.Url) error {

	newURL := bson.M{
		"id":     u.ID,
		"origin": u.Origin,
		"short":  u.Short,
		"usage":  u.Usage,
	}

	if _, err := UrlsCollection.InsertOne(DbCtx, newURL); err != nil {
		return fmt.Errorf("Could not insert new item into database, error: %v", err)
	}
	return nil
}

// ReadOneURL returns url if is found matching short-long
func ReadOneURL(sh string) (models.Url, error) {

	opts := options.FindOne().SetSort(bson.D{{Key: "usage", Value: 1}})

	var u bson.M
	if err := UrlsCollection.FindOne(
		DbCtx,
		bson.D{{Key: "short", Value: sh}},
		opts,
	).Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Url{}, mongo.ErrNoDocuments
		}
		return models.Url{}, fmt.Errorf("Failed to find item with short %v, error: %v", sh, err)
	}

	if _, err := UrlsCollection.UpdateOne(
		DbCtx,
		bson.D{{Key: "short", Value: sh}},
		bson.D{{Key: "$inc", Value: bson.D{{Key: "usage", Value: 1}}}},
	); err != nil {
		return models.Url{}, fmt.Errorf("Could not increment status usage of url %v, error: %v", u["origin"], err)
	}

	return models.Url{
		ID:     u["id"].(string),
		Origin: u["origin"].(string),
		Short:  u["short"].(string),
		Usage:  u["usage"].(int64),
	}, nil
}

// ReadAllURLs returns slice of urls in UrlsCollection
func ReadAllURLs() ([]models.Url, error) {
	var urls []models.Url

	opts := options.Find().SetSort(bson.D{{Key: "usage", Value: 1}})

	urlsCursor, err := UrlsCollection.Find(DbCtx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("Could not read from collection, error: %v", err)
	}

	// stores values in urlsMap
	var urlsMap []bson.M
	if err := urlsCursor.All(DbCtx, &urlsMap); err != nil {
		return nil, fmt.Errorf("Invalid value in database, error: %v", err)
	}

	// stores results in urls
	for _, urlItem := range urlsMap {

		u := models.Url{
			ID:     urlItem["id"].(string),
			Origin: urlItem["origin"].(string),
			Short:  urlItem["short"].(string),
			Usage:  urlItem["usage"].(int64),
		}

		urls = append(urls, u)
	}
	return urls, nil
}

// UpdateURL finds one url and change it for the newURL
func UpdateURL(id string, newURL models.Url) error {

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{Key: "id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"id":     newURL.ID,
		"origin": newURL.Origin,
		"short":  newURL.Short,
		"usage":  newURL.Usage,
	}}}

	if _, err := UrlsCollection.UpdateOne(DbCtx, filter, update, opts); err != nil {
		return fmt.Errorf("Update failed, error: %v", err)
	}
	return nil
}

// DeleteURL finds the url with ID and remove it from the UrlsCollection
func DeleteURL(id string) error {

	_, err := UrlsCollection.DeleteOne(DbCtx, bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("Failed to remove item with id %v, error: %v", id, err)
	}
	return nil
}
