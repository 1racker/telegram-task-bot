package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	DBURL string
	WeeklyReportDay string
	TZ string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found; reading environment variables directly")
	}

	cfg := &Config{
		TelegramToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		DBURL: os.Getenv("DATABASE_URL"),
		WeeklyReportDay: os.Getenv("WEEKLY_REPORT_DAY"),
		TZ: os.Getenv("TIMEZONE"),
	}

	if cfg.DBURL == "" {
		cfg.DBURL = "tasks.db"
	}

	if cfg.WeeklyReportDay == "" {
		cfg.WeeklyReportDay = "SUN"
	}

	if cfg.TZ == "" {
		cfg.TZ = "Europe/Bucharest"
	}

	return cfg
}