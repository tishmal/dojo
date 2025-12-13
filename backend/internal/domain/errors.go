// internal/domain/errors.go
package domain

import "errors"

// Ошибки пользователя
var (
	ErrUserNotFound = errors.New("пользователь не найден")
	ErrUserAlreadyExists = errors.New("пользователь уже существует")
	ErrInsufficientGold = errors.New("недостаточно золота")
	ErrInsufficientEnergy = errors.New("недостаточно энергии")
	ErrNoSenseiRequests = errors.New("закончились запросы к Сенсею")
	ErrLicenseInactive = errors.New("лицензия охотника неактивна")
)

// Ошибки заданий
var (
	ErrTaskNotFound = errors.New("задание не найдено")
	ErrTaskNotActive = errors.New("задание не активно")
	ErrTaskNotInProgress = errors.New("задание не начато")
	ErrTaskExpired = errors.New("срок задания истек")
	ErrTaskAlreadyStarted = errors.New("задание уже начато")
)

// Ошибки авторизации
var (
	ErrInvalidTelegramData = errors.New("невалидные данные авторизации")
	ErrUnauthorized = errors.New("требуется авторизация")
)

// Ошибки рейдов
var (
	ErrRaidNotFound = errors.New("рейд не найден")
	ErrCannotRaidSelf = errors.New("нельзя рейдить самого себя")
	ErrPlayerNotInactive = errors.New("игрок активен")
)

// Ошибки ИИ
var (
	ErrAIServiceUnavailable = errors.New("ИИ-сервис недоступен")
	ErrAIAnalysisFailed = errors.New("не удалось проанализировать задание")
)