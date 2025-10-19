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
	MarkDone(id uint) error
	Postpone(id uint, newTime time.Time) error
	Delete(id uint) error
	GetStats(userID int64) (total, done int64, err error)
}

type GormTaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &GormTaskRepository{db: db}
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
 userID, startOfDay, endOfDy).Order("priority ASC").Find(&tasks).Error
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

func (r *GormTaskRepository) MarkDone(id uint) error {
	return r.db.Model(&Task{}).Where("id = ?", id).Update("status", "done").Error
}

func(r *GormTaskRepository) Postpone(id uint, newTime time.Time) error {
	return r.db.Model(&Task{}).Where("id = ?", id).
	Updates(map[string]interface{}{"status": "postponed", "start_time": newTime}).Error
}

func(r *GormTaskRepository) Delete(id uint) error {
	return r.db.Delete(&Task{}, id). Error
}

func(r *GormTaskRepository) GetStats (userID int64) (total, done int64, err error) {
	err = r.db.Model(&Task{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return
	}
	err = r.db.Model(&Task{}).Where("user_id = ? AND status = ?", userID, "done").Count(&done).Error
	return
}
