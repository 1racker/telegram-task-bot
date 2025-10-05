package reports

import (
	"fmt"
	"sort"
	"time"

	"gorm.io/gorm"
	"github.com/1racker/telegram-task-bot/storage"
)

type DayStat struct {
	Date string
	Total int
	Done int
	Postponed int
	Deleted int
	Started int
}

func GenerateWeeklyReport(db *gorm.DB, userID int64) string {
	today := time.Now() 
	from := today.AddDate(0, 0, -6) 

	var tasks []storage.Task 
	db.Where("user_id = ? AND created_at >= ?", userID, from).Find(&tasks) 

	if len(tasks) == 0 {
		return " За последнюю неделю задач не найдено."
	}

	byDay := map[string]*DayStat{} 

	total := 0
	done := 0
	postponed := 0
	deleted := 0
	started := 0

	var completionMinutes []float64 

	for _, t := range tasks { 
		dayKey := t.CreatedAt.Format("2006-01-02") 
		if _, ok := byDay[dayKey]; !ok {
			byDay[dayKey] = &DayStat{Date: dayKey} 
		}

		byDay[dayKey].Total++ 
		total++ 

		switch t.Status { 
		case "done": 
			byDay[dayKey].Done++ 
			done++ 
			if t.DoneAt != nil { 
				completionMinutes = append(completionMinutes, t.DoneAt.Sub(t.CreatedAt).Minutes()) 
			}
		case "postponed": 
			byDay[dayKey].Postponed++ 
			postponed++ 
		case "deleted": 
			byDay[dayKey].Deleted++ 
			deleted++ 
		case "in_progress": 
			byDay[dayKey].Started++ 
			started++ 
		default:
			
		}
	}

	var days []string 
	for d := range byDay { 
		days = append(days, d) 
	}
	sort.Strings(days)

	bestDay := ""
	worstDay := ""
	bestCount := -1
	worstCount := 1_000_000

	details := "" 
	for _, d := range days { 
		ds := byDay[d] 
		details += fmt.Sprintf("%s → всего: %d |  %d |  %d |  %d |  %d\n",
			ds.Date, ds.Total, ds.Done, ds.Postponed, ds.Deleted, ds.Started)

		if ds.Done > bestCount {
			bestCount = ds.Done
			bestDay = ds.Date
		}
		if ds.Done < worstCount {
			worstCount = ds.Done
			worstDay = ds.Date
		}
	}

	percentDone := 0
	if total > 0 {
		percentDone = (done * 100) / total
	}

	avgCompletion := 0.0
	if len(completionMinutes) > 0 {
		sum := 0.0
		for _, v := range completionMinutes {
			sum += v
		}
		avgCompletion = sum / float64(len(completionMinutes))
	}

	report := fmt.Sprintf("Отчёт за последние 7 дней:\n\nВсего задач: %d\n Выполнено: %d\n⏳ Отложено: %d\n Удалено: %d\n В процессе: %d\n\nПроцент выполнения: %d%%\n\nСамый продуктивный день: %s (%d задач)\nНаименее продуктивный день: %s (%d задач)\n\nСреднее время выполнения: %.1f минут\n\n Подробно по дням:\n\n%s",
		total, done, postponed, deleted, started, percentDone, bestDay, bestCount, worstDay, worstCount, avgCompletion, details)

	return report 
}