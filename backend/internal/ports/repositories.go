// internal/ports/repositories.go
package ports

import (
	"context"
	"dojo/internal/domain"
)

// UserRepository - интерфейс работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	List(ctx context.Context, limit, offset int) ([]*domain.User, error)
	GetInactivePlayers(ctx context.Context, limit int) ([]*domain.User, error)
	UpdateActivity(ctx context.Context, userID int64) error
}

// TaskRepository - интерфейс работы с заданиями
type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	GetByID(ctx context.Context, id int64) (*domain.Task, error)
	GetByUserID(ctx context.Context, userID int64) ([]*domain.Task, error)
	GetActiveByUserID(ctx context.Context, userID int64) ([]*domain.Task, error)
	GetDailyTasks(ctx context.Context, userID int64) ([]*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id int64) error
	GetUrgentTasks(ctx context.Context, userID int64) ([]*domain.Task, error)
	ExpireOldTasks(ctx context.Context) error
}

// AIService - интерфейс ИИ-сервиса
type AIService interface {
	AnalyzeTask(ctx context.Context, title, description string) (*TaskAnalysis, error)
	Chat(ctx context.Context, userID int64, message string, history []ChatMessage) (string, error)
	GenerateUrgentCall(ctx context.Context, userID int64) (*UrgentCallSuggestion, error)
}

// TaskAnalysis - результат анализа ИИ
type TaskAnalysis struct {
	TaskType   domain.TaskType
	Difficulty int
	XPReward   int
	GoldReward int
	EnergyCost int
	Explanation string
}

// ChatMessage - сообщение в чате
type ChatMessage struct {
	Role    string
	Content string
}

// UrgentCallSuggestion - предложение срочного вызова
type UrgentCallSuggestion struct {
	Title       string
	Description string
	TaskType    domain.TaskType
	Duration    int // минуты
}