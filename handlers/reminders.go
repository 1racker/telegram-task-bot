package handlers

import (
	"time"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/1racker/telegram-task-bot/storage"
	tb "gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

func RegisterReminderHandlers(bot *tb.Bot, db *gorm.DB) {
	bot.Handle(&tb.InlineButton{Unique: "btn_start"}, func(c tb.Context) error {
		data := c.Data()
		parts := strings.Split(data, ":")
		if len(parts) != 2 {
			return c.Respond(&tb.CallbackResponse{Text: "Wrong data"})
		}
		id, _ := strconv.Atoi(parts[1])
		var task storage.Task
		if err := db.First(&task, id).Error; err != nil {
			return c.Respond(&tb.CallbackResponse{Text: "Task not found"})
		}
		now := time.Now()
		task.Status = "in_progress"
		task.StartedAt = &now
		if err := db.Save(&task).Error; err != nil {
			log.Printf("db save error: %v", err)
			return c.Respond(&tb.CallbackResponse{Text: "Error updating a task"})
		}
		return c.Respond(&tb.CallbackResponse{Text: fmt.Sprintf("Task '%s' marked as started", task.Title)})
	})
	bot.Handle(tb.InlineButton{Unique: "btn_done"}, func(c tb.Context) error {
		data := c.Data()
		parts := strings.Split(data, ":")
		if len(parts) != 2 {
			return c.Respond(&tb.CallbackResponse{Text: "Wrong data"})
		}
		id, _ := strconv.Atoi(parts[1])
		var task storage.Task
		if err := db.First(&task, id).Error; err != nil {
			return c.Respond(&tb.CallbackResponse{Text: "Task not found"})
		}
		now := time.Now()
		task.Status = "done"
		task.DoneAt = &now
		if err := db.Save(&task).Error; err != nil {
			log.Printf("db save error: %v", err)
			return c.Respond(&tb.CallbackResponse{Text: "Error updating a task"})
		}
		return c.Respond(&tb.CallbackResponse{Text: fmt.Sprintf("Task '%s' completed", task.Title)})
	})

	bot.Handle(&tb.InlineButton{Unique: "btn_postpone"}, func(c tb.Context) error {
		data := c.Data()
		parts := strings.Split(data, ":")
		if len(parts) != 2 {
			return c.Respond(&tb.CallbackResponse{Text: "Wrong data"})
		}
		id, _ := strconv.Atoi(parts[1])
		var task storage.Task
		if err := db.First(&task, id).Error; err != nil {
			return c.Respond(&tb.CallbackResponse{Text: "Task not found"})
		}
		task.Postpones++
		task.Status = "postponed"
		if err := db.Save(&task).Error; err != nil {
			log.Printf("db save error: %v", err)
			return c.Respond(&tb.CallbackResponse{Text: "Error when postponing a task"})
		}
		return c.Respond(&tb.CallbackResponse{Text: fmt.Sprintf("Task '%s' postponed", task.Title)})
	})

	bot.Handle(&tb.InlineButton{Unique: "btn_delete"}, func(c tb.Context) error{
		data := c.Data()
		parts := strings.Split(data, ":")
		if len(parts) != 2 {
			return c.Respond(&tb.CallbackResponse{Text: "Wrong data"})
		}
		id, _ := strconv.Atoi(parts[1])
		var task storage.Task
		if err := db.First(&task, id).Error; err != nil {
			return c.Respond(&tb.CallbackResponse{Text: "Task not found"})
		}
		task.Status = "deleted"
		if err := db.Save(&task).Error; err != nil {
			log.Printf("db save error: %v", err)
			return c.Respond(&tb.CallbackResponse{Text: "Error when deleting a task"})
		}
		return c.Respond(&tb.CallbackResponse{Text: fmt.Sprint("Task '%s' deleted", task.Title)})
	})
}