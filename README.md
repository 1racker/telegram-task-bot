# Telegram Task Manager Bot (Go)

A simple PET-project in Go: Telegram bot for task management and planning

## Features
-Add tasks with start time and duration
-Inline-buttons for marking statuses ("Started", "Completed", "Postponed", "Declined")
-Daily reminders via cron (the user sets the notification time)
-Automatic generation of a weekly report:
1.Amount of completed,postponed and rejected tasks
2.Percentage of completed tasks
3.Determination of the most productive and least productive days
4.(Optional:simple text diagram in the message)

## Technologies
- Go (golang)
- Libraries [telebot](https://github.com/tucnak/telebot) for Telegram API  
- [robfig/cron](https://github.com/robfig/cron) for reminders

## Launch
1.Download Go `>=1.21`
2.Clone repository
3.Create `.env` with token of a bot:
```env
TELEGRAM_TOKEN=your_token
4. go run main.go
