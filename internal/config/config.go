package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	ENV           string
	Host          string
	DSN           string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	LogLevel      string
}

func Load() *Config {
	_ = godotenv.Load()

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		redisDB = 0
	}

	return &Config{
		Port:          os.Getenv("APP_PORT"),
		ENV:           os.Getenv("APP_ENV"),
		Host:          os.Getenv("APP_HOST"),
		DSN:           dsn,
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       redisDB,
		LogLevel:      os.Getenv("LOG_LEVEL"),
	}
}

func (c *Config) Validate() error {
	var missing []string

	if c.Port == "" {
		missing = append(missing, "APP_PORT")
	}
	if c.DSN == "" {
		missing = append(missing, "DB_USER/DB_PASS/DB_HOST/DB_PORT/DB_NAME")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required env vars: %s",
			strings.Join(missing, ", "))
	}

	return nil
}
