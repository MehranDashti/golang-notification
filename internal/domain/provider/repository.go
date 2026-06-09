package provider

import (
	"context"
	base "notification/internal/repository/mongo"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const CollectionName = "provider_configs"

type ProviderConfigRepository struct {
	*base.BaseRepository[*ProviderConfig]
}

func NewProviderConfigRepository(db *mongo.Database) *ProviderConfigRepository {
	return &ProviderConfigRepository{
		BaseRepository: base.NewBaseRepository[*ProviderConfig](db, CollectionName),
	}
}

func (r *ProviderConfigRepository) GetActiveByChannel(ctx context.Context, channel ChannelType) (*ProviderConfig, error) {
	return r.FindOne(ctx,
		base.WithFilter(bson.M{
			"channel":   channel,
			"is_active": true,
		}),
	)
}
