package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"notification/internal/config"
	"notification/internal/database"
	"notification/internal/domain/notification"
	providerDomain "notification/internal/domain/provider"
	"notification/internal/kafka"
	"notification/internal/provider/email"
	"notification/internal/provider/push"
	"notification/internal/provider/sms"
	"notification/pkg/logger"
)

func main() {
	cfg := config.Load()
	logger.SetUp(cfg.ENV, cfg.LogLevel)

	if err := cfg.Validate(); err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	mongoClient := database.Connect(cfg.MongoURI)
	db := mongoClient.Database(cfg.MongoDatabase)

	providerRepo := providerDomain.NewProviderConfigRepository(db)
	providerService := providerDomain.NewProviderConfigService(providerRepo)

	smsDispatcher := sms.NewDispatcher(providerService)
	emailDispatcher := email.NewDispatcher(providerService)
	pushDispatcher := push.NewDispatcher(providerService)

	notifRepo := notification.NewNotificationRepository(db)

	consumer := kafka.NewConsumer(
		kafka.ConsumerConfig{
			Brokers:     strings.Split(cfg.KafkaBrokers, ","),
			GroupID:     cfg.KafkaGroupID,
			MaxAttempts: cfg.KafkaMaxRetries,
		},
		map[string]kafka.Dispatcher{
			"sms":   smsDispatcher,
			"email": emailDispatcher,
			"push":  pushDispatcher,
		},
		notifRepo,
	)

	ctx, cancel := context.WithCancel(context.Background())
	go consumer.Start(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("consumer shutting down...")
	cancel()

	if err := consumer.Close(); err != nil {
		slog.Error("failed to close consumer", "error", err)
	}
	if err := mongoClient.Disconnect(context.Background()); err != nil {
		slog.Error("failed to disconnect mongodb", "error", err)
	}
}
