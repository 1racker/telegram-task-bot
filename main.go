package main

import (
	"log"
	"github.com/1racker/telegram-task-bot/config"
	"github.com/1racker/telegram-task-bot/cron"
	"github.com/1racker/telegram-task-bot/storage"
	"github.com/1racker/telegram-task-bot/handlers"
	"time"
	tb "gopkg.in/telebot.v3"
)

func main() {
	cfg := config.Load()

	if cfg.TelegramToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is empty")
	}
	log.Printf("Token length: %d characters", len(cfg.TelegramToken))
	log.Printf("Token starts with: %s", cfg.TelegramToken[:10])

	db := storage.InitDB(cfg.DBPath)

	repo := storage.NewTaskRepository(db)

	pref := tb.Settings{
		Token: cfg.TelegramToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	}

	bot ,err := tb.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	handlers.RegisterTasks(bot, repo)
	handlers.RegisterReminderHandlers(bot, repo)
	handlers.RegisterStatsHandlers(bot, repo)

	cron.StartScheduler(bot, repo, cfg.TZ, cfg.WeeklyReportDay)

	log.Println("Bot launched...")
	bot.Start()
}