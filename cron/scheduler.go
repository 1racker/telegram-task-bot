package cron

import (
	"fmt"
	"log"
	"time"

	"github.com/1racker/telegram-task-bot/reports"
	"github.com/1racker/telegram-task-bot/storage"
	"github.com/robfig/cron/v3"
	
	tb "gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

func StartScheduler(bot *tb.Bot, db *gorm.DB, tz string, weeklyDay string) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Printf("timezone load error: %v - using Local", err)
		loc = time.Local
	}

	c := cron.New(cron.WithLocation(loc))

	_, err = c.AddFunc("@every 1m", func() {
		now := time.Now().In(loc)

		var userIDs []int64
		rows, err := db.Model(&storage.Task{}).Distinct("user_id").Rows()
		if err != nil {
			log.Printf("cron: error querying user: %v", err)
			return
		}
		defer rows.Close()

		var uid int64
		for rows.Next() {
			_ = rows.Scan(&uid)
			userIDs = append(userIDs, uid)
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
		var userIDs []int64
		rows, err := db.Model(&storage.Task{}).Distinct("user_id").Rows()
		if err != nil {
			log.Printf("cron: error querying users for weekly: %v", err)
			return
		}
		defer rows.Close()
		var uid int64
		for rows.Next() {
			_ = rows.Scan(&uid)
			userIDs = append(userIDs, uid)
		}

		for _, userID := range userIDs {
			report := reports.GenerateWeeklyReport(db, userID)
			recipient := &tb.User{ID: userID}
			if _, err := bot.Send(recipient, report); err != nil {
				log.Printf("cron: send weekly report error to %d: %v", userID, err)
			}
		}
	})
	if err != nil {
		log.Printf("cron add weekly error: %v", err)
	}
	c.Start()
}