package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	DBPath string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using environment variables")
	}

	return &Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		DBPath: os.Getenv("DB_PATH"),
	}
}