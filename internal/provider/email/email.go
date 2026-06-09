package email

import (
	"context"
	"fmt"

	providerDomain "notification/internal/domain/provider"
	"notification/internal/provider"
)

type Provider interface {
	provider.Provider
}

func NewFromConfig(cfg *providerDomain.ProviderConfig) (Provider, error) {
	switch cfg.Provider {
	case providerDomain.ProviderSMTP:
		return NewSMTP(SMTPConfig{
			Host:     cfg.Credentials["host"],
			Port:     cfg.Credentials["port"],
			Username: cfg.Credentials["username"],
			Password: cfg.Credentials["password"],
			From:     cfg.Credentials["from"],
		}), nil
	default:
		return nil, fmt.Errorf("email: unknown provider %q", cfg.Provider)
	}
}

type Dispatcher struct {
	providerSvc *providerDomain.ProviderConfigService
}

func NewDispatcher(providerSvc *providerDomain.ProviderConfigService) *Dispatcher {
	return &Dispatcher{providerSvc: providerSvc}
}

func (d *Dispatcher) Send(ctx context.Context, msg provider.Message) (*provider.Result, error) {
	cfg, err := d.providerSvc.GetActiveByChannel(ctx, providerDomain.ChannelEmail)
	if err != nil {
		return nil, err
	}

	p, err := NewFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	return p.Send(ctx, msg)
}
