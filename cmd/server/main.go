package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"notification/internal/config"
	"notification/internal/database"
	"notification/internal/domain/notification"
	providerDomain "notification/internal/domain/provider"
	rest_handler "notification/internal/handler/rest"
	"notification/internal/kafka"
	"notification/internal/router"
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
	database.Migrate(mongoClient, cfg.MongoDatabase)
	db := mongoClient.Database(cfg.MongoDatabase)

	producer := kafka.NewProducer(kafka.ProducerConfig{
		Brokers: strings.Split(cfg.KafkaBrokers, ","),
	})

	// Health
	healthHandler := rest_handler.NewHealthHandler(mongoClient)

	// provider domain
	providerRepo := providerDomain.NewProviderConfigRepository(db)
	providerService := providerDomain.NewProviderConfigService(providerRepo)
	providerHandler := rest_handler.NewProviderHandler(providerService)

	// dispatchers
	// smsDispatcher := sms.NewDispatcher(providerService)
	// emailDispatcher := email.NewDispatcher(providerService)
	// pushDispatcher := push.NewDispatcher(providerService)

	// notification
	notifRepo := notification.NewNotificationRepository(db)
	notifService := notification.NewNotificationService(notifRepo, producer)
	notifHandler := rest_handler.NewNotificationHandler(notifService)

	r := router.Setup(
		healthHandler,
		notifHandler,
		providerHandler,
	)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
		// Timeouts — protect against slow clients
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Disconnect(ctx); err != nil {
		slog.Error("failed to disconnect mongodb", "error", err)
	}
	// shutdown
	defer func() {
		if err := producer.Close(); err != nil {
			slog.Error("failed to close kafka producer", "error", err)
		}
	}()
}
