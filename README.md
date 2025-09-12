# File Server Go

Простой файловый сервер на Go для загрузки, скачивания и управления файлами.

## Возможности

- 📤 Загрузка файлов через веб-интерфейс
- 📥 Скачивание файлов
- 📋 Просмотр списка загруженных файлов
- 🌐 Простой веб-интерфейс
- 🔒 CORS поддержка
- 📝 Логирование запросов

## Быстрый старт

### Установка и запуск

```bash
# Клонируйте репозиторий
git clone <your-repo-url>
cd file-server-go

# Установите зависимости
make deps

# Запустите сервер
make run
```

Сервер будет доступен по адресу: http://localhost:8080

### Использование Make команд

```bash
make build      # Собрать приложение
make run        # Запустить приложение
make test       # Запустить тесты
make fmt        # Форматировать код
make vet        # Проверить код
make clean      # Очистить сборку
make build-all  # Собрать для всех платформ
```

## Конфигурация

Приложение настраивается через переменные окружения:

- `PORT` - порт сервера (по умолчанию: 8080)
- `UPLOAD_DIR` - папка для загруженных файлов (по умолчанию: ./uploads)

Пример:
```bash
export PORT=3000
export UPLOAD_DIR=/path/to/uploads
make run
```

## API

### Загрузка файла
```bash
curl -X POST -F "file=@example.txt" http://localhost:8080/upload
```

### Список файлов
```bash
curl http://localhost:8080/files
```

### Скачивание файла
```bash
curl -O http://localhost:8080/download/example.txt
```

## Структура проекта

```
file-server-go/
├── cmd/server/          # Точка входа приложения
├── internal/            # Приватный код приложения
│   ├── config/         # Конфигурация
│   ├── handlers/       # HTTP обработчики
│   ├── middleware/     # HTTP middleware
│   └── server/         # HTTP сервер
├── uploads/            # Загруженные файлы
├── Makefile           # Make команды
└── README.md          # Документация
```

## Разработка

Для разработки с автоматической перезагрузкой установите [Air](https://github.com/cosmtrek/air):

```bash
go install github.com/cosmtrek/air@latest
make dev
```

## Лицензия

MIT License