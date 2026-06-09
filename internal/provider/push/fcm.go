package push

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"notification/internal/provider"
	"time"
)

type FCMConfig struct {
	ServerKey string
}

type fcmProvider struct {
	cfg    FCMConfig
	client *http.Client
}

func NewFCM(cfg FCMConfig) Provider {
	return &fcmProvider{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (f *fcmProvider) Name() string { return "fcm" }

func (f *fcmProvider) Send(ctx context.Context, msg provider.Message) (*provider.Result, error) {
	payload := map[string]any{
		"to": msg.To,
		"notification": map[string]string{
			"title": msg.Title,
			"body":  msg.Body,
		},
		"data": msg.Metadata,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("fcm: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://fcm.googleapis.com/fcm/send", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("fcm: build request: %w", err)
	}
	req.Header.Set("Authorization", "key="+f.cfg.ServerKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fcm: send request: %w", err)
	}
	defer func() {
		if err := req.Body.Close(); err != nil {
			slog.Warn("failed to close request body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fcm: unexpected status %d", resp.StatusCode)
	}

	return &provider.Result{Provider: f.Name()}, nil
}
