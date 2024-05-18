package db

import (
	"context"
	"fmt"
	"time"

	"github.com/Admiral-Simo/task/models"
	"gorm.io/gorm"
)

type TaskStore interface {
	CreateTask(ctx context.Context, task *models.Task) error
	GetAllTasks(ctx context.Context) ([]*models.Task, error)
	GetAllTodaysTasks(ctx context.Context) ([]*models.Task, error)
	// UpdateTask(ctx context.Context, id int, title string) error
	// DeleteTask(ctx context.Context, id int) error
	MarkDoneTask(ctx context.Context, id int64) error
	MarkUnDoneTask(ctx context.Context, id int64) error
	SetPriority(ctx context.Context, id int64, priority string) error
}

type GormTaskStore struct {
	db *gorm.DB
}

func NewGormTaskStore(db *gorm.DB) *GormTaskStore {
	return &GormTaskStore{
		db: db,
	}
}

func (s *GormTaskStore) GetAllTasks(ctx context.Context) ([]*models.Task, error) {
	var tasks []*models.Task
	if err := s.db.WithContext(ctx).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *GormTaskStore) CreateTask(ctx context.Context, task *models.Task) error {
	if err := s.db.WithContext(ctx).Create(task).Error; err != nil {
		return err
	}
	return nil
}

func (s *GormTaskStore) GetAllTodaysTasks(ctx context.Context) ([]*models.Task, error) {
	var tasks []*models.Task
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	if err := s.db.WithContext(ctx).Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *GormTaskStore) MarkDoneTask(ctx context.Context, id int64) error {
	result := s.db.WithContext(ctx).Model(models.Task{}).Where("id = ?", id).Updates(map[string]interface{}{"done": true})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("task with ID %d not found", id)
	}
	return nil
}

func (s *GormTaskStore) MarkUnDoneTask(ctx context.Context, id int64) error {
	result := s.db.WithContext(ctx).Model(models.Task{}).Where("id = ?", id).Updates(map[string]interface{}{"done": false})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("task with ID %d not found", id)
	}
	return nil
}

func (s *GormTaskStore) SetPriority(ctx context.Context, id int64, priority string) error {
	// checking wheter the priority argument is correct or not
	if priority != "H" && priority != "M" && priority != "L" {
		return fmt.Errorf("only values to be identified as priorities is 'H': high 'M': medium 'L': low")
	}
	result := s.db.WithContext(ctx).Model(models.Task{}).Where("id = ?", id).Updates(map[string]interface{}{"priority": priority})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("task with ID %d not found", id)
	}
	return nil
}
