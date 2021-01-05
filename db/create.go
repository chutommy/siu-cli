// Package db provides CRUD operations
package db

import (
	"fmt"

	"github.com/chutified/siu/models"
	"go.mongodb.org/mongo-driver/bson"
)

// Create creates a new motion and returns an error if any is occurred.
func Create(m models.Motion) error {
	// new motion
	motionBSON := bson.D{
		{Key: "id", Value: m.ID},
		{Key: "name", Value: m.Name},
		{Key: "url", Value: m.URL},
		{Key: "shortcut", Value: m.Shortcut},
		{Key: "usage", Value: m.Usage},
	}

	if _, err := motionsCollection.InsertOne(ctx, motionBSON); err != nil {
		return fmt.Errorf("could not insert new item \n\n\t%v\n\t into database collection: %w", m, err)
	}

	return nil
}
