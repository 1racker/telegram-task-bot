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