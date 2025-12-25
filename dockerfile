# ---------- BUILD STAGE ----------
FROM golang:1.25.5-alpine AS builder

WORKDIR /dynamic-pricing-tool-ru

# Копируем весь проект
COPY . .

# Приводим модули в порядок
RUN go mod tidy

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./command/server

# ---------- FINAL STAGE ----------
FROM alpine:3.18

WORKDIR /app

# Копируем бинарник
COPY --from=builder /build/app .

# Запуск
CMD ["./app"]
