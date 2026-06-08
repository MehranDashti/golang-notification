package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Migration interface {
	Name() string
	Run(ctx context.Context, db *mongo.Database) error
}

func ensureIndexes(ctx context.Context, db *mongo.Database, collection string, indexes []mongo.IndexModel) error {
	_, err := db.Collection(collection).Indexes().CreateMany(ctx, indexes)
	return err
}
