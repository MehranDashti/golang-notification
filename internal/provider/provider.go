package provider

import "context"

type Message struct {
	To       string
	Title    string
	Body     string
	Metadata map[string]string
}

type Result struct {
	ProviderID string
	Provider   string
}

type Provider interface {
	Send(ctx context.Context, msg Message) (*Result, error)
	Name() string
}
