# ====================== BUILD STAGE ======================
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

ENV GOPRIVATE=gitlab.crja72.ru \
    GOSUMDB=off

WORKDIR /app

# Добавляем ARG для инвалидации кэша при изменении кода
ARG CACHEBUST=1

# Копируем модуль api в корень проекта (чтобы ../../api работал)
COPY api/ ./api/

# Создаем структуру каталогов и копируем go.mod и go.sum task_service
COPY services/task_service/go.mod services/task_service/go.sum ./services/task_service/

# Переходим в рабочую директорию сервиса
WORKDIR /app/services/task_service

# Теперь go mod download видит ../../api/go.mod и использует локальный replace
RUN go mod download

# Копируем весь код task_service
COPY services/task_service/ ./

# Собираем бинарники — пути от корня модуля
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /app/server \
    ./cmd/server/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /app/migrate \
    ./cmd/migrate/main.go

# ====================== RUNTIME STAGE ======================
FROM alpine:latest AS runtime

# Устанавливаем только netcat и tzdata, ca-certificates копируем из builder
RUN apk add --no-cache netcat-openbsd tzdata 2>/dev/null || apk add --no-cache busybox-extras tzdata

WORKDIR /app

# Копируем ca-certificates из builder образа
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/server /app/server
COPY --from=builder /app/migrate /app/migrate
COPY --from=builder /app/services/task_service/migrations /migrations

# ====================== TARGETS ======================
FROM runtime AS base
EXPOSE 50051

FROM base AS server
ENV TASK_PORT=50051
CMD ["/app/server"]

FROM base AS migration
CMD ["/app/migrate", "--path", "/migrations", "up"]