// cmd/api/main.go
package main

import (
	"log"
	"os"

	"dojo/internal/adapters/postgres"
	"dojo/internal/core"
	"dojo/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	postgresGorm "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Получаем переменные окружения
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://dojo:dev_password@localhost:5432/dojo?sslmode=disable"
	}
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	// Подключаемся к PostgreSQL
	db, err := gorm.Open(postgresGorm.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}
	
	// Автомиграция (создает таблицы если их нет)
	log.Println("Запуск автомиграции...")
	if err := db.AutoMigrate(&domain.User{}, &domain.Task{}); err != nil {
		log.Fatal("Ошибка миграции:", err)
	}
	log.Println("Миграция завершена!")
	
	// Инициализируем репозитории
	userRepo := postgres.NewUserRepository(db)
	taskRepo := postgres.NewTaskRepository(db)
	
	// Инициализируем сервисы (без ИИ пока)
	userService := core.NewUserService(userRepo, nil)
	taskService := core.NewTaskService(taskRepo, userRepo, nil)
	
	// Создаем Fiber приложение
	app := fiber.New(fiber.Config{
		AppName: "Dojo API v1.0",
	})
	
	// Middleware
	app.Use(logger.New()) // Логирование запросов
	app.Use(cors.New())   // CORS для фронтенда
	
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"app": "Dojo API",
		})
	})
	
	// API роуты
	api := app.Group("/api")
	
	// Временный эндпоинт для создания тестового пользователя
	api.Post("/auth/test", func(c *fiber.Ctx) error {
		user, err := userService.GetOrCreateUser(
			c.Context(),
			12345678, // Тестовый Telegram ID
			"test_user",
			"Test",
			"",
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		
		return c.JSON(fiber.Map{
			"user": user,
			"token": "test_token_123", // В проде будет JWT
		})
	})
	
	// Middleware для проверки авторизации (упрощенная версия)
	authMiddleware := func(c *fiber.Ctx) error {
		// TODO: Здесь будет проверка Telegram WebApp данных
		// Пока просто ставим тестовый ID
		c.Locals("user_id", int64(1))
		return c.Next()
	}
	
	// Защищенные роуты
	protected := api.Group("", authMiddleware)
	
	// Профиль
	protected.Get("/profile", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(int64)
		user, err := userService.GetProfile(c.Context(), userID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(user)
	})
	
	// Задания
	protected.Get("/tasks", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(int64)
		tasks, err := taskService.GetActiveTasks(c.Context(), userID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"tasks": tasks})
	})
	
	protected.Post("/tasks", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(int64)
		
		var req struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			TaskType    string `json:"task_type"`
		}
		
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Неверный формат"})
		}
		
		task, err := taskService.CreateCustomTask(
			c.Context(),
			userID,
			req.Title,
			req.Description,
			domain.TaskType(req.TaskType),
		)
		
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		
		return c.Status(201).JSON(task)
	})
	
	protected.Post("/tasks/:id/start", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(int64)
		taskID, _ := c.ParamsInt("id")
		
		err := taskService.StartTask(c.Context(), int64(taskID), userID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		
		return c.JSON(fiber.Map{"success": true})
	})
	
	protected.Post("/tasks/:id/complete", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(int64)
		taskID, _ := c.ParamsInt("id")
		
		result, err := taskService.CompleteTask(c.Context(), int64(taskID), userID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		
		return c.JSON(result)
	})
	
	// Запускаем сервер
	log.Printf("🚀 Сервер запущен на порту %s", port)
	log.Fatal(app.Listen(":" + port))
}