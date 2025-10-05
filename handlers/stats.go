package handlers

import (
	"gorm.io/gorm"
	tb "gopkg.in/telebot.v3"
	"github.com/1racker/telegram-task-bot/reports"
)

func RegisterStatsHandlers(bot *tb.Bot, db *gorm.DB) {
	bot.Handle("/report", func(c tb.Context) error {
		userID := c.Sender().ID
		report := reports.GenerateWeeklyReport(db, userID)
		return c.Send(report)
	})
}