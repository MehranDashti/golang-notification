package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"notification/internal/config"
	"notification/internal/database"
	"notification/internal/domain/notification"
	providerDomain "notification/internal/domain/provider"
	rest_handler "notification/internal/handler/rest"
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

	// Notification
	notifRepo := notification.NewNotificationRepository(db)
	notifService := notification.NewNotificationService(notifRepo)
	notifHandler := rest_handler.NewNotificationHandler(notifService)

	// Health
	healthHandler := rest_handler.NewHealthHandler(mongoClient)

	// provider domain
	providerRepo := providerDomain.NewProviderConfigRepository(db)
	providerService := providerDomain.NewProviderConfigService(providerRepo)
	providerHandler := rest_handler.NewProviderHandler(providerService)

	// dispatchers — each loads active config from DB at send time
	// smsDispatcher := sms.NewDispatcher(providerSvc)
	// emailDispatcher := email.NewDispatcher(providerSvc)
	// pushDispatcher := push.NewDispatcher(providerSvc)

	// notification service gets all three
	// notifService := notification.NewNotificationService(notifRepo, smsDispatcher, emailDispatcher, pushDispatcher)

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
}
