package provider

import (
	"context"
	"notification/internal/apperror"
)

type ProviderConfigService struct {
	repo *ProviderConfigRepository
}

func NewProviderConfigService(repo *ProviderConfigRepository) *ProviderConfigService {
	return &ProviderConfigService{repo: repo}
}

func (s *ProviderConfigService) Create(ctx context.Context, cfg *ProviderConfig) error {
	if err := s.repo.Insert(ctx, cfg); err != nil {
		return apperror.InternalWithDetails("cannot create provider config", err)
	}
	return nil
}

// func (s *ProviderConfigService) GetActiveByChannel(ctx context.Context, channel ChannelType) (*ProviderConfig, error) {
//     cfg, err := s.repo.GetActiveByChannel(ctx, channel)
//     if err != nil {
//         if errors.Is(err, mongo.ErrNoDocuments) {
//             return nil, apperror.NotFound("no active provider found for channel: " + string(channel))
//         }
//         return nil, apperror.InternalWithDetails("cannot get provider config", err)
//     }
//     return cfg, nil
// }

// func (s *ProviderConfigService) SetActive(ctx context.Context, id string) error {
//     if err := s.repo.UpdateOne(ctx, id, bson.M{
//         "$set": bson.M{"is_active": true},
//     }); err != nil {
//         return apperror.InternalWithDetails("cannot activate provider", err)
//     }
//     return nil
// }
