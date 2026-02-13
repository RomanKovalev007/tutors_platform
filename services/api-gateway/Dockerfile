FROM golang:1.24.3-alpine AS builder

RUN apk add --no-cache git

ENV GOPRIVATE=gitlab.crja72.ru \
    GOSUMDB=off

WORKDIR /app

# Копируем модуль api в корень проекта (чтобы ../../api работал)
COPY api/ ./api/

# Создаем структуру каталогов и копируем go.mod и go.sum api-gateway
COPY services/api-gateway/go.mod services/api-gateway/go.sum ./services/api-gateway/

# Переходим в рабочую директорию сервиса
WORKDIR /app/services/api-gateway

# Теперь go mod download видит ../../api/go.mod и использует локальный replace
RUN go mod download

# Копируем код самого api-gateway
COPY services/api-gateway/ ./

# Билдим (main в cmd/app)
RUN go build -o /app/main ./cmd/app

# Финальный образ
FROM alpine:latest

WORKDIR /app

# Копируем ca-certificates из builder образа
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]