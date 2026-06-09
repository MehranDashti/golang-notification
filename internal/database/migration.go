package database

import (
	"context"
	"log/slog"
	"time"

	"notification/internal/database/migrations"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Migrate(client *mongo.Client, dbName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	db := client.Database(dbName)

	runners := []migrations.Migration{
		migrations.NotificationMigration{},
		migrations.ProviderMigration{},
	}

	for _, m := range runners {
		if err := m.Run(ctx, db); err != nil {
			slog.Error("migration failed",
				"migration", m.Name(),
				"error", err,
			)
			return
		}
		slog.Info("migration completed", "migration", m.Name())
	}
}
