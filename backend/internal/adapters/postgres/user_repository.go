// internal/adapters/postgres/user_repository.go
package postgres

import (
	"context"
	"dojo/internal/domain"
	"dojo/internal/ports"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Where("telegram_id = ?", telegramID).
		First(&user).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("level DESC, xp DESC").
		Find(&users).Error
	
	return users, err
}

func (r *UserRepository) GetInactivePlayers(ctx context.Context, limit int) ([]*domain.User, error) {
	var users []*domain.User
	
	inactiveThreshold := time.Now().AddDate(0, 0, -3)
	
	err := r.db.WithContext(ctx).
		Where("last_active_at < ?", inactiveThreshold).
		Where("gold > ?", 10).
		Order("gold DESC").
		Limit(limit).
		Find(&users).Error
	
	return users, err
}

func (r *UserRepository) UpdateActivity(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", userID).
		Update("last_active_at", time.Now()).Error
}