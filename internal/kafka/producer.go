package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
)

type ProducerConfig struct {
	Brokers []string
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(cfg ProducerConfig) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(cfg.Brokers...),
			Topic:                  TopicNotifications,
			Balancer:               &kafka.Hash{},
			WriteTimeout:           10 * time.Second,
			RequiredAcks:           kafka.RequireAll,
			AllowAutoTopicCreation: true,
		},
	}
}

func (p *Producer) Publish(ctx context.Context, msg NotificationMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("producer: marshal message: %w", err)
	}

	if err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(msg.UserID),
		Value: data,
	}); err != nil {
		return fmt.Errorf("producer: write message: %w", err)
	}

	slog.Info("notification published to kafka",
		"notification_id", msg.NotificationID,
		"channel", msg.Channel,
	)

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
