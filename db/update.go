// Package db provides CRUD operations
package db

import (
	"fmt"

	"github.com/chutified/siu/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Update replace the motion with certain Name with a new motion.
func Update(id string, newMotion models.Motion) error {
	opts := options.Update().SetUpsert(true)
	// if names equal
	filter := bson.D{{Key: "id", Value: id}}

	// change with this bson
	update := bson.D{{Key: "$set", Value: bson.M{
		"id":       newMotion.ID,
		"name":     newMotion.Name,
		"url":      newMotion.URL,
		"shortcut": newMotion.Shortcut,
		"usage":    newMotion.Usage,
	}}}

	if _, err := motionsCollection.UpdateOne(ctx, filter, update, opts); err != nil {
		return fmt.Errorf("failed to update the motion: %w", err)
	}

	return nil
}
