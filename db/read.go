// Package db provides CRUD operations
package db

import (
	"fmt"

	"github.com/chutified/siu/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ReadOne finds the motion by any string and returns it.
func ReadOne(search string) (models.Motion, error) {
	// sort desc by the usage
	opts := options.Find().SetSort(bson.D{{Key: "usage", Value: -1}})

	cursor, err := motionsCollection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return models.Motion{}, fmt.Errorf("could not read from the collection: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var m bson.M
		if err = cursor.Decode(&m); err != nil {
			return models.Motion{}, fmt.Errorf("could not read an item from the collection: %w", err)
		}

		idM := m["id"].(string)
		nameM := m["name"].(string)
		urlM := m["url"].(string)
		shortcutM := m["shortcut"].(string)

		// checks each field
		if idM == search || nameM == search || urlM == search || shortcutM == search {
			return models.Motion{
				ID:       idM,
				Name:     nameM,
				URL:      urlM,
				Shortcut: shortcutM,
				Usage:    m["usage"].(int32),
			}, nil
		}
	}

	return models.Motion{}, mongo.ErrNoDocuments
}

// ReadAll return all motions in the collection.
func ReadAll() ([]models.Motion, error) {
	// sort desc by usage
	opts := options.Find().SetSort(bson.D{{Key: "usage", Value: -1}})

	cursor, err := motionsCollection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return []models.Motion{}, fmt.Errorf("could not read from the collection: %w", err)
	}
	defer cursor.Close(ctx)

	// create a var to store
	var motions []models.Motion

	for cursor.Next(ctx) {
		var m bson.M
		if err = cursor.Decode(&m); err != nil {
			return []models.Motion{}, fmt.Errorf("could not read an item from the collection: %w", err)
		}

		// get each motion
		motions = append(motions, models.Motion{
			ID:       m["id"].(string),
			Name:     m["name"].(string),
			URL:      m["url"].(string),
			Shortcut: m["shortcut"].(string),
			Usage:    m["usage"].(int32),
		})
	}

	return motions, nil
}
