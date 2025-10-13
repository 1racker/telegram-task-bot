package handlers

import (
	"bytes"

	"github.com/1racker/telegram-task-bot/reports"
	"github.com/1racker/telegram-task-bot/storage"
	tb "gopkg.in/telebot.v3"
)

func RegisterStatsHandlers(bot *tb.Bot, repo storage.TaskRepository) {
	bot.Handle("/report", func(c tb.Context) error {
		userID := c.Sender().ID
		report, chartData, err := reports.GenerateWeeklyReport(repo, userID)
		if err != nil {
			return c.Send("Error generating report: " +err.Error())
		}
		if chartData != nil {
			photo := &tb.Photo{
				File: tb.FromReader(bytes.NewReader(chartData)),
				Caption: "Weekly Task Statistics Chart",
			}
			if err := c.Send(photo); err != nil {
				return c.Send(report)
			}
		}
		return c.Send(report)
	})
}