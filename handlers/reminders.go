package handlers 

import (
	"fmt"          
	"log"           
	"os"           
	"strconv"       
	"strings"       

	"gorm.io/gorm"  
	tb "gopkg.in/telebot.v3" 
	"github.com/1racker/telegram-task-bot/reports" 
	"github.com/1racker/telegram-task-bot/storage" 
)

func RegisterReminderHandlers(bot *tb.Bot, db *gorm.DB) { 
	bot.Handle(tb.OnCallback, func(c tb.Context) error { 
		data := c.Data()                         
		if data == "" {                         
			return c.Respond(&tb.CallbackResponse{Text: "No data for processing"}) 
		}
		parts := strings.Split(data, ":")        
		if len(parts) != 2 {               
			return c.Respond(&tb.CallbackResponse{Text: "Wrong format of data"}) 
		}

		action := parts[0]                       
		idStr := parts[1]                        
		id, err := strconv.Atoi(idStr)           
		if err != nil {                          
			return c.Respond(&tb.CallbackResponse{Text: "Wrong task ID"}) 
		}

		_ = c.Respond(&tb.CallbackResponse{})

		var task storage.Task                   
		if err := db.First(&task, id).Error; err != nil { 
			log.Printf("handlers: task not found id=%d err=%v", id, err) 
			return c.Send(fmt.Sprintf("Task with ID %d not found.", id)) 
		}

		switch action {                         
		case "start":                           
			now := timeNow()                    
			task.Status = "in_progress"        
			task.StartedAt = &now               
			if err := db.Save(&task).Error; err != nil { 
				log.Printf("handlers: failed to mark start task id=%d err=%v", task.ID, err) 
				return c.Send("Error while updating task.") 
			}
			return c.Send(fmt.Sprintf("Task '%s' marker as started.", task.Title)) 

		case "done":                            
			now := timeNow()                   
			task.Status = "done"                
			task.DoneAt = &now                 
			if err := db.Save(&task).Error; err != nil { 
				log.Printf("handlers: failed to mark done task id=%d err=%v", task.ID, err)
				return c.Send("Error while updating task.")
			}
			return c.Send(fmt.Sprintf("Task '%s' marked as completed.", task.Title)) 

		case "postpone":                        
			task.Postpones++                    
			task.Status = "postponed"           
			if err := db.Save(&task).Error; err != nil { 
				log.Printf("handlers: failed to postpone task id=%d err=%v", task.ID, err)
				return c.Send("Error while postponing task.")
			}
			return c.Send(fmt.Sprintf(" Task '%s' postponed (%d).", task.Title, task.Postpones))

		case "delete":                           
			task.Status = "deleted"             
			if err := db.Save(&task).Error; err != nil { 
				log.Printf("handlers: failed to delete task id=%d err=%v", task.ID, err)
				return c.Send("Error while deleting task.")
			}
			return c.Send(fmt.Sprintf("Task '%s' marked as deleted.", task.Title)) 

		default:                                
			return c.Send("Unknown action.") 
		}
	})
}

func SendWeeklyReportToUser(bot *tb.Bot, db *gorm.DB, userID int64) error {
	reportText, chartPath, err := reports.GenerateWeeklyReportWithChart(db, userID) 
	if err != nil { 
		log.Printf("handlers: GenerateWeeklyReportWithChart error: %v", err) 
		_, _ = bot.Send(&tb.User{ID: userID}, "Error generating report.") 
		return err 
	}

	if _, err := bot.Send(&tb.User{ID: userID}, reportText); err != nil { 
		log.Printf("handlers: failed to send report text to %d: %v", userID, err) 
	}

	if chartPath != "" { 
		if _, err := os.Stat(chartPath); err == nil { 
			photo := &tb.Photo{File: tb.FromDisk(chartPath), Caption: "Chart: all (blue) vs completed (green)"} 
			if _, err := bot.Send(&tb.User{ID: userID}, photo); err != nil { 
				log.Printf("handlers: failed to send chart to %d: %v", userID, err) 
				return err 
			}
		} else { 
			log.Printf("handlers: chart file not found: %s", chartPath) 
		}
	}

	return nil 
}

func timeNow() time.Time {
	return time.Now()
}
