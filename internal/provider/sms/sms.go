package sms

// import (
//     "context"
//     "fmt"

//     providerDomain "notification/internal/domain/provider"
//     "notification/internal/provider"
// )

// type Provider interface {
//     provider.Provider
// }

// func NewFromConfig(cfg *providerDomain.ProviderConfig) (Provider, error) {
//     switch cfg.Provider {
//     case providerDomain.ProviderKavenegar:
//         return NewKavenegar(KavenegarConfig{
//             APIKey: cfg.Credentials["api_key"],
//             Sender: cfg.Credentials["sender"],
//         }), nil
//     case providerDomain.ProviderTwilio:
//         return NewTwilio(TwilioConfig{
//             AccountSID: cfg.Credentials["account_sid"],
//             AuthToken:  cfg.Credentials["auth_token"],
//             From:       cfg.Credentials["from"],
//         }), nil
//     default:
//         return nil, fmt.Errorf("sms: unknown provider %q", cfg.Provider)
//     }
// }

// // Dispatcher loads active provider from DB and sends
// type Dispatcher struct {
//     providerSvc *providerDomain.ProviderConfigService
// }

// func NewDispatcher(providerSvc *providerDomain.ProviderConfigService) *Dispatcher {
//     return &Dispatcher{providerSvc: providerSvc}
// }

// func (d *Dispatcher) Send(ctx context.Context, msg provider.Message) (*provider.Result, error) {
//     cfg, err := d.providerSvc.GetActiveByChannel(ctx, providerDomain.ChannelSMS)
//     if err != nil {
//         return nil, err
//     }

//     p, err := NewFromConfig(cfg)
//     if err != nil {
//         return nil, err
//     }

//     return p.Send(ctx, msg)
// }
