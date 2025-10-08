package storage

import (
	"time"

	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *Task) error
	GetByID(id uint) (*Task, error)
	Update(task *Task) error
	GetTodayTasks(userID int64) ([]Task, error)
	GetWeeklyTasks(userID int64, from, to time.Time) ([]Task, error)
	GetDistinctUserIDs() ([]int64, error)
}

type GormTaskRepository struct {
	db *gorm.DB
}

func (r *GormTaskRepository) Create(task *Task) error {
	return r.db.Create(task).Error
}

func (r *GormTaskRepository) GetByID(id uint) (*Task, error) {
	var task Task
	err := r.db.First(&task, id).Error
	return &task, err
}

func (r *GormTaskRepository) Update(task *Task) error {
	return r.db.Save(task).Error
}

func (r *GormTaskRepository) GetTodayTasks(userID int64) ([]Task, error) {
	var tasks []Task
	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDy := startOfDay.Add(24 * time.Hour)
	err := r.db.Where("user_id = ? AND start_time >= ? AND start_time < ?",
 userID, startOfDay, endOfDy).Order("priority as c").Find(&tasks).Error
 return tasks, err
}

func (r *GormTaskRepository) GetWeeklyTasks(userID int64, from, to time.Time) ([]Task, error) {
	var tasks []Task
	err := r.db.Where("user_id = ? AND created_at <= ?",
	userID, from, to).Find(&tasks).Error
 return tasks, err
}

func (r *GormTaskRepository) GetDistinctUserIDs() ([]int64, error) {
	var userIDs []int64
	err := r.db.Model(&Task{}).Distinct("user_id").Pluck("user_id", &userIDs).Error
	return userIDs, err
}