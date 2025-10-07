package reports

import (
	"fmt"
	"time"
	"os"
	"log"

	"gorm.io/gorm"
	"github.com/1racker/telegram-task-bot/storage"
	"github.com/wcharczuk/go-chart/v2"
)


func GenerateWeeklyReport(db *gorm.DB, userID int64) (string, string, error) {
	now := time.Now() 
	from := now.AddDate(0, 0, -6) 

	var tasks []storage.Task 
	if err := db.Where("user_id = ? AND created_at >= ?", userID, from).Find(&tasks).Error; err != nil {
		return "", "", fmt.Errorf("db query error: %w", err)
	} 

	if len(tasks) == 0 {
		return " No tasks for the previous week.", "", nil
	}

	byDay := map[string]*struct {
		Total int
		Done int
		Postponed int
		Deleted int
		Started int
	} {

	}

	loc := time.Now().Location()
	var dayKeys []string
	for i := 0; i < 7; i++ {
		d := from.AddDate(0, 0, i)
		key := d.Format("2006-01-02")
		dayKeys = append(dayKeys, key)

		byDay[key] = struct {
			Total int
			Done int
			Postponed int
			Deleted int
			Started int
		}{0, 0, 0, 0, 0}
	}

	var completionMinutes []float64 

	for _, t := range tasks { 
		dayKey := t.CreatedAt.In(loc).Format("2006-01-02") 
		agg := byDay[dayKey]
		agg.Total++
		switch t.Status { 
		case "done": 
			agg.Done++  
			if t.DoneAt != nil { 
				completionMinutes = append(completionMinutes, t.DoneAt.Sub(t.CreatedAt).Minutes()) 
			}
		case "postponed": 
			agg.Postponed++ 
		case "deleted": 
			agg.Deleted++ 
		case "in_progress": 
			agg.Started++ 
		default:
			byDay[dayKey] = agg
		}
	}

	total := 0
	done := 0
	postponed := 0
	deleted := 0
	started := 0

	for _, k := range dayKeys { 
		v := byDay[k]
		total += v.Total
		done += v.Done
		postponed += v.Postponed
		deleted += v.Deleted
		started += v.Started
	}

	percentDone := 0
	if total > 0 {
		percentDone = done * 100 / total
	}

	avgCompletion := 0.0
	if len(completionMinutes) > 0 {
		sum := 0.0
		for _, m := range completionMinutes {
			sum += m
		}
		avgCompletion = sum / float64(len(completionMinutes))
	}

	bestDay := ""
	worstDay := ""
	bestCount := -1
	worstCount := int(1_000_000)

	for _, k := range dayKeys {
		d := byDay[k]
		if d.Done > bestCount {
			bestCount = d.Done
			bestDay = k
		}
		if d.Done < worstCount {
			worstCount = d.Done
			worstCount = k
		}
	}

	details := "" 
	for _, k := range daysKeys { 
		d := byDay[k] 
		details += fmt.Sprintf("%s → total: %d |  %d |  %d |  %d |  %d\n",
			k, d.Total, d.Done, d.Postponed, d.Deleted, d.Started)
	}

reportText := fmt.Sprintf(
		"Report for the last 7 days:\n\nAll tasks: %d\nCompleted: %d\nPostponed: %d\nDeleted: %d\nIn progress: %d\n\nPercentage of completed: %d%%\n\nThe most productive day: %s (%d tasks)\nThe less productive day: %s (%d tasks)\n\nAverage compliting time: %.1f minutes\n\nDay by day:\n\n%s",
		total, done, postponed, deleted, started, percentDone, bestDay, bestCount, worstDay, worstCount, avgCompletion, details,
	)

	var xValues []time.Time
	var totalValues []float64
	var doneValues []float64

	for _, k := range dayKeys {
		tm, err := time.ParseInLocation("2006-01-02", k, loc)
		if err != nil {
			log.Printf("reports: failed to parse day key %s: %v", k, err)
			continue
		}
		xValues = append(xValues, tm)
		totalValues = append(totalValues, float64(byDay[k].Total))
		doneValues = append(doneValues, float64(byDay[k].Done))
	}

	if len(xValues) == 0 {
		return reportText, "", nil
	}

	totalSeries := chart.TimeSeries{
		Name:    "Total",         
		XValues: xValues,          
		YValues: totalValues,     
		Style: chart.Style{
			Show:        true,
			StrokeWidth: 2.0,
		},
	}

	doneSeries := chart.TimeSeries{
		Name:    "Done",           
		XValues: xValues,          
		YValues: doneValues,       
		Style: chart.Style{
			Show:        true,
			StrokeWidth: 3.0,
			StrokeColor: chart.ColorGreen, 
			FillColor:   chart.ColorGreen.WithAlpha(64),
		},
	}

	g := chart.Chart{
		Title: fmt.Sprintf("Tasks: total vs done — last 7 days (user %d)", userID), 
		XAxis: chart.XAxis{
			Name:           "Date",
			ValueFormatter: chart.TimeDateValueFormatter, 
		},
		Series: []chart.Series{
			totalSeries, 
			doneSeries,  
		},
	}

	outDir := "reports"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		log.Printf("reports: failed to create reports dir: %v", err)
		return reportText, "", nil
	}

	fileName := fmt.Sprintf("weekly_%d_%s.png", userID, time.Now().Format("20060102"))
	filePath := outDir + "/" + fileName

	f, err := os.Create(filePath)
	if err != nil {
		log.Printf("reports: failed to create file %s: %v", filePath, err)
		return reportText, "", nil
	}
	defer f.Close()

	if err := g.Render(chart.PNG, f); err != nil {
		log.Printf("reports: failed to render chart: %v", err)
		return reportText, "", nil
	}

	return reportText, filePath, nil
}