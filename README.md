# File Server Go

[English Version](README_en.md)

Простой файловый сервер на Go с поддержкой аутентификации, управления пользователями и приложениями для работы с файлами и директориями.

## Возможности

- 📤 Загрузка файлов через веб-интерфейс или API
- 📥 Скачивание файлов с автоматическим определением типа контента
- 📋 Просмотр списка файлов и директорий
- 📁 Создание и удаление директорий
- 🔐 JWT аутентификация для пользователей
- 🔑 API ключи для приложений (accessKeyId/accessKeySecret)
- 👥 Система управления пользователями с ролями
- 🎛️ Система управления приложениями
- 🌐 Современный веб-интерфейс с поддержкой локализации (RU/EN)
- 🔒 CORS поддержка
- 📝 Логирование запросов
- 🛡️ Безопасность: валидация путей, ограничения размера файлов

## Роли пользователей

- **Admin** - полный доступ ко всем функциям, управление пользователями и приложениями
- **Writer** - загрузка, создание, удаление файлов и директорий
- **Reader** - только просмотр и скачивание файлов

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
make dev        # Запуск с hot reload (требует air)
```

## Конфигурация

Приложение настраивается через переменные окружения:

- `PORT` - порт сервера (по умолчанию: 8080)
- `UPLOAD_DIR` - папка для загруженных файлов (по умолчанию: ./uploads)
- `JWT_KEY` - секретный ключ для JWT токенов (генерируется автоматически)
- `ADMIN_USERNAME` - имя администратора (по умолчанию: admin)
- `ADMIN_PASSWORD` - пароль администратора (по умолчанию: admin)

Пример:
```bash
export PORT=3000
export UPLOAD_DIR=/path/to/uploads
export JWT_KEY=your-secret-key
export ADMIN_USERNAME=admin
export ADMIN_PASSWORD=secure-password
make run
```

## API

### Аутентификация пользователей

Для получения JWT токена:
```bash
# Получение токена
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'

# Ответ будет содержать токен:
# {"token":"your-jwt-token"}
```

### Управление приложениями (только для админов)

#### Регистрация нового приложения
```bash
# Создание приложения с правами на чтение и запись
curl -X POST http://localhost:8080/register-app \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Application",
    "description": "Application for file management",
    "permissions": ["read", "write"]
  }'

# Ответ содержит учетные данные:
# {
#   "access_key_id": "ak_xxxxxxxxxxxxxxxx",
#   "access_key_secret": "sk_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
#   "name": "My Application",
#   "description": "Application for file management",
#   "permissions": ["read", "write"]
# }
```

#### Получение списка приложений
```bash
curl -X GET http://localhost:8080/applications \
  -H "Authorization: Bearer your-jwt-token"
```

#### Удаление приложения
```bash
curl -X DELETE http://localhost:8080/delete-app?id=1 \
  -H "Authorization: Bearer your-jwt-token"
```

### Загрузка файлов через API приложений

#### Загрузка файла с использованием accessKeyId и accessKeySecret
```bash
# Прямая загрузка файла (рекомендуемый способ)
curl -X POST http://localhost:8080/upload/raw/ \
  -H "Authorization: App ak_your_access_key_id:sk_your_access_key_secret" \
  -H "Content-Type: application/octet-stream" \
  --data-binary @path/to/your/file.txt \
  --get --data-urlencode "filename=uploaded_file.txt"

# Загрузка с TTL (файл будет удален через 3600 секунд)
curl -X POST http://localhost:8080/upload/raw/ \
  -H "Authorization: App ak_your_access_key_id:sk_your_access_key_secret" \
  -H "Content-Type: application/octet-stream" \
  --data-binary @path/to/your/file.txt \
  --get --data-urlencode "filename=temp_file.txt" \
  --data-urlencode "ttl=3600"

# Загрузка через PUT метод
curl -X PUT http://localhost:8080/upload/raw/my_document.pdf \
  -H "Authorization: App ak_your_access_key_id:sk_your_access_key_secret" \
  --data-binary @document.pdf
```

#### Примеры на разных языках программирования

**Python:**
```python
import requests

# Учетные данные приложения
access_key_id = "ak_your_access_key_id"
access_key_secret = "sk_your_access_key_secret"

# Загрузка файла
with open('file.txt', 'rb') as f:
    response = requests.post(
        'http://localhost:8080/upload/raw/',
        headers={
            'Authorization': f'App {access_key_id}:{access_key_secret}',
            'Content-Type': 'application/octet-stream'
        },
        params={'filename': 'uploaded_file.txt'},
        data=f
    )
    
print(response.json())
```

**JavaScript (Node.js):**
```javascript
const fs = require('fs');
const axios = require('axios');

const accessKeyId = 'ak_your_access_key_id';
const accessKeySecret = 'sk_your_access_key_secret';

const fileData = fs.readFileSync('file.txt');

axios.post('http://localhost:8080/upload/raw/', fileData, {
    headers: {
        'Authorization': `App ${accessKeyId}:${accessKeySecret}`,
        'Content-Type': 'application/octet-stream'
    },
    params: {
        filename: 'uploaded_file.txt'
    }
}).then(response => {
    console.log(response.data);
}).catch(error => {
    console.error(error.response.data);
});
```

**Go:**
```go
package main

import (
    "bytes"
    "fmt"
    "io"
    "net/http"
    "os"
)

func main() {
    accessKeyID := "ak_your_access_key_id"
    accessKeySecret := "sk_your_access_key_secret"
    
    file, err := os.Open("file.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    var buf bytes.Buffer
    io.Copy(&buf, file)
    
    req, err := http.NewRequest("POST", "http://localhost:8080/upload/raw/?filename=uploaded_file.txt", &buf)
    if err != nil {
        panic(err)
    }
    
    req.Header.Set("Authorization", fmt.Sprintf("App %s:%s", accessKeyID, accessKeySecret))
    req.Header.Set("Content-Type", "application/octet-stream")
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    fmt.Println("Status:", resp.Status)
}
```

### Работа с файлами через JWT токен

#### Загрузка файла через форму
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
```

#### Список файлов и директорий
```bash
# Получить список файлов в корневой директории
curl http://localhost:8080/files \
  -H "Authorization: Bearer your-jwt-token"

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
  -H "Authorization: Bearer your-jwt-token"

# Удаление файла
curl -X DELETE http://localhost:8080/delete-file?path=path/to/file.txt \
  -H "Authorization: Bearer your-jwt-token"
```

### Управление пользователями (только для админов)

#### Создание пользователя
```bash
curl -X POST http://localhost:8080/create-user \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "password": "password123",
    "role": "writer"
  }'
```

#### Получение списка пользователей
```bash
curl -X GET http://localhost:8080/users \
  -H "Authorization: Bearer your-jwt-token"
```

#### Изменение пароля
```bash
# Администратор может изменить пароль любому пользователю
curl -X POST http://localhost:8080/change-password \
  -H "Authorization: Bearer admin-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "targetuser",
    "new_password": "newpassword123"
  }'

# Обычный пользователь может изменить только свой пароль
curl -X POST http://localhost:8080/change-password \
  -H "Authorization: Bearer user-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "currentuser",
    "old_password": "oldpassword",
    "new_password": "newpassword123"
  }'
```

#### Удаление пользователя
```bash
curl -X DELETE http://localhost:8080/delete-user?id=2 \
  -H "Authorization: Bearer your-jwt-token"
```

## Структура проекта

```
file-server-go/
├── cmd/server/                    # Точка входа приложения
├── internal/                     # Приватный код приложения
│   ├── assets/                   # Встроенные ресурсы (веб-интерфейс)
│   ├── auth/                     # Аутентификация и JWT
│   ├── config/                   # Конфигурация
│   ├── database/                 # Работа с базой данных
│   ├── handlers/                 # HTTP обработчики
│   │   ├── app_handler.go        # Управление приложениями
│   │   ├── auth_handler.go       # Аутентификация
│   │   ├── file_handler.go       # Работа с файлами
│   │   └── user_handler.go       # Управление пользователями
│   ├── middleware/               # HTTP middleware
│   ├── models/                   # Модели данных
│   ├── server/                   # HTTP сервер
│   ├── services/                 # Бизнес-логика
│   └── utils/                    # Утилитные функции
├── uploads/                      # Загруженные файлы
├── Dockerfile                    # Docker конфигурация
├── Makefile                      # Make команды
└── README.md                     # Документация
```

## Безопасность

- Все запросы к API требуют аутентификации (JWT токен или API ключи)
- Поддерживается валидация путей для предотвращения path traversal атак
- Ограничение размера загружаемых файлов (по умолчанию 32MB)
- Проверка MIME-типов файлов для безопасного отображения
- Санитизация имен файлов и директорий
- Система ролей для разграничения доступа
- Безопасное хранение паролей и API ключей
- Централизованные проверки безопасности через утилитные функции

## База данных

Проект использует SQLite базу данных со следующими таблицами:
- `users` - пользователи системы
- `applications` - зарегистрированные приложения
- `files` - метаданные файлов (для TTL функциональности)

## Docker

```bash
# Сборка образа
docker build -t file-server-go .

# Запуск контейнера
docker run -p 8080:8080 -v $(pwd)/uploads:/app/uploads file-server-go
```

## Разработка

Для разработки с автоматической перезагрузкой установите [Air](https://github.com/cosmtrek/air):

```bash
go install github.com/cosmtrek/air@latest
make dev
```

## Лицензия

MIT License