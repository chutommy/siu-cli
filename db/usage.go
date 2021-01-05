package db

import (
	"fmt"

	"github.com/chutified/siu/models"
	"go.mongodb.org/mongo-driver/bson"
)

// IncMotionUsage increments usage of the motion.
func IncMotionUsage(m models.Motion) error {
	// if names equal
	filter := bson.D{{Key: "id", Value: m.ID}}

	// increment usage
	update := bson.D{{
		Key: "$inc", Value: bson.M{
			"usage": 1,
		}}}

	if _, err := motionsCollection.UpdateOne(ctx, filter, update); err != nil {
		return fmt.Errorf("could not update: %w", err)
	}

	return nil
}
