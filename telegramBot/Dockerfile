FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# Копируем файлы модулей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем приложение
RUN go build -o telegram-bot .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из builder stage
COPY --from=builder /app/telegram-bot .

# Копируем .env файл (опционально)
COPY --from=builder /app/.env ./

# Создаем не-root пользователя для безопасности
RUN adduser -D -s /bin/sh appuser
USER appuser

CMD ["./telegram-bot"]