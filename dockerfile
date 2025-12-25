# ---------- BUILD STAGE ----------
FROM golang:1.25.5-alpine AS builder

# Рабочая директория внутри контейнера
WORKDIR /dynamic-pricing-tool-ru

# Копируем файлы модулей и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod tidy

# Копируем весь проект
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dynamic-pricing-tool-ru .

# ---------- FINAL STAGE ----------
FROM alpine:3.18

# Рабочая директория для финального образа
WORKDIR /dynamic-pricing-tool-ru

# Копируем бинарник из билдера
COPY --from=builder /dynamic-pricing-tool-ru/dynamic-pricing-tool-ru .

# Команда запуска
CMD ["./dynamic-pricing-tool-ru"]
