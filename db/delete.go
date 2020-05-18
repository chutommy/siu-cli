// Package db provides CRUD operations
package db

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// Delete removes one document with specific id
func Delete(id string) error {
	connect()
	defer disconnect()

	if _, err := motionsCollection.DeleteOne(ctx, bson.M{"id": id}); err != nil {
		return fmt.Errorf("Failed to delete a document in the database: %v", err)
	}
	return nil
}
