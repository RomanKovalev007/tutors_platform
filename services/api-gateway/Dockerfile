FROM golang:1.24.3-alpine AS builder

RUN apk add --no-cache git

ENV GOPRIVATE=gitlab.crja72.ru \
    GOSUMDB=off

WORKDIR /app

# 1. Копируем go.mod и go.sum api-gateway
COPY services/api-gateway/go.mod services/api-gateway/go.sum ./

# 2. КОПИРУЕМ ВЕСЬ МОДУЛЬ api ДО go mod download (чтобы replace работал)
COPY api/ ./api/

# 3. Теперь go mod download видит api/go.mod и использует локальный replace
RUN --mount=type=secret,id=netrc,target=/root/.netrc \
    go mod download

# 4. Копируем код самого api-gateway
COPY services/api-gateway/ .

# 5. Билдим (main в cmd/app)
RUN go build -o main ./cmd/app

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]