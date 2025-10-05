package main

import (
	"log"
	"github.com/1racker/telegram-task-bot/config"
	"github.com/1racker/telegram-task-bot/cron"
	"github.com/1racker/telegram-task-bot/storage"
	"github.com/1racker/telegram-task-bot/handlers"
	"time"
	"gopkg.in/telebot.v3"
)

func main() {
	cfg := config.Load()

	db := storage.InitDB("tasks.db")

	pref := telebot.Settings{
		Token: cfg.TelegramToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot ,err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	handlers.RegisterTasks(bot, db)
	handlers.RegisterReminderHandlers(bot, db)
	handlers.RegisterStatsHandlers(bot, db)

	cron.StartScheduler(bot, db, cfg.TZ, cfg.WeeklyReportDay)

	log.Println("Bot launched...")
	bot.Start()
}