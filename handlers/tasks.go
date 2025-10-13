package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/1racker/telegram-task-bot/storage"
	tb "gopkg.in/telebot.v3"
)

type TaskValidator interface {
	ValidateTaskInput(title, category, timeStr string, prioruty, duration int) error
}

type DefaultTaskValidator struct {}

func (v *DefaultTaskValidator) ValidateTaskInput(title, category, timeStr string, priority, duration int) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if strings.TrimSpace(category) == "" {
		return fmt.Errorf("category cannot be empty")
	}
	if priority < 1 || priority > 3 {
		return fmt.Errorf("priority must be between 1-3")
	}
	if duration <= 0 {
		return fmt.Errorf("duration must be positive")
	}

	hourMin := strings.Split(timeStr, ":")
	if len(hourMin) != 2 {
		return fmt.Errorf("invalid time format, use HH:MM")
	}
	hour, err := strconv.Atoi(hourMin[0])
	if err != nil || hour < 0 || hour > 23 {
		return fmt.Errorf("hour must be between 0-23")
	}
	min, err := strconv.Atoi(hourMin[1])
	if err != nil || min < 0 || min > 59 {
		return fmt.Errorf("minutes must be between 0-59")
	}
	return nil
}

func RegisterTasks(bot *tb.Bot, repo storage.TaskRepository) {
	validator := &DefaultTaskValidator{}

	bot.Handle("/add", func(c tb.Context) error {
		payload := strings.TrimSpace(c.Message().Payload)

		if payload == "" {
			help := "Usage:\n/add Title|Category|Priority(1-3)|HH:MM|Duration(min)\nExample:\n/add Do test|Study|1|19:00|60"
			return c.Send(help)
		}

		parts := strings.Split(payload, "|")
		if len(parts) < 5 {
			return c.Send("Incorrect format. Use: /add Title|Category|Priority|HH:MM|Duration")
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

		if err := validator.ValidateTaskInput(title, category, timeStr, priority, duration); err != nil {
			return c.Send(fmt.Sprintf("Validation error: %v", err))
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

		if err := repo.Create(&task); err != nil {
			log.Printf("db create error: %v", err)
			return c.Send("Error while saving the data.")
		}
		return c.Send(fmt.Sprintf("Task '%s' added for %s", title, start.Format("15:04")))
	})
	bot.Handle("/today", func(c tb.Context) error {
		userID := c.Sender().ID
		tasks, err := repo.GetTodayTasks(userID)
		if err != nil {
			 log.Printf("db query error: %v", err)
			 return c.Send("Error while executing task")

		 }
		 if len(tasks) == 0 {
			return c.Send("No tasks for today")
		}

				for _, t := range tasks {
			text := fmt.Sprintf("%s\nCategory: %s\nPriority: %d\nTime: %s\nDuration: %d min\nStatus: %s",
			t.Title, t.Category, t.Priority, t.StartTime.Format("15:04"), t.Duration, t.Status)

			startBtn := tb.InlineButton{Unique: "btn_start", Text: "Start", Data: fmt.Sprintf("start:%d", t.ID)}
			doneBtn := tb.InlineButton{Unique: "btn_done", Text: "Complete", Data: fmt.Sprintf("done:%d", t.ID)}
			postponeBtn := tb.InlineButton{Unique: "btn_postpone", Text: "Postpone", Data: fmt.Sprintf("postpone:%d", t.ID)}
			deleteBtn := tb.InlineButton{Unique: "btn_delete", Text: "Delete", Data: fmt.Sprintf("delete:%d", t.ID)}

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

