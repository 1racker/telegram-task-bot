package cron

import (
	"fmt"
	"log"
	"time"
	"bytes"

	"github.com/1racker/telegram-task-bot/reports"
	"github.com/1racker/telegram-task-bot/storage"
	"github.com/robfig/cron/v3"
	
	tb "gopkg.in/telebot.v3"
)

func StartScheduler(bot *tb.Bot, repo storage.TaskRepository, tz string, weeklyDay string) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Printf("timezone load error: %v - using Local", err)
		loc = time.Local
	}

	c := cron.New(cron.WithLocation(loc))

	_, err = c.AddFunc("@every 1m", func() {
		now := time.Now().In(loc)

		userIDs, err := repo.GetDistinctUserIDs()
		if err != nil {
			log.Printf("cron: error querying users: %v", err)
			return
		}

		for _, userID := range userIDs {
			if now.Hour() == 21 &&now.Minute() == 0 {
				recipient := &tb.User{ID: userID}
				_, err := bot.Send(recipient, "Reminder: din`t forget to fulfil plan for tomorrow (use /add)")
				if err != nil {
					log.Printf("cron: sebd daily reminder error to %d: %v", userID, err)
				}
			}
		}

	})
	if err != nil {
		log.Printf("cron addFunc error: %v", err)
	}

	weeklyExpr := fmt.Sprintf("0 9 * * %s", weeklyDay)
	_, err = c.AddFunc(weeklyExpr, func()  {
		userIDs, err := repo.GetDistinctUserIDs()
		if err != nil {
			log.Printf("cron: error querying users for weekly: %v", err)
			return
		}

		for _, userID := range userIDs {
			report, chartData, err := reports.GenerateWeeklyReport(repo, userID)
			if err != nil {
				log.Printf("cron: error generating weekly report for user %d: %v", userID, err)
				continue
			}

			recipient := &tb.User{ID: userID}
			if chartData != nil {
				photo := &tb.Photo {
					File: tb.FromReader(bytes.NewReader(chartData)),
					Caption: "Your Weekly Task Statisitics",
				}
				if _, err := bot.Send(recipient, photo); err != nil {
					log.Printf("cron: send weekly chart error to %d: %v", userID, err)
				}
			}
			if _, err := bot.Send(recipient, report); err != nil {
				log.Printf("cron: send weekly report error to %d: %v", userID, err)
			}
		}
	})
	if err != nil {
		log.Printf("cron add weekly error: %v", err)
	}
	c.Start()
	log.Printf("Scheduler started")
}