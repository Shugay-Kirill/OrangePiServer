package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken      string
	Debug              bool
	MaxLengthAPIOutput int
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		TelegramToken:      getEnv("TELEGRAM_BOT_TOKEN", ""),
		Debug:              getEnvAsBool("DEBUG", false),
		MaxLengthAPIOutput: getEnvAsInt("MAX_LENGTH_MESSEGE_API", 200),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
