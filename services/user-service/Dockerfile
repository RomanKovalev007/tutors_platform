FROM golang:1.24.3-alpine AS builder

RUN apk add --no-cache git

ENV GOPRIVATE=gitlab.crja72.ru \
    GOSUMDB=off

WORKDIR /app

# Копируем модуль api в корень проекта (чтобы ../../api работал)
COPY api/ ./api/

# Создаем структуру каталогов и копируем go.mod и go.sum user-service
COPY services/user-service/go.mod services/user-service/go.sum ./services/user-service/

# Переходим в рабочую директорию сервиса
WORKDIR /app/services/user-service

# Теперь go mod download видит ../../api/go.mod и использует локальный replace
RUN go mod download

# Копируем код самого user-service
COPY services/user-service/ ./

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
COPY --from=builder /app/services/user-service/migrations /migrations

EXPOSE 50052

CMD ["./main"]