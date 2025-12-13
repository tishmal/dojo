// internal/adapters/postgres/task_repository.go
package postgres

import (
	"context"
	"dojo/internal/domain"
	"dojo/internal/ports"
	"time"

	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) ports.TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *TaskRepository) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
	var task domain.Task
	err := r.db.WithContext(ctx).First(&task, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrTaskNotFound
		}
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) GetByUserID(ctx context.Context, userID int64) ([]*domain.Task, error) {
	var tasks []*domain.Task
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&tasks).Error
	
	return tasks, err
}

func (r *TaskRepository) GetActiveByUserID(ctx context.Context, userID int64) ([]*domain.Task, error) {
	var tasks []*domain.Task
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("status IN ?", []domain.TaskStatus{
			domain.TaskStatusActive,
			domain.TaskStatusInProgress,
		}).
		Order("created_at DESC").
		Find(&tasks).Error
	
	return tasks, err
}

func (r *TaskRepository) GetDailyTasks(ctx context.Context, userID int64) ([]*domain.Task, error) {
	var tasks []*domain.Task
	
	today := time.Now().Truncate(24 * time.Hour)
	
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("frequency = ?", domain.FrequencyDaily).
		Where("created_at >= ?", today).
		Find(&tasks).Error
	
	return tasks, err
}

func (r *TaskRepository) Update(ctx context.Context, task *domain.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

func (r *TaskRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Task{}, id).Error
}

func (r *TaskRepository) GetUrgentTasks(ctx context.Context, userID int64) ([]*domain.Task, error) {
	var tasks []*domain.Task
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("is_urgent = ?", true).
		Where("status IN ?", []domain.TaskStatus{
			domain.TaskStatusActive,
			domain.TaskStatusInProgress,
		}).
		Order("urgent_until ASC").
		Find(&tasks).Error
	
	return tasks, err
}

func (r *TaskRepository) ExpireOldTasks(ctx context.Context) error {
	now := time.Now()
	
	result := r.db.WithContext(ctx).
		Model(&domain.Task{}).
		Where("is_urgent = ?", true).
		Where("urgent_until < ?", now).
		Where("status IN ?", []domain.TaskStatus{
			domain.TaskStatusActive,
			domain.TaskStatusInProgress,
		}).
		Update("status", domain.TaskStatusExpired)
	
	return result.Error
}