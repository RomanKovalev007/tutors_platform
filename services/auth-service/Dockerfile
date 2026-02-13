FROM golang:1.24.3-alpine AS builder

RUN apk add --no-cache git

ENV GOPRIVATE=gitlab.crja72.ru \
    GOSUMDB=off

WORKDIR /app

# Копируем модуль api в корень проекта (чтобы ../../api работал)
COPY api/ ./api/

# Создаем структуру каталогов и копируем go.mod и go.sum auth-service
COPY services/auth-service/go.mod services/auth-service/go.sum ./services/auth-service/

# Переходим в рабочую директорию сервиса
WORKDIR /app/services/auth-service

# Теперь go mod download видит ../../api/go.mod и использует локальный replace
RUN go mod download

# Копируем код самого auth-service
COPY services/auth-service/ ./

# Билдим
RUN go build -o /app/main ./cmd/app

# Финальный образ
FROM alpine:latest

# Устанавливаем только netcat, ca-certificates копируем из builder
RUN apk add --no-cache netcat-openbsd 2>/dev/null || apk add --no-cache busybox-extras

WORKDIR /app

# Копируем ca-certificates из builder образа
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main .
COPY --from=builder /app/services/auth-service/migrations /migrations

EXPOSE 50051

CMD ["./main"]
