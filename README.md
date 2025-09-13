# File Server Go

Простой файловый сервер на Go с поддержкой аутентификации для управления файлами и директориями.

## Возможности

- 📤 Загрузка файлов через веб-интерфейс или API
- 📥 Скачивание файлов
- 📋 Просмотр списка файлов и директорий
- 📁 Создание и удаление директорий
- 🔐 JWT аутентификация
- 🌐 Современный веб-интерфейс
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
- `JWT_SECRET` - секретный ключ для JWT токенов
- `AUTH_USERNAME` - имя пользователя для аутентификации
- `AUTH_PASSWORD` - пароль для аутентификации

Пример:
```bash
export PORT=3000
export UPLOAD_DIR=/path/to/uploads
export JWT_SECRET=your-secret-key
export AUTH_USERNAME=admin
export AUTH_PASSWORD=password
make run
```

## API

### Аутентификация

Для получения JWT токена:
```bash
# Получение токена
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'

# Ответ будет содержать токен:
# {"token":"your-jwt-token"}
```

### Работа с файлами

#### Загрузка файла через curl
```bash
# Загрузка файла в корневую директорию
curl -X POST http://localhost:8080/upload \
  -H "Authorization: Bearer your-jwt-token" \
  -F "file=@path/to/your/file.txt"

# Загрузка файла в определенную директорию
curl -X POST http://localhost:8080/upload \
  -H "Authorization: Bearer your-jwt-token" \
  -F "file=@path/to/your/file.txt" \
  -F "path=directory/subdirectory"

# Прямая загрузка файла (raw upload)
curl -X PUT http://localhost:8080/upload/raw/filename.txt \
  -H "Authorization: Bearer your-jwt-token" \
  --data-binary @path/to/your/file.txt
```

#### Загрузка файла через wget
```bash
# Получение токена и сохранение в переменную
TOKEN=$(wget -qO- --post-data='{"username":"admin","password":"password"}' \
  --header='Content-Type: application/json' \
  http://localhost:8080/login | grep -o '"token":"[^"]*' | cut -d'"' -f4)

# Загрузка файла
wget --method=POST \
  --header="Authorization: Bearer $TOKEN" \
  --body-file=/path/to/your/file.txt \
  http://localhost:8080/upload/raw/filename.txt
```

#### Список файлов и директорий
```bash
# Получить список файлов в корневой директории
curl http://localhost:8080/files \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json"

# Получить список файлов в определенной директории
curl -X POST http://localhost:8080/files \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{"path": "directory/subdirectory"}'
```

#### Скачивание файла
```bash
# Скачивание файла
curl -O -H "Authorization: Bearer your-jwt-token" \
  http://localhost:8080/download/path/to/file.txt
```

#### Работа с директориями
```bash
# Создание директории
curl -X POST http://localhost:8080/create-dir \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "dirname": "new_directory",
    "path": "parent_directory"
  }'

# Удаление директории
curl -X DELETE http://localhost:8080/delete-dir?path=directory/to/delete \
  -H "Authorization: Bearer your-jwt-token" \

# Удаление файла
curl -X DELETE http://localhost:8080/delete-file?path=path/to/file.txt\
  -H "Authorization: Bearer your-jwt-token" \
```

## Структура проекта

```
file-server-go/
├── cmd/server/          # Точка входа приложения
├── internal/           # Приватный код приложения
│   ├── auth/          # Аутентификация и JWT
│   ├── config/        # Конфигурация
│   ├── handlers/      # HTTP обработчики
│   ├── middleware/    # HTTP middleware
│   └── server/        # HTTP сервер
├── web/              # Веб-интерфейс
│   └── templates/    # HTML шаблоны
├── uploads/          # Загруженные файлы
├── Makefile         # Make команды
└── README.md        # Документация
```

## Безопасность

- Все запросы к API требуют JWT аутентификации
- Поддерживается валидация путей для предотвращения path traversal атак
- Ограничение размера загружаемых файлов (по умолчанию 32MB)
- Проверка MIME-типов файлов
- Санитизация имен файлов и директорий

## Разработка

Для разработки с автоматической перезагрузкой установите [Air](https://github.com/cosmtrek/air):

```bash
go install github.com/cosmtrek/air@latest
make dev
```

## Лицензия

MIT License