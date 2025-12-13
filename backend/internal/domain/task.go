// internal/domain/task.go
package domain

import "time"

type TaskStatus string

const (
	TaskStatusActive    TaskStatus = "active"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusExpired   TaskStatus = "expired"
)

type TaskFrequency string

const (
	FrequencyDaily  TaskFrequency = "daily"
	FrequencyCustom TaskFrequency = "custom"
	FrequencyUrgent TaskFrequency = "urgent"
)

type TaskType string

const (
	TypeStrength     TaskType = "strength"
	TypeAgility      TaskType = "agility"
	TypeIntelligence TaskType = "intelligence"
	TypeInsight      TaskType = "insight"
)

// Task - модель задания
type Task struct {
	ID          int64         `json:"id" gorm:"primaryKey"`
	UserID      int64         `json:"user_id" gorm:"index;not null"`
	
	Title       string        `json:"title" gorm:"not null"`
	Description string        `json:"description"`
	
	TaskType    TaskType      `json:"task_type" gorm:"not null"`
	Frequency   TaskFrequency `json:"frequency" gorm:"default:custom"`
	Status      TaskStatus    `json:"status" gorm:"default:active"`
	
	XPReward    int           `json:"xp_reward" gorm:"default:10"`
	GoldReward  int           `json:"gold_reward" gorm:"default:5"`
	StatBoost   int           `json:"stat_boost" gorm:"default:1"`
	
	EnergyCost  int           `json:"energy_cost" gorm:"default:10"`
	GoldCost    int           `json:"gold_cost" gorm:"default:0"`
	
	IsUrgent    bool          `json:"is_urgent" gorm:"default:false"`
	UrgentUntil *time.Time    `json:"urgent_until,omitempty"`
	Penalty     int           `json:"penalty" gorm:"default:0"`
	
	StartedAt   *time.Time    `json:"started_at,omitempty"`
	CompletedAt *time.Time    `json:"completed_at,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	
	AIAnalyzed  bool          `json:"ai_analyzed" gorm:"default:false"`
	AIDifficulty int          `json:"ai_difficulty" gorm:"default:1"`
}

// CanStart - проверяет можно ли начать
func (t *Task) CanStart() error {
	if t.Status != TaskStatusActive {
		return ErrTaskNotActive
	}
	
	if t.IsUrgent && t.UrgentUntil != nil {
		if time.Now().After(*t.UrgentUntil) {
			return ErrTaskExpired
		}
	}
	
	return nil
}

// Start - начинает задание
func (t *Task) Start() error {
	if err := t.CanStart(); err != nil {
		return err
	}
	
	now := time.Now()
	t.Status = TaskStatusInProgress
	t.StartedAt = &now
	return nil
}

// CanComplete - проверяет можно ли завершить
func (t *Task) CanComplete() error {
	if t.Status != TaskStatusInProgress {
		return ErrTaskNotInProgress
	}
	return nil
}

// Complete - завершает задание
func (t *Task) Complete() error {
	if err := t.CanComplete(); err != nil {
		return err
	}
	
	now := time.Now()
	t.Status = TaskStatusCompleted
	t.CompletedAt = &now
	return nil
}

// Fail - провал задания
func (t *Task) Fail() {
	t.Status = TaskStatusFailed
	now := time.Now()
	t.CompletedAt = &now
}

// Expire - истекло
func (t *Task) Expire() {
	t.Status = TaskStatusExpired
}

// IsExpired - проверка истечения
func (t *Task) IsExpired() bool {
	if !t.IsUrgent || t.UrgentUntil == nil {
		return false
	}
	return time.Now().After(*t.UrgentUntil)
}

// GetRewards - возвращает награды
func (t *Task) GetRewards() TaskRewards {
	return TaskRewards{
		XP:        t.XPReward,
		Gold:      t.GoldReward,
		StatBoost: t.StatBoost,
		StatType:  string(t.TaskType),
	}
}

// CalculateRewards - расчет наград по сложности
func (t *Task) CalculateRewards() {
	difficulty := t.AIDifficulty
	if difficulty == 0 {
		difficulty = 1
	}
	
	t.XPReward = 10 * difficulty
	t.GoldReward = 5 * difficulty
	t.StatBoost = 1
	
	if t.IsUrgent {
		t.XPReward = int(float64(t.XPReward) * 1.5)
		t.GoldReward = int(float64(t.GoldReward) * 1.5)
	}
}

// SetUrgent - делает срочным
func (t *Task) SetUrgent(duration time.Duration, penalty int) {
	t.IsUrgent = true
	t.Frequency = FrequencyUrgent
	urgentUntil := time.Now().Add(duration)
	t.UrgentUntil = &urgentUntil
	t.Penalty = penalty
	
	t.CalculateRewards()
}

// GetTimeRemaining - оставшееся время
func (t *Task) GetTimeRemaining() time.Duration {
	if !t.IsUrgent || t.UrgentUntil == nil {
		return 0
	}
	
	remaining := time.Until(*t.UrgentUntil)
	if remaining < 0 {
		return 0
	}
	return remaining
}

type TaskRewards struct {
	XP        int    `json:"xp"`
	Gold      int    `json:"gold"`
	StatBoost int    `json:"stat_boost"`
	StatType  string `json:"stat_type"`
}

// Конструкторы
func NewDailyTask(userID int64, title string, taskType TaskType) *Task {
	return &Task{
		UserID:     userID,
		Title:      title,
		TaskType:   taskType,
		Frequency:  FrequencyDaily,
		Status:     TaskStatusActive,
		XPReward:   15,
		GoldReward: 5,
		StatBoost:  1,
		EnergyCost: 10,
	}
}

func NewCustomTask(userID int64, title, description string, taskType TaskType) *Task {
	return &Task{
		UserID:      userID,
		Title:       title,
		Description: description,
		TaskType:    taskType,
		Frequency:   FrequencyCustom,
		Status:      TaskStatusActive,
		EnergyCost:  10,
		GoldCost:    10,
		AIAnalyzed:  false,
	}
}

func NewUrgentCall(userID int64, title string, taskType TaskType, duration time.Duration) *Task {
	task := &Task{
		UserID:     userID,
		Title:      title,
		TaskType:   taskType,
		Frequency:  FrequencyUrgent,
		Status:     TaskStatusActive,
		EnergyCost: 15,
	}
	
	task.SetUrgent(duration, 5)
	return task
}