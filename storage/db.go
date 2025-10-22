package storage

import (
	"log"
	"time"
	"gorm.io/driver/postgres"
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
var DB *gorm.DB

func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("error connecting to PostgreSQL: %v", err)
	}

	if err := db.AutoMigrate(&Task{}); err != nil {
		log.Fatalf("auto migration failed: %v", err)
	}
	
	log.Println("Connected to PostgreSQL succesfully")
	DB = db
	return db
}