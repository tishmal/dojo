// internal/core/user_service.go
package core

import (
	"context"
	"dojo/internal/domain"
	"dojo/internal/ports"
	"time"
)

type UserService struct {
	userRepo  ports.UserRepository
	aiService ports.AIService
}

func NewUserService(userRepo ports.UserRepository, aiService ports.AIService) *UserService {
	return &UserService{
		userRepo:  userRepo,
		aiService: aiService,
	}
}

// GetOrCreateUser - получить или создать пользователя
func (s *UserService) GetOrCreateUser(ctx context.Context, telegramID int64, username, firstName, photoURL string) (*domain.User, error) {
	user, err := s.userRepo.GetByTelegramID(ctx, telegramID)
	if err == nil {
		user.UpdateActivity()
		user.Username = username
		user.FirstName = firstName
		user.PhotoURL = photoURL
		
		if err := s.userRepo.Update(ctx, user); err != nil {
			return nil, err
		}
		
		return user, nil
	}
	
	if err == domain.ErrUserNotFound {
		user = &domain.User{
			TelegramID:       telegramID,
			Username:         username,
			FirstName:        firstName,
			PhotoURL:         photoURL,
			Level:            1,
			XP:               0,
			Gold:             100,
			Energy:           100,
			MaxEnergy:        100,
			Strength:         1,
			Agility:          1,
			Intelligence:     1,
			Insight:          1,
			LicenseActive:    true,
			SenseiRequests:   5,
			LastActiveAt:     time.Now(),
		}
		
		user.XPToNextLvl = user.CalculateXPToNextLevel()
		user.RenewLicense()
		user.SenseiResetsAt = time.Now().AddDate(0, 0, 7)
		
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
		
		return user, nil
	}
	
	return nil, err
}

// GetProfile - получить профиль
func (s *UserService) GetProfile(ctx context.Context, userID int64) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	user.UpdateActivity()
	s.userRepo.Update(ctx, user)
	
	return user, nil
}

// RenewLicense - продлить лицензию
func (s *UserService) RenewLicense(ctx context.Context, userID int64) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	
	user.RenewLicense()
	return s.userRepo.Update(ctx, user)
}

// ChatWithSensei - чат с Сенсеем
func (s *UserService) ChatWithSensei(ctx context.Context, userID int64, message string, history []ports.ChatMessage) (string, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}
	
	if err := user.UseSenseiRequest(); err != nil {
		return "", err
	}
	
	response, err := s.aiService.Chat(ctx, userID, message, history)
	if err != nil {
		user.SenseiRequests++
		s.userRepo.Update(ctx, user)
		return "", err
	}
	
	if err := s.userRepo.Update(ctx, user); err != nil {
		return "", err
	}
	
	return response, nil
}

// RaidInactivePlayer - рейд неактивного игрока
func (s *UserService) RaidInactivePlayer(ctx context.Context, attackerID, targetID int64, cost int) (*RaidResult, error) {
	attacker, err := s.userRepo.GetByID(ctx, attackerID)
	if err != nil {
		return nil, err
	}
	
	target, err := s.userRepo.GetByID(ctx, targetID)
	if err != nil {
		return nil, err
	}
	
	if attackerID == targetID {
		return nil, domain.ErrCannotRaidSelf
	}
	
	if !target.IsInactive() {
		return nil, domain.ErrPlayerNotInactive
	}
	
	if !attacker.IsLicenseValid() {
		return nil, domain.ErrLicenseInactive
	}
	
	if err := attacker.SpendGold(cost); err != nil {
		return nil, err
	}
	
	loot := int(float64(target.Gold) * 0.2)
	if loot < 10 {
		loot = 10
	}
	
	target.AddGold(-loot)
	attacker.AddGold(loot)
	
	bonusXP := target.Level * 5
	attacker.AddXP(bonusXP)
	
	if err := s.userRepo.Update(ctx, attacker); err != nil {
		return nil, err
	}
	
	if err := s.userRepo.Update(ctx, target); err != nil {
		return nil, err
	}
	
	result := &RaidResult{
		GoldLooted: loot,
		XPGained:   bonusXP,
		TargetName: target.Username,
		TargetRank: target.GetRank(),
	}
	
	return result, nil
}

type RaidResult struct {
	GoldLooted int    `json:"gold_looted"`
	XPGained   int    `json:"xp_gained"`
	TargetName string `json:"target_name"`
	TargetRank string `json:"target_rank"`
}