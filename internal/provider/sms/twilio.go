package sms

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"notification/internal/provider"
	"strings"
	"time"
)

type TwilioConfig struct {
	AccountSID string
	AuthToken  string
	From       string
}

type twilioProvider struct {
	cfg    TwilioConfig
	client *http.Client
}

func NewTwilio(cfg TwilioConfig) Provider {
	return &twilioProvider{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (t *twilioProvider) Name() string {
	return "twilio"
}

func (t *twilioProvider) Send(ctx context.Context, msg provider.Message) (*provider.Result, error) {
	endpoint := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", t.cfg.AccountSID)

	params := url.Values{}
	params.Set("From", t.cfg.From)
	params.Set("To", msg.To)
	params.Set("Body", msg.Body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("twilio: build request: %w", err)
	}
	req.SetBasicAuth(t.cfg.AccountSID, t.cfg.AuthToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("twilio: send request: %w", err)
	}
	defer func() {
		if err := req.Body.Close(); err != nil {
			slog.Warn("failed to close request body", "error", err)
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("twilio: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		SID string `json:"sid"`
	}
	if err := json.NewDecoder(req.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("twilio: decode response: %w", err)
	}

	return &provider.Result{
		Provider:   t.Name(),
		ProviderID: result.SID,
	}, nil
}
