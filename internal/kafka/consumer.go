package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"notification/internal/provider"

	"github.com/segmentio/kafka-go"
)

type Dispatcher interface {
	Send(ctx context.Context, msg provider.Message) (*provider.Result, error)
}

type StatusUpdater interface {
	UpdateStatus(ctx context.Context, id string, status string, errMsg string) error
}

type ConsumerConfig struct {
	Brokers     []string
	GroupID     string
	MaxAttempts int
}

type Consumer struct {
	reader      *kafka.Reader
	dlqWriter   *kafka.Writer
	dispatchers map[string]Dispatcher
	updater     StatusUpdater
	maxAttempts int
}

func NewConsumer(
	cfg ConsumerConfig,
	dispatchers map[string]Dispatcher,
	updater StatusUpdater,
) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Brokers,
		Topic:          TopicNotifications,
		GroupID:        cfg.GroupID,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
	})

	dlqWriter := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Brokers...),
		Topic:                  TopicNotificationsDLQ,
		AllowAutoTopicCreation: true,
	}

	return &Consumer{
		reader:      reader,
		dlqWriter:   dlqWriter,
		dispatchers: dispatchers,
		updater:     updater,
		maxAttempts: cfg.MaxAttempts,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	slog.Info("kafka consumer started")
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				slog.Info("consumer: context cancelled, stopping")
				return
			}
			slog.Error("consumer: fetch message failed", "error", err)
			continue
		}
		c.process(ctx, msg)
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			slog.Error("consumer: commit failed", "error", err)
		}
	}
}

func (c *Consumer) process(ctx context.Context, msg kafka.Message) {
	var n NotificationMessage
	if err := json.Unmarshal(msg.Value, &n); err != nil {
		slog.Error("consumer: unmarshal failed", "error", err)
		return
	}

	slog.Info("consumer: processing",
		"notification_id", n.NotificationID,
		"channel", n.Channel,
		"attempt", n.Attempt,
	)

	_ = c.updater.UpdateStatus(ctx, n.NotificationID, "processing", "")

	dispatcher, ok := c.dispatchers[n.Channel]
	if !ok {
		slog.Error("consumer: no dispatcher", "channel", n.Channel)
		_ = c.updater.UpdateStatus(ctx, n.NotificationID, "failed",
			"no dispatcher for channel: "+n.Channel)
		return
	}

	result, err := dispatcher.Send(ctx, provider.Message{
		To:       n.Metadata["to"],
		Title:    n.Title,
		Body:     n.Body,
		Metadata: n.Metadata,
	})
	if err != nil {
		n.Attempt++
		slog.Warn("consumer: dispatch failed",
			"notification_id", n.NotificationID,
			"attempt", n.Attempt,
			"error", err,
		)
		if n.Attempt >= c.maxAttempts {
			c.sendToDLQ(ctx, n, err.Error())
			_ = c.updater.UpdateStatus(ctx, n.NotificationID, "failed",
				fmt.Sprintf("max attempts reached: %s", err.Error()))
			return
		}
		c.retry(ctx, n)
		return
	}

	slog.Info("consumer: notification sent",
		"notification_id", n.NotificationID,
		"provider", result.Provider,
		"provider_id", result.ProviderID,
	)
	_ = c.updater.UpdateStatus(ctx, n.NotificationID, "sent", "")
}

func (c *Consumer) retry(ctx context.Context, n NotificationMessage) {
	backoff := time.Duration(n.Attempt*n.Attempt) * 5 * time.Second
	slog.Info("consumer: retrying",
		"notification_id", n.NotificationID,
		"backoff", backoff,
	)
	time.Sleep(backoff)

	data, _ := json.Marshal(n)
	w := &kafka.Writer{
		Addr:                   kafka.TCP(c.reader.Config().Brokers...),
		Topic:                  TopicNotifications,
		AllowAutoTopicCreation: true,
	}
	defer func() {
		if err := w.Close(); err != nil {
			slog.Warn("consumer: retry writer close failed", "error", err)
		}
	}()

	if err := w.WriteMessages(ctx, kafka.Message{
		Key:   []byte(n.UserID),
		Value: data,
	}); err != nil {
		slog.Error("consumer: retry publish failed", "error", err)
	}
}

func (c *Consumer) sendToDLQ(ctx context.Context, n NotificationMessage, reason string) {
	slog.Error("consumer: sending to DLQ",
		"notification_id", n.NotificationID,
		"reason", reason,
	)
	data, _ := json.Marshal(n)
	if err := c.dlqWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(n.NotificationID),
		Value: data,
	}); err != nil {
		slog.Error("consumer: DLQ write failed", "error", err)
	}
}

func (c *Consumer) Close() error {
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("consumer: close reader: %w", err)
	}
	return c.dlqWriter.Close()
}
