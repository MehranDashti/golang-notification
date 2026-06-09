package dto

import (
	"time"
)

type NotificationResponse struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Channel   string            `json:"channel"`
	Status    string            `json:"status"`
	Title     string            `json:"title"`
	Body      string            `json:"body"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	SentAt    *time.Time        `json:"sent_at,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

type ProviderConfigResponse struct {
	ID          string            `json:"id"`
	Channel     string            `json:"channel"`
	Provider    string            `json:"provider"`
	Credentials map[string]string `json:"credentials"`
	IsActive    bool              `json:"is_active"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
