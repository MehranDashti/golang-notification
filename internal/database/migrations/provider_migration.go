package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ProviderMigration struct{}

func (m ProviderMigration) Name() string { return "provider" }

func (m ProviderMigration) Run(ctx context.Context, db *mongo.Database) error {
	return ensureIndexes(ctx, db, "provider_configs", []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "channel", Value: 1}},
			Options: options.Index().SetName("idx_channel"),
		},
		{
			Keys:    bson.D{{Key: "is_active", Value: 1}},
			Options: options.Index().SetName("idx_is_active"),
		},
		{
			// the most common query: find active provider for a channel
			Keys: bson.D{
				{Key: "channel", Value: 1},
				{Key: "is_active", Value: 1},
			},
			Options: options.Index().SetName("idx_channel_is_active"),
		},
	})
}
