# ---------- build stage ----------
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# зависимости
COPY go.mod go.sum ./
RUN go mod download

# исходники
COPY . .

# сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o pricing-tool-ru cmd/server/main.go


# ---------- runtime stage ----------
FROM alpine:latest

WORKDIR /app

# сертификаты (на всякий случай)
RUN apk add --no-cache ca-certificates

# бинарник
COPY --from=builder /app/pricing-tool-ru .

# .env будет монтироваться через docker-compose
EXPOSE 5004

CMD ["./pricing-tool-ru"]
