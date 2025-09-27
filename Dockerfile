# Этап сборки
FROM golang:1.25-alpine AS builder

RUN apk add sqlite-dev gcc musl-dev build-base

WORKDIR /app

# Копируем файлы проекта
COPY go.mod ./
RUN go mod download

COPY . .

# Собираем приложение
ENV CGO_ENABLED=1
RUN go build -o /app/server ./cmd/server

# Финальный этап
FROM alpine:latest AS file-server

WORKDIR /app

# Создаем директорию для загрузок
RUN mkdir -p /app/uploads

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/server .

# Указываем порт
EXPOSE 8080

# Создаем пользователя без прав root
RUN adduser -D appuser
RUN chown -R appuser:appuser /app
USER appuser

# Запускаем сервер
CMD ["./server"]
