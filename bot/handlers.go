package bot

import (
	"log"
	"telegram-task-bot/config"
	"telegram-task-bot/storage"

	"github.com/daixiang0/gci/pkg/config"
	thbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

func StartBot(cfg *config.Config, db *gorm.DB) error {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return err
	}

	log.Printf("🤖 bot launched as %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdateChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Comand() {
		case "start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,"Hi👋! I`m your task-manager bot")
			bot.Send(msg)

		case "addtask":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "adding of the task in progress ✍️")
			bot.Send(msg)

		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command ❓")
			bot.Send(msg)
		}
	}
	return nil
}