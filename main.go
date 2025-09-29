package main

import (
	"log"
	"telegram-task-bot/bot"
	"telegram-task-bot/config"

	"telegram-task-bot.storage"
)

func main() {
	cfg := config.LoadConfig()

	db := storage.InitDB(cfg)

	err := bot.StartBot(cfg, db)
	if err != nil {
		log.Fatal(err)
	}
}