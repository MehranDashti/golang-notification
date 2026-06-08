package sms

import (
	"context"
	provider "notification/internal/provider"
)


type Provider interface {
	Send(ctx context.Context, msg provider.Message) (*provider.Result, error)
	Name() string
}
