.PHONY: build run clean test fmt vet

# Переменные
BINARY_NAME=file-server
MAIN_PATH=./cmd/server

# Сборка приложения
build:
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

# Запуск приложения
run:
	go run $(MAIN_PATH)

# Очистка
clean:
	go clean
	rm -f bin/$(BINARY_NAME)

# Тестирование
test:
	go test -v ./...

# Форматирование кода
fmt:
	go fmt ./...

# Проверка кода
vet:
	go vet ./...

# Установка зависимостей
deps:
	go mod download
	go mod tidy

# Сборка для разных платформ
build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux $(MAIN_PATH)

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows.exe $(MAIN_PATH)

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-mac $(MAIN_PATH)

# Сборка для всех платформ
build-all: build-linux build-windows build-mac

# Запуск с hot reload (требует air: go install github.com/cosmtrek/air@latest)
dev:
	air