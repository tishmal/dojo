// internal/domain/user.go
package domain

import "time"

// User - основная модель пользователя в системе Додзё
type User struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	TelegramID   int64     `json:"telegram_id" gorm:"uniqueIndex;not null"`
	Username     string    `json:"username"`
	FirstName    string    `json:"first_name"`
	PhotoURL     string    `json:"photo_url"`
	
	// Прогресс и уровень
	Level        int       `json:"level" gorm:"default:1"`
	XP           int       `json:"xp" gorm:"default:0"`
	XPToNextLvl  int       `json:"xp_to_next_level"`
	
	// Ресурсы
	Gold         int       `json:"gold" gorm:"default:0"`
	Energy       int       `json:"energy" gorm:"default:100"`
	MaxEnergy    int       `json:"max_energy" gorm:"default:100"`
	
	// Характеристики
	Strength     int       `json:"strength" gorm:"default:1"`
	Agility      int       `json:"agility" gorm:"default:1"`
	Intelligence int       `json:"intelligence" gorm:"default:1"`
	Insight      int       `json:"insight" gorm:"default:1"`
	
	// Система лицензии
	LicenseActive     bool      `json:"license_active" gorm:"default:true"`
	LastLicenseCheck  time.Time `json:"last_license_check"`
	LicenseExpiresAt  time.Time `json:"license_expires_at"`
	
	// Система Сенсея
	SenseiRequests    int       `json:"sensei_requests" gorm:"default:5"`
	SenseiResetsAt    time.Time `json:"sensei_resets_at"`
	
	// Метаданные
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	LastActiveAt time.Time `json:"last_active_at"`
}

// CalculateXPToNextLevel - расчет XP для следующего уровня
func (u *User) CalculateXPToNextLevel() int {
	baseXP := 100
	return baseXP * u.Level * u.Level / 2
}

// AddXP - добавляет опыт, возвращает true если level up
func (u *User) AddXP(xp int) bool {
	u.XP += xp
	u.XPToNextLvl = u.CalculateXPToNextLevel()
	
	if u.XP >= u.XPToNextLvl {
		u.LevelUp()
		return true
	}
	return false
}

// LevelUp - повышает уровень
func (u *User) LevelUp() {
	u.Level++
	u.XP = u.XP - u.XPToNextLvl
	u.XPToNextLvl = u.CalculateXPToNextLevel()
	
	u.MaxEnergy += 10
	u.Energy = u.MaxEnergy
	u.Gold += u.Level * 10
}

// AddGold - добавляет золото
func (u *User) AddGold(amount int) {
	u.Gold += amount
	if u.Gold < 0 {
		u.Gold = 0
	}
}

// SpendGold - тратит золото
func (u *User) SpendGold(amount int) error {
	if u.Gold < amount {
		return ErrInsufficientGold
	}
	u.Gold -= amount
	return nil
}

// HasEnoughEnergy - проверяет энергию
func (u *User) HasEnoughEnergy(required int) bool {
	return u.Energy >= required
}

// SpendEnergy - тратит энергию
func (u *User) SpendEnergy(amount int) error {
	if !u.HasEnoughEnergy(amount) {
		return ErrInsufficientEnergy
	}
	u.Energy -= amount
	return nil
}

// RestoreEnergy - восстанавливает энергию
func (u *User) RestoreEnergy(amount int) {
	u.Energy += amount
	if u.Energy > u.MaxEnergy {
		u.Energy = u.MaxEnergy
	}
}

// IsLicenseValid - проверяет лицензию
func (u *User) IsLicenseValid() bool {
	return u.LicenseActive && time.Now().Before(u.LicenseExpiresAt)
}

// RenewLicense - продлевает лицензию
func (u *User) RenewLicense() {
	u.LicenseActive = true
	u.LastLicenseCheck = time.Now()
	u.LicenseExpiresAt = time.Now().AddDate(0, 0, 7)
}

// RevokeLicense - отзывает лицензию
func (u *User) RevokeLicense() {
	u.LicenseActive = false
}

// CanUseSensei - проверяет запросы к Сенсею
func (u *User) CanUseSensei() bool {
	if u.SenseiRequests <= 0 {
		if time.Now().After(u.SenseiResetsAt) {
			u.SenseiRequests = 5
			u.SenseiResetsAt = time.Now().AddDate(0, 0, 7)
		}
	}
	return u.SenseiRequests > 0
}

// UseSenseiRequest - использует запрос
func (u *User) UseSenseiRequest() error {
	if !u.CanUseSensei() {
		return ErrNoSenseiRequests
	}
	u.SenseiRequests--
	return nil
}

// BuySenseiRequest - покупка запроса
func (u *User) BuySenseiRequest(cost int) error {
	if err := u.SpendGold(cost); err != nil {
		return err
	}
	u.SenseiRequests++
	return nil
}

// IncreaseAttribute - увеличивает характеристику
func (u *User) IncreaseAttribute(taskType string, amount int) {
	switch taskType {
	case "strength":
		u.Strength += amount
	case "agility":
		u.Agility += amount
	case "intelligence":
		u.Intelligence += amount
	case "insight":
		u.Insight += amount
	}
}

// GetRank - возвращает ранг
func (u *User) GetRank() string {
	switch {
	case u.Level >= 50:
		return "S"
	case u.Level >= 40:
		return "A"
	case u.Level >= 30:
		return "B"
	case u.Level >= 20:
		return "C"
	case u.Level >= 10:
		return "D"
	default:
		return "E"
	}
}

// IsInactive - проверяет неактивность (для рейдов)
func (u *User) IsInactive() bool {
	return time.Since(u.LastActiveAt) > 3*24*time.Hour
}

// UpdateActivity - обновляет активность
func (u *User) UpdateActivity() {
	u.LastActiveAt = time.Now()
}