package handlers 

import (
	"fmt"          
	"log"           
	"strconv"       
	"strings"       

	tb "gopkg.in/telebot.v3" 
	"github.com/1racker/telegram-task-bot/storage" 
)

func RegisterReminderHandlers(bot *tb.Bot, repo storage.TaskRepository) { 
	bot.Handle(&tb.InlineButton{Unique: "btn_start"}, func(c tb.Context) error { 
		data := c.Data()                         
		if data == "" {                         
			return c.Respond(&tb.CallbackResponse{Text: "No data for processing"}) 
		}
		parts := strings.Split(data, ":")        
		if len(parts) != 2 {               
			return c.Respond(&tb.CallbackResponse{Text: "Invalid format of data"}) 
		}

		id, err := strconv.Atoi(parts[1])           
		if err != nil {                          
			return c.Respond(&tb.CallbackResponse{Text: "Invalid task ID"}) 
		}

		task, err := repo.GetByID(uint(id))
		if err != nil {
			return c.Respond(&tb.CallbackResponse{Text: "Task not found"})
		}
		
		task.Postpones++
		task.Status = "postponed"
		
		if err := repo.Update(task); err != nil {
			log.Printf("db save error: %v", err)
			return c.Respond(&tb.CallbackResponse{Text: "Error postponing task"})
		}
		return c.Respond(&tb.CallbackResponse{Text: fmt.Sprintf("Task '%s' postponed", task.Title)})
	})

	bot.Handle(&tb.InlineButton{Unique: "btn_delete"}, func(c tb.Context) error {
		data := c.Data()
		parts := strings.Split(data, ":")
		if len(parts) != 2 {
			return c.Respond(&tb.CallbackResponse{Text: "Invalid data"})
		}
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			return c.Respond(&tb.CallbackResponse{Text: "Invalid task ID"})
		}
		
		task, err := repo.GetByID(uint(id))
		if err != nil {
			return c.Respond(&tb.CallbackResponse{Text: "Task not found"})
		}
		
		task.Status = "deleted"
		
		if err := repo.Update(task); err != nil {
			log.Printf("db save error: %v", err)
			return c.Respond(&tb.CallbackResponse{Text: "Error deleting task"})
		}
		return c.Respond(&tb.CallbackResponse{Text: fmt.Sprintf("Task '%s' deleted", task.Title)})
	})
}