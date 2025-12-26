# Stage 1: Build
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# Устанавливаем минимальные зависимости
RUN apk add --no-cache git ca-certificates tzdata

# Копируем модули для кэширования
COPY go.mod go.sum ./
RUN go mod download

# Копируем все исходники
COPY . .

# Собираем приложение из cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/pricing-tool ./cmd/server/

# Stage 2: Runtime
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata curl

WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /app/pricing-tool .
COPY --from=builder /app/cmd/server/.env.example .env.example

# Создаем директории для логов
RUN mkdir -p /app/logs

# Создаем непривилегированного пользователя
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Экспортируем порт приложения
EXPOSE 5004

# Запускаем приложение
CMD ["./pricing-tool"]