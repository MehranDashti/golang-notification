package sms

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"notification/internal/provider"
	"time"
)

type KavenegarConfig struct {
	APIKey string
	Sender string
}

type KavenegarProvider struct {
	cfg    KavenegarConfig
	client *http.Client
}

func newKavenegar(cfg KavenegarConfig) Provider {
	return &KavenegarProvider{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (k *KavenegarProvider) Name() string {
	return "Kavenagar"
}

func (k *KavenegarProvider) Send(ctx context.Context, msg provider.Message) (*provider.Result, error) {
	endpoint := fmt.Sprintf("https://api.kavenegar.com/v1/%s/sms/send.json", k.cfg.APIKey)

	params := url.Values{}
	params.Set("sender", k.cfg.Sender)
	params.Set("receptor", msg.To)
	params.Set("message", msg.Body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("kavenegar: build request: %w", err)
	}
	req.URL.RawQuery = params.Encode()
	resp, err := k.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("kavenegar: send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Warn("failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(req.Body)
		return nil, fmt.Errorf("kavenegar: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	return &provider.Result{
		Provider:   k.Name(),
		ProviderID: "",
	}, nil
}
