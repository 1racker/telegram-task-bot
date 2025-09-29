package storage

import "time"

type User struct {
	ID uint `gorm:"primaryKey"`
	TelegramID int64 `gorm:"uniqueIndex"`
	ReminderTime string
}

type Task struct {
	ID uint `gorm:"primaryKey"`
	UserID uint
	Title string
	Category string
	Priority int
	Date time.Time
	StartTime string
	Duration int
	Status string
}