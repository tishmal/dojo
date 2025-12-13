# ============================================
# Стадия 1: Базовый образ для зависимостей
# ============================================
FROM golang:1.21-alpine AS base

# Устанавливаем необходимые пакеты
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Копируем go.mod и go.sum для кэширования слоя
COPY go.mod go.sum ./
RUN go mod download

# ============================================
# Стадия 2: Разработка (с hot-reload)
# ============================================
FROM base AS development

# Устанавливаем Air для hot-reload
RUN go install github.com/cosmtrek/air@latest

# Копируем конфиг Air
COPY .air.toml ./

# Копируем весь код
COPY . .

# Открываем порт
EXPOSE 8080

# Запускаем через Air (hot-reload)
CMD ["air", "-c", ".air.toml"]

# ============================================
# Стадия 3: Сборка production бинарника
# ============================================
FROM base AS builder

# Копируем весь код
COPY . .

# Компилируем с оптимизациями
# -ldflags="-s -w" уменьшает размер бинарника
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o main \
    ./cmd/api/main.go

# ============================================
# Стадия 4: Production образ (минимальный)
# ============================================
FROM alpine:latest AS production

# Устанавливаем CA сертификаты для HTTPS
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем скомпилированный бинарник
COPY --from=builder /app/main .

# Создаем непривилегированного пользователя
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /root

USER appuser

EXPOSE 8080

# Healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./main"]