# Этап сборки
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем файлы проекта
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Финальный этап
FROM alpine:latest

WORKDIR /app

# Создаем директорию для загрузок
RUN mkdir -p /app/uploads

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/server .
COPY --from=builder /app/web/static /app/web/static

# Указываем порт
EXPOSE 8080

# Создаем пользователя без прав root
RUN adduser -D appuser
RUN chown -R appuser:appuser /app
USER appuser

# Запускаем сервер
CMD ["./server"]
