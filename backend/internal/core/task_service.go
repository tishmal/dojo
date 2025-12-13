// internal/core/task_service.go
package core

import (
	"context"
	"dojo/internal/domain"
	"dojo/internal/ports"
)

type TaskService struct {
	taskRepo  ports.TaskRepository
	userRepo  ports.UserRepository
	aiService ports.AIService
}

func NewTaskService(
	taskRepo ports.TaskRepository,
	userRepo ports.UserRepository,
	aiService ports.AIService,
) *TaskService {
	return &TaskService{
		taskRepo:  taskRepo,
		userRepo:  userRepo,
		aiService: aiService,
	}
}

// GetActiveTasks - получить активные задания
func (s *TaskService) GetActiveTasks(ctx context.Context, userID int64) ([]*domain.Task, error) {
	s.userRepo.UpdateActivity(ctx, userID)
	return s.taskRepo.GetActiveByUserID(ctx, userID)
}

// CreateCustomTask - создать пользовательское задание
func (s *TaskService) CreateCustomTask(ctx context.Context, userID int64, title, description string, taskType domain.TaskType) (*domain.Task, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	task := domain.NewCustomTask(userID, title, description, taskType)
	
	// Списываем золото за создание
	if err := user.SpendGold(task.GoldCost); err != nil {
		return nil, err
	}
	
	// Анализируем через ИИ
	if s.aiService != nil {
		analysis, err := s.aiService.AnalyzeTask(ctx, title, description)
		if err == nil {
			task.TaskType = analysis.TaskType
			task.AIDifficulty = analysis.Difficulty
			task.XPReward = analysis.XPReward
			task.GoldReward = analysis.GoldReward
			task.EnergyCost = analysis.EnergyCost
			task.AIAnalyzed = true
		} else {
			task.CalculateRewards()
		}
	} else {
		task.CalculateRewards()
	}
	
	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}
	
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	
	return task, nil
}

// StartTask - начать задание
func (s *TaskService) StartTask(ctx context.Context, taskID, userID int64) error {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}
	
	if task.UserID != userID {
		return domain.ErrUnauthorized
	}
	
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	
	if err := user.SpendEnergy(task.EnergyCost); err != nil {
		return err
	}
	
	if err := task.Start(); err != nil {
		return err
	}
	
	if err := s.taskRepo.Update(ctx, task); err != nil {
		return err
	}
	
	return s.userRepo.Update(ctx, user)
}

// CompleteTask - завершить задание
func (s *TaskService) CompleteTask(ctx context.Context, taskID, userID int64) (*TaskCompletionResult, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	
	if task.UserID != userID {
		return nil, domain.ErrUnauthorized
	}
	
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	// Штраф за отсутствие лицензии
	if task.Frequency != domain.FrequencyDaily && !user.IsLicenseValid() {
		task.XPReward = task.XPReward / 2
		task.GoldReward = task.GoldReward / 2
	}
	
	if err := task.Complete(); err != nil {
		return nil, err
	}
	
	leveledUp := user.AddXP(task.XPReward)
	user.AddGold(task.GoldReward)
	user.IncreaseAttribute(string(task.TaskType), task.StatBoost)
	
	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}
	
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	
	result := &TaskCompletionResult{
		Task:      task,
		LeveledUp: leveledUp,
		NewLevel:  user.Level,
		Rewards:   task.GetRewards(),
	}
	
	return result, nil
}

// DeclineUrgentCall - отказ от срочного вызова
func (s *TaskService) DeclineUrgentCall(ctx context.Context, taskID, userID int64) error {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}
	
	if task.UserID != userID {
		return domain.ErrUnauthorized
	}
	
	if !task.IsUrgent {
		return domain.ErrTaskNotActive
	}
	
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	
	user.AddGold(-task.Penalty)
	task.Fail()
	
	if err := s.taskRepo.Update(ctx, task); err != nil {
		return err
	}
	
	return s.userRepo.Update(ctx, user)
}

type TaskCompletionResult struct {
	Task      *domain.Task
	LeveledUp bool
	NewLevel  int
	Rewards   domain.TaskRewards
}