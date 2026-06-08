package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"notification/internal/config"
	"notification/internal/database"
	rest_handler "notification/internal/handler/rest"
	"notification/internal/router"
)

func main() {
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	db := database.Connect(cfg.MongoURI)

	slog.Info("server starting",
		"port", cfg.Port,
	)

	//Handlers
	healthHandler := rest_handler.NewHealthHandler(db)

	r := router.Setup(
		healthHandler,
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
}
