# Tutors Platform - Microservices Backend

Микросервисное бэкенд-приложение для взаимодействия репетиторов и учеников. Репетиторы могут создавать группы учеников, отправлять им домашние задания, проверять их и выставлять оценки. Ученики могут состоять в нескольких группах, выполнять домашние задания и смотреть результат оценки от преподавателя.

## Архитектура

```
┌─────────────────────────────────────────────────────────────────────┐
│                          EXTERNAL CLIENT                             │
└─────────────────────────────────┬───────────────────────────────────┘
                                  │ HTTP (port 80)
                                  ▼
                        ┌──────────────────┐
                        │      NGINX       │
                        │  Reverse Proxy   │
                        └────────┬─────────┘
                                 │
                                 ▼
                        ┌──────────────────┐
                        │   API GATEWAY    │
                        │   Port: 8080     │
                        └────┬─┬─┬─────┬───┘
           ┌────────────────┘ │ │     └─────────────────┐
           │ gRPC             │ │ gRPC                  │ gRPC
           ▼                  ▼ ▼                       ▼
┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐
│  AUTH SERVICE    │ │  USER SERVICE    │ │  GROUP SERVICE   │
│  Port: 50051     │ │  Port: 50051     │ │  Port: 50051     │
│                  │ │                  │ │                  │
│ ├─ PostgreSQL    │ │ ├─ PostgreSQL    │ │ ├─ PostgreSQL    │
│ ├─ Redis         │ │ ├─ Redis Cache   │ │ ├─ Redis Cache   │
│ └─ Kafka         │ │ └─ Kafka         │ │ └─ User Client   │
└──────────────────┘ └──────────────────┘ └──────────────────┘
                                                   │
                                                   │ gRPC
                                                   ▼
                                        ┌──────────────────┐
                                        │   TASK SERVICE   │
                                        │   Port: 50051    │
                                        │                  │
                                        │ ├─ PostgreSQL    │
                                        │ ├─ Redis Cache   │
                                        │ └─ Group Client  │
                                        └──────────────────┘
```

### Микросервисы

| Сервис | Описание | Порт |
|--------|----------|------|
| **API Gateway** | HTTP REST API, маршрутизация запросов к gRPC сервисам | 8080 |
| **Auth Service** | Аутентификация, JWT токены, управление пользователями | 50051 |
| **User Service** | Профили пользователей, данные репетиторов и учеников | 50051 |
| **Group Service** | Управление группами, членство в группах | 50051 |
| **Task Service** | Задания, сдача работ, оценивание | 50051 |

### Инфраструктура

| Компонент | Описание | Порт |
|-----------|----------|------|
| **Nginx** | Reverse proxy, единая точка входа | 80 |
| **PostgreSQL** | Отдельная БД для каждого сервиса | 5432-5436 |
| **Redis Auth** | Хранение токенов для auth-service | 6379 |
| **Redis Cache** | Кэширование для всех сервисов | 6379 |
| **Kafka** | Асинхронное взаимодействие между сервисами | 9092 |
| **Zookeeper** | Координация Kafka | 2181 |
| **Prometheus** | Сбор метрик | 9090 |
| **Grafana** | Визуализация метрик | 3000 |

## Технологии

- **Go 1.24** - основной язык разработки
- **gRPC** - межсервисное взаимодействие
- **gRPC-Gateway** - HTTP REST API поверх gRPC
- **PostgreSQL 15** - хранение данных
- **Redis 7** - кэширование и хранение токенов
- **Kafka** - асинхронная коммуникация
- **Docker & Docker Compose** - контейнеризация
- **Nginx** - reverse proxy
- **Prometheus & Grafana** - мониторинг

## Запуск

### Предварительные требования

- Docker и Docker Compose
- Make (опционально)

### Запуск всех сервисов

```bash
docker-compose up -d
```

### Проверка статуса

```bash
docker-compose ps
```

### Остановка

```bash
docker-compose down
```

### Просмотр логов

```bash
docker-compose logs -f [service-name]
```

## API Endpoints

### Аутентификация (Auth Service)

| Метод | Endpoint | Описание |
|-------|----------|----------|
| POST | `/v1/auth/register` | Регистрация пользователя |
| POST | `/v1/auth/login` | Вход в систему |
| POST | `/v1/auth/logout` | Выход из системы |
| POST | `/v1/auth/refresh` | Обновление токенов |
| PATCH | `/v1/auth/change-password` | Смена пароля |
| POST | `/v1/auth/forgot-password` | Запрос на сброс пароля |
| POST | `/v1/auth/reset-password` | Сброс пароля |
| GET | `/v1/auth/users` | Список пользователей |
| GET | `/v1/auth/users/{id}` | Получение пользователя |
| PATCH | `/v1/auth/users/{id}` | Обновление пользователя |
| DELETE | `/v1/auth/users/{id}` | Удаление пользователя |

### Профили пользователей (User Service)

| Метод | Endpoint | Описание |
|-------|----------|----------|
| POST | `/v1/users` | Создание профиля |
| GET | `/v1/users/{user_id}` | Получение профиля |
| PATCH | `/v1/users/{user_id}` | Обновление профиля |
| DELETE | `/v1/users/{user_id}` | Удаление профиля |
| GET | `/v1/users/{user_id}/types` | Получение типа пользователя |
| GET | `/v1/users/{user_id}/full` | Полный профиль |
| POST | `/v1/tutors` | Создание профиля репетитора |
| GET | `/v1/tutors/{user_id}` | Получение профиля репетитора |
| PATCH | `/v1/tutors/{user_id}` | Обновление профиля репетитора |
| DELETE | `/v1/tutors/{user_id}` | Удаление профиля репетитора |
| POST | `/v1/students` | Создание профиля ученика |
| GET | `/v1/students/{user_id}` | Получение профиля ученика |
| PATCH | `/v1/students/{user_id}` | Обновление профиля ученика |
| DELETE | `/v1/students/{user_id}` | Удаление профиля ученика |

### Группы (Group Service)

| Метод | Endpoint | Описание |
|-------|----------|----------|
| POST | `/v1/groups` | Создание группы |
| GET | `/v1/groups` | Список групп |
| GET | `/v1/groups/{id}` | Получение группы |
| PATCH | `/v1/groups/{id}` | Обновление группы |
| DELETE | `/v1/groups/{id}` | Удаление группы |
| GET | `/v1/groups/{group_id}/members` | Участники группы |
| POST | `/v1/groups/{group_id}/members` | Добавление участников |
| POST | `/v1/groups/{group_id}/members:remove` | Удаление участников |

### Задания (Task Service)

| Метод | Endpoint | Описание |
|-------|----------|----------|
| POST | `/v1/tasks` | Создание задания |
| GET | `/v1/tasks/{task_id}` | Получение задания |
| PATCH | `/v1/tasks/{task_id}` | Обновление задания |
| DELETE | `/v1/tasks/{task_id}` | Удаление задания |
| GET | `/v1/groups/{group_id}/tasks` | Задания группы |
| GET | `/v1/tasks/created-by-me` | Мои задания |
| POST | `/v1/tasks/{task_id}/submissions` | Сдача работы |
| GET | `/v1/tasks/{task_id}/submissions` | Работы по заданию |
| GET | `/v1/submissions/{submission_id}` | Получение работы |
| PATCH | `/v1/submissions/{submission_id}` | Обновление работы |
| DELETE | `/v1/submissions/{submission_id}` | Удаление работы |
| POST | `/v1/submissions/{submission_id}/grade` | Оценивание работы |
| POST | `/v1/submissions/{submission_id}/reset-grade` | Сброс оценки |

### Мониторинг

| Endpoint | Описание |
|----------|----------|
| `/health` | Health check API Gateway |
| `/metrics` | Prometheus метрики |
| `/prometheus/` | Prometheus UI |
| `/grafana/` | Grafana дашборды |

## Отказоустойчивость

### Паттерны

- **Retry** - автоматические повторные попытки при временных ошибках
- **Circuit Breaker** - защита от каскадных отказов
- **Fallback** - резервные стратегии при недоступности сервисов
- **Timeout** - ограничение времени ожидания

### Конфигурация Circuit Breaker

```go
CircuitBreakerConfig{
    FailureThreshold: 5,      // Порог ошибок для открытия
    SuccessThreshold: 2,      // Порог успехов для закрытия
    Timeout:          30s,    // Время до перехода в half-open
}
```

### Конфигурация Retry

```go
RetryConfig{
    MaxAttempts:     3,             // Максимум попыток
    InitialInterval: 100ms,         // Начальный интервал
    MaxInterval:     2s,            // Максимальный интервал
    Multiplier:      2.0,           // Множитель интервала
    Jitter:          0.1,           // Джиттер (10%)
}
```

## Кэширование

### Redis использование

| Сервис | Redis | Назначение |
|--------|-------|------------|
| Auth Service | redis-auth | Хранение refresh токенов |
| User Service | redis-cache | Кэширование профилей |
| Group Service | redis-cache | Кэширование групп |
| Task Service | redis-cache | Кэширование заданий |

### TTL по умолчанию

- Refresh токены: 14 дней
- Кэш профилей: 1 час
- Кэш групп: 30 минут

## Мониторинг

### Метрики Prometheus

```
# HTTP метрики
http_requests_total{method, path, status}
http_request_duration_seconds{method, path}
http_requests_in_flight

# gRPC метрики
grpc_client_requests_total{service, method, code}
grpc_client_request_duration_seconds{service, method}

# Circuit Breaker
circuit_breaker_state{service, name}

# Кэш
cache_hits_total{service, cache_type}
cache_misses_total{service, cache_type}
```

### Grafana Dashboard

Доступ: `http://localhost/grafana/`

Логин: `admin` / `admin`

Дашборд включает:
- Request Rate
- Latency (p95)
- Error Rate by Service
- Circuit Breaker State
- Cache Hit Rate

## Тестирование

### Запуск тестов

```bash
# Auth Service
cd services/auth-service && go test -v ./...

# User Service
cd services/user-service && go test -v ./...

# Task Service
cd services/task_service && go test -v ./...

# Group Service
cd services/group-service && go test -v ./...
```

### Проверка покрытия

```bash
cd services/auth-service
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

### Требования к покрытию

- Минимальное покрытие: **30%** для каждого сервиса (кроме API Gateway)
- CI/CD pipeline автоматически проверяет покрытие

## CI/CD

### GitLab CI/CD Pipeline

Стадии:
1. **build** - сборка Go приложений
2. **test** - запуск тестов
3. **coverage** - проверка покрытия (>= 30%)
4. **docker** - сборка и публикация Docker образов

### Триггеры

- Изменения в соответствующих директориях сервисов
- Merge requests
- Push в main branch

## Структура проекта

```
.
├── api/                          # API определения (proto)
│   ├── proto/                    # Protocol Buffer файлы
│   ├── gen/go/                   # Сгенерированный Go код
│   └── Makefile                  # Генерация proto
│
├── services/
│   ├── api-gateway/              # HTTP REST → gRPC gateway
│   ├── auth-service/             # Аутентификация
│   ├── user-service/             # Профили пользователей
│   ├── group-service/            # Группы
│   └── task_service/             # Задания
│
├── nginx/                        # Nginx конфигурация
├── prometheus/                   # Prometheus конфигурация
├── grafana/                      # Grafana конфигурация
│
├── docker-compose.yaml           # Основной compose файл
├── docker-compose.services.yaml  # Сервисы приложения
├── docker-compose.infrastructure.yaml  # Инфраструктура
│
├── .gitlab-ci.yml               # CI/CD pipeline
└── README.md                    # Документация
```

## Переменные окружения

### API Gateway

```env
HTTP_PORT=8080
HTTP_READ_TIMEOUT=30s
HTTP_WRITE_TIMEOUT=30s
AUTH_GRPC=auth-go:50051
USER_GRPC=user-go:50051
GROUP_GRPC=group-go:50051
TASKS_GRPC=tasks-go:50051
```

### Auth Service

```env
GRPC_PORT=50051
POSTGRES_HOST=postgres-auth
POSTGRES_PORT=5432
POSTGRES_DB=auth_db
REDIS_HOST=redis-auth
ACCESS_TTL_M=5
REFRESH_TTL_H=336
SECRET=your-secret-key
KAFKA_BROKERS=kafka:9092
```

### User Service

```env
GRPC_PORT=50051
POSTGRES_HOST=postgres-user
POSTGRES_PORT=5432
POSTGRES_DB=user_db
REDIS_CACHE_HOST=redis-cache
KAFKA_BROKERS=kafka:9092
```

### Group Service

```env
GRPC_PORT=50051
USER_SERVICE_ADDRESS=user-go:50051
POSTGRES_HOST=postgres-group
POSTGRES_PORT=5432
POSTGRES_DB=group-db
REDIS_CACHE_HOST=redis-cache
```

### Task Service

```env
TASK_PORT=50051
GROUP_SERV_ADDR=group-go:50051
POSTGRES_HOST=postgres-task
POSTGRES_PORT=5432
POSTGRES_DB=task_postgres
REDIS_CACHE_HOST=redis-cache
```

## Безопасность

- JWT токены для аутентификации (access + refresh)
- Access токен: 5 минут TTL
- Refresh токен: 14 дней TTL
- Все сервисы изолированы в Docker network
- Наружу открыт только порт Nginx (80)
- Rate limiting на уровне Nginx
- Security headers (X-Frame-Options, X-Content-Type-Options, etc.)

## Контакты

При возникновении вопросов или проблем создайте issue в репозитории.
