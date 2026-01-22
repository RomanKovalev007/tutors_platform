# ====================== BUILD STAGE ======================
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

ENV GOPRIVATE=gitlab.crja72.ru \
    GOSUMDB=off

WORKDIR /app

# Добавляем ARG для инвалидации кэша при изменении кода
ARG CACHEBUST=1

# Копируем go.mod и go.sum в корень
COPY services/task_service/go.mod services/task_service/go.sum ./

RUN go mod download

# Копируем весь код task_service в корень /app (как было раньше)
COPY services/task_service/ ./

# Копируем общий модуль api, если он есть
COPY api ./api

# Выполняем go mod tidy
RUN go mod tidy

# Копируем миграции
COPY services/task_service/migrations /migrations

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

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /app/migrate /app/migrate
COPY --from=builder /migrations /migrations

# ====================== TARGETS ======================
FROM runtime AS base
EXPOSE 50051

FROM base AS server
ENV TASK_PORT=50051
CMD ["/app/server"]

FROM base AS migration
CMD ["/app/migrate", "--path", "/migrations", "up"]