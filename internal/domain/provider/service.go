package provider

import (
	"context"
	"errors"
	"notification/internal/apperror"

	base "notification/internal/repository/mongo"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProviderConfigService struct {
	repo *ProviderConfigRepository
}

func NewProviderConfigService(repo *ProviderConfigRepository) *ProviderConfigService {
	return &ProviderConfigService{repo: repo}
}

func (s *ProviderConfigService) List(ctx context.Context, limit int, offset int) ([]*ProviderConfig, int64, error) {
	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, apperror.InternalWithDetails("Cannot count provider config", err)
	}

	configs, err := s.repo.FindAll(ctx,
		base.WithSort("created_at", false),
		base.WithLimit(limit),
		base.WithOffset(offset),
	)
	if err != nil {
		return nil, 0, apperror.InternalWithDetails("cannot get list of provider configs", err)
	}
	return configs, total, nil
}

func (s *ProviderConfigService) Create(ctx context.Context, cfg *ProviderConfig) error {
	if err := s.repo.Insert(ctx, cfg); err != nil {
		return apperror.InternalWithDetails("cannot create provider config", err)
	}
	return nil
}

func (s *ProviderConfigService) GetByID(ctx context.Context, id string) (*ProviderConfig, error) {
	cfg, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.NotFound("provider config not found")
		}
		return nil, apperror.InternalWithDetails("cannot get provider config", err)
	}
	return cfg, err
}

func (s *ProviderConfigService) Update(ctx context.Context, id string, credentials map[string]string, isActive *bool) error {
	fields := bson.M{}
	if len(credentials) > 0 {
		fields["credentials"] = credentials
	}
	if isActive != nil {
		fields["is_active"] = isActive
	}
	if len(fields) == 0 {
		return apperror.Internal("there is nothing for update")
	}

	if err := s.repo.UpdateOne(ctx, id, bson.M{"$set": fields}); err != nil {
		return apperror.InternalWithDetails("cannot update provider config", err)
	}

	return nil
}

func (s *ProviderConfigService) SetActive(ctx context.Context, id string) error {
	if err := s.repo.UpdateOne(ctx, id, bson.M{"$set": bson.M{"is_active": true}}); err != nil {
		return apperror.InternalWithDetails("cannot active provider config", err)
	}
	return nil
}

func (s *ProviderConfigService) SetInActive(ctx context.Context, id string) error {
	if err := s.repo.UpdateOne(ctx, id, bson.M{"$set": bson.M{"is_active": false}}); err != nil {
		return apperror.InternalWithDetails("cannot inactive provider config", err)
	}
	return nil
}

func (s *ProviderConfigService) GetActiveByChannel(ctx context.Context, channel ChannelType) (*ProviderConfig, error) {
	cfg, err := s.repo.GetActiveByChannel(ctx, channel)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.NotFound("no active provider found for channel: " + string(channel))
		}
		return nil, apperror.InternalWithDetails("cannot get provider config", err)
	}
	return cfg, nil
}
