# Этап 1: Сборка бинарника
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код проекта
COPY . .

# Собираем легковесный бинарник для Linux без CGO
RUN CGO_ENABLED=0 GOOS=linux go build -o /iplatform-api ./cmd/api/main.go

# Этап 2: Финальный минимальный образ для запуска
FROM alpine:3.19

WORKDIR /

# Копируем скомпилированный файл из предыдущего этапа
COPY --from=builder /iplatform-api /iplatform-api

# Открываем порт наружу
EXPOSE 8080

# Точка входа в контейнер
ENTRYPOINT ["/iplatform-api"]
