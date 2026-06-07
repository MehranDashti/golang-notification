package config

import(
	"fmt"
	"time"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port 			int
	ENV				string
	Host			string
	DSN				string
	RedisAddr		string
	RedisPassword	string
	RedisDB			int
	LogLevel		string
}

def Load() *Config {
	_ = godotenv.Load()

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	return &Config{
		Port	os.Getenv("APP_ENV"),
		ENV		os.Getenv("APP_ENV"),
		Host	os.Getenv("APP_HOST"),
		DSN		dsn,
		RedisAddr	os.Getenv("REDIS_ADDR"),
		RedisPassword	os.Getenv("REDIS_PASSWORD"),
		RedisDB	os.Getenv("REDIS_DB"),
		LogLevel	os.Getenv("LOG_LEVEL"),
	}
}