// Package db provides CRUD operations
package db

import (
	"errors"
	"fmt"

	"github.com/chutommy/siu/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrInvalidValueType = errors.New("internal error, invalid value type")

// ReadOne finds the motion by any string and returns it.
func ReadOne(search string) (models.Motion, error) {
	// sort desc by the usage
	opts := options.Find().SetSort(bson.D{{Key: "usage", Value: -1}})

	cursor, err := motionsCollection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return models.Motion{}, fmt.Errorf("could not read from the collection: %w", err)
	}

	defer func() {
		err = cursor.Close(ctx)
	}()

	for cursor.Next(ctx) {
		var m bson.M
		if err = cursor.Decode(&m); err != nil {
			return models.Motion{}, fmt.Errorf("could not read an item from the collection: %w", err)
		}

		var idM, nameM, urlM, shortcutM string

		var ok bool

		idM, ok = m["id"].(string)
		if !ok {
			return models.Motion{}, ErrInvalidValueType
		}

		nameM, ok = m["name"].(string)
		if !ok {
			return models.Motion{}, ErrInvalidValueType
		}

		urlM, ok = m["url"].(string)
		if !ok {
			return models.Motion{}, ErrInvalidValueType
		}

		shortcutM, ok = m["shortcut"].(string)
		if !ok {
			return models.Motion{}, ErrInvalidValueType
		}

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
func ReadAll() (m []models.Motion, err error) {
	// sort desc by usage
	opts := options.Find().SetSort(bson.D{{Key: "usage", Value: -1}})

	cursor, err := motionsCollection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return []models.Motion{}, fmt.Errorf("could not read from the collection: %w", err)
	}

	defer func() {
		err = cursor.Close(ctx)
	}()

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
