package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	ENV             string
	Host            string
	MongoURI        string
	MongoDatabase   string
	RedisAddr       string
	RedisPassword   string
	RedisDB         int
	LogLevel        string
	KafkaBrokers    string
	KafkaGroupID    string
	KafkaMaxRetries int
}

func Load() *Config {
	_ = godotenv.Load()

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		redisDB = 0
	}
	maxRetries, err := strconv.Atoi(os.Getenv("KAFKA_MAX_RETRIES"))
	if err != nil {
		maxRetries = 3
	}

	return &Config{
		Port:            os.Getenv("APP_PORT"),
		ENV:             os.Getenv("APP_ENV"),
		Host:            os.Getenv("APP_HOST"),
		MongoURI:        os.Getenv("MONGO_URI"),
		MongoDatabase:   os.Getenv("MONGO_DATABASE"),
		RedisAddr:       os.Getenv("REDIS_ADDR"),
		RedisPassword:   os.Getenv("REDIS_PASSWORD"),
		RedisDB:         redisDB,
		LogLevel:        os.Getenv("LOG_LEVEL"),
		KafkaBrokers:    os.Getenv("KAFKA_BROKERS"),
		KafkaGroupID:    os.Getenv("KAFKA_GROUP_ID"),
		KafkaMaxRetries: maxRetries,
	}
}

func (c *Config) Validate() error {
	var missing []string

	if c.Port == "" {
		missing = append(missing, "APP_PORT")
	}
	if c.MongoURI == "" {
		missing = append(missing, "MONGO_URI")
	}
	if c.MongoDatabase == "" {
		missing = append(missing, "MONGO_DATABASE")
	}
	if c.KafkaBrokers == "" {
		missing = append(missing, "KAFKA_BROKERS")
	}
	if c.KafkaGroupID == "" {
		missing = append(missing, "KAFKA_GROUP_ID")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required env vars: %s",
			strings.Join(missing, ", "))
	}

	return nil
}
