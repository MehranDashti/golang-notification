package database

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	client *mongo.Client
	once   sync.Once
)

func Connect(uri string) *mongo.Client {
	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		c, err := mongo.Connect(options.Client().ApplyURI(uri))
		if err != nil {
			panic("failed to connect to mongodb: " + err.Error())
		}

		if err := c.Ping(ctx, nil); err != nil {
			panic("failed to ping mongodb: " + err.Error())
		}

		client = c
	})
	return client
}
