package main

import (
	"log"
	"net/http"
	"os"


	"github.com/1racker/telegram-task-bot/config"
	"github.com/1racker/telegram-task-bot/cron"
	"github.com/1racker/telegram-task-bot/handlers"
	"github.com/1racker/telegram-task-bot/storage"
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

	webhookURL := os.Getenv("WEBHOOK_URL")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	pref := tb.Settings{
		Token: cfg.TelegramToken,
		Poller: &tb.Webhook{
			Listen: ":" + port,
			Endpoint: &tb.WebhookEndpoint{
				PublicURL: webhookURL,
			},
		},
	}

	bot ,err := tb.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	handlers.RegisterTasks(bot, repo)
	handlers.RegisterReminderHandlers(bot, repo)
	handlers.RegisterStatsHandlers(bot, repo)

	cron.StartScheduler(bot, repo, cfg.TZ, cfg.WeeklyReportDay)

	log.Printf("Bot launched on port %s via webhook %s", port, webhookURL)
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	go http.ListenAndServe(":"+port, nil)
}