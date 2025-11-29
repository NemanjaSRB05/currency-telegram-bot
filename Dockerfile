
# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /bot ./cmd/bot

# Run stage  
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем бинарник из builder stage
COPY --from=builder /bot .

# Копируем миграции
COPY --from=builder /app/migrations ./migrations

# Создаем не-root пользователя для безопасности
RUN adduser -D -s /bin/sh appuser
USER appuser

EXPOSE 8080

# Используем ENTRYPOINT для поддержки команд
ENTRYPOINT ["./bot"]
CMD []