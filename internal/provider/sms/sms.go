package sms

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
	case providerDomain.ProviderKavenegar:
		return newKavenegar(KavenegarConfig{
			APIKey: cfg.Credentials["api_key"],
			Sender: cfg.Credentials["sender"],
		}), nil
	case providerDomain.ProviderTwilio:
		return NewTwilio(TwilioConfig{
			AccountSID: cfg.Credentials["account_sid"],
			AuthToken:  cfg.Credentials["auth_token"],
			From:       cfg.Credentials["from"],
		}), nil
	default:
		return nil, fmt.Errorf("sms: unknown provider %q", cfg.Provider)
	}
}

type Dispatcher struct {
	providerService *providerDomain.ProviderConfigService
}

func NewDispatcher(providerService *providerDomain.ProviderConfigService) *Dispatcher {
	return &Dispatcher{
		providerService: providerService,
	}
}

func (d *Dispatcher) Send(ctx context.Context, msg provider.Message) (*provider.Result, error) {
	cfg, err := d.providerService.GetActiveByChannel(ctx, providerDomain.ChannelSms)
	if err != nil {
		return nil, err
	}
	p, err := NewFromConfig(cfg)
	if err != nil {
		return nil, err
	}
	return p.Send(ctx, msg)
}
