package reports

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/1racker/telegram-task-bot/storage"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

type DayStat struct {
	Date string
	Total int
	Done int
	Postponed int
	Deleted int
	Started int
}

func GenerateWeeklyReport(repo storage.TaskRepository, userID int64) (string, []byte, error) {
	today := time.Now() 
	from := today.AddDate(0, 0, -6) 

	tasks, err := repo.GetWeeklyTasks(userID, from, today)
	if err != nil {
		return "Error retrieving data for the last week.", nil, err
	}

	if len(tasks) == 0 {
		return " No tasks for the previous week.", nil, nil
	}

	byDay := make(map[string]*DayStat)

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
			if t.StartedAt != nil && t.DoneAt != nil { 
				completionMinutes = append(completionMinutes, t.DoneAt.Sub(*t.StartedAt).Minutes()) 
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
	worstCount := int(1_000_000)

	details := "" 
	for _, d := range days { 
		ds := byDay[d] 
		details += fmt.Sprintf("%s → total: %d | done: %d | postponed: %d | deleted: %d | in progress: %d\n",
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
		percentDone = done * 100 / total
	}

	avgCompletion := 0.0
	if len(completionMinutes) > 0 {
		sum := 0.0
		for _, v := range completionMinutes {
			sum += v
		}
		avgCompletion = sum / float64(len(completionMinutes))
	}


report := fmt.Sprintf("Weekly Report (Last 7 days):\n\nTotal tasks: %d\n Done: %d\nPostponed: %d\nDeleted: %d\nIn progress: %d\n\nCompletion rate: %d%%\n\nMost productive day: %s (%d tasks)\nLeast productive day: %s (%d tasks)\n\nAverage completion time: %.1f minutes\n\nDaily breakdown:\n\n%s",
		total, done, postponed, deleted, started, percentDone, bestDay, bestCount, worstDay, worstCount, avgCompletion, details)

		chartBytes, err := generateProductivityChart(days, byDay)
		if err != nil {
			log.Printf("Error generating chart: %v", err)
			return report, nil, nil
		}

	return report, chartBytes, nil
}

func generateProductivityChart(days []string, byDay map[string]*DayStat) ([]byte, error) {
    var series []chart.Series
    
   colors := []drawing.Color{
	drawing.ColorFromHex("2ecc71"),
	drawing.ColorFromHex("3498db"),
	drawing.ColorFromHex("f39c12"),
	drawing.ColorFromHex("95a5a6"),
	drawing.ColorFromHex("e74c3c"),
   }
   statusNames := []string{"Done", "In Progress", "Postponed", "New", "Deleted"}

    var xValues []float64
    var xLabels []string
    for i, day := range days {
        xValues = append(xValues, float64(i))
        t, _ := time.Parse("2006-01-02", day)
        xLabels = append(xLabels, t.Format("01-02"))
    }
    
    for statusIdx, statusName := range statusNames {
        var yValues []float64
        for _, day := range days {
			stats := byDay[day]
			var value float64
			switch  statusName {
			case "Done":
				value = float64(stats.Done)
			case "In Progress":
				value = float64(stats.Started)
			case "Postponed":
				value = float64(stats.Postponed)
			case "Deleted":
			value = float64(stats.Deleted)
			case "New":
				value = float64(stats.Total-stats.Done-stats.Postponed-stats.Deleted-stats.Started)
			}
            yValues = append(yValues, value)
        }
        
        series = append(series, chart.ContinuousSeries{
            Name:    statusName,
            Style:   chart.Style{StrokeColor: colors[statusIdx], FillColor: colors[statusIdx].WithAlpha(65), StrokeWidth: 3,},
            XValues: xValues,
            YValues: yValues,
        })
    }
    
	var ticks []chart.Tick
	for i, label := range xLabels {
		ticks = append(ticks, chart.Tick{
			Value: float64(i),
			Label: label,
		})
	}

    graph := chart.Chart{
        Title: "Weekly Task Overview",
        Background: chart.Style{
            Padding: chart.Box{Top: 40, Left: 20, Right: 20, Bottom: 20},
        },
        Series: series,
        XAxis: chart.XAxis{
            Ticks: ticks,
        },
        YAxis: chart.YAxis{
            Name: "Number of Tasks",
        },
    }
    
	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}