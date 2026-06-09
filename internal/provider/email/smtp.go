package email

import (
	"context"
	"fmt"
	"net/smtp"
	"notification/internal/provider"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type smtpProvider struct {
	cfg SMTPConfig
}

func NewSMTP(cfg SMTPConfig) Provider {
	return &smtpProvider{cfg: cfg}
}

func (s *smtpProvider) Name() string { return "smtp" }

func (s *smtpProvider) Send(_ context.Context, msg provider.Message) (*provider.Result, error) {
	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)

	body := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		s.cfg.From, msg.To, msg.Title, msg.Body)

	addr := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)
	if err := smtp.SendMail(addr, auth, s.cfg.From, []string{msg.To}, []byte(body)); err != nil {
		return nil, fmt.Errorf("smtp: send mail: %w", err)
	}

	return &provider.Result{Provider: s.Name()}, nil
}
