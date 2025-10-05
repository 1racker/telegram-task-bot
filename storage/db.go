package storage

import (
	"log"
	"time"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Task struct {
	ID uint `gorm:"primaryKey"`
	UserID int64
	Title string
	Category string
	Priority int
	StartTime time.Time
	Duration int
	Status string
	CreatedAt time.Time
	UpdatedAt time.Time
	StartedAt *time.Time
	DoneAt *time.Time
	Postpones int 
}

func InitDB(path string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		log.Fatal("error connecting to db: %v", err)
	}

	if err := db.AutoMigrate(&Task{}); err != nil {
		log.Fatalf("failed to migrate db: %v", err)
	}
	
	return db
}