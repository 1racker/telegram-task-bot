package storage

import (
	"log"
	"telegram-task-bot/config"

	"github.com/daixiang0/gci/pkg/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatal("error connecting to db:", err)
	}
	db.AutoMigrate(&User{}, &Task{})
	return db
}