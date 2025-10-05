package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/1racker/telegram-task-bot/storage"
	tb "gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

func RegisterTasks(bot *tb.Bot, db *gorm.DB) {
	bot.Handle("/add", func(c tb.Context) error {
		payload := strings.TrimSpace(c.Message().Payload)

		if payload == "" {
			help := "Использование:\n/add Название|Категория|Приоритет(1-3)|HH:MM|Длительность(мин)\nПример:\n/add Сделать тест|Учеба|1|19:00|60"
			return c.Send(help)
		}

		parts := strings.Split(payload, "|")
		if len(parts) < 5 {
			return c.Send("Incorrect format. Используй /add Название|Категория|Приоритет|HH:MM|Длительность")
		}

		title := strings.TrimSpace(parts[0])
		category := strings.TrimSpace(parts[1])
		priority, err := strconv.Atoi(strings.TrimSpace(parts[2]))
		if err != nil {
			priority = 2
		}
		timeStr := strings.TrimSpace(parts[3])
		duration, err := strconv.Atoi(strings.TrimSpace(parts[4]))
		if err != nil {
			duration = 60
		}
		now := time.Now()
		hourMin := strings.Split(timeStr, ":")
		hour, _ := strconv.Atoi(hourMin[0])
		min, _ := strconv.Atoi(hourMin[1])
		start := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location())

		task := storage.Task{
			UserID: c.Sender().ID,
			Title: title,
			Category: category,
			Priority: priority,
			StartTime: start,
			Duration: duration,
			Status: "new",
		}

		if err := db.Create(&task).Error; err != nil {
			log.Printf("db create error: %v", err)
			return c.Send("Error while saving the data.")
		}
		return c.Send(fmt.Sprintf("Задача '%s' добавлена на %s", title, start.Format("15:04")))
	})
	bot.Handle("/today", func(c tb.Context) error {
		userID := c.Sender(). ID
		startOfDay := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
		endOfDay := startOfDay.Add(24 * time.Hour)

		var tasks []storage.Task
		if err := db.Where("user_id = ? AND start_time >= ? AND start_time < ?", userID, startOfDay, endOfDay).Order("priority asc").Find(&tasks).Error; err != nil {
		return c.Send("Error while executing task")
		}

		if len(tasks) == 0 {
			return c.Send("На сегодня задач нет")
		}

		for _, t := range tasks {
			text := fmt.Sprintf("%s\nКатегория: %s\nПриоритет: %d\nВремя: %s\nДлительность: %d мин\nСтатус: %s",
			t.Title, t.Category, t.Priority, t.StartTime.Format("15:04"), t.Duration, t.Status)

			startBtn := tb.InlineButton{Unique: "btn_start", Text: "Начать", Data: fmt.Sprintf("start:%d", t.ID)}
			doneBtn := tb.InlineButton{Unique: "btn_done", Text: "Завершить", Data: fmt.Sprintf("done:%d", t.ID)}
			postponeBtn := tb.InlineButton{Unique: "btn_postpone", Text: "Отложить", Data: fmt.Sprintf("postpone:%d", t.ID)}
			deleteBtn := tb.InlineButton{Unique: "btn_delete", Text: "Удалить", Data: fmt.Sprintf("delete:%d", t.ID)}

			inlineKeys := [][]tb.InlineButton{
				{startBtn, doneBtn},
				{postponeBtn, deleteBtn},
			}

			if err := c.Send(text, &tb.ReplyMarkup{InlineKeyboard: inlineKeys}); err != nil {
				log.Printf("send message error: %v", err)
			}
		}
		return nil
	})
}