.PHONY: help dev prod stop clean logs build test migrate

# Цвета для вывода
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
RESET  := $(shell tput -Txterm sgr0)

help: ## Показать помощь
	@echo '${GREEN}Доступные команды:${RESET}'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  ${YELLOW}%-15s${RESET} %s\n", $$1, $$2}'

dev: ## Запустить в режиме разработки (hot-reload)
	@echo "${GREEN}🚀 Запуск в режиме разработки...${RESET}"
	docker-compose -f docker-compose.dev.yml up --build

prod: ## Запустить в режиме production
	@echo "${GREEN}🚀 Запуск в production режиме...${RESET}"
	docker-compose --profile prod up -d --build

stop: ## Остановить все контейнеры
	@echo "${YELLOW}⏸️  Остановка контейнеров...${RESET}"
	docker-compose down
	docker-compose -f docker-compose.dev.yml down

clean: ## Удалить контейнеры и volumes
	@echo "${YELLOW}🧹 Очистка...${RESET}"
	docker-compose down -v
	docker-compose -f docker-compose.dev.yml down -v
	rm -rf backend/tmp

logs: ## Показать логи
	docker-compose logs -f

logs-api: ## Показать логи только API
	docker-compose logs -f api-dev

logs-db: ## Показать логи только БД
	docker-compose logs -f postgres

build: ## Пересобрать образы
	@echo "${GREEN}🔨 Сборка образов...${RESET}"
	docker-compose build

test: ## Запустить тесты
	@echo "${GREEN}🧪 Запуск тестов...${RESET}"
	cd backend && go test -v ./...

migrate: ## Запустить миграции вручную
	@echo "${GREEN}📦 Запуск миграций...${RESET}"
	docker-compose exec api-dev go run cmd/api/main.go migrate

psql: ## Подключиться к PostgreSQL
	docker-compose exec postgres psql -U dojo -d dojo

shell-api: ## Зайти в контейнер API
	docker-compose exec api-dev sh

shell-db: ## Зайти в контейнер БД
	docker-compose exec postgres sh

setup: ## Первая настройка проекта
	@echo "${GREEN}⚙️  Настройка проекта...${RESET}"
	cp .env.example .env
	@echo "${YELLOW}Отредактируй файл .env перед запуском!${RESET}"

install: setup ## Установка зависимостей
	@echo "${GREEN}📦 Установка зависимостей...${RESET}"
	cd backend && go mod download
	cd frontend && npm install