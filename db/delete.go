// Package db provides CRUD operations
package db

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// Delete removes one document with a specific id.
func Delete(id string) error {
	if _, err := motionsCollection.DeleteOne(ctx, bson.M{"id": id}); err != nil {
		return fmt.Errorf("failed to delete a document in the database: %w", err)
	}

	return nil
}
